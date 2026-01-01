package consumer

import (
	"context"
	"time"

	"github.com/aritradevelops/billbharat/backend/notification/internal/core/notifier"
	"github.com/aritradevelops/billbharat/backend/notification/internal/persistence/dao"
	"github.com/aritradevelops/billbharat/backend/notification/internal/persistence/repository"
	"github.com/aritradevelops/billbharat/backend/shared/events"
	"github.com/aritradevelops/billbharat/backend/shared/logger"
)

type Consumer struct {
	ctx          context.Context
	eventManager events.EventManager
	repository   repository.Repository
	notifier     notifier.Notifier
}

func New(ctx context.Context, eventManager events.EventManager, repository repository.Repository, notifier notifier.Notifier) *Consumer {
	return &Consumer{
		eventManager: eventManager,
		repository:   repository,
		ctx:          ctx,
		notifier:     notifier,
	}
}

func (c *Consumer) Start() {
	c.eventManager.OnManageNotificationEvent(c.ctx, c.handleNotificationEvent)
	c.eventManager.OnManageUserEvent(c.ctx, c.handleUserEvent)
	c.eventManager.OnManageBusinessEvent(c.ctx, c.handleBusinessEvent)
	c.eventManager.OnManageBusinessUserEvent(c.ctx, c.handleBusinessUserEvent)
}

func (c *Consumer) handleNotificationEvent(payload events.EventPayload[events.MangageNotificationEventPayload]) error {
	logger.Info().Interface("payload", payload).Msg("manage notification event received")
	ctx, cancel := context.WithTimeout(c.ctx, time.Second*10)
	defer cancel()
	err := c.notifier.Notify(ctx, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to notify")
		return err
	}
	return nil
}

func (c *Consumer) handleUserEvent(payload events.EventPayload[events.ManageUserEventPayload]) error {
	logger.Info().Interface("payload", payload).Msg("manage user event received")
	ctx, cancel := context.WithTimeout(c.ctx, time.Second*10)
	defer cancel()
	err := c.repository.SyncUser(ctx, dao.User(payload.Data))
	if err != nil {
		logger.Error().Err(err).Msg("failed to sync user")
		return err
	}
	return nil
}

func (c *Consumer) handleBusinessEvent(payload events.EventPayload[events.MangageBusinessEventPayload]) error {
	logger.Info().Interface("payload", payload).Msg("manage business event received")
	ctx, cancel := context.WithTimeout(c.ctx, time.Second*10)
	defer cancel()
	err := c.repository.SyncBusiness(ctx, dao.Business(payload.Data))
	if err != nil {
		logger.Error().Err(err).Msg("failed to sync business")
		return err
	}
	return nil
}

func (c *Consumer) handleBusinessUserEvent(payload events.EventPayload[events.MangageBusinessUserEventPayload]) error {
	logger.Info().Interface("payload", payload).Msg("manage business user event received")
	ctx, cancel := context.WithTimeout(c.ctx, time.Second*10)
	defer cancel()
	err := c.repository.SyncBusinessUser(ctx, dao.BusinessUser(payload.Data))
	if err != nil {
		logger.Error().Err(err).Msg("failed to sync business user")
		return err
	}
	return nil
}
