package app

import (
	"embed"
	"io/fs"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"grest.dev/grest"
)

func Server() ServerInterface {
	if server == nil {
		server = &serverUtil{}
		server.configure()
	}
	return server
}

type ServerInterface interface {
	AddRoute(path, method string, handler fiber.Handler, operation OpenAPIOperationInterface)
	AddStaticRoute(path string, fsConfig filesystem.Config)
	AddOpenAPIDoc(path string, f embed.FS)
	AddMiddleware(handler fiber.Handler)
	Start() error
	Test(req *http.Request, msTimeout ...int) (resp *http.Response, err error)
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
		ErrorHandler:          ErrorHandler,
		ReadBufferSize:        16384,
		DisableStartupMessage: true,
	})
	s.AddMiddleware(Recover)
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

func (s *serverUtil) AddOpenAPIDoc(path string, f embed.FS) {
	docs, err := fs.Sub(f, "docs")
	if err != nil {
		Logger().Fatal().Err(err).Send()
	}
	s.AddStaticRoute(path, filesystem.Config{
		Root: http.FS(docs),
	})
}

func (s *serverUtil) AddMiddleware(handler fiber.Handler) {
	s.Fiber.Use(handler)
}

func (s *serverUtil) Start() error {
	s.Fiber.Use(NotFoundHandler)
	if !s.DisableStartupMessage {
		grest.StartupMessage(s.Addr)
	}
	if s.IsUseTLS {
		return s.Fiber.ListenTLS(s.Addr, s.CertFile, s.KeyFile)
	}
	return s.Fiber.Listen(s.Addr)
}

func (s *serverUtil) Test(req *http.Request, msTimeout ...int) (*http.Response, error) {
	return s.Fiber.Test(req, msTimeout...)
}

func ParseQuery(c *fiber.Ctx) url.Values {
	query := url.Values{}
	_, qs, _ := strings.Cut(c.OriginalURL(), "?")
	for qs != "" {
		q := ""
		q, qs, _ = strings.Cut(qs, "&")
		if q == "" || strings.Contains(q, ";") {
			continue
		}

		key, value, _ := strings.Cut(q, "=")
		if k, err := url.QueryUnescape(key); err == nil {
			key = k
		}
		if v, err := url.QueryUnescape(value); err == nil {
			value = v
		}
		query.Add(key, value)
	}
	return query
}

func VersionHandler(c *fiber.Ctx) error {
	return c.JSON(map[string]any{
		"version": APP_VERSION,
	})
}
