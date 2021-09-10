package rpc

import (
	"context"
	"fmt"
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

func (s *AnalyticsRpcServer) Listen(ctx context.Context, cancel context.CancelFunc) {
	s.BaseGrpcServer.Listen(ctx, cancel, s.server)
}

func NewAnalyticsRpcServer(
	analyticsService *services.AnalyticsService,
	c *config.Config,
	tracer *opentracing.Tracer,
	logger *zap.Logger,
) (s *AnalyticsRpcServer) {
	uri := fmt.Sprintf("%s:%s", c.GrpcAddr, c.GrpcPort)
	bs := grpcserver.NewGrpcServer(tracer, logger, uri)
	s = &AnalyticsRpcServer{
		analyticsService: analyticsService,
		BaseGrpcServer:   bs,
		server:           grpc.NewServer(),
	}
	anrpc.RegisterAnalyticsServer(s.server, s)
	return
}
