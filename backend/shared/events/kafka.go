package events

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/aritradevelops/billbharat/backend/shared/logger"
	"github.com/segmentio/kafka-go"
)

type KafkaOpts struct {
	Servers []string
	GroupId string
}

// kafka implements event manager
type Kafka struct {
	servers []string
	groupId string
	writer  *kafka.Writer
}

func NewKafkaEventManager(opts KafkaOpts) EventManager {
	return &Kafka{
		servers: opts.Servers,
		groupId: opts.GroupId,
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  opts.Servers,
			Balancer: &kafka.RoundRobin{},
			Dialer: &kafka.Dialer{
				Timeout: 10 * time.Second,
			},
			Async: true,
		}),
	}
}

func (k *Kafka) EmitManageUserEvent(ctx context.Context, data EventPayload[ManageUserEventPayload]) error {
	return k.produce(ctx, ManageUserEvent, data)
}

func (k *Kafka) EmitManageNotificationEvent(ctx context.Context, data EventPayload[MangageNotificationEventPayload]) error {
	return k.produce(ctx, ManageNotification, data)
}

func (k *Kafka) EmitManageBusinessEvent(ctx context.Context, data EventPayload[MangageBusinessEventPayload]) error {
	return k.produce(ctx, ManageBusinessEvent, data)
}

func (k *Kafka) EmitManageBusinessUserEvent(ctx context.Context, data EventPayload[MangageBusinessUserEventPayload]) error {
	return k.produce(ctx, ManageBusinessUserEvent, data)
}

func (k *Kafka) OnManageUserEvent(ctx context.Context, handler func(EventPayload[ManageUserEventPayload]) error) {
	go startKafkaConsumer(ctx, k.newReader(ManageUserEvent), handler)
}

func (k *Kafka) OnManageNotificationEvent(ctx context.Context, handler func(EventPayload[MangageNotificationEventPayload]) error) {
	go startKafkaConsumer(ctx, k.newReader(ManageNotification), handler)
}

func (k *Kafka) OnManageBusinessEvent(ctx context.Context, handler func(EventPayload[MangageBusinessEventPayload]) error) {
	go startKafkaConsumer(ctx, k.newReader(ManageBusinessEvent), handler)
}

func (k *Kafka) OnManageBusinessUserEvent(ctx context.Context, handler func(EventPayload[MangageBusinessUserEventPayload]) error) {
	go startKafkaConsumer(ctx, k.newReader(ManageBusinessUserEvent), handler)
}

func (k *Kafka) newReader(e Event) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  k.servers,
		GroupID:  k.groupId,
		Topic:    string(e),
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
}

func (k *Kafka) produce(ctx context.Context, event Event, data any) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return k.writer.WriteMessages(ctx, kafka.Message{
		Topic: string(event),
		Value: dataBytes,
	})
}

func startKafkaConsumer[T any](
	ctx context.Context,
	reader *kafka.Reader,
	handler func(T) error,
) {
	defer reader.Close()
	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if err == io.EOF {
				return
			}
			logger.Error().Err(err).Msg("fetch failed")
			continue
		}

		var event T
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			logger.Error().Err(err).Msg("unmarshal failed")
			continue
		}

		if err := handler(event); err != nil {
			logger.Error().Err(err).Msg("handler failed")
			continue
		}

		if err := reader.CommitMessages(ctx, msg); err != nil {
			logger.Error().Err(err).Msg("commit failed")
		}
	}
}
