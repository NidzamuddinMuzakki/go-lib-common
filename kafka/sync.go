package kafka

import (
	"context"
	"time"

	commonSentry "bitbucket.org/moladinTech/go-lib-common/sentry"

	"github.com/Shopify/sarama"
)

type SyncPublisher struct {
	producer sarama.SyncProducer
	sentry   commonSentry.ISentry
}

func NewSyncPublisher(
	brokers []string,
	config *sarama.Config,
	sentry commonSentry.ISentry,
) (*SyncPublisher, error) {
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &SyncPublisher{producer: producer, sentry: sentry}, nil
}

func (sp *SyncPublisher) Publish(ctx context.Context, topic Topic, message IMessage) (int32, int64, error) {
	const logCtx = "kafka.sync.SyncPublisher.Publish"

	if sp.sentry != nil {
		span := sp.sentry.StartSpan(ctx, logCtx)
		ctx = sp.sentry.SpanContext(*span)
		defer sp.sentry.Finish(span)
	}

	messageHeaders := message.GetHeaders()
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
	partition, offset, err := sp.producer.SendMessage(msg)
	if err != nil {
		return 0, 0, err
	}

	return partition, offset, nil
}
