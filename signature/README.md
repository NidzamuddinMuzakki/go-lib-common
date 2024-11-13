# Signature

## Introduction
This package is used for generating and verifying signature.

## Supported Algorithm
- SHA256 `signature.SHA256`
- BCrypt `signature.BCrypt`

## Usage
```go
package main

import (
	"context"
    "fmt"

    "bitbucket.org/moladinTech/go-lib-common/signature"
)

func main() {
	ctx := context.Background()

	s, err := signature.NewSignature(
		signature.WithAlgorithm(signature.BCrypt),
	)
	if err != nil {
		panic(err)
	}

	key := "{service-sender}:{service-receiver}:{x-request-id}:{x-request-at}:{secret-key}"
	
	hashed, err := s.Generate(ctx,
		key)
	if err != nil {
		panic(err)
	}

	match := s.Verify(ctx, key, hashed)
	if !match {
		panic("not match")
	}

	fmt.Printf("Success")
}
```

