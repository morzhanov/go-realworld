package events

import (
	"fmt"

	anrpc "github.com/morzhanov/go-realworld/api/rpc/analytics"
	"github.com/morzhanov/go-realworld/internal/analytics/services"
	"github.com/morzhanov/go-realworld/internal/common/events"
	"github.com/morzhanov/go-realworld/internal/common/sender"
)

type AnalyticsEventsController struct {
	events.BaseEventsController
	service *services.AnalyticsService
}

func (c *AnalyticsEventsController) Listen() {
	c.BaseEventsController.Listen(
		func(m *sender.EventMessage) { c.processRequest(m) },
	)
}

func (c *AnalyticsEventsController) processRequest(in *sender.EventMessage) error {
	switch in.Key {
	case "logData":
		return c.logData(in)
	case "getLogs":
		return c.getLogs(in)
	default:
		return fmt.Errorf("Wrong event name")
	}
}

func (c *AnalyticsEventsController) logData(in *sender.EventMessage) error {
	res := anrpc.LogDataRequest{}
	if _, err := sender.ParseEventsResponse(in.Value, &res); err != nil {
		return err
	}
	if err := c.service.LogData(&res); err != nil {
		return err
	}
	return nil
}

func (c *AnalyticsEventsController) getLogs(in *sender.EventMessage) error {
	res := anrpc.GetLogRequest{}
	payload, err := sender.ParseEventsResponse(in.Value, &res)
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

func NewAnalyticsEventsController(s *services.AnalyticsService, sender *sender.Sender) *AnalyticsEventsController {
	// TODO: provide topic from config
	return &AnalyticsEventsController{
		service:              s,
		BaseEventsController: *events.NewEventsController(sender, "analytics"),
	}
}
