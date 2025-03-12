package repository

import (
	"github.com/username/anime-streaming/internal/models"
	"gorm.io/gorm"
)

// SeasonRepository handles database operations for seasons
type SeasonRepository struct {
	db *gorm.DB
}

// NewSeasonRepository creates a new SeasonRepository
func NewSeasonRepository(db *gorm.DB) *SeasonRepository {
	return &SeasonRepository{db: db}
}

// Create creates a new season
func (r *SeasonRepository) Create(season *models.Season) error {
	return r.db.Create(season).Error
}

// FindByID finds a season by ID
func (r *SeasonRepository) FindByID(id uint) (*models.Season, error) {
	var season models.Season
	if err := r.db.First(&season, id).Error; err != nil {
		return nil, err
	}
	return &season, nil
}

// Update updates a season
func (r *SeasonRepository) Update(season *models.Season) error {
	return r.db.Save(season).Error
}

// Delete deletes a season
func (r *SeasonRepository) Delete(id uint) error {
	return r.db.Delete(&models.Season{}, id).Error
}

// List lists all seasons
func (r *SeasonRepository) List() ([]models.Season, error) {
	var seasons []models.Season
	err := r.db.Order("year DESC, CASE name " +
		"WHEN 'Winter' THEN 1 " +
		"WHEN 'Spring' THEN 2 " +
		"WHEN 'Summer' THEN 3 " +
		"WHEN 'Fall' THEN 4 " +
		"END").Find(&seasons).Error
	return seasons, err
}

// FindByNameAndYear finds a season by name and year
func (r *SeasonRepository) FindByNameAndYear(name string, year int) (*models.Season, error) {
	var season models.Season
	if err := r.db.Where("name = ? AND year = ?", name, year).First(&season).Error; err != nil {
		return nil, err
	}
	return &season, nil
}

// GetCurrentSeason gets the current active season
func (r *SeasonRepository) GetCurrentSeason() (*models.Season, error) {
	var season models.Season
	if err := r.db.Where("status = ?", "Active").First(&season).Error; err != nil {
		return nil, err
	}
	return &season, nil
}
