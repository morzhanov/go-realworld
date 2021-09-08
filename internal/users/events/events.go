package events

import (
	"context"
	"errors"
	"go.uber.org/zap"

	urpc "github.com/morzhanov/go-realworld/api/rpc/users"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/events"
	"github.com/morzhanov/go-realworld/internal/common/events/eventscontroller"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/morzhanov/go-realworld/internal/users/services"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
)

type UsersEventsController struct {
	eventscontroller.BaseEventsController
	service *services.UsersService
	logger  *zap.Logger
}

func (c *UsersEventsController) processRequest(in *kafka.Message) error {
	switch string(in.Key) {
	case "getUser":
		return c.getUser(in)
	case "getUserByUsername":
		return c.getUserByUsername(in)
	case "validatePassword":
		return c.validatePassword(in)
	case "createUser":
		return c.createUser(in)
	case "deleteUser":
		return c.deleteUser(in)
	default:
		return errors.New("wrong event name")
	}
}

func (c *UsersEventsController) getUser(in *kafka.Message) error {
	span := c.CreateSpan(in)
	defer span.Finish()

	res := urpc.GetUserDataRequest{}
	payload, err := events.ParseEventsResponse(in.Value, &res)
	if err != nil {
		return err
	}
	d, err := c.service.GetUserData(res.UserId)
	if err != nil {
		return err
	}
	return c.BaseEventsController.SendResponse(payload.EventId, &d, &span)
}

func (c *UsersEventsController) getUserByUsername(in *kafka.Message) error {
	span := c.CreateSpan(in)
	defer span.Finish()

	res := urpc.GetUserDataByUsernameRequest{}
	payload, err := events.ParseEventsResponse(in.Value, &res)
	if err != nil {
		return err
	}
	d, err := c.service.GetUserData(res.Username)
	if err != nil {
		return err
	}
	return c.BaseEventsController.SendResponse(payload.EventId, &d, &span)
}

func (c *UsersEventsController) validatePassword(in *kafka.Message) error {
	span := c.CreateSpan(in)
	defer span.Finish()

	res := urpc.ValidateUserPasswordRequest{}
	payload, err := events.ParseEventsResponse(in.Value, &res)
	if err != nil {
		return err
	}
	err = c.service.ValidateUserPassword(&res)
	if err != nil {
		return err
	}
	return c.BaseEventsController.SendResponse(payload.EventId, nil, &span)
}

func (c *UsersEventsController) createUser(in *kafka.Message) error {
	span := c.CreateSpan(in)
	defer span.Finish()

	res := urpc.CreateUserRequest{}
	payload, err := events.ParseEventsResponse(in.Value, &res)
	if err != nil {
		return err
	}
	d, err := c.service.CreateUser(&res)
	if err != nil {
		return err
	}
	return c.BaseEventsController.SendResponse(payload.EventId, &d, &span)
}

func (c *UsersEventsController) deleteUser(in *kafka.Message) error {
	span := c.CreateSpan(in)
	defer span.Finish()

	res := urpc.DeleteUserRequest{}
	payload, err := events.ParseEventsResponse(in.Value, &res)
	if err != nil {
		return err
	}
	err = c.service.DeleteUser(res.UserId)
	if err != nil {
		return err
	}
	return c.BaseEventsController.SendResponse(payload.EventId, nil, &span)
}

func (c *UsersEventsController) Listen(ctx context.Context) error {
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

func NewUsersEventsController(
	s *services.UsersService,
	c *config.Config,
	sender *sender.Sender,
	tracer *opentracing.Tracer,
	logger *zap.Logger,
) (*UsersEventsController, error) {
	controller, err := eventscontroller.NewEventsController(
		sender,
		tracer,
		c.KafkaTopic,
		c.KafkaUri,
	)
	return &UsersEventsController{
		service:              s,
		BaseEventsController: *controller,
		logger:               logger,
	}, err
}
