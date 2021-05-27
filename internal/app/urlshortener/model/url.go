package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

const (
	repoName = "sql_URL"
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
func (s SQLURLRepo) Create(url *URL) (finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(repoName, "create", startTime, finalErr) }()

	return s.MasterDB.Create(url).Error
}

// FindByID find a url in database by ID.
func (s SQLURLRepo) FindByID(id uint64) (_ *URL, finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(repoName, "find_by_id", startTime, finalErr) }()

	var result URL

	if err := s.SlaveDB.Where("id = ?", id).Take(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}
