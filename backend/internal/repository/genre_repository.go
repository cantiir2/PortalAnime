package repository

import (
	"github.com/username/anime-streaming/internal/models"
	"gorm.io/gorm"
)

// GenreRepository handles database operations for genres
type GenreRepository struct {
	db *gorm.DB
}

// NewGenreRepository creates a new GenreRepository
func NewGenreRepository(db *gorm.DB) *GenreRepository {
	return &GenreRepository{db: db}
}

// Create creates a new genre
func (r *GenreRepository) Create(genre *models.Genre) error {
	return r.db.Create(genre).Error
}

// FindByID finds a genre by ID
func (r *GenreRepository) FindByID(id uint) (*models.Genre, error) {
	var genre models.Genre
	if err := r.db.First(&genre, id).Error; err != nil {
		return nil, err
	}
	return &genre, nil
}

// FindByName finds a genre by name
func (r *GenreRepository) FindByName(name string) (*models.Genre, error) {
	var genre models.Genre
	if err := r.db.Where("name = ?", name).First(&genre).Error; err != nil {
		return nil, err
	}
	return &genre, nil
}

// Update updates a genre
func (r *GenreRepository) Update(genre *models.Genre) error {
	return r.db.Save(genre).Error
}

// Delete deletes a genre
func (r *GenreRepository) Delete(id uint) error {
	return r.db.Delete(&models.Genre{}, id).Error
}

// List lists all genres
func (r *GenreRepository) List() ([]models.Genre, error) {
	var genres []models.Genre
	err := r.db.Order("name").Find(&genres).Error
	return genres, err
}

// AddContentToGenre adds content to a genre
func (r *GenreRepository) AddContentToGenre(genreID, contentID uint) error {
	return r.db.Exec(
		"INSERT INTO content_genres (genre_id, content_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
		genreID, contentID,
	).Error
}

// RemoveContentFromGenre removes content from a genre
func (r *GenreRepository) RemoveContentFromGenre(genreID, contentID uint) error {
	return r.db.Exec(
		"DELETE FROM content_genres WHERE genre_id = ? AND content_id = ?",
		genreID, contentID,
	).Error
} 