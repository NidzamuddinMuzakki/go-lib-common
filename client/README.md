# Client

## Introduction
This package is used for request configuration to access resources on the server.
What's got in this package.
1. AWS - Contains a configuration request to aws.
2. Moladin Evo - Contains a configuration request to moladin evo.
3. Notification - Contains a file requests for notifications-related

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