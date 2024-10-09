package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"grest.dev/cmd/codegentemplate/app"
)

func Ctx() *ctxHandler {
	if ch == nil {
		ch = &ctxHandler{}
	}
	return ch
}

var ch *ctxHandler

type ctxHandler struct{}

func (*ctxHandler) New(c *fiber.Ctx) error {
	action := app.Action{
		Method: c.Method(),
		Path:   c.Path(),
	}
	lang := c.Get("Accept-Language")
	if lang == "" || lang == "*" || strings.Contains(lang, ",") || strings.Contains(lang, ";") {
		lang = "en"
	}
	ctx := app.Ctx{
		Lang:   lang,
		Action: action,
	}
	c.Locals("ctx", &ctx)
	return c.Next()
}
