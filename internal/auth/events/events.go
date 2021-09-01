package events

import (
	"context"
	"fmt"

	arpc "github.com/morzhanov/go-realworld/api/rpc/auth"
	"github.com/morzhanov/go-realworld/internal/auth/services"
	"github.com/morzhanov/go-realworld/internal/common/events"
	"github.com/morzhanov/go-realworld/internal/common/sender"
)

type AuthEventsController struct {
	events.BaseEventsController
	service *services.AuthService
}

func (c *AuthEventsController) Listen() {
	c.BaseEventsController.Listen(
		func(m *sender.EventMessage) { c.processRequest(m) },
	)
}

func (c *AuthEventsController) processRequest(in *sender.EventMessage) error {
	switch in.Key {
	case "validateEventsRequest":
		return c.validateEventsRequest(in)
	case "login":
		return c.login(in)
	case "signup":
		return c.signup(in)
	default:
		return fmt.Errorf("Wrong event name")
	}
}

func (c *AuthEventsController) validateEventsRequest(in *sender.EventMessage) error {
	res := arpc.ValidateEventsRequestInput{}
	payload, err := sender.ParseEventsResponse(in.Value, &res)
	if err != nil {
		return err
	}
	d, err := c.service.ValidateEventsRequest(&res)
	if err != nil {
		return err
	}
	c.BaseEventsController.SendResponse(payload.EventId, &d)
	return nil
}

func (c *AuthEventsController) login(in *sender.EventMessage) error {
	res := arpc.LoginInput{}
	payload, err := sender.ParseEventsResponse(in.Value, &res)
	if err != nil {
		return err
	}
	ctx := context.WithValue(context.Background(), "transport", sender.EventsTransport)
	d, err := c.service.Login(ctx, &res)
	if err != nil {
		return err
	}
	c.BaseEventsController.SendResponse(payload.EventId, &d)
	return nil
}

func (c *AuthEventsController) signup(in *sender.EventMessage) error {
	res := arpc.SignupInput{}
	payload, err := sender.ParseEventsResponse(in.Value, &res)
	if err != nil {
		return err
	}

	ctx := context.WithValue(context.Background(), "transport", sender.EventsTransport)
	d, err := c.service.Signup(ctx, &res)
	if err != nil {
		return err
	}
	c.BaseEventsController.SendResponse(payload.EventId, &d)
	return nil
}

func NewAuthEventsController(s *services.AuthService, sender *sender.Sender) *AuthEventsController {
	// TODO: provide topic from config
	return &AuthEventsController{
		service:              s,
		BaseEventsController: *events.NewEventsController(sender, "auth"),
	}
}
