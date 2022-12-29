package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

func New() *validator.Validate {
	validate := validator.New()
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")

	enTranslations.RegisterDefaultTranslations(validate, trans)

	return validate
}

func ToErrResponse(err error) string {
	var errors []string
	if fieldErrors, ok := err.(validator.ValidationErrors); ok {
		for _, err := range fieldErrors {
			switch err.Tag() {
			case "required":
				errors = append(errors, fmt.Sprintf("%s is a required field", err.Field()))
			case "len":
				errors = append(errors, fmt.Sprintf("%s must be a %s length", err.Field(), err.Param()))
			case "min":
				errors = append(errors, fmt.Sprintf("%s must be a minimum of %s in length", err.Field(), err.Param()))
			case "max":
				errors = append(errors, fmt.Sprintf("%s must be a maximum of %s in length", err.Field(), err.Param()))
			case "url":
				errors = append(errors, fmt.Sprintf("%s must be a valid URL", err.Field()))
			case "oneof":
				errors = append(errors, fmt.Sprintf("%s must be an oneof [%s]", err.Field(), err.Param()))
			case "required_without":
				errors = append(errors, fmt.Sprintf("%s is a required if %s is empty", err.Field(), err.Param()))
			case "required_without_all":
				errors = append(errors, fmt.Sprintf("%s is a required if %s are empty", err.Field(), err.Param()))
			case "required_with":
				errors = append(errors, fmt.Sprintf("%s is a required if %s is not empty", err.Field(), err.Param()))
			case "excluded_with":
				errors = append(errors, fmt.Sprintf("%s is a exclude if %s is empty", err.Field(), err.Param()))
			case "ltecsfield":
				errors = append(errors, fmt.Sprintf("%s is less than to another %s field", err.Field(), err.Param()))
			default:
				errors = append(errors, fmt.Sprintf("something wrong on %s; %s", err.Field(), err.Tag()))
			}
		}
	}
	return strings.Join(errors, ",")
}
