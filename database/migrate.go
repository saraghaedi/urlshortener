package database

import (
	"github.com/jinzhu/gorm"
	"github.com/saraghaedi/urlshortener/model"
)

// Migrate function will do the database migration
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.URL{}).Error
}
