package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	prpc "github.com/morzhanov/go-realworld/api/rpc/pictures"
	"github.com/morzhanov/go-realworld/internal/common/helper"
	"github.com/morzhanov/go-realworld/internal/common/metrics"
	"github.com/morzhanov/go-realworld/internal/common/tracing"
	"github.com/morzhanov/go-realworld/internal/pictures/services"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type PicturesRestController struct {
	service *services.PictureService
	router  *gin.Engine
	tracer  *opentracing.Tracer
}

func (c *PicturesRestController) handleCreateUserPicture(ctx *gin.Context) {
	span := tracing.StartSpanFromHttpRequest(*c.tracer, ctx.Request)
	defer span.Finish()

	input := prpc.CreateUserPictureRequest{}
	if err := helper.ParseRestBody(ctx, &input); err != nil {
		helper.HandleRestError(ctx, err)
		return
	}
	userId := ctx.Param("userId")
	input.UserId = userId

	res, err := c.service.CreateUserPicture(&input)
	if err != nil {
		helper.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *PicturesRestController) handleGetUserPictures(ctx *gin.Context) {
	span := tracing.StartSpanFromHttpRequest(*c.tracer, ctx.Request)
	defer span.Finish()

	userId := ctx.Param("userId")
	res, err := c.service.GetUserPictures(userId)
	if err != nil {
		helper.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *PicturesRestController) handleGetUserPicture(ctx *gin.Context) {
	span := tracing.StartSpanFromHttpRequest(*c.tracer, ctx.Request)
	defer span.Finish()

	userId := ctx.Param("userId")
	id := ctx.Param("id")

	res, err := c.service.GetUserPicture(userId, id)
	if err != nil {
		helper.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *PicturesRestController) handleDeleteUserPicture(ctx *gin.Context) {
	span := tracing.StartSpanFromHttpRequest(*c.tracer, ctx.Request)
	defer span.Finish()

	userId := ctx.Param("userId")
	id := ctx.Param("id")

	err := c.service.DeleteUserPicture(userId, id)
	if err != nil {
		helper.HandleRestError(ctx, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (c *PicturesRestController) Listen(ctx context.Context, port string, logger *zap.Logger) {
	helper.StartRestServer(ctx, port, c.router, logger)
}

func NewPicturesRestController(s *services.PictureService, tracer *opentracing.Tracer, mc *metrics.MetricsCollector) *PicturesRestController {
	router := gin.Default()
	c := PicturesRestController{s, router, tracer}

	router.POST("/pictures", c.handleCreateUserPicture)
	router.GET("/pictures/:userId", c.handleGetUserPictures)
	router.GET("/pictures/:userId/:id", c.handleGetUserPicture)
	router.DELETE("/pictures/:userId/:id", c.handleDeleteUserPicture)
	mc.RegisterMetricsEndpoint(router)
	return &c
}
