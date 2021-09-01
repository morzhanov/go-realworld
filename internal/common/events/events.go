package events

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/segmentio/kafka-go"
)

type BaseEventsController struct {
	sender *sender.Sender
	conn   *kafka.Conn
}

// TODO: common logic
func check(err error) {
	// TODO: handle error
	if err != nil {
		log.Fatal(err)
	}
}

func createKafkaConnection(topic string, partition int) *kafka.Conn {
	// TODO: provide kafka uri
	uri := "192.168.0.180:32181"
	conn, _ := kafka.DialLeader(context.Background(), "tcp", uri, topic, partition)
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	return conn
}

func (c *BaseEventsController) Listen(processRequest func(*sender.EventMessage)) {
	c.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	batch := c.conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

	b := make([]byte, 10e3) // 10KB max per message
	for {
		_, err := batch.Read(b)
		if err != nil {
			break
		}

		input := sender.EventMessage{}
		err = json.Unmarshal(b, &input)
		check(err)

		go processRequest(&input)
	}

	err := batch.Close()
	check(err)
	err = c.conn.Close()
	check(err)
}

func (c *BaseEventsController) SendResponse(eventId string, data interface{}) {
	c.sender.SendEventsResponse(eventId, data)
}

func NewEventsController(s *sender.Sender, topic string) *BaseEventsController {
	conn := createKafkaConnection(topic, 0)
	return &BaseEventsController{
		sender: s,
		conn:   conn,
	}
}
