package rpc

import (
	"context"

	prpc "github.com/morzhanov/go-realworld/api/rpc/pictures"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/helper"
	"github.com/morzhanov/go-realworld/internal/pictures/services"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PicturesRpcServer struct {
	prpc.UnimplementedPicturesServer
	picturesService *services.PictureService
	port            string
	server          *grpc.Server
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

func (s *PicturesRpcServer) Listen(ctx context.Context) error {
	return helper.StartGrpcServer(ctx, s.server, s.port)
}

func NewPicturesRpcServer(
	picturesService *services.PictureService,
	c *config.Config,
) (server *PicturesRpcServer) {
	server = &PicturesRpcServer{picturesService: picturesService, port: c.PicturesRestPort}
	prpc.RegisterPicturesServer(grpc.NewServer(), server)
	return
}
