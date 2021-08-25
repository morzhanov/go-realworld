package rpc

import (
	"context"

	prpc "github.com/morzhanov/go-realworld/api/rpc/pictures"
	"github.com/morzhanov/go-realworld/internal/pictures/services"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PicturesRpcServer struct {
	prpc.UnimplementedPicturesServer
	picturesService *services.PictureService
}

func (s *PicturesRpcServer) GetUserPictures(ctx context.Context, in *prpc.GetUserPicturesRequest) (res *prpc.PicturesMessage, err error) {
	return s.picturesService.GetUserPictures(in.UserId)
}

func (s *PicturesRpcServer) GetUserPicture(ctx context.Context, in *prpc.GetUserPictureRequest) (res *prpc.PictureMessage, err error) {
	return s.picturesService.GetUserPicture(in.UserId, in.PictureId)
}

func (s *PicturesRpcServer) CreateUserPicture(ctx context.Context, in *prpc.CreateUserPictureRequest) (res *prpc.PictureMessage, err error) {
	return s.picturesService.CreateUserPicture(in)
}

func (s *PicturesRpcServer) DeleteUserPicture(ctx context.Context, in *prpc.DeleteUserPictureRequest) (res *emptypb.Empty, err error) {
	err = s.picturesService.DeleteUserPicture(in.UserId, in.PictureId)
	return res, err
}

func NewPicturesRpcService(picturesService services.PictureService) (s *grpc.Server) {
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
	prpc.RegisterPicturesServer(s, &PicturesRpcServer{picturesService: &picturesService})
	return s
}
