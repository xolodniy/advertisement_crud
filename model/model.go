package model

import (
	"fmt"

	"advertisement_crud/etc/config"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Model struct {
	db *gorm.DB
}

func New(config config.Database) *Model {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.WithError(err).Fatal("failed to connect database")
	}
	sqlDB, err := db.DB()
	if err != nil {
		logrus.WithError(err).Fatal("failed to extract raw db")
	}
	if err := sqlDB.Ping(); err != nil {
		logrus.WithError(err).Fatal("database is not responding")
	}
	return &Model{
		db: db,
	}
}
