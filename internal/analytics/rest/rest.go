package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/morzhanov/go-realworld/internal/analytics/dto"
	. "github.com/morzhanov/go-realworld/internal/analytics/models"
	. "github.com/morzhanov/go-realworld/internal/analytics/services"
)

type AnalyticsRestController struct {
	service *AnalyticsService
	router  *gin.Engine
}

func handleError(c *gin.Context, err error) {
	c.String(http.StatusInternalServerError, err.Error())
}

func (c *AnalyticsRestController) handleLogData(ctx *gin.Context) {
	jsonData, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		handleError(ctx, err)
		return
	}

	input := AnalyticsEntry{}
	if err = json.Unmarshal(jsonData, &input); err != nil {
		handleError(ctx, err)
		return
	}

	if err = c.service.LogData(&input); err != nil {
		handleError(ctx, err)
		return
	}
	ctx.Status(http.StatusCreated)
}

func (c *AnalyticsRestController) handleGetData(ctx *gin.Context) {
	offset, err := strconv.Atoi(ctx.Param("offse"))
	if err != nil {
		handleError(ctx, err)
	}

	res, err := c.service.GetLog(&dto.GetLogsInput{Offset: offset})
	if err != nil {
		handleError(ctx, err)
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
