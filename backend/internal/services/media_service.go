package services

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/username/anime-streaming/internal/models"
	"github.com/username/anime-streaming/internal/repository"
)

// MediaService handles business logic for media files
type MediaService struct {
	contentRepo       *repository.ContentRepository
	episodeRepo       *repository.EpisodeRepository
	mediaPath         string
	contentTypeHelper *models.ContentTypeHelper
}

// VideoQuality represents video quality options
type VideoQuality struct {
	Name       string
	Resolution string
	Bitrate    string
}

var (
	// Available video qualities
	VideoQualities = []VideoQuality{
		{Name: "240p", Resolution: "426x240", Bitrate: "400k"},
		{Name: "360p", Resolution: "640x360", Bitrate: "800k"},
		{Name: "480p", Resolution: "854x480", Bitrate: "1500k"},
		{Name: "720p", Resolution: "1280x720", Bitrate: "2500k"},
		{Name: "1080p", Resolution: "1920x1080", Bitrate: "4000k"},
	}
)

// NewMediaService creates a new MediaService
func NewMediaService(
	contentRepo *repository.ContentRepository,
	episodeRepo *repository.EpisodeRepository,
	mediaPath string,
) *MediaService {
	return &MediaService{
		contentRepo:       contentRepo,
		episodeRepo:       episodeRepo,
		mediaPath:         mediaPath,
		contentTypeHelper: models.NewContentTypeHelper(),
	}
}

// UploadContentCover handles cover image upload for content
func (s *MediaService) UploadContentCover(contentID uint, file *multipart.FileHeader) error {
	log.Printf("Uploading cover image for content %d", contentID)

	// Ensure directory exists
	coverDir := filepath.Join(s.mediaPath, "thumbnails", "content")
	if err := os.MkdirAll(coverDir, 0755); err != nil {
		return fmt.Errorf("failed to create cover directory: %v", err)
	}

	// Generate filename with content ID and original extension
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_cover%s", contentID, ext)
	dst := filepath.Join(coverDir, filename)

	log.Printf("Saving cover image to: %s", dst)

	// Save the file
	if err := s.saveUploadedFile(file, dst); err != nil {
		return fmt.Errorf("failed to save cover image: %v", err)
	}

	// Update content record with cover image path
	relativePath := filepath.Join("media", "thumbnails", "content", filename)
	// Convert to Windows path format if needed
	relativePath = strings.ReplaceAll(relativePath, "/", "\\")

	log.Printf("Updating content record with cover image path: %s", relativePath)

	// Update the content record in the database
	content, err := s.contentRepo.FindByID(contentID)
	if err != nil {
		return fmt.Errorf("failed to find content: %v", err)
	}

	content.CoverImage = relativePath
	if err := s.contentRepo.Update(content); err != nil {
		return fmt.Errorf("failed to update content with cover image: %v", err)
	}

	log.Printf("Successfully updated content %d with cover image path", contentID)
	return nil
}

func (s *MediaService) GetMediaPath() string {
	return s.mediaPath
}

// UploadEpisodeThumbnail uploads a thumbnail for an episode
func (s *MediaService) UploadEpisodeThumbnail(episodeID uint, file *multipart.FileHeader) error {
	// Verify episode exists
	episode, err := s.episodeRepo.FindByID(episodeID)
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	thumbnailDir := filepath.Join(s.mediaPath, "thumbnails", "episodes")
	if err := os.MkdirAll(thumbnailDir, 0755); err != nil {
		return err
	}

	// Generate filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_thumbnail%s", episodeID, ext)
	filepath := filepath.Join(thumbnailDir, filename)

	// Save file
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// Update episode thumbnail path
	episode.ThumbnailURL = filename
	return s.episodeRepo.Update(episode)
}

// UploadVideo handles video file upload and processing
func (s *MediaService) UploadVideo(contentID uint, episodeID *uint, file *multipart.FileHeader) (string, error) {
	log.Printf("Starting video upload for contentID: %d, episodeID: %v", contentID, episodeID)

	// Generate unique filename
	filename := fmt.Sprintf("%d_%d.mp4", contentID, time.Now().Unix())

	// Set upload directories
	originalDir := filepath.Join(s.mediaPath, "videos", "original")
	if err := os.MkdirAll(originalDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create original video directory: %v", err)
	}

	// Full path for original video
	originalPath := filepath.Join(originalDir, filename)

	// Save the original file
	if err := s.saveUploadedFile(file, originalPath); err != nil {
		return "", fmt.Errorf("failed to save original video: %v", err)
	}
	log.Printf("Original video saved successfully at: %s", originalPath)

	// Start transcoding in background
	go func() {
		if err := s.transcodeVideo(filename); err != nil {
			log.Printf("Failed to transcode video %s: %v", filename, err)
			return
		}
		log.Printf("Video transcoding completed for %s", filename)
	}()

	// Return relative path for database storage
	relativePath := filepath.Join("videos", "original", filename)
	log.Printf("Returning relative path: %s", relativePath)

	return relativePath, nil
}

