package controller

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"advertisement_crud/model"

	"github.com/gin-gonic/gin"
)

type PostPhotoResponse struct {
	ID int `json:"id"`
}

func (c *Controller) postPhoto(ctx *gin.Context) {
	errWrongFile := errors.New("ошибка передачи файла")

	form, err := ctx.FormFile("photo")
	if err != nil {
		respondError(ctx, errWrongFile)
		return
	}
	var maxFileSizeMB = 10
	if form.Size > int64(maxFileSizeMB<<20) {
		respondError(ctx, fmt.Errorf("Максимальный размер файла не должен превышать %d мегабайт", maxFileSizeMB))
		return
	}
	file, err := form.Open()
	if err != nil {
		respondError(ctx, errWrongFile)
		return
	}
	defer file.Close()
	docBody, err := ioutil.ReadAll(file)
	if err != nil {
		respondError(ctx, errWrongFile)
		return
	}
	mime := http.DetectContentType(docBody)
	switch mime {
	case "image/png", "image/jpeg", "image/bmp":
	// do nothing
	default:
		respondError(ctx, errors.New("Поддерживаемые форматы файлов: 'png', 'jpeg', 'bmp'"))
		return
	}

	id, err := c.app.CreatePhoto(model.Photo{
		Mime: mime,
		Body: docBody,
	})
	if err != nil {
		respondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, PostPhotoResponse{ID: id})
}
