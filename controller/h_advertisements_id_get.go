package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type GetAdvertisementResponse struct {
	Caption     string   `json:"caption"                example:"Car for sale"`
	Price       int      `json:"price"                  example:"30100"`
	Description string   `json:"description,omitempty"  example:"low price"`
	Photos      []string `json:"photos,omitempty"`
}

// @Summary Объявления
// @Description Получить объявление по идентификатору
// @Description Допустимы опциональные поля в параметре 'fields': 'description', 'photos'
// @Accept  json
// @Produce  json
// @Param id     path  int    true  "Идентификатор объявления" default(3)
// @Param fields query string false "Дополнительные поля"
// @Success 200 {object} GetAdvertisementResponse
// @Router /api/v1/advertisements/id{id} [get]
func (c *Controller) getAdvertisement(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id < 1 {
		respondError(ctx, errors.New("Неправильно задан идентификатор объявления"))
		return
	}
	fields := strings.Split(ctx.Query("fields"), ",")
	for i := len(fields) - 1; i >= 0; i-- {
		switch fields[i] {
		case "description", "photos":
			// do nothing
		case "":
			// omit empty field
			fields = append(fields[:i], fields[i+1:]...)
		default:
			respondError(ctx, fmt.Errorf("Неподдерживаемое fields поле: '%s'", fields[i]))
			return
		}
	}

	advertisement, err := c.app.GetAdvertisement(id)
	if err != nil {
		respondError(ctx, err)
		return
	}

	response := GetAdvertisementResponse{
		Caption: advertisement.Caption,
		Price:   advertisement.Price,
	}
	for _, field := range fields {
		switch field {
		case "description":
			response.Description = advertisement.Description
		case "photos":
			response.Photos = make([]string, len(advertisement.Photos))
			for i, photo := range advertisement.Photos {
				response.Photos[i] = fmt.Sprintf("%s/photos/id%d", c.config.FQDN, photo.ID)
			}
		}
	}
	ctx.JSON(http.StatusOK, response)
}
