package mq

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type MQ struct {
	Conn      *kafka.Conn
	Brokers   []string
	Topic     string
	Partition int
}

func NewMq(topic string, partition int) (*MQ, error) {
	const connUri = "192.168.0.180:9092"
	conn, err := kafka.DialLeader(context.Background(), "tcp", connUri, topic, partition)
	if err != nil {
		return nil, err
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	return &MQ{conn, []string{connUri}, topic, partition}, nil
}
