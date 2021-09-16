package eventscontroller

import (
	"context"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/mq"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/morzhanov/go-realworld/internal/common/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type BaseEventsController struct {
	tracer          *opentracing.Tracer
	sender          *sender.Sender
	mq              *mq.MQ
	Logger          *zap.Logger
	ConsumerGroupId string
}

func (c *BaseEventsController) CreateSpan(in *kafka.Message) opentracing.Span {
	return tracing.StartSpanFromEventsRequest(*c.tracer, in)
}

func (c *BaseEventsController) Listen(
	ctx context.Context,
	processRequest func(*kafka.Message),
) {
	r := c.mq.CreateReader(c.ConsumerGroupId)
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
	logger *zap.Logger,
	conf *config.Config,
) (*BaseEventsController, error) {
	msgQ, err := mq.NewMq(conf, conf.KafkaTopic)
	if err != nil {
		return nil, err
	}
	c := &BaseEventsController{
		sender:          s,
		tracer:          tracer,
		mq:              msgQ,
		Logger:          logger,
		ConsumerGroupId: conf.KafkaConsumerGroupId,
	}
	return c, err
}
