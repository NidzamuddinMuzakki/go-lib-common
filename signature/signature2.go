package signature

import (
	"context"
	"errors"

	bcryptLib "golang.org/x/crypto/bcrypt"
)

type Generator2 interface {
	Generate2(ctx context.Context, key string) (string, error)
}

type Verifier2 interface {
	Verify2(ctx context.Context, key, sign string) (bool, error)
}

type GenerateAndVerify2 interface {
	Generator2
	Verifier2
}

type Signature2 struct {
	Algorithm     *algorithm
	Cost          int
	ExpiredinHour int
}

type Option2 func(*Signature2)

func WithAlgorithm2(algorithm algorithm) Option2 {
	return func(s *Signature2) {
		s.Algorithm = &algorithm
	}
}

func WithCost2(cost int) Option2 {
	return func(s *Signature2) {
		s.Cost = cost
	}
}

func WithExpired(expired int) Option2 {
	return func(s *Signature2) {
		s.ExpiredinHour = expired
	}
}

var (
	ErrHashAlgorithmUnavailable2 = errors.New("signature: algorithm unavailable")
	ErrSignatureInvalid          = errors.New("signature: Invalid")
	ErrSignatureExpired          = errors.New("signature: Expired")
)

func NewSignature2(options ...Option2) (GenerateAndVerify, error) {
	s := Signature2{}
	for _, option := range options {
		option(&s)
	}

	if s.Algorithm == nil {
		return nil, ErrHashAlgorithmUnavailable
	}
	if s.ExpiredinHour == 0 {
		s.ExpiredinHour = 1
	}
	switch *s.Algorithm {
	case Sha256:
		return newSha256(), nil
	case BCrypt:
		if s.Cost == 0 {
			s.Cost = bcryptLib.DefaultCost
		}
		return newBCrypt(s.Cost), nil
	default:
		return newSha256(), nil
	}
}
