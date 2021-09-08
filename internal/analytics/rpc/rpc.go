package rpc

import (
	"context"
	"github.com/morzhanov/go-realworld/internal/common/grpc/grpcserver"
	"go.uber.org/zap"

	anrpc "github.com/morzhanov/go-realworld/api/rpc/analytics"
	"github.com/morzhanov/go-realworld/internal/analytics/services"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AnalyticsRpcServer struct {
	anrpc.UnimplementedAnalyticsServer
	*grpcserver.BaseGrpcServer
	analyticsService *services.AnalyticsService
	server           *grpc.Server
}

func (s *AnalyticsRpcServer) LogData(ctx context.Context, in *anrpc.LogDataRequest) (res *emptypb.Empty, err error) {
	span := s.PrepareContext(ctx)
	defer span.Finish()
	err = s.analyticsService.LogData(in)
	return res, err
}

func (s *AnalyticsRpcServer) GetLog(ctx context.Context, in *anrpc.GetLogRequest) (res *anrpc.AnalyticsEntryMessage, err error) {
	span := s.PrepareContext(ctx)
	defer span.Finish()
	return s.analyticsService.GetLog(in)
}

func (s *AnalyticsRpcServer) Listen(ctx context.Context) error {
	return s.BaseGrpcServer.Listen(ctx, s.server)
}

func NewAnalyticsRpcServer(
	analyticsService *services.AnalyticsService,
	c *config.Config,
	tracer *opentracing.Tracer,
	logger *zap.Logger,
) (s *AnalyticsRpcServer) {
	bs := grpcserver.NewGrpcServer(tracer, logger, c.GrpcPort)
	s = &AnalyticsRpcServer{
		analyticsService: analyticsService,
		BaseGrpcServer:   bs,
	}
	anrpc.RegisterAnalyticsServer(grpc.NewServer(), s)
	return
}
