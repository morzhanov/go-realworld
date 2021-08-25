package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	prpc "github.com/morzhanov/go-realworld/api/rpc/pictures"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/morzhanov/go-realworld/internal/pictures/services"
	"github.com/segmentio/kafka-go"
)

type PicturesEventsController struct {
	service *services.PictureService
	Conn    *kafka.Conn
}

// TODO: move to common package
type EventMessage struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type BaseEventPayload struct {
	EventId string `json:"eventId"`
}

type ErrorMessage struct {
	Error string `json:"error"`
}

type SuccessMessage struct{}

type GetUserPicturesInput struct {
	BaseEventPayload
	UserId string `json:"userId"`
}

type GetUserPictureInput struct {
	BaseEventPayload
	UserId    string `json:"userId"`
	PictireId string `json:"pictureId"`
}

type DeleteUserPictureInput struct {
	BaseEventPayload
	UserId    string `json:"userId"`
	PictireId string `json:"pictureId"`
}

func check(err error) {
	// TODO: handle error
	if err != nil {
		log.Fatal(err)
	}
}

// TODO: looks like common function for all controllers
// TODO: topic name should be service name
func CreateKafkaConnection(topic string, partition int) *kafka.Conn {
	// TODO: provide kafka uri
	uri := "192.168.0.180:32181"
	conn, _ := kafka.DialLeader(context.Background(), "tcp", uri, topic, partition)
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	return conn
}

// TODO: looks like common function for all controllers
func (c *PicturesEventsController) Listen() {
	c.Conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	batch := c.Conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

	b := make([]byte, 10e3) // 10KB max per message
	for {
		_, err := batch.Read(b)
		if err != nil {
			break
		}
		go c.processRequest(&b)
	}

	if err := batch.Close(); err != nil {
		log.Fatal("failed to close batch:", err)
	}
	if err := c.Conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}
}

func (c *PicturesEventsController) processRequest(b *[]byte) {
	input := EventMessage{}
	err := json.Unmarshal(*b, &input)
	check(err)

	switch input.Key {
	case "getPictures":
		res := prpc.GetUserPicturesRequest{}
		payload, err := sender.ParseEventsResponse(input.Value, &res)
		check(err)
		d, err := c.service.GetUserPictures(res.UserId)
		check(err)
		c.sendResponse(payload.EventId, &d)
	case "getPicture":
		res := prpc.GetUserPictureRequest{}
		payload, err := sender.ParseEventsResponse(input.Value, &res)
		check(err)
		d, err := c.service.GetUserPicture(res.UserId, res.PictureId)
		check(err)
		c.sendResponse(payload.EventId, &d)
	case "createPicture":
		res := prpc.CreateUserPictureRequest{}
		payload, err := sender.ParseEventsResponse(input.Value, &res)
		check(err)
		d, err := c.service.CreateUserPicture(&res)
		check(err)
		c.sendResponse(payload.EventId, &d)
	case "deletePicture":
		res := prpc.DeleteUserPictureRequest{}
		_, err := sender.ParseEventsResponse(input.Value, &res)
		check(err)
		err = c.service.DeleteUserPicture(res.UserId, res.PictureId)
		check(err)
	}
}

// TODO: seems like common logic
// TODO: send the response with client module and use payload.EventId as event_uuid
func (c *PicturesEventsController) sendResponse(eventUuid string, value interface{}) {
	// TODO: send response via client kafka
	payload, err := json.Marshal(&value)
	check(err)
	fmt.Printf("payload %v\n", payload)
}

func NewPicturesEventsController(s *services.PictureService) *PicturesEventsController {
	// TODO: provide topic and partition from config
	conn := CreateKafkaConnection("topic", 0)
	return &PicturesEventsController{service: s, Conn: conn}
}
