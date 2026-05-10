package v1

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rabin/tictactoe/internal/game"
	"github.com/rabin/tictactoe/internal/game/ai"
	"github.com/rabin/tictactoe/internal/storage"
)

// GameHandler handles game-related requests
type GameHandler struct {
	gameRepo storage.GameRepositoryInterface
	userRepo storage.UserRepositoryInterface
}

// NewGameHandler creates a new game handler
func NewGameHandler(gameRepo storage.GameRepositoryInterface, userRepo storage.UserRepositoryInterface) *GameHandler {
	return &GameHandler{
		gameRepo: gameRepo,
		userRepo: userRepo,
	}
}

// CreateGameRequest represents a request to create a new game
type CreateGameRequest struct {
	Difficulty string `json:"difficulty" binding:"required"`
}

// MakeMoveRequest represents a request to make a move
type MakeMoveRequest struct {
	Position int `json:"position" binding:"min=0,max=8"`
}

// CreateGame handles game creation
func (h *GameHandler) CreateGame(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
			"code":  "UNAUTHORIZED",
		})
		return
	}

	var req CreateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"code":  "INVALID_REQUEST",
		})
		return
	}

	// Validate difficulty
	if err := game.ValidateDifficulty(req.Difficulty); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid difficulty level. Must be easy, medium, hard, or impossible",
			"code":  "INVALID_DIFFICULTY",
		})
		return
	}

	// Create new game
	userIDInt := userID.(int64)
	newGame := game.NewGame(userIDInt, game.Difficulty(req.Difficulty))

	// Serialize board state
	boardState, err := json.Marshal(newGame.Board.ToSlice())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create game",
			"code":  "INTERNAL_ERROR",
		})
		return
	}

	// Save to database
	dbGame, err := h.gameRepo.Create(int(userIDInt), req.Difficulty, boardState)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create game",
			"code":  "INTERNAL_ERROR",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":           dbGame.ID,
		"user_id":      dbGame.UserID,
		"difficulty":   dbGame.Difficulty,
		"board":        newGame.Board.ToSlice(),
		"current_turn": "X",
		"status":       "active",
		"created_at":   dbGame.CreatedAt,
	})
}

// GetGame retrieves a game by ID
func (h *GameHandler) GetGame(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
			"code":  "UNAUTHORIZED",
		})
		return
	}

	gameIDStr := c.Param("id")
	gameID, err := strconv.Atoi(gameIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID",
			"code":  "INVALID_GAME_ID",
		})
		return
	}

	// Get game from database
	dbGame, err := h.gameRepo.GetByID(gameID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
			"code":  "GAME_NOT_FOUND",
		})
		return
	}

	// Check ownership
	if dbGame.UserID != int(userID.(int64)) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have permission to access this game",
			"code":  "FORBIDDEN",
		})
		return
	}

	// Parse board state
	var board []game.CellValue
	if err := json.Unmarshal(dbGame.BoardState, &board); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse game state",
			"code":  "INTERNAL_ERROR",
		})
		return
	}

	response := gin.H{
		"id":           dbGame.ID,
		"user_id":      dbGame.UserID,
		"difficulty":   dbGame.Difficulty,
		"board":        board,
		"current_turn": dbGame.CurrentTurn.String,
		"created_at":   dbGame.CreatedAt,
	}

	if dbGame.Result.Valid {
		response["result"] = dbGame.Result.String
		response["status"] = dbGame.Result.String
	} else {
		response["status"] = "active"
	}

	if dbGame.CompletedAt.Valid {
		response["completed_at"] = dbGame.CompletedAt.Time
	}

	c.JSON(http.StatusOK, response)
}

