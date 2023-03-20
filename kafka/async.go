package kafka

import (
	"context"
	"time"

	commonSentry "bitbucket.org/moladinTech/go-lib-common/sentry"

	"github.com/Shopify/sarama"
)

type AsyncPublisher struct {
	producer sarama.AsyncProducer
	sentry   commonSentry.ISentry
}

func NewAsyncPublisher(
	brokers []string,
	config *sarama.Config,
	sentry commonSentry.ISentry,
) (*AsyncPublisher, error) {
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &AsyncPublisher{producer: producer, sentry: sentry}, nil
}

func (asp *AsyncPublisher) Publish(ctx context.Context, topic Topic, message IMessage) (int32, int64, error) {
	const logCtx = "kafka.async.AsyncPublisher.Publish"

	if asp.sentry != nil {
		span := asp.sentry.StartSpan(ctx, logCtx)
		ctx = asp.sentry.SpanContext(*span)
		defer asp.sentry.Finish(span)
	}

	messageHeaders := message.GetHeaders(ctx)
	headers := make([]sarama.RecordHeader, 0, len(messageHeaders))
	for key, value := range messageHeaders {
		headers = append(headers, sarama.RecordHeader{
			Key:   []byte(key),
			Value: []byte(value),
		})
	}
	value, err := message.GetValue()
	if err != nil {
		return 0, 0, err
	}

	msg := &sarama.ProducerMessage{
		Topic:     topic.String(),
		Value:     sarama.StringEncoder(value),
		Headers:   headers,
		Metadata:  message.GetMeta(),
		Timestamp: time.Now().UTC(),
	}
	asp.producer.Input() <- msg

	return 0, 0, nil
}
