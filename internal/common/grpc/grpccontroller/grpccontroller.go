package grpccontroller

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func StartGrpcServer(ctx context.Context, s *grpc.Server, port string, log *zap.Logger) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	log.Info("Grpc server started", zap.String("port", port))
	<-ctx.Done()
	lis.Close()
	return nil
}

// TODO: create base rest controller with span injection for methods
