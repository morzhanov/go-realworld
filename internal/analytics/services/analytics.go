package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	. "github.com/morzhanov/go-realworld/internal/analytics/models"
	. "github.com/morzhanov/go-realworld/internal/analytics/mq"
	"github.com/segmentio/kafka-go"
)

type AnalyticsService struct {
	mq *MQ
}

func (s *AnalyticsService) LogData(data *AnalyticsEntry) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if _, err = s.mq.Conn.Write(bytes); err != nil {
		return err
	}
	return nil
}

func (s *AnalyticsService) GetLog(data *GetLogsInput) (res *AnalyticsEntry, err error) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   s.mq.Brokers,
		Topic:     s.mq.Topic,
		Partition: s.mq.Partition,
		MaxWait:   1 * time.Second,
	})
	defer r.Close()

	r.SetOffset(int64(data.Offset))

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	m, err := r.ReadMessage(ctx)
	if err != nil && err.Error() != "context deadline exceeded" {
		return nil, err
	}
	if m.Value == nil {
		return nil, errors.New(fmt.Sprintf("No message found on the %v offset", data.Offset))
	}

	res = &AnalyticsEntry{}
	if err = json.Unmarshal(m.Value, res); err != nil {
		return nil, err
	}
	return res, nil
}

func NewAnalyticsService(mq *MQ) *AnalyticsService {
	return &AnalyticsService{mq}
}
