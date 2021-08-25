package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	urpc "github.com/morzhanov/go-realworld/api/rpc/users"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	. "github.com/morzhanov/go-realworld/internal/users/services"
)

type UsersRestController struct {
	service *UsersService
	router  *gin.Engine
}

func handleError(c *gin.Context, err error) {
	c.String(http.StatusInternalServerError, err.Error())
}

func (c *UsersRestController) handleGetUserData(ctx *gin.Context) {
	id := ctx.Param("id")

	res, err := c.service.GetUserData(id)
	if err != nil {
		handleError(ctx, err)
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
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *UsersRestController) handleValidateUserPassword(ctx *gin.Context) {
	input := urpc.ValidateUserPasswordRequest{}
	if err := sender.ParseRestBody(ctx, &input); err != nil {
		handleError(ctx, err)
		return
	}

	err := c.service.ValidateUserPassword(&input)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (c *UsersRestController) handleCreateUser(ctx *gin.Context) {
	input := urpc.CreateUserRequest{}
	if err := sender.ParseRestBody(ctx, &input); err != nil {
		handleError(ctx, err)
		return
	}

	res, err := c.service.CreateUser(&input)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *UsersRestController) handleDeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.service.DeleteUser(id)
	if err != nil {
		handleError(ctx, err)
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
