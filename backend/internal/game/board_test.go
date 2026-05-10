package game

import (
	"testing"
)

func TestNewBoard(t *testing.T) {
	board := NewBoard()

	// Test that all positions are empty
	for i := 0; i < 9; i++ {
		cell, err := board.Get(i)
		if err != nil {
			t.Errorf("unexpected error getting position %d: %v", i, err)
		}
		if cell != Empty {
			t.Errorf("expected position %d to be empty, got %v", i, cell)
		}
	}

	// Test that board is not full
	if board.IsFull() {
		t.Error("expected new board to not be full")
	}

	// Test that game is not over
	if board.IsGameOver() {
		t.Error("expected new board game to not be over")
	}

	// Test that there's no winner
	if winner := board.CheckWinner(); winner != Empty {
		t.Errorf("expected no winner on new board, got %v", winner)
	}
}

func TestBoardSetGet(t *testing.T) {
	board := NewBoard()

	// Test setting and getting a value
	err := board.Set(0, X)
	if err != nil {
		t.Errorf("unexpected error setting position: %v", err)
	}

	cell, err := board.Get(0)
	if err != nil {
		t.Errorf("unexpected error getting position: %v", err)
	}

	if cell != X {
		t.Errorf("expected X at position 0, got %v", cell)
	}

	// Test invalid position
	err = board.Set(10, X)
	if err == nil {
		t.Error("expected error when setting invalid position 10")
	}

	_, err = board.Get(-1)
	if err == nil {
		t.Error("expected error when getting invalid position -1")
	}
}

func TestBoardIsEmpty(t *testing.T) {
	board := NewBoard()

	// Test empty position
	if !board.IsEmpty(0) {
		t.Error("expected position 0 to be empty")
	}

	// Set a value
	board.Set(0, X)

	// Test occupied position
	if board.IsEmpty(0) {
		t.Error("expected position 0 to not be empty")
	}

	// Test invalid position
	if board.IsEmpty(-1) {
		t.Error("expected invalid position to return false")
	}
}

func TestBoardGetEmptyPositions(t *testing.T) {
	board := NewBoard()

	// Test all positions empty
	empty := board.GetEmptyPositions()
	if len(empty) != 9 {
		t.Errorf("expected 9 empty positions, got %d", len(empty))
	}

	// Fill some positions
	board.Set(0, X)
	board.Set(4, O)
	board.Set(8, X)

	empty = board.GetEmptyPositions()
	if len(empty) != 6 {
		t.Errorf("expected 6 empty positions, got %d", len(empty))
	}

	// Verify correct positions are empty
	expectedEmpty := []int{1, 2, 3, 5, 6, 7}
	for i, pos := range empty {
		if pos != expectedEmpty[i] {
			t.Errorf("expected position %d at index %d, got %d", expectedEmpty[i], i, pos)
		}
	}
}

func TestBoardIsFull(t *testing.T) {
	board := NewBoard()

	if board.IsFull() {
		t.Error("expected new board to not be full")
	}

	// Fill all positions
	board.Set(0, X)
	board.Set(1, O)
	board.Set(2, X)
	board.Set(3, O)
	board.Set(4, X)
	board.Set(5, O)
	board.Set(6, X)
	board.Set(7, O)
	board.Set(8, X)

	if !board.IsFull() {
		t.Error("expected board to be full")
	}
}

func TestBoardClone(t *testing.T) {
	board := NewBoard()
	board.Set(0, X)
	board.Set(4, O)

	clone := board.Clone()

	// Verify clone has same values
	for i := 0; i < 9; i++ {
		orig, _ := board.Get(i)
		cloned, _ := clone.Get(i)
		if orig != cloned {
			t.Errorf("position %d: expected %v, got %v", i, orig, cloned)
		}
	}

	// Verify modifying clone doesn't affect original
	clone.Set(1, X)
	orig, _ := board.Get(1)
	if orig != Empty {
		t.Error("modifying clone affected original board")
	}
}

func TestBoardToSliceFromSlice(t *testing.T) {
	board := NewBoard()
	board.Set(0, X)
	board.Set(4, O)
	board.Set(8, X)

	// Test ToSlice
	slice := board.ToSlice()
	if len(slice) != 9 {
		t.Errorf("expected slice length 9, got %d", len(slice))
	}
	if slice[0] != X || slice[4] != O || slice[8] != X {
		t.Error("ToSlice returned incorrect values")
	}

	// Test FromSlice
	newBoard := NewBoard()
	err := newBoard.FromSlice(slice)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	for i := 0; i < 9; i++ {
		orig, _ := board.Get(i)
		new, _ := newBoard.Get(i)
		if orig != new {
			t.Errorf("position %d: expected %v, got %v", i, orig, new)
		}
	}

	// Test FromSlice with invalid length
	err = newBoard.FromSlice([]CellValue{X, O})
	if err == nil {
		t.Error("expected error with invalid slice length")
	}
}

