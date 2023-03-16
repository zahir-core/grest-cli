package src

import (
	"grest.dev/cmd/codegentemplate/app"
	// codegentemplate import : DONT REMOVE THIS COMMENT
)

func Router() *routerImpl {
	if router == nil {
		router = &routerImpl{}
		router.Configure()
		router.isConfigured = true
	}
	return router
}

var router *routerImpl

type routerImpl struct {
	isConfigured bool
}

func (r *routerImpl) Configure() {
	app.Server().AddRoute("/api/version", "GET", app.VersionHandler, nil)

	// codegentemplate AddRoute : DONT REMOVE THIS COMMENT
}
