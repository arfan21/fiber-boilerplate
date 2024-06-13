package validation

import (
	"encoding/json"
	"reflect"
	"strings"
	"sync"

	"github.com/arfan21/fiber-boilerplate/pkg/constant"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	uni        *ut.UniversalTranslator
	validate   *validator.Validate
	translator ut.Translator

	uniOnce        sync.Once
	validateOnce   sync.Once
	translatorOnce sync.Once
)

func Validate[T any](modelValidate T) error {
	uniOnce.Do(func() {
		en := en.New()
		uni = ut.New(en, en)
	})

	validateOnce.Do(func() {
		validate = validator.New()
	})

	translatorOnce.Do(func() {
		translatorUni, _ := uni.GetTranslator("en")
		translator = translatorUni
		en_translations.RegisterDefaultTranslations(validate, translator)
	})

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	err := validate.Struct(modelValidate)
	if err != nil {
		var messages []map[string]interface{}

		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.Field()

			messages = append(messages, map[string]interface{}{
				"field":   fieldName,
				"message": err.Translate(translator),
			})
		}
		jsonMessage, errJson := json.Marshal(messages)
		if errJson != nil {
			return errJson
		}

		return &constant.ErrValidation{Message: string(jsonMessage)}
	}

	return nil
}
