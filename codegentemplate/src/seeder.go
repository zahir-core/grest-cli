package src

import "grest.dev/cmd/codegentemplate/app"

func Seeder() *seederImpl {
	if seeder == nil {
		seeder = &seederImpl{}
		seeder.Configure()
		if app.APP_ENV == "local" || app.IS_MAIN_SERVER {
			seeder.Run()
		}
		seeder.isConfigured = true
	}
	return seeder
}

var seeder *seederImpl

type seederImpl struct {
	isConfigured bool
}

func (s *seederImpl) Configure() {

}

func (s *seederImpl) Run() {

}
