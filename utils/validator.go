package utils

import (
	"github.com/noxue/validator"
	zh_translations "github.com/noxue/validator/translations/zh"
	"github.com/noxue/universal-translator"
	"github.com/noxue/locales/zh"
	"github.com/gin-gonic/gin"
	"log"
)

var Validate *validator.Validate
var Trans ut.Translator

func initValidator() {
	Validate = validator.New()

	zh1 := zh.New()
	uni := ut.New(zh1, zh1)
	Trans, _ = uni.GetTranslator("zh")

	zh_translations.RegisterDefaultTranslations(Validate, Trans)
}

func ValidateStruct(t interface{}) (errs []gin.H){
	err := Validate.Struct(t)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			log.Println(err)
			return
		}
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs,gin.H{err.Field():err.Translate(Trans)})
		}

		return
	}
	return
}