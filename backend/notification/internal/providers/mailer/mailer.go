package mailer

import (
	"fmt"

	"github.com/aritradevelops/billbharat/backend/shared/logger"
	"github.com/aritradevelops/billbharat/backend/shared/notification"
)

type Mailer interface {
	Send(data notification.EmailData, subject string, body string, alternativeBody *string) error
}

type Emailer struct {
	// some config
}

type Logger struct {
	// some config
}

func New(env string) Mailer {
	if env == "production" {
		return &Emailer{}
	}
	return &Logger{}
}

func (e *Emailer) Send(email notification.EmailData, subject string, body string, alternativeBody *string) error {
	return fmt.Errorf("not implemented")
}

func (l *Logger) Send(email notification.EmailData, subject string, body string, alternativeBody *string) error {
	logger.Info().Interface("email", email).Interface("subject", subject).Interface("body", body).Interface("alternativeBody", alternativeBody).Msg("sending email")
	return nil
}
