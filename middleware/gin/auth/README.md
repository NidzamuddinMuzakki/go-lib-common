# Auth

## Introduction
This package is used for auth middleware.
What's got in this package.
1. AuthToken - used for authentication with moladin evo account
2. AuthXApiKey - used for authentication with your app key
3. Auth - used for authentication with moladin evo or authentication api key with new generate

```go
// api key with new generate
token := []byte(xServiceName + xServiceName)
validateKey := sha256.Sum256(token)
if xApiKey != hex.EncodeToString(validateKey[:]) {
  c.JSON(http.StatusUnauthorized, response.Response{
    Message: http.StatusText(http.StatusUnauthorized),
    Status:  response.StatusFail,
  })
  c.Abort()
  return
}
```


## Using Package
```go
    auth := NewAuth(
        validator.New(), // is required field | import from go-lib-common/validator
        WithSentry(sentry), // is required field | import from go-lib-common/sentry
        WithMoladinEvoClient(moladinEvo), // is required field | import from go-lib-common/client/moladin_evo
        WithConfigApiKey(url) // is option field | used for the AuthXApiKey function
        WithPermittedRoles(channel) // is optional field | used for AuthToken and Auth functions | roles that are in moladin evo to be permitted
    )
```

### Using AuthToken

```go
gin.Default().Use(
  auth.AuthToken(),
)
``` 

### Using AuthXApiKey
```go
gin.Default().Use(
  auth.AuthXApiKey(),
)
```

### Using Auth
```go
gin.Default().Use(
  auth.Auth(),
)
```