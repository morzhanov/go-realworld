package rest

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	anrpc "github.com/morzhanov/go-realworld/api/rpc/analytics"
	"github.com/morzhanov/go-realworld/internal/analytics/services"
	"github.com/morzhanov/go-realworld/internal/common/helper"
	"github.com/morzhanov/go-realworld/internal/common/tracing"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type AnalyticsRestController struct {
	service *services.AnalyticsService
	router  *gin.Engine
	tracer  *opentracing.Tracer
}

func (c *AnalyticsRestController) handleLogData(ctx *gin.Context) {
	// TODO: maybe we should somehow generalize this step via middleware or something
	// TODO: because we'are using the same code for tracing in all controllers
	span := tracing.StartSpanFromHttpRequest(*c.tracer, ctx.Request)
	defer span.Finish()

	input := anrpc.LogDataRequest{}
	if err := helper.ParseRestBody(ctx, &input); err != nil {
		helper.HandleRestError(ctx, err)
		return
	}

	if err := c.service.LogData(&input); err != nil {
		helper.HandleRestError(ctx, err)
		return
	}
	ctx.Status(http.StatusCreated)
}

func (c *AnalyticsRestController) handleGetData(ctx *gin.Context) {
	span := tracing.StartSpanFromHttpRequest(*c.tracer, ctx.Request)
	defer span.Finish()

	offset, err := strconv.Atoi(ctx.Param("offset"))
	if err != nil {
		helper.HandleRestError(ctx, err)
	}

	res, err := c.service.GetLog(&anrpc.GetLogRequest{Offset: int32(offset)})
	if err != nil {
		helper.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *AnalyticsRestController) Listen(ctx context.Context, port string, logger *zap.Logger) {
	helper.StartRestServer(ctx, port, c.router, logger)
}

func NewAnalyticsRestController(s *services.AnalyticsService, tracer *opentracing.Tracer) (c *AnalyticsRestController) {
	router := gin.Default()
	c = &AnalyticsRestController{s, router, tracer}

	router.GET("/analytics", c.handleLogData)
	router.GET("/analytics/:offset", c.handleGetData)
	return
}
