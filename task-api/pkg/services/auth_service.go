package services

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"task-management-api/pkg/config"
	"task-management-api/pkg/models"
	"task-management-api/pkg/repositories"
	"task-management-api/pkg/utils"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo *repositories.UserRepository
	config   *config.Config
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo *repositories.UserRepository, config *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   config,
	}
}

// Register creates a new user account
func (s *AuthService) Register(req *models.RegisterRequest) (*models.User, error) {
	// Validate input
	if !utils.IsValidEmail(req.Email) {
		return nil, errors.New("invalid email format")
	}

	if !utils.IsValidPassword(req.Password) {
		return nil, errors.New("password must be at least 6 characters")
	}

	if !utils.ValidateUserName(req.Name) {
		return nil, errors.New("name must be between 2 and 255 characters")
	}

	// Check if user already exists
	exists, err := s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	// Create user
	user := &models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		Role:     "user",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Validate input
	if !utils.IsValidEmail(req.Email) {
		return nil, errors.New("invalid email format")
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate tokens
	accessToken, refreshToken, err := utils.GenerateTokens(user, s.config)
	if err != nil {
		return nil, fmt.Errorf("error generating tokens: %w", err)
	}

	// Clear password before returning
	user.Password = ""

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

// ValidateToken validates a JWT token
func (s *AuthService) ValidateToken(tokenString string) (map[string]interface{}, error) {
	claims, err := utils.ValidateToken(tokenString, s.config.JWT.Secret)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

// RefreshToken refreshes an access token using a refresh token
func (s *AuthService) RefreshToken(refreshTokenString string) (string, error) {
	claims, err := utils.ValidateToken(refreshTokenString, s.config.JWT.Secret)
	if err != nil {
		return "", err
	}

	// Verify it's a refresh token
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return "", errors.New("invalid token type")
	}

	// Get user
	userID, err := utils.ExtractUserIDFromToken(claims)
	if err != nil {
		return "", err
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return "", err
	}

	// Generate new access token
	accessToken, _, err := utils.GenerateTokens(user, s.config)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(userID string) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Clear password
	user.Password = ""
	return user, nil
}

// UpdateUserProfile updates a user's profile
func (s *AuthService) UpdateUserProfile(userID string, req *models.UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		if !utils.ValidateUserName(req.Name) {
			return nil, errors.New("name must be between 2 and 255 characters")
		}
		user.Name = req.Name
	}

	if req.Email != "" {
		if !utils.IsValidEmail(req.Email) {
			return nil, errors.New("invalid email format")
		}

		// Check if email already exists (excluding current user)
		otherUser, err := s.userRepo.GetByEmail(req.Email)
		if err == nil && otherUser.ID != userID {
			return nil, errors.New("email already in use")
		}

		user.Email = req.Email
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	// Clear password
	user.Password = ""
	return user, nil
}
