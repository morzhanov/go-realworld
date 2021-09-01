package events

import (
	"fmt"

	urpc "github.com/morzhanov/go-realworld/api/rpc/users"
	"github.com/morzhanov/go-realworld/internal/common/events"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/morzhanov/go-realworld/internal/users/services"
)

type UsersEventsController struct {
	events.BaseEventsController
	service *services.UsersService
}

func (c *UsersEventsController) Listen() {
	c.BaseEventsController.Listen(
		func(m *sender.EventMessage) { c.processRequest(m) },
	)
}

func (c *UsersEventsController) processRequest(in *sender.EventMessage) error {
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

func (c *UsersEventsController) getUser(in *sender.EventMessage) error {
	res := urpc.GetUserDataRequest{}
	payload, err := sender.ParseEventsResponse(in.Value, &res)
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

func (c *UsersEventsController) getUserByUsername(in *sender.EventMessage) error {
	res := urpc.GetUserDataByUsernameRequest{}
	payload, err := sender.ParseEventsResponse(in.Value, &res)
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

func (c *UsersEventsController) validatePassword(in *sender.EventMessage) error {
	res := urpc.ValidateUserPasswordRequest{}
	payload, err := sender.ParseEventsResponse(in.Value, &res)
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

func (c *UsersEventsController) createUser(in *sender.EventMessage) error {
	res := urpc.CreateUserRequest{}
	payload, err := sender.ParseEventsResponse(in.Value, &res)
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

func (c *UsersEventsController) deleteUser(in *sender.EventMessage) error {
	res := urpc.DeleteUserRequest{}
	payload, err := sender.ParseEventsResponse(in.Value, &res)
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

func NewUsersEventsController(s *services.UsersService, sender *sender.Sender) *UsersEventsController {
	// TODO: provide topic from config
	return &UsersEventsController{
		service:              s,
		BaseEventsController: *events.NewEventsController(sender, "users"),
	}
}
