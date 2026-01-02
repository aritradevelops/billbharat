package mailer

import (
	"fmt"
	"time"

	"github.com/aritradevelops/billbharat/backend/notification/internal/config"
	"github.com/aritradevelops/billbharat/backend/shared/notification"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mailer interface {
	Send(data notification.EmailData, subject string, body string, alternativeBody *string) error
}

type Emailer struct {
	// some config
}

type Mailhog struct {
	domain   string
	host     string
	port     int
	username string
	password string
	from     string
	fromName string
}

func New(env string, config config.Mailer) Mailer {
	if env == "production" {
		return &Emailer{}
	}
	return &Mailhog{
		domain:   config.Domain,
		host:     config.Host,
		port:     config.Port,
		username: config.Username,
		password: config.Password,
		from:     config.From,
		fromName: config.FromName,
	}
}

func (e *Emailer) Send(email notification.EmailData, subject string, body string, alternativeBody *string) error {
	return fmt.Errorf("not implemented")
}

func (l *Mailhog) Send(email notification.EmailData, subject string, body string, alternativeBody *string) error {

	server := mail.NewSMTPClient()
	server.Host = l.host
	server.Port = l.port
	server.Username = l.username
	server.Password = l.password
	server.Encryption = mail.EncryptionNone
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	client, err := server.Connect()
	if err != nil {
		return err
	}

	message := mail.NewMSG()
	message.SetFrom(fmt.Sprintf("%s <%s>", l.fromName, l.from))
	message.AddTo(email.To...)
	message.SetSubject(subject)
	message.SetBody(mail.TextHTML, body)
	if alternativeBody != nil {
		message.AddAlternative(mail.TextPlain, *alternativeBody)
	}

	if err := message.Send(client); err != nil {
		return err
	}

	return nil
}
