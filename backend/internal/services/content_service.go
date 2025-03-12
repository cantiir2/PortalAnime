package services

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/username/anime-streaming/internal/models"
	"github.com/username/anime-streaming/internal/repository"
	"gorm.io/gorm"
)

// ContentService handles business logic for content
type ContentService struct {
	contentRepo       *repository.ContentRepository
	genreRepo         *repository.GenreRepository
	categoryRepo      *repository.CategoryRepository
	mediaPath         string
	contentTypeHelper *models.ContentTypeHelper
}

// NewContentService creates a new ContentService
func NewContentService(
	contentRepo *repository.ContentRepository,
	genreRepo *repository.GenreRepository,
	categoryRepo *repository.CategoryRepository,
	mediaPath string,
) *ContentService {
	return &ContentService{
		contentRepo:       contentRepo,
		genreRepo:         genreRepo,
		categoryRepo:      categoryRepo,
		mediaPath:         mediaPath,
		contentTypeHelper: models.NewContentTypeHelper(),
	}
}

// CreateContent creates a new content entry
func (s *ContentService) CreateContent(content *models.Content) error {
	// Validate content type against existing categories
	categories, err := s.categoryRepo.List()
	if err != nil {
		return fmt.Errorf("failed to fetch categories: %v", err)
	}

	if !s.contentTypeHelper.IsValidType(content.Type, categories) {
		return errors.New("invalid content type: must match an existing category")
	}

	// Set cover image path if provided
	if content.CoverImage != "" {
		// Remove any existing path prefixes
		content.CoverImage = filepath.Base(content.CoverImage)
		// Set the correct path
		content.CoverImage = filepath.Join("thumbnails", "content", content.CoverImage)
	}

	return s.contentRepo.Create(content)
}

// GetContentByID retrieves a content by its ID
func (s *ContentService) GetContentByID(id uint) (*models.Content, error) {
	// Preload all relationships
	return s.contentRepo.FindByID(id, "Episodes", "Genres", "Categories", "StreamLinks", "DownloadLinks")
}

// UpdateContent updates content information
func (s *ContentService) UpdateContent(content *models.Content) error {
	// Validate content type against existing categories
	categories, err := s.categoryRepo.List()
	if err != nil {
		return fmt.Errorf("failed to fetch categories: %v", err)
	}

	if !s.contentTypeHelper.IsValidType(content.Type, categories) {
		return errors.New("invalid content type: must match an existing category")
	}

	fmt.Println("check content", content)

	// Update cover image path if provided
	if content.CoverImage != "" {
		// Remove any existing path prefixes
		content.CoverImage = strings.TrimPrefix(content.CoverImage, "media\\thumbnails\\")
		content.CoverImage = strings.TrimPrefix(content.CoverImage, "media/thumbnails/")
		content.CoverImage = strings.TrimPrefix(content.CoverImage, "thumbnails\\")
		content.CoverImage = strings.TrimPrefix(content.CoverImage, "thumbnails/")
		// Set the correct path
		content.CoverImage = filepath.Join("thumbnails", "content", filepath.Base(content.CoverImage))
	}

	return s.contentRepo.Update(content)
}

// DeleteContent deletes content
func (s *ContentService) DeleteContent(id uint) error {
	log.Printf("Service: Attempting to delete content with ID: %d", id)

	// Verify content exists first
	content, err := s.contentRepo.FindByID(id)
	if err != nil {
		log.Printf("Service: Content not found: %v", err)
		return fmt.Errorf("content not found: %v", err)
	}

	// Delete related records first (if any)
	// Note: This might not be necessary if you've set up CASCADE DELETE in your database
	if err := s.contentRepo.DeleteStreamLinks(id); err != nil {
		log.Printf("Service: Failed to delete stream links: %v", err)
	}
	if err := s.contentRepo.DeleteDownloadLinks(id); err != nil {
		log.Printf("Service: Failed to delete download links: %v", err)
	}

	// Delete the content
	if err := s.contentRepo.Delete(content.ID); err != nil {
		log.Printf("Service: Failed to delete content: %v", err)
		return fmt.Errorf("failed to delete content: %v", err)
	}

	log.Printf("Service: Content %d deleted successfully", id)
	return nil
}

