package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/username/anime-streaming/internal/models"
	"github.com/username/anime-streaming/internal/services"
)

// ContentHandler handles content related requests
type ContentHandler struct {
	contentService *services.ContentService
	mediaService   *services.MediaService
}

// NewContentHandler creates a new ContentHandler
func NewContentHandler(contentService *services.ContentService, mediaService *services.MediaService) *ContentHandler {
	return &ContentHandler{
		contentService: contentService,
		mediaService:   mediaService,
	}
}

// CreateContentRequest represents the request body for creating content
type CreateContentRequest struct {
	Title         string             `form:"title" binding:"required"`
	Description   string             `form:"description"`
	Type          models.ContentType `form:"type" binding:"required"`
	ReleaseDate   *time.Time         `form:"releaseDate"`
	GenreIds      []uint             `form:"genreIds[]"` // Ubah ini untuk menerima array
	StreamLinks   string             `form:"streamLinks"`
	DownloadLinks string             `form:"downloadLinks"`
	Episodes      string             `form:"episodes"`
	Rating        float32            `form:"rating"`
	SeasonID      *uint              `form:"season_id"`
}

type StreamLinkRequest struct {
	Name       string `json:"name"`
	Type       string `json:"type"` // 'embed' atau 'self-hosted'
	Quality    string `json:"quality"`
	URL        string `json:"url,omitempty"`
	VideoField string `json:"videoField,omitempty"`
	FileKey    string `json:"fileKey,omitempty"`
	FileType   string `json:"fileType,omitempty"`
}

type DownloadLinkRequest struct {
	Name    string `json:"name"`
	Quality string `json:"quality"`
	URL     string `json:"url"`
}

