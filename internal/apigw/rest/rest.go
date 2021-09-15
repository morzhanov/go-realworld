package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	anrpc "github.com/morzhanov/go-realworld/api/rpc/analytics"
	arpc "github.com/morzhanov/go-realworld/api/rpc/auth"
	prpc "github.com/morzhanov/go-realworld/api/rpc/pictures"
	"github.com/morzhanov/go-realworld/internal/apigw/services"
	"github.com/morzhanov/go-realworld/internal/common/metrics"
	"github.com/morzhanov/go-realworld/internal/common/rest/restcontroller"
	"github.com/morzhanov/go-realworld/internal/common/sender"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"net/http"
)

type APIGatewayRestController struct {
	*restcontroller.BaseRestController
	service *services.APIGatewayService
	sender  *sender.Sender
}

func (c *APIGatewayRestController) handleLogin(ctx *gin.Context) {
	transport, err := c.sender.StringToTransport(ctx.Param("transport"))
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}
	input := arpc.LoginInput{}
	if err := c.ParseRestBody(ctx, &input); err != nil {
		c.HandleRestError(ctx, err)
		return
	}

	span := c.GetSpan(ctx)
	res, err := c.service.Login(transport, &input, span)
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *APIGatewayRestController) handleSignup(ctx *gin.Context) {
	transport, err := c.sender.StringToTransport(ctx.Param("transport"))
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}
	input := arpc.SignupInput{}
	if err := c.ParseRestBody(ctx, &input); err != nil {
		c.HandleRestError(ctx, err)
		return
	}

	span := c.GetSpan(ctx)
	res, err := c.service.Signup(transport, &input, span)
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *APIGatewayRestController) handleCreatePicture(ctx *gin.Context) {
	transport, err := c.sender.StringToTransport(ctx.Param("transport"))
	span := c.GetSpan(ctx)
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}
	validationRes, err := c.service.CheckAuth(ctx, transport, "pictures", "createUserPicture", span)
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}

	input := prpc.CreateUserPictureRequest{}
	if err := c.ParseRestBody(ctx, &input); err != nil {
		c.HandleRestError(ctx, err)
		return
	}
	input.UserId = validationRes.UserId

	res, err := c.service.CreatePicture(transport, &input, span)
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *APIGatewayRestController) handleGetPictures(ctx *gin.Context) {
	transport, err := c.sender.StringToTransport(ctx.Param("transport"))
	span := c.GetSpan(ctx)
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}
	validationRes, err := c.service.CheckAuth(ctx, transport, "pictures", "getUserPictures", span)
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}

	res, err := c.service.GetPictures(transport, validationRes.UserId, span)
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res.Pictures)
}

func (c *APIGatewayRestController) handleGetPicture(ctx *gin.Context) {
	transport, err := c.sender.StringToTransport(ctx.Param("transport"))
	picId := ctx.Param("id")
	span := c.GetSpan(ctx)
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}
	validationRes, err := c.service.CheckAuth(ctx, transport, "pictures", "getUserPicture", span)
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}

	res, err := c.service.GetPicture(transport, validationRes.UserId, picId, span)
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *APIGatewayRestController) handleDeletePicture(ctx *gin.Context) {
	transport, err := c.sender.StringToTransport(ctx.Param("transport"))
	span := c.GetSpan(ctx)
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}
	validationRes, err := c.service.CheckAuth(ctx, transport, "pictures", "deletePicture", span)
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}

	input := prpc.DeleteUserPictureRequest{}
	if err := c.ParseRestBody(ctx, &input); err != nil {
		c.HandleRestError(ctx, err)
		return
	}

	err = c.service.DeletePicture(transport, validationRes.UserId, input.PictureId, span)
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (c *APIGatewayRestController) handleGetAnalytics(ctx *gin.Context) {
	transport, err := c.sender.StringToTransport(ctx.Param("transport"))
	span := c.GetSpan(ctx)
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}
	_, err = c.service.CheckAuth(ctx, transport, "analytics", "getLogs", span)
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}

	input := anrpc.GetLogRequest{}
	if err := c.ParseRestBody(ctx, &input); err != nil {
		c.HandleRestError(ctx, err)
		return
	}

	res, err := c.service.GetAnalytics(transport, &input, span)
	if err != nil {
		c.HandleRestError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *APIGatewayRestController) Listen(
	ctx context.Context,
	cancel context.CancelFunc,
	port string,
) {
	c.BaseRestController.Listen(ctx, cancel, port)
}

func NewAPIGatewayRestController(
	s *services.APIGatewayService,
	tracer *opentracing.Tracer,
	logger *zap.Logger,
	mc *metrics.MetricsCollector,
	sender *sender.Sender,
) *APIGatewayRestController {
	bc := restcontroller.NewRestController(
		tracer,
		logger,
		mc,
	)
	c := APIGatewayRestController{
		service:            s,
		sender:             sender,
		BaseRestController: bc,
	}

	bc.Router.POST("/:transport/login", bc.Handler(c.handleLogin))
	bc.Router.POST("/:transport/signup", bc.Handler(c.handleSignup))
	bc.Router.POST("/:transport/pictures", bc.Handler(c.handleCreatePicture))
	bc.Router.GET("/:transport/pictures", bc.Handler(c.handleGetPictures))
	bc.Router.GET("/:transport/pictures/:id", bc.Handler(c.handleGetPicture))
	bc.Router.DELETE("/:transport/pictures/:id", bc.Handler(c.handleDeletePicture))
	bc.Router.GET("/:transport/analytics", bc.Handler(c.handleGetAnalytics))
	return &c
}
