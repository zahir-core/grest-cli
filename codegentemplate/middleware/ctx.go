package middleware

import (
	"github.com/gofiber/fiber/v2"

	"grest.dev/cmd/codegentemplate/app"
)

func NewCtx(c *fiber.Ctx) error {
	lang := c.Get("Accept-Language")
	if lang == "" {
		lang = "en"
	}
	ctx := app.Ctx{
		Lang: lang,
	}
	c.Locals("ctx", &ctx)
	return c.Next()
}
