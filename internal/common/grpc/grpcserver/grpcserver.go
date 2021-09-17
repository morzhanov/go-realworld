package grpcserver

import (
	"context"
	"github.com/morzhanov/go-realworld/internal/common/errors"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/morzhanov/go-realworld/internal/common/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/common/log"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type baseGrpcServer struct {
	Tracer opentracing.Tracer
	Logger *zap.Logger
	Uri    string
}

type BaseGrpcServer interface {
	PrepareContext(ctx context.Context) (context.Context, opentracing.Span)
	Listen(ctx context.Context, cancel context.CancelFunc, server *grpc.Server)
}

func (s *baseGrpcServer) PrepareContext(ctx context.Context) (context.Context, opentracing.Span) {
	span := tracing.StartSpanFromGrpcRequest(s.Tracer, ctx)
	ctx = context.WithValue(ctx, "transport", sender.RpcTransport)
	return ctx, span
}

func (s *baseGrpcServer) Listen(ctx context.Context, cancel context.CancelFunc, server *grpc.Server) {
	lis, err := net.Listen("tcp", s.Uri)
	if err != nil {
		cancel()
		errors.LogInitializationError(err, "grpc server", s.Logger)
		return
	}

	if err := server.Serve(lis); err != nil {
		cancel()
		errors.LogInitializationError(err, "grpc server", s.Logger)
		return
	}
	log.Info("Grpc server started", zap.String("port", s.Uri))
	<-ctx.Done()
	if err := lis.Close(); err != nil {
		cancel()
		errors.LogInitializationError(err, "grpc server", s.Logger)
		return
	}
}

func NewGrpcServer(
	tracer opentracing.Tracer,
	logger *zap.Logger,
	uri string,
) BaseGrpcServer {
	return &baseGrpcServer{tracer, logger, uri}
}
