package mq_test

import (
	"github.com/morzhanov/go-realworld/internal/common/mq"
	"github.com/segmentio/kafka-go"
)

type MqMock struct {
	createReaderMock func(groupId string) *kafka.Reader
	connMock         func() *kafka.Conn
	kafkaUriMock     func() string
	topicMock        func() string
}

func (m *MqMock) CreateReader(groupId string) *kafka.Reader {
	return m.createReaderMock(groupId)
}
func (m *MqMock) Conn() *kafka.Conn {
	return m.connMock()
}
func (m *MqMock) KafkaUri() string {
	return m.kafkaUriMock()
}
func (m *MqMock) Topic() string {
	return m.topicMock()
}

func NewMqMock(
	createReader func(groupId string) *kafka.Reader,
	conn func() *kafka.Conn,
	kafkaUri func() string,
	topic func() string,
) mq.MQ {
	m := MqMock{}
	if createReader != nil {
		m.createReaderMock = createReader
	} else {
		m.createReaderMock = func(groupId string) *kafka.Reader { return nil }
	}
	if conn != nil {
		m.connMock = conn
	} else {
		m.connMock = func() *kafka.Conn { return nil }
	}
	if kafkaUri != nil {
		m.kafkaUriMock = kafkaUri
	} else {
		m.kafkaUriMock = func() string { return "uri" }
	}
	if topic != nil {
		m.topicMock = topic
	} else {
		m.topicMock = func() string { return "topic" }
	}
	return &m
}
