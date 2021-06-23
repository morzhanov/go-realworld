package rpc

import (
	"context"

	grpc "github.com/morzhanov/go-realworld/api/rpc/auth"
	"github.com/morzhanov/go-realworld/internal/auth/dto"
	"github.com/morzhanov/go-realworld/internal/auth/services"
	core_grpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthRpcServer struct {
	grpc.UnimplementedAuthServer
	authService *services.AuthService
}

func (s *AuthRpcServer) ValidateRpcRequest(ctx context.Context, in *grpc.ValidateRpcRequestInput) (res *emptypb.Empty, err error) {
	err = s.authService.ValidateRpcRequest(&dto.ValidateRpcRequestInput{AccessToken: in.AccessToken, Method: in.Method})
	return res, err
}

func NewAuthRpcService(authService services.AuthService) (s *core_grpc.Server) {
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
	grpc.RegisterAuthServer(s, &AuthRpcServer{authService: &authService})
	return s
}
