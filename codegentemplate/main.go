package main

import (
	"embed"
	"os"

	"grest.dev/cmd/codegentemplate/app"
	"grest.dev/cmd/codegentemplate/src"
)

//go:embed all:docs
var f embed.FS

func main() {
	app.Config()
	if app.IS_GENERATE_OPEN_API_DOC {
		src.Router()
		app.OpenAPI().Configure().Generate()
		os.Exit(0)
	}

	app.Logger()
	app.Cache()
	app.Validator()
	app.Translator()
	app.FS()
	app.DB()
	defer app.DB().Close()
	app.Server()

	if app.APP_ENV != "production" {
		app.Server().AddOpenAPIDoc("/api/docs", f)
	}

	src.Middleware()
	src.Router()
	src.Migrator()
	src.Seeder()
	src.Scheduler()
	err := app.Server().Start()
	if err != nil {
		app.Logger().Fatal().Err(err).Send()
	}
}
