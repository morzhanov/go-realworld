package eventscontroller

import (
	"context"
	"encoding/json"
	"github.com/morzhanov/go-realworld/internal/common/tracing"
	"go.uber.org/zap"
	"time"

	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
)

type BaseEventsController struct {
	tracer *opentracing.Tracer
	sender *sender.Sender
	conn   *kafka.Conn
	Logger *zap.Logger
}

func createKafkaConnection(topic string, partition int, kafkaUri string) (*kafka.Conn, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", kafkaUri, topic, partition)
	if err != nil {
		return nil, err
	}
	if err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return nil, err
	}
	if err := conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return nil, err
	}
	return conn, nil
}

func (c *BaseEventsController) CreateSpan(in *kafka.Message) opentracing.Span {
	return tracing.StartSpanFromEventsRequest(*c.tracer, in)
}

func (c *BaseEventsController) Listen(
	ctx context.Context,
	cancel context.CancelFunc,
	processRequest func(*kafka.Message),
) {
loop:
	for {
		batch := c.conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max
		b := make([]byte, 10e3)              // 10KB max per message

		_, err := batch.Read(b)
		if err != nil {
			break
		}
		input := kafka.Message{}
		err = json.Unmarshal(b, &input)
		if err != nil {
			cancel()
			c.Logger.Fatal(err.Error())
		}
		go processRequest(&input)

		if err := batch.Close(); err != nil {
			cancel()
			c.Logger.Fatal(err.Error())
		}
		if err := c.conn.Close(); err != nil {
			cancel()
			c.Logger.Fatal(err.Error())
		}

		select {
		case <-ctx.Done():
			break loop
		default:
			continue
		}
	}
}

func (c *BaseEventsController) SendResponse(eventId string, data interface{}, span *opentracing.Span) error {
	return c.sender.SendEventsResponse(eventId, data, span)
}

func NewEventsController(
	s *sender.Sender,
	tracer *opentracing.Tracer,
	topic string,
	kafkaUri string,
	logger *zap.Logger,
) (*BaseEventsController, error) {
	conn, err := createKafkaConnection(topic, 0, kafkaUri)
	return &BaseEventsController{sender: s, conn: conn, tracer: tracer, Logger: logger}, err
}
