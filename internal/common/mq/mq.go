package mq

import (
	"context"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/segmentio/kafka-go"
)

type MQ struct {
	Conn     *kafka.Conn
	KafkaUri string
	Topic    string
}

func (mq *MQ) createTopic() error {
	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             mq.Topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}
	return mq.Conn.CreateTopics(topicConfigs...)
}

func (mq *MQ) CreateReader(groupId string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{mq.KafkaUri},
		Topic:    mq.Topic,
		GroupID:  groupId,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
}

func NewMq(c *config.Config, topic string) (msgQ *MQ, err error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", c.KafkaUri, c.KafkaTopic, 0)
	if err != nil {
		return nil, err
	}
	msgQ = &MQ{
		conn,
		c.KafkaUri,
		topic,
	}
	if err := msgQ.createTopic(); err != nil {
		return nil, err
	}
	return
}
