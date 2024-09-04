package cf

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"gorm.io/gorm"
)

var TokenJWT string
var GinGroup *gin.Engine
var ORMDB *gorm.DB
var Trans ut.Translator

func Init(ginGroup *gin.Engine, ormDB *gorm.DB, JwtSecret string) {
	TokenJWT = JwtSecret
	GinGroup = ginGroup
	ORMDB = ormDB
	setupValidatorTranslations()
}

func setupValidatorTranslations() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zh := zh.New()
		uni := ut.New(zh, zh)
		Trans, _ = uni.GetTranslator("zh")
		zh_translations.RegisterDefaultTranslations(v, Trans)
	}
}
