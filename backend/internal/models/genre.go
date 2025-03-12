package models

import (
	"time"

	"gorm.io/gorm"
)

// Genre represents a content genre (action, comedy, etc.)
type Genre struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"unique;not null"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Contents []Content `json:"contents,omitempty" gorm:"many2many:content_genres;"`
}

// TableName specifies the table name for Genre
func (Genre) TableName() string {
	return "genres"
}
