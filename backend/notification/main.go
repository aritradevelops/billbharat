package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/aritradevelops/billbharat/backend/notification/internal/config"
	"github.com/aritradevelops/billbharat/backend/notification/internal/core/consumer"
	"github.com/aritradevelops/billbharat/backend/notification/internal/core/notifier"
	"github.com/aritradevelops/billbharat/backend/notification/internal/persistence/database"
	"github.com/aritradevelops/billbharat/backend/notification/internal/persistence/repository"
	"github.com/aritradevelops/billbharat/backend/shared/events"
	"github.com/aritradevelops/billbharat/backend/shared/logger"
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
	notifier := notifier.New(ctx, conf.Deployment.Env, repo)
	consumer := consumer.New(ctx, eventManager, repo, notifier)
	consumer.Start()

	<-ctx.Done()
	logger.Info().Msg("shutting down")
}
