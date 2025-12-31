package smsprovider

import (
	"fmt"

	"github.com/aritradevelops/billbharat/backend/shared/logger"
	"github.com/aritradevelops/billbharat/backend/shared/notification"
)

type SMSProvider interface {
	Send(data notification.SMSData, subject string, body string) error
}

type SMSProviderImpl struct {
	// some config
}
type Logger struct {
	// some config
}

func New(env string) SMSProvider {
	if env == "production" {
		return &SMSProviderImpl{}
	}
	return &Logger{}
}

func (s *SMSProviderImpl) Send(data notification.SMSData, subject string, body string) error {
	return fmt.Errorf("not implemented")
}

func (s *Logger) Send(data notification.SMSData, subject string, body string) error {
	logger.Info().Interface("data", data).Interface("subject", subject).Interface("body", body).Msg("sending sms")
	return nil
}
