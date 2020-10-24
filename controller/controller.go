package controller

type Controller struct {
	app IApplication
}

func New(app IApplication) *Controller {
	return &Controller{
		app: app,
	}
}

type IApplication interface{}
