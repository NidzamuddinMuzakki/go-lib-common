package signature

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignature(t *testing.T) {
	t.Parallel()

	sha256, err := NewSignature(WithAlgorithm(Sha256))
	if err != nil {
		t.Fatal(err)
	}

	bcrypt, err := NewSignature(WithAlgorithm(BCrypt), WithCost(10))
	if err != nil {
		t.Fatal(err)
	}

	testSiganture_ErrorInitializeAlgorithmEmpty(t)
	testSiganture_InitializeToDefaultWhenCostEmpty(t)
	testSignature_(t, sha256)
	testSignature_(t, bcrypt)
}

func testSignature_(t *testing.T, signature GenerateAndVerify) {
	testGenerate_ShouldSuccess(t, signature)
	testVerify_ShouldMatchingSiganture(t, signature)
	testVerify_ErrorNotMatchSignature(t, signature)
}

// Initialize signature
func testSiganture_ErrorInitializeAlgorithmEmpty(t *testing.T) {
	_, err := NewSignature()
	assert.Equal(t, ErrHashAlgorithmUnavailable, err)
}

func testSiganture_InitializeToDefaultWhenCostEmpty(t *testing.T) {
	_, err := NewSignature(WithAlgorithm(BCrypt))
	assert.Equal(t, nil, err)
}

// Function Generate
func testGenerate_ShouldSuccess(t *testing.T, signature GenerateAndVerify) {
	ctx := context.Background()
	_, err := signature.Generate(ctx, "dummy")

	assert.Equal(t, nil, err)
}

// Function Verify
func testVerify_ShouldMatchingSiganture(t *testing.T, signature GenerateAndVerify) {
	ctx := context.Background()

	hashed, _ := signature.Generate(ctx, "dummy")
	match := signature.Verify(ctx, "dummy", hashed)

	assert.Equal(t, true, match)
}

func testVerify_ErrorNotMatchSignature(t *testing.T, signature GenerateAndVerify) {
	ctx := context.Background()

	hashed, _ := signature.Generate(ctx, "dummy")
	match := signature.Verify(ctx, "dumm", hashed)

	assert.Equal(t, false, match)
}
