package grpcserver

import (
	"context"
	"fmt"
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

func (s *BaseGrpcServer) Listen(ctx context.Context, server *grpc.Server) error {
	lis, err := net.Listen("tcp", s.Port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	if err := server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	log.Info("Grpc server started", zap.String("port", s.Port))
	<-ctx.Done()
	return lis.Close()
}

func NewGrpcServer(
	tracer *opentracing.Tracer,
	logger *zap.Logger,
	port string,
) *BaseGrpcServer {
	return &BaseGrpcServer{tracer, logger, port}
}
