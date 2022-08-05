package validate

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

var errTranslatorNotFound = errors.New("translator not found")

type validate struct {
	v     *validator.Validate
	trans ut.Translator
}

func New() (*validate, error) {
	eng := en.New()
	uni := ut.New(eng, eng)

	t, found := uni.GetTranslator("en")
	if !found {
		return nil, errTranslatorNotFound
	}

	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	if err := enTranslations.RegisterDefaultTranslations(v, t); err != nil {
		return nil, fmt.Errorf("enTranslations.RegisterDefaultTranslations: %w", err)
	}

	return &validate{v: v, trans: t}, nil
}

func (v *validate) TranslateError(err error) map[string]string {
	errs := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		errs[err.Field()] = err.Translate(v.trans)
	}
	return errs
}

func (v *validate) RequestBody(b io.ReadCloser, dest any) error {
	if err := json.NewDecoder(b).Decode(dest); err != nil {
		return err
	}

	if err := v.v.Struct(dest); err != nil {
		return err
	}

	return nil
}
