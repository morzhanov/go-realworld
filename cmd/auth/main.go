package main

import (
	"context"
	_ "github.com/jnewmano/grpc-json-proxy/codec"
	"github.com/morzhanov/go-realworld/internal/auth/events"
	"github.com/morzhanov/go-realworld/internal/auth/grpc"
	"github.com/morzhanov/go-realworld/internal/auth/rest"
	"github.com/morzhanov/go-realworld/internal/auth/services"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/errors"
	"github.com/morzhanov/go-realworld/internal/common/events/eventslistener"
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

	log.Println("creating auth service")
	c, err := config.NewConfig("./configs/", ".env.auth")
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
		errors.LogInitializationError(err, "api config", l)
	}
	l.Info("apiConfig created...")

	s := sender.NewSender(apiConfig, l)
	l.Info("sender created...")

	el := eventslistener.NewEventListener(ctx, c, l)
	l.Info("events listener created...")

	service := services.NewAuthService(s, el, c)
	l.Info("service created...")
	rpcServer := grpc.NewAuthRpcServer(service, c, t, l)
	l.Info("grpc server created...")
	restController := rest.NewAuthRestController(service, t, l, mc)
	l.Info("rest controller created...")
	eventsController, err := events.NewAuthEventsController(service, c, s, t, l)
	if err != nil {
		cancel()
		errors.LogInitializationError(err, "events controller", l)
	}
	l.Info("events controller created...")

	go s.Connect(c, cancel)

	go rpcServer.Listen(ctx, cancel)
	go restController.Listen(ctx, cancel, c.RestPort)
	go eventsController.Listen(ctx)
	l.Info("all controllers started...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	l.Info("auth service successfully started!")
	<-quit
	l.Info("received os.Interrupt, exiting...")
}
