package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/morzhanov/go-realworld/internal/analytics/events"
	"github.com/morzhanov/go-realworld/internal/analytics/rest"
	"github.com/morzhanov/go-realworld/internal/analytics/rpc"
	"github.com/morzhanov/go-realworld/internal/analytics/services"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/helper"
	"github.com/morzhanov/go-realworld/internal/common/mq"
	"github.com/morzhanov/go-realworld/internal/common/sender"
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
	messageQ, err := mq.NewMq(c.AnalyticsKafkaTopic, 0)
	if err != nil {
		cancel()
		helper.HandleInitializationError(err, "message queue")
	}

	service := services.NewAnalyticsService(messageQ)
	rpcServer := rpc.NewAnalyticsRpcServer(service, c)
	restController := rest.NewAnalyticsRestController(service)
	eventsController := events.NewAnalyticsEventsController(service, c, sender)

	go rpcServer.Listen(ctx)
	go restController.Listen(ctx, c.AnalyticsRestPort)
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
