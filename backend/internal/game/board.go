package game

import "fmt"

// CellValue represents a cell on the board
type CellValue string

const (
	Empty CellValue = ""
	X     CellValue = "X"
	O     CellValue = "O"
)

// Board represents a 3x3 tic-tac-toe board
// Positions are numbered 0-8:
//
//	0 | 1 | 2
//	---------
//	3 | 4 | 5
//	---------
//	6 | 7 | 8
type Board struct {
	cells [9]CellValue
}

// NewBoard creates a new empty board
func NewBoard() *Board {
	return &Board{}
}

// Get returns the value at the given position (0-8)
func (b *Board) Get(position int) (CellValue, error) {
	if position < 0 || position > 8 {
		return Empty, fmt.Errorf("position must be between 0 and 8")
	}
	return b.cells[position], nil
}

// Set places a mark at the given position
func (b *Board) Set(position int, value CellValue) error {
	if position < 0 || position > 8 {
		return fmt.Errorf("position must be between 0 and 8")
	}
	b.cells[position] = value
	return nil
}

// IsEmpty checks if a position is empty
func (b *Board) IsEmpty(position int) bool {
	if position < 0 || position > 8 {
		return false
	}
	return b.cells[position] == Empty
}

// IsFull checks if the board has no empty positions
func (b *Board) IsFull() bool {
	for _, cell := range b.cells {
		if cell == Empty {
			return false
		}
	}
	return true
}

// GetEmptyPositions returns a slice of all empty position indices
func (b *Board) GetEmptyPositions() []int {
	var positions []int
	for i, cell := range b.cells {
		if cell == Empty {
			positions = append(positions, i)
		}
	}
	return positions
}

// Clone creates a deep copy of the board
func (b *Board) Clone() *Board {
	clone := &Board{}
	copy(clone.cells[:], b.cells[:])
	return clone
}

// ToSlice returns the board as a slice
func (b *Board) ToSlice() []CellValue {
	result := make([]CellValue, 9)
	copy(result, b.cells[:])
	return result
}

// FromSlice initializes the board from a slice
func (b *Board) FromSlice(slice []CellValue) error {
	if len(slice) != 9 {
		return fmt.Errorf("slice must have exactly 9 elements")
	}
	copy(b.cells[:], slice)
	return nil
}

// winningCombinations defines all possible winning combinations
var winningCombinations = [][3]int{
	// Rows
	{0, 1, 2},
	{3, 4, 5},
	{6, 7, 8},
	// Columns
	{0, 3, 6},
	{1, 4, 7},
	{2, 5, 8},
	// Diagonals
	{0, 4, 8},
	{2, 4, 6},
}

// CheckWinner checks if there's a winner and returns the winning player
// Returns Empty if no winner
func (b *Board) CheckWinner() CellValue {
	for _, combo := range winningCombinations {
		a, b1, c := b.cells[combo[0]], b.cells[combo[1]], b.cells[combo[2]]
		if a != Empty && a == b1 && b1 == c {
			return a
		}
	}
	return Empty
}

// IsGameOver checks if the game is over (winner or draw)
func (b *Board) IsGameOver() bool {
	return b.CheckWinner() != Empty || b.IsFull()
}

// ValidateMove checks if a move is valid
// Returns an error if the move is invalid
func (b *Board) ValidateMove(position int) error {
	// Check if position is in valid range
	if position < 0 || position > 8 {
		return fmt.Errorf("invalid position: must be between 0 and 8")
	}

	// Check if game is already over
	if b.IsGameOver() {
		return fmt.Errorf("game already completed")
	}

	// Check if position is already occupied
	if !b.IsEmpty(position) {
		return fmt.Errorf("position already occupied")
	}

	return nil
}
