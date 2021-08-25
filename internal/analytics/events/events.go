package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	anrpc "github.com/morzhanov/go-realworld/api/rpc/analytics"
	"github.com/morzhanov/go-realworld/internal/analytics/models"
	"github.com/morzhanov/go-realworld/internal/analytics/services"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/segmentio/kafka-go"
)

type AnalyticsEventsController struct {
	service *services.AnalyticsService
	Conn    *kafka.Conn
}

// TODO: move to common package
type EventMessage struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ErrorMessage struct {
	Error string `json:"error"`
}

type SuccessMessage struct{}

type LogDataInput struct {
	sender.BaseEventPayload
	models.AnalyticsEntry
}

type GetLogsInput struct {
	sender.BaseEventPayload
	anrpc.GetLogRequest
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
func (c *AnalyticsEventsController) Listen() {
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

func (c *AnalyticsEventsController) processRequest(b *[]byte) {
	input := EventMessage{}
	err := json.Unmarshal(*b, &input)
	check(err)

	switch input.Key {
	case "logData":
		res := anrpc.LogDataRequest{}
		_, err := sender.ParseEventsResponse(input.Value, &res)
		check(err)
		err = c.service.LogData(&res)
		check(err)
	case "getLogs":
		res := anrpc.GetLogRequest{}
		payload, err := sender.ParseEventsResponse(input.Value, &res)
		check(err)
		d, err := c.service.GetLog(&res)
		check(err)
		c.sendResponse(payload.EventId, &d)
	}
}

// TODO: seems like common
// TODO: send the response with client module and use payload.EventId as event_uuid
func (c *AnalyticsEventsController) sendResponse(eventUuid string, value interface{}) {
	// TODO: send response via client kafka
	payload, err := json.Marshal(&value)
	check(err)
	fmt.Printf("payload %v\n", payload)
}

func NewAnalyticsEventsController(s *services.AnalyticsService) *AnalyticsEventsController {
	// TODO: provide topic and partition from config
	conn := CreateKafkaConnection("topic", 0)
	return &AnalyticsEventsController{service: s, Conn: conn}
}
