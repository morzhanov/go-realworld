package rpc

import (
	"context"

	grpc "github.com/morzhanov/go-realworld/api/rpc/auth"
	"github.com/morzhanov/go-realworld/internal/auth/services"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	core_grpc "google.golang.org/grpc"
)

type AuthRpcServer struct {
	grpc.UnimplementedAuthServer
	authService *services.AuthService
}

func (s *AuthRpcServer) ValidateRpcRequest(ctx context.Context, in *grpc.ValidateRpcRequestInput) (res *grpc.ValidationResponse, err error) {
	return s.authService.ValidateRpcRequest(in)
}

func (s *AuthRpcServer) Login(ctx context.Context, in *grpc.LoginInput) (res *grpc.AuthResponse, err error) {
	ctx = context.WithValue(ctx, "transport", sender.RpcTransport)
	return s.authService.Login(ctx, in)
}

func (s *AuthRpcServer) Signup(ctx context.Context, in *grpc.SignupInput) (res *grpc.AuthResponse, err error) {
	ctx = context.WithValue(ctx, "transport", sender.RpcTransport)
	return s.authService.Signup(ctx, in)
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