// Create creates or updates content
func (h *ContentHandler) Create(c *gin.Context) {
	var input CreateContentRequest
	if err := c.ShouldBind(&input); err != nil {
		log.Printf("Failed to bind request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if content already exists by title
	existingContent, err := h.contentService.GetContentByTitle(input.Title)
	isUpdate := err == nil && existingContent != nil

	content := &models.Content{
		Title:       input.Title,
		Description: input.Description,
		Type:        string(input.Type),
		ReleaseDate: input.ReleaseDate,
		Rating:      input.Rating,
		SeasonID:    input.SeasonID,
	}

	if isUpdate {
		content.ID = existingContent.ID
		// Update existing content
		if err := h.contentService.UpdateContent(content); err != nil {
			log.Printf("Failed to update content: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update content: %v", err)})
			return
		}
	} else {
		// Create new content
		if err := h.contentService.CreateContent(content); err != nil {
			log.Printf("Failed to create content: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create content: %v", err)})
			return
		}
	}

	// Handle cover image upload if provided
	if file, err := c.FormFile("coverImage"); err == nil {
		log.Printf("Cover image found, uploading for content ID: %d, filename: %s", content.ID, file.Filename)
		if err := h.mediaService.UploadContentCover(content.ID, file); err != nil {
			log.Printf("Failed to upload cover image: %v", err)
			// Don't return error here, continue with other operations
		} else {
			log.Printf("Cover image uploaded successfully for content ID: %d", content.ID)

			// Reload content to get updated cover image path
			updatedContent, err := h.contentService.GetContentByID(content.ID)
			if err == nil {
				content.CoverImage = updatedContent.CoverImage
				log.Printf("Updated content with cover image path: %s", content.CoverImage)
			}
		}
	} else {
		log.Printf("No cover image provided or error getting cover: %v", err)
	}

	// Process genres
	if len(input.GenreIds) > 0 {
		log.Printf("Processing genres for content %d: %v", content.ID, input.GenreIds)

		// Clear existing genres first
		if err := h.contentService.GetDB().Model(content).Association("Genres").Clear(); err != nil {
			log.Printf("Failed to clear existing genres: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to clear genres: %v", err)})
			return
		}

		// Add new genres
		var genres []models.Genre
		if err := h.contentService.GetDB().Where("id IN ?", input.GenreIds).Find(&genres).Error; err != nil {
			log.Printf("Failed to fetch genres: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch genres: %v", err)})
			return
		}

		if err := h.contentService.GetDB().Model(content).Association("Genres").Replace(genres); err != nil {
			log.Printf("Failed to update genres: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update genres: %v", err)})
			return
		}

		log.Printf("Successfully updated genres for content %d", content.ID)
	}

	// Process episodes if provided
	if input.Episodes != "" {
		var episodes []struct {
			Title         string                `json:"title"`
			Description   string                `json:"description"`
			EpisodeNumber int                   `json:"episodeNumber"`
			SeasonNumber  int                   `json:"seasonNumber"`
			StreamLinks   []StreamLinkRequest   `json:"streamLinks"`
			DownloadLinks []DownloadLinkRequest `json:"downloadLinks"`
			VideoField    string                `json:"videoField"`
		}

		if err := json.Unmarshal([]byte(input.Episodes), &episodes); err != nil {
			log.Printf("Failed to parse episodes: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to parse episodes: %v", err)})
			return
		}

		// Delete existing episodes and their links if updating
		if isUpdate {
			if err := h.contentService.DeleteContentEpisodes(content.ID); err != nil {
				log.Printf("Failed to delete existing episodes: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete existing episodes: %v", err)})
				return
			}
		}

		// Create new episodes
		for _, ep := range episodes {
			episode := &models.Episode{
				ContentID:     content.ID,
				Title:         ep.Title,
				Description:   ep.Description,
				Type:          "episode",
				EpisodeNumber: ep.EpisodeNumber,
				SeasonNumber:  ep.SeasonNumber,
			}

			// Create episode
			if err := h.contentService.GetDB().Create(episode).Error; err != nil {
				log.Printf("Failed to create episode: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create episode: %v", err)})
				return
			}

			// Process stream links
			for _, sl := range ep.StreamLinks {
				// Handle video upload for self-hosted streams
				if sl.Type == "self-hosted" && sl.VideoField != "" {
					file, err := c.FormFile(sl.VideoField)
					if err != nil {
						log.Printf("No video file found for field %s: %v", sl.VideoField, err)
						continue
					}

					log.Printf("Processing video upload for episode %d", ep.EpisodeNumber)
					videoPath, err := h.mediaService.UploadVideo(content.ID, &episode.ID, file)
					if err != nil {
						log.Printf("Failed to upload video: %v", err)
						c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload video: %v", err)})
						return
					}

					// Update episode with video path
					episode.VideoPath = videoPath
					if err := h.contentService.GetDB().Save(episode).Error; err != nil {
						log.Printf("Failed to update episode with video path: %v", err)
						c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update episode: %v", err)})
						return
					}

					// Create stream link with video path
					streamLink := &models.StreamLink{
						ContentID:     content.ID,
						Name:          sl.Name,
						Quality:       sl.Quality,
						Type:          sl.Type,
						Server:        "local",
						URL:           videoPath,
						EpisodeNumber: ep.EpisodeNumber,
						SeasonNumber:  ep.SeasonNumber,
					}

					if err := h.contentService.AddStreamLink(content.ID, streamLink); err != nil {
						log.Printf("Failed to add stream link: %v", err)
						c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add stream link: %v", err)})
						return
					}
				} else {
					// For external streams
					streamLink := &models.StreamLink{
						ContentID:     content.ID,
						Name:          sl.Name,
						Quality:       sl.Quality,
						Type:          sl.Type,
						Server:        "external",
						URL:           sl.URL,
						EpisodeNumber: ep.EpisodeNumber,
						SeasonNumber:  ep.SeasonNumber,
					}

					// Handle embed type differently
					if sl.Type == "embed" {
						// For embed type, keep the full IFRAME tag
						streamLink.Type = "embed"
						// If URL doesn't contain IFRAME tag but is from known embed providers, wrap it in IFRAME
						if !strings.Contains(strings.ToUpper(sl.URL), "IFRAME") {
							// Add more providers as needed
							if strings.Contains(sl.URL, "mp4upload.com") {
								streamLink.URL = fmt.Sprintf(`<IFRAME SRC="%s" FRAMEBORDER=0 MARGINWIDTH=0 MARGINHEIGHT=0 SCROLLING=NO WIDTH=1280 HEIGHT=720 allowfullscreen></IFRAME>`, sl.URL)
							} else {
								streamLink.URL = sl.URL
							}
						}
					}

					if err := h.contentService.AddStreamLink(content.ID, streamLink); err != nil {
						log.Printf("Failed to add stream link: %v", err)
						c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add stream link: %v", err)})
						return
					}
				}
			}

			// Process download links
			for _, dl := range ep.DownloadLinks {
				downloadLink := &models.DownloadLink{
					ContentID:     content.ID,
					Name:          dl.Name,
					Quality:       dl.Quality,
					URL:           dl.URL,
					Server:        dl.Name,
					EpisodeNumber: ep.EpisodeNumber,
					SeasonNumber:  ep.SeasonNumber,
				}
				if err := h.contentService.AddDownloadLink(content.ID, downloadLink); err != nil {
					log.Printf("Failed to add download link: %v", err)
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Content %s successfully", map[bool]string{true: "updated", false: "created"}[isUpdate]),
		"content": content,
	})
}

