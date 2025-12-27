package eventbroker

import (
	"context"

	"github.com/aritradeveops/billbharat/backend/auth/internal/pkg/logger"
	kafka "github.com/segmentio/kafka-go"
)

type Producer interface {
	Produce(ctx context.Context, topic string, event Event) error
}

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(servers []string) Producer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  servers,
		Balancer: &kafka.LeastBytes{},
	})
	return &KafkaProducer{
		writer: writer,
	}
}

func (p *KafkaProducer) Produce(ctx context.Context, topic string, event Event) error {
	err := p.writer.WriteMessages(ctx, kafka.Message{
		Topic: event.Topic(),
		Value: event.Data(),
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to produce event")
		return err
	}
	return nil
}
