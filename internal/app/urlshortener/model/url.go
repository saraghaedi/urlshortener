package model

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

// ErrRecordNotFound represents an error for not finding an error in DB
var ErrRecordNotFound = errors.New("record not found")

const (
	sqlRepoName   = "sql_url"
	nosqlRepoName = "redis_url"
)

// URL represents url table structure.
type URL struct {
	gorm.Model
	URL   string `json:"url"`
	Count int64  `json:"count"`
}

// URLRepo represents repository model.
type URLRepo interface {
	Create(url *URL) error
	FindByID(id uint64) (*URL, error)
	Update(id uint64, additionalCount int64) error
}

// SQLURLRepo represents repository model for SQL databases.
type SQLURLRepo struct {
	MasterDB *gorm.DB
	SlaveDB  *gorm.DB
}

// Create creates a new shorted url in database.
func (s SQLURLRepo) Create(url *URL) (finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(sqlRepoName, "create", startTime, finalErr) }()

	return s.MasterDB.Create(url).Error
}

// FindByID finds a url in database by ID.
func (s SQLURLRepo) FindByID(id uint64) (_ *URL, finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(sqlRepoName, "find_by_id", startTime, finalErr) }()

	var result URL

	if err := s.SlaveDB.Where("id = ?", id).Take(&result).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, ErrRecordNotFound
		}

		return nil, err
	}

	return &result, nil
}

// Update updates count column in database.
func (s SQLURLRepo) Update(id uint64, additionalCount int64) (finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(sqlRepoName, "update", startTime, finalErr) }()

	var u URL

	if err := s.SlaveDB.Where("id = ?", id).Find(&u).Error; err != nil {
		return err
	}

	return s.MasterDB.Model(URL{}).Where("id = ?", id).Update("count", u.Count+additionalCount).Error
}
