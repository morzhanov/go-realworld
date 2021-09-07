package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	arpc "github.com/morzhanov/go-realworld/api/rpc/auth"
	"github.com/morzhanov/go-realworld/internal/auth/services"
	"github.com/morzhanov/go-realworld/internal/common/helper"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/morzhanov/go-realworld/internal/common/tracing"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type AuthRestController struct {
	service *services.AuthService
	router  *gin.Engine
	tracer  *opentracing.Tracer
}

func (c *AuthRestController) handleAuthValidation(ctx *gin.Context) {
	span := tracing.StartSpanFromHttpRequest(*c.tracer, ctx.Request)
	defer span.Finish()

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
	span := tracing.StartSpanFromHttpRequest(*c.tracer, ctx.Request)
	defer span.Finish()

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
	span := tracing.StartSpanFromHttpRequest(*c.tracer, ctx.Request)
	defer span.Finish()

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

func (c *AuthRestController) Listen(ctx context.Context, port string, logger *zap.Logger) {
	helper.StartRestServer(ctx, port, c.router, logger)
}

func NewAuthRestController(s *services.AuthService, tracer *opentracing.Tracer) *AuthRestController {
	router := gin.Default()
	c := AuthRestController{s, router, tracer}

	router.GET("/auth", c.handleAuthValidation)
	router.POST("/login", c.handleLogin)
	router.POST("/signup", c.handleSignup)
	return &c
}
