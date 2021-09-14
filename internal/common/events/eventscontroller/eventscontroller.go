package eventscontroller

import (
	"context"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/morzhanov/go-realworld/internal/common/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type BaseEventsController struct {
	tracer   *opentracing.Tracer
	sender   *sender.Sender
	conn     *kafka.Conn
	topic    string
	kafkaUri string
	Logger   *zap.Logger
}

func (c *BaseEventsController) createTopic() error {
	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             c.topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}
	return c.conn.CreateTopics(topicConfigs...)
}

func (c *BaseEventsController) createKafkaConnection(partition int) error {
	conn, err := kafka.DialLeader(context.Background(), "tcp", c.kafkaUri, c.topic, partition)
	if err != nil {
		return err
	}
	c.conn = conn
	if err := c.createTopic(); err != nil {
		return err
	}
	return nil
}

func (c *BaseEventsController) CreateSpan(in *kafka.Message) opentracing.Span {
	return tracing.StartSpanFromEventsRequest(*c.tracer, in)
}

func (c *BaseEventsController) Listen(
	ctx context.Context,
	processRequest func(*kafka.Message),
) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{c.kafkaUri},
		Topic:    c.topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			c.Logger.Error(err.Error())
			continue
		}
		go processRequest(&m)
		select {
		case <-ctx.Done():
			break
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
	c := &BaseEventsController{sender: s, tracer: tracer, Logger: logger, topic: topic, kafkaUri: kafkaUri}
	err := c.createKafkaConnection(0)
	return c, err
}
