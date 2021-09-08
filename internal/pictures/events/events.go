package events

import (
	"context"
	"errors"
	"go.uber.org/zap"

	prpc "github.com/morzhanov/go-realworld/api/rpc/pictures"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/events"
	"github.com/morzhanov/go-realworld/internal/common/events/eventscontroller"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/morzhanov/go-realworld/internal/pictures/services"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
)

type PicturesEventsController struct {
	eventscontroller.BaseEventsController
	service *services.PictureService
	logger  *zap.Logger
}

func (c *PicturesEventsController) processRequest(in *kafka.Message) error {
	switch string(in.Key) {
	case "getPictures":
		return c.getPictures(in)
	case "getPicture":
		return c.getPicture(in)
	case "createPicture":
		return c.createPicture(in)
	case "deletePicture":
		return c.deletePicture(in)
	default:
		return errors.New("wrong event name")
	}
}

func (c *PicturesEventsController) getPictures(in *kafka.Message) error {
	span := c.CreateSpan(in)
	defer span.Finish()

	res := prpc.GetUserPicturesRequest{}
	payload, err := events.ParseEventsResponse(in.Value, &res)
	if err != nil {
		return err
	}
	d, err := c.service.GetUserPictures(res.UserId)
	if err != nil {
		return err
	}
	return c.BaseEventsController.SendResponse(payload.EventId, &d, &span)
}

func (c *PicturesEventsController) getPicture(in *kafka.Message) error {
	span := c.CreateSpan(in)
	defer span.Finish()

	res := prpc.GetUserPictureRequest{}
	payload, err := events.ParseEventsResponse(in.Value, &res)
	if err != nil {
		return err
	}
	d, err := c.service.GetUserPicture(res.UserId, res.PictureId)
	if err != nil {
		return err
	}
	return c.BaseEventsController.SendResponse(payload.EventId, &d, &span)
}

func (c *PicturesEventsController) createPicture(in *kafka.Message) error {
	span := c.CreateSpan(in)
	defer span.Finish()

	res := prpc.CreateUserPictureRequest{}
	payload, err := events.ParseEventsResponse(in.Value, &res)
	if err != nil {
		return err
	}
	d, err := c.service.CreateUserPicture(&res)
	if err != nil {
		return err
	}
	return c.BaseEventsController.SendResponse(payload.EventId, &d, &span)
}

func (c *PicturesEventsController) deletePicture(in *kafka.Message) error {
	span := c.CreateSpan(in)
	defer span.Finish()

	res := prpc.DeleteUserPictureRequest{}
	if _, err := events.ParseEventsResponse(in.Value, &res); err != nil {
		return err
	}
	return c.service.DeleteUserPicture(res.UserId, res.PictureId)
}

func (c *PicturesEventsController) Listen(ctx context.Context) error {
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

func NewPicturesEventsController(
	s *services.PictureService,
	c *config.Config,
	sender *sender.Sender,
	tracer *opentracing.Tracer,
	logger *zap.Logger,
) (*PicturesEventsController, error) {
	controller, err := eventscontroller.NewEventsController(
		sender,
		tracer,
		c.KafkaTopic,
		c.KafkaUri,
	)
	return &PicturesEventsController{
		service:              s,
		BaseEventsController: *controller,
		logger:               logger,
	}, err
}
