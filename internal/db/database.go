package db

import (
	"github.com/aarongmx/finanzas-personales/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Transaction{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
