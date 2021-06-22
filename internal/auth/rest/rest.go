package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/morzhanov/go-realworld/internal/auth/dto"
	. "github.com/morzhanov/go-realworld/internal/auth/services"
)

type AuthRestController struct {
	service *AuthService
	router  *gin.Engine
}

func handleError(c *gin.Context, err error) {
	c.String(http.StatusInternalServerError, err.Error())
}

func (c *AuthRestController) handleAuthValidation(ctx *gin.Context) {
	// TODO: review how to properly get and validate token and proxy response
}

func (c *AuthRestController) handleLogin(ctx *gin.Context) {
	jsonData, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		handleError(ctx, err)
		return
	}

	input := dto.LoginInput{}
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

func (c *AuthRestController) handleSignup(ctx *gin.Context) {
	jsonData, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		handleError(ctx, err)
		return
	}

	input := dto.SignupInput{}
	if err = json.Unmarshal(jsonData, &input); err != nil {
		handleError(ctx, err)
		return
	}

	res, err := c.service.Signup(&input)
	if err != nil {
		handleError(ctx, err)
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
