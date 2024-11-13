# Usage

- Install Sarama in your project

```go
go get -u github.com/Shopify/sarama@v1.38.1
```

- Create new Publisher (Sync / ASync)

```go
// main.go
import (
	"bitbucket.org/moladinTech/go-lib-common/kafka"
    "bitbucket.org/moladinTech/go-lib-common/registry"
    "bitbucket.org/moladinTech/go-lib-common/sentry"
	
    "github.com/Shopify/sarama"
)

brokers := []string{"1.1.1.1:1234"}
saramaConfig := sarama.NewConfig()

// sync
publisherSync, err := kafka.NewSyncPublisher(brokers, saramaConfig, sentry)
if err != nil; {
	...
}
// async
publisherAsync, err := kafka.NewAsyncPublisher(brokers, saramaConfig, sentry)
if err != nil; {
	...
}

// Initialization
message := NewMessage[T](event, meta, bodyType, body[T])

partition, offset, err := kafka.publisherSync.Publish(ctx, topic, message)
if err != nil {
	....
}

// when using async the partition, offset, and error always return 0, 0, nil
partition, offset, err := publisherAsync.Publish(ctx, topic, message)
if err != nil {
    ....
}
```

## Optional
- You can add multiple publishers to common registry too

```go
commonRegistry := registry.NewRegistry(
    registry.AddPublisher('sync', publisherSync),
    registry.AddPublisher('async', publisherAsync),
)

commonRegistry.GetPublisher('sync').Publish(...)
```

- Using custom Publisher you must implement `IPublisher` interface then publish the message using publisher (Sync / Async)

- Using custom Message you must implement `IMessage` interface then publish the message using publisher (Sync / Async)
