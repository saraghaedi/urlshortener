package model

import (
	"github.com/jinzhu/gorm"
)

type URL struct {
	gorm.Model
	URL string `json:"url"`
}

type URLRepo interface {
	Create(url *URL) error
	FindByID(id uint64) (*URL, error)
}

type SQLURLRepo struct {
	Driver   string
	MasterDB *gorm.DB
	SlaveDB  *gorm.DB
}

func (s SQLURLRepo) Create(url *URL) error {
	return s.MasterDB.Create(url).Error
}

func (s SQLURLRepo) FindByID(id uint64) (*URL, error) {
	var result URL

	if err := s.SlaveDB.Where("id = ?", id).Take(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}
