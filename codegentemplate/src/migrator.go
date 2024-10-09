package src

import (
	"log/slog"

	"grest.dev/cmd/codegentemplate/app"
	// import : DONT REMOVE THIS COMMENT
)

func Migrator() *migratorUtil {
	if migrator == nil {
		migrator = &migratorUtil{}
		migrator.Configure()
		if app.APP_ENV == "local" || app.IS_MAIN_SERVER {
			migrator.Run()
		}
		migrator.isConfigured = true
	}
	return migrator
}

var migrator *migratorUtil

type migratorUtil struct {
	isConfigured bool
}

func (*migratorUtil) Configure() {
	// RegisterTable : DONT REMOVE THIS COMMENT
}

func (*migratorUtil) Run() {
	tx, err := app.DB().Conn("main")
	if err != nil {
		app.Logger().Error("Failed to connect to main db", slog.Any("err", err))
	}
	if err = app.DB().MigrateTable(tx, "main", app.Setting{}); err != nil {
		app.Logger().Error("Failed to connect to migrate db table", slog.Any("err", err))
	}
}
