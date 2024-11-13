# NotificationService

## Introduction
This package is used for configuration request to [notifcation services](https://notification-api.production.mofi.id/swagger/index.html).
What's got in this package.
1. Health - The Health Checks API helps you monitor notification service health
2. SendSlack - Used to send messages to slack channels
3. SendEmail - Used to send messages to email

## Using Package
```go
    notifService := NewNotificationService(
        validator.New(), // is required field | import from go-lib-common/validator
        WithSentry(sentry), // is required field | import from go-lib-common/sentry
        WithTimeoutInSeconds(10), // is required field | Setting Notification Slack Timeout
        WithURL("https://notification-api.production.mofi.id") // is required field | Setting Notification URL
        WithServiceEnv(10), // is optional field | Setting Notification Environment
        WithServiceName("go-lib-common"), // is optional field | Setting Notification Services Name
    )
```

### Using Health
Health has 1 parameters
1. context

```go
notifService.Health(ctx)
```

### Using Send Slack
SendSlack has 2 parameters
1. context
2. messages

```go
notifService.SendSlack(ctx, "test-notification-channel", "Hello World")
```

### Using Send Email
SendEmail has 2 parameters
1. context
2. messages

```go
notifService.SendEmail(ctx, []string{"employee@moladin.com"}, "Test Subject", "Hello World")
```

