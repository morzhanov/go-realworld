package rpc

import (
	"context"

	grpc "github.com/morzhanov/go-realworld/api/rpc/pictures"
	"github.com/morzhanov/go-realworld/internal/pictures/dto"
	"github.com/morzhanov/go-realworld/internal/pictures/services"
	core_grpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PicturesRpcServer struct {
	grpc.UnimplementedPicturesServer
	picturesService *services.PictureService
}

func (s *PicturesRpcServer) GetUserPictures(ctx context.Context, in *grpc.GetUserPicturesRequest) (res *grpc.PicturesMessage, err error) {
	dto, err := s.picturesService.GetUserPictures(in.UserId)

	pictures := make([]*grpc.PictureMessage, len(dto))
	for i, picture := range dto {
		pictures[i] = &grpc.PictureMessage{
			Id:     picture.ID,
			Title:  picture.Title,
			Base64: picture.Base64,
			UserId: picture.UserId,
		}
	}

	res = &grpc.PicturesMessage{Pictures: pictures}
	return res, err
}

func (s *PicturesRpcServer) GetUserPicture(ctx context.Context, in *grpc.GetUserPictureRequest) (res *grpc.PictureMessage, err error) {
	dto, err := s.picturesService.GetUserPicture(in.UserId, in.PictureId)

	res = &grpc.PictureMessage{
		Id:     dto.ID,
		Title:  dto.Title,
		Base64: dto.Base64,
		UserId: dto.UserId,
	}
	return res, err
}

func (s *PicturesRpcServer) CreateUserPicture(ctx context.Context, in *grpc.CreateUserPictureRequest) (res *grpc.PictureMessage, err error) {
	dto, err := s.picturesService.CreateUserPicture(
		in.UserId,
		&dto.CreatePicturesDto{Title: in.Title, Base64: in.Base64},
	)

	res = &grpc.PictureMessage{
		Id:     dto.ID,
		Title:  dto.Title,
		Base64: dto.Base64,
		UserId: dto.UserId,
	}
	return res, err
}

func (s *PicturesRpcServer) DeleteUserPicture(ctx context.Context, in *grpc.DeleteUserPictureRequest) (res *emptypb.Empty, err error) {
	err = s.picturesService.DeleteUserPicture(in.UserId, in.PictureId)
	return res, err
}

func NewPicturesRpcService(picturesService services.PictureService) (s *core_grpc.Server) {
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
	grpc.RegisterPicturesServer(s, &PicturesRpcServer{picturesService: &picturesService})
	return s
}
