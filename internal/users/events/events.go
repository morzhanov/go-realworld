package events

import (
	"context"
	"fmt"

	urpc "github.com/morzhanov/go-realworld/api/rpc/users"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/events"
	"github.com/morzhanov/go-realworld/internal/common/events/eventscontroller"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/morzhanov/go-realworld/internal/users/services"
)

type UsersEventsController struct {
	eventscontroller.BaseEventsController
	service *services.UsersService
}

func (c *UsersEventsController) processRequest(in *events.EventMessage) error {
	switch in.Key {
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
		return fmt.Errorf("Wrong event name")
	}
}

func (c *UsersEventsController) getUser(in *events.EventMessage) error {
	res := urpc.GetUserDataRequest{}
	payload, err := events.ParseEventsResponse(in.Value, &res)
	if err != nil {
		return err
	}
	d, err := c.service.GetUserData(res.UserId)
	if err != nil {
		return err
	}
	c.BaseEventsController.SendResponse(payload.EventId, &d)
	return nil
}

func (c *UsersEventsController) getUserByUsername(in *events.EventMessage) error {
	res := urpc.GetUserDataByUsernameRequest{}
	payload, err := events.ParseEventsResponse(in.Value, &res)
	if err != nil {
		return err
	}
	d, err := c.service.GetUserData(res.Username)
	if err != nil {
		return err
	}
	c.BaseEventsController.SendResponse(payload.EventId, &d)
	return nil
}

func (c *UsersEventsController) validatePassword(in *events.EventMessage) error {
	res := urpc.ValidateUserPasswordRequest{}
	payload, err := events.ParseEventsResponse(in.Value, &res)
	if err != nil {
		return err
	}
	err = c.service.ValidateUserPassword(&res)
	if err != nil {
		return err
	}
	c.BaseEventsController.SendResponse(payload.EventId, nil)
	return nil
}

func (c *UsersEventsController) createUser(in *events.EventMessage) error {
	res := urpc.CreateUserRequest{}
	payload, err := events.ParseEventsResponse(in.Value, &res)
	if err != nil {
		return err
	}
	d, err := c.service.CreateUser(&res)
	if err != nil {
		return err
	}
	c.BaseEventsController.SendResponse(payload.EventId, &d)
	return nil
}

func (c *UsersEventsController) deleteUser(in *events.EventMessage) error {
	res := urpc.DeleteUserRequest{}
	payload, err := events.ParseEventsResponse(in.Value, &res)
	if err != nil {
		return err
	}
	err = c.service.DeleteUser(res.UserId)
	if err != nil {
		return err
	}
	c.BaseEventsController.SendResponse(payload.EventId, nil)
	return nil
}

func (c *UsersEventsController) Listen(ctx context.Context) {
	c.BaseEventsController.Listen(
		ctx,
		func(m *events.EventMessage) { c.processRequest(m) },
	)
}

func NewUsersEventsController(
	s *services.UsersService,
	c *config.Config,
	sender *sender.Sender,
) *UsersEventsController {
	controller := eventscontroller.NewEventsController(sender, c.UsersKafkaTopic, c.KafkaUri)
	return &UsersEventsController{
		service:              s,
		BaseEventsController: *controller,
	}
}
