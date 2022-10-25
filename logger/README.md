# Logger

## Introduction
This package is used for configuration logger
What's got in this package.
1. Debug - used to provide global zerolog log Debug.
2. Info - used to provide global zerolog log Info.
3. Warn - used to provide global zerolog log Warn.
4. Error - used to provide global zerolog log Error.
5. Fatal - used to provide global zerolog log Fatal.

## Using Package

```go
	logger.Init(logger.Config{
		AppName: "My App Name",
		Debug:   true, //true or false
	})
```

### Using Debug
```go
    logger.Debug(ctx, "debug run engine", logger.Tag{Key: "error", Value: err.Error()})
```

### Using Info
```go
    logger.Info(ctx, "Get data finished", logger.Tag{Key: "data finished", Value: myData})
```

### Using Warn
```go
    logger.Fatal(ctx, "warn data is empty", logger.Tag{Key: "warning", Value: myData})
```

### Using Error
```go
    err := errors.New("error")
    logger.Fatal(ctx, "erros on get data", logger.Tag{Key: "error", Value: err.Error()})
```

### Using Fatal
```go
    err := errors.New("error")
    logger.Fatal(ctx, "failed to run engine", logger.Tag{Key: "error", Value: err.Error()})
```