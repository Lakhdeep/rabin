package ai

import (
	"testing"
	"time"

	"github.com/rabin/tictactoe/internal/game"
)

func TestEasyAI(t *testing.T) {
	ai := &EasyAI{}

	t.Run("returns valid move on empty board", func(t *testing.T) {
		board := game.NewBoard()
		move, err := ai.GetMove(board)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if move < 0 || move > 8 {
			t.Errorf("move %d is out of range", move)
		}

		if !board.IsEmpty(move) {
			t.Errorf("move %d is not an empty position", move)
		}
	})

	t.Run("returns valid move on partially filled board", func(t *testing.T) {
		board := game.NewBoard()
		board.Set(0, game.X)
		board.Set(4, game.O)
		board.Set(8, game.X)

		move, err := ai.GetMove(board)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if move < 0 || move > 8 {
			t.Errorf("move %d is out of range", move)
		}

		if !board.IsEmpty(move) {
			t.Errorf("move %d is not an empty position", move)
		}

		// Verify it's one of the valid positions
		validMoves := []int{1, 2, 3, 5, 6, 7}
		isValid := false
		for _, valid := range validMoves {
			if move == valid {
				isValid = true
				break
			}
		}
		if !isValid {
			t.Errorf("move %d is not in valid positions %v", move, validMoves)
		}
	})

	t.Run("returns error on full board", func(t *testing.T) {
		board := game.NewBoard()
		// Fill all positions
		for i := 0; i < 9; i++ {
			if i%2 == 0 {
				board.Set(i, game.X)
			} else {
				board.Set(i, game.O)
			}
		}

		_, err := ai.GetMove(board)
		if err == nil {
			t.Error("expected error on full board")
		}
	})

	t.Run("generates different moves over time", func(t *testing.T) {
		// Test that EasyAI doesn't always return the same move
		board := game.NewBoard()
		moves := make(map[int]bool)

		// Try getting moves multiple times
		for i := 0; i < 20; i++ {
			move, err := ai.GetMove(board)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			moves[move] = true
			time.Sleep(1 * time.Millisecond) // Small delay to ensure different seed
		}

		// Should have gotten at least 2 different moves
		if len(moves) < 2 {
			t.Errorf("expected at least 2 different moves, got %d", len(moves))
		}
	})
}

func TestMediumAI(t *testing.T) {
	ai := &MediumAI{}

	t.Run("returns valid move", func(t *testing.T) {
		board := game.NewBoard()
		move, err := ai.GetMove(board)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if move < 0 || move > 8 {
			t.Errorf("move %d is out of range", move)
		}

		if !board.IsEmpty(move) {
			t.Errorf("move %d is not an empty position", move)
		}
	})

	t.Run("can block winning move", func(t *testing.T) {
		// Test that MediumAI sometimes blocks winning moves
		// X X _
		// _ O _
		// _ _ _
		board := game.NewBoard()
		board.Set(0, game.X)
		board.Set(1, game.X)
		board.Set(4, game.O)

		blockedCount := 0
		for i := 0; i < 10; i++ {
			testBoard := board.Clone()
			move, err := ai.GetMove(testBoard)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			// Position 2 blocks X from winning
			if move == 2 {
				blockedCount++
			}
		}

		// Should block at least some of the time (due to 50% optimal strategy)
		if blockedCount == 0 {
			t.Error("MediumAI never blocked winning move in 10 attempts")
		}
	})
}

func TestNewAI(t *testing.T) {
	tests := []struct {
		name       string
		difficulty game.Difficulty
		wantType   string
		wantError  bool
	}{
		{
			name:       "easy difficulty",
			difficulty: game.DifficultyEasy,
			wantType:   "*ai.EasyAI",
			wantError:  false,
		},
		{
			name:       "medium difficulty",
			difficulty: game.DifficultyMedium,
			wantType:   "*ai.MediumAI",
			wantError:  false,
		},
		{
			name:       "hard difficulty",
			difficulty: game.DifficultyHard,
			wantType:   "*ai.HardAI",
			wantError:  false,
		},
		{
			name:       "impossible difficulty",
			difficulty: game.DifficultyImpossible,
			wantType:   "*ai.ImpossibleAI",
			wantError:  false,
		},
		{
			name:       "invalid difficulty",
			difficulty: game.Difficulty("invalid"),
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ai, err := NewAI(tt.difficulty)

			if tt.wantError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if ai == nil {
				t.Fatal("expected AI instance, got nil")
			}
		})
	}
}

