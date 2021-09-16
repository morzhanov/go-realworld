package events

import (
	"context"
	"fmt"
	"go.uber.org/zap"

	urpc "github.com/morzhanov/go-realworld/api/grpc/users"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/events"
	"github.com/morzhanov/go-realworld/internal/common/events/eventscontroller"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/morzhanov/go-realworld/internal/users/services"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
)

type usersEventsController struct {
	eventscontroller.BaseEventsController
	service services.UsersService
	sender  sender.Sender
}

type UsersEventsController interface {
	Listen(ctx context.Context)
}

func (c *usersEventsController) processRequest(in *kafka.Message) error {
	switch string(in.Key) {
	case c.sender.GetAPI().Users.Events["getUser"].Event:
		return c.getUser(in)
	case c.sender.GetAPI().Users.Events["getUserByUsername"].Event:
		return c.getUserByUsername(in)
	case c.sender.GetAPI().Users.Events["validatePassword"].Event:
		return c.validatePassword(in)
	case c.sender.GetAPI().Users.Events["createUser"].Event:
		return c.createUser(in)
	case c.sender.GetAPI().Users.Events["deleteUser"].Event:
		return c.deleteUser(in)
	default:
		return fmt.Errorf("wrong event name: %s", in.Key)
	}
}

func (c *usersEventsController) getUser(in *kafka.Message) error {
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

func (c *usersEventsController) getUserByUsername(in *kafka.Message) error {
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

func (c *usersEventsController) validatePassword(in *kafka.Message) error {
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

func (c *usersEventsController) createUser(in *kafka.Message) error {
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

func (c *usersEventsController) deleteUser(in *kafka.Message) error {
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

func (c *usersEventsController) Listen(ctx context.Context) {
	c.BaseEventsController.Listen(
		ctx,
		func(m *kafka.Message) {
			err := c.processRequest(m)
			if err != nil {
				c.BaseEventsController.Logger.Error(err.Error())
			}
		},
	)
}

func NewUsersEventsController(
	s services.UsersService,
	c *config.Config,
	sender sender.Sender,
	tracer opentracing.Tracer,
	logger *zap.Logger,
) (UsersEventsController, error) {
	controller, err := eventscontroller.NewEventsController(
		sender,
		tracer,
		logger,
		c,
	)
	return &usersEventsController{
		service:              s,
		BaseEventsController: *controller,
		sender:               sender,
	}, err
}
