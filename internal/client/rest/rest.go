package rest

import (
	"github.com/gin-gonic/gin"
	. "github.com/morzhanov/go-realworld/internal/client/services"
)

type ClientRestController struct {
	service *ClientService
	router  *gin.Engine
}

// TODO: kuber gateway/ingress should inject userId after request is authenticated
// TODO: or we can call auth.ValidateXXXXRequest explicitly with each request
func (c *ClientRestController) handleLogin(ctx *gin.Context) {}

func (c *ClientRestController) handleSignup(ctx *gin.Context) {}

func (c *ClientRestController) handleCreatePicture(ctx *gin.Context) {}

func (c *ClientRestController) handleGetPictures(ctx *gin.Context) {}

func (c *ClientRestController) handleGetPicture(ctx *gin.Context) {}

func (c *ClientRestController) handleDeletePicture(ctx *gin.Context) {}

func (c *ClientRestController) handleGetAnalytics(ctx *gin.Context) {}

func NewRestController(s *ClientService) *ClientRestController {
	router := gin.Default()

	c := ClientRestController{s, router}

	router.POST("/login", c.handleLogin)
	router.POST("/signup", c.handleSignup)
	router.POST("/pictures", c.handleCreatePicture)
	router.GET("/pictures", c.handleGetPictures)
	router.GET("/pictures/:id", c.handleGetPicture)
	router.DELETE("/pictures/:id", c.handleDeletePicture)
	router.GET("/analytics", c.handleGetAnalytics)

	return &c
}
