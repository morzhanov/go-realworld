package events

import (
	"context"
	"fmt"
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
	sender  *sender.Sender
}

func (c *AnalyticsEventsController) processRequest(in *kafka.Message) error {
	switch string(in.Key) {
	case c.sender.API.Analytics.Events["logData"].Event:
		return c.logData(in)
	case c.sender.API.Analytics.Events["getLogs"].Event:
		return c.getLogs(in)
	default:
		return fmt.Errorf("wrong event name: %s", in.Key)
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

func (c *AnalyticsEventsController) Listen(ctx context.Context) {
	c.BaseEventsController.Listen(
		ctx,
		func(m *kafka.Message) {
			err := c.processRequest(m)
			if err != nil {
				c.BaseEventsController.Logger.Error(err.Error())
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
		logger,
	)
	return &AnalyticsEventsController{
		service:              s,
		BaseEventsController: *controller,
		sender:               sender,
	}, err
}
