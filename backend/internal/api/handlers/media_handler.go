package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/username/anime-streaming/internal/services"
)

// MediaHandler handles media related requests
type MediaHandler struct {
	mediaService *services.MediaService
}

// NewMediaHandler creates a new MediaHandler
func NewMediaHandler(mediaService *services.MediaService) *MediaHandler {
	return &MediaHandler{
		mediaService: mediaService,
	}
}

// UploadContentCover handles cover image upload for content
func (h *MediaHandler) UploadContentCover(c *gin.Context) {
	contentID, err := strconv.ParseUint(c.Param("contentId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}

	file, err := c.FormFile("cover")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	if err := h.mediaService.UploadContentCover(uint(contentID), file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cover image uploaded successfully"})
}

// UploadEpisodeThumbnail handles thumbnail upload for episodes
func (h *MediaHandler) UploadEpisodeThumbnail(c *gin.Context) {
	episodeID, err := strconv.ParseUint(c.Param("episodeId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid episode ID"})
		return
	}

	file, err := c.FormFile("thumbnail")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	if err := h.mediaService.UploadEpisodeThumbnail(uint(episodeID), file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Thumbnail uploaded successfully"})
}

// UploadVideo handles video file upload
func (h *MediaHandler) UploadVideo(c *gin.Context) {
	contentID, err := strconv.ParseUint(c.Param("contentId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}

	file, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No video file uploaded"})
		return
	}

	// Validasi ukuran file
	if file.Size > 500*1024*1024 { // 500MB limit
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size too large. Maximum size is 500MB"})
		return
	}

	// Log untuk debugging
	log.Printf("Received video upload for content %d: %s (size: %d bytes)",
		contentID, file.Filename, file.Size)

	// Upload video
	videoPath, err := h.mediaService.UploadVideo(uint(contentID), nil, file)
	if err != nil {
		log.Printf("Failed to upload video: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Video uploaded successfully: %s", videoPath)

	c.JSON(http.StatusOK, gin.H{
		"message": "Video uploaded successfully",
		"path":    videoPath,
	})
}

// StreamVideo streams a video file
func (h *MediaHandler) StreamVideo(c *gin.Context) {
	contentID, err := strconv.ParseUint(c.Param("contentId"), 10, 32)
	if err != nil {
		log.Printf("Invalid content ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid content id"})
		return
	}

	var episodeID *uint
	episodeIDStr := c.Param("episodeId")
	if episodeIDStr != "" {
		epID, err := strconv.ParseUint(episodeIDStr, 10, 32)
		if err != nil {
			log.Printf("Invalid episode ID: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid episode id"})
			return
		}
		epIDUint := uint(epID)
		episodeID = &epIDUint
	}

	quality := c.DefaultQuery("quality", "original")
	log.Printf("Streaming request - Content: %d, Episode: %v, Quality: %s", contentID, episodeID, quality)

	// Get video path
	videoPath, err := h.mediaService.GetVideoPath(uint(contentID), episodeID, quality)
	log.Printf("Video path: %s", videoPath)
	if err != nil {
		log.Printf("Failed to get video path: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
		return
	}

	// Check if file exists
	file, err := os.Open(videoPath)
	if err != nil {
		log.Printf("Failed to open video file: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "video file not found"})
		return
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("Failed to get file info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get file info"})
		return
	}

	// Set content type and basic headers
	c.Header("Content-Type", "video/mp4")
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	// Handle range request
	rangeHeader := c.GetHeader("Range")
	if rangeHeader != "" {
		ranges, err := parseRange(rangeHeader, fileInfo.Size())
		if err != nil {
			log.Printf("Failed to parse range header: %v", err)
			c.Status(http.StatusRequestedRangeNotSatisfiable)
			return
		}

		if len(ranges) > 1 {
			log.Printf("Multiple ranges not supported")
			c.Status(http.StatusRequestedRangeNotSatisfiable)
			return
		}

		r := ranges[0]
		if _, err := file.Seek(r.start, io.SeekStart); err != nil {
			log.Printf("Failed to seek file: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusPartialContent)
		c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", r.start, r.end, fileInfo.Size()))
		c.Header("Content-Length", strconv.FormatInt(r.length, 10))

		// Stream the file
		_, err = io.CopyN(c.Writer, file, r.length)
		if err != nil {
			log.Printf("Failed to stream file: %v", err)
			return
		}
	} else {
		// Stream entire file
		log.Printf("Streaming entire file: %s", videoPath)
		http.ServeContent(c.Writer, c.Request, fileInfo.Name(), fileInfo.ModTime(), file)
	}
}

type httpRange struct {
	start, end int64
	length     int64
}

// parseRange parses a Range header string as per RFC 7233
func parseRange(s string, size int64) ([]httpRange, error) {
	if !strings.HasPrefix(s, "bytes=") {
		return nil, fmt.Errorf("invalid range format")
	}

	s = strings.TrimPrefix(s, "bytes=")
	var ranges []httpRange

	for _, ra := range strings.Split(s, ",") {
		ra = strings.TrimSpace(ra)
		if ra == "" {
			continue
		}

		i := strings.Index(ra, "-")
		if i < 0 {
			return nil, fmt.Errorf("invalid range format")
		}

		start, end := strings.TrimSpace(ra[:i]), strings.TrimSpace(ra[i+1:])
		r := httpRange{start: 0, end: size - 1}

		if start == "" {
			// suffix-length
			if end == "" {
				return nil, fmt.Errorf("invalid range format")
			}
			i, err := strconv.ParseInt(end, 10, 64)
			if err != nil || i > size {
				return nil, fmt.Errorf("invalid range format")
			}
			r.start = size - i
		} else {
			i, err := strconv.ParseInt(start, 10, 64)
			if err != nil || i >= size || i < 0 {
				return nil, fmt.Errorf("invalid range format")
			}
			r.start = i
			if end != "" {
				i, err := strconv.ParseInt(end, 10, 64)
				if err != nil || i >= size || i < 0 {
					return nil, fmt.Errorf("invalid range format")
				}
				r.end = i
			}
		}

		if r.start > r.end {
			return nil, fmt.Errorf("invalid range format")
		}

		r.length = r.end - r.start + 1
		ranges = append(ranges, r)
	}

	if len(ranges) == 0 {
		return nil, fmt.Errorf("invalid range format")
	}

	return ranges, nil
}
