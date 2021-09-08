package rpc

import (
	"context"
	"github.com/morzhanov/go-realworld/internal/common/grpc/grpcserver"

	prpc "github.com/morzhanov/go-realworld/api/rpc/pictures"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/pictures/services"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PicturesRpcServer struct {
	prpc.UnimplementedPicturesServer
	*grpcserver.BaseGrpcServer
	picturesService *services.PictureService
	server          *grpc.Server
}

func (s *PicturesRpcServer) GetUserPictures(ctx context.Context, in *prpc.GetUserPicturesRequest) (*prpc.PicturesMessage, error) {
	span := s.PrepareContext(ctx)
	defer span.Finish()
	return s.picturesService.GetUserPictures(in.UserId)
}

func (s *PicturesRpcServer) GetUserPicture(ctx context.Context, in *prpc.GetUserPictureRequest) (*prpc.PictureMessage, error) {
	span := s.PrepareContext(ctx)
	defer span.Finish()
	return s.picturesService.GetUserPicture(in.UserId, in.PictureId)
}

func (s *PicturesRpcServer) CreateUserPicture(ctx context.Context, in *prpc.CreateUserPictureRequest) (*prpc.PictureMessage, error) {
	span := s.PrepareContext(ctx)
	defer span.Finish()
	return s.picturesService.CreateUserPicture(in)
}

func (s *PicturesRpcServer) DeleteUserPicture(ctx context.Context, in *prpc.DeleteUserPictureRequest) (*emptypb.Empty, error) {
	span := s.PrepareContext(ctx)
	defer span.Finish()
	err := s.picturesService.DeleteUserPicture(in.UserId, in.PictureId)
	return nil, err
}

func (s *PicturesRpcServer) Listen(ctx context.Context) error {
	return s.BaseGrpcServer.Listen(ctx, s.server)
}

func NewPicturesRpcServer(
	picturesService *services.PictureService,
	c *config.Config,
	tracer *opentracing.Tracer,
	logger *zap.Logger,
) (s *PicturesRpcServer) {
	bs := grpcserver.NewGrpcServer(tracer, logger, c.GrpcPort)
	s = &PicturesRpcServer{
		picturesService: picturesService,
		BaseGrpcServer:  bs,
	}
	prpc.RegisterPicturesServer(grpc.NewServer(), s)
	return
}
