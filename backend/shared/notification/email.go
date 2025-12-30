package notification

import (
	"strings"

	"github.com/aritradevelops/billbharat/backend/shared/logger"
)

// EmailChannel implements Channel interface
type EmailChannel struct {
	To        []string  `json:"to"`
	CC        []string  `json:"cc"`
	BCC       []string  `json:"bcc"`
	EventType EventType `json:"event_type"`
}

func (e *EmailChannel) Send() error {
	if e.EventType == Broadcast {
		logger.Info().Msgf("Sending to all recipients: %v", strings.Join(e.To, ", "))
	} else {
		for _, to := range e.To {
			logger.Info().Msgf("Sending to recipient: %s", to)
		}
	}
	return nil
}

func NewEmailChannel(to []string, cc []string, bcc []string, eventType EventType) *EmailChannel {
	return &EmailChannel{
		To:        to,
		CC:        cc,
		BCC:       bcc,
		EventType: eventType,
	}
}
