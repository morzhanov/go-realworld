package main

import (
	"context"
	_ "github.com/jnewmano/grpc-json-proxy/codec"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/db"
	"github.com/morzhanov/go-realworld/internal/common/errors"
	"github.com/morzhanov/go-realworld/internal/common/logger"
	"github.com/morzhanov/go-realworld/internal/common/metrics"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/morzhanov/go-realworld/internal/common/tracing"
	"github.com/morzhanov/go-realworld/internal/users/events"
	"github.com/morzhanov/go-realworld/internal/users/grpc"
	"github.com/morzhanov/go-realworld/internal/users/rest"
	"github.com/morzhanov/go-realworld/internal/users/services"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("creating users service")
	c, err := config.NewConfig("./configs/", ".env.users")
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

	dbs, err := db.NewDb(c)
	if err != nil {
		cancel()
		errors.LogInitializationError(err, "database", l)
	}
	l.Info("database connection created...")

	service := services.NewUsersService(dbs)
	l.Info("service created...")
	rpcServer := grpc.NewUsersRpcServer(service, c, t, l)
	l.Info("grpc server created...")
	restController := rest.NewUsersRestController(service, t, l, mc)
	l.Info("rest controller created...")
	eventsController, err := events.NewUsersEventsController(service, c, s, t, l)
	if err != nil {
		cancel()
		errors.LogInitializationError(err, "events controller", l)
	}
	l.Info("events controller created...")

	go rpcServer.Listen(ctx, cancel)
	go restController.Listen(ctx, cancel, c.RestPort)
	go eventsController.Listen(ctx)
	l.Info("all controllers started...")

	if err := db.RunMigrations(dbs, "users"); err != nil {
		cancel()
		errors.LogInitializationError(err, "migrations", l)
	}
	l.Info("all database migrations applied...")

	go s.Connect(c, cancel)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	l.Info("users service successfully started!")

	<-quit
	l.Info("received os.Interrupt, exiting...")
}
