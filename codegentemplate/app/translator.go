package app

import (
	"grest.dev/grest"

	"grest.dev/cmd/codegentemplate/app/i18n"
)

func Translator() TranslatorInterface {
	if translator == nil {
		translator = &translatorUtil{}
		translator.configure()
	}
	return translator
}

type TranslatorInterface interface {
	Trans(lang, key string, params ...map[string]string) string
}

var translator *translatorUtil

// translatorUtil implement translatorInterface embed from grest.translator for simplicity
type translatorUtil struct {
	grest.Translator
}

func (t *translatorUtil) configure() {
	t.AddTranslation("en-US", i18n.EnUS())
	t.AddTranslation("id-ID", i18n.IdID())
}
