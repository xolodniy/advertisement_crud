package model

import (
	"advertisement_crud/etc/common"
	"errors"
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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
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

func (m *Model) Migrate() {
	tx := m.db.Begin()

	err := tx.Exec(`
		CREATE TABLE advertisements (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP NOT NULL DEFAULT now(),
			updated_at TIMESTAMP NOT NULL DEFAULT now(),
			deleted_at TIMESTAMP NULL,
	
			caption     TEXT NOT NULL,
			description TEXT NOT NULL,
			price       INTEGER NOT NULL
		);
	`).Error
	if err != nil {
		logrus.WithError(err).Fatal("can't migrate advertisements table")
	}

	err = tx.Exec(`
		CREATE TABLE photos (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP NOT NULL DEFAULT now(),
			updated_at TIMESTAMP NOT NULL DEFAULT now(),
			deleted_at TIMESTAMP NULL,
			
			mime TEXT NOT NULL,
			body MEDIUMBLOB
		);
		`).Error
	if err != nil {
		logrus.WithError(err).Fatal("can't migrate photos table")
	}

	err = tx.Exec(`
	CREATE TABLE photos_advertisements (
		photo_id         INTEGER REFERENCES photos,
		advertisement_id INTEGER REFERENCES advertisements,
	
		PRIMARY KEY(photo_id, advertisement_id)
	);
	`).Error
	if err != nil {
		logrus.WithError(err).Fatal("can't migrate m2m table")
	}
	if err := tx.Commit().Error; err != nil {
		logrus.WithError(err).Fatal("can't commit migration")
	}
}

func (m *Model) IsMigrated() bool {
	return m.db.Migrator().HasTable(&Advertisement{})
}

type Advertisement struct {
	gorm.Model

	Caption     string
	Description string
	Price       int

	Photos []Photo `gorm:"many2many:photos_advertisements;jointable_foreignkey:advertisement_id"`
}

type Photo struct {
	gorm.Model

	Mime string
	Body []byte
}

type PhotosAdvertisements struct {
	PhotoID         int `gorm:"primary_key"`
	AdvertisementID int `gorm:"primary_key"`
}

func (m *Model) GetAdvertisement(id int) (Advertisement, error) {
	var response Advertisement
	if err := m.db.Preload("Photos", func(db *gorm.DB) *gorm.DB {
		return db.Omit("body")
	}).First(&response, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Advertisement{}, common.ErrNotFound
		}
		logrus.WithError(err).WithField("id", id).Error("can't get advertisement from database")
		return Advertisement{}, common.ErrInternal
	}
	return response, nil
}

func (m *Model) GetAdvertisements(order string) ([]Advertisement, error) {
	var response []Advertisement
	err := m.db.
		Order(order).
		Preload("Photos", func(db *gorm.DB) *gorm.DB {
			return db.Omit("body")
		}).
		Find(&response).
		Error
	if err != nil {
		logrus.WithField("order", order).WithError(err).Error("can't get advertisements")
		return nil, common.ErrInternal
	}
	return response, nil
}

func (m *Model) CreateAdvertisement(advertisement Advertisement, photoIDs []int) (int, error) {
	tx := m.db.Begin()
	if err := tx.Create(&advertisement).Error; err != nil {
		tx.Rollback()
		logrus.WithError(err).WithField("advertisement", fmt.Sprintf("%+v", advertisement)).
			Error("can't create advertisement")
		return 0, common.ErrInternal
	}
	for _, id := range photoIDs {
		err := tx.Create(&PhotosAdvertisements{
			PhotoID:         id,
			AdvertisementID: int(advertisement.ID),
		}).Error
		if err != nil {
			tx.Rollback()
			logrus.WithError(err).WithFields(logrus.Fields{
				"advertisement": fmt.Sprintf("%+v", advertisement),
				"photoID":       id,
			}).Error("can't create link between advertisement and photo")
			return 0, common.ErrInternal
		}
	}
	if err := tx.Commit().Error; err != nil {
		logrus.WithError(err).Error("can't commit transaction")
		return 0, common.ErrInternal
	}
	return int(advertisement.ID), nil
}

func (m *Model) CheckPhotoExist(id int) bool {
	var i int64
	if err := m.db.Model(Photo{}).Where("id = ?", id).Count(&i).Error; err != nil {
		logrus.WithError(err).WithField("id", id).Error("can't count photos in database")
		return false
	}
	return i == 1
}

func (m *Model) CreatePhoto(photo Photo) (int, error) {
	if err := m.db.Create(&photo).Error; err != nil {
		logrus.WithError(err).Error("can't create photo")
		return 0, common.ErrInternal
	}
	return int(photo.ID), nil
}

func (m *Model) GetPhoto(id int) (Photo, error) {
	var response Photo
	if err := m.db.First(&response, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Photo{}, common.ErrNotFound
		}
		logrus.WithError(err).WithField("id", id).Error("can't get photo from database")
		return Photo{}, common.ErrInternal
	}
	return response, nil
}
