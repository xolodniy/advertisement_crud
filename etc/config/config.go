package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
)

type Main struct {
	Port     int    `json:"port"     binding:"min=1,max=65535"`
	LogLevel string `json:"logLevel" binding:"required"`
	FQDN     string `json:"fqdn"     binding:"required"`

	Database Database
}

// Database configuration
type Database struct {
	Host     string `json:"host"     binding:"required"`
	Port     int    `json:"port"     binding:"min=1,max=65535"`
	User     string `json:"user"     binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name"     binding:"required"`
}

func New(path string) Main {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.WithError(err).WithField("configPath", path).Fatal("can't read config file in the selected path")
	}
	var config Main
	if err := json.Unmarshal(body, &config); err != nil {
		logrus.WithError(err).Fatal("can't unmarshal config file as a json object")
	}

	if err := binding.Validator.ValidateStruct(config); err != nil {
		logrus.WithError(err).Fatal("config validation failed")
	}

	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		logrus.Fatal("invalid 'logLevel' parameter in configuration. Available values: ", logrus.AllLevels)
	}
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})

	return config
}