// Get handles getting a single content
func (h *ContentHandler) Get(c *gin.Context) {
	idStr := c.Param("contentId")
	log.Printf("Getting content with ID: %s", idStr)

	// Tambahkan validasi untuk ID kosong
	if idStr == "" {
		log.Printf("Content ID is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content ID is required"})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("Invalid content ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID format"})
		return
	}

	// Get content with preloaded relationships including genres
	content, err := h.contentService.GetContentByID(uint(id))
	if err != nil {
		log.Printf("Content not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Content not found"})
		return
	}

	// Log response untuk debugging
	log.Printf("Content retrieved successfully: ID=%d, Title=%s, Genres=%v",
		content.ID,
		content.Title,
		content.Genres)

	c.JSON(http.StatusOK, content)
}

// List handles listing all content with pagination and filters
func (h *ContentHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "12"))

	// Check for categoryId parameter
	if categoryIDStr := c.Query("categoryId"); categoryIDStr != "" {
		categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
			return
		}
		contents, total, err := h.contentService.GetContentByCategory(uint(categoryID), page, pageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"contents": contents,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		})
		return
	}

	// Build filters for other cases
	filters := make(map[string]interface{})
	if contentType := c.Query("type"); contentType != "" {
		filters["type"] = contentType
	}

	contents, total, err := h.contentService.ListContent(page, pageSize, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"contents": contents,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// Update handles content updates
func (h *ContentHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Get existing content
	existingContent, err := h.contentService.GetContentByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Content not found"})
		return
	}

	var input CreateContentRequest
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update content fields
	existingContent.Title = input.Title
	existingContent.Description = input.Description
	existingContent.Type = string(input.Type)
	existingContent.ReleaseDate = input.ReleaseDate
	existingContent.Rating = input.Rating
	existingContent.SeasonID = input.SeasonID

	// Handle cover image update if provided
	if file, err := c.FormFile("coverImage"); err == nil {
		log.Printf("Updating cover image for contentID: %d with file: %s", existingContent.ID, file.Filename)
		if err := h.mediaService.UploadContentCover(existingContent.ID, file); err != nil {
			log.Printf("Failed to update cover image: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update cover image: %v", err)})
			return
		}
		log.Printf("Successfully updated cover image for content %d", existingContent.ID)

		// Reload content to get updated cover image path
		updatedContent, err := h.contentService.GetContentByID(existingContent.ID)
		if err == nil {
			existingContent.CoverImage = updatedContent.CoverImage
			log.Printf("Updated content with cover image path: %s", existingContent.CoverImage)
		}
	} else {
		log.Printf("No cover image update provided or error getting cover: %v", err)
	}

	// Update content in database
	if err := h.contentService.UpdateContent(existingContent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Process episodes if any
	if input.Episodes != "" {
		var episodes []struct {
			Title         string                `json:"title"`
			Description   string                `json:"description"`
			EpisodeNumber int                   `json:"episodeNumber"`
			SeasonNumber  int                   `json:"seasonNumber"`
			StreamLinks   []StreamLinkRequest   `json:"streamLinks"`
			DownloadLinks []DownloadLinkRequest `json:"downloadLinks"`
			VideoField    string                `json:"videoField"`
		}

		if err := json.Unmarshal([]byte(input.Episodes), &episodes); err != nil {
			log.Printf("Failed to parse episodes: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to parse episodes: %v", err)})
			return
		}

		// Delete existing episodes and their links
		if err := h.contentService.DeleteContentEpisodes(existingContent.ID); err != nil {
			log.Printf("Failed to delete existing episodes: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete existing episodes: %v", err)})
			return
		}

		// Create new episodes
		for _, ep := range episodes {
			episode := &models.Episode{
				ContentID:     existingContent.ID,
				Title:         ep.Title,
				Description:   ep.Description,
				Type:          "episode",
				EpisodeNumber: ep.EpisodeNumber,
				SeasonNumber:  ep.SeasonNumber,
			}

			// Create episode
			if err := h.contentService.GetDB().Create(episode).Error; err != nil {
				log.Printf("Failed to create episode: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create episode: %v", err)})
				return
			}

			// Process stream links
			for _, sl := range ep.StreamLinks {
				// Handle video upload for self-hosted streams
				if sl.Type == "self-hosted" && sl.VideoField != "" {
					file, err := c.FormFile(sl.VideoField)
					if err != nil {
						log.Printf("No video file found for field %s: %v", sl.VideoField, err)
						continue
					}

					log.Printf("Processing video upload for episode %d", ep.EpisodeNumber)
					videoPath, err := h.mediaService.UploadVideo(existingContent.ID, &episode.ID, file)
					if err != nil {
						log.Printf("Failed to upload video: %v", err)
						c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload video: %v", err)})
						return
					}

					// Update episode with video path
					episode.VideoPath = videoPath
					if err := h.contentService.GetDB().Save(episode).Error; err != nil {
						log.Printf("Failed to update episode with video path: %v", err)
						c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update episode: %v", err)})
						return
					}

					// Create stream link with video path
					streamLink := &models.StreamLink{
						ContentID:     existingContent.ID,
						Name:          sl.Name,
						Quality:       sl.Quality,
						Type:          sl.Type,
						Server:        "local",
						URL:           videoPath,
						EpisodeNumber: ep.EpisodeNumber,
						SeasonNumber:  ep.SeasonNumber,
					}

					if err := h.contentService.AddStreamLink(existingContent.ID, streamLink); err != nil {
						log.Printf("Failed to add stream link: %v", err)
						c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add stream link: %v", err)})
						return
					}
				} else {
					// For external streams
					streamLink := &models.StreamLink{
						ContentID:     existingContent.ID,
						Name:          sl.Name,
						Quality:       sl.Quality,
						Type:          sl.Type,
						Server:        "external",
						URL:           sl.URL,
						EpisodeNumber: ep.EpisodeNumber,
						SeasonNumber:  ep.SeasonNumber,
					}

					// Handle embed type differently
					if sl.Type == "embed" {
						// For embed type, keep the full IFRAME tag
						streamLink.Type = "embed"
						// If URL doesn't contain IFRAME tag but is from known embed providers, wrap it in IFRAME
						if !strings.Contains(strings.ToUpper(sl.URL), "IFRAME") {
							// Add more providers as needed
							if strings.Contains(sl.URL, "mp4upload.com") {
								streamLink.URL = fmt.Sprintf(`<IFRAME SRC="%s" FRAMEBORDER=0 MARGINWIDTH=0 MARGINHEIGHT=0 SCROLLING=NO WIDTH=1280 HEIGHT=720 allowfullscreen></IFRAME>`, sl.URL)
							} else {
								streamLink.URL = sl.URL
							}
						}
					}

					if err := h.contentService.AddStreamLink(existingContent.ID, streamLink); err != nil {
						log.Printf("Failed to add stream link: %v", err)
						c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add stream link: %v", err)})
						return
					}
				}
			}

			// Process download links
			for _, dl := range ep.DownloadLinks {
				downloadLink := &models.DownloadLink{
					ContentID:     existingContent.ID,
					Name:          dl.Name,
					Quality:       dl.Quality,
					URL:           dl.URL,
					Server:        "external",
					EpisodeNumber: ep.EpisodeNumber,
					SeasonNumber:  ep.SeasonNumber,
				}

				if err := h.contentService.AddDownloadLink(existingContent.ID, downloadLink); err != nil {
					log.Printf("Failed to add download link: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add download link: %v", err)})
					return
				}
			}
		}
	}

	c.JSON(http.StatusOK, existingContent)
}

