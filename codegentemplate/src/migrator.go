package src

import (
	"grest.dev/cmd/codegentemplate/app"
	// codegentemplate import : DONT REMOVE THIS COMMENT
)

func Migrator() *migratorImpl {
	if migrator == nil {
		migrator = &migratorImpl{}
		migrator.Configure()
		if app.APP_ENV == "local" || app.IS_MAIN_SERVER {
			migrator.Run()
		}
		migrator.isConfigured = true
	}
	return migrator
}

var migrator *migratorImpl

type migratorImpl struct {
	isConfigured bool
}

func (*migratorImpl) Configure() {
	// codegentemplate RegisterTable : DONT REMOVE THIS COMMENT
}

func (*migratorImpl) Run() {
	tx, err := app.DB().Conn("main")
	if err != nil {
		app.Logger().Fatal().Err(err).Send()
	} else {
		err = app.DB().MigrateTable(tx, "main", app.Setting{})
	}
	if err != nil {
		app.Logger().Fatal().Err(err).Send()
	}
}
