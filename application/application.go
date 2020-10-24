package application

type Application struct {
	model IModel
}

func New(app IModel) *Application {
	return &Application{
		model: app,
	}
}

type IModel interface{}
