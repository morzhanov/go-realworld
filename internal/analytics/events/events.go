package events

import (
	"context"
	"fmt"

	anrpc "github.com/morzhanov/go-realworld/api/rpc/analytics"
	"github.com/morzhanov/go-realworld/internal/analytics/services"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/events"
	"github.com/morzhanov/go-realworld/internal/common/events/eventscontroller"
	"github.com/morzhanov/go-realworld/internal/common/sender"
)

type AnalyticsEventsController struct {
	eventscontroller.BaseEventsController
	service *services.AnalyticsService
}

func (c *AnalyticsEventsController) processRequest(in *events.EventMessage) error {
	switch in.Key {
	case "logData":
		return c.logData(in)
	case "getLogs":
		return c.getLogs(in)
	default:
		return fmt.Errorf("Wrong event name")
	}
}

func (c *AnalyticsEventsController) logData(in *events.EventMessage) error {
	res := anrpc.LogDataRequest{}
	if _, err := events.ParseEventsResponse(in.Value, &res); err != nil {
		return err
	}
	if err := c.service.LogData(&res); err != nil {
		return err
	}
	return nil
}

func (c *AnalyticsEventsController) getLogs(in *events.EventMessage) error {
	res := anrpc.GetLogRequest{}
	payload, err := events.ParseEventsResponse(in.Value, &res)
	if err != nil {
		return err
	}
	d, err := c.service.GetLog(&res)
	if err != nil {
		return err
	}
	c.BaseEventsController.SendResponse(payload.EventId, &d)
	return nil
}

func (c *AnalyticsEventsController) Listen(ctx context.Context) {
	c.BaseEventsController.Listen(
		ctx,
		func(m *events.EventMessage) { c.processRequest(m) },
	)
}

func NewAnalyticsEventsController(
	s *services.AnalyticsService,
	c *config.Config,
	sender *sender.Sender,
) *AnalyticsEventsController {
	controller := eventscontroller.NewEventsController(sender, c.AnalyticsKafkaTopic, c.KafkaUri)
	return &AnalyticsEventsController{
		service:              s,
		BaseEventsController: *controller,
	}
}
