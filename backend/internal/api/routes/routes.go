package routes

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/username/anime-streaming/internal/api/handlers"
	"github.com/username/anime-streaming/internal/api/middleware"
	"github.com/username/anime-streaming/internal/config"
	"github.com/username/anime-streaming/internal/models"
	"github.com/username/anime-streaming/internal/repository"
	"github.com/username/anime-streaming/internal/services"
	"gorm.io/gorm"
)

// SetupRouter sets up the routing for the application
func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(cfg.CorsAllowedOrigins, ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Custom handler for serving media files with logging
	router.GET("/media/*filepath", func(c *gin.Context) {
		filepath := c.Param("filepath")
		fullPath := cfg.MediaPath + "/" + filepath
		log.Printf("Accessing media file: %s", filepath)
		log.Printf("Full path: %s", fullPath)

		// Check if file exists
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			log.Printf("File not found: %s", fullPath)
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}

		c.File(fullPath)
	})

	// Serve favicon
	router.StaticFile("/favicon.ico", "./assets/favicon.ico")

	// Set maximum multipart memory
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	// Tambahkan middleware untuk meningkatkan batas ukuran body
	router.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 32<<20) // 32MB
		c.Next()
	})

	// Setup CORS
	// router.Use(func(c *gin.Context) {
	// 	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	// 	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	// 	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	// 	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
	// 	c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
	// 	c.Writer.Header().Set("Access-Control-Max-Age", "86400")

	// 	if c.Request.Method == "OPTIONS" {
	// 		c.AbortWithStatus(http.StatusNoContent)
	// 		return
	// 	}

	// 	c.Next()
	// })

	// Disable automatic redirection
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	contentRepo := repository.NewContentRepository(db)
	episodeRepo := repository.NewEpisodeRepository(db)
	genreRepo := repository.NewGenreRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	watchHistoryRepo := repository.NewWatchHistoryRepository(db)
	seasonRepo := repository.NewSeasonRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo, cfg.JWTSecret)
	contentService := services.NewContentService(contentRepo, genreRepo, categoryRepo, cfg.MediaPath)
	episodeService := services.NewEpisodeService(episodeRepo, contentRepo, cfg.MediaPath)
	watchHistoryService := services.NewWatchHistoryService(watchHistoryRepo, contentRepo, episodeRepo)
	mediaService := services.NewMediaService(contentRepo, episodeRepo, cfg.MediaPath)
	seasonService := services.NewSeasonService(seasonRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userService)
	contentHandler := handlers.NewContentHandler(contentService, mediaService)
	episodeHandler := handlers.NewEpisodeHandler(episodeService)
	watchHistoryHandler := handlers.NewWatchHistoryHandler(watchHistoryService)
	mediaHandler := handlers.NewMediaHandler(mediaService)
	seasonHandler := handlers.NewSeasonHandler(seasonService)

	// Auth middleware
	authMiddleware := middleware.AuthMiddleware(userService)
	adminMiddleware := middleware.AdminMiddleware()

	api := router.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Content routes
		contents := api.Group("/contents")
		{
			// Public content routes (no parameters)
			contents.GET("", contentHandler.List)
			contents.GET("/search", contentHandler.Search)
			contents.GET("/genre/:genreId", contentHandler.GetByGenre)
			contents.GET("/category/:categoryId", contentHandler.GetByCategory)

			// Protected content routes (no parameters)
			protectedContents := contents.Use(authMiddleware)
			{
				protectedContents.POST("/create", middleware.AuthMiddleware(userService), middleware.AdminMiddleware(), contentHandler.Create)
			}

			// Content detail routes (with contentId)
			contentDetail := contents.Group("/:contentId")
			{
				// Get single content
				contentDetail.GET("", contentHandler.Get)

				// Protected content detail routes
				protectedDetail := contentDetail.Use(authMiddleware)
				{
					protectedDetail.PUT("", adminMiddleware, contentHandler.Update)
					protectedDetail.DELETE("", adminMiddleware, contentHandler.Delete)
					protectedDetail.POST("/upload-video", adminMiddleware, mediaHandler.UploadVideo)
				}

				// Episodes routes
				episodes := contentDetail.Group("/episodes")
				{
					episodes.GET("", episodeHandler.List)
					episodes.GET("/next", episodeHandler.GetNext)
					episodes.GET("/latest", episodeHandler.GetLatest)
					episodes.GET("/:episodeId", episodeHandler.Get)

					// Protected episode routes
					protectedEpisodes := episodes.Use(authMiddleware)
					{
						protectedEpisodes.POST("", adminMiddleware, episodeHandler.Create)
						protectedEpisodes.PUT("/:episodeId", adminMiddleware, episodeHandler.Update)
						protectedEpisodes.DELETE("/:episodeId", adminMiddleware, episodeHandler.Delete)
					}
				}
			}
		}

		// Watch history routes
		history := api.Group("/watch-history", authMiddleware)
		{
			history.POST("", watchHistoryHandler.UpdateProgress)
			history.GET("/content/:contentId", watchHistoryHandler.GetProgress)
			history.GET("/user", watchHistoryHandler.GetUserHistory)
			history.GET("/continue-watching", watchHistoryHandler.GetContinueWatching)
			history.DELETE("/:id", watchHistoryHandler.Delete)
		}

		// Media routes
		media := api.Group("/media")
		{
			// Public media routes
			media.GET("/stream/:contentId", mediaHandler.StreamVideo)
			media.GET("/stream/:contentId/episodes/:episodeId", mediaHandler.StreamVideo)

			// Protected media routes
			protectedMedia := media.Use(authMiddleware, adminMiddleware)
			{
				protectedMedia.POST("/content/:contentId/cover", mediaHandler.UploadContentCover)
				protectedMedia.POST("/episode/:episodeId/thumbnail", mediaHandler.UploadEpisodeThumbnail)
				protectedMedia.POST("/content/:contentId/video", mediaHandler.UploadVideo)
				protectedMedia.POST("/content/:contentId/episodes/:episodeId/video", mediaHandler.UploadVideo)
			}
		}

		// User routes
		users := api.Group("/users", authMiddleware)
		{
			users.PUT("/profile", authHandler.UpdateProfile)
			users.POST("/change-password", authHandler.ChangePassword)
		}

		// Genre routes
		genreService := services.NewGenreService(db)
		genreHandler := handlers.NewGenreHandler(genreService)

		genreRoutes := api.Group("/genres")
		{
			genreRoutes.GET("", genreHandler.List)
			genreRoutes.GET("/:id", genreHandler.Get)
			genreRoutes.POST("", middleware.AuthMiddleware(userService), genreHandler.Create)
			genreRoutes.PUT("/:id", middleware.AuthMiddleware(userService), genreHandler.Update)
			genreRoutes.DELETE("/:id", middleware.AuthMiddleware(userService), genreHandler.Delete)
		}

		// Category routes
		categories := api.Group("/categories")
		{
			categories.GET("", func(c *gin.Context) {
				categories, err := categoryRepo.List()
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, categories)
			})

			// Protected category routes
			protectedCategories := categories.Use(authMiddleware, adminMiddleware)
			{
				protectedCategories.POST("", func(c *gin.Context) {
					// Add debug logging
					roleInterface, _ := c.Get("userRole")
					userID, _ := c.Get("userID")
					log.Printf("User ID: %v, Role: %v attempting to create category", userID, roleInterface)

					var category models.Category
					if err := c.ShouldBindJSON(&category); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}
					if err := categoryRepo.Create(&category); err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					c.JSON(http.StatusCreated, category)
				})

				protectedCategories.PUT("/:id", func(c *gin.Context) {
					// Konversi ID dari string ke uint
					id, err := strconv.ParseUint(c.Param("id"), 10, 32)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
						return
					}

					var category models.Category
					if err := c.ShouldBindJSON(&category); err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}

					// Set ID dari parameter ke model
					category.ID = uint(id)

					if err := categoryRepo.Update(&category); err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					c.JSON(http.StatusOK, category)
				})

				protectedCategories.DELETE("/:id", func(c *gin.Context) {
					// Konversi ID dari string ke uint
					id, err := strconv.ParseUint(c.Param("id"), 10, 32)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
						return
					}

					if err := categoryRepo.Delete(uint(id)); err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
				})
			}
		}

		// Season routes
		seasons := api.Group("/seasons")
		{
			seasons.GET("", seasonHandler.List)
			seasons.GET("/current", seasonHandler.GetCurrent)
			seasons.GET("/:id", seasonHandler.Get)

			// Protected season routes
			protectedSeasons := seasons.Use(authMiddleware, adminMiddleware)
			{
				protectedSeasons.POST("", seasonHandler.Create)
				protectedSeasons.PUT("/:id", seasonHandler.Update)
				protectedSeasons.DELETE("/:id", seasonHandler.Delete)
			}
		}

		// Admin routes
		admin := api.Group("/admin").Use(authMiddleware, adminMiddleware)
		{
			admin.GET("/check", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			})
			admin.GET("/verify", func(c *gin.Context) {
				roleInterface, _ := c.Get("userRole")
				userID, _ := c.Get("userID")

				// Add debug logging
				log.Printf("Verifying admin access - UserID: %v, Role: %v", userID, roleInterface)

				c.JSON(http.StatusOK, gin.H{
					"user_id": userID,
					"role":    roleInterface,
					"message": "Admin access verified",
				})
			})

			// Add route to change user role
			admin.PUT("/users/:id/role", func(c *gin.Context) {
				id, err := strconv.ParseUint(c.Param("id"), 10, 32)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
					return
				}

				var input struct {
					Role models.Role `json:"role" binding:"required"`
				}

				if err := c.ShouldBindJSON(&input); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				// Validate role
				if input.Role != models.RoleAdmin && input.Role != models.RoleUser {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
					return
				}

				user, err := userService.GetUserByID(uint(id))
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
					return
				}

				user.Role = input.Role
				if err := userService.UpdateUser(user); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, user)
			})
		}
	}

	return router
}
