package rpc

import (
	"context"

	anrpc "github.com/morzhanov/go-realworld/api/rpc/analytics"
	"github.com/morzhanov/go-realworld/internal/analytics/services"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AnalyticsRpcServer struct {
	anrpc.UnimplementedAnalyticsServer
	analyticsService *services.AnalyticsService
}

func (s *AnalyticsRpcServer) LogData(ctx context.Context, in *anrpc.LogDataRequest) (res *emptypb.Empty, err error) {
	err = s.analyticsService.LogData(in)
	return res, err
}

func (s *AnalyticsRpcServer) GetLog(ctx context.Context, in *anrpc.GetLogRequest) (res *anrpc.AnalyticsEntryMessage, err error) {
	return s.analyticsService.GetLog(in)
}

func NewAnalyticsRpcService(analyticsService services.AnalyticsService) (s *grpc.Server) {
	// TODO: get port from env vars
	// TODO: call those lines in the main app to start rpc server
	// port := ":5000"

	// lis, err := net.Listen("tcp", port)
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }

	// if err := s.Serve(lis); err != nil {
	// 	log.Fatalf("failed to serve: %v", err)
	// }
	// log.Printf("Server started at: localhost%v", port)

	s = grpc.NewServer()
	anrpc.RegisterAnalyticsServer(s, &AnalyticsRpcServer{analyticsService: &analyticsService})
	return s
}
