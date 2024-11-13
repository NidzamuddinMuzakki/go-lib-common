package notification

import "context"

type INotification interface {
	Send(ctx context.Context, message string) error
	Health(ctx context.Context) error
	GetFormattedMessage(ctx context.Context, logCtx string, message any) string
}
