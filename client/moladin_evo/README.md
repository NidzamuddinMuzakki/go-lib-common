# MoladinEvo

## Introduction
This package is used for configuration request to moladin evo.
What's got in this package.
1. Health - The Health Checks API helps you monitor moladin evo service health
2. UserDetail - Knowing the details of user moladin evo by using tokens as data authentication

## Using Package
```go
    moladin_evo := NewMoladinEvo(
        context,
        validator.New(), // import from go-lib-common/validator
        WithBaseUrl(url), // is required field | moladin evo url e.g https://dev-ucr-api.moladin.com
        WithServicesName("myAppName"), // is required field | Your app name
        WithSentry(url) // is required field | import from go-lib-common/sentry
    )
```

### Using Health
Health has 1 parameters
1. context

```go
moladin_evo.Health(ctx)
```

### Using UserDetail
UserDetail has 2 parameters
1. context
2. token - token from moladin CRM

```go
moladin_evo.UserDetail(ctx,token)
```

