package events

import (
	"context"
	"errors"

	urpc "github.com/morzhanov/go-realworld/api/rpc/users"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/events"
	"github.com/morzhanov/go-realworld/internal/common/events/eventscontroller"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/morzhanov/go-realworld/internal/common/tracing"
	"github.com/morzhanov/go-realworld/internal/users/services"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
)

type UsersEventsController struct {
	eventscontroller.BaseEventsController
	service *services.UsersService
	tracer  *opentracing.Tracer
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
	span := tracing.StartSpanFromEventsRequest(*c.tracer, in)
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
	c.BaseEventsController.SendResponse(payload.EventId, &d, &span)
	return nil
}

func (c *UsersEventsController) getUserByUsername(in *kafka.Message) error {
	span := tracing.StartSpanFromEventsRequest(*c.tracer, in)
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
	c.BaseEventsController.SendResponse(payload.EventId, &d, &span)
	return nil
}

func (c *UsersEventsController) validatePassword(in *kafka.Message) error {
	span := tracing.StartSpanFromEventsRequest(*c.tracer, in)
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
	c.BaseEventsController.SendResponse(payload.EventId, nil, &span)
	return nil
}

func (c *UsersEventsController) createUser(in *kafka.Message) error {
	span := tracing.StartSpanFromEventsRequest(*c.tracer, in)
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
	c.BaseEventsController.SendResponse(payload.EventId, &d, &span)
	return nil
}

func (c *UsersEventsController) deleteUser(in *kafka.Message) error {
	span := tracing.StartSpanFromEventsRequest(*c.tracer, in)
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
	c.BaseEventsController.SendResponse(payload.EventId, nil, &span)
	return nil
}

func (c *UsersEventsController) Listen(ctx context.Context) {
	c.BaseEventsController.Listen(
		ctx,
		func(m *kafka.Message) { c.processRequest(m) },
	)
}

func NewUsersEventsController(
	s *services.UsersService,
	c *config.Config,
	sender *sender.Sender,
) *UsersEventsController {
	controller := eventscontroller.NewEventsController(sender, c.KafkaTopic, c.KafkaUri)
	return &UsersEventsController{
		service:              s,
		BaseEventsController: *controller,
	}
}
