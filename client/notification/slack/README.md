# Slack

## Introduction
This package is used for configuration request to slack.
What's got in this package.
1. Health - The Health Checks API helps you monitor notification service health
2. Send - Used to send messages to slack channels

## Using Package
```go
    slack := NewSlack(
        validator.New(), // is required field | import from go-lib-common/validator
        WithSentry(sentry), // is required field | import from go-lib-common/sentry
        WithSlackConfigNotificationSlackTimeoutInSeconds(10), // is required field | Setting Notification Slack Timeout
        WithSlackConfigURL(url) // is required field | url notification services e.g https://notification-api.development.jinny.id
        WithSlackConfigChannel(channel) // is required field | the channel to which to go for messaging
    )
```

### Using Health
Health has 1 parameters
1. context

```go
slack.Health(ctx)
```

### Using Send
UserDetail has 2 parameters
1. context
2. messages

```go
slack.Send(ctx, "Hello World")
```

