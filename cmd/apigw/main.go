package main

import (
	"context"
	"github.com/morzhanov/go-realworld/internal/apigw/rest"
	"github.com/morzhanov/go-realworld/internal/apigw/services"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/events/eventslistener"
	"github.com/morzhanov/go-realworld/internal/common/helper"
	"github.com/morzhanov/go-realworld/internal/common/logger"
	"github.com/morzhanov/go-realworld/internal/common/metrics"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/morzhanov/go-realworld/internal/common/tracing"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("creating API Gateway service")
	c, err := config.NewConfig("./configs/", ".env.apigw")
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

	el := eventslistener.NewEventListener(c.KafkaTopic, 0, c, l, cancel)
	l.Info("events listener created...")

	service := services.NewAPIGatewayService(s, el)
	l.Info("service created...")
	restController := rest.NewAPIGatewayRestController(service, t, l, mc)
	l.Info("rest controller created...")

	go restController.Listen(ctx, cancel, c.RestPort)
	l.Info("all controllers started...")

	go s.Connect(c, cancel)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	l.Info("API Gateway service successfully started!")
	<-quit
	l.Info("received os.Interrupt, exiting...")
}