func TestGetMoveWithTimeout(t *testing.T) {
	t.Run("completes within timeout", func(t *testing.T) {
		ai := &EasyAI{}
		board := game.NewBoard()

		move, err := GetMoveWithTimeout(ai, board, 1*time.Second)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if move < 0 || move > 8 {
			t.Errorf("move %d is out of range", move)
		}
	})

	t.Run("handles timeout gracefully", func(t *testing.T) {
		// Create a slow AI that will timeout
		ai := &ImpossibleAI{}
		board := game.NewBoard()

		// Use very short timeout - should still get a fallback move
		move, err := GetMoveWithTimeout(ai, board, 1*time.Nanosecond)

		// Even with timeout, should get a valid move (fallback to random)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if move < 0 || move > 8 {
			t.Errorf("move %d is out of range", move)
		}
	})
}

func TestMinimaxAlgorithm(t *testing.T) {
	ai := &ImpossibleAI{}

	t.Run("blocks opponent winning move", func(t *testing.T) {
		// X X _
		// _ O _
		// _ _ _
		// AI should block at position 2
		board := game.NewBoard()
		board.Set(0, game.X)
		board.Set(1, game.X)
		board.Set(4, game.O)

		move, err := ai.GetMove(board)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if move != 2 {
			t.Errorf("expected AI to block at position 2, got %d", move)
		}
	})

	t.Run("takes winning move when available", func(t *testing.T) {
		// X _ X
		// O O _
		// _ _ _
		// AI should win at position 5
		board := game.NewBoard()
		board.Set(0, game.X)
		board.Set(2, game.X)
		board.Set(3, game.O)
		board.Set(4, game.O)

		move, err := ai.GetMove(board)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if move != 5 {
			t.Errorf("expected AI to win at position 5, got %d", move)
		}
	})

	t.Run("prioritizes immediate win over blocking", func(t *testing.T) {
		// X X _
		// O O _
		// _ _ _
		// AI has winning move at 5, and X has winning move at 2
		// AI should take its own win
		board := game.NewBoard()
		board.Set(0, game.X)
		board.Set(1, game.X)
		board.Set(3, game.O)
		board.Set(4, game.O)

		move, err := ai.GetMove(board)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if move != 5 {
			t.Errorf("expected AI to take winning move at position 5, got %d", move)
		}
	})

	t.Run("chooses optimal opening move on empty board", func(t *testing.T) {
		// On empty board, center (4) or corners (0,2,6,8) are optimal
		board := game.NewBoard()

		move, err := ai.GetMove(board)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Center or corners are optimal opening moves in tic-tac-toe
		optimalMoves := []int{0, 2, 4, 6, 8}
		isOptimal := false
		for _, optimal := range optimalMoves {
			if move == optimal {
				isOptimal = true
				break
			}
		}
		if !isOptimal {
			t.Errorf("expected AI to choose optimal move (0,2,4,6,8), got %d", move)
		}
	})

	t.Run("forces draw from losing position", func(t *testing.T) {
		// X _ _
		// _ O _
		// _ _ _
		// From this position, AI can force a draw
		board := game.NewBoard()
		board.Set(0, game.X)
		board.Set(4, game.O)

		// Play out the game - AI should never lose
		for !board.IsGameOver() {
			// AI move
			if len(board.GetEmptyPositions()) > 0 {
				move, err := ai.GetMove(board)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				board.Set(move, game.O)
			}

			if board.IsGameOver() {
				break
			}

			// Make suboptimal X move (random)
			emptyPos := board.GetEmptyPositions()
			if len(emptyPos) > 0 {
				board.Set(emptyPos[0], game.X)
			}
		}

		winner := board.CheckWinner()
		if winner == game.X {
			t.Error("AI should never lose against any opponent")
		}
	})
}

func TestHardAI(t *testing.T) {
	ai := &HardAI{depth: 4}

	t.Run("blocks winning move", func(t *testing.T) {
		// X X _
		// _ O _
		// _ _ _
		board := game.NewBoard()
		board.Set(0, game.X)
		board.Set(1, game.X)
		board.Set(4, game.O)

		move, err := ai.GetMove(board)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if move != 2 {
			t.Errorf("expected AI to block at position 2, got %d", move)
		}
	})

	t.Run("takes winning move", func(t *testing.T) {
		// X _ X
		// O O _
		// _ _ _
		board := game.NewBoard()
		board.Set(0, game.X)
		board.Set(2, game.X)
		board.Set(3, game.O)
		board.Set(4, game.O)

		move, err := ai.GetMove(board)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if move != 5 {
			t.Errorf("expected AI to win at position 5, got %d", move)
		}
	})
}

