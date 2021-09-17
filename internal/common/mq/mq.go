package mq

import (
	"context"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/segmentio/kafka-go"
)

type mQ struct {
	conn     *kafka.Conn
	kafkaUri string
	topic    string
}

type MQ interface {
	CreateReader(groupId string) *kafka.Reader
	Conn() *kafka.Conn
	KafkaUri() string
	Topic() string
}

func (mq *mQ) createTopic() error {
	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             mq.topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}
	return mq.conn.CreateTopics(topicConfigs...)
}

func (mq *mQ) CreateReader(groupId string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{mq.kafkaUri},
		Topic:    mq.topic,
		GroupID:  groupId,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
}

func (mq *mQ) Conn() *kafka.Conn {
	return mq.conn
}

func (mq *mQ) KafkaUri() string {
	return mq.kafkaUri
}

func (mq *mQ) Topic() string {
	return mq.topic
}

func NewMq(c *config.Config, topic string) (res MQ, err error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", c.KafkaUri, c.KafkaTopic, 0)
	if err != nil {
		return nil, err
	}
	msgQ := mQ{
		conn,
		c.KafkaUri,
		topic,
	}
	if err := msgQ.createTopic(); err != nil {
		return nil, err
	}
	return &msgQ, nil
}
