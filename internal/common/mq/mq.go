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

func NewMq(topic string, partition int) *MQ {
	const KAFKA_CONNECTION_URI = "192.168.0.180:19092"
	conn, _ := kafka.DialLeader(context.Background(), "tcp", KAFKA_CONNECTION_URI, topic, partition)
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	return &MQ{conn, []string{KAFKA_CONNECTION_URI}, topic, partition}
}