package repository

import (
	"log"

	"github.com/username/anime-streaming/internal/models"
	"gorm.io/gorm"
)

// ContentRepository handles database operations for content (movies/series)
type ContentRepository struct {
	db *gorm.DB
}

// NewContentRepository creates a new ContentRepository
func NewContentRepository(db *gorm.DB) *ContentRepository {
	return &ContentRepository{db: db}
}

// Create creates a new content
func (r *ContentRepository) Create(content *models.Content) error {
	return r.db.Create(content).Error
}

// FindByID finds content by ID with optional preloading of relationships
func (r *ContentRepository) FindByID(id uint, preload ...string) (*models.Content, error) {
	var content models.Content
	query := r.db

	for _, relation := range preload {
		query = query.Preload(relation)
	}

	if err := query.First(&content, id).Error; err != nil {
		return nil, err
	}
	return &content, nil
}

// Update updates content
func (r *ContentRepository) Update(content *models.Content) error {
	return r.db.Save(content).Error
}

// Delete deletes content
func (r *ContentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Content{}, id).Error
}

// List lists all content with pagination and optional filtering
func (r *ContentRepository) List(page, pageSize int, filters map[string]interface{}, preload ...string) ([]models.Content, int64, error) {
	var contents []models.Content
	var count int64
	query := r.db.Model(&models.Content{})

	// Apply filters
	if filters != nil {
		for key, value := range filters {
			// Log the filter being applied
			log.Printf("Applying filter: %s = %v", key, value)
			query = query.Where(key+" = ?", value)
		}
	}

	// Count total items
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Apply preloading
	for _, relation := range preload {
		query = query.Preload(relation)
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&contents).Error; err != nil {
		return nil, 0, err
	}

	return contents, count, nil
}

// Search searches content by title
func (r *ContentRepository) Search(term string, page, pageSize int, preload ...string) ([]models.Content, int64, error) {
	var contents []models.Content
	var count int64
	query := r.db.Model(&models.Content{}).Where("title ILIKE ?", "%"+term+"%")

	// Count total items
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Apply preloading
	for _, relation := range preload {
		query = query.Preload(relation)
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&contents).Error; err != nil {
		return nil, 0, err
	}

	return contents, count, nil
}

// FindByGenre finds content by genre
func (r *ContentRepository) FindByGenre(genreID uint, page, pageSize int, preload ...string) ([]models.Content, int64, error) {
	var contents []models.Content
	var count int64

	subQuery := r.db.Table("content_genres").Where("genre_id = ?", genreID).Select("id")
	query := r.db.Model(&models.Content{}).Where("id IN (?)", subQuery)

	// Count total items
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Apply preloading
	for _, relation := range preload {
		query = query.Preload(relation)
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&contents).Error; err != nil {
		return nil, 0, err
	}

	return contents, count, nil
}

// FindByCategory finds content by category
func (r *ContentRepository) FindByCategory(categoryID uint, page, pageSize int, preload ...string) ([]models.Content, int64, error) {
	var contents []models.Content
	var count int64

	log.Printf("Check Category in findByCategory repo: %s", categoryID)
	subQuery := r.db.Table("categories").Where("id = ?", categoryID).Select("name")
	query := r.db.Model(&models.Content{}).Where("type IN (?)", subQuery)

	// Count total items
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Apply preloading
	for _, relation := range preload {
		query = query.Preload(relation)
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&contents).Error; err != nil {
		return nil, 0, err
	}

	return contents, count, nil
}

// AddStreamLink adds a stream link to content
func (r *ContentRepository) AddStreamLink(contentID uint, streamLink *models.StreamLink) error {
	streamLink.ContentID = contentID
	return r.db.Create(streamLink).Error
}

// AddDownloadLink adds a download link to content
func (r *ContentRepository) AddDownloadLink(contentID uint, downloadLink *models.DownloadLink) error {
	downloadLink.ContentID = contentID
	return r.db.Create(downloadLink).Error
}

// Get DB
func (r *ContentRepository) GetDB() *gorm.DB {
	return r.db
}

// DeleteStreamLinks deletes all stream links for a content
func (r *ContentRepository) DeleteStreamLinks(contentID uint) error {
	return r.db.Where("content_id = ?", contentID).Delete(&models.StreamLink{}).Error
}

// DeleteDownloadLinks deletes all download links for a content
func (r *ContentRepository) DeleteDownloadLinks(contentID uint) error {
	return r.db.Where("content_id = ?", contentID).Delete(&models.DownloadLink{}).Error
}
