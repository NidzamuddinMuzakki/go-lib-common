package encryption_test

import (
	"fmt"
	"testing"

	"bitbucket.org/moladinTech/go-lib-common/encryption"
	"bitbucket.org/moladinTech/go-lib-common/validator"

	"encoding/hex"
	"github.com/stretchr/testify/require"
)

func TestGenerateSalt_ShouldSucceed(t *testing.T) {
	t.Run("Should Succeed Generate Salt", func(t *testing.T) {
		encPkg := encryption.NewEncryption(
			validator.New(),
			encryption.WithAppName("test-service"),
		)
		salt := encPkg.GenerateSalt("secret-key-test")
		require.Equal(t, []byte{0x8, 0x3e, 0x2d, 0xbb, 0x1e, 0x54, 0xc9, 0x13, 0x84, 0xdd, 0x3a, 0x2d, 0x7d, 0xd3, 0x7c, 0x8e}, salt)
	})
}

func TestNewEncryption_ErrorOnValidation(t *testing.T) {
	t.Run("Error On Validation New Encryption", func(t *testing.T) {
		require.Panics(t, func() {
			encryption.NewEncryption(
				validator.New(),
			)
		})
	})
}

func TestEncrypt_ShouldSucceed(t *testing.T) {
	t.Run("Should Succeed Encryption", func(t *testing.T) {
		encPkg := encryption.NewEncryption(
			validator.New(),
			encryption.WithAppName("test-service"),
		)
		salt := encPkg.GenerateSalt("secret-key-test")
		enc, _ := encPkg.Encrypt("data encryption", salt)
		encToStr := hex.EncodeToString(enc)
		require.Equal(t, "7665ae375d061ccd2ec9c9e280cf51ab", encToStr)
	})
}

func TestDecrypt_ShouldSucceed(t *testing.T) {
	t.Run("Should Succeed Decryption", func(t *testing.T) {
		encPkg := encryption.NewEncryption(
			validator.New(),
			encryption.WithAppName("test-service"),
		)
		salt := encPkg.GenerateSalt("secret-key-test")
		dec, _ := encPkg.Decrypt("7665ae375d061ccd2ec9c9e280cf51ab", salt)
		decToStr := fmt.Sprintf("%s", dec)
		require.Equal(t, "data encryption", decToStr)
	})
}
