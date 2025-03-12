package services

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/username/anime-streaming/internal/models"
	"github.com/username/anime-streaming/internal/repository"
)

// EpisodeService handles business logic for episodes
type EpisodeService struct {
	episodeRepo *repository.EpisodeRepository
	contentRepo *repository.ContentRepository
	mediaPath   string
}

// NewEpisodeService creates a new EpisodeService
func NewEpisodeService(
	episodeRepo *repository.EpisodeRepository,
	contentRepo *repository.ContentRepository,
	mediaPath string,
) *EpisodeService {
	return &EpisodeService{
		episodeRepo: episodeRepo,
		contentRepo: contentRepo,
		mediaPath:   mediaPath,
	}
}

// CreateEpisode creates a new episode
func (s *EpisodeService) CreateEpisode(episode *models.Episode) error {
	// Verify content exists
	_, err := s.contentRepo.FindByID(episode.ContentID)
	if err != nil {
		return err
	}

	// Set video path and thumbnail URL
	if episode.VideoPath != "" {
		// Remove any existing path prefixes
		episode.VideoPath = filepath.Base(episode.VideoPath)
		// Set the correct path
		episode.VideoPath = filepath.Join("videos", "original", episode.VideoPath)
	}
	if episode.ThumbnailURL != "" {
		// Remove any existing path prefixes
		episode.ThumbnailURL = filepath.Base(episode.ThumbnailURL)
		// Set the correct path
		episode.ThumbnailURL = filepath.Join("thumbnails", "episodes", episode.ThumbnailURL)
	}

	return s.episodeRepo.Create(episode)
}

// GetEpisodeByID retrieves an episode by ID
func (s *EpisodeService) GetEpisodeByID(id uint) (*models.Episode, error) {
	return s.episodeRepo.FindByID(id)
}

// UpdateEpisode updates episode information
func (s *EpisodeService) UpdateEpisode(episode *models.Episode) error {
	// Verify content exists and is a series
	content, err := s.contentRepo.FindByID(episode.ContentID)
	if err != nil {
		return err
	}
	if content.Type != "Series" && content.Type != "Anime" {
		return errors.New("content must be a series or anime")
	}

	// Update paths if provided
	if episode.VideoPath != "" {
		// Remove any existing path prefixes
		episode.VideoPath = strings.TrimPrefix(episode.VideoPath, "media\\videos\\")
		episode.VideoPath = strings.TrimPrefix(episode.VideoPath, "media/videos/")
		episode.VideoPath = strings.TrimPrefix(episode.VideoPath, "videos\\")
		episode.VideoPath = strings.TrimPrefix(episode.VideoPath, "videos/")
		// Set the correct path
		episode.VideoPath = filepath.Join("videos", "original", filepath.Base(episode.VideoPath))
	}
	if episode.ThumbnailURL != "" {
		// Remove any existing path prefixes
		episode.ThumbnailURL = strings.TrimPrefix(episode.ThumbnailURL, "media\\thumbnails\\")
		episode.ThumbnailURL = strings.TrimPrefix(episode.ThumbnailURL, "media/thumbnails/")
		episode.ThumbnailURL = strings.TrimPrefix(episode.ThumbnailURL, "thumbnails\\")
		episode.ThumbnailURL = strings.TrimPrefix(episode.ThumbnailURL, "thumbnails/")
		// Set the correct path
		episode.ThumbnailURL = filepath.Join("thumbnails", "episodes", filepath.Base(episode.ThumbnailURL))
	}

	return s.episodeRepo.Update(episode)
}

// DeleteEpisode deletes an episode
func (s *EpisodeService) DeleteEpisode(id uint) error {
	return s.episodeRepo.Delete(id)
}

// ListEpisodes lists all episodes for a content with optional season filter
func (s *EpisodeService) ListEpisodes(contentID uint, season *int) ([]models.Episode, error) {
	return s.episodeRepo.ListByContentID(contentID, season)
}

// GetNextEpisode gets the next episode in a series
func (s *EpisodeService) GetNextEpisode(contentID uint, currentSeason, currentEpisode int) (*models.Episode, error) {
	return s.episodeRepo.GetNextEpisode(contentID, currentSeason, currentEpisode)
}

// GetLatestEpisode gets the latest episode for a series
func (s *EpisodeService) GetLatestEpisode(contentID uint) (*models.Episode, error) {
	return s.episodeRepo.GetLatestEpisode(contentID)
}

// GetContent gets content by ID
func (s *EpisodeService) GetContent(contentID uint) (*models.Content, error) {
	content, err := s.contentRepo.FindByID(contentID)
	if err != nil {
		return nil, err
	}
	return content, nil
}
