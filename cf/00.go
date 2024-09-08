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

var GinGroup *gin.RouterGroup
var ORMDB *gorm.DB
var Trans ut.Translator
var Config *CfgRegist

func Init(
	ginGroup *gin.RouterGroup,
	ormDB *gorm.DB,
	config *CfgRegist,
) {
	GinGroup = ginGroup
	ORMDB = ormDB
	Config = config

	RegisterValidatorTranslations()
}

func RegisterValidatorTranslations() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zh := zh.New()
		uni := ut.New(zh, zh)
		Trans, _ = uni.GetTranslator("zh")
		zh_translations.RegisterDefaultTranslations(v, Trans)
	}
}
