package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/username/anime-streaming/internal/models"
	"github.com/username/anime-streaming/internal/services"
)

// GenreHandler handles genre related requests
type GenreHandler struct {
	genreService *services.GenreService
}

// NewGenreHandler creates a new GenreHandler
func NewGenreHandler(genreService *services.GenreService) *GenreHandler {
	return &GenreHandler{
		genreService: genreService,
	}
}

// Create handles genre creation
func (h *GenreHandler) Create(c *gin.Context) {
	var genre models.Genre
	if err := c.ShouldBindJSON(&genre); err != nil {
		log.Printf("Failed to bind genre data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.genreService.CreateGenre(&genre); err != nil {
		log.Printf("Failed to create genre: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Genre created successfully: %s", genre.Name)
	c.JSON(http.StatusCreated, genre)
}

// List handles listing all genres
func (h *GenreHandler) List(c *gin.Context) {
	genres, err := h.genreService.ListGenres()
	if err != nil {
		log.Printf("Failed to list genres: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Retrieved %d genres", len(genres))
	c.JSON(http.StatusOK, genres)
}

// Get handles getting a single genre
func (h *GenreHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Printf("Invalid genre ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID"})
		return
	}

	genre, err := h.genreService.GetGenreByID(uint(id))
	if err != nil {
		log.Printf("Genre not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Genre not found"})
		return
	}

	log.Printf("Retrieved genre: %s", genre.Name)
	c.JSON(http.StatusOK, genre)
}

// Update handles genre updates
func (h *GenreHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Printf("Invalid genre ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID"})
		return
	}

	var genre models.Genre
	if err := c.ShouldBindJSON(&genre); err != nil {
		log.Printf("Failed to bind genre data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	genre.ID = uint(id)
	if err := h.genreService.UpdateGenre(&genre); err != nil {
		log.Printf("Failed to update genre: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Genre updated successfully: %s", genre.Name)
	c.JSON(http.StatusOK, genre)
}

// Delete handles genre deletion
func (h *GenreHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		log.Printf("Invalid genre ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid genre ID"})
		return
	}

	if err := h.genreService.DeleteGenre(uint(id)); err != nil {
		log.Printf("Failed to delete genre: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Genre deleted successfully: %d", id)
	c.JSON(http.StatusOK, gin.H{"message": "Genre deleted successfully"})
}
