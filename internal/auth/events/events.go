package events

import (
	"context"
	"fmt"
	"go.uber.org/zap"

	arpc "github.com/morzhanov/go-realworld/api/grpc/auth"
	"github.com/morzhanov/go-realworld/internal/auth/services"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/events"
	"github.com/morzhanov/go-realworld/internal/common/events/eventscontroller"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
)

type authEventsController struct {
	eventscontroller.BaseEventsController
	service services.AuthService
	sender  sender.Sender
}

type AuthEventsController interface {
	Listen(ctx context.Context)
}

func (c *authEventsController) processRequest(in *kafka.Message) error {
	switch string(in.Key) {
	case c.sender.GetAPI().Auth().Events["validateEventsRequest"].Event:
		return c.validateEventsRequest(in)
	case c.sender.GetAPI().Auth().Events["login"].Event:
		return c.login(in)
	case c.sender.GetAPI().Auth().Events["signup"].Event:
		return c.signup(in)
	default:
		return fmt.Errorf("wrong event name: %s", in.Key)
	}
}

func (c *authEventsController) validateEventsRequest(in *kafka.Message) error {
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

func (c *authEventsController) login(in *kafka.Message) error {
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

func (c *authEventsController) signup(in *kafka.Message) error {
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

func (c *authEventsController) Listen(ctx context.Context) {
	c.BaseEventsController.Listen(
		ctx,
		func(m *kafka.Message) {
			err := c.processRequest(m)
			if err != nil {
				c.BaseEventsController.Logger().Error(err.Error())
			}
		},
	)
}

func NewAuthEventsController(
	s services.AuthService,
	c *config.Config,
	sender sender.Sender,
	tracer opentracing.Tracer,
	logger *zap.Logger,
) (AuthEventsController, error) {
	controller, err := eventscontroller.NewEventsController(
		sender,
		tracer,
		logger,
		c,
	)
	return &authEventsController{
		service:              s,
		BaseEventsController: controller,
		sender:               sender,
	}, err
}
