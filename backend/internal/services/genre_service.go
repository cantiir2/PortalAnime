package services

import (
	"errors"
	"log"

	"github.com/username/anime-streaming/internal/models"
	"gorm.io/gorm"
)

// GenreService handles genre related operations
type GenreService struct {
	db *gorm.DB
}

// NewGenreService creates a new GenreService
func NewGenreService(db *gorm.DB) *GenreService {
	return &GenreService{
		db: db,
	}
}

// CreateGenre creates a new genre
func (s *GenreService) CreateGenre(genre *models.Genre) error {
	log.Printf("Creating new genre: %s", genre.Name)
	result := s.db.Create(genre)
	if result.Error != nil {
		log.Printf("Failed to create genre: %v", result.Error)
		return result.Error
	}
	return nil
}

// GetGenreByID retrieves a genre by ID
func (s *GenreService) GetGenreByID(id uint) (*models.Genre, error) {
	var genre models.Genre
	result := s.db.First(&genre, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("genre not found")
		}
		return nil, result.Error
	}
	return &genre, nil
}

// ListGenres retrieves all genres
func (s *GenreService) ListGenres() ([]models.Genre, error) {
	var genres []models.Genre
	result := s.db.Find(&genres)
	if result.Error != nil {
		return nil, result.Error
	}
	return genres, nil
}

// UpdateGenre updates an existing genre
func (s *GenreService) UpdateGenre(genre *models.Genre) error {
	result := s.db.Save(genre)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("genre not found")
	}
	return nil
}

// DeleteGenre deletes a genre
func (s *GenreService) DeleteGenre(id uint) error {
	result := s.db.Delete(&models.Genre{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("genre not found")
	}
	return nil
}

// AddGenreToContent adds a genre to a content
func (s *GenreService) AddGenreToContent(contentID, genreID uint) error {
	var content models.Content
	if err := s.db.First(&content, contentID).Error; err != nil {
		return errors.New("content not found")
	}

	var genre models.Genre
	if err := s.db.First(&genre, genreID).Error; err != nil {
		return errors.New("genre not found")
	}

	return s.db.Model(&content).Association("Genres").Append(&genre)
}

// RemoveGenreFromContent removes a genre from a content
func (s *GenreService) RemoveGenreFromContent(contentID, genreID uint) error {
	var content models.Content
	if err := s.db.First(&content, contentID).Error; err != nil {
		return errors.New("content not found")
	}

	var genre models.Genre
	if err := s.db.First(&genre, genreID).Error; err != nil {
		return errors.New("genre not found")
	}

	return s.db.Model(&content).Association("Genres").Delete(&genre)
}

// GetContentGenres retrieves all genres for a content
func (s *GenreService) GetContentGenres(contentID uint) ([]models.Genre, error) {
	var content models.Content
	if err := s.db.First(&content, contentID).Error; err != nil {
		return nil, errors.New("content not found")
	}

	var genres []models.Genre
	if err := s.db.Model(&content).Association("Genres").Find(&genres); err != nil {
		return nil, err
	}

	return genres, nil
}
