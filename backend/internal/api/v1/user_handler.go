package v1

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rabin/tictactoe/internal/storage"
)

// UserHandler handles user-related requests
type UserHandler struct {
	userRepo storage.UserRepositoryInterface
}

// NewUserHandler creates a new user handler
func NewUserHandler(userRepo storage.UserRepositoryInterface) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

// GetUserStats retrieves user statistics by user ID
// This endpoint is public and does not require authentication
func (h *UserHandler) GetUserStats(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
			"code":  "INVALID_USER_ID",
		})
		return
	}

	// Get user from database
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		if err == sql.ErrNoRows || err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
				"code":  "USER_NOT_FOUND",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user statistics",
			"code":  "INTERNAL_ERROR",
		})
		return
	}

	// Calculate win rate
	var winRate float64
	if user.TotalGames > 0 {
		winRate = float64(user.Wins) / float64(user.TotalGames) * 100
	} else {
		winRate = 0.0
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":     user.ID,
		"username":    user.Username,
		"total_games": user.TotalGames,
		"wins":        user.Wins,
		"losses":      user.Losses,
		"draws":       user.Draws,
		"win_rate":    winRate,
	})
}
