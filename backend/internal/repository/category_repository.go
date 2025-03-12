package repository

import (
	"github.com/username/anime-streaming/internal/models"
	"gorm.io/gorm"
)

// CategoryRepository handles database operations for categories
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new CategoryRepository
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// Create creates a new category
func (r *CategoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

// FindByID finds a category by ID
func (r *CategoryRepository) FindByID(id uint) (*models.Category, error) {
	var category models.Category
	if err := r.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// FindByName finds a category by name
func (r *CategoryRepository) FindByName(name string) (*models.Category, error) {
	var category models.Category
	if err := r.db.Where("name = ?", name).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// Update updates a category
func (r *CategoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

// Delete deletes a category
func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}

// List lists all categories
func (r *CategoryRepository) List() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Order("name").Find(&categories).Error
	return categories, err
}

// AddContentToCategory adds content to a category
func (r *CategoryRepository) AddContentToCategory(categoryID, contentID uint) error {
	return r.db.Exec(
		"INSERT INTO content_categories (category_id, content_id) VALUES (?, ?) ON CONFLICT DO NOTHING",
		categoryID, contentID,
	).Error
}

// RemoveContentFromCategory removes content from a category
func (r *CategoryRepository) RemoveContentFromCategory(categoryID, contentID uint) error {
	return r.db.Exec(
		"DELETE FROM content_categories WHERE category_id = ? AND content_id = ?",
		categoryID, contentID,
	).Error
} 