package main

import (
	"context"
	"fmt"

	"github.com/aritradevelops/billbharat/backend/product/internal/config"
	"github.com/aritradevelops/billbharat/backend/product/internal/core/consumer"
	"github.com/aritradevelops/billbharat/backend/product/internal/core/jwtutil"
	"github.com/aritradevelops/billbharat/backend/product/internal/core/service"
	"github.com/aritradevelops/billbharat/backend/product/internal/persistence/database"
	"github.com/aritradevelops/billbharat/backend/product/internal/persistence/repository"
	"github.com/aritradevelops/billbharat/backend/product/internal/ports/httpd"
	"github.com/aritradevelops/billbharat/backend/product/internal/ports/httpd/handlers"
	"github.com/aritradevelops/billbharat/backend/shared/events"
	figure "github.com/common-nighthawk/go-figure"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		fmt.Println("failed to load config", err)
		return
	}

	logo := figure.NewColorFigure(conf.Service.Name, "", "blue", true)
	logo.Print()

	db := database.NewPostgres(conf.Database.Uri, conf.Database.Timeout)

	err = db.Connect()
	if err != nil {
		fmt.Println("failed to connect to database", err)
		return
	}
	defer db.Disconnect()

	repo, err := repository.New(db)
	if err != nil {
		fmt.Println("failed to create repository", err)
		return
	}

	eventManager := events.NewKafkaEventManager(events.KafkaOpts{
		Servers: conf.EventBroker.Servers,
		GroupId: conf.EventBroker.GroupID,
	})

	jwtManager := jwtutil.NewJwtManager(conf.Jwt.Secret, conf.Jwt.Lifetime.Duration())

	srv := service.New(repo)

	handler := handlers.New(db, srv, conf.Deployment.Env)

	server := httpd.NewServer(conf.Http.Host, conf.Http.Port, handler, jwtManager)
	server.SetupRoutes()

	consumer := consumer.New(context.Background(), eventManager, repo)
	consumer.Start()

	if err := server.Start(); err != nil {
		fmt.Println("server failed to start", err)
		return
	}
}
