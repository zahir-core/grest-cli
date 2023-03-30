package app

import (
	"mime/multipart"

	"grest.dev/grest"
)

func Telegram(message ...string) TelegramInterface {
	if telegram == nil {
		telegram = &telegramUtil{}
		telegram.configure()
		if len(message) > 0 {
			telegram.AddMessage(message[0])
		}
	}
	return telegram
}

type TelegramInterface interface {
	AddMessage(text string)
	AddAttachment(file *multipart.FileHeader)
	Send() error
}

var telegram *telegramUtil

// telegramUtil implement TelegramInterface embed from grest.Telegram for simplicity
type telegramUtil struct {
	grest.Telegram
}

func (t *telegramUtil) configure() {
	t.BotToken = TELEGRAM_ALERT_TOKEN
	t.ChatID = TELEGRAM_ALERT_USER_ID
}
