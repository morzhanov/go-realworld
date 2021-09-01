package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	arpc "github.com/morzhanov/go-realworld/api/rpc/auth"
	. "github.com/morzhanov/go-realworld/internal/auth/services"
	"github.com/morzhanov/go-realworld/internal/common/helper"
	"github.com/morzhanov/go-realworld/internal/common/sender"
)

type AuthRestController struct {
	service *AuthService
	router  *gin.Engine
}

func (c *AuthRestController) handleAuthValidation(ctx *gin.Context) {
	input := arpc.ValidateRestRequestInput{}
	if err := helper.ParseRestBody(ctx, &input); err != nil {
		helper.HandleRestError(ctx, err)
		return
	}

	res, err := c.service.ValidateRestRequest(&input)
	if err != nil {
		helper.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *AuthRestController) handleLogin(ctx *gin.Context) {
	input := arpc.LoginInput{}
	if err := helper.ParseRestBody(ctx, &input); err != nil {
		helper.HandleRestError(ctx, err)
		return
	}

	reqCtx := context.WithValue(context.Background(), "transport", sender.RestTransport)
	res, err := c.service.Login(reqCtx, &input)
	if err != nil {
		helper.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *AuthRestController) handleSignup(ctx *gin.Context) {
	input := arpc.SignupInput{}
	if err := helper.ParseRestBody(ctx, &input); err != nil {
		helper.HandleRestError(ctx, err)
		return
	}

	reqCtx := context.WithValue(context.Background(), "transport", sender.RestTransport)
	res, err := c.service.Signup(reqCtx, &input)
	if err != nil {
		helper.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func NewRestController(s *AuthService) *AuthRestController {
	router := gin.Default()
	c := AuthRestController{s, router}

	router.GET("/auth", c.handleAuthValidation)
	router.POST("/login", c.handleLogin)
	router.POST("/signup", c.handleSignup)
	return &c
}
