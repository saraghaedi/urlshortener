package database

import (
	"github.com/jinzhu/gorm"
	"github.com/saraghaedi/urlshortener/model"
)

// Migrate runs the database migrations.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.URL{}).Error
}
