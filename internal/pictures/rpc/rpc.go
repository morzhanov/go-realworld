package rpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/morzhanov/go-realworld/internal/common/grpc/grpcserver"
	"reflect"

	prpc "github.com/morzhanov/go-realworld/api/rpc/pictures"
	"github.com/morzhanov/go-realworld/internal/common/config"
	"github.com/morzhanov/go-realworld/internal/pictures/services"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PicturesRpcServer struct {
	prpc.UnimplementedPicturesServer
	*grpcserver.BaseGrpcServer
	picturesService *services.PictureService
	server          *grpc.Server
}

func (s *PicturesRpcServer) GetUserPictures(ctx context.Context, in *prpc.GetUserPicturesRequest) (*prpc.PicturesMessage, error) {
	ctx, span := s.PrepareContext(ctx)
	defer span.Finish()
	return s.picturesService.GetUserPictures(in.UserId)
}

func (s *PicturesRpcServer) GetUserPicture(ctx context.Context, in *prpc.GetUserPictureRequest) (*prpc.PictureMessage, error) {
	ctx, span := s.PrepareContext(ctx)
	defer span.Finish()
	res, err := s.picturesService.GetUserPicture(in.UserId, in.PictureId)
	if err == nil && reflect.ValueOf(res).IsNil() {
		return nil, errors.New("picture not found")
	}
	return res, err
}

func (s *PicturesRpcServer) CreateUserPicture(ctx context.Context, in *prpc.CreateUserPictureRequest) (*prpc.PictureMessage, error) {
	ctx, span := s.PrepareContext(ctx)
	defer span.Finish()
	return s.picturesService.CreateUserPicture(in)
}

func (s *PicturesRpcServer) DeleteUserPicture(ctx context.Context, in *prpc.DeleteUserPictureRequest) (*emptypb.Empty, error) {
	ctx, span := s.PrepareContext(ctx)
	defer span.Finish()
	err := s.picturesService.DeleteUserPicture(in.UserId, in.PictureId)
	return &emptypb.Empty{}, err
}

func (s *PicturesRpcServer) Listen(ctx context.Context, cancel context.CancelFunc) {
	s.BaseGrpcServer.Listen(ctx, cancel, s.server)
}

func NewPicturesRpcServer(
	picturesService *services.PictureService,
	c *config.Config,
	tracer *opentracing.Tracer,
	logger *zap.Logger,
) (s *PicturesRpcServer) {
	uri := fmt.Sprintf("%s:%s", c.GrpcAddr, c.GrpcPort)
	bs := grpcserver.NewGrpcServer(tracer, logger, uri)
	s = &PicturesRpcServer{
		picturesService: picturesService,
		BaseGrpcServer:  bs,
		server:          grpc.NewServer(),
	}
	prpc.RegisterPicturesServer(s.server, s)
	reflection.Register(s.server)
	return
}
