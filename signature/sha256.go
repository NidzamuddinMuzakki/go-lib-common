package signature

import (
	"context"
	sha256Lib "crypto/sha256"
)

type sha256 struct {
}

func newSha256() *sha256 {
	return &sha256{}
}

func (s *sha256) Generate(_ context.Context, key string) (string, error) {
	h := sha256Lib.New()
	_, err := h.Write([]byte(key))
	if err != nil {
		return "", err
	}

	hashed := h.Sum(nil)

	return string(hashed), nil
}

func (s *sha256) Verify(_ context.Context, key, sign string) bool {
	h := sha256Lib.New()
	h.Write([]byte(key))

	hashed := h.Sum(nil)

	if string(hashed) == sign {
		return true
	}

	return false
}
