package rpc

import (
	"context"

	arpc "github.com/morzhanov/go-realworld/api/rpc/auth"
	"github.com/morzhanov/go-realworld/internal/auth/services"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/helper"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	grpc "google.golang.org/grpc"
)

type AuthRpcServer struct {
	arpc.UnimplementedAuthServer
	authService *services.AuthService
	port        string
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

func (s *AuthRpcServer) Listen(ctx context.Context) error {
	return helper.StartGrpcServer(ctx, s.server, s.port)
}

func NewAuthRpcServer(
	authService *services.AuthService,
	c *config.Config,
) (server *AuthRpcServer) {
	server = &AuthRpcServer{authService: authService, port: c.AuthGrpcPort}
	arpc.RegisterAuthServer(grpc.NewServer(), server)
	return
}
