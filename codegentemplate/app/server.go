package app

import (
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"grest.dev/grest"
)

func Server() *serverUtil {
	if server == nil {
		server = &serverUtil{}
		server.configure()
	}
	return server
}

var server *serverUtil

type serverUtil struct {
	Addr                  string
	IsUseTLS              bool
	CertFile              string
	KeyFile               string
	DisableStartupMessage bool
	Fiber                 *fiber.App
}

func (s *serverUtil) configure() {
	s.Addr = ":" + APP_PORT
	s.Fiber = fiber.New(fiber.Config{
		ErrorHandler:          s.Error,
		ReadBufferSize:        16384,
		DisableStartupMessage: true,
	})
	s.AddMiddleware(s.Recover)
}

// use grest to add route so it can generate swagger api documentation automatically
func (s *serverUtil) AddRoute(path, method string, handler fiber.Handler, operation OpenAPIOperationInterface) {
	if method == "ALL" {
		for _, m := range []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE", "CONNECT", "OPTIONS", "TRACE"} {
			s.AddRoute(path, m, handler, operation)
		}
	} else {
		s.Fiber.Add(method, strings.ReplaceAll(strings.ReplaceAll(path, "{", ":"), "}", ""), handler)
		if IS_GENERATE_OPEN_API_DOC && operation != nil {
			OpenAPI().AddRoute(path, method, operation)
		}
	}
}

func (s *serverUtil) AddStaticRoute(path string, fsConfig filesystem.Config) {
	s.Fiber.Use(path, filesystem.New(fsConfig))
}

func (s *serverUtil) AddStaticFSRoute(path, dirName string, f embed.FS) {
	dirFS, err := fs.Sub(f, dirName)
	if err != nil {
		Logger().Error("Failed to add "+dirName, slog.Any("err", err))
	}
	s.AddStaticRoute(path, filesystem.Config{
		Root: http.FS(dirFS),
	})
}

func (s *serverUtil) AddMiddleware(handler fiber.Handler) {
	s.Fiber.Use(handler)
}

func (serverUtil) Version(c *fiber.Ctx) error {
	return c.JSON(map[string]any{
		"version": APP_VERSION,
	})
}

func (s *serverUtil) NotFound(c *fiber.Ctx) error {
	lang := c.Get("Accept-Language")
	if lang == "" || lang == "*" || strings.Contains(lang, ",") || strings.Contains(lang, ";") {
		lang = "en"
	}
	err := Error().New(http.StatusNotFound, Translator().Trans(lang, "404_not_found"))
	return c.Status(err.StatusCode()).JSON(err.Body())
}

// Error handles errors by processing them and returning an appropriate response.
// It retrieves the language from the context (c) and assigns it to lang.
// It checks if the error is an instance of grest.Error.
// If it is not, it sets the error code and message based on the received error.
// If the error status code is not in the 4xx or 5xx range, it sets the code to http.StatusInternalServerError.
// If the error status code is http.StatusInternalServerError, it translates the error message and assigns it to e.Message.
// It returns a JSON response with the error status code and body.
func (serverUtil) Error(c *fiber.Ctx, err error) error {
	lang := "en"
	ctx, ctxOK := c.Locals("ctx").(*Ctx)
	if ctxOK {
		lang = ctx.Lang
		ctx.Err = err
	}
	e, ok := err.(*grest.Error)
	if !ok {
		e = &grest.Error{}
		code := http.StatusInternalServerError
		fiberError, isFiberError := err.(*fiber.Error)
		if isFiberError {
			code = fiberError.Code
		}
		e.Code = code
		e.Message = err.Error()
	}
	if e.StatusCode() < 400 || e.StatusCode() > 599 {
		e.Code = http.StatusInternalServerError
	}
	if e.StatusCode() == http.StatusInternalServerError {
		if e.Detail == nil {
			e.Detail = map[string]string{"message": e.Error()}
		}
		e.Message = Translator().Trans(lang, "500_internal_error")
	}
	return c.Status(e.StatusCode()).JSON(e.Body())
}

// Recover recovers from a panic during Fiber request processing.
// It uses a defer statement to catch and recover from panics.
// Inside the deferred function, there is a placeholder for saving logs and sending alerts.
func (serverUtil) Recover(c *fiber.Ctx) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
			ctx := &Ctx{}
			if c != nil {
				ctx, _ = c.Locals(CtxKey).(*Ctx)
			}
			attrs := Logger().Attrs(*ctx)
			attrs = append(attrs, slog.Any("panic", r))
			attrs = append(attrs, slog.Any("stack", string(debug.Stack())))
			Logger().Error(err.Error(), attrs...)
		}
	}()
	return c.Next()
}

// RecoverAsync catch and recover from panics on asyncronous processing (goroutine).
// It saving logs and sending alerts.
func (serverUtil) RecoverAsync(ctx Ctx, msg string) {
	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("%v", r)
		}
		attrs := Logger().Attrs(ctx)
		attrs = append(attrs, slog.Any("panic", r))
		attrs = append(attrs, slog.Any("stack", string(debug.Stack())))
		Logger().Error(msg+" : "+err.Error(), attrs...)
	}
}

func (s *serverUtil) Test(req *http.Request, msTimeout ...int) (*http.Response, error) {
	return s.Fiber.Test(req, msTimeout...)
}

func (s *serverUtil) Start() error {
	s.Fiber.Use(s.NotFound)
	if !s.DisableStartupMessage {
		grest.StartupMessage(s.Addr)
	}
	if s.IsUseTLS {
		return s.Fiber.ListenTLS(s.Addr, s.CertFile, s.KeyFile)
	}
	return s.Fiber.Listen(s.Addr)
}