// saveUploadedFile saves the uploaded file to disk
func (s *MediaService) saveUploadedFile(file *multipart.FileHeader, dst string) error {
	log.Printf("Saving uploaded file: %s to %s", file.Filename, dst)

	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// Ensure the destination directory exists
	dir := filepath.Dir(dst)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %v", err)
	}

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return fmt.Errorf("failed to copy file: %v", err)
	}

	log.Printf("Successfully saved file to %s", dst)
	return nil
}

// transcodeVideo transcodes the original video into different qualities
func (s *MediaService) transcodeVideo(filename string) error {
	log.Printf("Starting transcoding for video: %s", filename)
	originalPath := filepath.Join(s.mediaPath, "videos", "original", filename)

	for _, quality := range VideoQualities {
		outputDir := filepath.Join(s.mediaPath, "videos", "transcoded", quality.Name)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory for %s: %v", quality.Name, err)
		}

		outputPath := filepath.Join(outputDir, filename)
		if strings.HasSuffix(outputPath, ".mp4") {
			outputPath = strings.TrimSuffix(outputPath, ".mp4") + ".m3u8"
		}

		log.Printf("Transcoding to %s quality. Output: %s", quality.Name, outputPath)

		cmd := exec.Command(
			"ffmpeg",
			"-i", originalPath,
			"-c:v", "libx264",
			"-c:a", "aac",
			"-b:v", quality.Bitrate,
			"-s", quality.Resolution,
			"-f", "hls",
			"-hls_time", "10",
			"-hls_list_size", "0",
			"-hls_segment_filename", filepath.Join(outputDir, fmt.Sprintf("%s_%%03d.ts", strings.TrimSuffix(filename, ".mp4"))),
			outputPath,
		)

		stderr, err := cmd.StderrPipe()
		if err != nil {
			return fmt.Errorf("failed to get stderr pipe: %v", err)
		}

		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start transcoding: %v", err)
		}

		// Read and log stderr in real-time
		go func() {
			buf := make([]byte, 1024)
			for {
				n, err := stderr.Read(buf)
				if n > 0 {
					log.Printf("FFmpeg: %s", string(buf[:n]))
				}
				if err != nil {
					break
				}
			}
		}()

		if err := cmd.Wait(); err != nil {
			return fmt.Errorf("transcoding failed: %v", err)
		}

		log.Printf("Successfully transcoded to %s quality", quality.Name)
	}

	log.Printf("Transcoding completed for all qualities: %s", filename)
	return nil
}

// GetVideoPath gets the video file path for streaming
func (s *MediaService) GetVideoPath(contentID uint, episodeID *uint, quality string) (string, error) {
	log.Printf("Getting video path for content %d, episode %v, quality %s", contentID, episodeID, quality)

	if episodeID != nil {
		// Get episode
		episode, err := s.episodeRepo.FindByID(*episodeID)
		if err != nil {
			return "", fmt.Errorf("failed to find episode: %v", err)
		}

		// Verify episode belongs to content
		if episode.ContentID != contentID {
			return "", fmt.Errorf("episode does not belong to content")
		}

		log.Printf("Found episode: ID=%d, ContentID=%d, VideoPath=%s", episode.ID, episode.ContentID, episode.VideoPath)

		// If quality is not original, check for transcoded version
		if quality != "original" {
			transcodedPath := filepath.Join(s.mediaPath, "videos", quality, fmt.Sprintf("%d_%d.mp4", contentID, episode.ID))
			if _, err := os.Stat(transcodedPath); err == nil {
				log.Printf("Found transcoded video: %s", transcodedPath)
				return transcodedPath, nil
			}
			log.Printf("Transcoded version not found, falling back to original")
		}

		// Return original video path
		if episode.VideoPath != "" {
			// Remove any "media\" prefix if it exists
			videoPath := strings.TrimPrefix(episode.VideoPath, "media\\")
			fullPath := filepath.Join(s.mediaPath, videoPath)
			log.Printf("Original video path: %s", fullPath)
			return fullPath, nil
		}
	}

	// For movies or if episode path is empty
	fullPath := filepath.Join(s.mediaPath, "videos", "original", fmt.Sprintf("%d.mp4", contentID))
	log.Printf("Returning original video path: %s", fullPath)
	return fullPath, nil
}
