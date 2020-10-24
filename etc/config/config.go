package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

type Main struct {
	Port     int    `json:"port"`
	LogLevel string `json:"logLevel"`
}

func New(path string) *Main {
	body, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.WithError(err).WithField("configPath", path).Fatal("can't read config file in the selected path")
	}
	var config Main
	if err := json.Unmarshal(body, &config); err != nil {
		logrus.WithError(err).Fatal("can't unmarshal config file as a json object")
	}

	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		logrus.Fatal("invalid 'logLevel' parameter in configuration. Available values: ", logrus.AllLevels)
	}
	logrus.SetLevel(level)
	logrus.SetReportCaller(true) // adds line number to log message
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})

	return &config
}
