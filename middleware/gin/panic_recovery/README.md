# Panic Recovery

## Introduction
This package is used for panic recovery handler.
What's got in this package.
1. PanicRecoveryMiddleware - used for panic recovery handler


## Using Package
```go
    panicRecovery := NewPanicRecovery(
        validator.New(), // is required field | import from go-lib-common/validator
        WithSentry(sentry), // is required field | import from go-lib-common/sentry
        WithConfigEnv(environment), // is required field | what is your environment? development or staging or production
    )
```

### Using PanicRecoveryMiddleware
```go
panicRecovery.PanicRecoveryMiddleware(),
``` 