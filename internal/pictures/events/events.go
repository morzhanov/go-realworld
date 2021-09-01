package events

import (
	"log"

	prpc "github.com/morzhanov/go-realworld/api/rpc/pictures"
	"github.com/morzhanov/go-realworld/internal/common/events"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/morzhanov/go-realworld/internal/pictures/services"
)

type PicturesEventsController struct {
	events.BaseEventsController
	service *services.PictureService
}

func check(err error) {
	// TODO: handle error
	if err != nil {
		log.Fatal(err)
	}
}

func (c *PicturesEventsController) Listen() {
	c.BaseEventsController.Listen(
		func(m *sender.EventMessage) { c.processRequest(m) },
	)
}

func (c *PicturesEventsController) processRequest(in *sender.EventMessage) {
	switch in.Key {
	case "getPictures":
		res := prpc.GetUserPicturesRequest{}
		payload, err := sender.ParseEventsResponse(in.Value, &res)
		check(err)
		d, err := c.service.GetUserPictures(res.UserId)
		check(err)
		c.BaseEventsController.SendResponse(payload.EventId, &d)
	case "getPicture":
		res := prpc.GetUserPictureRequest{}
		payload, err := sender.ParseEventsResponse(in.Value, &res)
		check(err)
		d, err := c.service.GetUserPicture(res.UserId, res.PictureId)
		check(err)
		c.BaseEventsController.SendResponse(payload.EventId, &d)
	case "createPicture":
		res := prpc.CreateUserPictureRequest{}
		payload, err := sender.ParseEventsResponse(in.Value, &res)
		check(err)
		d, err := c.service.CreateUserPicture(&res)
		check(err)
		c.BaseEventsController.SendResponse(payload.EventId, &d)
	case "deletePicture":
		res := prpc.DeleteUserPictureRequest{}
		_, err := sender.ParseEventsResponse(in.Value, &res)
		check(err)
		err = c.service.DeleteUserPicture(res.UserId, res.PictureId)
		check(err)
	}
}

func NewPicturesEventsController(s *services.PictureService, sender *sender.Sender) *PicturesEventsController {
	// TODO: provide topic from config
	return &PicturesEventsController{
		service:              s,
		BaseEventsController: *events.NewEventsController(sender, "pictures"),
	}
}
