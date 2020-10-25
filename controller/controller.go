package controller

import (
	"fmt"
	"net/http"

	"advertisement_crud/etc/common"
	"advertisement_crud/etc/config"
	"advertisement_crud/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	app    IApplication
	router *gin.Engine
	config config.Main
}

func New(app IApplication, config config.Main) *Controller {
	return &Controller{
		app:    app,
		router: gin.New(),
		config: config,
	}
}

type IApplication interface {
	CreateAdvertisement(advertisement model.Advertisement, photoIDs []int) (int, error)
	GetAdvertisements(order string) ([]model.Advertisement, error)
	GetAdvertisement(id int) (model.Advertisement, error)
	CreatePhoto(photo model.Photo) (int, error)
	GetPhoto(id int) (model.Photo, error)
}

func respondError(ctx *gin.Context, err error) {
	switch err {
	case common.ErrInternal:
		ctx.JSON(http.StatusInternalServerError, err.Error())
	case common.ErrNotFound:
		ctx.JSON(http.StatusNotFound, err.Error())
	default:
		ctx.JSON(http.StatusBadRequest, err.Error())
	}
}

func (c *Controller) ServeHTTP() {
	c.router.POST("/api/v1/advertisements", c.createAdvertisement)
	c.router.GET("/api/v1/advertisements", c.getAdvertisements)
	c.router.GET("/api/v1/advertisements/id:id", c.getAdvertisement)
	c.router.POST("/api/v1/photos", c.postPhoto)
	c.router.GET("/api/v1/photos/id:id", c.getPhoto)

	srv := &http.Server{
		Addr:    fmt.Sprint(":", c.config.Port),
		Handler: c.router,
	}

	logrus.Infof("server run on :%d", c.config.Port)
	srv.ListenAndServe()
}
