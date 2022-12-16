package signature

import (
	"context"
	"errors"
	bcryptLib "golang.org/x/crypto/bcrypt"
)

type Generator interface {
	Generate(ctx context.Context, key string) (string, error)
}

type Verifier interface {
	Verify(ctx context.Context, key, sign string) bool
}

type GenerateAndVerify interface {
	Generator
	Verifier
}

type algorithm int

const (
	Sha256 algorithm = iota
	BCrypt
)

type Signature struct {
	Algorithm *algorithm
	Cost      int
}

type Option func(*Signature)

func WithAlgorithm(algorithm algorithm) Option {
	return func(s *Signature) {
		s.Algorithm = &algorithm
	}
}

func WithCost(cost int) Option {
	return func(s *Signature) {
		s.Cost = cost
	}
}

var (
	ErrHashAlgorithmUnavailable = errors.New("signature: algorithm unavailable")
)

func NewSignature(options ...Option) (GenerateAndVerify, error) {
	s := Signature{}
	for _, option := range options {
		option(&s)
	}

	if s.Algorithm == nil {
		return nil, ErrHashAlgorithmUnavailable
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
