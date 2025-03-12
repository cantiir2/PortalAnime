package repository

import (
	"log"

	"github.com/username/anime-streaming/internal/models"
	"gorm.io/gorm"
)

// EpisodeRepository handles database operations for episodes
type EpisodeRepository struct {
	db *gorm.DB
}

// NewEpisodeRepository creates a new EpisodeRepository
func NewEpisodeRepository(db *gorm.DB) *EpisodeRepository {
	return &EpisodeRepository{db: db}
}

// Create creates a new episode
func (r *EpisodeRepository) Create(episode *models.Episode) error {
	return r.db.Create(episode).Error
}

// FindByID finds an episode by ID
func (r *EpisodeRepository) FindByID(id uint) (*models.Episode, error) {
	var episode models.Episode
	if err := r.db.First(&episode, id).Error; err != nil {
		log.Printf("Failed to find episode with ID %d: %v", id, err)
		return nil, err
	}
	log.Printf("Found episode: ID=%d, ContentID=%d, VideoPath=%s", episode.ID, episode.ContentID, episode.VideoPath)
	return &episode, nil
}

// Update updates an episode
func (r *EpisodeRepository) Update(episode *models.Episode) error {
	return r.db.Save(episode).Error
}

// Delete deletes an episode
func (r *EpisodeRepository) Delete(id uint) error {
	return r.db.Delete(&models.Episode{}, id).Error
}

// ListByContentID lists all episodes for a specific content ID with optional season filter
func (r *EpisodeRepository) ListByContentID(contentID uint, season *int) ([]models.Episode, error) {
	var episodes []models.Episode
	query := r.db.Where("content_id = ?", contentID)

	if season != nil {
		query = query.Where("season_number = ?", *season)
	}

	err := query.Order("season_number, episode_number").Find(&episodes).Error
	return episodes, err
}

// GetNextEpisode gets the next episode in a series
func (r *EpisodeRepository) GetNextEpisode(contentID uint, currentSeason, currentEpisode int) (*models.Episode, error) {
	var nextEpisode models.Episode

	// Try to find the next episode in the same season
	err := r.db.Where(
		"content_id = ? AND season_number = ? AND episode_number > ?",
		contentID, currentSeason, currentEpisode,
	).Order("episode_number").First(&nextEpisode).Error

	if err == nil {
		return &nextEpisode, nil
	}

	if err == gorm.ErrRecordNotFound {
		// Try to find the first episode of the next season
		err = r.db.Where(
			"content_id = ? AND season_number > ?",
			contentID, currentSeason,
		).Order("season_number, episode_number").First(&nextEpisode).Error

		if err != nil {
			return nil, err
		}

		return &nextEpisode, nil
	}

	return nil, err
}

// GetLatestEpisode gets the latest episode for a content
func (r *EpisodeRepository) GetLatestEpisode(contentID uint) (*models.Episode, error) {
	var episode models.Episode
	err := r.db.Where("content_id = ?", contentID).
		Order("season_number DESC, episode_number DESC").
		First(&episode).Error

	if err != nil {
		return nil, err
	}

	return &episode, nil
}
