package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	urpc "github.com/morzhanov/go-realworld/api/rpc/users"
	"github.com/morzhanov/go-realworld/internal/common/helper"
	"github.com/morzhanov/go-realworld/internal/common/metrics"
	"github.com/morzhanov/go-realworld/internal/common/tracing"
	"github.com/morzhanov/go-realworld/internal/users/services"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type UsersRestController struct {
	service *services.UsersService
	router  *gin.Engine
	tracer  *opentracing.Tracer
}

func (c *UsersRestController) handleGetUserData(ctx *gin.Context) {
	span := tracing.StartSpanFromHttpRequest(*c.tracer, ctx.Request)
	defer span.Finish()

	id := ctx.Param("id")

	res, err := c.service.GetUserData(id)
	if err != nil {
		helper.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *UsersRestController) handleGetUserDataByUsername(ctx *gin.Context) {
	span := tracing.StartSpanFromHttpRequest(*c.tracer, ctx.Request)
	defer span.Finish()

	username := ctx.Query("username")

	if username == "" {
		ctx.String(http.StatusBadRequest, "username should be provided in query params")
		return
	}

	res, err := c.service.GetUserDataByUsername(username)
	if err != nil {
		helper.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *UsersRestController) handleValidateUserPassword(ctx *gin.Context) {
	span := tracing.StartSpanFromHttpRequest(*c.tracer, ctx.Request)
	defer span.Finish()

	input := urpc.ValidateUserPasswordRequest{}
	if err := helper.ParseRestBody(ctx, &input); err != nil {
		helper.HandleRestError(ctx, err)
		return
	}

	if err := c.service.ValidateUserPassword(&input); err != nil {
		helper.HandleRestError(ctx, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (c *UsersRestController) handleCreateUser(ctx *gin.Context) {
	span := tracing.StartSpanFromHttpRequest(*c.tracer, ctx.Request)
	defer span.Finish()

	input := urpc.CreateUserRequest{}
	if err := helper.ParseRestBody(ctx, &input); err != nil {
		helper.HandleRestError(ctx, err)
		return
	}

	res, err := c.service.CreateUser(&input)
	if err != nil {
		helper.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *UsersRestController) handleDeleteUser(ctx *gin.Context) {
	span := tracing.StartSpanFromHttpRequest(*c.tracer, ctx.Request)
	defer span.Finish()

	id := ctx.Param("id")
	if err := c.service.DeleteUser(id); err != nil {
		helper.HandleRestError(ctx, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (c *UsersRestController) Listen(ctx context.Context, port string, logger *zap.Logger) {
	helper.StartRestServer(ctx, port, c.router, logger)
}

func NewUsersRestController(s *services.UsersService, tracer *opentracing.Tracer, mc *metrics.MetricsCollector) *UsersRestController {
	router := gin.Default()
	c := UsersRestController{s, router, tracer}

	router.GET("/users/:id", c.handleGetUserData)
	router.GET("/users", c.handleGetUserDataByUsername)
	router.POST("/users/validate-password", c.handleValidateUserPassword)
	router.POST("/users", c.handleCreateUser)
	router.DELETE("/users/:id", c.handleDeleteUser)
	mc.RegisterMetricsEndpoint(router)
	return &c
}
