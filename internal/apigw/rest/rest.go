package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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

	res, err := c.service.Login(&input)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *APIGatewayRestController) handleSignup(ctx *gin.Context) {
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

	res, err := c.service.Signup(&input)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *APIGatewayRestController) handleCreatePicture(ctx *gin.Context) {
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

	res, err := c.service.CreatePicture(userId, &input)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *APIGatewayRestController) handleGetPictures(ctx *gin.Context) {
	// TODO: get userId from kuber ingress or via auth service varify token endpoint
	userId := "user-id"

	res, err := c.service.GetPictures(userId)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *APIGatewayRestController) handleGetPicture(ctx *gin.Context) {
	id := ctx.Param("id")
	// TODO: get userId from kuber ingress or via auth service varify token endpoint
	userId := "user-id"

	res, err := c.service.GetPicture(userId, id)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *APIGatewayRestController) handleDeletePicture(ctx *gin.Context) {
	id := ctx.Param("id")
	// TODO: get userId from kuber ingress or via auth service varify token endpoint
	userId := "user-id"

	err := c.service.DeletePicture(userId, id)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (c *APIGatewayRestController) handleGetAnalytics(ctx *gin.Context) {
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

	res, err := c.service.GetAnalytics(&input)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func NewAPIGatewayRestController(s *APIGatewayService) *APIGatewayRestController {
	router := gin.Default()
	c := APIGatewayRestController{s, router}

	router.POST("/login", c.handleLogin)
	router.POST("/signup", c.handleSignup)
	router.POST("/pictures", c.handleCreatePicture)
	router.GET("/pictures", c.handleGetPictures)
	router.GET("/pictures/:id", c.handleGetPicture)
	router.DELETE("/pictures/:id", c.handleDeletePicture)
	router.GET("/analytics", c.handleGetAnalytics)
	return &c
}
