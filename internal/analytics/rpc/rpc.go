package rpc

import (
	"context"

	anrpc "github.com/morzhanov/go-realworld/api/rpc/analytics"
	"github.com/morzhanov/go-realworld/internal/analytics/services"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/helper"
	"github.com/morzhanov/go-realworld/internal/common/tracing"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AnalyticsRpcServer struct {
	anrpc.UnimplementedAnalyticsServer
	analyticsService *services.AnalyticsService
	port             string
	server           *grpc.Server
	tracer           *opentracing.Tracer
}

func (s *AnalyticsRpcServer) LogData(ctx context.Context, in *anrpc.LogDataRequest) (res *emptypb.Empty, err error) {
	// TODO: maybe we should somehow generalize this step via middleware or something
	// TODO: because we'are using the same code for tracing in all controllers
	span := tracing.StartSpanFromGrpcRequest(*s.tracer, ctx)
	defer span.Finish()
	err = s.analyticsService.LogData(in)
	return res, err
}

func (s *AnalyticsRpcServer) GetLog(ctx context.Context, in *anrpc.GetLogRequest) (res *anrpc.AnalyticsEntryMessage, err error) {
	span := tracing.StartSpanFromGrpcRequest(*s.tracer, ctx)
	defer span.Finish()
	return s.analyticsService.GetLog(in)
}

func (s *AnalyticsRpcServer) Listen(ctx context.Context, logger *zap.Logger) error {
	return helper.StartGrpcServer(ctx, s.server, s.port, logger)
}

func NewAnalyticsRpcServer(
	analyticsService *services.AnalyticsService,
	c *config.Config,
	tracer *opentracing.Tracer,
) (server *AnalyticsRpcServer) {
	server = &AnalyticsRpcServer{analyticsService: analyticsService, port: c.GrpcPort, tracer: tracer}
	anrpc.RegisterAnalyticsServer(grpc.NewServer(), server)
	return
}
