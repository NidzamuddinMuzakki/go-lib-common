package notification

import "context"

type INotification interface {
	Send(ctx context.Context, message string) error
	Health(ctx context.Context) error
	GetFormattedMessage(logCtx string, ctx context.Context, message any) string
}
