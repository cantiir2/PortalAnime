package services

import (
	"errors"
	"log"

	"github.com/username/anime-streaming/internal/models"
	"github.com/username/anime-streaming/internal/repository"
)

// SeasonService handles business logic for seasons
type SeasonService struct {
	seasonRepo *repository.SeasonRepository
}

// NewSeasonService creates a new SeasonService
func NewSeasonService(seasonRepo *repository.SeasonRepository) *SeasonService {
	return &SeasonService{
		seasonRepo: seasonRepo,
	}
}

// CreateSeason creates a new season
func (s *SeasonService) CreateSeason(season *models.Season) error {
	log.Printf("Creating new season: %s %d", season.Name, season.Year)

	// Validate season name
	if !isValidSeasonName(season.Name) {
		return errors.New("invalid season name. Must be Winter, Spring, Summer, or Fall")
	}

	// Validate status
	if !isValidStatus(season.Status) {
		return errors.New("invalid status. Must be Coming Soon, Active, or Ended")
	}

	// Check if season already exists
	existing, err := s.seasonRepo.FindByNameAndYear(season.Name, season.Year)
	if err == nil && existing != nil {
		return errors.New("season already exists")
	}

	return s.seasonRepo.Create(season)
}

// GetSeasonByID retrieves a season by ID
func (s *SeasonService) GetSeasonByID(id uint) (*models.Season, error) {
	return s.seasonRepo.FindByID(id)
}

// ListSeasons retrieves all seasons
func (s *SeasonService) ListSeasons() ([]models.Season, error) {
	return s.seasonRepo.List()
}

// UpdateSeason updates an existing season
func (s *SeasonService) UpdateSeason(season *models.Season) error {
	// Validate season name
	if !isValidSeasonName(season.Name) {
		return errors.New("invalid season name. Must be Winter, Spring, Summer, or Fall")
	}

	// Validate status
	if !isValidStatus(season.Status) {
		return errors.New("invalid status. Must be Coming Soon, Active, or Ended")
	}

	// Check if season exists
	existing, err := s.seasonRepo.FindByID(season.ID)
	if err != nil {
		return errors.New("season not found")
	}

	// Check if another season with same name and year exists
	if existing.Name != season.Name || existing.Year != season.Year {
		duplicate, err := s.seasonRepo.FindByNameAndYear(season.Name, season.Year)
		if err == nil && duplicate != nil && duplicate.ID != season.ID {
			return errors.New("another season with same name and year already exists")
		}
	}

	return s.seasonRepo.Update(season)
}

// DeleteSeason deletes a season
func (s *SeasonService) DeleteSeason(id uint) error {
	return s.seasonRepo.Delete(id)
}

// GetCurrentSeason gets the current active season
func (s *SeasonService) GetCurrentSeason() (*models.Season, error) {
	return s.seasonRepo.GetCurrentSeason()
}

// Helper functions
func isValidSeasonName(name string) bool {
	validNames := map[string]bool{
		"Winter": true,
		"Spring": true,
		"Summer": true,
		"Fall":   true,
	}
	return validNames[name]
}

func isValidStatus(status string) bool {
	validStatuses := map[string]bool{
		"Coming Soon": true,
		"Active":      true,
		"Ended":       true,
	}
	return validStatuses[status]
}
