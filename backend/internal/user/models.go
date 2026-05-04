package user

import "time"

// User represents a user in the application
type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Never expose password hash in JSON
	TotalGames   int       `json:"total_games"`
	Wins         int       `json:"wins"`
	Losses       int       `json:"losses"`
	Draws        int       `json:"draws"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a successful login response
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// UserStats represents user statistics
type UserStats struct {
	UserID     int     `json:"user_id"`
	Username   string  `json:"username"`
	TotalGames int     `json:"total_games"`
	Wins       int     `json:"wins"`
	Losses     int     `json:"losses"`
	Draws      int     `json:"draws"`
	WinRate    float64 `json:"win_rate"`
}

// CalculateWinRate calculates the win rate percentage
func (u *User) CalculateWinRate() float64 {
	if u.TotalGames == 0 {
		return 0.0
	}
	return (float64(u.Wins) / float64(u.TotalGames)) * 100
}
