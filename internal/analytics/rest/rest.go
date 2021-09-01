package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	anrpc "github.com/morzhanov/go-realworld/api/rpc/analytics"
	. "github.com/morzhanov/go-realworld/internal/analytics/services"
	"github.com/morzhanov/go-realworld/internal/common/helper"
)

type AnalyticsRestController struct {
	service *AnalyticsService
	router  *gin.Engine
}

func (c *AnalyticsRestController) handleLogData(ctx *gin.Context) {
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

func NewRestController(s *AnalyticsService) *AnalyticsRestController {
	router := gin.Default()
	c := AnalyticsRestController{s, router}

	router.GET("/analytics", c.handleLogData)
	router.GET("/analytics/:offset", c.handleGetData)
	return &c
}
