package rpc

import (
	"context"

	anrpc "github.com/morzhanov/go-realworld/api/rpc/analytics"
	"github.com/morzhanov/go-realworld/internal/analytics/services"
	"github.com/morzhanov/go-realworld/internal/common/helper"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AnalyticsRpcServer struct {
	anrpc.UnimplementedAnalyticsServer
	analyticsService *services.AnalyticsService
	server           *grpc.Server
}

func (s *AnalyticsRpcServer) LogData(ctx context.Context, in *anrpc.LogDataRequest) (res *emptypb.Empty, err error) {
	err = s.analyticsService.LogData(in)
	return res, err
}

func (s *AnalyticsRpcServer) GetLog(ctx context.Context, in *anrpc.GetLogRequest) (res *anrpc.AnalyticsEntryMessage, err error) {
	return s.analyticsService.GetLog(in)
}

func (s *AnalyticsRpcServer) Listen() error {
	port := viper.GetString("ANALYTICS_GRPC_PORT")
	return helper.StartGrpcServer(s.server, port)
}

func NewAnalyticsRpcService(analyticsService services.AnalyticsService) (s *grpc.Server) {
	s = grpc.NewServer()
	anrpc.RegisterAnalyticsServer(s, &AnalyticsRpcServer{analyticsService: &analyticsService})
	return s
}
