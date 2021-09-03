package events

import (
	"context"
	"errors"

	arpc "github.com/morzhanov/go-realworld/api/rpc/auth"
	"github.com/morzhanov/go-realworld/internal/auth/services"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/events"
	"github.com/morzhanov/go-realworld/internal/common/events/eventscontroller"
	"github.com/morzhanov/go-realworld/internal/common/sender"
)

type AuthEventsController struct {
	eventscontroller.BaseEventsController
	service *services.AuthService
}

func (c *AuthEventsController) processRequest(in *events.EventMessage) error {
	switch in.Key {
	case "validateEventsRequest":
		return c.validateEventsRequest(in)
	case "login":
		return c.login(in)
	case "signup":
		return c.signup(in)
	default:
		return errors.New("wrong event name")
	}
}

func (c *AuthEventsController) validateEventsRequest(in *events.EventMessage) error {
	res := arpc.ValidateEventsRequestInput{}
	payload, err := events.ParseEventsResponse(in.Value, &res)
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

func (c *AuthEventsController) login(in *events.EventMessage) error {
	res := arpc.LoginInput{}
	payload, err := events.ParseEventsResponse(in.Value, &res)
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

func (c *AuthEventsController) signup(in *events.EventMessage) error {
	res := arpc.SignupInput{}
	payload, err := events.ParseEventsResponse(in.Value, &res)
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

func (c *AuthEventsController) Listen(ctx context.Context) {
	c.BaseEventsController.Listen(
		ctx,
		func(m *events.EventMessage) { c.processRequest(m) },
	)
}

func NewAuthEventsController(
	s *services.AuthService,
	c *config.Config,
	sender *sender.Sender,
) *AuthEventsController {
	controller := eventscontroller.NewEventsController(sender, c.AuthKafkaTopic, c.KafkaUri)
	return &AuthEventsController{
		service:              s,
		BaseEventsController: *controller,
	}
}
