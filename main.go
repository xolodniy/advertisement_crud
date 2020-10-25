package main

import (
	"advertisement_crud/application"
	"advertisement_crud/controller"
	"advertisement_crud/etc/config"
	"advertisement_crud/model"
)

func main() {
	conf := config.New("/etc/advertisement_crud/config.json")
	m := model.New(conf.Database)
	if !m.IsMigrated() {
		m.Migrate()
	}

	a := application.New(m)
	c := controller.New(a, conf)

	c.ServeHTTP()
}
