package models

import (
	"time"

	"gorm.io/gorm"
)

// WatchHistory represents a user's watch history record
type WatchHistory struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	UserID         uint           `gorm:"not null" json:"user_id"`
	ContentID      uint           `gorm:"not null" json:"content_id"`
	EpisodeID      *uint          `json:"episode_id"`
	WatchProgress  int            `gorm:"default:0" json:"watch_progress"` // in seconds
	WatchedAt      time.Time      `json:"watched_at"`
	CompletedWatch bool           `gorm:"default:false" json:"completed_watch"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships (optional for eager loading)
	Content  *Content  `gorm:"foreignKey:ContentID" json:"content,omitempty"`
	Episode  *Episode  `gorm:"foreignKey:EpisodeID" json:"episode,omitempty"`
	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for WatchHistory
func (WatchHistory) TableName() string {
	return "watch_history"
} 