// ListContent lists all content with pagination and filtering
func (s *ContentService) ListContent(page, pageSize int, filters map[string]interface{}) ([]models.Content, int64, error) {
	return s.contentRepo.List(page, pageSize, filters, "Episodes", "Genres", "Categories", "StreamLinks", "DownloadLinks")
}

// SearchContent searches content by title
func (s *ContentService) SearchContent(term string, page, pageSize int) ([]models.Content, int64, error) {
	return s.contentRepo.Search(term, page, pageSize, "Episodes", "Genres", "Categories", "StreamLinks", "DownloadLinks")
}

// GetContentByGenre gets content by genre
func (s *ContentService) GetContentByGenre(genreID uint, page, pageSize int) ([]models.Content, int64, error) {
	return s.contentRepo.FindByGenre(genreID, page, pageSize, "Episodes", "Genres", "Categories", "StreamLinks", "DownloadLinks")
}

// GetContentByCategory gets content by category
func (s *ContentService) GetContentByCategory(categoryID uint, page, pageSize int) ([]models.Content, int64, error) {
	return s.contentRepo.FindByCategory(categoryID, page, pageSize, "Episodes", "Genres", "Categories", "Season", "StreamLinks", "DownloadLinks")
}

// AddGenreToContent adds a genre to content
func (s *ContentService) AddGenreToContent(contentID, genreID uint) error {
	// Verify content exists
	if _, err := s.contentRepo.FindByID(contentID); err != nil {
		return err
	}

	// Verify genre exists
	if _, err := s.genreRepo.FindByID(genreID); err != nil {
		return err
	}

	return s.genreRepo.AddContentToGenre(genreID, contentID)
}

// RemoveGenreFromContent removes a genre from content
func (s *ContentService) RemoveGenreFromContent(contentID, genreID uint) error {
	return s.genreRepo.RemoveContentFromGenre(genreID, contentID)
}

// AddCategoryToContent adds a category to content
func (s *ContentService) AddCategoryToContent(contentID, categoryID uint) error {
	// Verify content exists
	if _, err := s.contentRepo.FindByID(contentID); err != nil {
		return err
	}

	// Verify category exists
	if _, err := s.categoryRepo.FindByID(categoryID); err != nil {
		return err
	}

	return s.categoryRepo.AddContentToCategory(categoryID, contentID)
}

// RemoveCategoryFromContent removes a category from content
func (s *ContentService) RemoveCategoryFromContent(contentID, categoryID uint) error {
	return s.categoryRepo.RemoveContentFromCategory(categoryID, contentID)
}

// AddStreamLink adds a stream link to content
func (s *ContentService) AddStreamLink(contentID uint, streamLink *models.StreamLink) error {
	return s.contentRepo.AddStreamLink(contentID, streamLink)
}

// AddDownloadLink adds a download link to content
func (s *ContentService) AddDownloadLink(contentID uint, downloadLink *models.DownloadLink) error {
	return s.contentRepo.AddDownloadLink(contentID, downloadLink)
}

// GetDB returns the database instance
func (s *ContentService) GetDB() *gorm.DB {
	return s.contentRepo.GetDB()
}

// DeleteContentEpisodes deletes all episodes for a given content ID
func (s *ContentService) DeleteContentEpisodes(contentID uint) error {
	// Delete all episodes for this content
	result := s.contentRepo.GetDB().Where("content_id = ?", contentID).Delete(&models.Episode{})
	if result.Error != nil {
		return result.Error
	}

	// Delete associated stream links and download links
	if err := s.contentRepo.DeleteStreamLinks(contentID); err != nil {
		return err
	}

	if err := s.contentRepo.DeleteDownloadLinks(contentID); err != nil {
		return err
	}

	return nil
}

// GetContentByTitle retrieves content by title
func (s *ContentService) GetContentByTitle(title string) (*models.Content, error) {
	var content models.Content
	err := s.contentRepo.GetDB().Where("title = ?", title).First(&content).Error
	if err != nil {
		return nil, err
	}
	return &content, nil
}
