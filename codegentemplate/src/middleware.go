package src

import (
	"grest.dev/cmd/codegentemplate/app"
	"grest.dev/cmd/codegentemplate/middleware"
)

func Middleware() *middlewareUtil {
	if mdlwr == nil {
		mdlwr = &middlewareUtil{}
		mdlwr.Configure()
		mdlwr.isConfigured = true
	}
	return mdlwr
}

var mdlwr *middlewareUtil

type middlewareUtil struct {
	isConfigured bool
}

func (*middlewareUtil) Configure() {
	app.Server().AddMiddleware(middleware.Ctx().New)
	app.Server().AddMiddleware(middleware.DB().New)
}
