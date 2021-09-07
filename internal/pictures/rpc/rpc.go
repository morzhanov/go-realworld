package rpc

import (
	"context"

	prpc "github.com/morzhanov/go-realworld/api/rpc/pictures"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/common/helper"
	"github.com/morzhanov/go-realworld/internal/common/tracing"
	"github.com/morzhanov/go-realworld/internal/pictures/services"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PicturesRpcServer struct {
	prpc.UnimplementedPicturesServer
	picturesService *services.PictureService
	port            string
	server          *grpc.Server
	tracer          *opentracing.Tracer
}

func (s *PicturesRpcServer) GetUserPictures(ctx context.Context, in *prpc.GetUserPicturesRequest) (res *prpc.PicturesMessage, err error) {
	span := tracing.StartSpanFromGrpcRequest(*s.tracer, ctx)
	defer span.Finish()
	return s.picturesService.GetUserPictures(in.UserId)
}

func (s *PicturesRpcServer) GetUserPicture(ctx context.Context, in *prpc.GetUserPictureRequest) (res *prpc.PictureMessage, err error) {
	span := tracing.StartSpanFromGrpcRequest(*s.tracer, ctx)
	defer span.Finish()
	return s.picturesService.GetUserPicture(in.UserId, in.PictureId)
}

func (s *PicturesRpcServer) CreateUserPicture(ctx context.Context, in *prpc.CreateUserPictureRequest) (res *prpc.PictureMessage, err error) {
	span := tracing.StartSpanFromGrpcRequest(*s.tracer, ctx)
	defer span.Finish()
	return s.picturesService.CreateUserPicture(in)
}

func (s *PicturesRpcServer) DeleteUserPicture(ctx context.Context, in *prpc.DeleteUserPictureRequest) (res *emptypb.Empty, err error) {
	span := tracing.StartSpanFromGrpcRequest(*s.tracer, ctx)
	defer span.Finish()
	err = s.picturesService.DeleteUserPicture(in.UserId, in.PictureId)
	return res, err
}

func (s *PicturesRpcServer) Listen(ctx context.Context, logger *zap.Logger) error {
	return helper.StartGrpcServer(ctx, s.server, s.port, logger)
}

func NewPicturesRpcServer(
	picturesService *services.PictureService,
	c *config.Config,
	tracer *opentracing.Tracer,

) (server *PicturesRpcServer) {
	server = &PicturesRpcServer{picturesService: picturesService, port: c.PicturesRestPort, tracer: tracer}
	prpc.RegisterPicturesServer(grpc.NewServer(), server)
	return
}
