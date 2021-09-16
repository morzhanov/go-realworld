package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	anrpc "github.com/morzhanov/go-realworld/api/grpc/analytics"
	"github.com/morzhanov/go-realworld/internal/analytics/services"
	"github.com/morzhanov/go-realworld/internal/common/metrics"
	"github.com/morzhanov/go-realworld/internal/common/rest/restcontroller"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
)

type AnalyticsRestController struct {
	*restcontroller.BaseRestController
	service *services.AnalyticsService
}

func (c *AnalyticsRestController) handleLogData(ctx *gin.Context) {
	input := anrpc.LogDataRequest{}
	if err := c.ParseRestBody(ctx, &input); err != nil {
		c.HandleRestError(ctx, err)
		return
	}
	if err := c.service.LogData(&input); err != nil {
		c.HandleRestError(ctx, err)
		return
	}
	ctx.Status(http.StatusCreated)
}

func (c *AnalyticsRestController) handleGetData(ctx *gin.Context) {
	res, err := c.service.GetLog(&emptypb.Empty{})
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *AnalyticsRestController) Listen(
	ctx context.Context,
	cancel context.CancelFunc,
	port string,
) {
	c.BaseRestController.Listen(ctx, cancel, port)
}

func NewAnalyticsRestController(
	s *services.AnalyticsService,
	tracer *opentracing.Tracer,
	logger *zap.Logger,
	mc *metrics.MetricsCollector,
) *AnalyticsRestController {
	bc := restcontroller.NewRestController(
		tracer,
		logger,
		mc,
	)
	c := AnalyticsRestController{
		service:            s,
		BaseRestController: bc,
	}

	bc.Router.POST("/analytics", bc.Handler(c.handleLogData))
	bc.Router.GET("/analytics", bc.Handler(c.handleGetData))
	return &c
}
