package eventbroker

import (
	"context"
	"time"

	"github.com/aritradevelops/billbharat/backend/shared/logger"
	kafka "github.com/segmentio/kafka-go"
)

const (
	BatchTimeout = 0 * time.Second
)

type Producer interface {
	Produce(ctx context.Context, event Event) error
}

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(servers []string) Producer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      servers,
		Balancer:     &kafka.RoundRobin{},
		BatchTimeout: BatchTimeout,
		Dialer: &kafka.Dialer{
			Timeout: 10 * time.Second,
		},
		WriteTimeout: 10 * time.Second,
		BatchSize:    0,
		Async:        true,
	})
	return &KafkaProducer{
		writer: writer,
	}
}

func (p *KafkaProducer) Produce(ctx context.Context, event Event) error {
	err := p.writer.WriteMessages(ctx, kafka.Message{
		Topic: event.Topic(),
		Value: event.Data(),
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to produce event")
		return err
	} else {
		logger.Info().Str("topic", event.Topic()).Msg("event produced successfully")
	}
	return nil
}
