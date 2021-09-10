package grpcserver

import (
	"context"
	"github.com/morzhanov/go-realworld/internal/common/helper"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/morzhanov/go-realworld/internal/common/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/common/log"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type BaseGrpcServer struct {
	Tracer *opentracing.Tracer
	Logger *zap.Logger
	Port   string
}

func (s *BaseGrpcServer) PrepareContext(ctx context.Context) opentracing.Span {
	span := tracing.StartSpanFromGrpcRequest(*s.Tracer, ctx)
	ctx = context.WithValue(ctx, "transport", sender.RpcTransport)
	return span
}

func (s *BaseGrpcServer) Listen(ctx context.Context, cancel context.CancelFunc, server *grpc.Server) {
	lis, err := net.Listen("tcp", s.Port)
	if err != nil {
		cancel()
		helper.HandleInitializationError(err, "grpc server", s.Logger)
	}

	if err := server.Serve(lis); err != nil {
		cancel()
		helper.HandleInitializationError(err, "grpc server", s.Logger)
	}
	log.Info("Grpc server started", zap.String("port", s.Port))
	<-ctx.Done()
	if err := lis.Close(); err != nil {
		cancel()
		helper.HandleInitializationError(err, "grpc server", s.Logger)
	}
}

func NewGrpcServer(
	tracer *opentracing.Tracer,
	logger *zap.Logger,
	port string,
) *BaseGrpcServer {
	return &BaseGrpcServer{tracer, logger, port}
}
