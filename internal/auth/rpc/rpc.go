package rpc

import (
	"context"
	"github.com/morzhanov/go-realworld/internal/common/grpc/grpcserver"
	"go.uber.org/zap"

	arpc "github.com/morzhanov/go-realworld/api/rpc/auth"
	"github.com/morzhanov/go-realworld/internal/auth/services"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

type AuthRpcServer struct {
	arpc.UnimplementedAuthServer
	*grpcserver.BaseGrpcServer
	authService *services.AuthService
	server      *grpc.Server
}

func (s *AuthRpcServer) ValidateRpcRequest(ctx context.Context, in *arpc.ValidateRpcRequestInput) (*arpc.ValidationResponse, error) {
	span := s.PrepareContext(ctx)
	defer span.Finish()
	return s.authService.ValidateRpcRequest(in)
}

func (s *AuthRpcServer) Login(ctx context.Context, in *arpc.LoginInput) (*arpc.AuthResponse, error) {
	span := s.PrepareContext(ctx)
	defer span.Finish()
	return s.authService.Login(ctx, in, &span)
}

func (s *AuthRpcServer) Signup(ctx context.Context, in *arpc.SignupInput) (res *arpc.AuthResponse, err error) {
	span := s.PrepareContext(ctx)
	defer span.Finish()
	return s.authService.Signup(ctx, in, &span)
}

func (s *AuthRpcServer) Listen(ctx context.Context, cancel context.CancelFunc) {
	s.BaseGrpcServer.Listen(ctx, cancel, s.server)
}

func NewAuthRpcServer(
	authService *services.AuthService,
	c *config.Config,
	tracer *opentracing.Tracer,
	logger *zap.Logger,
) (s *AuthRpcServer) {
	bs := grpcserver.NewGrpcServer(tracer, logger, c.GrpcPort)
	s = &AuthRpcServer{
		authService:    authService,
		BaseGrpcServer: bs,
	}
	arpc.RegisterAuthServer(grpc.NewServer(), s)
	return
}
