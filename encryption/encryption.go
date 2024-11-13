//go:generate mockery --name=IEncryption
package encryption

import (
	"fmt"

	commonValidator "bitbucket.org/moladinTech/go-lib-common/validator"

	"crypto/md5"
	"encoding/hex"
	"github.com/andreburgaud/crypt2go/ecb"
	"github.com/andreburgaud/crypt2go/padding"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/blowfish"
)

type IEncryption interface {
	GenerateSalt(key string) []byte
	Encrypt(data string, salt []byte) ([]byte, error)
	Decrypt(data string, salt []byte) ([]byte, error)
}

type EncryptionPackage struct {
	AppName string `validate:"required"`
}

func WithAppName(appName string) Option {
	return func(s *EncryptionPackage) {
		s.AppName = appName
	}
}

type Option func(*EncryptionPackage)

func NewEncryption(
	validator *validator.Validate,
	options ...Option,
) *EncryptionPackage {
	encryptionPackage := &EncryptionPackage{}

	for _, option := range options {
		option(encryptionPackage)
	}

	err := validator.Struct(encryptionPackage)
	if err != nil {
		panic(commonValidator.ToErrResponse(err))
	}

	return encryptionPackage
}

// GenerateSalt Convert the key to a md5 hash
func (p *EncryptionPackage) GenerateSalt(key string) []byte {
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s:%s", p.AppName, key)))
	return h.Sum(nil)
}

// Encrypt the data using the salt
func (p *EncryptionPackage) Encrypt(data string, salt []byte) ([]byte, error) {
	// Create a new blowfish cipher
	bytes := []byte(data)
	block, err := blowfish.NewCipher(salt)
	if err != nil {
		return nil, err
	}

	// Pad the data
	mode := ecb.NewECBEncrypter(block)
	padder := padding.NewPkcs5Padding()
	bytes, err = padder.Pad(bytes)
	if err != nil {
		return nil, err
	}

	// Encrypt the data
	ct := make([]byte, len(bytes))
	mode.CryptBlocks(ct, bytes)
	return ct, nil
}

func (p *EncryptionPackage) Decrypt(data string, salt []byte) ([]byte, error) {
	ciphertext, err := hex.DecodeString(data)
	if err != nil {
		return nil, err
	}

	// Create a new blowfish cipher
	block, err := blowfish.NewCipher(salt)
	if err != nil {
		return nil, err
	}

	// Decrypt the data
	mode := ecb.NewECBDecrypter(block)
	pt := make([]byte, len(ciphertext))
	mode.CryptBlocks(pt, ciphertext)

	// Unpad the data
	padder := padding.NewPkcs5Padding()
	pt, err = padder.Unpad(pt)
	if err != nil {
		return nil, err
	}

	return pt, nil
}
