import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { GameProvider, useGame } from '../context/GameContext';
import GameBoard from '../components/GameBoard';
import './Game.css';

const GameContent = () => {
  const navigate = useNavigate();
  const {
    gameId,
    board,
    difficulty,
    currentTurn,
    status,
    loading,
    error,
    createGame,
    makeMove,
    resetGame,
    setDifficulty,
    clearError,
  } = useGame();

  const [selectedDifficulty, setSelectedDifficulty] = useState('medium');

  const handleDifficultyChange = (newDifficulty) => {
    setSelectedDifficulty(newDifficulty);
    setDifficulty(newDifficulty);
  };

  const handleStartGame = async () => {
    try {
      await createGame(selectedDifficulty);
    } catch (err) {
      console.error('Failed to start game:', err);
    }
  };

  const handleCellClick = async (position) => {
    try {
      await makeMove(position);
    } catch (err) {
      console.error('Failed to make move:', err);
    }
  };

  const handlePlayAgain = () => {
    resetGame();
  };

  const handleBackToDashboard = () => {
    resetGame();
    navigate('/dashboard');
  };

  const getStatusMessage = () => {
    if (loading) {
      return 'AI is thinking...';
    }

    if (status === 'won') {
      return '🎉 You won!';
    }

    if (status === 'lost') {
      return '😔 You lost!';
    }

    if (status === 'draw') {
      return '🤝 It is a draw!';
    }

    if (status === 'active' && currentTurn === 'X') {
      return 'Your turn (X)';
    }

    if (status === 'active' && currentTurn === 'O') {
      return 'AI is thinking...';
    }

    return 'Select difficulty and start a game';
  };

  const getDifficultyInfo = (level) => {
    const info = {
      easy: 'Random moves - Perfect for beginners',
      medium: 'Mixed strategy - Good challenge',
      hard: 'Smart moves - Difficult opponent',
      impossible: 'Perfect play - Unbeatable AI',
    };
    return info[level] || '';
  };

  const isGameOver = status === 'won' || status === 'lost' || status === 'draw';

  return (
    <div className="game-container">
      <div className="game-content">
        <div className="game-header">
          <h1>Tic-Tac-Toe</h1>
          <button onClick={handleBackToDashboard} className="btn btn-secondary">
            ← Back to Dashboard
          </button>
        </div>

        <div className={`game-status ${isGameOver ? 'game-over' : ''}`}>
          <h2>{getStatusMessage()}</h2>
          {status === 'active' && (
            <div className="turn-indicator">
              <span className={currentTurn === 'X' ? 'active' : ''}>You (X)</span>
              <span className="vs">vs</span>
              <span className={currentTurn === 'O' ? 'active' : ''}>AI (O)</span>
            </div>
          )}
        </div>

        {error && (
          <div className="error-message">
            {error}
            <button onClick={clearError} className="close-btn">×</button>
          </div>
        )}

        {!gameId && (
          <div className="game-setup card">
            <h3>Choose Difficulty</h3>
            <div className="difficulty-selector">
              {['easy', 'medium', 'hard', 'impossible'].map((level) => (
                <button
                  key={level}
                  onClick={() => handleDifficultyChange(level)}
                  className={`difficulty-btn ${selectedDifficulty === level ? 'selected' : ''}`}
                  title={getDifficultyInfo(level)}
                >
                  <span className="difficulty-name">{level}</span>
                  <span className="difficulty-desc">{getDifficultyInfo(level)}</span>
                </button>
              ))}
            </div>
            <button
              onClick={handleStartGame}
              className="btn btn-primary btn-large"
              disabled={loading}
            >
              {loading ? 'Starting...' : 'Start Game'}
            </button>
          </div>
        )}

        {gameId && (
          <div className="game-board-section">
            {loading && (
              <div className="loading-overlay">
                <div className="spinner"></div>
              </div>
            )}
            <GameBoard
              board={board}
              onCellClick={handleCellClick}
              disabled={loading || isGameOver}
              status={status}
            />
          </div>
        )}

        {isGameOver && (
          <div className="game-actions">
            <button onClick={handlePlayAgain} className="btn btn-primary">
              Play Again
            </button>
            <button onClick={handleBackToDashboard} className="btn btn-secondary">
              Back to Dashboard
            </button>
          </div>
        )}

        {gameId && (
          <div className="game-info">
            <div className="info-item">
              <span className="label">Difficulty:</span>
              <span className="value">{difficulty}</span>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

const Game = () => {
  return (
    <GameProvider>
      <GameContent />
    </GameProvider>
  );
};

export default Game;
