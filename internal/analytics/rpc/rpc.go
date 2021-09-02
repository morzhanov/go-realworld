package rpc

import (
	"context"

	anrpc "github.com/morzhanov/go-realworld/api/rpc/analytics"
	"github.com/morzhanov/go-realworld/internal/analytics/services"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/helper"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AnalyticsRpcServer struct {
	anrpc.UnimplementedAnalyticsServer
	analyticsService *services.AnalyticsService
	port             string
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
	return helper.StartGrpcServer(s.server, s.port)
}

func NewAnalyticsRpcService(
	analyticsService services.AnalyticsService,
	c *config.Config,
) (server *AnalyticsRpcServer) {
	server = &AnalyticsRpcServer{analyticsService: &analyticsService, port: c.AnalyticsGrpcPort}
	anrpc.RegisterAnalyticsServer(grpc.NewServer(), server)
	return
}
