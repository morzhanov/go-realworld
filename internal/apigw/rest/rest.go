package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	anrpc "github.com/morzhanov/go-realworld/api/rpc/analytics"
	arpc "github.com/morzhanov/go-realworld/api/rpc/auth"
	prpc "github.com/morzhanov/go-realworld/api/rpc/pictures"
	. "github.com/morzhanov/go-realworld/internal/apigw/services"
	"github.com/morzhanov/go-realworld/internal/common/sender"
)

type APIGatewayRestController struct {
	service *APIGatewayService
	router  *gin.Engine
}

func handleError(c *gin.Context, err error) {
	c.String(http.StatusInternalServerError, err.Error())
}

func (c *APIGatewayRestController) handleLogin(ctx *gin.Context) {
	transport, err := strconv.Atoi(ctx.Param("transport"))
	if err != nil {
		handleError(ctx, err)
		return
	}
	input := arpc.LoginInput{}
	if err := sender.ParseRestBody(ctx, &input); err != nil {
		handleError(ctx, err)
		return
	}

	res, err := c.service.Login(sender.Transport(transport), &input)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *APIGatewayRestController) handleSignup(ctx *gin.Context) {
	transport, err := strconv.Atoi(ctx.Param("transport"))
	if err != nil {
		handleError(ctx, err)
		return
	}
	input := arpc.SignupInput{}
	if err := sender.ParseRestBody(ctx, &input); err != nil {
		handleError(ctx, err)
		return
	}

	res, err := c.service.Signup(sender.Transport(transport), &input)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *APIGatewayRestController) handleCreatePicture(ctx *gin.Context) {
	transport, err := strconv.Atoi(ctx.Param("transport"))
	if err != nil {
		handleError(ctx, err)
		return
	}
	validationRes, err := c.service.CheckAuth(ctx, sender.Transport(transport), "pictures", "createUserPicture")

	input := prpc.CreateUserPictureRequest{}
	if err := sender.ParseRestBody(ctx, &input); err != nil {
		handleError(ctx, err)
		return
	}
	input.UserId = validationRes.UserId

	res, err := c.service.CreatePicture(sender.Transport(transport), &input)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *APIGatewayRestController) handleGetPictures(ctx *gin.Context) {
	transport, err := strconv.Atoi(ctx.Param("transport"))
	if err != nil {
		handleError(ctx, err)
		return
	}
	validationRes, err := c.service.CheckAuth(ctx, sender.Transport(transport), "pictures", "getUserPictures")

	input := prpc.GetUserPicturesRequest{}
	if err := sender.ParseRestBody(ctx, &input); err != nil {
		handleError(ctx, err)
		return
	}

	res, err := c.service.GetPictures(sender.Transport(transport), validationRes.UserId)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *APIGatewayRestController) handleGetPicture(ctx *gin.Context) {
	transport, err := strconv.Atoi(ctx.Param("transport"))
	if err != nil {
		handleError(ctx, err)
		return
	}
	validationRes, err := c.service.CheckAuth(ctx, sender.Transport(transport), "pictures", "getUserPicture")

	input := prpc.GetUserPictureRequest{}
	if err := sender.ParseRestBody(ctx, &input); err != nil {
		handleError(ctx, err)
		return
	}

	res, err := c.service.GetPicture(sender.Transport(transport), validationRes.UserId, input.PictureId)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *APIGatewayRestController) handleDeletePicture(ctx *gin.Context) {
	transport, err := strconv.Atoi(ctx.Param("transport"))
	if err != nil {
		handleError(ctx, err)
		return
	}
	validationRes, err := c.service.CheckAuth(ctx, sender.Transport(transport), "pictures", "deletePicture")

	input := prpc.DeleteUserPictureRequest{}
	if err := sender.ParseRestBody(ctx, &input); err != nil {
		handleError(ctx, err)
		return
	}

	err = c.service.DeletePicture(sender.Transport(transport), validationRes.UserId, input.PictureId)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (c *APIGatewayRestController) handleGetAnalytics(ctx *gin.Context) {
	transport, err := strconv.Atoi(ctx.Param("transport"))
	if err != nil {
		handleError(ctx, err)
		return
	}
	_, err = c.service.CheckAuth(ctx, sender.Transport(transport), "analytics", "getLogs")

	input := anrpc.GetLogRequest{}
	if err := sender.ParseRestBody(ctx, &input); err != nil {
		handleError(ctx, err)
		return
	}

	res, err := c.service.GetAnalytics(sender.Transport(transport), &input)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func NewAPIGatewayRestController(s *APIGatewayService) *APIGatewayRestController {
	router := gin.Default()
	c := APIGatewayRestController{s, router}

	router.POST("/:transport/login", c.handleLogin)
	router.POST("/:transport/signup", c.handleSignup)
	router.POST("/:transport/pictures", c.handleCreatePicture)
	router.GET("/:transport/pictures", c.handleGetPictures)
	router.GET("/:transport/pictures/:id", c.handleGetPicture)
	router.DELETE("/:transport/pictures/:id", c.handleDeletePicture)
	router.GET("/:transport/analytics", c.handleGetAnalytics)
	return &c
}
