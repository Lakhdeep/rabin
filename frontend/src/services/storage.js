/**
 * Save JWT token to localStorage
 * @param {string} token - JWT token
 */
export const saveToken = (token) => {
  localStorage.setItem('token', token);
};

/**
 * Get JWT token from localStorage
 * @returns {string|null} JWT token or null if not found
 */
export const getToken = () => {
  return localStorage.getItem('token');
};

/**
 * Clear JWT token from localStorage
 */
export const clearToken = () => {
  localStorage.removeItem('token');
};

/**
 * Save user data to localStorage
 * @param {Object} user - User data
 */
export const saveUser = (user) => {
  localStorage.setItem('user', JSON.stringify(user));
};

/**
 * Get user data from localStorage
 * @returns {Object|null} User data or null if not found
 */
export const getUser = () => {
  const user = localStorage.getItem('user');
  return user ? JSON.parse(user) : null;
};

/**
 * Clear user data from localStorage
 */
export const clearUser = () => {
  localStorage.removeItem('user');
};

/**
 * Clear all auth data from localStorage
 */
export const clearAuth = () => {
  clearToken();
  clearUser();
};
