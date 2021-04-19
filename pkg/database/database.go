package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // PostgreSQL driver should have blank import
)

// New opens a new connection to the database.
func New(driver string, conStr string) (*gorm.DB, error) {
	db, err := gorm.Open(driver, conStr)
	if err != nil {
		return nil, err
	}

	err = db.DB().Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
