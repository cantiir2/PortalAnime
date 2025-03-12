package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/username/anime-streaming/internal/models"
	"github.com/username/anime-streaming/internal/services"
)

// SeasonHandler handles season related requests
type SeasonHandler struct {
	seasonService *services.SeasonService
}

// NewSeasonHandler creates a new SeasonHandler
func NewSeasonHandler(seasonService *services.SeasonService) *SeasonHandler {
	return &SeasonHandler{
		seasonService: seasonService,
	}
}

// Create handles season creation
func (h *SeasonHandler) Create(c *gin.Context) {
	var season models.Season
	if err := c.ShouldBindJSON(&season); err != nil {
		log.Printf("Failed to bind season data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.seasonService.CreateSeason(&season); err != nil {
		log.Printf("Failed to create season: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Season created successfully: %s %d", season.Name, season.Year)
	c.JSON(http.StatusCreated, season)
}

// List handles listing all seasons
func (h *SeasonHandler) List(c *gin.Context) {
	seasons, err := h.seasonService.ListSeasons()
	if err != nil {
		log.Printf("Failed to list seasons: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Retrieved %d seasons", len(seasons))
	c.JSON(http.StatusOK, seasons)
}

// Get handles getting a single season
func (h *SeasonHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Printf("Invalid season ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season ID"})
		return
	}

	season, err := h.seasonService.GetSeasonByID(uint(id))
	if err != nil {
		log.Printf("Season not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Season not found"})
		return
	}

	log.Printf("Retrieved season: %s %d", season.Name, season.Year)
	c.JSON(http.StatusOK, season)
}

// Update handles season updates
func (h *SeasonHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Printf("Invalid season ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season ID"})
		return
	}

	var season models.Season
	if err := c.ShouldBindJSON(&season); err != nil {
		log.Printf("Failed to bind season data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	season.ID = uint(id)
	if err := h.seasonService.UpdateSeason(&season); err != nil {
		log.Printf("Failed to update season: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Season updated successfully: %s %d", season.Name, season.Year)
	c.JSON(http.StatusOK, season)
}

// Delete handles season deletion
func (h *SeasonHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Printf("Invalid season ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season ID"})
		return
	}

	if err := h.seasonService.DeleteSeason(uint(id)); err != nil {
		log.Printf("Failed to delete season: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Season deleted successfully: %d", id)
	c.JSON(http.StatusOK, gin.H{"message": "Season deleted successfully"})
}

// GetCurrent handles getting the current active season
func (h *SeasonHandler) GetCurrent(c *gin.Context) {
	season, err := h.seasonService.GetCurrentSeason()
	if err != nil {
		log.Printf("Failed to get current season: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "No active season found"})
		return
	}

	log.Printf("Retrieved current season: %s %d", season.Name, season.Year)
	c.JSON(http.StatusOK, season)
}
