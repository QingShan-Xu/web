// Package rt 包含了验证器的初始化和使用方法。
package rt

import (
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"

	"github.com/go-playground/validator/v10"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	// 初始化中文翻译。
	zhLocale := zh.New()
	uni = ut.New(zhLocale, zhLocale)

	trans, _ = uni.GetTranslator("zh")

	validate = validator.New()
	_ = zhTranslations.RegisterDefaultTranslations(validate, trans)
}

// ValidateStruct 验证结构体并返回错误信息。
// data: 需要验证的结构体。
// 返回验证错误的翻译信息。
func ValidateStruct(data interface{}) validator.ValidationErrorsTranslations {
	if err := validate.Struct(data); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			return errs.Translate(trans)
		}
	}
	return nil
}
