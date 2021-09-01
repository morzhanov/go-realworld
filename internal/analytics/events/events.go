package events

import (
	"log"

	anrpc "github.com/morzhanov/go-realworld/api/rpc/analytics"
	"github.com/morzhanov/go-realworld/internal/analytics/services"
	"github.com/morzhanov/go-realworld/internal/common/events"
	"github.com/morzhanov/go-realworld/internal/common/sender"
)

type AnalyticsEventsController struct {
	events.BaseEventsController
	service *services.AnalyticsService
}

// TODO: common logic
func check(err error) {
	// TODO: handle error
	if err != nil {
		log.Fatal(err)
	}
}

func (c *AnalyticsEventsController) Listen() {
	c.BaseEventsController.Listen(
		func(m *sender.EventMessage) { c.processRequest(m) },
	)
}

func (c *AnalyticsEventsController) processRequest(in *sender.EventMessage) {
	switch in.Key {
	case "logData":
		res := anrpc.LogDataRequest{}
		_, err := sender.ParseEventsResponse(in.Value, &res)
		check(err)
		err = c.service.LogData(&res)
		check(err)
	case "getLogs":
		res := anrpc.GetLogRequest{}
		payload, err := sender.ParseEventsResponse(in.Value, &res)
		check(err)
		d, err := c.service.GetLog(&res)
		check(err)
		c.BaseEventsController.SendResponse(payload.EventId, &d)
	}
}

func NewAnalyticsEventsController(s *services.AnalyticsService, sender *sender.Sender) *AnalyticsEventsController {
	// TODO: provide topic from config
	return &AnalyticsEventsController{
		service:              s,
		BaseEventsController: *events.NewEventsController(sender, "analytics"),
	}
}
