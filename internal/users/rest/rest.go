package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	urpc "github.com/morzhanov/go-realworld/api/rpc/users"
	"github.com/morzhanov/go-realworld/internal/common/helper"
	. "github.com/morzhanov/go-realworld/internal/users/services"
)

type UsersRestController struct {
	service *UsersService
	router  *gin.Engine
}

func (c *UsersRestController) handleGetUserData(ctx *gin.Context) {
	id := ctx.Param("id")

	res, err := c.service.GetUserData(id)
	if err != nil {
		helper.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *UsersRestController) handleGetUserDataByUsername(ctx *gin.Context) {
	username := ctx.Query("username")

	if username == "" {
		ctx.String(http.StatusBadRequest, "username should be provided in query params")
	}

	res, err := c.service.GetUserDataByUsername(username)
	if err != nil {
		helper.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *UsersRestController) handleValidateUserPassword(ctx *gin.Context) {
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
	id := ctx.Param("id")
	if err := c.service.DeleteUser(id); err != nil {
		helper.HandleRestError(ctx, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func NewRestController(s *UsersService) *UsersRestController {
	router := gin.Default()
	c := UsersRestController{s, router}

	router.GET("/users/:id", c.handleGetUserData)
	router.GET("/users", c.handleGetUserDataByUsername)
	router.POST("/users/validate-password", c.handleValidateUserPassword)
	router.POST("/users", c.handleCreateUser)
	router.POST("/users/:id", c.handleDeleteUser)
	return &c
}
