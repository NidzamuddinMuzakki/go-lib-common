//go:generate mockery --name=IEncryption
package encryption

import (
	"crypto/md5"
	"fmt"

	"github.com/andreburgaud/crypt2go/ecb"
	"github.com/andreburgaud/crypt2go/padding"
	"golang.org/x/crypto/blowfish"

	"bitbucket.org/moladinTech/go-lib-common/sentry"
	commonValidator "bitbucket.org/moladinTech/go-lib-common/validator"
	"github.com/go-playground/validator/v10"
)

type IEncryption interface {
	Encrypt(data string) []byte
}

type EncryptionPackage struct {
	AppName string         `validate:"required"`
	Sentry  sentry.ISentry `validate:"required"`
}

func WithAppName(appName string) Option {
	return func(s *EncryptionPackage) {
		s.AppName = appName
	}
}
func WithSentry(sentry sentry.ISentry) Option {
	return func(s *EncryptionPackage) {
		s.Sentry = sentry
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

func (p *EncryptionPackage) Encrypt(data string) []byte {
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s:%s", p.AppName, p.AppName)))
	key := h.Sum(nil)
	bytes := []byte(data)
	block, err := blowfish.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	mode := ecb.NewECBEncrypter(block)
	padder := padding.NewPkcs5Padding()
	bytes, err = padder.Pad(bytes)
	if err != nil {
		panic(err.Error())
	}
	ct := make([]byte, len(bytes))
	mode.CryptBlocks(ct, bytes)
	return ct
}
