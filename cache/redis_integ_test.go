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

func TestRedis(t *testing.T) {
	redis := NewRedis("localhost:6379", "", 0)

	testStructs, testKeys := generateTestKeyAndVal()

	testBatchGetNil(t, redis, testKeys)
	testBatchSet(t, redis, testStructs)
	testBatchGet(t, redis, testKeys)

	testKeys = append(testKeys, Key("none key"))
	testBatchGet(t, redis, testKeys)
	testBatchGetUseSet(t, redis, testKeys)

	testData := Data{
		Key:   Key("testKey"),
		Value: TestStruct{"test"},
	}

	testSet(t, redis, testData)
	testGet(t, redis, testData.Key)

}

func testBatchSet(t *testing.T, redis *Redis, testStructs []Data) {
	err := redis.BatchSet(context.TODO(), testStructs, time.Second)

	assert.Equal(t, nil, err)

}

func testBatchGet(t *testing.T, redis *Redis, testKeys []Key) {
	dest := []TestStruct{}

	err := redis.BatchGet(context.Background(), testKeys, &dest)

	assert.Equal(t, nil, err)
	assert.Equal(t, 10, len(dest))

	for i := 0; i < 10; i++ {
		assert.Equal(t, fmt.Sprint(i), dest[i].Name)
	}
}

func testBatchGetUseSet(t *testing.T, redis *Redis, testKeys []Key) {
	dest := make(map[string]struct{})

	err := redis.BatchGet(context.Background(), testKeys, dest)

	assert.Equal(t, nil, err)
	assert.Equal(t, 10, len(dest))

	for i := 0; i < 10; i++ {
		_, ok := dest[fmt.Sprint(i)]
		assert.Equal(t, true, ok)
	}
}

func testBatchGetNil(t *testing.T, redis *Redis, testKeys []Key) {
	dest := []TestStruct{}

	err := redis.BatchGet(context.Background(), testKeys, &dest)

	assert.Equal(t, nil, err)
}

func testSet(t *testing.T, redis *Redis, test Data) {
	err := redis.Set(context.TODO(), test, time.Second)

	assert.Equal(t, nil, err)
}

func testGet(t *testing.T, redis *Redis, key Key) {
	dest := TestStruct{}
	err := redis.Get(context.Background(), key, &dest)

	assert.Equal(t, nil, err)
	assert.Equal(t, dest.Name, "test")

}
