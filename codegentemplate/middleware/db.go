package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"grest.dev/cmd/codegentemplate/app"
)

func SetDB(c *fiber.Ctx) error {
	ctx, ok := c.Locals(app.CtxKey).(*app.Ctx)
	if !ok {
		return app.NewError(http.StatusInternalServerError, "ctx is not found")
	}
	ctx.TxBegin()
	err := c.Next()
	if err != nil || (c.Response().StatusCode() >= http.StatusBadRequest || c.Response().StatusCode() < http.StatusOK) {
		ctx.TxRollback()
	} else {
		ctx.TxCommit()
	}
	return nil
}
