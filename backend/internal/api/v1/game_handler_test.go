package v1

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rabin/tictactoe/internal/game"
	"github.com/rabin/tictactoe/internal/storage"
)

// Mock GameRepository for testing
type mockGameRepo struct {
	games map[int]*storage.Game
	moves map[int][]storage.GameMove
	nextID int
}

func newMockGameRepo() *mockGameRepo {
	return &mockGameRepo{
		games: make(map[int]*storage.Game),
		moves: make(map[int][]storage.GameMove),
		nextID: 1,
	}
}

func (m *mockGameRepo) Create(userID int, difficulty string, boardState []byte) (*storage.Game, error) {
	game := &storage.Game{
		ID:          m.nextID,
		UserID:      userID,
		Difficulty:  difficulty,
		BoardState:  boardState,
		CurrentTurn: sql.NullString{String: "X", Valid: true},
		CreatedAt:   time.Now(),
	}
	m.games[m.nextID] = game
	m.nextID++
	return game, nil
}

func (m *mockGameRepo) GetByID(id int) (*storage.Game, error) {
	if game, ok := m.games[id]; ok {
		return game, nil
	}
	return nil, sql.ErrNoRows
}

func (m *mockGameRepo) Update(game *storage.Game) error {
	m.games[game.ID] = game
	return nil
}

func (m *mockGameRepo) ListByUserID(userID int) ([]storage.Game, error) {
	var games []storage.Game
	for _, g := range m.games {
		if g.UserID == userID {
			games = append(games, *g)
		}
	}
	return games, nil
}

func (m *mockGameRepo) AddMove(gameID, position int, player string, moveNumber int) error {
	move := storage.GameMove{
		GameID:     gameID,
		Position:   position,
		Player:     player,
		MoveNumber: moveNumber,
		CreatedAt:  time.Now(),
	}
	m.moves[gameID] = append(m.moves[gameID], move)
	return nil
}

func (m *mockGameRepo) GetMoves(gameID int) ([]storage.GameMove, error) {
	return m.moves[gameID], nil
}

func TestCreateGame(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    map[string]string
		authenticated  bool
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name: "Valid game creation with easy difficulty",
			requestBody: map[string]string{
				"difficulty": "easy",
			},
			authenticated:  true,
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["id"] == nil {
					t.Error("Expected game ID in response")
				}
				if resp["difficulty"] != "easy" {
					t.Errorf("Expected difficulty easy, got %v", resp["difficulty"])
				}
				if resp["current_turn"] != "X" {
					t.Error("Expected current turn to be X")
				}
				if resp["status"] != "active" {
					t.Error("Expected status to be active")
				}
			},
		},
		{
			name: "Valid game creation with medium difficulty",
			requestBody: map[string]string{
				"difficulty": "medium",
			},
			authenticated:  true,
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["difficulty"] != "medium" {
					t.Errorf("Expected difficulty medium, got %v", resp["difficulty"])
				}
			},
		},
		{
			name: "Valid game creation with hard difficulty",
			requestBody: map[string]string{
				"difficulty": "hard",
			},
			authenticated:  true,
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["difficulty"] != "hard" {
					t.Errorf("Expected difficulty hard, got %v", resp["difficulty"])
				}
			},
		},
		{
			name: "Valid game creation with impossible difficulty",
			requestBody: map[string]string{
				"difficulty": "impossible",
			},
			authenticated:  true,
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["difficulty"] != "impossible" {
					t.Errorf("Expected difficulty impossible, got %v", resp["difficulty"])
				}
			},
		},
		{
			name: "Invalid difficulty rejected",
			requestBody: map[string]string{
				"difficulty": "super-easy",
			},
			authenticated:  true,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["code"] != "INVALID_DIFFICULTY" {
					t.Error("Expected INVALID_DIFFICULTY code")
				}
			},
		},
		{
			name: "Unauthenticated request rejected",
			requestBody: map[string]string{
				"difficulty": "easy",
			},
			authenticated:  false,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Missing difficulty rejected",
			requestBody:    map[string]string{},
			authenticated:  true,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockGameRepo := newMockGameRepo()
			mockUserRepo := newMockUserRepo()
			handler := NewGameHandler(mockGameRepo, mockUserRepo)

			router := gin.New()
			if tt.authenticated {
				router.Use(func(c *gin.Context) {
					c.Set("user_id", int64(1))
					c.Next()
				})
			}
			router.POST("/games", handler.CreateGame)

			// Create request
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/games", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Execute
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.checkResponse != nil && w.Code == tt.expectedStatus {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				tt.checkResponse(t, response)
			}
		})
	}
}

