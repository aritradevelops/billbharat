package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/aritradevelops/billbharat/backend/notification/internal/config"
	"github.com/aritradevelops/billbharat/backend/notification/internal/persistence/dao"
	"github.com/aritradevelops/billbharat/backend/notification/internal/persistence/database"
	"github.com/aritradevelops/billbharat/backend/notification/internal/persistence/repository"
	"github.com/aritradevelops/billbharat/backend/shared/events"
	"github.com/aritradevelops/billbharat/backend/shared/logger"
	"github.com/aritradevelops/billbharat/backend/shared/notification"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		logger.Error().Err(err).Msg("failed to load config")
		return
	}
	db := database.NewMongoDB(conf.Database.Uri)

	if err := db.Connect(); err != nil {
		logger.Error().Err(err).Msg("failed to connect to database")
		return
	}
	defer db.Disconnect()

	repo := repository.NewRepository(db)

	eventManager := events.NewKafkaEventManager(events.KafkaOpts{
		Servers: conf.EventBroker.Servers,
		GroupId: conf.EventBroker.GroupID,
	})
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go eventManager.OnManageUserEvent(ctx, func(payload events.EventPayload[events.ManageUserEventPayload]) error {
		logger.Info().Interface("payload", payload).Msg("manage user event received")
		err := repo.SyncUser(ctx, dao.User(payload.Data))
		if err != nil {
			logger.Error().Err(err).Msg("failed to sync user")
		}
		return nil
	})

	go eventManager.OnManageNotificationEvent(ctx, func(payload events.EventPayload[events.MangageNotificationEventPayload]) error {
		logger.Info().Interface("payload", payload).Msg("manage notification event received")

		for _, channel := range payload.Data.Channels {
			switch channel.Type {
			case notification.EMAIL:
				var emailChannel notification.EmailChannel

				dbByte, _ := json.Marshal(channel.Data)
				err := json.Unmarshal(dbByte, &emailChannel)
				if err != nil {
					logger.Error().Err(err).Msg("failed to unmarshal email channel")
					return err
				}
				err = emailChannel.Send()
				if err != nil {
					logger.Error().Err(err).Msg("failed to send email")
					return err
				}
			}
		}
		return nil
	})

	<-ctx.Done()
	logger.Info().Msg("shutting down")
}
