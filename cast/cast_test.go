//go:build unit
// +build unit

package cast_test

import (
	"testing"

	"bitbucket.org/moladinTech/go-lib-common/cast"
	"github.com/go-playground/assert/v2"
)

func TestConvert(t *testing.T) {
	type ModelTwo struct {
		Wuhu string
	}
	type ModelOne struct {
		Name string
	}
	m := []ModelOne{
		{
			Name: `mantab jiwa`,
		},
	}
	n := make([]ModelTwo, 0, len(m))

	cast.Convert(&m, &n, func(k *ModelOne) ModelTwo {
		return ModelTwo{
			Wuhu: k.Name,
		}
	})
	assert.Equal(t, len(m), len(n))       // 1 == 1
	assert.Equal(t, m[0].Name, n[0].Wuhu) // mantab jiwa = mantab jiwa
}

func TestConvertWithAllocator(t *testing.T) {
	type ModelTwo struct {
		Wuhu string
	}
	type ModelOne struct {
		Name string
	}
	m := []ModelOne{
		{
			Name: `mantab jiwa`,
		},
	}
	n := []ModelTwo{}

	cast.ConvertAndAllocate(&m, &n, func(k *ModelOne) ModelTwo {
		return ModelTwo{
			Wuhu: k.Name,
		}
	})
	assert.Equal(t, len(m), len(n))       // 1 == 1
	assert.Equal(t, m[0].Name, n[0].Wuhu) // mantab jiwa = mantab jiwa
}

func BenchmarkConvert(b *testing.B) {
	type ModelTwo struct {
		Wuhu string
	}
	type ModelOne struct {
		Name string
	}

	for i := 0; i < b.N; i++ {
		m := make([]ModelOne, 0, 100)
		for j := 0; j < 100; j++ {
			m = append(m, ModelOne{
				Name: `mantab jiwa`,
			})
		}
		n := make([]ModelTwo, 0, len(m))

		cast.Convert(&m, &n, func(k *ModelOne) ModelTwo {
			return ModelTwo{
				Wuhu: k.Name,
			}
		})
	}
}

func BenchmarkConvertWithAllocate(b *testing.B) {
	type ModelTwo struct {
		Wuhu string
	}
	type ModelOne struct {
		Name string
	}

	for i := 0; i < b.N; i++ {
		m := make([]ModelOne, 0, 100)
		for j := 0; j < 100; j++ {
			m = append(m, ModelOne{
				Name: `mantab jiwa`,
			})
		}
		n := []ModelTwo{}

		cast.ConvertAndAllocate(&m, &n, func(k *ModelOne) ModelTwo {
			return ModelTwo{
				Wuhu: k.Name,
			}
		})
	}
}

// user using standard convert but forgot to allocate the memory
func BenchmarkConvertWithoutAllocate(b *testing.B) {
	type ModelTwo struct {
		Wuhu string
	}
	type ModelOne struct {
		Name string
	}

	for i := 0; i < b.N; i++ {
		m := make([]ModelOne, 0, 100)
		for j := 0; j < 100; j++ {
			m = append(m, ModelOne{
				Name: `mantab jiwa`,
			})
		}
		n := []ModelTwo{}

		cast.Convert(&m, &n, func(k *ModelOne) ModelTwo {
			return ModelTwo{
				Wuhu: k.Name,
			}
		})
	}
}
