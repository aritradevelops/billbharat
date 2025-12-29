package eventbroker

import (
	"encoding/json"

	"github.com/aritradevelops/billbharat/backend/shared/logger"
)

type Event interface {
	Topic() string
	Data() []byte
}

type JSONEvent struct {
	topic string
	data  any
}

func NewJSONEvent(topic string, data any) Event {
	return &JSONEvent{
		topic: topic,
		data:  data,
	}
}

func (e *JSONEvent) Topic() string {
	return e.topic
}

func (e *JSONEvent) Data() []byte {
	data, err := json.Marshal(e.data)
	if err != nil {
		logger.Error().Err(err).Msg("failed to marshal event data")
		panic(err)
	}
	return data
}
