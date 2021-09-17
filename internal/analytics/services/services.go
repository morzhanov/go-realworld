package services

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"

	anrpc "github.com/morzhanov/go-realworld/api/grpc/analytics"
	"github.com/morzhanov/go-realworld/internal/common/mq"
)

type analyticsService struct {
	mq        mq.MQ
	dataTopic string
}

type AnalyticsService interface {
	LogData(data *anrpc.LogDataRequest) (err error)
	GetLog(_ *emptypb.Empty) (res *anrpc.GetLogsMessage, err error)
}

func (s *analyticsService) LogData(data *anrpc.LogDataRequest) (err error) {
	defer func() { err = errors.Wrap(err, "analyticsService:logData") }()
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if _, err = s.mq.Conn().Write(bytes); err != nil {
		return err
	}
	return nil
}

func (s *analyticsService) GetLog(_ *emptypb.Empty) (res *anrpc.GetLogsMessage, err error) {
	defer func() { err = errors.Wrap(err, "analyticsService:getlog") }()
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

func NewAnalyticsService(mq mq.MQ, dataTopic string) AnalyticsService {
	return &analyticsService{mq, dataTopic}
}
