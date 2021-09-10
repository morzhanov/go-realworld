package rpc

import (
	"context"
	"fmt"
	"github.com/morzhanov/go-realworld/internal/common/grpc/grpcserver"

	urpc "github.com/morzhanov/go-realworld/api/rpc/users"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/users/services"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UsersRpcServer struct {
	urpc.UnimplementedUsersServer
	*grpcserver.BaseGrpcServer
	usersService *services.UsersService
	server       *grpc.Server
}

func (s *UsersRpcServer) GetUserData(ctx context.Context, in *urpc.GetUserDataRequest) (res *urpc.UserMessage, err error) {
	span := s.PrepareContext(ctx)
	defer span.Finish()
	return s.usersService.GetUserData(in.UserId)
}

func (s *UsersRpcServer) GetUserDataByUsername(ctx context.Context, in *urpc.GetUserDataByUsernameRequest) (res *urpc.UserMessage, err error) {
	span := s.PrepareContext(ctx)
	defer span.Finish()
	return s.usersService.GetUserDataByUsername(in.Username)
}

func (s *UsersRpcServer) ValidateUserPassword(ctx context.Context, in *urpc.ValidateUserPasswordRequest) (res *emptypb.Empty, err error) {
	span := s.PrepareContext(ctx)
	defer span.Finish()
	err = s.usersService.ValidateUserPassword(in)
	return res, err
}

func (s *UsersRpcServer) CreateUser(ctx context.Context, in *urpc.CreateUserRequest) (res *urpc.UserMessage, err error) {
	span := s.PrepareContext(ctx)
	defer span.Finish()
	return s.usersService.CreateUser(in)
}

func (s *UsersRpcServer) DeleteUser(ctx context.Context, in *urpc.DeleteUserRequest) (res *emptypb.Empty, err error) {
	span := s.PrepareContext(ctx)
	defer span.Finish()
	err = s.usersService.DeleteUser(in.UserId)
	return res, err
}

func (s *UsersRpcServer) Listen(ctx context.Context, cancel context.CancelFunc) {
	s.BaseGrpcServer.Listen(ctx, cancel, s.server)
}

func NewUsersRpcServer(
	usersService *services.UsersService,
	c *config.Config,
	tracer *opentracing.Tracer,
	logger *zap.Logger,
) (s *UsersRpcServer) {
	uri := fmt.Sprintf("%s:%s", c.GrpcAddr, c.GrpcPort)
	bs := grpcserver.NewGrpcServer(tracer, logger, uri)
	s = &UsersRpcServer{
		usersService:   usersService,
		BaseGrpcServer: bs,
		server:         grpc.NewServer(),
	}
	urpc.RegisterUsersServer(s.server, s)
	return
}