func TestMakeMove(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupGame      func(*mockGameRepo) int // Returns game ID
		requestBody    map[string]int
		authenticated  bool
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name: "Valid first move",
			setupGame: func(repo *mockGameRepo) int {
				// Create empty game
				board := game.NewBoard()
				boardState, _ := json.Marshal(board.ToSlice())
				g, _ := repo.Create(1, "easy", boardState)
				return g.ID
			},
			requestBody: map[string]int{
				"position": 0,
			},
			authenticated:  true,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["board"] == nil {
					t.Error("Expected board in response")
				}
				if resp["ai_move"] == nil {
					t.Error("Expected ai_move in response")
				}
			},
		},
		{
			name: "Invalid position out of range",
			setupGame: func(repo *mockGameRepo) int {
				board := game.NewBoard()
				boardState, _ := json.Marshal(board.ToSlice())
				g, _ := repo.Create(1, "easy", boardState)
				return g.ID
			},
			requestBody: map[string]int{
				"position": 9,
			},
			authenticated:  true,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				// Position validation is handled by binding, which returns INVALID_REQUEST
				code, _ := resp["code"].(string)
				if code != "INVALID_REQUEST" && code != "INVALID_MOVE" {
					t.Errorf("Expected INVALID_REQUEST or INVALID_MOVE code, got %s", code)
				}
			},
		},
		{
			name: "Move on occupied position rejected",
			setupGame: func(repo *mockGameRepo) int {
				board := game.NewBoard()
				board.Set(0, game.X)
				boardState, _ := json.Marshal(board.ToSlice())
				g, _ := repo.Create(1, "easy", boardState)
				return g.ID
			},
			requestBody: map[string]int{
				"position": 0,
			},
			authenticated:  true,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["code"] != "INVALID_MOVE" {
					t.Error("Expected INVALID_MOVE code")
				}
			},
		},
		{
			name: "Unauthenticated move rejected",
			setupGame: func(repo *mockGameRepo) int {
				board := game.NewBoard()
				boardState, _ := json.Marshal(board.ToSlice())
				g, _ := repo.Create(1, "easy", boardState)
				return g.ID
			},
			requestBody: map[string]int{
				"position": 0,
			},
			authenticated:  false,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Move in another user's game rejected",
			setupGame: func(repo *mockGameRepo) int {
				board := game.NewBoard()
				boardState, _ := json.Marshal(board.ToSlice())
				g, _ := repo.Create(2, "easy", boardState) // Different user
				return g.ID
			},
			requestBody: map[string]int{
				"position": 0,
			},
			authenticated:  true,
			expectedStatus: http.StatusForbidden,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["code"] != "FORBIDDEN" {
					t.Error("Expected FORBIDDEN code")
				}
			},
		},
		{
			name: "Move in completed game rejected",
			setupGame: func(repo *mockGameRepo) int {
				board := game.NewBoard()
				boardState, _ := json.Marshal(board.ToSlice())
				g, _ := repo.Create(1, "easy", boardState)
				g.Result = sql.NullString{String: "win", Valid: true}
				repo.Update(g)
				return g.ID
			},
			requestBody: map[string]int{
				"position": 0,
			},
			authenticated:  true,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["code"] != "GAME_COMPLETED" {
					t.Error("Expected GAME_COMPLETED code")
				}
			},
		},
		{
			name: "Nonexistent game returns 404",
			setupGame: func(repo *mockGameRepo) int {
				return 999 // Non-existent game ID
			},
			requestBody: map[string]int{
				"position": 0,
			},
			authenticated:  true,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockGameRepo := newMockGameRepo()
			mockUserRepo := newMockUserRepo()
			gameID := tt.setupGame(mockGameRepo)
			handler := NewGameHandler(mockGameRepo, mockUserRepo)

			router := gin.New()
			if tt.authenticated {
				router.Use(func(c *gin.Context) {
					c.Set("user_id", int64(1))
					c.Next()
				})
			}
			router.POST("/games/:id/move", handler.MakeMove)

			// Create request
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/games/"+strconv.Itoa(gameID)+"/move", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Execute
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.checkResponse != nil && w.Code == tt.expectedStatus {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				tt.checkResponse(t, response)
			}
		})
	}
}

func TestGameCompletionAndScoreUpdates(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	mockGameRepo := newMockGameRepo()
	mockUserRepo := newMockUserRepo()
	mockUserRepo.users["test@example.com"] = &storage.User{
		ID:         1,
		Email:      "test@example.com",
		Username:   "testuser",
		TotalGames: 0,
		Wins:       0,
		Losses:     0,
		Draws:      0,
	}
	handler := NewGameHandler(mockGameRepo, mockUserRepo)

	// Create a game that will result in a win
	board := game.NewBoard()
	// Set up a near-win scenario: X just needs position 2 to win (0-1-2)
	board.Set(0, game.X)
	board.Set(1, game.X)
	board.Set(3, game.O)
	board.Set(4, game.O)
	boardState, _ := json.Marshal(board.ToSlice())
	g, _ := mockGameRepo.Create(1, "easy", boardState)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
		c.Next()
	})
	router.POST("/games/:id/move", handler.MakeMove)

	// Make winning move
	body, _ := json.Marshal(map[string]int{"position": 2})
	req, _ := http.NewRequest(http.MethodPost, "/games/"+strconv.Itoa(g.ID)+"/move", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	// Check that result is set
	if response["result"] != "win" {
		t.Errorf("Expected result to be 'win', got %v", response["result"])
	}

	// Verify user statistics were updated
	user, _ := mockUserRepo.GetByID(1)
	if user.Wins != 1 {
		t.Errorf("Expected user wins to be 1, got %d", user.Wins)
	}
	if user.TotalGames != 1 {
		t.Errorf("Expected user total_games to be 1, got %d", user.TotalGames)
	}

	// Verify game was marked as completed
	updatedGame, _ := mockGameRepo.GetByID(g.ID)
	if !updatedGame.Result.Valid || updatedGame.Result.String != "win" {
		t.Error("Expected game result to be set to 'win'")
	}
	if !updatedGame.CompletedAt.Valid {
		t.Error("Expected game completed_at to be set")
	}
}

