package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/morzhanov/go-realworld/internal/analytics/events"
	"github.com/morzhanov/go-realworld/internal/analytics/rest"
	"github.com/morzhanov/go-realworld/internal/analytics/rpc"
	"github.com/morzhanov/go-realworld/internal/analytics/services"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/helper"
	"github.com/morzhanov/go-realworld/internal/common/logger"
	"github.com/morzhanov/go-realworld/internal/common/mq"
	"github.com/morzhanov/go-realworld/internal/common/sender"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := config.NewConfig("../../configs/.env.analytics")
	if err != nil {
		cancel()
		log.Fatal(err)
	}
	l, err := logger.NewLogger(c.ServiceName)
	if err != nil {
		cancel()
		log.Fatal(err)
	}

	apiConfig, err := config.NewApiConfig()
	if err != nil {
		cancel()
		helper.HandleInitializationError(err, "api config", l)
	}
	sender, err := sender.NewSender(c, apiConfig)
	if err != nil {
		cancel()
		helper.HandleInitializationError(err, "sender", l)
	}
	messageQ, err := mq.NewMq(c.KafkaTopic, 0)
	if err != nil {
		cancel()
		helper.HandleInitializationError(err, "message queue", l)
	}

	service := services.NewAnalyticsService(messageQ)
	rpcServer := rpc.NewAnalyticsRpcServer(service, c)
	restController := rest.NewAnalyticsRestController(service)
	eventsController := events.NewAnalyticsEventsController(service, c, sender)

	go rpcServer.Listen(ctx)
	go restController.Listen(ctx, c.RestPort)
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
