package models

import (
	"time"

	"gorm.io/gorm"
)

// Season represents an anime season
type Season struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:50;not null" json:"name"` // Winter, Spring, Summer, Fall
	Year      int            `gorm:"not null" json:"year"`
	Status    string         `gorm:"size:20;not null" json:"status"` // Coming Soon, Active, Ended
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Contents []Content `gorm:"foreignKey:SeasonID" json:"contents,omitempty"`
}

// TableName specifies the table name for Season
func (Season) TableName() string {
	return "seasons"
}
