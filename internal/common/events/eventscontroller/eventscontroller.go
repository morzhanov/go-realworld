package eventscontroller

import (
	"context"
	"encoding/json"
	"time"

	"github.com/morzhanov/go-realworld/internal/common/events"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

type BaseEventsController struct {
	sender *sender.Sender
	conn   *kafka.Conn
}

func createKafkaConnection(topic string, partition int) *kafka.Conn {
	kafkaUri := viper.GetString("KAFKA_URI")
	conn, _ := kafka.DialLeader(context.Background(), "tcp", kafkaUri, topic, partition)
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	return conn
}

func (c *BaseEventsController) Listen(processRequest func(*events.EventMessage)) error {
	c.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	batch := c.conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

	b := make([]byte, 10e3) // 10KB max per message
	for {
		_, err := batch.Read(b)
		if err != nil {
			break
		}

		input := events.EventMessage{}
		err = json.Unmarshal(b, &input)
		if err != nil {
			return err
		}
		go processRequest(&input)
	}

	if err := batch.Close(); err != nil {
		return err
	}
	if err := c.conn.Close(); err != nil {
		return err
	}
	return nil
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
