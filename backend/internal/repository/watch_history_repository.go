package repository

import (
	"time"

	"github.com/username/anime-streaming/internal/models"
	"gorm.io/gorm"
)

// WatchHistoryRepository handles database operations for watch history
type WatchHistoryRepository struct {
	db *gorm.DB
}

// NewWatchHistoryRepository creates a new WatchHistoryRepository
func NewWatchHistoryRepository(db *gorm.DB) *WatchHistoryRepository {
	return &WatchHistoryRepository{db: db}
}

// Create creates a new watch history record
func (r *WatchHistoryRepository) Create(history *models.WatchHistory) error {
	history.WatchedAt = time.Now()
	return r.db.Create(history).Error
}

// FindByID finds a watch history record by ID
func (r *WatchHistoryRepository) FindByID(id uint) (*models.WatchHistory, error) {
	var history models.WatchHistory
	if err := r.db.First(&history, id).Error; err != nil {
		return nil, err
	}
	return &history, nil
}

// FindByUserAndContent finds a watch history record by user ID and content ID
func (r *WatchHistoryRepository) FindByUserAndContent(userID, contentID uint) (*models.WatchHistory, error) {
	var history models.WatchHistory
	err := r.db.Where("user_id = ? AND content_id = ?", userID, contentID).
		Order("watched_at DESC").
		First(&history).Error
	
	if err != nil {
		return nil, err
	}
	
	return &history, nil
}

// FindByUserAndEpisode finds a watch history record by user ID and episode ID
func (r *WatchHistoryRepository) FindByUserAndEpisode(userID uint, episodeID uint) (*models.WatchHistory, error) {
	var history models.WatchHistory
	err := r.db.Where("user_id = ? AND episode_id = ?", userID, episodeID).
		Order("watched_at DESC").
		First(&history).Error
	
	if err != nil {
		return nil, err
	}
	
	return &history, nil
}

// Update updates a watch history record
func (r *WatchHistoryRepository) Update(history *models.WatchHistory) error {
	history.WatchedAt = time.Now()
	return r.db.Save(history).Error
}

// Delete deletes a watch history record
func (r *WatchHistoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.WatchHistory{}, id).Error
}

// GetUserHistory gets a user's watch history with pagination
func (r *WatchHistoryRepository) GetUserHistory(userID uint, page, pageSize int) ([]models.WatchHistory, int64, error) {
	var histories []models.WatchHistory
	var count int64
	
	if err := r.db.Model(&models.WatchHistory{}).
		Where("user_id = ?", userID).
		Count(&count).Error; err != nil {
		return nil, 0, err
	}
	
	offset := (page - 1) * pageSize
	err := r.db.Where("user_id = ?", userID).
		Preload("Content").
		Preload("Episode").
		Order("watched_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&histories).Error
	
	return histories, count, err
}

// GetContinueWatching gets content that a user has started but not completed
func (r *WatchHistoryRepository) GetContinueWatching(userID uint, limit int) ([]models.WatchHistory, error) {
	var histories []models.WatchHistory
	
	err := r.db.Where("user_id = ? AND completed_watch = false", userID).
		Preload("Content").
		Preload("Episode").
		Order("watched_at DESC").
		Limit(limit).
		Find(&histories).Error
	
	return histories, err
}

// UpdateProgress updates the watch progress for a movie or episode
func (r *WatchHistoryRepository) UpdateProgress(userID, contentID uint, episodeID *uint, progress int, completed bool) error {
	var history models.WatchHistory
	
	// Try to find an existing record
	query := r.db.Where("user_id = ? AND content_id = ?", userID, contentID)
	if episodeID != nil {
		query = query.Where("episode_id = ?", *episodeID)
	} else {
		query = query.Where("episode_id IS NULL")
	}
	
	err := query.Order("watched_at DESC").First(&history).Error
	
	if err == gorm.ErrRecordNotFound {
		// Create a new record
		history = models.WatchHistory{
			UserID:         userID,
			ContentID:      contentID,
			EpisodeID:      episodeID,
			WatchProgress:  progress,
			CompletedWatch: completed,
			WatchedAt:      time.Now(),
		}
		return r.Create(&history)
	} else if err != nil {
		return err
	}
	
	// Update the existing record
	history.WatchProgress = progress
	history.CompletedWatch = completed
	history.WatchedAt = time.Now()
	
	return r.Update(&history)
} 