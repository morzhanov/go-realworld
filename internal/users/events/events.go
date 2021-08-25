package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	urpc "github.com/morzhanov/go-realworld/api/rpc/users"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/morzhanov/go-realworld/internal/users/services"
	"github.com/segmentio/kafka-go"
)

type UsersEventsController struct {
	service *services.UsersService
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

type GetUserDataInput struct {
	BaseEventPayload
	UserId string `json:"userId"`
}

type GetUserDataByUsernameInput struct {
	BaseEventPayload
	Username string `json:"username"`
}

type DeleteUserInput struct {
	BaseEventPayload
	UserId string `json:"userId"`
}

type ErrorMessage struct {
	Error string `json:"error"`
}

type SuccessMessage struct{}

func check(err error) {
	// TODO: handle error
	if err != nil {
		log.Fatal(err)
	}
}

// TODO: topic name should be service name
func CreateKafkaConnection(topic string, partition int) *kafka.Conn {
	// TODO: provide kafka uri
	uri := "192.168.0.180:32181"
	conn, _ := kafka.DialLeader(context.Background(), "tcp", uri, topic, partition)
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	return conn
}

func (c *UsersEventsController) Listen() {
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

func (c *UsersEventsController) processRequest(b *[]byte) {
	input := EventMessage{}
	err := json.Unmarshal(*b, &input)
	check(err)

	switch input.Key {
	case "getUser":
		res := urpc.GetUserDataRequest{}
		payload, err := sender.ParseEventsResponse(input.Value, &res)
		check(err)
		d, err := c.service.GetUserData(res.UserId)
		check(err)
		c.sendResponse(payload.EventId, &d)
	case "getUserByUsername":
		res := urpc.GetUserDataByUsernameRequest{}
		payload, err := sender.ParseEventsResponse(input.Value, &res)
		check(err)
		d, err := c.service.GetUserData(res.Username)
		check(err)
		c.sendResponse(payload.EventId, &d)
	case "validatePassword":
		res := urpc.ValidateUserPasswordRequest{}
		payload, err := sender.ParseEventsResponse(input.Value, &res)
		check(err)
		err = c.service.ValidateUserPassword(&res)
		check(err)
		c.sendResponse(payload.EventId, nil)
	case "createUser":
		res := urpc.CreateUserRequest{}
		payload, err := sender.ParseEventsResponse(input.Value, &res)
		check(err)
		d, err := c.service.CreateUser(&res)
		check(err)
		c.sendResponse(payload.EventId, &d)
	case "deleteUser":
		res := urpc.DeleteUserRequest{}
		payload, err := sender.ParseEventsResponse(input.Value, &res)
		check(err)
		err = c.service.DeleteUser(res.UserId)
		check(err)
		c.sendResponse(payload.EventId, nil)
	}
}

// TODO: seems like common
// TODO: send the response with client module and use payload.EventId as event_uuid
func (c *UsersEventsController) sendResponse(eventUuid string, value interface{}) {
	// TODO: send response via client kafka
	payload, err := json.Marshal(&value)
	check(err)
	fmt.Printf("payload %v\n", payload)
}

func NewUsersEventsController(s *services.UsersService) *UsersEventsController {
	// TODO: provide topic and partition from config
	conn := CreateKafkaConnection("topic", 0)
	return &UsersEventsController{service: s, Conn: conn}
}
