package events

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	anrpc "github.com/morzhanov/go-realworld/api/grpc/analytics"
	"github.com/morzhanov/go-realworld/internal/analytics/services"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/events"
	"github.com/morzhanov/go-realworld/internal/common/events/eventscontroller"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
)

type analyticsEventsController struct {
	eventscontroller.BaseEventsController
	service *services.AnalyticsService
	sender  sender.Sender
}

type AnalyticsEventsController interface {
	Listen(ctx context.Context)
}

func (c *analyticsEventsController) processRequest(in *kafka.Message) error {
	switch string(in.Key) {
	case c.sender.GetAPI().Analytics.Events["logData"].Event:
		return c.logData(in)
	case c.sender.GetAPI().Analytics.Events["getLogs"].Event:
		return c.getLogs(in)
	default:
		return fmt.Errorf("wrong event name: %s", in.Key)
	}
}

func (c *analyticsEventsController) logData(in *kafka.Message) error {
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

func (c *analyticsEventsController) getLogs(in *kafka.Message) error {
	span := c.CreateSpan(in)
	defer span.Finish()

	payload, err := events.ParseEventsResponse(in.Value, &emptypb.Empty{})
	if err != nil {
		return err
	}
	d, err := c.service.GetLog(&emptypb.Empty{})
	if err != nil {
		return err
	}
	return c.BaseEventsController.SendResponse(payload.EventId, &d, &span)
}

func (c *analyticsEventsController) Listen(ctx context.Context) {
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
	sender sender.Sender,
	tracer opentracing.Tracer,
	logger *zap.Logger,
) (AnalyticsEventsController, error) {
	ec, err := eventscontroller.NewEventsController(
		sender,
		tracer,
		logger,
		c,
	)
	return &analyticsEventsController{
		service:              s,
		BaseEventsController: *ec,
		sender:               sender,
	}, err
}
