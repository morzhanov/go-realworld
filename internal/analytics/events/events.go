package events

import (
	"context"
	"errors"
	"go.uber.org/zap"

	anrpc "github.com/morzhanov/go-realworld/api/rpc/analytics"
	"github.com/morzhanov/go-realworld/internal/analytics/services"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/events"
	"github.com/morzhanov/go-realworld/internal/common/events/eventscontroller"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
)

type AnalyticsEventsController struct {
	eventscontroller.BaseEventsController
	service *services.AnalyticsService
	logger  *zap.Logger
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
	span := c.CreateSpan(in)
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
	span := c.CreateSpan(in)
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
	return c.BaseEventsController.SendResponse(payload.EventId, &d, &span)
}

func (c *AnalyticsEventsController) Listen(ctx context.Context) error {
	return c.BaseEventsController.Listen(
		ctx,
		func(m *kafka.Message) {
			err := c.processRequest(m)
			if err != nil {
				c.logger.Error(err.Error())
			}
		},
	)
}

func NewAnalyticsEventsController(
	s *services.AnalyticsService,
	c *config.Config,
	sender *sender.Sender,
	tracer *opentracing.Tracer,
	logger *zap.Logger,
) (*AnalyticsEventsController, error) {
	controller, err := eventscontroller.NewEventsController(
		sender,
		tracer,
		c.KafkaTopic,
		c.KafkaUri,
	)
	return &AnalyticsEventsController{
		service:              s,
		BaseEventsController: *controller,
		logger:               logger,
	}, err
}
