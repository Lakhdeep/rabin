package v1

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rabin/tictactoe/internal/auth"
	"github.com/rabin/tictactoe/internal/storage"
	"github.com/rabin/tictactoe/internal/user"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	userRepo   storage.UserRepositoryInterface
	jwtService *auth.JWTService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userRepo storage.UserRepositoryInterface, jwtService *auth.JWTService) *AuthHandler {
	return &AuthHandler{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req user.RegisterRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"code":  "INVALID_REQUEST",
		})
		return
	}

	// Validate registration fields
	if err := user.ValidateRegistration(req.Email, req.Username, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "VALIDATION_ERROR",
		})
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process registration",
			"code":  "INTERNAL_ERROR",
		})
		return
	}

	// Create user
	newUser, err := h.userRepo.Create(req.Email, req.Username, hashedPassword)
	if err != nil {
		// Check for duplicate email/username errors
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			if strings.Contains(err.Error(), "email") {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Email already registered",
					"code":  "DUPLICATE_EMAIL",
				})
				return
			}
			if strings.Contains(err.Error(), "username") {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Username already taken",
					"code":  "DUPLICATE_USERNAME",
				})
				return
			}
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
			"code":  "INTERNAL_ERROR",
		})
		return
	}

	// Return created user (without password hash)
	c.JSON(http.StatusCreated, gin.H{
		"id":         newUser.ID,
		"email":      newUser.Email,
		"username":   newUser.Username,
		"created_at": newUser.CreatedAt,
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req user.LoginRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"code":  "INVALID_REQUEST",
		})
		return
	}

	// Get user by email
	existingUser, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
			"code":  "INVALID_CREDENTIALS",
		})
		return
	}

	// Verify password
	if err := auth.VerifyPassword(existingUser.PasswordHash, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
			"code":  "INVALID_CREDENTIALS",
		})
		return
	}

	// Generate JWT token
	token, err := h.jwtService.GenerateToken(existingUser.ID, existingUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
			"code":  "INTERNAL_ERROR",
		})
		return
	}

	// Return token and user info
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       existingUser.ID,
			"email":    existingUser.Email,
			"username": existingUser.Username,
		},
	})
}

// GetCurrentUser handles getting the current authenticated user
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
			"code":  "UNAUTHORIZED",
		})
		return
	}

	// Get user from database
	existingUser, err := h.userRepo.GetByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
				"code":  "USER_NOT_FOUND",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user",
			"code":  "INTERNAL_ERROR",
		})
		return
	}

	// Return user with statistics
	c.JSON(http.StatusOK, gin.H{
		"id":          existingUser.ID,
		"email":       existingUser.Email,
		"username":    existingUser.Username,
		"total_games": existingUser.TotalGames,
		"wins":        existingUser.Wins,
		"losses":      existingUser.Losses,
		"draws":       existingUser.Draws,
		"created_at":  existingUser.CreatedAt,
	})
}
