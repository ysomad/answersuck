package validation

import (
	"errors"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type ginTranslator struct {
	validate *validator.Validate
	trans    ut.Translator
}

func NewGinTranslator() (*ginTranslator, error) {
	eng := en.New()
	uni := ut.New(eng, eng)
	trans, found := uni.GetTranslator("en")
	if !found {
		return nil, errors.New("validation translator not found")
	}

	v := binding.Validator.Engine().(*validator.Validate)

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	return &ginTranslator{
		validate: v,
		trans:    trans,
	}, nil
}

func (gt *ginTranslator) register() error {
	return enTranslations.RegisterDefaultTranslations(gt.validate, gt.trans)
}

// TranslateError returns translated validation errors received from gin.c.ShouldBindJSON err
func (gt *ginTranslator) TranslateError(err error) map[string]string {
	_ = gt.register()

	errs := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {
		errs[err.Field()] = strings.Join(strings.Split(err.Translate(gt.trans), " "), " ")
	}

	return errs
}
