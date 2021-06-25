package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/morzhanov/go-realworld/internal/pictures/dto"
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

type CreateUserPictureInput struct {
	BaseEventPayload
	UserId string `json:"userId"`
	dto.CreatePicturesDto
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

/*
Request event schema:
	Key: "resource:action"
	Value: "{payload}"
*/
/*
Response event schema:
	Key: "result:event_uuid"
	Value: "{payload}"
*/
func (c *PicturesEventsController) processRequest(b *[]byte) {
	input := EventMessage{}
	err := json.Unmarshal(*b, &input)
	check(err)

	switch input.Key {
	case "pictures:get":
		payload := GetUserPicturesInput{}
		err := json.Unmarshal([]byte(input.Value), &payload)
		check(err)

		res, err := c.service.GetUserPictures(payload.UserId)
		check(err)

		c.sendResponse(payload.EventId, res)
	case "pictures:get_one":
		payload := GetUserPictureInput{}
		err := json.Unmarshal([]byte(input.Value), &payload)
		check(err)

		res, err := c.service.GetUserPicture(payload.UserId, payload.PictireId)
		check(err)

		c.sendResponse(payload.EventId, res)
	case "pictures:create":
		payload := CreateUserPictureInput{}
		err := json.Unmarshal([]byte(input.Value), &payload)
		check(err)

		res, err := c.service.CreateUserPicture(payload.UserId, &payload.CreatePicturesDto)
		check(err)
		c.sendResponse(payload.EventId, res)
	case "pictures:delete":
		payload := DeleteUserPictureInput{}
		err := json.Unmarshal([]byte(input.Value), &payload)
		check(err)

		err = c.service.DeleteUserPicture(payload.UserId, payload.PictireId)
		check(err)

		c.sendResponse(payload.EventId, &SuccessMessage{})
	}
}

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