func TestGetUserStats(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		setupUser      func(*mockUserRepo)
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:   "Get stats for user with games",
			userID: "1",
			setupUser: func(repo *mockUserRepo) {
				repo.users["user1"] = &storage.User{
					ID:         1,
					Username:   "testuser",
					Email:      "test@example.com",
					TotalGames: 10,
					Wins:       7,
					Losses:     2,
					Draws:      1,
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["user_id"] != float64(1) {
					t.Errorf("Expected user_id 1, got %v", resp["user_id"])
				}
				if resp["username"] != "testuser" {
					t.Errorf("Expected username testuser, got %v", resp["username"])
				}
				if resp["total_games"] != float64(10) {
					t.Errorf("Expected total_games 10, got %v", resp["total_games"])
				}
				if resp["wins"] != float64(7) {
					t.Errorf("Expected wins 7, got %v", resp["wins"])
				}
				if resp["losses"] != float64(2) {
					t.Errorf("Expected losses 2, got %v", resp["losses"])
				}
				if resp["draws"] != float64(1) {
					t.Errorf("Expected draws 1, got %v", resp["draws"])
				}
				// Win rate should be 70%
				winRate := resp["win_rate"].(float64)
				if winRate < 69.9 || winRate > 70.1 {
					t.Errorf("Expected win_rate ~70%%, got %v", winRate)
				}
			},
		},
		{
			name:   "Get stats for user with no games",
			userID: "2",
			setupUser: func(repo *mockUserRepo) {
				repo.users["user2"] = &storage.User{
					ID:         2,
					Username:   "newuser",
					Email:      "new@example.com",
					TotalGames: 0,
					Wins:       0,
					Losses:     0,
					Draws:      0,
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["total_games"] != float64(0) {
					t.Errorf("Expected total_games 0, got %v", resp["total_games"])
				}
				if resp["win_rate"] != float64(0) {
					t.Errorf("Expected win_rate 0 for user with no games, got %v", resp["win_rate"])
				}
			},
		},
		{
			name:   "Get stats for nonexistent user",
			userID: "999",
			setupUser: func(repo *mockUserRepo) {
				// No user setup
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["code"] != "USER_NOT_FOUND" {
					t.Errorf("Expected code USER_NOT_FOUND, got %v", resp["code"])
				}
			},
		},
		{
			name:           "Invalid user ID format",
			userID:         "abc",
			setupUser:      func(repo *mockUserRepo) {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["code"] != "INVALID_USER_ID" {
					t.Errorf("Expected code INVALID_USER_ID, got %v", resp["code"])
				}
			},
		},
		{
			name:   "Perfect win rate (all wins)",
			userID: "3",
			setupUser: func(repo *mockUserRepo) {
				repo.users["user3"] = &storage.User{
					ID:         3,
					Username:   "winner",
					Email:      "winner@example.com",
					TotalGames: 5,
					Wins:       5,
					Losses:     0,
					Draws:      0,
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				winRate := resp["win_rate"].(float64)
				if winRate < 99.9 || winRate > 100.1 {
					t.Errorf("Expected win_rate 100%%, got %v", winRate)
				}
			},
		},
		{
			name:   "Zero win rate (all losses)",
			userID: "4",
			setupUser: func(repo *mockUserRepo) {
				repo.users["user4"] = &storage.User{
					ID:         4,
					Username:   "loser",
					Email:      "loser@example.com",
					TotalGames: 5,
					Wins:       0,
					Losses:     5,
					Draws:      0,
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				winRate := resp["win_rate"].(float64)
				if winRate != 0 {
					t.Errorf("Expected win_rate 0%%, got %v", winRate)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockUserRepo := newMockUserRepo()
			tt.setupUser(mockUserRepo)
			handler := NewUserHandler(mockUserRepo)

			router := gin.New()
			router.GET("/users/:id/stats", handler.GetUserStats)

			// Create request
			req, _ := http.NewRequest(http.MethodGet, "/users/"+tt.userID+"/stats", nil)

			// Execute
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.checkResponse != nil && w.Code == tt.expectedStatus {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				tt.checkResponse(t, response)
			}
		})
	}
}
