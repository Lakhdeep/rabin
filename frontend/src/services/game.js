import api from './api';

/**
 * Create a new game
 * @param {Object} gameData - Game creation data
 * @param {string} gameData.difficulty - AI difficulty (easy, medium, hard, impossible)
 * @returns {Promise} Created game data
 */
export const createGame = async (gameData) => {
  const response = await api.post('/games', gameData);
  return response.data;
};

/**
 * Get game by ID
 * @param {number} gameId - Game ID
 * @returns {Promise} Game data
 */
export const getGame = async (gameId) => {
  const response = await api.get(`/games/${gameId}`);
  return response.data;
};

/**
 * Make a move in a game
 * @param {number} gameId - Game ID
 * @param {Object} moveData - Move data
 * @param {number} moveData.position - Position on board (0-8)
 * @returns {Promise} Updated game state with AI response
 */
export const makeMove = async (gameId, moveData) => {
  const response = await api.post(`/games/${gameId}/move`, moveData);
  return response.data;
};

/**
 * Get all games for current user
 * @returns {Promise} List of games
 */
export const getUserGames = async () => {
  const response = await api.get('/games');
  return response.data;
};
