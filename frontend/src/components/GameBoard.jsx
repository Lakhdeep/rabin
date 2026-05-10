import { useState, useEffect } from 'react';
import './GameBoard.css';

const GameBoard = ({ board, onCellClick, disabled, status }) => {
  const [winningLine, setWinningLine] = useState(null);

  // Winning combinations
  const winningCombinations = [
    [0, 1, 2], // Top row
    [3, 4, 5], // Middle row
    [6, 7, 8], // Bottom row
    [0, 3, 6], // Left column
    [1, 4, 7], // Middle column
    [2, 5, 8], // Right column
    [0, 4, 8], // Diagonal top-left to bottom-right
    [2, 4, 6], // Diagonal top-right to bottom-left
  ];

  // Check for winning line when game ends
  useEffect(() => {
    if (status === 'won' || status === 'lost') {
      for (const combo of winningCombinations) {
        const [a, b, c] = combo;
        if (board[a] && board[a] === board[b] && board[a] === board[c]) {
          setWinningLine(combo);
          break;
        }
      }
    } else {
      setWinningLine(null);
    }
  }, [status, board]);

  const handleCellClick = (index) => {
    // Don't allow clicks if disabled, cell is occupied, or game is over
    if (disabled || board[index] || status !== 'active') {
      return;
    }
    onCellClick(index);
  };

  const getCellClassName = (index) => {
    let className = 'cell';
    
    // Add occupied class if cell has a value
    if (board[index]) {
      className += ' occupied';
    }

    // Add player-specific class
    if (board[index] === 'X') {
      className += ' cell-x';
    } else if (board[index] === 'O') {
      className += ' cell-o';
    }

    // Add winning class if part of winning line
    if (winningLine && winningLine.includes(index)) {
      className += ' winning';
    }

    // Add disabled class
    if (disabled || status !== 'active') {
      className += ' disabled';
    }

    return className;
  };

  return (
    <div className="game-board">
      {board.map((value, index) => (
        <button
          key={index}
          className={getCellClassName(index)}
          onClick={() => handleCellClick(index)}
          disabled={disabled || board[index] || status !== 'active'}
          aria-label={`Cell ${index + 1}${value ? `, occupied by ${value}` : ', empty'}`}
        >
          <span className="cell-value">{value}</span>
        </button>
      ))}
    </div>
  );
};

export default GameBoard;
