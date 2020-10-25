package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type GetAdvertisementsResponse struct {
	Advertisements []Advertisement `json:"advertisements"`
}

type Advertisement struct {
	ID      int    `json:"id"      example:"144122"`
	Caption string `json:"caption" example:"Car for sale"`
	Price   int    `json:"price"   example:"12300"`
	Photo   string `json:"photo"   example:"my.net/api/v1/14412"`
}

// @Summary Объявления
// @Description Список объявлений
// @Accept  json
// @Produce  json
// @Param order     query string false "Сортировка" Enums(price, created_at)
// @Param direction query string false "Направление сортировки" Enums(asc, desc) default(desc)
// @Success 200 {object} GetAdvertisementsResponse
// @Router /api/v1/advertisements [get]
func (c *Controller) getAdvertisements(ctx *gin.Context) {
	order := strings.ToLower(ctx.DefaultQuery("order", "id"))
	direction := strings.ToLower(ctx.DefaultQuery("direction", "desc"))
	switch order {
	case "id", "price", "created_at":
		// do nothing
	default:
		respondError(ctx, errors.New("Допустимые значения сортировки: 'price', 'created_at"))
		return
	}
	switch direction {
	case "asc", "desc":
	// do nothing
	default:
		respondError(ctx, errors.New("Допустимые значения направления сортировки: 'asc', 'desc'"))
		return
	}

	advertisements, err := c.app.GetAdvertisements(order + " " + direction)
	if err != nil {
		respondError(ctx, err)
		return
	}

	response := GetAdvertisementsResponse{
		Advertisements: make([]Advertisement, len(advertisements)),
	}
	for i := range advertisements {
		response.Advertisements[i] = Advertisement{
			ID:      int(advertisements[i].ID),
			Caption: advertisements[i].Caption,
			Price:   advertisements[i].Price,
		}
		if len(advertisements[i].Photos) > 0 {
			response.Advertisements[i].Photo = fmt.Sprintf("%s/photos/id%d", c.config.FQDN, advertisements[i].Photos[0].ID)
		}
	}
	ctx.JSON(http.StatusOK, response)
}
