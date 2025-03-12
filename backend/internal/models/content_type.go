package models

// ContentTypeHelper provides helper functions for content types
type ContentTypeHelper struct{}

// IsMovie checks if the given type is a movie
func (h *ContentTypeHelper) IsMovie(contentType string) bool {
	return contentType == "Movie"
}

// IsSeries checks if the given type is a series or anime
func (h *ContentTypeHelper) IsSeries(contentType string) bool {
	return contentType == "Series" || contentType == "Anime"
}

// IsValidType checks if the given type is valid
func (h *ContentTypeHelper) IsValidType(contentType string, categories []Category) bool {
	for _, category := range categories {
		if contentType == category.Name {
			return true
		}
	}
	return false
}

// NewContentTypeHelper creates a new ContentTypeHelper
func NewContentTypeHelper() *ContentTypeHelper {
	return &ContentTypeHelper{}
}
