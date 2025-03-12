package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/username/anime-streaming/internal/models"
	"github.com/username/anime-streaming/internal/services"
)

// EpisodeHandler handles episode related requests
type EpisodeHandler struct {
	episodeService *services.EpisodeService
}

// NewEpisodeHandler creates a new EpisodeHandler
func NewEpisodeHandler(episodeService *services.EpisodeService) *EpisodeHandler {
	return &EpisodeHandler{
		episodeService: episodeService,
	}
}

// Create handles episode creation
func (h *EpisodeHandler) Create(c *gin.Context) {
	var episode models.Episode
	if err := c.ShouldBindJSON(&episode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.episodeService.CreateEpisode(&episode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, episode)
}

// Get handles getting a single episode
func (h *EpisodeHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	episode, err := h.episodeService.GetEpisodeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Episode not found"})
		return
	}

	c.JSON(http.StatusOK, episode)
}

// List handles listing episodes for a content
func (h *EpisodeHandler) List(c *gin.Context) {
	contentIdStr := c.Param("contentId")
	log.Printf("Listing episodes for content ID: %s", contentIdStr)

	contentID, err := strconv.ParseUint(contentIdStr, 10, 32)
	if err != nil {
		log.Printf("Invalid content ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID format"})
		return
	}

	// Verify content exists first
	content, err := h.episodeService.GetContent(uint(contentID))
	if err != nil {
		log.Printf("Content not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Content not found"})
		return
	}

	var season *int
	if seasonStr := c.Query("season"); seasonStr != "" {
		seasonInt, err := strconv.Atoi(seasonStr)
		if err != nil {
			log.Printf("Invalid season number: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season number"})
			return
		}
		season = &seasonInt
	}

	episodes, err := h.episodeService.ListEpisodes(uint(contentID), season)
	if err != nil {
		log.Printf("Failed to list episodes: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Found %d episodes for content ID %d", len(episodes), contentID)
	c.JSON(http.StatusOK, gin.H{
		"episodes": episodes,
		"content":  content,
	})
}

// Update handles episode updates
func (h *EpisodeHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	episode, err := h.episodeService.GetEpisodeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Episode not found"})
		return
	}

	if err := c.ShouldBindJSON(episode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.episodeService.UpdateEpisode(episode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, episode)
}

// Delete handles episode deletion
func (h *EpisodeHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.episodeService.DeleteEpisode(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Episode deleted successfully"})
}

// GetNext handles getting the next episode
func (h *EpisodeHandler) GetNext(c *gin.Context) {
	contentID, err := strconv.ParseUint(c.Param("contentId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}

	currentSeason, err := strconv.Atoi(c.Query("season"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season number"})
		return
	}

	currentEpisode, err := strconv.Atoi(c.Query("episode"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid episode number"})
		return
	}

	episode, err := h.episodeService.GetNextEpisode(uint(contentID), currentSeason, currentEpisode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No next episode found"})
		return
	}

	c.JSON(http.StatusOK, episode)
}

// GetLatest handles getting the latest episode
func (h *EpisodeHandler) GetLatest(c *gin.Context) {
	contentID, err := strconv.ParseUint(c.Param("contentId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}

	episode, err := h.episodeService.GetLatestEpisode(uint(contentID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No episodes found"})
		return
	}

	c.JSON(http.StatusOK, episode)
}
