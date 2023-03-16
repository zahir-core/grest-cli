package app

import (
	"grest.dev/grest"

	"grest.dev/cmd/codegentemplate/app/i18n"
)

func Translator() TranslatorInterface {
	if translator == nil {
		translator = &translatorImpl{}
		translator.configure()
	}
	return translator
}

type TranslatorInterface interface {
	grest.TranslatorInterface
}

var translator *translatorImpl

// translatorImpl implement translatorInterface embed from grest.translator for simplicity
type translatorImpl struct {
	grest.Translator
}

func (t *translatorImpl) configure() {
	t.AddTranslation("en-US", i18n.EnUS())
	t.AddTranslation("id-ID", i18n.IdID())
}