func TestValidateMove(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*Board)
		position  int
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid move on empty board",
			setup:     func(b *Board) {},
			position:  0,
			wantError: false,
		},
		{
			name:      "valid move on partially filled board",
			setup:     func(b *Board) { b.Set(0, X) },
			position:  1,
			wantError: false,
		},
		{
			name:      "invalid position - negative",
			setup:     func(b *Board) {},
			position:  -1,
			wantError: true,
			errorMsg:  "invalid position",
		},
		{
			name:      "invalid position - too large",
			setup:     func(b *Board) {},
			position:  9,
			wantError: true,
			errorMsg:  "invalid position",
		},
		{
			name:      "position already occupied",
			setup:     func(b *Board) { b.Set(4, X) },
			position:  4,
			wantError: true,
			errorMsg:  "position already occupied",
		},
		{
			name: "game already completed - winner",
			setup: func(b *Board) {
				// Create winning position for X
				b.Set(0, X)
				b.Set(1, X)
				b.Set(2, X)
			},
			position:  3,
			wantError: true,
			errorMsg:  "game already completed",
		},
		{
			name: "game already completed - draw",
			setup: func(b *Board) {
				// Create a draw position
				b.Set(0, X) // X O X
				b.Set(1, O) // O X X
				b.Set(2, X) // O X O
				b.Set(3, O)
				b.Set(4, X)
				b.Set(5, X)
				b.Set(6, O)
				b.Set(7, X)
				b.Set(8, O)
			},
			position:  0, // Already occupied anyway
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := NewBoard()
			tt.setup(board)

			err := board.ValidateMove(tt.position)

			if tt.wantError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.wantError && err != nil && tt.errorMsg != "" {
				if !contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			len(s) > len(substr)+1 && findSubstr(s, substr)))
}

func findSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestCheckWinner(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*Board)
		expected CellValue
		desc     string
	}{
		{
			name:     "no winner - empty board",
			setup:    func(b *Board) {},
			expected: Empty,
			desc:     "empty board should have no winner",
		},
		{
			name: "no winner - partial board",
			setup: func(b *Board) {
				b.Set(0, X)
				b.Set(4, O)
			},
			expected: Empty,
			desc:     "partially filled board with no winning combination",
		},
		{
			name: "row 0 win - X",
			setup: func(b *Board) {
				b.Set(0, X)
				b.Set(1, X)
				b.Set(2, X)
			},
			expected: X,
			desc:     "top row (0-1-2) should win for X",
		},
		{
			name: "row 1 win - O",
			setup: func(b *Board) {
				b.Set(3, O)
				b.Set(4, O)
				b.Set(5, O)
			},
			expected: O,
			desc:     "middle row (3-4-5) should win for O",
		},
		{
			name: "row 2 win - X",
			setup: func(b *Board) {
				b.Set(6, X)
				b.Set(7, X)
				b.Set(8, X)
			},
			expected: X,
			desc:     "bottom row (6-7-8) should win for X",
		},
		{
			name: "column 0 win - O",
			setup: func(b *Board) {
				b.Set(0, O)
				b.Set(3, O)
				b.Set(6, O)
			},
			expected: O,
			desc:     "left column (0-3-6) should win for O",
		},
		{
			name: "column 1 win - X",
			setup: func(b *Board) {
				b.Set(1, X)
				b.Set(4, X)
				b.Set(7, X)
			},
			expected: X,
			desc:     "middle column (1-4-7) should win for X",
		},
		{
			name: "column 2 win - O",
			setup: func(b *Board) {
				b.Set(2, O)
				b.Set(5, O)
				b.Set(8, O)
			},
			expected: O,
			desc:     "right column (2-5-8) should win for O",
		},
		{
			name: "diagonal top-left to bottom-right - X",
			setup: func(b *Board) {
				b.Set(0, X)
				b.Set(4, X)
				b.Set(8, X)
			},
			expected: X,
			desc:     "diagonal (0-4-8) should win for X",
		},
		{
			name: "diagonal top-right to bottom-left - O",
			setup: func(b *Board) {
				b.Set(2, O)
				b.Set(4, O)
				b.Set(6, O)
			},
			expected: O,
			desc:     "diagonal (2-4-6) should win for O",
		},
		{
			name: "no winner - draw position",
			setup: func(b *Board) {
				// X O X
				// O X X
				// O X O
				b.Set(0, X)
				b.Set(1, O)
				b.Set(2, X)
				b.Set(3, O)
				b.Set(4, X)
				b.Set(5, X)
				b.Set(6, O)
				b.Set(7, X)
				b.Set(8, O)
			},
			expected: Empty,
			desc:     "full board with no winner should return Empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := NewBoard()
			tt.setup(board)

			winner := board.CheckWinner()
			if winner != tt.expected {
				t.Errorf("%s: expected %v, got %v", tt.desc, tt.expected, winner)
			}
		})
	}
}

func TestIsGameOver(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*Board)
		expected bool
	}{
		{
			name:     "not over - empty board",
			setup:    func(b *Board) {},
			expected: false,
		},
		{
			name: "not over - partial board",
			setup: func(b *Board) {
				b.Set(0, X)
				b.Set(4, O)
			},
			expected: false,
		},
		{
			name: "over - X wins",
			setup: func(b *Board) {
				b.Set(0, X)
				b.Set(1, X)
				b.Set(2, X)
			},
			expected: true,
		},
		{
			name: "over - O wins",
			setup: func(b *Board) {
				b.Set(0, O)
				b.Set(3, O)
				b.Set(6, O)
			},
			expected: true,
		},
		{
			name: "over - draw",
			setup: func(b *Board) {
				b.Set(0, X)
				b.Set(1, O)
				b.Set(2, X)
				b.Set(3, O)
				b.Set(4, X)
				b.Set(5, X)
				b.Set(6, O)
				b.Set(7, X)
				b.Set(8, O)
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := NewBoard()
			tt.setup(board)

			if got := board.IsGameOver(); got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}
