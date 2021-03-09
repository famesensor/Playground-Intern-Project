package validators

import (
	"log"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var validatorEntity *validator.Validate = validator.New()

var translator = en.New()
var uni = ut.New(translator, translator)
var trans ut.Translator

type validateStruct = map[string]interface{}

func InitializeTranslator() {
	transS, found := uni.GetTranslator("en")

	if !found {
		log.Fatal("Translator not found")
	}

	trans = transS

	if err := en_translations.RegisterDefaultTranslations(validatorEntity, trans); err != nil {
		log.Fatal(err)
	}

	_ = validatorEntity.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	_ = validatorEntity.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} must be a valid email", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})

	_ = validatorEntity.RegisterTranslation("required_if", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} must be a valid field", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})
}

func ValidateStructErrors(v interface{}) (bool, validateStruct) {
	err := validatorEntity.Struct(v)
	if err != nil {
		validateData := make(validateStruct)
		for _, err := range err.(validator.ValidationErrors) {
			validateData[err.Field()] = err.Translate(trans)
		}

		return false, validateData
	}
	return true, nil
}
