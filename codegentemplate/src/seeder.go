package src

import (
	"log/slog"

	"grest.dev/cmd/codegentemplate/app"
)

func Seeder() *seederUtil {
	if seeder == nil {
		seeder = &seederUtil{}
		seeder.Configure()
		if app.APP_ENV == "local" || app.IS_MAIN_SERVER {
			seeder.Run()
		}
		seeder.isConfigured = true
	}
	return seeder
}

var seeder *seederUtil

type seederUtil struct {
	isConfigured bool
}

func (s *seederUtil) Configure() {
	// example
	// app.DB().RegisterSeeder("main", "2024-10-09_16.30-country-data", country.Seeder().Run)
}

func (s *seederUtil) Run() {
	tx, err := app.DB().Conn("main")
	if err != nil {
		app.Logger().Error("Failed to connect to main db", slog.Any("err", err))
	}
	if err = app.DB().RunSeeder(tx, "main", app.Setting{}); err != nil {
		app.Logger().Error("Failed to connect to run seeder", slog.Any("err", err))
	}
}
