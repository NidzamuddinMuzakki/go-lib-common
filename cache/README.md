# Cache

## Introduction
This package is used for caching with TTL.

## Supported Driver
- In Memory `cache.InMemory`
- Redis `cache.Redis`

## Usage
```go
package main

import (
	"context"
	"fmt"
	"time"

	"bitbucket.org/moladinTech/go-lib-common/cache"
)

func main() {
	c, err := cache.NewCache(
		cache.WitDriver(cache.InMemoryDriver),
	)
	if err != nil {
		panic(err)
	}

	type Parent struct {
		Name string
	}

	type Child struct {
		Parent Parent
		Name   string
	}

	ctx := context.Background()
	data := cache.Data{
		Key: "test",
		Value: Child{
			Name: "Test 1",
			Parent: Parent{
				Name: "Test 2",
			},
		},
	}

	err = c.Set(ctx, data, time.Minute*1)
	if err != nil {
		panic(err)
	}

	var result Child
	err = c.Get(ctx, data.Key, &result)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}

```