// Delete handles content deletion
func (h *ContentHandler) Delete(c *gin.Context) {
	// Get ID from URL parameter
	idStr := c.Param("contentId") // Ubah dari "id" menjadi "contentId" sesuai dengan route parameter
	if idStr == "" {
		log.Printf("Content ID is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content ID is required"})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("Invalid content ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID format"})
		return
	}

	log.Printf("Attempting to delete content with ID: %d", id)

	// Verify content exists first
	content, err := h.contentService.GetContentByID(uint(id))
	if err != nil {
		log.Printf("Content not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Content not found"})
		return
	}

	// Delete content
	if err := h.contentService.DeleteContent(content.ID); err != nil {
		log.Printf("Failed to delete content: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete content"})
		return
	}

	log.Printf("Content %d deleted successfully", id)
	c.JSON(http.StatusOK, gin.H{
		"message": "Content deleted successfully",
		"id":      id,
	})
}

// Search handles content search
func (h *ContentHandler) Search(c *gin.Context) {
	term := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "12"))

	contents, total, err := h.contentService.SearchContent(term, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"contents": contents,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetByGenre handles getting content by genre
func (h *ContentHandler) GetByGenre(c *gin.Context) {
	genreID, err := strconv.ParseUint(c.Param("genreId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "12"))

	contents, total, err := h.contentService.GetContentByGenre(uint(genreID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"contents": contents,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetByCategory handles getting content by category
func (h *ContentHandler) GetByCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("categoryId"), 10, 32)
	log.Printf("Checking Log GetByCategory: %s", categoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "12"))

	contents, total, err := h.contentService.GetContentByCategory(uint(categoryID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"contents": contents,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func init() {
	// Buat direktori media jika belum ada
	mediaDirs := []string{
		"media/thumbnails/content",
		"media/thumbnails/episodes",
		"media/videos/original",
		"media/videos/transcoded",
	}

	for _, dir := range mediaDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Failed to create directory %s: %v", dir, err)
		}
	}
}
