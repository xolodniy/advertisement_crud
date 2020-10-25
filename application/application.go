package application

import (
	"fmt"

	"advertisement_crud/model"
)

type Application struct {
	model IModel
}

func New(model IModel) *Application {
	return &Application{
		model: model,
	}
}

type IModel interface {
	GetAdvertisement(id int) (model.Advertisement, error)
	GetAdvertisements(order string) ([]model.Advertisement, error)
	CreateAdvertisement(advertisement model.Advertisement, photoIDs []int) (int, error)
	CheckPhotoExist(id int) bool
	CreatePhoto(photo model.Photo) (int, error)
	GetPhoto(id int) (model.Photo, error)
}

func (a *Application) CreateAdvertisement(advertisement model.Advertisement, photoIDs []int) (int, error) {
	for _, id := range photoIDs {
		if !a.model.CheckPhotoExist(id) {
			return 0, fmt.Errorf("Не найдена фотография с идентификатором '%d'", id)
		}
	}
	return a.model.CreateAdvertisement(advertisement, photoIDs)
}

func (a *Application) GetAdvertisements(order string) ([]model.Advertisement, error) {
	return a.model.GetAdvertisements(order)
}
func (a *Application) GetAdvertisement(id int) (model.Advertisement, error) {
	return a.model.GetAdvertisement(id)
}

func (a *Application) CreatePhoto(photo model.Photo) (int, error) {
	return a.model.CreatePhoto(photo)
}

func (a *Application) GetPhoto(id int) (model.Photo, error) {
	return a.model.GetPhoto(id)
}
