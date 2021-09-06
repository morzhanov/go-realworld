package events

import (
	"context"
	"errors"

	anrpc "github.com/morzhanov/go-realworld/api/rpc/analytics"
	"github.com/morzhanov/go-realworld/internal/analytics/services"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/events"
	"github.com/morzhanov/go-realworld/internal/common/events/eventscontroller"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/morzhanov/go-realworld/internal/common/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
)

type AnalyticsEventsController struct {
	eventscontroller.BaseEventsController
	service *services.AnalyticsService
	tracer  *opentracing.Tracer
}

func (c *AnalyticsEventsController) processRequest(in *kafka.Message) error {
	switch string(in.Key) {
	case "logData":
		return c.logData(in)
	case "getLogs":
		return c.getLogs(in)
	default:
		return errors.New("wrong event name")
	}
}

func (c *AnalyticsEventsController) logData(in *kafka.Message) error {
	span := tracing.StartSpanFromEventsRequest(*c.tracer, in)
	defer span.Finish()

	res := anrpc.LogDataRequest{}
	if _, err := events.ParseEventsResponse(in.Value, &res); err != nil {
		return err
	}
	if err := c.service.LogData(&res); err != nil {
		return err
	}
	return nil
}

func (c *AnalyticsEventsController) getLogs(in *kafka.Message) error {
	span := tracing.StartSpanFromEventsRequest(*c.tracer, in)
	defer span.Finish()

	res := anrpc.GetLogRequest{}
	payload, err := events.ParseEventsResponse(in.Value, &res)
	if err != nil {
		return err
	}
	d, err := c.service.GetLog(&res)
	if err != nil {
		return err
	}
	c.BaseEventsController.SendResponse(payload.EventId, &d, &span)
	return nil
}

func (c *AnalyticsEventsController) Listen(ctx context.Context) {
	c.BaseEventsController.Listen(
		ctx,
		func(m *kafka.Message) { c.processRequest(m) },
	)
}

func NewAnalyticsEventsController(
	s *services.AnalyticsService,
	c *config.Config,
	sender *sender.Sender,
) *AnalyticsEventsController {
	controller := eventscontroller.NewEventsController(sender, c.KafkaTopic, c.KafkaUri)
	return &AnalyticsEventsController{
		service:              s,
		BaseEventsController: *controller,
	}
}
