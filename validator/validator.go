package validator

import (
	"fmt"
	"strings"

	modelResp "bitbucket.org/moladinTech/go-lib-common/response/model"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

var ErrValidator = map[string]string{}

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
			case "required_if":
				params := strings.Split(err.Param(), " ")
				formattedParams := params[0]
				for i, param := range params {
					if i > 0 {
						if i%2 != 0 {
							formattedParams += fmt.Sprintf(" is %s", param)
						} else {
							formattedParams += fmt.Sprintf(" and %s", param)
						}
					}
				}
				errors = append(errors, fmt.Sprintf("%s is a required if %s", err.Field(), formattedParams))
			case "required_unless":
				paramString := err.Param()
				formattedParams := strings.Replace(paramString, " ", " is not ", -1)
				errors = append(errors, fmt.Sprintf("%s is a required if %s", err.Field(), formattedParams))
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

func ToErrResponseV2(err error) (validationResponses []modelResp.ValidationResponse) {

	if fieldErrors, ok := err.(validator.ValidationErrors); ok {
		for _, err := range fieldErrors {
			switch err.Tag() {
			case "required":
				validationResponses = append(validationResponses, modelResp.ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s is a required field", err.Field()),
				})
			case "len":
				validationResponses = append(validationResponses, modelResp.ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s must be a %s length", err.Field(), err.Param()),
				})
			case "min":
				validationResponses = append(validationResponses, modelResp.ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s must be a minimum of %s in length", err.Field(), err.Param()),
				})
			case "max":
				validationResponses = append(validationResponses, modelResp.ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s must be a maximum of %s in length", err.Field(), err.Param()),
				})
			case "url":
				validationResponses = append(validationResponses, modelResp.ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s must be a valid URL", err.Field()),
				})
			case "oneof":
				validationResponses = append(validationResponses, modelResp.ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s must be an oneof [%s]", err.Field(), err.Param()),
				})
			case "required_if":
				params := strings.Split(err.Param(), " ")
				formattedParams := params[0]
				for i, param := range params {
					if i > 0 {
						if i%2 != 0 {
							formattedParams += fmt.Sprintf(" is %s", param)
						} else {
							formattedParams += fmt.Sprintf(" and %s", param)
						}
					}
				}
				validationResponses = append(validationResponses, modelResp.ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s is a required if %s", err.Field(), formattedParams),
				})
			case "required_unless":
				paramString := err.Param()
				formattedParams := strings.Replace(paramString, " ", " is not ", -1)
				validationResponses = append(validationResponses, modelResp.ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s is a required if %s", err.Field(), formattedParams),
				})
			case "required_without":
				validationResponses = append(validationResponses, modelResp.ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s is a required if %s is empty", err.Field(), err.Param()),
				})
			case "required_without_all":
				validationResponses = append(validationResponses, modelResp.ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s is a required if %s are empty", err.Field(), err.Param()),
				})
			case "required_with":
				validationResponses = append(validationResponses, modelResp.ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s is a required if %s is not empty", err.Field(), err.Param()),
				})
			case "excluded_with":
				validationResponses = append(validationResponses, modelResp.ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s is a exclude if %s is empty", err.Field(), err.Param()),
				})
			case "ltecsfield":
				validationResponses = append(validationResponses, modelResp.ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s is less than to another %s field", err.Field(), err.Param()),
				})
			default:
				errValidator, ok := ErrValidator[err.Tag()]
				if ok {
					count := strings.Count(errValidator, "%s")
					if count == 1 {
						validationResponses = append(validationResponses, modelResp.ValidationResponse{
							Field:   err.Field(),
							Message: fmt.Sprintf(errValidator, err.Field()),
						})
					} else {
						validationResponses = append(validationResponses, modelResp.ValidationResponse{
							Field:   err.Field(),
							Message: fmt.Sprintf(errValidator, err.Field(), err.Param()),
						})
					}
				} else {
					validationResponses = append(validationResponses, modelResp.ValidationResponse{
						Field:   err.Field(),
						Message: fmt.Sprintf("something wrong on %s; %s", err.Field(), err.Tag()),
					})
				}
			}
		}
	}
	return validationResponses
}
