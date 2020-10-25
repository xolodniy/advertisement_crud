package controller

import (
	"net/http"

	"advertisement_crud/model"

	"github.com/gin-gonic/gin"
)

type CreateAdvertisementReq struct {
	Caption     string `json:"caption"     binding:"required"`
	Description string `json:"description" binding:"required"`
	Price       int    `json:"price"       binding:"min=1"`
	PhotoIDs    []int  `json:"photoIDs" binding:"min=1,dive,min=1"`
}

type CreateAdvertisementResponse struct {
	ID int `json:"id"`
}

func (c *Controller) createAdvertisement(ctx *gin.Context) {
	var req CreateAdvertisementReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	id, err := c.app.CreateAdvertisement(model.Advertisement{
		Caption:     req.Caption,
		Description: req.Description,
		Price:       req.Price,
	}, req.PhotoIDs)
	if err != nil {
		respondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, CreateAdvertisementResponse{
		ID: id,
	})
}
