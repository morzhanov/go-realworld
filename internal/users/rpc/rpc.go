package rpc

import (
	"context"

	urpc "github.com/morzhanov/go-realworld/api/rpc/users"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/helper"
	"github.com/morzhanov/go-realworld/internal/common/tracing"
	"github.com/morzhanov/go-realworld/internal/users/services"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UsersRpcServer struct {
	urpc.UnimplementedUsersServer
	usersService *services.UsersService
	port         string
	server       *grpc.Server
	tracer       *opentracing.Tracer
}

func (s *UsersRpcServer) GetUserData(ctx context.Context, in *urpc.GetUserDataRequest) (res *urpc.UserMessage, err error) {
	span := tracing.StartSpanFromGrpcRequest(*s.tracer, ctx)
	defer span.Finish()
	return s.usersService.GetUserData(in.UserId)
}

func (s *UsersRpcServer) GetUserDataByUsername(ctx context.Context, in *urpc.GetUserDataByUsernameRequest) (res *urpc.UserMessage, err error) {
	span := tracing.StartSpanFromGrpcRequest(*s.tracer, ctx)
	defer span.Finish()
	return s.usersService.GetUserDataByUsername(in.Username)
}

func (s *UsersRpcServer) ValidateUserPassword(ctx context.Context, in *urpc.ValidateUserPasswordRequest) (res *emptypb.Empty, err error) {
	span := tracing.StartSpanFromGrpcRequest(*s.tracer, ctx)
	defer span.Finish()
	err = s.usersService.ValidateUserPassword(in)
	return res, err
}

func (s *UsersRpcServer) CreateUser(ctx context.Context, in *urpc.CreateUserRequest) (res *urpc.UserMessage, err error) {
	span := tracing.StartSpanFromGrpcRequest(*s.tracer, ctx)
	defer span.Finish()
	return s.usersService.CreateUser(in)
}

func (s *UsersRpcServer) DeleteUser(ctx context.Context, in *urpc.DeleteUserRequest) (res *emptypb.Empty, err error) {
	span := tracing.StartSpanFromGrpcRequest(*s.tracer, ctx)
	defer span.Finish()
	err = s.usersService.DeleteUser(in.UserId)
	return res, err
}

func (s *UsersRpcServer) Listen(ctx context.Context, logger *zap.Logger) error {
	return helper.StartGrpcServer(ctx, s.server, s.port, logger)
}

func NewUsersRpcServer(
	usersService *services.UsersService,
	c *config.Config,
	tracer *opentracing.Tracer,
) (server *UsersRpcServer) {
	server = &UsersRpcServer{usersService: usersService, port: c.AnalyticsGrpcPort, tracer: tracer}
	urpc.RegisterUsersServer(grpc.NewServer(), server)
	return
}
