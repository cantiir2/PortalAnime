package models

import (
	"time"

	"gorm.io/gorm"
)

// Category represents a content category (movies, TV shows, anime, etc.)
type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null;unique" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Contents []Content `gorm:"many2many:content_categories;" json:"contents,omitempty"`
}

// TableName specifies the table name for Category
func (Category) TableName() string {
	return "categories"
} 