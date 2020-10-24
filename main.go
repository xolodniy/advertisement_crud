package main

import (
	"advertisement_crud/etc/config"
	"advertisement_crud/model"

	"github.com/sirupsen/logrus"
)

func main() {
	conf := config.New("/etc/advertisement_crud/config.json")
	_ = model.New(conf.Database)
	logrus.Info("Hello world")
}
