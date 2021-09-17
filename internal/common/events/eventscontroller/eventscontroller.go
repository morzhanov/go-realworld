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

type baseEventsController struct {
	tracer          opentracing.Tracer
	sender          sender.Sender
	mq              mq.MQ
	logger          *zap.Logger
	consumerGroupId string
}

type BaseEventsController interface {
	CreateSpan(in *kafka.Message) opentracing.Span
	Listen(ctx context.Context, processRequest func(*kafka.Message))
	SendResponse(eventId string, data interface{}, span *opentracing.Span) error
	Logger() *zap.Logger
	ConsumerGroupId() string
}

func (c *baseEventsController) CreateSpan(in *kafka.Message) opentracing.Span {
	return tracing.StartSpanFromEventsRequest(c.tracer, in)
}

func (c *baseEventsController) Listen(
	ctx context.Context,
	processRequest func(*kafka.Message),
) {
	r := c.mq.CreateReader(c.consumerGroupId)
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			c.logger.Error(err.Error())
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

func (c *baseEventsController) SendResponse(eventId string, data interface{}, span *opentracing.Span) error {
	return c.sender.SendEventsResponse(eventId, data, span)
}

func (c *baseEventsController) Logger() *zap.Logger {
	return c.logger
}

func (c *baseEventsController) ConsumerGroupId() string {
	return c.consumerGroupId
}

func NewEventsController(
	s sender.Sender,
	tracer opentracing.Tracer,
	logger *zap.Logger,
	conf *config.Config,
) (BaseEventsController, error) {
	msgQ, err := mq.NewMq(conf, conf.KafkaTopic)
	if err != nil {
		return nil, err
	}
	c := &baseEventsController{
		sender:          s,
		tracer:          tracer,
		mq:              msgQ,
		logger:          logger,
		consumerGroupId: conf.KafkaConsumerGroupId,
	}
	return c, err
}
