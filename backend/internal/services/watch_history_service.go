package services

import (
	"errors"

	"github.com/username/anime-streaming/internal/models"
	"github.com/username/anime-streaming/internal/repository"
)

// WatchHistoryService handles business logic for watch history
type WatchHistoryService struct {
	watchHistoryRepo  *repository.WatchHistoryRepository
	contentRepo       *repository.ContentRepository
	episodeRepo       *repository.EpisodeRepository
	contentTypeHelper *models.ContentTypeHelper
}

// NewWatchHistoryService creates a new WatchHistoryService
func NewWatchHistoryService(
	watchHistoryRepo *repository.WatchHistoryRepository,
	contentRepo *repository.ContentRepository,
	episodeRepo *repository.EpisodeRepository,
) *WatchHistoryService {
	return &WatchHistoryService{
		watchHistoryRepo:  watchHistoryRepo,
		contentRepo:       contentRepo,
		episodeRepo:       episodeRepo,
		contentTypeHelper: models.NewContentTypeHelper(),
	}
}

// UpdateProgress updates the watch progress for a movie or episode
func (s *WatchHistoryService) UpdateProgress(userID, contentID uint, episodeID *uint, progress int) error {
	// Verify content exists
	content, err := s.contentRepo.FindByID(contentID)
	if err != nil {
		return err
	}

	// For episodes, verify episode exists and belongs to the content
	if episodeID != nil {
		episode, err := s.episodeRepo.FindByID(*episodeID)
		if err != nil {
			return err
		}
		if episode.ContentID != contentID {
			return err
		}

		// Mark as completed if progress is near the end (e.g., 90% or more)
		completed := float64(progress)/float64(episode.Duration) >= 0.9
		return s.watchHistoryRepo.UpdateProgress(userID, contentID, episodeID, progress, completed)
	}

	// For movies
	if !s.contentTypeHelper.IsMovie(content.Type) {
		return errors.New("content must be a movie type")
	}

	// Mark as completed if progress is near the end (e.g., 90% or more)
	completed := float64(progress)/float64(*content.Duration*60) >= 0.9 // content.Duration is in minutes
	return s.watchHistoryRepo.UpdateProgress(userID, contentID, nil, progress, completed)
}

// GetProgress gets the watch progress for a movie or episode
func (s *WatchHistoryService) GetProgress(userID, contentID uint, episodeID *uint) (*models.WatchHistory, error) {
	if episodeID != nil {
		return s.watchHistoryRepo.FindByUserAndEpisode(userID, *episodeID)
	}
	return s.watchHistoryRepo.FindByUserAndContent(userID, contentID)
}

// GetUserHistory gets a user's watch history with pagination
func (s *WatchHistoryService) GetUserHistory(userID uint, page, pageSize int) ([]models.WatchHistory, int64, error) {
	return s.watchHistoryRepo.GetUserHistory(userID, page, pageSize)
}

// GetContinueWatching gets content that a user has started but not completed
func (s *WatchHistoryService) GetContinueWatching(userID uint, limit int) ([]models.WatchHistory, error) {
	return s.watchHistoryRepo.GetContinueWatching(userID, limit)
}

// DeleteHistory deletes a watch history record
func (s *WatchHistoryService) DeleteHistory(id uint) error {
	return s.watchHistoryRepo.Delete(id)
}
