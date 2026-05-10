package ai

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/rabin/tictactoe/internal/game"
)

// AIStrategy defines the interface for AI move selection
type AIStrategy interface {
	GetMove(board *game.Board) (int, error)
}

// GetMoveWithTimeout wraps an AI strategy with timeout protection
func GetMoveWithTimeout(ai AIStrategy, board *game.Board, timeout time.Duration) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	type result struct {
		move int
		err  error
	}

	resultChan := make(chan result, 1)

	go func() {
		move, err := ai.GetMove(board)
		resultChan <- result{move, err}
	}()

	select {
	case res := <-resultChan:
		return res.move, res.err
	case <-ctx.Done():
		// Timeout occurred, fallback to random move
		easy := &EasyAI{}
		return easy.GetMove(board)
	}
}

// NewAI creates an AI strategy based on difficulty level
func NewAI(difficulty game.Difficulty) (AIStrategy, error) {
	switch difficulty {
	case game.DifficultyEasy:
		return &EasyAI{}, nil
	case game.DifficultyMedium:
		return &MediumAI{}, nil
	case game.DifficultyHard:
		return &HardAI{depth: 4}, nil
	case game.DifficultyImpossible:
		return &ImpossibleAI{}, nil
	default:
		return nil, fmt.Errorf("invalid difficulty: %s", difficulty)
	}
}

// EasyAI implements random move selection
type EasyAI struct{}

// GetMove returns a random valid move
func (ai *EasyAI) GetMove(board *game.Board) (int, error) {
	emptyPositions := board.GetEmptyPositions()
	if len(emptyPositions) == 0 {
		return -1, fmt.Errorf("no valid moves available")
	}

	// Use time-based seed for better randomness
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return emptyPositions[r.Intn(len(emptyPositions))], nil
}

// MediumAI implements 50% optimal, 50% random strategy
type MediumAI struct{}

// GetMove returns either a minimax move or random move with 50% probability each
func (ai *MediumAI) GetMove(board *game.Board) (int, error) {
	emptyPositions := board.GetEmptyPositions()
	if len(emptyPositions) == 0 {
		return -1, fmt.Errorf("no valid moves available")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 50% chance to use minimax, 50% chance for random
	if r.Float64() < 0.5 {
		// Use optimal strategy
		impossible := &ImpossibleAI{}
		return impossible.GetMove(board)
	}

	// Use random strategy
	easy := &EasyAI{}
	return easy.GetMove(board)
}

// HardAI implements depth-limited minimax
type HardAI struct {
	depth int
}

// GetMove returns the best move using depth-limited minimax
func (ai *HardAI) GetMove(board *game.Board) (int, error) {
	emptyPositions := board.GetEmptyPositions()
	if len(emptyPositions) == 0 {
		return -1, fmt.Errorf("no valid moves available")
	}

	bestMove := -1
	bestScore := math.MinInt32

	for _, pos := range emptyPositions {
		// Try this move
		boardClone := board.Clone()
		boardClone.Set(pos, game.O)

		// Get score with depth limit
		score := ai.minimaxWithDepth(boardClone, 0, ai.depth, false, math.MinInt32, math.MaxInt32)

		if score > bestScore {
			bestScore = score
			bestMove = pos
		}
	}

	if bestMove == -1 {
		return -1, fmt.Errorf("could not find valid move")
	}

	return bestMove, nil
}

// minimaxWithDepth implements minimax with alpha-beta pruning and depth limit
func (ai *HardAI) minimaxWithDepth(board *game.Board, depth, maxDepth int, isMaximizing bool, alpha, beta int) int {
	// Check terminal states
	winner := board.CheckWinner()
	if winner == game.O {
		return 10 - depth // Prefer faster wins
	}
	if winner == game.X {
		return depth - 10 // Prefer slower losses
	}
	if board.IsFull() {
		return 0 // Draw
	}

	// Check depth limit
	if depth >= maxDepth {
		return 0 // Heuristic evaluation (neutral)
	}

	emptyPositions := board.GetEmptyPositions()

	if isMaximizing {
		// AI (O) turn
		maxScore := math.MinInt32
		for _, pos := range emptyPositions {
			boardClone := board.Clone()
			boardClone.Set(pos, game.O)
			score := ai.minimaxWithDepth(boardClone, depth+1, maxDepth, false, alpha, beta)
			maxScore = max(maxScore, score)
			alpha = max(alpha, score)
			if beta <= alpha {
				break // Beta cutoff
			}
		}
		return maxScore
	} else {
		// Player (X) turn
		minScore := math.MaxInt32
		for _, pos := range emptyPositions {
			boardClone := board.Clone()
			boardClone.Set(pos, game.X)
			score := ai.minimaxWithDepth(boardClone, depth+1, maxDepth, true, alpha, beta)
			minScore = min(minScore, score)
			beta = min(beta, score)
			if beta <= alpha {
				break // Alpha cutoff
			}
		}
		return minScore
	}
}

// ImpossibleAI implements full minimax with alpha-beta pruning
type ImpossibleAI struct{}

// GetMove returns the optimal move using full minimax
func (ai *ImpossibleAI) GetMove(board *game.Board) (int, error) {
	emptyPositions := board.GetEmptyPositions()
	if len(emptyPositions) == 0 {
		return -1, fmt.Errorf("no valid moves available")
	}

	bestMove := -1
	bestScore := math.MinInt32

	for _, pos := range emptyPositions {
		// Try this move
		boardClone := board.Clone()
		boardClone.Set(pos, game.O)

		// Get score
		score := ai.minimax(boardClone, 0, false, math.MinInt32, math.MaxInt32)

		if score > bestScore {
			bestScore = score
			bestMove = pos
		}
	}

	if bestMove == -1 {
		return -1, fmt.Errorf("could not find valid move")
	}

	return bestMove, nil
}

// minimax implements the minimax algorithm with alpha-beta pruning
func (ai *ImpossibleAI) minimax(board *game.Board, depth int, isMaximizing bool, alpha, beta int) int {
	// Check terminal states
	winner := board.CheckWinner()
	if winner == game.O {
		return 10 - depth // Prefer faster wins
	}
	if winner == game.X {
		return depth - 10 // Prefer slower losses
	}
	if board.IsFull() {
		return 0 // Draw
	}

	emptyPositions := board.GetEmptyPositions()

	if isMaximizing {
		// AI (O) turn
		maxScore := math.MinInt32
		for _, pos := range emptyPositions {
			boardClone := board.Clone()
			boardClone.Set(pos, game.O)
			score := ai.minimax(boardClone, depth+1, false, alpha, beta)
			maxScore = max(maxScore, score)
			alpha = max(alpha, score)
			if beta <= alpha {
				break // Beta cutoff
			}
		}
		return maxScore
	} else {
		// Player (X) turn
		minScore := math.MaxInt32
		for _, pos := range emptyPositions {
			boardClone := board.Clone()
			boardClone.Set(pos, game.X)
			score := ai.minimax(boardClone, depth+1, true, alpha, beta)
			minScore = min(minScore, score)
			beta = min(beta, score)
			if beta <= alpha {
				break // Alpha cutoff
			}
		}
		return minScore
	}
}

// Helper functions for min/max
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
