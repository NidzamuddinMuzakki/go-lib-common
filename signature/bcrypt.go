package signature

import (
	"context"
	sha256Lib "crypto/sha256"

	bcryptLib "golang.org/x/crypto/bcrypt"
)

type bcrypt struct {
	cost int
}

func newBCrypt(cost int) *bcrypt {
	return &bcrypt{
		cost: cost,
	}
}

func (b *bcrypt) Generate(_ context.Context, key string) (string, error) {
	h := sha256Lib.New()
	h.Write([]byte(key))

	hashed := h.Sum(nil)

	encrypted, err := bcryptLib.GenerateFromPassword(hashed, b.cost)
	if err != nil {
		return "", err
	}

	return string(encrypted), nil
}

func (b *bcrypt) Verify(_ context.Context, key string, sign string) bool {
	h := sha256Lib.New()
	h.Write([]byte(key))

	hashed := h.Sum(nil)

	err := bcryptLib.CompareHashAndPassword([]byte(sign), hashed)

	return err == nil
}
