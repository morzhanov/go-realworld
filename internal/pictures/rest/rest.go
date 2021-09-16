package rest

import (
	"context"
	"github.com/morzhanov/go-realworld/internal/common/rest/restcontroller"
	"net/http"

	"github.com/gin-gonic/gin"
	prpc "github.com/morzhanov/go-realworld/api/grpc/pictures"
	"github.com/morzhanov/go-realworld/internal/common/metrics"
	"github.com/morzhanov/go-realworld/internal/pictures/services"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type picturesRestController struct {
	*restcontroller.BaseRestController
	service *services.PictureService
}

type PicturesRestController interface {
	Listen(ctx context.Context, cancel context.CancelFunc, port string)
}

func (c *picturesRestController) handleCreateUserPicture(ctx *gin.Context) {
	input := prpc.CreateUserPictureRequest{}
	if err := c.ParseRestBody(ctx, &input); err != nil {
		c.HandleRestError(ctx, err)
		return
	}

	res, err := c.service.CreateUserPicture(&input)
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *picturesRestController) handleGetUserPictures(ctx *gin.Context) {
	userId := ctx.Param("userId")
	res, err := c.service.GetUserPictures(userId)
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *picturesRestController) handleGetUserPicture(ctx *gin.Context) {
	userId := ctx.Param("userId")
	id := ctx.Param("id")

	res, err := c.service.GetUserPicture(userId, id)
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}
	if res == nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *picturesRestController) handleDeleteUserPicture(ctx *gin.Context) {
	userId := ctx.Param("userId")
	id := ctx.Param("id")

	err := c.service.DeleteUserPicture(userId, id)
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (c *picturesRestController) Listen(
	ctx context.Context,
	cancel context.CancelFunc,
	port string,
) {
	c.BaseRestController.Listen(ctx, cancel, port)
}

func NewPicturesRestController(
	s *services.PictureService,
	tracer opentracing.Tracer,
	logger *zap.Logger,
	mc *metrics.MetricsCollector,
) PicturesRestController {
	bc := restcontroller.NewRestController(
		tracer,
		logger,
		mc,
	)
	c := picturesRestController{
		service:            s,
		BaseRestController: bc,
	}

	bc.Router.POST("/pictures", bc.Handler(c.handleCreateUserPicture))
	bc.Router.GET("/pictures/:userId", bc.Handler(c.handleGetUserPictures))
	bc.Router.GET("/pictures/:userId/:id", bc.Handler(c.handleGetUserPicture))
	bc.Router.DELETE("/pictures/:userId/:id", bc.Handler(c.handleDeleteUserPicture))
	return &c
}
