package models

import (
	"time"

	"gorm.io/gorm"
)

// Episode represents an episode of a series
type Episode struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ContentID     uint           `gorm:"not null" json:"content_id"`
	Title         string         `gorm:"size:255;not null" json:"title"`
	Description   string         `gorm:"type:text" json:"description"`
	Type          string         `gorm:"size:20;not null;default:'episode'" json:"type"`
	EpisodeNumber int            `gorm:"not null" json:"episode_number"`
	SeasonNumber  int            `gorm:"default:1" json:"season_number"`
	VideoPath     string         `gorm:"size:255;not null" json:"video_path"`
	Duration      int            `gorm:"default:0" json:"duration"` // in seconds
	ThumbnailURL  string         `gorm:"size:255" json:"thumbnail_url"`
	ReleaseDate   *time.Time     `json:"release_date"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Episode
func (Episode) TableName() string {
	return "episodes"
}