func TestImpossibleAI_IsUnbeatable(t *testing.T) {
	ai := &ImpossibleAI{}

	t.Run("achieves draw or win against any opponent", func(t *testing.T) {
		// Test that AI never loses across multiple game scenarios
		// Testing all possible first moves by opponent
		firstMoves := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}

		for _, firstMove := range firstMoves {
			board := game.NewBoard()
			board.Set(firstMove, game.X) // Opponent moves first

			// Play out the game with AI playing optimally as O
			for !board.IsGameOver() {
				// AI's turn (O)
				move, err := ai.GetMove(board)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				board.Set(move, game.O)

				if board.IsGameOver() {
					break
				}

				// Opponent's turn - play optimally
				// Use greedy strategy: take win, block loss, or random
				emptyPos := board.GetEmptyPositions()
				moved := false

				// Check if X can win
				for _, pos := range emptyPos {
					testBoard := board.Clone()
					testBoard.Set(pos, game.X)
					if testBoard.CheckWinner() == game.X {
						board.Set(pos, game.X)
						moved = true
						break
					}
				}

				if !moved && len(emptyPos) > 0 {
					// Check if O can win and block it
					for _, pos := range emptyPos {
						testBoard := board.Clone()
						testBoard.Set(pos, game.O)
						if testBoard.CheckWinner() == game.O {
							board.Set(pos, game.X)
							moved = true
							break
						}
					}
				}

				if !moved && len(emptyPos) > 0 {
					// Just take first available
					board.Set(emptyPos[0], game.X)
				}
			}

			winner := board.CheckWinner()
			if winner == game.X {
				t.Errorf("AI lost when opponent started at position %d", firstMove)
			}
		}
	})

	t.Run("never loses against random opponent", func(t *testing.T) {
		easyAI := &EasyAI{}

		// Run multiple games
		for i := 0; i < 10; i++ {
			board := game.NewBoard()

			for !board.IsGameOver() {
				// Random player (X) goes first
				emptyPos := board.GetEmptyPositions()
				if len(emptyPos) > 0 {
					move, _ := easyAI.GetMove(board)
					board.Set(move, game.X)
				}

				if board.IsGameOver() {
					break
				}

				// Impossible AI (O)
				move, err := ai.GetMove(board)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				board.Set(move, game.O)
			}

			winner := board.CheckWinner()
			if winner == game.X {
				t.Errorf("Game %d: ImpossibleAI should never lose, but X won", i)
			}
		}
	})

	t.Run("returns error on full board", func(t *testing.T) {
		board := game.NewBoard()
		// Fill board
		for i := 0; i < 9; i++ {
			if i%2 == 0 {
				board.Set(i, game.X)
			} else {
				board.Set(i, game.O)
			}
		}

		_, err := ai.GetMove(board)
		if err == nil {
			t.Error("expected error on full board")
		}
	})
}

func TestAIPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping performance test in short mode")
	}

	tests := []struct {
		name    string
		ai      AIStrategy
		maxTime time.Duration
	}{
		{
			name:    "EasyAI completes quickly",
			ai:      &EasyAI{},
			maxTime: 10 * time.Millisecond,
		},
		{
			name:    "MediumAI completes quickly",
			ai:      &MediumAI{},
			maxTime: 100 * time.Millisecond,
		},
		{
			name:    "HardAI completes within 1 second",
			ai:      &HardAI{depth: 4},
			maxTime: 1 * time.Second,
		},
		{
			name:    "ImpossibleAI completes within 1 second",
			ai:      &ImpossibleAI{},
			maxTime: 1 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := game.NewBoard()

			start := time.Now()
			_, err := tt.ai.GetMove(board)
			elapsed := time.Since(start)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if elapsed > tt.maxTime {
				t.Errorf("AI took %v, expected less than %v", elapsed, tt.maxTime)
			}
		})
	}

	t.Run("ImpossibleAI performance on various board states", func(t *testing.T) {
		ai := &ImpossibleAI{}

		testBoards := []struct {
			name  string
			setup func(*game.Board)
		}{
			{
				name:  "empty board",
				setup: func(b *game.Board) {},
			},
			{
				name: "early game (2 moves)",
				setup: func(b *game.Board) {
					b.Set(0, game.X)
					b.Set(4, game.O)
				},
			},
			{
				name: "mid game (5 moves)",
				setup: func(b *game.Board) {
					b.Set(0, game.X)
					b.Set(4, game.O)
					b.Set(1, game.X)
					b.Set(3, game.O)
					b.Set(8, game.X)
				},
			},
			{
				name: "late game (7 moves)",
				setup: func(b *game.Board) {
					b.Set(0, game.X)
					b.Set(4, game.O)
					b.Set(1, game.X)
					b.Set(3, game.O)
					b.Set(8, game.X)
					b.Set(2, game.O)
					b.Set(6, game.X)
				},
			},
		}

		for _, tc := range testBoards {
			t.Run(tc.name, func(t *testing.T) {
				board := game.NewBoard()
				tc.setup(board)

				start := time.Now()
				_, err := ai.GetMove(board)
				elapsed := time.Since(start)

				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if elapsed > 1*time.Second {
					t.Errorf("%s: AI took %v, expected less than 1s", tc.name, elapsed)
				}
			})
		}
	})
}
