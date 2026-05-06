package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"task-management-api/pkg/middleware"
	"task-management-api/pkg/models"
	"task-management-api/pkg/services"
)

// AuthController handles authentication endpoints
type AuthController struct {
	authService *services.AuthService
}

// NewAuthController creates a new auth controller
func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration request"
// @Success 201 {object} models.SuccessResponse{data=models.User}
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/register [post]
func (ctrl *AuthController) Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	user, err := ctrl.authService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "registration_error",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Data:    user,
		Message: "User registered successfully",
		Code:    http.StatusCreated,
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login request"
// @Success 200 {object} models.SuccessResponse{data=models.LoginResponse}
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	response, err := ctrl.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "login_error",
			Message: err.Error(),
			Code:    http.StatusUnauthorized,
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Data:    response,
		Message: "Login successful",
		Code:    http.StatusOK,
	})
}

// Refresh godoc
// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Refresh Token with Bearer prefix"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/refresh [post]
func (ctrl *AuthController) Refresh(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "Authorization header required",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "Invalid authorization header format",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	refreshToken := parts[1]

	accessToken, err := ctrl.authService.RefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "refresh_error",
			Message: err.Error(),
			Code:    http.StatusUnauthorized,
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Data: map[string]string{
			"access_token": accessToken,
		},
		Message: "Token refreshed successfully",
		Code:    http.StatusOK,
	})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current authenticated user's profile
// @Tags auth
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} models.SuccessResponse{data=models.User}
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/me [get]
func (ctrl *AuthController) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	user, err := ctrl.authService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "error",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Data:    user,
		Message: "Profile retrieved successfully",
		Code:    http.StatusOK,
	})
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update current authenticated user's profile
// @Tags auth
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body models.UpdateUserRequest true "Update request"
// @Success 200 {object} models.SuccessResponse{data=models.User}
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/me [put]
func (ctrl *AuthController) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	user, err := ctrl.authService.UpdateUserProfile(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "update_error",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Data:    user,
		Message: "Profile updated successfully",
		Code:    http.StatusOK,
	})
}
