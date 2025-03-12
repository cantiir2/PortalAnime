package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/username/anime-streaming/internal/models"
	"github.com/username/anime-streaming/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles business logic for users
type UserService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

// NewUserService creates a new UserService
func NewUserService(userRepo *repository.UserRepository, jwtSecret string) *UserService {
	return &UserService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Register creates a new user account
func (s *UserService) Register(username, email, password string) error {
	// Check if this is the first user (make them admin)
	count, err := s.userRepo.CountUsers()
	if err != nil {
		return err
	}

	// Add debug logging
	log.Printf("Current user count: %d", count)

	// Check if username exists
	if _, err := s.userRepo.FindByUsername(username); err == nil {
		return errors.New("username already exists")
	}

	// Check if email exists
	if _, err := s.userRepo.FindByEmail(email); err == nil {
		return errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create user
	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Role:     models.RoleUser, // Default role
	}

	// If this is the first user, make them admin
	if count == 0 {
		log.Printf("Creating first user as admin: %s", email)
		user.Role = models.RoleAdmin
	}

	// Add debug logging
	log.Printf("Creating user with role: %s", user.Role)

	return s.userRepo.Create(user)
}

// Login authenticates a user and returns a JWT token
func (s *UserService) Login(email, password string) (string, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	if err := s.userRepo.Update(user); err != nil {
		return "", err
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 hour expiry
	})

	return token.SignedString([]byte(s.jwtSecret))
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(user *models.User) error {
	return s.userRepo.Update(user)
}

// DeleteUser deletes a user account
func (s *UserService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}

// ListUsers lists all users with pagination
func (s *UserService) ListUsers(page, pageSize int) ([]models.User, int64, error) {
	return s.userRepo.List(page, pageSize)
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.userRepo.Update(user)
}

// ValidateToken validates a JWT token and returns the user ID and role
func (s *UserService) ValidateToken(tokenString string) (uint, models.Role, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return 0, "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["user_id"].(float64))
		roleStr := claims["role"].(string)
		role := models.Role(roleStr)

		// Add debug logging
		log.Printf("Token validation - UserID: %v, Role: %v", userID, role)

		// Validate role
		if role != models.RoleAdmin && role != models.RoleUser {
			return 0, "", fmt.Errorf("invalid role in token")
		}

		return userID, role, nil
	}

	return 0, "", fmt.Errorf("invalid token")
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.FindByEmail(email)
}
