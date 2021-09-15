package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	anrpc "github.com/morzhanov/go-realworld/api/grpc/analytics"
	. "github.com/morzhanov/go-realworld/internal/analytics/models"
	"github.com/morzhanov/go-realworld/internal/common/mq"
	"github.com/segmentio/kafka-go"
)

type AnalyticsService struct {
	mq *mq.MQ
}

func (s *AnalyticsService) LogData(data *anrpc.LogDataRequest) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if _, err = s.mq.Conn.Write(bytes); err != nil {
		return err
	}
	return nil
}

func (s *AnalyticsService) GetLog(data *anrpc.GetLogRequest) (res *anrpc.AnalyticsEntryMessage, err error) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   s.mq.Brokers,
		Topic:     s.mq.Topic,
		Partition: s.mq.Partition,
		MaxWait:   1 * time.Second,
	})
	if err := r.SetOffset(int64(data.Offset)); err != nil {
		return nil, err
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	m, err := r.ReadMessage(ctx)
	if err := r.Close(); err != nil {
		return nil, err
	}
	if err != nil && err.Error() != "context deadline exceeded" {
		return nil, err
	}
	if m.Value == nil {
		return nil, fmt.Errorf("no message found on the %v offset", data.Offset)
	}

	result := &AnalyticsEntry{}
	if err = json.Unmarshal(m.Value, result); err != nil {
		return nil, err
	}
	return &anrpc.AnalyticsEntryMessage{
		Id:        result.ID,
		UserId:    result.UserID,
		Operation: result.Operation,
		Data:      result.Data,
	}, nil
}

func NewAnalyticsService(mq *mq.MQ) *AnalyticsService {
	return &AnalyticsService{mq}
}
