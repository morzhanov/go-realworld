package rpc

import (
	"context"

	arpc "github.com/morzhanov/go-realworld/api/rpc/auth"
	"github.com/morzhanov/go-realworld/internal/auth/services"
	"github.com/morzhanov/go-realworld/internal/common/helper"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/spf13/viper"
	grpc "google.golang.org/grpc"
)

type AuthRpcServer struct {
	arpc.UnimplementedAuthServer
	authService *services.AuthService
	server      *grpc.Server
}

func (s *AuthRpcServer) ValidateRpcRequest(ctx context.Context, in *arpc.ValidateRpcRequestInput) (res *arpc.ValidationResponse, err error) {
	return s.authService.ValidateRpcRequest(in)
}

func (s *AuthRpcServer) Login(ctx context.Context, in *arpc.LoginInput) (res *arpc.AuthResponse, err error) {
	ctx = context.WithValue(ctx, "transport", sender.RpcTransport)
	return s.authService.Login(ctx, in)
}

func (s *AuthRpcServer) Signup(ctx context.Context, in *arpc.SignupInput) (res *arpc.AuthResponse, err error) {
	ctx = context.WithValue(ctx, "transport", sender.RpcTransport)
	return s.authService.Signup(ctx, in)
}

func (s *AuthRpcServer) Listen() error {
	port := viper.GetString("AUTH_GRPC_PORT")
	return helper.StartGrpcServer(s.server, port)
}

func NewAuthRpcService(authService services.AuthService) (s *grpc.Server) {
	s = grpc.NewServer()
	arpc.RegisterAuthServer(s, &AuthRpcServer{authService: &authService, server: s})
	return s
}
