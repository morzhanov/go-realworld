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
	"github.com/morzhanov/go-realworld/internal/common/metrics"
	"github.com/morzhanov/go-realworld/internal/common/mq"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/morzhanov/go-realworld/internal/common/tracing"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("creating analytics service")
	c, err := config.NewConfig("./configs/", ".env.analytics")
	if err != nil {
		cancel()
		log.Fatal(err)
	}
	log.Println("config parsed...")
	l, err := logger.NewLogger(c.ServiceName)
	if err != nil {
		cancel()
		log.Fatal(err)
	}
	l.Info("logger created...")
	t, err := tracing.NewTracer(ctx, c, l)
	if err != nil {
		cancel()
		log.Fatal(err)
	}
	l.Info("tracer created...")

	mc := metrics.NewMetricsCollector(c)
	mc.RecordBaseMetrics(ctx)
	l.Info("metrics collector created...")

	apiConfig, err := config.NewApiConfig()
	if err != nil {
		cancel()
		helper.HandleInitializationError(err, "api config", l)
	}
	l.Info("apiConfig created...")

	s := sender.NewSender(apiConfig, l)
	l.Info("sender created...")

	messageQ, err := mq.NewMq(c.KafkaTopic, 0)
	if err != nil {
		cancel()
		helper.HandleInitializationError(err, "message queue", l)
	}
	l.Info("message queue created...")

	service := services.NewAnalyticsService(messageQ)
	l.Info("service created...")
	rpcServer := rpc.NewAnalyticsRpcServer(service, c, t, l)
	l.Info("grpc server created...")
	restController := rest.NewAnalyticsRestController(service, t, l, mc)
	l.Info("rest controller created...")
	eventsController, err := events.NewAnalyticsEventsController(service, c, s, t, l)
	if err != nil {
		cancel()
		helper.HandleInitializationError(err, "events controller", l)
	}
	l.Info("events controller created...")

	go rpcServer.Listen(ctx, cancel)
	go restController.Listen(ctx, cancel, c.RestPort)
	go eventsController.Listen(ctx, cancel)
	l.Info("all controllers started...")

	go s.Connect(c, cancel)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	l.Info("analytics service successfully started!")
loop:
	for {
		select {
		case <-quit:
			l.Info("received os.Interrupt, exiting...")
			break loop
		default:
			time.Sleep(time.Second * 5)
		}
	}
}
