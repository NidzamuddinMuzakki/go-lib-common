# Panic Tracer

## Introduction
This package is used for sentry and logging for tracer your services.
What's got in this package.
1. Tracer - used for tracer your services


## Using Package
```go
    tracer := NewTracer(
        validator.New(), // is required field | import from go-lib-common/validator
        WithSentry(sentry), // is required field | import from go-lib-common/sentry
    )
```

### Using Tracer
```go
tracer.Tracer(),
``` 