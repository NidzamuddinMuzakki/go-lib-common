# Sentry

## Introduction
This package is used for sentry configurations 
What's got in this package.
1. SetStartTransaction - used to creates a transaction span to time runs of an expensive operation on items from a channel. Timing for each operation is sent to Sentry and grouped by transaction name
2. Trace - used to describe each operation
3. StartSpan - used to starting the first span in a transaction
4. Finish - used to finishing the span in a transaction
5. SetTag - used to set tag in transaction
6. CaptureException - used to capture an event in Go, you can pass any struct implementing an error interface to CaptureException
7. GetGinMiddleware - used to middleware default sentry gin.
8. Flush - used to waits until any buffered events are sent to the Sentry server, blocking for at most the given timeout
9. SetUserInfo - used to describe user information
10. HandlingPanic - used to capture unhandled panics in our Go SDK is through the Recover method.

## Using Package
```go
    sentry := NewSentry(
        validator.New(), // is required field | import from go-lib-common/validator
        WithDsn(Dsn), // is required field | dsn sentry
        WithDebug(Debug), // is optional field (default: false) | debug sentry
        WithEnv(Env), // is required field | environment your app
        WithSampleRate(SampleRate), // is required field | To send a representative sample of your errors to Sentry, set the SampleRate option in your SDK configuration to a number between 0 (0% of errors sent) and 1 (100% of errors sent)
        WithBlacklistTransactions([]string{"POST /v1/sub-sku/cron/late-check-in", "GET /v1/sub-sku/:subSKUID/bpkb/number"}),
    )
```

### Using SetStartTransaction

```go
    sentry.SetStartTransaction(
        context,
        "middleware.Test",
        fmt.Sprintf("%s %s", c.Request.Method, c.FullPath()),
        func(ctx context.Context) (string, uint8) {
            err := reconService.SyncMoneyOut(ctx)
            if err != nil {
                return "500", uint8(sentry.STATUS_INTERNAL_SERVER_ERROR)
            }
            return "200", uint8(sentry.STATUS_OK)
        },
    )
```

### Using SetStartTransaction
```go
    sentry.SetStartTransaction(
        context,
        "middleware.Test",
        fmt.Sprintf("%s %s", c.Request.Method, c.FullPath()),
        func(ctx context.Context) (string, uint8) {
            err := reconService.SyncMoneyOut(ctx)
            if err != nil {
                return "500", uint8(sentry.STATUS_INTERNAL_SERVER_ERROR)
            }
            return "200", uint8(sentry.STATUS_OK)
        },
    )
```

### Using Trace
```go
    sentry.Trace(ctx, "reconService.SyncMoneyOut", func(ctx context.Context) {
		err := reconService.SyncMoneyOut(ctx)
	})
```

### Using StartSpan and Finish
```go
    span := sentry.StartSpan(ctx, spanName)
	defer sentry.Finish(span)
```

### Using SetTag
```go
    span := sentry.StartSpan(ctx, spanName)
	defer sentry.Finish(span)
    sentry.SetTag(span,"myKey","myValue")
```

### Using CaptureException
```go
    eventID := sentry.CaptureException(erros.New("error"))
```

### Using GetGinMiddleware
```go
   gin.Default().Use(
        sentry.GetGinMiddleware(),
   )
```

### Using Flush
```go
    sentry.Flush(10), //in second
```

### Using HandlingPanic
```go
    sentry.HandlingPanic(erros.New("error"))
```

### Using SetUserInfo
```go
   sentry.SetUserInfo(sentry.UserInfoSentry{
        ID: "1",
        Username: "john",
        Email: "john@example.com",
    })
```
