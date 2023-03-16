package app

import (
	enLocale "github.com/go-playground/locales/en"
	idLocale "github.com/go-playground/locales/id"
	enTranslation "github.com/go-playground/validator/v10/translations/en"
	idTranslation "github.com/go-playground/validator/v10/translations/id"
	"grest.dev/grest"
)

func Validator() ValidatorInterface {
	if validator == nil {
		validator = &validatorImpl{}
		validator.configure()
	}
	return validator
}

type ValidatorInterface interface {
	grest.ValidatorInterface
}

var validator *validatorImpl

// validatorImpl implement ValidatorInterface embed from grest.Validator for simplicity
type validatorImpl struct {
	grest.Validator
}

func (v *validatorImpl) configure() {
	v.New()
	v.RegisterCustomTypeFunc(v.ValidateValuer,
		NullBool{},
		NullInt64{},
		NullFloat64{},
		NullString{},
		NullDateTime{},
		NullDate{},
		NullTime{},
		NullText{},
		NullJSON{},
		NullUUID{},
	)
	v.RegisterTranslator("en", enLocale.New(), enTranslation.RegisterDefaultTranslations)
	v.RegisterTranslator("id", idLocale.New(), idTranslation.RegisterDefaultTranslations)
	v.RegisterTranslator("id-ID", idLocale.New(), idTranslation.RegisterDefaultTranslations)
}
