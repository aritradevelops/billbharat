package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"text/template"

	"github.com/aritradevelops/billbharat/backend/notification/internal/persistence/repository"
	"github.com/aritradevelops/billbharat/backend/notification/internal/providers/mailer"
	"github.com/aritradevelops/billbharat/backend/notification/internal/providers/smsprovider"
	"github.com/aritradevelops/billbharat/backend/notification/internal/providers/templatestore"
	"github.com/aritradevelops/billbharat/backend/shared/events"
	"github.com/aritradevelops/billbharat/backend/shared/logger"
	"github.com/aritradevelops/billbharat/backend/shared/notification"
)

type Notifier interface {
	Notify(ctx context.Context, payload events.EventPayload[events.MangageNotificationEventPayload]) error
}

type NotifierImpl struct {
	ctx           context.Context
	mailer        mailer.Mailer
	smsProvider   smsprovider.SMSProvider
	templateStore templatestore.TemplateStorage
}

func New(ctx context.Context, env string, repo repository.Repository) Notifier {
	return &NotifierImpl{
		ctx:           ctx,
		mailer:        mailer.New(env),
		smsProvider:   smsprovider.New(env),
		templateStore: templatestore.New(env, repo),
	}
}

func (n *NotifierImpl) Notify(ctx context.Context, payload events.EventPayload[events.MangageNotificationEventPayload]) error {
	for _, msg := range payload.Data.Payload {
		switch msg.Channel {
		case notification.EMAIL:
			err := n.handleEmailNotification(ctx, payload.Data.Event, payload.Data.Kind, msg, payload.Data.Tokens)
			if err != nil {
				return err
			}
		case notification.SMS:
			err := n.handleSMSSNotification(ctx, payload.Data.Event, payload.Data.Kind, msg, payload.Data.Tokens)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (n *NotifierImpl) handleEmailNotification(ctx context.Context, event notification.Event, kind notification.Kind, data events.NotificationChannelPayload, tokens any) error {
	var emailData notification.EmailData
	if data.Data == nil {
		logger.Error().Msg("email data is nil")
		return nil
	}
	dbByte, _ := json.Marshal(data.Data)
	err := json.Unmarshal(dbByte, &emailData)
	if err != nil {
		logger.Error().Err(err).Msg("failed to unmarshal email channel")
		return err
	}

	htmlTemplate, err := n.templateStore.FindTemplate(n.ctx, templatestore.FindTemplateParams{
		Event:    event,
		Channel:  data.Channel,
		Locale:   "en",
		Scope:    "default",
		Mimetype: "text/html",
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to find html template")
		return err
	}

	textTemplate, err := n.templateStore.FindTemplate(n.ctx, templatestore.FindTemplateParams{
		Event:    event,
		Channel:  data.Channel,
		Locale:   "en",
		Scope:    "default",
		Mimetype: "text/plain",
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to find text template")
	}
	body, err := n.compileTemplate(htmlTemplate.Body, tokens)
	if err != nil {
		logger.Error().Err(err).Msg("failed to compile html template")
		return err
	}

	subject, err := n.compileTemplate(htmlTemplate.Subject, tokens)
	if err != nil {
		logger.Error().Err(err).Msg("failed to compile html template")
		return err
	}
	var textBody *string = nil
	if textTemplate.Body != "" {
		tb, err := n.compileTemplate(textTemplate.Body, tokens)
		if err != nil {
			logger.Error().Err(err).Msg("failed to compile text template")
			return err
		}
		textBody = &tb
	}

	err = n.mailer.Send(emailData, subject, body, textBody)
	if err != nil {
		logger.Error().Err(err).Msg("failed to send email")
		return err
	}
	return nil
}

func (n *NotifierImpl) handleSMSSNotification(ctx context.Context, event notification.Event, kind notification.Kind, data events.NotificationChannelPayload, tokens any) error {
	var smsData notification.SMSData
	if data.Data == nil {
		logger.Error().Msg("sms data is nil")
		return nil
	}
	dbByte, _ := json.Marshal(data.Data)
	err := json.Unmarshal(dbByte, &smsData)
	if err != nil {
		logger.Error().Err(err).Msg("failed to unmarshal sms channel")
		return err
	}

	smsTemplate, err := n.templateStore.FindTemplate(n.ctx, templatestore.FindTemplateParams{
		Event:    event,
		Channel:  data.Channel,
		Locale:   "en",
		Scope:    "default",
		Mimetype: "text/plain",
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to find sms template")
		return err
	}

	body, err := n.compileTemplate(smsTemplate.Body, tokens)
	if err != nil {
		logger.Error().Err(err).Msg("failed to compile sms template")
		return err
	}
	subject, err := n.compileTemplate(smsTemplate.Subject, tokens)
	if err != nil {
		logger.Error().Err(err).Msg("failed to compile sms subject")
		return err
	}

	err = n.smsProvider.Send(smsData, subject, body)
	if err != nil {
		logger.Error().Err(err).Msg("failed to send sms")
		return err
	}
	return nil
}

func (n *NotifierImpl) compileTemplate(tmpl string, tokens any) (string, error) {
	template, err := template.New("any").Parse(tmpl)
	if err != nil {
		logger.Error().Err(err).Msg("failed to parse sms template")
		return "", err
	}
	body := bytes.Buffer{}
	err = template.Execute(&body, tokens)
	if err != nil {
		logger.Error().Err(err).Msg("failed to execute sms template")
		return "", err
	}
	return body.String(), nil
}
