# Config

## Introduction
This package is used for the connection of each configuration.
What's got in this package.
1. BindFromFile - used for config load from filename then assign to destination
2. BindFromConsul - used for config load from remote consul then assign to destination
3. BindAndWatchFromConsul - used for load and watch config from remote consul then assign to destination
4. LoadConsulIntervalFromEnv - used for get interval value for loading config from consul

## Using Package
### Using BindFromFile
You must have a struct named ColdFlat that will be used for the binding value 

```go
type ColdFlat struct {
	AppEnv      string `json:"appEnv" yaml:"appEnv"`
	AppName     string `json:"appName" yaml:"appName"`
}
```

and also you have a json file named config.cold.json
```json
{
  "appEnv": "development",
  "appName": "moladin-go-skeleton-service",
}
```

After that you can call the function BindFromFile
```go
err := config.BindFromFile(&cfg.Cold, "config.cold.json", ".")
```

### Using BindFromConsul
You are required to create a property CONSUL_HTTP_TOKEN in the .env and config.cold.json files
```file
CONSUL_HTTP_TOKEN=xxx
```

```go
err := config.BindFromConsul(
    &ColdFlat,
    consulURL, // Consul URL
    "moladin-go-skeleton-service/backend/cold", // path file config.cold.json in consul
)
```

### Using BindAndWatchFromConsul
You are required to create a property CONSUL_HTTP_TOKEN in the .env and config.cold.json files
```file
CONSUL_HTTP_TOKEN=xxx
```

```go
err := config.BindAndWatchFromConsul(
    &ColdFlat,
    consulURL, // Consul URL
    "moladin-go-skeleton-service/backend/cold", // path file config.cold.json in consul
    10, // interval reloading watch value in consul
)
```
### Using LoadConsulIntervalFromEnv

```go
interval, err := config.LoadConsulIntervalFromEnv()
```