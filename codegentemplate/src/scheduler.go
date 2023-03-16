package src

import (
	"github.com/robfig/cron/v3"

	"grest.dev/cmd/codegentemplate/app"
)

func Scheduler() *schedulerImpl {
	if scheduler == nil {
		scheduler = &schedulerImpl{}
		if app.APP_ENV == "local" || app.IS_MAIN_SERVER {
			scheduler.Configure()
		}
		scheduler.isConfigured = true
	}
	return scheduler
}

var scheduler *schedulerImpl

type schedulerImpl struct {
	isConfigured bool
}

func (s *schedulerImpl) Configure() {
	c := cron.New()

	// add scheduler func here, for example :
	// c.AddFunc("CRON_TZ=Asia/Jakarta 5 0 * * *", app.Auth().RemoveExpiredToken)

	c.Start()
}
