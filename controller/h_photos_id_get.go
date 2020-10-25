package controller

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (c *Controller) getPhoto(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id < 1 {
		respondError(ctx, errors.New("Неправильно задан идентификатор фотографии"))
		return
	}
	photo, err := c.app.GetPhoto(id)
	if err != nil {
		respondError(ctx, err)
		return
	}
	ctx.Header("Content-Type", photo.Mime)
	if _, err := ctx.Writer.Write(photo.Body); err != nil {
		logrus.WithError(err).Error("can't write photo to response body")
	}
}
