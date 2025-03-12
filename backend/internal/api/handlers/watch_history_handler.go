package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/username/anime-streaming/internal/services"
)

// WatchHistoryHandler handles watch history related requests
type WatchHistoryHandler struct {
	watchHistoryService *services.WatchHistoryService
}

// NewWatchHistoryHandler creates a new WatchHistoryHandler
func NewWatchHistoryHandler(watchHistoryService *services.WatchHistoryService) *WatchHistoryHandler {
	return &WatchHistoryHandler{
		watchHistoryService: watchHistoryService,
	}
}

// UpdateProgress handles updating watch progress
func (h *WatchHistoryHandler) UpdateProgress(c *gin.Context) {
	userID, _ := c.Get("userID")
	var input struct {
		ContentID uint  `json:"contentId" binding:"required"`
		EpisodeID *uint `json:"episodeId"`
		Progress  int   `json:"progress" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.watchHistoryService.UpdateProgress(userID.(uint), input.ContentID, input.EpisodeID, input.Progress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Progress updated successfully"})
}

// GetProgress handles getting watch progress
func (h *WatchHistoryHandler) GetProgress(c *gin.Context) {
	userID, _ := c.Get("userID")
	contentID, err := strconv.ParseUint(c.Param("contentId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}

	var episodeID *uint
	if episodeIDStr := c.Query("episodeId"); episodeIDStr != "" {
		episodeIDUint, err := strconv.ParseUint(episodeIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid episode ID"})
			return
		}
		episodeIDUint32 := uint(episodeIDUint)
		episodeID = &episodeIDUint32
	}

	progress, err := h.watchHistoryService.GetProgress(userID.(uint), uint(contentID), episodeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Progress not found"})
		return
	}

	c.JSON(http.StatusOK, progress)
}

// GetUserHistory handles getting user's watch history
func (h *WatchHistoryHandler) GetUserHistory(c *gin.Context) {
	userID, _ := c.Get("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "12"))

	history, total, err := h.watchHistoryService.GetUserHistory(userID.(uint), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"history":  history,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetContinueWatching handles getting content that a user has started but not completed
func (h *WatchHistoryHandler) GetContinueWatching(c *gin.Context) {
	userID, _ := c.Get("userID")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	history, err := h.watchHistoryService.GetContinueWatching(userID.(uint), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

// Delete handles deleting a watch history record
func (h *WatchHistoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.watchHistoryService.DeleteHistory(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Watch history deleted successfully"})
}
