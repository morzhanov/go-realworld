package events

import (
	"context"
	"fmt"

	prpc "github.com/morzhanov/go-realworld/api/rpc/pictures"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/events"
	"github.com/morzhanov/go-realworld/internal/common/events/eventscontroller"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/morzhanov/go-realworld/internal/pictures/services"
)

type PicturesEventsController struct {
	eventscontroller.BaseEventsController
	service *services.PictureService
}

func (c *PicturesEventsController) processRequest(in *events.EventMessage) error {
	switch in.Key {
	case "getPictures":
		return c.getPictures(in)
	case "getPicture":
		return c.getPicture(in)
	case "createPicture":
		return c.createPicture(in)
	case "deletePicture":
		return c.deletePicture(in)
	default:
		return fmt.Errorf("Wrong event name")
	}
}

func (c *PicturesEventsController) getPictures(in *events.EventMessage) error {
	res := prpc.GetUserPicturesRequest{}
	payload, err := events.ParseEventsResponse(in.Value, &res)
	if err != nil {
		return err
	}
	d, err := c.service.GetUserPictures(res.UserId)
	if err != nil {
		return err
	}
	c.BaseEventsController.SendResponse(payload.EventId, &d)
	return nil
}

func (c *PicturesEventsController) getPicture(in *events.EventMessage) error {
	res := prpc.GetUserPictureRequest{}
	payload, err := events.ParseEventsResponse(in.Value, &res)
	if err != nil {
		return err
	}
	d, err := c.service.GetUserPicture(res.UserId, res.PictureId)
	if err != nil {
		return err
	}
	c.BaseEventsController.SendResponse(payload.EventId, &d)
	return nil
}

func (c *PicturesEventsController) createPicture(in *events.EventMessage) error {
	res := prpc.CreateUserPictureRequest{}
	payload, err := events.ParseEventsResponse(in.Value, &res)
	if err != nil {
		return err
	}
	d, err := c.service.CreateUserPicture(&res)
	if err != nil {
		return err
	}
	c.BaseEventsController.SendResponse(payload.EventId, &d)
	return nil
}

func (c *PicturesEventsController) deletePicture(in *events.EventMessage) error {
	res := prpc.DeleteUserPictureRequest{}
	if _, err := events.ParseEventsResponse(in.Value, &res); err != nil {
		return err
	}
	return c.service.DeleteUserPicture(res.UserId, res.PictureId)
}

func (c *PicturesEventsController) Listen(ctx context.Context) {
	c.BaseEventsController.Listen(
		ctx,
		func(m *events.EventMessage) { c.processRequest(m) },
	)
}

func NewPicturesEventsController(
	s *services.PictureService,
	c *config.Config,
	sender *sender.Sender,
) *PicturesEventsController {
	controller := eventscontroller.NewEventsController(sender, c.PicturesKafkaTopic, c.KafkaUri)
	return &PicturesEventsController{
		service:              s,
		BaseEventsController: *controller,
	}
}
