package app

import "grest.dev/grest"

func Telegram() TelegramInterface {
	if telegram == nil {
		telegram = &telegramImpl{}
		telegram.configure()
	}
	return telegram
}

type TelegramInterface interface {
	grest.TelegramInterface
}

var telegram *telegramImpl

// telegramImpl implement TelegramInterface embed from grest.Telegram for simplicity
type telegramImpl struct {
	grest.Telegram
}

func (t *telegramImpl) configure() {
	t.BotToken = TELEGRAM_ALERT_TOKEN
	t.ChatID = TELEGRAM_ALERT_USER_ID
}
