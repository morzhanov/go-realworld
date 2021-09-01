package events

import (
	"log"

	urpc "github.com/morzhanov/go-realworld/api/rpc/users"
	"github.com/morzhanov/go-realworld/internal/common/events"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/morzhanov/go-realworld/internal/users/services"
)

type UsersEventsController struct {
	events.BaseEventsController
	service *services.UsersService
}

func check(err error) {
	// TODO: handle error
	if err != nil {
		log.Fatal(err)
	}
}

func (c *UsersEventsController) Listen() {
	c.BaseEventsController.Listen(
		func(m *sender.EventMessage) { c.processRequest(m) },
	)
}

func (c *UsersEventsController) processRequest(in *sender.EventMessage) {
	switch in.Key {
	case "getUser":
		res := urpc.GetUserDataRequest{}
		payload, err := sender.ParseEventsResponse(in.Value, &res)
		check(err)
		d, err := c.service.GetUserData(res.UserId)
		check(err)
		c.BaseEventsController.SendResponse(payload.EventId, &d)
	case "getUserByUsername":
		res := urpc.GetUserDataByUsernameRequest{}
		payload, err := sender.ParseEventsResponse(in.Value, &res)
		check(err)
		d, err := c.service.GetUserData(res.Username)
		check(err)
		c.BaseEventsController.SendResponse(payload.EventId, &d)
	case "validatePassword":
		res := urpc.ValidateUserPasswordRequest{}
		payload, err := sender.ParseEventsResponse(in.Value, &res)
		check(err)
		err = c.service.ValidateUserPassword(&res)
		check(err)
		c.BaseEventsController.SendResponse(payload.EventId, nil)
	case "createUser":
		res := urpc.CreateUserRequest{}
		payload, err := sender.ParseEventsResponse(in.Value, &res)
		check(err)
		d, err := c.service.CreateUser(&res)
		check(err)
		c.BaseEventsController.SendResponse(payload.EventId, &d)
	case "deleteUser":
		res := urpc.DeleteUserRequest{}
		payload, err := sender.ParseEventsResponse(in.Value, &res)
		check(err)
		err = c.service.DeleteUser(res.UserId)
		check(err)
		c.BaseEventsController.SendResponse(payload.EventId, nil)
	}
}

func NewUsersEventsController(s *services.UsersService, sender *sender.Sender) *UsersEventsController {
	// TODO: provide topic from config
	return &UsersEventsController{
		service:              s,
		BaseEventsController: *events.NewEventsController(sender, "users"),
	}
}
