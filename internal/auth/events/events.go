package events

import (
	"context"
	"errors"
	"go.uber.org/zap"

	arpc "github.com/morzhanov/go-realworld/api/rpc/auth"
	"github.com/morzhanov/go-realworld/internal/auth/services"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/events"
	"github.com/morzhanov/go-realworld/internal/common/events/eventscontroller"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
)

type AuthEventsController struct {
	eventscontroller.BaseEventsController
	service *services.AuthService
	logger  *zap.Logger
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
	span := c.CreateSpan(in)
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
	return c.BaseEventsController.SendResponse(payload.EventId, &d, &span)
}

func (c *AuthEventsController) login(in *kafka.Message) error {
	span := c.CreateSpan(in)
	defer span.Finish()

	res := arpc.LoginInput{}
	payload, err := events.ParseEventsResponse(in.Value, &res)
	if err != nil {
		return err
	}
	ctx := context.WithValue(context.Background(), "transport", sender.EventsTransport)
	d, err := c.service.Login(ctx, &res, &span)
	if err != nil {
		return err
	}
	return c.BaseEventsController.SendResponse(payload.EventId, &d, &span)
}

func (c *AuthEventsController) signup(in *kafka.Message) error {
	span := c.CreateSpan(in)
	defer span.Finish()

	res := arpc.SignupInput{}
	payload, err := events.ParseEventsResponse(in.Value, &res)
	if err != nil {
		return err
	}

	ctx := context.WithValue(context.Background(), "transport", sender.EventsTransport)
	d, err := c.service.Signup(ctx, &res, &span)
	if err != nil {
		return err
	}
	return c.BaseEventsController.SendResponse(payload.EventId, &d, &span)
}

func (c *AuthEventsController) Listen(ctx context.Context) error {
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

func NewAuthEventsController(
	s *services.AuthService,
	c *config.Config,
	sender *sender.Sender,
	tracer *opentracing.Tracer,
	logger *zap.Logger,
) (*AuthEventsController, error) {
	controller, err := eventscontroller.NewEventsController(
		sender,
		tracer,
		c.KafkaTopic,
		c.KafkaUri,
	)
	return &AuthEventsController{
		service:              s,
		BaseEventsController: *controller,
		logger:               logger,
	}, err
}
