package main

import (
	"context"
	"fmt"

	"github.com/aritradevelops/billbharat/backend/notification/internal/config"
	"github.com/aritradevelops/billbharat/backend/notification/internal/persistence/dao"
	"github.com/aritradevelops/billbharat/backend/notification/internal/persistence/database"
	"github.com/aritradevelops/billbharat/backend/notification/internal/persistence/repository"
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

	var repo repository.Repository

	if conf.Deployment.Env == "local" {
		repo = repository.NewFilesystemRepository()
	} else {
		repo = repository.NewRepository(db)
	}

	temp, err := repo.FindTemplate(context.Background(), repository.FindTemplateParams{
		Event:    dao.USER_SIGNUP_OTP,
		Channel:  dao.EMAIL,
		Locale:   "en",
		Scope:    "default",
		Mimetype: "text/html",
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to find template")
		return
	}
	logger.Info().Interface("template", temp).Msg("template found")
	fmt.Println(temp.Body)

}
