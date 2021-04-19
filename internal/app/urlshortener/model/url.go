package model

import (
	"github.com/jinzhu/gorm"
)

// URL is a model and its attributes comes from gorm.model.
type URL struct {
	gorm.Model
	URL string
}
