package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/morzhanov/go-realworld/internal/apigw/rest"
	"github.com/morzhanov/go-realworld/internal/apigw/services"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/events/eventslistener"
	"github.com/morzhanov/go-realworld/internal/common/helper"
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
	el := eventslistener.NewEventListener(c.ApiGatewayKafkaTopic, 0, c)

	service := services.NewAPIGatewayService(sender, el)
	restController := rest.NewAPIGatewayRestController(service)

	go restController.Listen(ctx, c.ApiGatewayRestPort)

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
