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
	"github.com/morzhanov/go-realworld/internal/common/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
)

type AuthEventsController struct {
	eventscontroller.BaseEventsController
	service *services.AuthService
	tracer  *opentracing.Tracer
}

func (c *AuthEventsController) processRequest(in *kafka.Message) error {
	switch string(in.Key) {
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

func (c *AuthEventsController) validateEventsRequest(in *kafka.Message) error {
	// TODO: maybe we should somehow generalize this step via middleware or something
	// TODO: because we'are using the same code for tracing in all controllers
	span := tracing.StartSpanFromEventsRequest(*c.tracer, in)
	defer span.Finish()

	res := arpc.ValidateEventsRequestInput{}
	payload, err := events.ParseEventsResponse(in.Value, &res)
	if err != nil {
		return err
	}
	d, err := c.service.ValidateEventsRequest(&res)
	if err != nil {
		return err
	}
	c.BaseEventsController.SendResponse(payload.EventId, &d, &span)
	return nil
}

func (c *AuthEventsController) login(in *kafka.Message) error {
	span := tracing.StartSpanFromEventsRequest(*c.tracer, in)
	defer span.Finish()

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
	c.BaseEventsController.SendResponse(payload.EventId, &d, &span)
	return nil
}

func (c *AuthEventsController) signup(in *kafka.Message) error {
	span := tracing.StartSpanFromEventsRequest(*c.tracer, in)
	defer span.Finish()

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
	c.BaseEventsController.SendResponse(payload.EventId, &d, &span)
	return nil
}

func (c *AuthEventsController) Listen(ctx context.Context) {
	c.BaseEventsController.Listen(
		ctx,
		func(m *kafka.Message) { c.processRequest(m) },
	)
}

func NewAuthEventsController(
	s *services.AuthService,
	c *config.Config,
	sender *sender.Sender,
	tracer *opentracing.Tracer,
) *AuthEventsController {
	controller := eventscontroller.NewEventsController(sender, c.KafkaTopic, c.KafkaUri)
	return &AuthEventsController{
		service:              s,
		BaseEventsController: *controller,
		tracer:               tracer,
	}
}
