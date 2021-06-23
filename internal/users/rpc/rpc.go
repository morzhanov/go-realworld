package rpc

import (
	"context"

	grpc "github.com/morzhanov/go-realworld/api/rpc/users"
	"github.com/morzhanov/go-realworld/internal/users/dto"
	"github.com/morzhanov/go-realworld/internal/users/services"
	core_grpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UsersRpcServer struct {
	grpc.UnimplementedUsersServer
	usersService *services.UsersService
}

func (s *UsersRpcServer) GetUserData(ctx context.Context, in *grpc.GetUserDataRequest) (res *grpc.UserMessage, err error) {
	dto, err := s.usersService.GetUserData(in.UserId)
	res = &grpc.UserMessage{Id: dto.ID, Username: dto.Username}
	return res, err
}

func (s *UsersRpcServer) GetUserDataByUsername(ctx context.Context, in *grpc.GetUserDataByUsernameRequest) (res *grpc.UserMessage, err error) {
	dto, err := s.usersService.GetUserDataByUsername(in.Username)
	res = &grpc.UserMessage{Id: dto.ID, Username: dto.Username}
	return res, err
}

func (s *UsersRpcServer) ValidateUserPassword(ctx context.Context, in *grpc.ValidateUserPasswordRequest) (res *emptypb.Empty, err error) {
	err = s.usersService.ValidateUserPassword(&dto.ValidateUserPasswordDto{Username: in.Username, Password: in.Password})
	return res, err
}

func (s *UsersRpcServer) CreateUser(ctx context.Context, in *grpc.CreateUserRequest) (res *grpc.UserMessage, err error) {
	dto, err := s.usersService.CreateUser(&dto.CreateUserDto{Username: in.Username, Password: in.Password})
	res = &grpc.UserMessage{Id: dto.ID, Username: dto.Username}
	return res, err
}

func (s *UsersRpcServer) DeleteUser(ctx context.Context, in *grpc.DeleteUserRequest) (res *emptypb.Empty, err error) {
	err = s.usersService.DeleteUser(in.UserId)
	return res, err
}

func NewUsersRpcService(usersService services.UsersService) (s *core_grpc.Server) {
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
	grpc.RegisterUsersServer(s, &UsersRpcServer{usersService: &usersService})
	return s
}
