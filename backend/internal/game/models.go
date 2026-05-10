package game

import (
	"fmt"
	"time"
)

// GameStatus represents the current state of a game
type GameStatus string

const (
	StatusActive GameStatus = "active"
	StatusWon    GameStatus = "won"
	StatusLost   GameStatus = "lost"
	StatusDraw   GameStatus = "draw"
)

// GameResult represents the outcome of a game from the user's perspective
type GameResult string

const (
	ResultWin  GameResult = "win"
	ResultLoss GameResult = "loss"
	ResultDraw GameResult = "draw"
	ResultNone GameResult = ""
)

// Difficulty represents AI difficulty level
type Difficulty string

const (
	DifficultyEasy       Difficulty = "easy"
	DifficultyMedium     Difficulty = "medium"
	DifficultyHard       Difficulty = "hard"
	DifficultyImpossible Difficulty = "impossible"
)

// ValidateDifficulty checks if a difficulty string is valid
func ValidateDifficulty(d string) error {
	switch Difficulty(d) {
	case DifficultyEasy, DifficultyMedium, DifficultyHard, DifficultyImpossible:
		return nil
	default:
		return fmt.Errorf("invalid difficulty: must be easy, medium, hard, or impossible")
	}
}

// Game represents a tic-tac-toe game
type Game struct {
	ID          int64
	UserID      int64
	Board       *Board
	CurrentTurn CellValue
	Difficulty  Difficulty
	Status      GameStatus
	Result      GameResult
	CreatedAt   time.Time
	CompletedAt *time.Time
}

// NewGame creates a new game with the given user ID and difficulty
func NewGame(userID int64, difficulty Difficulty) *Game {
	return &Game{
		UserID:      userID,
		Board:       NewBoard(),
		CurrentTurn: X, // User always starts as X
		Difficulty:  difficulty,
		Status:      StatusActive,
		Result:      ResultNone,
		CreatedAt:   time.Now(),
	}
}

// MakeMove attempts to make a move at the given position
// Returns an error if the move is invalid
func (g *Game) MakeMove(position int, player CellValue) error {
	// Validate that it's the correct player's turn
	if player != g.CurrentTurn {
		return fmt.Errorf("not your turn")
	}

	// Validate the move
	if err := g.Board.ValidateMove(position); err != nil {
		return err
	}

	// Make the move
	if err := g.Board.Set(position, player); err != nil {
		return err
	}

	// Check for game end conditions
	g.updateGameState()

	// Switch turns if game is still active
	if g.Status == StatusActive {
		g.switchTurn()
	}

	return nil
}

// switchTurn alternates between X and O
func (g *Game) switchTurn() {
	if g.CurrentTurn == X {
		g.CurrentTurn = O
	} else {
		g.CurrentTurn = X
	}
}

// updateGameState checks win/draw conditions and updates game status
func (g *Game) updateGameState() {
	winner := g.Board.CheckWinner()

	if winner == X {
		// User (X) won
		g.Status = StatusWon
		g.Result = ResultWin
		now := time.Now()
		g.CompletedAt = &now
	} else if winner == O {
		// AI (O) won
		g.Status = StatusLost
		g.Result = ResultLoss
		now := time.Now()
		g.CompletedAt = &now
	} else if g.Board.IsFull() {
		// Draw
		g.Status = StatusDraw
		g.Result = ResultDraw
		now := time.Now()
		g.CompletedAt = &now
	}
	// Otherwise game remains active
}

// IsActive returns true if the game is still in progress
func (g *Game) IsActive() bool {
	return g.Status == StatusActive
}

// IsUserTurn returns true if it's the user's turn (X)
func (g *Game) IsUserTurn() bool {
	return g.CurrentTurn == X
}

// IsAITurn returns true if it's the AI's turn (O)
func (g *Game) IsAITurn() bool {
	return g.CurrentTurn == O
}
