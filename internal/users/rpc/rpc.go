package rpc

import (
	"context"

	urpc "github.com/morzhanov/go-realworld/api/rpc/users"
	"github.com/morzhanov/go-realworld/internal/users/services"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UsersRpcServer struct {
	urpc.UnimplementedUsersServer
	usersService *services.UsersService
}

func (s *UsersRpcServer) GetUserData(ctx context.Context, in *urpc.GetUserDataRequest) (res *urpc.UserMessage, err error) {
	return s.usersService.GetUserData(in.UserId)
}

func (s *UsersRpcServer) GetUserDataByUsername(ctx context.Context, in *urpc.GetUserDataByUsernameRequest) (res *urpc.UserMessage, err error) {
	return s.usersService.GetUserDataByUsername(in.Username)
}

func (s *UsersRpcServer) ValidateUserPassword(ctx context.Context, in *urpc.ValidateUserPasswordRequest) (res *emptypb.Empty, err error) {
	err = s.usersService.ValidateUserPassword(in)
	return res, err
}

func (s *UsersRpcServer) CreateUser(ctx context.Context, in *urpc.CreateUserRequest) (res *urpc.UserMessage, err error) {
	return s.usersService.CreateUser(in)
}

func (s *UsersRpcServer) DeleteUser(ctx context.Context, in *urpc.DeleteUserRequest) (res *emptypb.Empty, err error) {
	err = s.usersService.DeleteUser(in.UserId)
	return res, err
}

func NewUsersRpcService(usersService services.UsersService) (s *grpc.Server) {
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
	urpc.RegisterUsersServer(s, &UsersRpcServer{usersService: &usersService})
	return s
}
