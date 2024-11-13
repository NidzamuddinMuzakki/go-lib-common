package validator_test

import (
	"reflect"
	"strings"
	"testing"

	"bitbucket.org/moladinTech/go-lib-common/response/model"
	"bitbucket.org/moladinTech/go-lib-common/validator"
	"github.com/stretchr/testify/require"
)

type Person struct {
	Name    string `validate:"required_without=Address"`
	Address string `validate:"required"`
	Weight  string `validate:"ltecsfield=Height"`
	Url     string `validate:"url"`
	Height  int    `validate:"max=100"`
	Gender  uint   `validate:"required_if=Url url Height 400"`
	Role    string `validate:"required_unless=Url uri"`
}

func TestNewValidator_ShouldSucceedToError(t *testing.T) {
	t.Parallel()

	t.Run("Should Succeed New Validator", func(t *testing.T) {
		vld := validator.New()
		err := vld.Struct(Person{Weight: "500", Url: "url", Height: 400})
		require.Error(t, err)

		errText := validator.ToErrResponse(err)
		require.Equal(t,
			strings.Join([]string{
				"Name is a required if Address is empty",
				"Address is a required field",
				"Weight is less than to another Height field",
				"Url must be a valid URL",
				"Height must be a maximum of 100 in length",
				"Gender is a required if Url is url and Height is 400",
				"Role is a required if Url is not uri",
			}, ","),
			errText,
		)
	})
}

func TestNewValidatorV2_ShouldSucceedToError(t *testing.T) {
	t.Parallel()

	t.Run("Should Succeed New Validator V2", func(t *testing.T) {
		vld := validator.New()
		err := vld.Struct(Person{Weight: "500", Url: "url", Height: 400})
		require.Error(t, err)

		errText := validator.ToErrResponseV2(err)
		errMap := []model.ValidationResponse{
			{
				Field:   "Name",
				Message: "Name is a required if Address is empty",
			},
			{
				Field:   "Address",
				Message: "Address is a required field",
			},
			{
				Field:   "Weight",
				Message: "Weight is less than to another Height field",
			},
			{
				Field:   "Url",
				Message: "Url must be a valid URL",
			},
			{
				Field:   "Height",
				Message: "Height must be a maximum of 100 in length",
			},
			{
				Field:   "Gender",
				Message: "Gender is a required if Url is url and Height is 400",
			},
			{
				Field:   "Role",
				Message: "Role is a required if Url is not uri",
			},
		}
		if !reflect.DeepEqual(errText, errMap) {
			t.Errorf("UpdateValues() failed, expected %v but got %v", errMap, errText)
		}
	})
}