// MakeMove handles making a move in a game
func (h *GameHandler) MakeMove(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
			"code":  "UNAUTHORIZED",
		})
		return
	}

	gameIDStr := c.Param("id")
	gameID, err := strconv.Atoi(gameIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid game ID",
			"code":  "INVALID_GAME_ID",
		})
		return
	}

	var req MakeMoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"code":  "INVALID_REQUEST",
		})
		return
	}

	// Get game from database
	dbGame, err := h.gameRepo.GetByID(gameID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Game not found",
			"code":  "GAME_NOT_FOUND",
		})
		return
	}

	// Check ownership
	if dbGame.UserID != int(userID.(int64)) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have permission to access this game",
			"code":  "FORBIDDEN",
		})
		return
	}

	// Check if game is already completed
	if dbGame.Result.Valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Game already completed",
			"code":  "GAME_COMPLETED",
		})
		return
	}

	// Parse board state
	var boardSlice []game.CellValue
	if err := json.Unmarshal(dbGame.BoardState, &boardSlice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse game state",
			"code":  "INTERNAL_ERROR",
		})
		return
	}

	// Reconstruct game object
	board := game.NewBoard()
	board.FromSlice(boardSlice)

	g := &game.Game{
		ID:          int64(dbGame.ID),
		UserID:      int64(dbGame.UserID),
		Board:       board,
		CurrentTurn: game.CellValue(dbGame.CurrentTurn.String),
		Difficulty:  game.Difficulty(dbGame.Difficulty),
		Status:      game.StatusActive,
		CreatedAt:   dbGame.CreatedAt,
	}

	// Get move count for this game
	moves, err := h.gameRepo.GetMoves(gameID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve game moves",
			"code":  "INTERNAL_ERROR",
		})
		return
	}
	moveNumber := len(moves) + 1

	// Make user's move
	if err := g.MakeMove(req.Position, game.X); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "INVALID_MOVE",
		})
		return
	}

	// Record user's move
	if err := h.gameRepo.AddMove(gameID, req.Position, "X", moveNumber); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to record move",
			"code":  "INTERNAL_ERROR",
		})
		return
	}

	// Check if game ended after user's move
	if !g.IsActive() {
		// Update game result and user statistics
		if err := h.updateGameResult(dbGame, g); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update game",
				"code":  "INTERNAL_ERROR",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"board":  g.Board.ToSlice(),
			"result": g.Result,
			"status": g.Status,
		})
		return
	}

	// AI's turn
	aiStrategy, err := ai.NewAI(g.Difficulty)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to initialize AI",
			"code":  "INTERNAL_ERROR",
		})
		return
	}

	aiMove, err := ai.GetMoveWithTimeout(aiStrategy, g.Board, 1*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "AI failed to make move",
			"code":  "AI_ERROR",
		})
		return
	}

	// Make AI's move
	if err := g.MakeMove(aiMove, game.O); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "AI move failed",
			"code":  "AI_ERROR",
		})
		return
	}

	// Record AI's move
	if err := h.gameRepo.AddMove(gameID, aiMove, "O", moveNumber+1); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to record AI move",
			"code":  "INTERNAL_ERROR",
		})
		return
	}

	// Update game result
	if err := h.updateGameResult(dbGame, g); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update game",
			"code":  "INTERNAL_ERROR",
		})
		return
	}

	response := gin.H{
		"board":   g.Board.ToSlice(),
		"ai_move": aiMove,
		"status":  g.Status,
	}

	if !g.IsActive() {
		response["result"] = g.Result
	} else {
		response["current_turn"] = g.CurrentTurn
	}

	c.JSON(http.StatusOK, response)
}

// ListGames lists all games for the authenticated user
func (h *GameHandler) ListGames(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
			"code":  "UNAUTHORIZED",
		})
		return
	}

	games, err := h.gameRepo.ListByUserID(int(userID.(int64)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve games",
			"code":  "INTERNAL_ERROR",
		})
		return
	}

	response := make([]gin.H, 0, len(games))
	for i := range games {
		dbGame := &games[i]
		gameData := gin.H{
			"id":         dbGame.ID,
			"difficulty": dbGame.Difficulty,
			"created_at": dbGame.CreatedAt,
		}

		if dbGame.Result.Valid {
			gameData["result"] = dbGame.Result.String
			gameData["status"] = dbGame.Result.String
		} else {
			gameData["status"] = "active"
		}

		if dbGame.CompletedAt.Valid {
			gameData["completed_at"] = dbGame.CompletedAt.Time
		}

		response = append(response, gameData)
	}

	c.JSON(http.StatusOK, gin.H{
		"games": response,
	})
}

// updateGameResult updates the game in database and user statistics
func (h *GameHandler) updateGameResult(dbGame *storage.Game, g *game.Game) error {
	// Serialize board
	boardState, err := json.Marshal(g.Board.ToSlice())
	if err != nil {
		return err
	}

	dbGame.BoardState = boardState

	if g.IsActive() {
		dbGame.CurrentTurn = sql.NullString{String: string(g.CurrentTurn), Valid: true}
	} else {
		// Game ended
		dbGame.Result = sql.NullString{String: string(g.Result), Valid: true}
		now := time.Now()
		dbGame.CompletedAt = sql.NullTime{Time: now, Valid: true}

		// Update user statistics
		if err := h.userRepo.UpdateStatistics(dbGame.UserID, string(g.Result)); err != nil {
			return err
		}
	}

	return h.gameRepo.Update(dbGame)
}
