import { createContext, useContext, useState, useEffect } from 'react';
import * as authService from '../services/auth';
import * as storage from '../services/storage';

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Check for existing token on app load
  useEffect(() => {
    const initAuth = async () => {
      try {
        const token = storage.getToken();
        const storedUser = storage.getUser();

        if (token && storedUser) {
          // Verify token is still valid by fetching current user
          try {
            const userData = await authService.getCurrentUser();
            setUser(userData);
          } catch (err) {
            // Token is invalid, clear storage
            storage.clearAuth();
            setUser(null);
          }
        }
      } catch (err) {
        console.error('Auth initialization error:', err);
      } finally {
        setLoading(false);
      }
    };

    initAuth();
  }, []);

  const register = async (userData) => {
    try {
      setLoading(true);
      setError(null);
      const response = await authService.register(userData);
      return response;
    } catch (err) {
      const errorMessage = err.response?.data?.error || 'Registration failed';
      setError(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const login = async (credentials) => {
    try {
      setLoading(true);
      setError(null);
      const response = await authService.login(credentials);
      
      // Store token and user data
      if (response.token) {
        storage.saveToken(response.token);
        
        // Store user data if available, otherwise fetch it
        if (response.user) {
          storage.saveUser(response.user);
          setUser(response.user);
        } else {
          const userData = await authService.getCurrentUser();
          storage.saveUser(userData);
          setUser(userData);
        }
      }
      
      return response;
    } catch (err) {
      const errorMessage = err.response?.data?.error || 'Login failed';
      setError(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const logout = () => {
    authService.logout();
    setUser(null);
    setError(null);
  };

  const clearError = () => {
    setError(null);
  };

  const value = {
    user,
    loading,
    error,
    register,
    login,
    logout,
    clearError,
    isAuthenticated: !!user,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

// Custom hook to use auth context
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
