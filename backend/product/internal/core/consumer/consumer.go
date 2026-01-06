package consumer

import (
	"context"
	"time"

	"github.com/aritradevelops/billbharat/backend/product/internal/persistence/dao"
	"github.com/aritradevelops/billbharat/backend/product/internal/persistence/repository"
	"github.com/aritradevelops/billbharat/backend/shared/events"
	"github.com/aritradevelops/billbharat/backend/shared/logger"
)

type Consumer struct {
	ctx          context.Context
	eventManager events.EventManager
	repository   repository.Repository
}

func New(ctx context.Context, eventManager events.EventManager, repository repository.Repository) *Consumer {
	return &Consumer{
		eventManager: eventManager,
		repository:   repository,
		ctx:          ctx,
	}
}

func (c *Consumer) Start() {
	c.eventManager.OnManageUserEvent(c.ctx, c.handleUserEvent)
	c.eventManager.OnManageBusinessEvent(c.ctx, c.handleBusinessEvent)
	c.eventManager.OnManageBusinessUserEvent(c.ctx, c.handleBusinessUserEvent)
}

func (c *Consumer) handleUserEvent(payload events.EventPayload[events.ManageUserEventPayload]) error {
	logger.Info().Interface("payload", payload).Msg("manage user event received")
	ctx, cancel := context.WithTimeout(c.ctx, time.Second*10)
	defer cancel()
	err := c.repository.SyncUser(ctx, dao.SyncUserParams(payload.Data))
	if err != nil {
		logger.Error().Err(err).Msg("failed to sync user")
		return err
	}
	logger.Info().Msg("user synced successfully")
	return nil
}

func (c *Consumer) handleBusinessEvent(payload events.EventPayload[events.MangageBusinessEventPayload]) error {
	logger.Info().Interface("payload", payload).Msg("manage business event received")
	ctx, cancel := context.WithTimeout(c.ctx, time.Second*10)
	defer cancel()
	err := c.repository.SyncBusiness(ctx, dao.SyncBusinessParams(payload.Data))
	if err != nil {
		logger.Error().Err(err).Msg("failed to sync business")
		return err
	}
	logger.Info().Msg("business synced successfully")
	return nil
}

func (c *Consumer) handleBusinessUserEvent(payload events.EventPayload[events.MangageBusinessUserEventPayload]) error {
	logger.Info().Interface("payload", payload).Msg("manage business user event received")
	ctx, cancel := context.WithTimeout(c.ctx, time.Second*10)
	defer cancel()
	err := c.repository.SyncBusinessUser(ctx, dao.SyncBusinessUserParams(payload.Data))
	if err != nil {
		logger.Error().Err(err).Msg("failed to sync business user")
		return err
	}
	logger.Info().Msg("business user synced successfully")
	return nil
}
