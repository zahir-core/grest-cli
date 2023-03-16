package app

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"grest.dev/grest"
)

func NewError(statusCode int, message string, detail ...any) *grest.Error {
	return grest.NewError(statusCode, message, detail...)
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	lang := "en"
	ctx, ctxOK := c.Locals("ctx").(*Ctx)
	if ctxOK {
		lang = ctx.Lang
	}
	e, ok := err.(*grest.Error)
	if !ok {
		code := http.StatusInternalServerError
		fiberError, isFiberError := err.(*fiber.Error)
		if isFiberError {
			code = fiberError.Code
		}
		e = NewError(code, err.Error())
	}
	if e.StatusCode() < 400 || e.StatusCode() > 599 {
		e.Code = http.StatusInternalServerError
	}
	if e.StatusCode() == http.StatusInternalServerError {
		// todo: add trace to log & send alert to telegram
		e.Detail = map[string]string{"message": e.Error()}
		e.Message = Translator().Trans(lang, "500_internal_error")
	}
	return c.Status(e.StatusCode()).JSON(e.Body())
}

func NotFoundHandler(c *fiber.Ctx) error {
	lang := c.Get("Accept-Language")
	if lang == "" {
		lang = c.Get("Content-Language") // backward compatibility
		if lang == "" {
			lang = "en"
		}
	}
	e := NewError(http.StatusNotFound, Translator().Trans(lang, "404_not_found"))
	return c.Status(e.StatusCode()).JSON(e.Body())
}

func Recover(c *fiber.Ctx) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// todo: save log & send alert to telegram
		}
	}()
	return c.Next()
}
