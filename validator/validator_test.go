package validator_test

import (
	"strings"
	"testing"

	"bitbucket.org/moladinTech/go-lib-common/validator"
	"github.com/stretchr/testify/require"
)

type Person struct {
	Name    string `validate:"required_without=Address"`
	Address string `validate:"required"`
	Weight  string `validate:"ltecsfield=Height"`
	Url     string `validate:"url"`
	Height  int    `validate:"max=100"`
}

func TestNewValidator_ShouldSucceedToError(t *testing.T) {
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
			}, ","),
			errText,
		)
	})
}
