package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/db"
	"github.com/morzhanov/go-realworld/internal/common/helper"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/morzhanov/go-realworld/internal/users/events"
	"github.com/morzhanov/go-realworld/internal/users/rest"
	"github.com/morzhanov/go-realworld/internal/users/rpc"
	"github.com/morzhanov/go-realworld/internal/users/services"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := config.NewConfig("../../configs/.env")
	if err != nil {
		cancel()
		helper.HandleInitializationError(err, "config")
	}
	apiConfig, err := config.NewApiConfig()
	if err != nil {
		cancel()
		helper.HandleInitializationError(err, "api config")
	}
	sender, err := sender.NewSender(c, apiConfig)
	if err != nil {
		cancel()
		helper.HandleInitializationError(err, "sender")
	}
	db, err := db.NewDb(c)
	if err != nil {
		cancel()
		helper.HandleInitializationError(err, "database")
	}

	service := services.NewUsersService(db)
	rpcServer := rpc.NewUsersRpcServer(service, c)
	restController := rest.NewUsersRestController(service)
	eventsController := events.NewUsersEventsController(service, c, sender)

	go rpcServer.Listen(ctx)
	go restController.Listen(ctx, c.UsersRestPort)
	go eventsController.Listen(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
loop:
	for {
		select {
		case <-quit:
			break loop
		default:
			time.Sleep(time.Second * 5)
		}
	}
}
