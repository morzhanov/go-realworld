package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	picturedto "github.com/morzhanov/go-realworld/internal/pictures/dto"
	. "github.com/morzhanov/go-realworld/internal/pictures/services"
)

type PicturesRestController struct {
	service *PictureService
	router  *gin.Engine
}

func handleError(c *gin.Context, err error) {
	c.String(http.StatusInternalServerError, err.Error())
}

func (c *PicturesRestController) handleCreateUserPicture(ctx *gin.Context) {
	userId := ctx.Param("userId")

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

	res, err := c.service.CreateUserPicture(userId, &input)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, res)
}

func (c *PicturesRestController) handleGetUserPictures(ctx *gin.Context) {
	userId := ctx.Param("userId")

	res, err := c.service.GetUserPictures(userId)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *PicturesRestController) handleGetUserPicture(ctx *gin.Context) {
	userId := ctx.Param("userId")
	id := ctx.Param("id")

	res, err := c.service.GetUserPicture(userId, id)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *PicturesRestController) handleDeleteUserPicture(ctx *gin.Context) {
	userId := ctx.Param("userId")
	id := ctx.Param("id")

	err := c.service.DeleteUserPicture(userId, id)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func NewRestController(s *PictureService) *PicturesRestController {
	router := gin.Default()
	c := PicturesRestController{s, router}

	router.POST("/pictures", c.handleCreateUserPicture)
	router.GET("/pictures/:userId", c.handleGetUserPictures)
	router.GET("/pictures/:userId/:id", c.handleGetUserPicture)
	router.DELETE("/pictures/:userId/:id", c.handleDeleteUserPicture)
	return &c
}
