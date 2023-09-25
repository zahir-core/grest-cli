package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	"grest.dev/cmd/codegentemplate/app"
)

func Log() *logHandler {
	if lh == nil {
		lh = &logHandler{}
	}
	return lh
}

var lh *logHandler

type logHandler struct{}

func (l *logHandler) New(c *fiber.Ctx) error {
	startAt := time.Now()
	err := c.Next()
	level := slog.LevelInfo
	if err != nil || c.Response().StatusCode() >= http.StatusInternalServerError || c.Response().StatusCode() < http.StatusOK {
		level = slog.LevelError
	} else if c.Response().StatusCode() >= http.StatusBadRequest {
		level = slog.LevelWarn
	}
	ctx, attrs := l.getAttrs(c, startAt)
	go l.send(ctx, level, attrs)
	return nil
}

func (*logHandler) getAttrs(c *fiber.Ctx, startAt time.Time) (app.Ctx, []any) {
	var attrs []any
	if app.LOG_WITH_DURATION {
		finishAt := time.Now()
		attrs = append(attrs, slog.Time("start_at", startAt))
		attrs = append(attrs, slog.Time("finish_at", finishAt))
		attrs = append(attrs, slog.Duration("duration", finishAt.Sub(startAt)))
	}
	attrs = append(attrs, slog.Int("status", c.Response().StatusCode()))
	if app.LOG_WITH_REQUEST_HEADER {
		headers := []any{}
		c.Request().Header.VisitAll(func(key, value []byte) {
			headers = append(headers, slog.String(string(key), string(value)))
		})
		attrs = append(attrs, slog.Group("header", headers...))
	}
	if app.LOG_WITH_REQUEST_BODY {
		var bodyRequest any
		c.BodyParser(&bodyRequest)
		attrs = append(attrs, slog.Any("body_request", bodyRequest))
	}
	if app.LOG_WITH_RESPONSE_BODY {
		var bodyResponse any
		json.Unmarshal(c.Response().Body(), &bodyResponse)
		attrs = append(attrs, slog.Any("body_response", bodyResponse))
	}
	ctx, _ := c.Locals(app.CtxKey).(*app.Ctx)
	return *ctx, attrs
}

func (*logHandler) send(ctx app.Ctx, level slog.Level, attrs []any) error {
	msg := ""
	if ctx.Err != nil {
		msg = app.Error().GetOriginalMessage(ctx.Err)
	}
	attrs = app.Logger().Attrs(ctx, attrs)
	if level == slog.LevelError {
		app.Logger().Error(msg, attrs...)
	} else if level == slog.LevelWarn {
		app.Logger().Warn(msg, attrs...)
	} else {
		app.Logger().Info(msg, attrs...)
	}
	return nil
}
