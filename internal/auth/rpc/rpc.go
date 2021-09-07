package rpc

import (
	"context"

	arpc "github.com/morzhanov/go-realworld/api/rpc/auth"
	"github.com/morzhanov/go-realworld/internal/auth/services"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/helper"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/morzhanov/go-realworld/internal/common/tracing"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type AuthRpcServer struct {
	arpc.UnimplementedAuthServer
	authService *services.AuthService
	port        string
	server      *grpc.Server
	tracer      *opentracing.Tracer
}

func (s *AuthRpcServer) ValidateRpcRequest(ctx context.Context, in *arpc.ValidateRpcRequestInput) (res *arpc.ValidationResponse, err error) {
	span := tracing.StartSpanFromGrpcRequest(*s.tracer, ctx)
	defer span.Finish()
	return s.authService.ValidateRpcRequest(in)
}

func (s *AuthRpcServer) Login(ctx context.Context, in *arpc.LoginInput) (res *arpc.AuthResponse, err error) {
	span := tracing.StartSpanFromGrpcRequest(*s.tracer, ctx)
	defer span.Finish()
	ctx = context.WithValue(ctx, "transport", sender.RpcTransport)
	return s.authService.Login(ctx, in)
}

func (s *AuthRpcServer) Signup(ctx context.Context, in *arpc.SignupInput) (res *arpc.AuthResponse, err error) {
	span := tracing.StartSpanFromGrpcRequest(*s.tracer, ctx)
	defer span.Finish()
	ctx = context.WithValue(ctx, "transport", sender.RpcTransport)
	return s.authService.Signup(ctx, in)
}

func (s *AuthRpcServer) Listen(ctx context.Context, logger *zap.Logger) error {
	return helper.StartGrpcServer(ctx, s.server, s.port, logger)
}

func NewAuthRpcServer(
	authService *services.AuthService,
	c *config.Config,
	tracer *opentracing.Tracer,
) (server *AuthRpcServer) {
	server = &AuthRpcServer{authService: authService, port: c.GrpcPort, tracer: tracer}
	arpc.RegisterAuthServer(grpc.NewServer(), server)
	return
}
