package model

import (
	"github.com/jinzhu/gorm"
)

// URL represent url table structure.
type URL struct {
	gorm.Model
	URL string `json:"url"`
}

// URLRepo represent repository model.
type URLRepo interface {
	Create(url *URL) error
	FindByID(id uint64) (*URL, error)
}

// SQLURLRepo represent repository model for SQL databases.
type SQLURLRepo struct {
	MasterDB *gorm.DB
	SlaveDB  *gorm.DB
}

// Create create a new shorted url in database.
func (s SQLURLRepo) Create(url *URL) error {
	return s.MasterDB.Create(url).Error
}

// FindByID find a url in database by ID.
func (s SQLURLRepo) FindByID(id uint64) (*URL, error) {
	var result URL

	if err := s.SlaveDB.Where("id = ?", id).Take(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}
