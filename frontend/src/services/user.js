import api from './api';

/**
 * Get user statistics by user ID
 * @param {number} userId - User ID
 * @returns {Promise} User statistics including win rate
 */
export const getUserStats = async (userId) => {
  const response = await api.get(`/users/${userId}/stats`);
  return response.data;
};
