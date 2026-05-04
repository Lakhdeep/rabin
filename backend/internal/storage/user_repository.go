package storage

import (
	"database/sql"
	"fmt"
	"time"
)

// User represents a user in the database
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

// UserRepository handles user database operations
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(email, username, passwordHash string) (*User, error) {
	query := `
		INSERT INTO users (email, username, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, email, username, password_hash, total_games, wins, losses, draws, created_at, updated_at
	`

	user := &User{}
	err := r.db.QueryRow(query, email, username, passwordHash).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.TotalGames,
		&user.Wins,
		&user.Losses,
		&user.Draws,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, email, username, password_hash, total_games, wins, losses, draws, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.TotalGames,
		&user.Wins,
		&user.Losses,
		&user.Draws,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int) (*User, error) {
	query := `
		SELECT id, email, username, password_hash, total_games, wins, losses, draws, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.TotalGames,
		&user.Wins,
		&user.Losses,
		&user.Draws,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}

// Update updates user information
func (r *UserRepository) Update(user *User) error {
	query := `
		UPDATE users
		SET email = $1, username = $2, total_games = $3, wins = $4, losses = $5, draws = $6, updated_at = NOW()
		WHERE id = $7
	`

	_, err := r.db.Exec(query, user.Email, user.Username, user.TotalGames, user.Wins, user.Losses, user.Draws, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete deletes a user
func (r *UserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// UpdateStatistics updates user game statistics atomically
func (r *UserRepository) UpdateStatistics(userID int, result string) error {
	var query string

	switch result {
	case "win":
		query = `
			UPDATE users
			SET total_games = total_games + 1, wins = wins + 1, updated_at = NOW()
			WHERE id = $1
		`
	case "loss":
		query = `
			UPDATE users
			SET total_games = total_games + 1, losses = losses + 1, updated_at = NOW()
			WHERE id = $1
		`
	case "draw":
		query = `
			UPDATE users
			SET total_games = total_games + 1, draws = draws + 1, updated_at = NOW()
			WHERE id = $1
		`
	default:
		return fmt.Errorf("invalid result: %s", result)
	}

	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to update statistics: %w", err)
	}

	return nil
}
