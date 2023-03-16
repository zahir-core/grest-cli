package src

import (
	"grest.dev/cmd/codegentemplate/app"
	"grest.dev/cmd/codegentemplate/middleware"
)

func Middleware() *middlewareImpl {
	if mdlwr == nil {
		mdlwr = &middlewareImpl{}
		mdlwr.Configure()
		mdlwr.isConfigured = true
	}
	return mdlwr
}

var mdlwr *middlewareImpl

type middlewareImpl struct {
	isConfigured bool
}

func (*middlewareImpl) Configure() {
	app.Server().AddMiddleware(middleware.NewCtx)
	app.Server().AddMiddleware(middleware.SetDB)
}
