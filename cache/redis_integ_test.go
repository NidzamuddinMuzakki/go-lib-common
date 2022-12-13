//go:build integration
// +build integration

package cache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name string `json:"name"`
}

func generateTestKeyAndVal() ([]Data, []Key) {
	testStructs := make([]Data, 0, 10)

	testKeys := make([]Key, 0, 10)

	for i := 0; i < 10; i++ {
		testKeys = append(testKeys, Key(fmt.Sprint(i)))
		testStructs = append(testStructs, Data{
			Key: Key(fmt.Sprint(i)),
			Value: TestStruct{
				Name: fmt.Sprint(i),
			},
		})
	}

	return testStructs, testKeys

}

func TestCache(t *testing.T) {
	redis, err := NewCache(WitHost("localhost:6379"), WitDatabase("0"), WitDriver(RedisDriver), WitPassword(""))
	if err != nil {
		t.Fatal(err)
	}

	inmemory, err := NewCache(WitDriver(InMemoryDriver))
	testCache_(t, redis)
	testCache_(t, inmemory)
}

func testCache_(t *testing.T, cache Cacher) {

	testStructs, testKeys := generateTestKeyAndVal()

	testBatchGetNil(t, cache, testKeys)
	testBatchSet(t, cache, testStructs)
	testBatchGet(t, cache, testKeys)

	testKeys = append(testKeys, Key("none key"))
	testBatchGet(t, cache, testKeys)
	testBatchGetUseSet(t, cache, testKeys)

	testData := Data{
		Key:   Key("testKey"),
		Value: TestStruct{"test"},
	}

	testSet(t, cache, testData)
	testGet(t, cache, testData.Key)
}

func testBatchSet(t *testing.T, cache Cacher, testStructs []Data) {
	err := cache.BatchSet(context.TODO(), testStructs, time.Second)

	assert.Equal(t, nil, err)

}

func testBatchGet(t *testing.T, cache Cacher, testKeys []Key) {
	dest := []TestStruct{}

	err := cache.BatchGet(context.Background(), testKeys, &dest)

	assert.Equal(t, nil, err)
	assert.Equal(t, 10, len(dest))

	for i := 0; i < 10; i++ {
		assert.Equal(t, fmt.Sprint(i), dest[i].Name)
	}
}

func testBatchGetUseSet(t *testing.T, cache Cacher, testKeys []Key) {
	dest := make(map[string]struct{})

	err := cache.BatchGet(context.Background(), testKeys, dest)

	assert.Equal(t, nil, err)
	assert.Equal(t, 10, len(dest))

	for i := 0; i < 10; i++ {
		_, ok := dest[fmt.Sprint(i)]
		assert.Equal(t, true, ok)
	}
}

func testBatchGetNil(t *testing.T, cache Cacher, testKeys []Key) {
	dest := []TestStruct{}

	err := cache.BatchGet(context.Background(), testKeys, &dest)

	assert.Equal(t, nil, err)
}

func testSet(t *testing.T, cache Cacher, test Data) {
	err := cache.Set(context.TODO(), test, time.Second)

	assert.Equal(t, nil, err)
}

func testGet(t *testing.T, cache Cacher, key Key) {
	dest := TestStruct{}
	err := cache.Get(context.Background(), key, &dest)

	assert.Equal(t, nil, err)
	assert.Equal(t, dest.Name, "test")

}
