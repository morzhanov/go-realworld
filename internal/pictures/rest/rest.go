package rest

import (
	"context"
	"github.com/morzhanov/go-realworld/internal/common/rest/restcontroller"
	"net/http"

	"github.com/gin-gonic/gin"
	prpc "github.com/morzhanov/go-realworld/api/rpc/pictures"
	"github.com/morzhanov/go-realworld/internal/common/metrics"
	"github.com/morzhanov/go-realworld/internal/pictures/services"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type PicturesRestController struct {
	*restcontroller.BaseRestController
	service *services.PictureService
}

func (c *PicturesRestController) handleCreateUserPicture(ctx *gin.Context) {
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

func (c *PicturesRestController) handleGetUserPictures(ctx *gin.Context) {
	userId := ctx.Param("userId")
	res, err := c.service.GetUserPictures(userId)
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *PicturesRestController) handleGetUserPicture(ctx *gin.Context) {
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

func (c *PicturesRestController) handleDeleteUserPicture(ctx *gin.Context) {
	userId := ctx.Param("userId")
	id := ctx.Param("id")

	err := c.service.DeleteUserPicture(userId, id)
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (c *PicturesRestController) Listen(
	ctx context.Context,
	cancel context.CancelFunc,
	port string,
) {
	c.BaseRestController.Listen(ctx, cancel, port)
}

func NewPicturesRestController(
	s *services.PictureService,
	tracer *opentracing.Tracer,
	logger *zap.Logger,
	mc *metrics.MetricsCollector,
) *PicturesRestController {
	bc := restcontroller.NewRestController(
		tracer,
		logger,
		mc,
	)
	c := PicturesRestController{
		service:            s,
		BaseRestController: bc,
	}

	bc.Router.POST("/pictures", bc.Handler(c.handleCreateUserPicture))
	bc.Router.GET("/pictures/:userId", bc.Handler(c.handleGetUserPictures))
	bc.Router.GET("/pictures/:userId/:id", bc.Handler(c.handleGetUserPicture))
	bc.Router.DELETE("/pictures/:userId/:id", bc.Handler(c.handleDeleteUserPicture))
	return &c
}
