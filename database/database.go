package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // PostgreSQL driver should have blank import
)

// New opens a new connection to the database.
func New() (*gorm.DB, error) {
	conStr := connectionString()

	db, err := gorm.Open(dbType, conStr)
	if err != nil {
		return nil, err
	}

	err = db.DB().Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
