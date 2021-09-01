package events

import (
	"context"
	"log"

	arpc "github.com/morzhanov/go-realworld/api/rpc/auth"
	"github.com/morzhanov/go-realworld/internal/auth/services"
	"github.com/morzhanov/go-realworld/internal/common/events"
	"github.com/morzhanov/go-realworld/internal/common/sender"
)

type AuthEventsController struct {
	events.BaseEventsController
	service *services.AuthService
}

// TODO: common logic
func check(err error) {
	// TODO: handle error
	if err != nil {
		log.Fatal(err)
	}
}

func (c *AuthEventsController) Listen() {
	c.BaseEventsController.Listen(
		func(m *sender.EventMessage) { c.processRequest(m) },
	)
}

func (c *AuthEventsController) processRequest(in *sender.EventMessage) {
	switch in.Key {
	case "validateEventsRequest":
		res := arpc.ValidateEventsRequestInput{}
		payload, err := sender.ParseEventsResponse(in.Value, &res)
		check(err)
		d, err := c.service.ValidateEventsRequest(&res)
		check(err)
		c.BaseEventsController.SendResponse(payload.EventId, &d)
	case "login":
		res := arpc.LoginInput{}
		payload, err := sender.ParseEventsResponse(in.Value, &res)
		check(err)

		ctx := context.WithValue(context.Background(), "transport", sender.EventsTransport)
		d, err := c.service.Login(ctx, &res)
		check(err)
		c.BaseEventsController.SendResponse(payload.EventId, &d)
	case "signup":
		res := arpc.SignupInput{}
		payload, err := sender.ParseEventsResponse(in.Value, &res)
		check(err)

		ctx := context.WithValue(context.Background(), "transport", sender.EventsTransport)
		d, err := c.service.Signup(ctx, &res)
		check(err)
		c.BaseEventsController.SendResponse(payload.EventId, &d)
	}
}

func NewAuthEventsController(s *services.AuthService, sender *sender.Sender) *AuthEventsController {
	// TODO: provide topic from config
	return &AuthEventsController{
		service:              s,
		BaseEventsController: *events.NewEventsController(sender, "auth"),
	}
}
