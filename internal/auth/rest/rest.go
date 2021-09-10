package rest

import (
	"context"
	"github.com/morzhanov/go-realworld/internal/common/rest/restcontroller"
	"net/http"

	"github.com/gin-gonic/gin"
	arpc "github.com/morzhanov/go-realworld/api/rpc/auth"
	"github.com/morzhanov/go-realworld/internal/auth/services"
	"github.com/morzhanov/go-realworld/internal/common/metrics"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type AuthRestController struct {
	*restcontroller.BaseRestController
	service *services.AuthService
}

func (c *AuthRestController) handleAuthValidation(ctx *gin.Context) {
	input := arpc.ValidateRestRequestInput{}
	if err := c.ParseRestBody(ctx, &input); err != nil {
		c.HandleRestError(ctx, err)
		return
	}

	res, err := c.service.ValidateRestRequest(&input)
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *AuthRestController) handleLogin(ctx *gin.Context) {
	input := arpc.LoginInput{}
	if err := c.ParseRestBody(ctx, &input); err != nil {
		c.HandleRestError(ctx, err)
		return
	}

	span := c.GetSpan(ctx)
	reqCtx := context.WithValue(context.Background(), "transport", sender.RestTransport)
	res, err := c.service.Login(reqCtx, &input, span)
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *AuthRestController) handleSignup(ctx *gin.Context) {
	input := arpc.SignupInput{}
	if err := c.ParseRestBody(ctx, &input); err != nil {
		c.HandleRestError(ctx, err)
		return
	}

	span := c.GetSpan(ctx)
	reqCtx := context.WithValue(context.Background(), "transport", sender.RestTransport)
	res, err := c.service.Signup(reqCtx, &input, span)
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *AuthRestController) Listen(
	ctx context.Context,
	cancel context.CancelFunc,
	port string,
) {
	c.BaseRestController.Listen(ctx, cancel, port)
}

func NewAuthRestController(
	s *services.AuthService,
	tracer *opentracing.Tracer,
	logger *zap.Logger,
	mc *metrics.MetricsCollector,
) *AuthRestController {
	bc := restcontroller.NewRestController(
		tracer,
		logger,
		mc,
	)
	c := AuthRestController{
		service:            s,
		BaseRestController: bc,
	}

	bc.Router.GET("/auth", bc.Handler(c.handleAuthValidation))
	bc.Router.POST("/login", bc.Handler(c.handleLogin))
	bc.Router.POST("/signup", bc.Handler(c.handleSignup))
	return &c
}
