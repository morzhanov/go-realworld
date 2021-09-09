package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/morzhanov/go-realworld/internal/auth/events"
	"github.com/morzhanov/go-realworld/internal/auth/rest"
	"github.com/morzhanov/go-realworld/internal/auth/rpc"
	"github.com/morzhanov/go-realworld/internal/auth/services"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/db"
	"github.com/morzhanov/go-realworld/internal/common/events/eventslistener"
	"github.com/morzhanov/go-realworld/internal/common/helper"
	"github.com/morzhanov/go-realworld/internal/common/logger"
	"github.com/morzhanov/go-realworld/internal/common/metrics"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/morzhanov/go-realworld/internal/common/tracing"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := config.NewConfig("./configs/", ".env.auth")
	if err != nil {
		cancel()
		log.Fatal(err)
	}
	l, err := logger.NewLogger(c.ServiceName)
	if err != nil {
		cancel()
		log.Fatal(err)
	}
	t, err := tracing.NewTracer(ctx, c, l)
	if err != nil {
		cancel()
		log.Fatal(err)
	}

	mc := metrics.NewMetricsCollector(c)
	mc.RecordBaseMetrics(ctx)

	apiConfig, err := config.NewApiConfig()
	if err != nil {
		cancel()
		helper.HandleInitializationError(err, "api config", l)
	}
	s, err := sender.NewSender(c, apiConfig)
	if err != nil {
		cancel()
		helper.HandleInitializationError(err, "sender", l)
	}
	dbs, err := db.NewDb(c)
	if err != nil {
		cancel()
		helper.HandleInitializationError(err, "database", l)
	}
	el := eventslistener.NewEventListener(c.KafkaTopic, 0, c, l)

	service := services.NewAuthService(dbs, s, el, c)
	rpcServer := rpc.NewAuthRpcServer(service, c, t, l)
	restController := rest.NewAuthRestController(service, t, l, mc)
	eventsController, err := events.NewAuthEventsController(service, c, s, t, l)
	if err != nil {
		cancel()
		helper.HandleInitializationError(err, "events controller", l)
	}

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
