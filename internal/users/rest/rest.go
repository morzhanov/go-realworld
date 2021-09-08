package rest

import (
	"context"
	"github.com/morzhanov/go-realworld/internal/common/rest/restcontroller"
	"net/http"

	"github.com/gin-gonic/gin"
	urpc "github.com/morzhanov/go-realworld/api/rpc/users"
	"github.com/morzhanov/go-realworld/internal/common/metrics"
	"github.com/morzhanov/go-realworld/internal/users/services"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type UsersRestController struct {
	service        *services.UsersService
	baseController *restcontroller.BaseRestController
}

func (c *UsersRestController) handleGetUserData(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := c.service.GetUserData(id)
	if err != nil {
		c.baseController.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *UsersRestController) handleGetUserDataByUsername(ctx *gin.Context) {
	username := ctx.Query("username")
	if username == "" {
		ctx.String(http.StatusBadRequest, "username should be provided in query params")
		return
	}

	res, err := c.service.GetUserDataByUsername(username)
	if err != nil {
		c.baseController.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *UsersRestController) handleValidateUserPassword(ctx *gin.Context) {
	input := urpc.ValidateUserPasswordRequest{}
	if err := c.baseController.ParseRestBody(ctx, &input); err != nil {
		c.baseController.HandleRestError(ctx, err)
		return
	}

	if err := c.service.ValidateUserPassword(&input); err != nil {
		c.baseController.HandleRestError(ctx, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (c *UsersRestController) handleCreateUser(ctx *gin.Context) {
	input := urpc.CreateUserRequest{}
	if err := c.baseController.ParseRestBody(ctx, &input); err != nil {
		c.baseController.HandleRestError(ctx, err)
		return
	}

	res, err := c.service.CreateUser(&input)
	if err != nil {
		c.baseController.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *UsersRestController) handleDeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.DeleteUser(id); err != nil {
		c.baseController.HandleRestError(ctx, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (c *UsersRestController) Listen(
	ctx context.Context,
	port string,
) error {
	return c.baseController.Listen(ctx, port)
}

func NewUsersRestController(
	s *services.UsersService,
	tracer *opentracing.Tracer,
	logger *zap.Logger,
	mc *metrics.MetricsCollector,
) *UsersRestController {
	bc := restcontroller.NewRestController(
		tracer,
		logger,
		mc,
	)
	c := UsersRestController{
		service:        s,
		baseController: bc,
	}

	bc.Router.GET("/users/:id", bc.Handler(c.handleGetUserData))
	bc.Router.GET("/users", bc.Handler(c.handleGetUserDataByUsername))
	bc.Router.POST("/users/validate-password", bc.Handler(c.handleValidateUserPassword))
	bc.Router.POST("/users", bc.Handler(c.handleCreateUser))
	bc.Router.DELETE("/users/:id", bc.Handler(c.handleDeleteUser))
	return &c
}
