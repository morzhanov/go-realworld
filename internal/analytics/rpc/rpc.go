package rpc

import (
	"context"

	grpc "github.com/morzhanov/go-realworld/api/rpc/analytics"
	"github.com/morzhanov/go-realworld/internal/analytics/dto"
	"github.com/morzhanov/go-realworld/internal/analytics/models"
	"github.com/morzhanov/go-realworld/internal/analytics/services"
	core_grpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AnalyticsRpcServer struct {
	grpc.UnimplementedAnalyticsServer
	analyticsService *services.AnalyticsService
}

func (s *AnalyticsRpcServer) LogData(ctx context.Context, in *grpc.LogDataRequest) (res *emptypb.Empty, err error) {
	err = s.analyticsService.LogData(&models.AnalyticsEntry{
		ID:        in.Id,
		UserID:    in.UserId,
		Operation: in.Operation,
		Data:      in.Data,
	})

	return res, err
}

func (s *AnalyticsRpcServer) GetLog(ctx context.Context, in *grpc.GetLogRequest) (res *grpc.AnalyticsEntryMessage, err error) {
	entry, err := s.analyticsService.GetLog(&dto.GetLogsInput{Offset: int(in.Offset)})

	res = &grpc.AnalyticsEntryMessage{
		Id:        entry.ID,
		UserId:    entry.UserID,
		Operation: entry.Operation,
		Data:      entry.Data,
	}
	return res, err
}

func NewAnalyticsRpcService(analyticsService services.AnalyticsService) (s *core_grpc.Server) {
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

	s = core_grpc.NewServer()
	grpc.RegisterAnalyticsServer(s, &AnalyticsRpcServer{analyticsService: &analyticsService})
	return s
}
