package main

import (
	"fmt"

	"github.com/aritradevelops/billbharat/backend/auth/internal/config"
	"github.com/aritradevelops/billbharat/backend/auth/internal/core/jwtutil"
	"github.com/aritradevelops/billbharat/backend/auth/internal/core/service"
	"github.com/aritradevelops/billbharat/backend/auth/internal/persistence/database"
	"github.com/aritradevelops/billbharat/backend/auth/internal/persistence/repository"
	"github.com/aritradevelops/billbharat/backend/auth/internal/ports/httpd"
	"github.com/aritradevelops/billbharat/backend/auth/internal/ports/httpd/handlers"
	"github.com/aritradevelops/billbharat/backend/shared/eventbroker"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		fmt.Println("failed to load config", err)
		return
	}

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
	jwtManager := jwtutil.NewJwtManager(conf.Jwt.Secret, conf.Jwt.Lifetime.Duration())

	eventBroker := eventbroker.NewKafkaProducer(conf.EventBroker.Servers)

	srv := service.New(repo, jwtManager, eventBroker)

	handler := handlers.New(db, srv, eventBroker)

	server := httpd.NewServer(conf.Http.Host, conf.Http.Port, handler, jwtManager)
	server.SetupRoutes()

	if err := server.Start(); err != nil {
		fmt.Println("server failed to start", err)
		return
	}
}
