package validation

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

var (
	errTranslatorNotFound = errors.New("translator not found")
)

type module struct {
	validate *validator.Validate
	trans    ut.Translator
}

func NewModule() (*module, error) {
	v := newValidate()

	t, err := newTranslator()
	if err != nil {
		return nil, fmt.Errorf("newTranslator: %w", err)
	}

	if err = enTranslations.RegisterDefaultTranslations(v, t); err != nil {
		return nil, fmt.Errorf("enTranslations.RegisterDefaultTranslations: %w", err)
	}

	return &module{
		validate: v,
		trans:    t,
	}, nil
}

func newTranslator() (ut.Translator, error) {
	eng := en.New()
	uni := ut.New(eng, eng)

	t, found := uni.GetTranslator("en")
	if !found {
		return nil, errTranslatorNotFound
	}

	return t, nil
}

func newValidate() *validator.Validate {
	v := validator.New()

	v.RegisterTagNameFunc(tagName)

	return v
}

func tagName(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

	if name == "-" {
		return ""
	}

	return name
}

func (v *module) ValidateStruct(s any) error {
	return v.validate.Struct(s)
}

func (v *module) TranslateError(err error) map[string]string {
	errs := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {
		errs[err.Field()] = err.Translate(v.trans)
	}

	return errs
}
