package kafka

import (
	commonContext "bitbucket.org/moladinTech/go-lib-common/context"
	"context"
	"encoding/json"
	"time"

	"bitbucket.org/moladinTech/go-lib-common/constant"
)

type IMessage interface {
	GetHeaders() map[string]string
	GetMeta() any
	GetValue() (string, error)
}

type EventName string

type MessageEvent struct {
	Name EventName `json:"name"`
}

type MessageMeta struct {
	Sender    string     `json:"sender"`
	SendingAt time.Time  `json:"sendingAt"`
	ExpiredAt *time.Time `json:"expiredAt"`
	Version   *string    `json:"version"`
}

type DataType string

const (
	JSON   DataType = "JSON"
	Byte   DataType = "BYTE"
	String DataType = "STRING"
)

type MessageBody[T any] struct {
	Type DataType `json:"type"`
	Data T        `json:"data"`
}

type Message[T any] struct {
	Event MessageEvent   `json:"event"`
	Meta  MessageMeta    `json:"meta"`
	Body  MessageBody[T] `json:"body"`
}

func NewMessage[T any](event MessageEvent, meta MessageMeta, bodyType DataType, body T) *Message[T] {
	return &Message[T]{
		Event: event,
		Meta:  meta,
		Body: MessageBody[T]{
			Type: bodyType,
			Data: body,
		},
	}
}

func (m *Message[T]) GetHeaders(ctx context.Context) map[string]string {
	headers := make(map[string]string, 1)
	headers[constant.XRequestIdHeader] = commonContext.GetValueAsString(ctx, constant.XRequestIdHeader)
	return headers
}

func (m *Message[T]) GetValue() (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (m *Message[T]) GetMeta() any {
	return m.Meta
}
