package rt

import (
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	zh := zh.New()
	uni = ut.New(zh, zh)

	trans, _ = uni.GetTranslator("zh")

	validate = validator.New(validator.WithRequiredStructEnabled())
	zh_translations.RegisterDefaultTranslations(validate, trans)
}

func ValidateStruct(data interface{}) validator.ValidationErrorsTranslations {
	err := validate.Struct(data)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		return errs.Translate(trans)
	}
	return nil
}
