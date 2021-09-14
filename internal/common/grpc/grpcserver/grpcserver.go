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
	Uri    string
}

func (s *BaseGrpcServer) PrepareContext(ctx context.Context) (context.Context, opentracing.Span) {
	span := tracing.StartSpanFromGrpcRequest(*s.Tracer, ctx)
	ctx = context.WithValue(ctx, "transport", sender.RpcTransport)
	return ctx, span
}

func (s *BaseGrpcServer) Listen(ctx context.Context, cancel context.CancelFunc, server *grpc.Server) {
	lis, err := net.Listen("tcp", s.Uri)
	if err != nil {
		cancel()
		helper.HandleInitializationError(err, "grpc server", s.Logger)
		return
	}

	if err := server.Serve(lis); err != nil {
		cancel()
		helper.HandleInitializationError(err, "grpc server", s.Logger)
		return
	}
	log.Info("Grpc server started", zap.String("port", s.Uri))
	<-ctx.Done()
	if err := lis.Close(); err != nil {
		cancel()
		helper.HandleInitializationError(err, "grpc server", s.Logger)
		return
	}
}

func NewGrpcServer(
	tracer *opentracing.Tracer,
	logger *zap.Logger,
	uri string,
) *BaseGrpcServer {
	return &BaseGrpcServer{tracer, logger, uri}
}
