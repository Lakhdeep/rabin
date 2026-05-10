import { createContext, useContext, useState } from 'react';
import * as gameService from '../services/game';

const GameContext = createContext(null);

export const GameProvider = ({ children }) => {
  const [gameId, setGameId] = useState(null);
  const [board, setBoard] = useState(Array(9).fill(''));
  const [difficulty, setDifficulty] = useState('medium');
  const [currentTurn, setCurrentTurn] = useState('X');
  const [status, setStatus] = useState('idle'); // idle, active, won, lost, draw
  const [result, setResult] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const createGame = async (selectedDifficulty) => {
    try {
      setLoading(true);
      setError(null);

      const response = await gameService.createGame({ 
        difficulty: selectedDifficulty 
      });

      setGameId(response.id);
      setBoard(response.board || Array(9).fill(''));
      setDifficulty(selectedDifficulty);
      setCurrentTurn('X');
      setStatus('active');
      setResult(null);

      return response;
    } catch (err) {
      const errorMessage = err.response?.data?.error || 'Failed to create game';
      setError(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const makeMove = async (position) => {
    if (!gameId || status !== 'active' || loading) {
      return;
    }

    try {
      setLoading(true);
      setError(null);

      const response = await gameService.makeMove(gameId, { position });

      // Update board with new state
      setBoard(response.board);

      // Check if game ended
      if (response.status !== 'active') {
        setStatus(response.status);
        setResult(response.result);
      } else {
        setCurrentTurn(response.current_turn || 'X');
      }

      return response;
    } catch (err) {
      const errorMessage = err.response?.data?.error || 'Failed to make move';
      setError(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const resetGame = () => {
    setGameId(null);
    setBoard(Array(9).fill(''));
    setDifficulty('medium');
    setCurrentTurn('X');
    setStatus('idle');
    setResult(null);
    setError(null);
    setLoading(false);
  };

  const clearError = () => {
    setError(null);
  };

  const value = {
    gameId,
    board,
    difficulty,
    currentTurn,
    status,
    result,
    loading,
    error,
    createGame,
    makeMove,
    resetGame,
    clearError,
    setDifficulty,
  };

  return <GameContext.Provider value={value}>{children}</GameContext.Provider>;
};

// Custom hook to use game context
export const useGame = () => {
  const context = useContext(GameContext);
  if (!context) {
    throw new Error('useGame must be used within a GameProvider');
  }
  return context;
};
