package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// Game represents a game in the database
type Game struct {
	ID          int             `json:"id"`
	UserID      int             `json:"user_id"`
	Difficulty  string          `json:"difficulty"`
	Result      sql.NullString  `json:"result"`
	BoardState  json.RawMessage `json:"board_state"`
	CurrentTurn sql.NullString  `json:"current_turn"`
	CreatedAt   time.Time       `json:"created_at"`
	CompletedAt sql.NullTime    `json:"completed_at"`
}

// GameMove represents a move in a game
type GameMove struct {
	ID         int       `json:"id"`
	GameID     int       `json:"game_id"`
	Position   int       `json:"position"`
	Player     string    `json:"player"`
	MoveNumber int       `json:"move_number"`
	CreatedAt  time.Time `json:"created_at"`
}

// GameRepository handles game database operations
type GameRepository struct {
	db *sql.DB
}

// NewGameRepository creates a new game repository
func NewGameRepository(db *sql.DB) *GameRepository {
	return &GameRepository{db: db}
}

// Create creates a new game
func (r *GameRepository) Create(userID int, difficulty string, boardState []byte) (*Game, error) {
	query := `
		INSERT INTO games (user_id, difficulty, board_state, current_turn)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, difficulty, result, board_state, current_turn, created_at, completed_at
	`

	game := &Game{}
	err := r.db.QueryRow(query, userID, difficulty, boardState, "X").Scan(
		&game.ID,
		&game.UserID,
		&game.Difficulty,
		&game.Result,
		&game.BoardState,
		&game.CurrentTurn,
		&game.CreatedAt,
		&game.CompletedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create game: %w", err)
	}

	return game, nil
}

// GetByID retrieves a game by ID
func (r *GameRepository) GetByID(id int) (*Game, error) {
	query := `
		SELECT id, user_id, difficulty, result, board_state, current_turn, created_at, completed_at
		FROM games
		WHERE id = $1
	`

	game := &Game{}
	err := r.db.QueryRow(query, id).Scan(
		&game.ID,
		&game.UserID,
		&game.Difficulty,
		&game.Result,
		&game.BoardState,
		&game.CurrentTurn,
		&game.CreatedAt,
		&game.CompletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("game not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	return game, nil
}

// Update updates a game's state
func (r *GameRepository) Update(game *Game) error {
	query := `
		UPDATE games
		SET board_state = $1, current_turn = $2, result = $3, completed_at = $4
		WHERE id = $5
	`

	_, err := r.db.Exec(query, game.BoardState, game.CurrentTurn, game.Result, game.CompletedAt, game.ID)
	if err != nil {
		return fmt.Errorf("failed to update game: %w", err)
	}

	return nil
}

// ListByUserID retrieves all games for a user
func (r *GameRepository) ListByUserID(userID int) ([]*Game, error) {
	query := `
		SELECT id, user_id, difficulty, result, board_state, current_turn, created_at, completed_at
		FROM games
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list games: %w", err)
	}
	defer rows.Close()

	var games []*Game
	for rows.Next() {
		game := &Game{}
		err := rows.Scan(
			&game.ID,
			&game.UserID,
			&game.Difficulty,
			&game.Result,
			&game.BoardState,
			&game.CurrentTurn,
			&game.CreatedAt,
			&game.CompletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan game: %w", err)
		}
		games = append(games, game)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating games: %w", err)
	}

	return games, nil
}

// AddMove adds a move to a game
func (r *GameRepository) AddMove(gameID, position int, player string, moveNumber int) error {
	query := `
		INSERT INTO game_moves (game_id, position, player, move_number)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(query, gameID, position, player, moveNumber)
	if err != nil {
		return fmt.Errorf("failed to add move: %w", err)
	}

	return nil
}

// GetMoves retrieves all moves for a game
func (r *GameRepository) GetMoves(gameID int) ([]*GameMove, error) {
	query := `
		SELECT id, game_id, position, player, move_number, created_at
		FROM game_moves
		WHERE game_id = $1
		ORDER BY move_number ASC
	`

	rows, err := r.db.Query(query, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get moves: %w", err)
	}
	defer rows.Close()

	var moves []*GameMove
	for rows.Next() {
		move := &GameMove{}
		err := rows.Scan(
			&move.ID,
			&move.GameID,
			&move.Position,
			&move.Player,
			&move.MoveNumber,
			&move.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan move: %w", err)
		}
		moves = append(moves, move)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating moves: %w", err)
	}

	return moves, nil
}
