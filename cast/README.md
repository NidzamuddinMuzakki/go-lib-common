# Cast

## Introduction
This package is used for a casting variables.
What's got in this package.
1. NewPointer

## Using Package

### Using NewPointer
```go
    val := cast.NewPointer[float64](10)
```

```go
    m := []modelDB.MasterSKU{
    {
        SKUID: `mantab jiwa`,
    },
	}
	n := make([]model.WithDetail, 0, len(m))

	modelDB.Convert(m, &n, func(k *modelDB.MasterSKU) model.WithDetail {
		return model.WithDetail{
			SKUID: k.SKUID,
		}
	})
	assert.Equal(t, len(m), len(n)) // 1 == 1
	assert.Equal(t, m[0].SKUID, n[0].SKUID) // mantab jiwa = mantab jiwa
```