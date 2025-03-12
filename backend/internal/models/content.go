package models

import (
	"time"

	"gorm.io/gorm"
)

// ContentType represents types of content
type ContentType string

// Content represents media content (movie or series)
type Content struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"size:255;not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	Type        string     `gorm:"size:20;not null" json:"type"`
	CoverImage  string     `gorm:"size:255" json:"cover_image"`
	ReleaseDate *time.Time `json:"release_date"`
	Duration    *int       `json:"duration"` // in minutes, for movies
	Rating      float32    `gorm:"default:0" json:"rating"`
	SeasonID    *uint      `json:"season_id"`

	// Tambahan field baru
	DownloadLinks []DownloadLink `gorm:"foreignKey:ContentID" json:"download_links"`
	StreamLinks   []StreamLink   `gorm:"foreignKey:ContentID" json:"stream_links"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Episodes   []Episode  `gorm:"foreignKey:ContentID" json:"episodes,omitempty"`
	Genres     []Genre    `gorm:"many2many:content_genres;" json:"genres,omitempty"`
	Categories []Category `gorm:"many2many:content_categories;" json:"categories,omitempty"`
	Season     *Season    `gorm:"foreignKey:SeasonID" json:"season,omitempty"`
}

// DownloadLink represents a download link for content
type DownloadLink struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ContentID     uint           `gorm:"not null" json:"content_id"`
	Name          string         `gorm:"size:100;not null" json:"name"`
	Quality       string         `gorm:"size:20" json:"quality"`
	URL           string         `gorm:"type:text;not null" json:"url"`
	Server        string         `gorm:"size:50;not null;default:'external'" json:"server"`
	EpisodeNumber int            `gorm:"default:1" json:"episode_number"`
	SeasonNumber  int            `gorm:"default:1" json:"season_number"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// StreamLink represents a streaming link for content
type StreamLink struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ContentID     uint           `gorm:"not null" json:"content_id"`
	Name          string         `gorm:"size:100;not null" json:"name"`
	Quality       string         `gorm:"size:20" json:"quality"`
	URL           string         `gorm:"type:text;not null" json:"url"`
	Type          string         `gorm:"size:20;default:'embed'" json:"type"` // 'embed' atau 'self-hosted'
	Server        string         `gorm:"size:50;not null;default:'local'" json:"server"`
	EpisodeNumber int            `gorm:"default:1" json:"episode_number"`
	SeasonNumber  int            `gorm:"default:1" json:"season_number"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Content
func (Content) TableName() string {
	return "contents"
}
