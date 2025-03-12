package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/username/anime-streaming/internal/models"
)

// MockEpisodeService is a mock for EpisodeService
type MockEpisodeService struct {
	mock.Mock
}

func (m *MockEpisodeService) ListEpisodes(contentID uint, season *int) ([]models.Episode, error) {
	args := m.Called(contentID, season)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Episode), args.Error(1)
}

func (m *MockEpisodeService) GetContent(contentID uint) (*models.Content, error) {
	args := m.Called(contentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Content), args.Error(1)
}

func TestEpisodeHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		contentID      string
		season         *int
		setupMock      func(*MockEpisodeService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "Valid request without season",
			contentID: "33",
			season:    nil,
			setupMock: func(m *MockEpisodeService) {
				content := &models.Content{
					ID:    33,
					Title: "Test Content",
					Type:  "Anime",
				}
				episodes := []models.Episode{
					{
						ID:            1,
						ContentID:     33,
						Title:         "Episode 1",
						EpisodeNumber: 1,
						SeasonNumber:  1,
					},
				}
				m.On("GetContent", uint(33)).Return(content, nil)
				m.On("ListEpisodes", uint(33), (*int)(nil)).Return(episodes, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"content":{"id":33,"title":"Test Content","type":"Anime"},"episodes":[{"id":1,"content_id":33,"title":"Episode 1","episode_number":1,"season_number":1}]}`,
		},
		{
			name:           "Invalid content ID",
			contentID:      "invalid",
			season:         nil,
			setupMock:      func(m *MockEpisodeService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid content ID format"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockEpisodeService)
			tt.setupMock(mockService)

			handler := NewEpisodeHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set up request parameters
			c.Params = gin.Params{
				{Key: "contentId", Value: tt.contentID},
			}
			if tt.season != nil {
				c.Request = httptest.NewRequest("GET", "/?season="+string(*tt.season), nil)
			} else {
				c.Request = httptest.NewRequest("GET", "/", nil)
			}

			// Call the handler
			handler.List(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}

			mockService.AssertExpectations(t)
		})
	}
}
