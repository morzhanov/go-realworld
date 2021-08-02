package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	analyticsdto "github.com/morzhanov/go-realworld/internal/analytics/dto"
	. "github.com/morzhanov/go-realworld/internal/apigw/services"
	authdto "github.com/morzhanov/go-realworld/internal/auth/dto"
	picturedto "github.com/morzhanov/go-realworld/internal/pictures/dto"
)

type APIGatewayRestController struct {
	service *APIGatewayService
	router  *gin.Engine
}

func handleError(c *gin.Context, err error) {
	c.String(http.StatusInternalServerError, err.Error())
}

// TODO: kuber gateway/ingress should inject userId after request is authenticated
// TODO: or we can call auth.ValidateXXXXRequest explicitly with each request
func (c *APIGatewayRestController) handleLogin(ctx *gin.Context) {
	transport, err := strconv.Atoi(ctx.Param("transport"))
	if err != nil {
		handleError(ctx, err)
		return
	}
	jsonData, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		handleError(ctx, err)
		return
	}

	input := authdto.LoginInput{}
	if err = json.Unmarshal(jsonData, &input); err != nil {
		handleError(ctx, err)
		return
	}

	res, err := c.service.Login(Transport(transport), &input)
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
	jsonData, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		handleError(ctx, err)
		return
	}

	input := authdto.SignupInput{}
	if err = json.Unmarshal(jsonData, &input); err != nil {
		handleError(ctx, err)
		return
	}

	res, err := c.service.Signup(Transport(transport), &input)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *APIGatewayRestController) handleCreatePicture(ctx *gin.Context) {
	transport, err := strconv.Atoi(ctx.Param("transport"))
	if err != nil {
		handleError(ctx, err)
		return
	}
	jsonData, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		handleError(ctx, err)
		return
	}

	input := picturedto.CreatePicturesDto{}
	if err = json.Unmarshal(jsonData, &input); err != nil {
		handleError(ctx, err)
		return
	}

	// TODO: get userId from kuber ingress or via auth service varify token endpoint
	userId := "user-id"

	res, err := c.service.CreatePicture(Transport(transport), userId, &input)
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
	// TODO: get userId from kuber ingress or via auth service varify token endpoint
	userId := "user-id"

	res, err := c.service.GetPictures(Transport(transport), userId)
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
	id := ctx.Param("id")
	// TODO: get userId from kuber ingress or via auth service varify token endpoint
	userId := "user-id"

	res, err := c.service.GetPicture(Transport(transport), userId, id)
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
	id := ctx.Param("id")
	// TODO: get userId from kuber ingress or via auth service varify token endpoint
	userId := "user-id"

	err = c.service.DeletePicture(Transport(transport), userId, id)
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
	jsonData, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		handleError(ctx, err)
		return
	}

	input := analyticsdto.GetLogsInput{}
	if err = json.Unmarshal(jsonData, &input); err != nil {
		handleError(ctx, err)
		return
	}

	res, err := c.service.GetAnalytics(Transport(transport), &input)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, res)
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
