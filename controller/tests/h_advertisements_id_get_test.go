package tests

import (
	"net/http"
	"net/http/httptest"

	"advertisement_crud/model"

	"gorm.io/gorm"
)

func (s *UnitTestSuite) TestGetAdvertisementByID() {
	// Arrange
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/advertisements/id12?fields=photos", nil)
	s.app.On("GetAdvertisement", 12).Return(model.Advertisement{
		Model: gorm.Model{
			ID: 12,
		},
		Caption:     "sell car",
		Description: "want to sell car",
		Price:       1200,
		Photos: []model.Photo{
			{Model: gorm.Model{ID: 13}},
			{Model: gorm.Model{ID: 14}},
			{Model: gorm.Model{ID: 15}},
		},
	}, nil)

	// Act
	s.controller.HandleRequest(w, req)

	// Assert
	s.app.AssertExpectations(s.T())

	body := s.ReadAllAndAssertErr(w.Result().Body)

	s.Equal(200, w.Code)
	expected := `
	{
		"caption": "sell car",
		"price": 1200,
		"photos": [
			"localhost:8080/photos/id13",
			"localhost:8080/photos/id14",
			"localhost:8080/photos/id15"
		]
	}`
	s.JSONEq(expected, body)
}
