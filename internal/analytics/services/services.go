package services

import (
	"context"
	"encoding/json"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"

	anrpc "github.com/morzhanov/go-realworld/api/grpc/analytics"
	"github.com/morzhanov/go-realworld/internal/common/mq"
)

type AnalyticsService struct {
	mq        *mq.MQ
	dataTopic string
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

func (s *AnalyticsService) GetLog(_ *emptypb.Empty) (res *anrpc.GetLogsMessage, err error) {
	r := s.mq.CreateReader(s.dataTopic)
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	m, err := r.ReadMessage(ctx)
	if m.Value == nil && err != nil {
		return nil, nil
	}
	if err := r.Close(); err != nil {
		return nil, err
	}
	result := &anrpc.AnalyticsEntryMessage{}
	if err = json.Unmarshal(m.Value, result); err != nil {
		return nil, err
	}
	res = &anrpc.GetLogsMessage{Logs: []*anrpc.AnalyticsEntryMessage{result}}
	return
}

func NewAnalyticsService(mq *mq.MQ, dataTopic string) *AnalyticsService {
	return &AnalyticsService{mq, dataTopic}
}
