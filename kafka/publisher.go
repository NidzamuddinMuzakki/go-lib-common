package kafka

import (
	"context"
)

type Topic string

func (t Topic) String() string {
	return string(t)
}

type IPublisher interface {
	Publish(ctx context.Context, topic Topic, message IMessage) (int32, int64, error)
}
