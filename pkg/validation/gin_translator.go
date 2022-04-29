package validation

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type GinTranslator struct {
	validate *validator.Validate
	trans    ut.Translator
}

func NewGinTranslator() (*GinTranslator, error) {
	eng := en.New()
	uni := ut.New(eng, eng)
	trans, found := uni.GetTranslator("en")
	if !found {
		return nil, errors.New("validation translator not found")
	}

	return &GinTranslator{
		validate: binding.Validator.Engine().(*validator.Validate),
		trans:    trans,
	}, nil
}

func (gt *GinTranslator) register() error {
	return enTranslations.RegisterDefaultTranslations(gt.validate, gt.trans)
}

// TranslateError returns translated validation errors received from gin.c.ShouldBindJSON err
func (gt *GinTranslator) TranslateError(err error) map[string]string {
	_ = gt.register()

	errs := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {
		//	errs[err.Field()] = err.Translate(gt.trans)
		errs[err.Field()] = strings.Join(strings.Split(err.Translate(gt.trans), " ")[1:], " ")
	}

	return errs
}
