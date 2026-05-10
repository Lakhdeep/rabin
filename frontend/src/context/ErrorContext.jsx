import { createContext, useContext, useState, useCallback } from 'react';
import ErrorToast from '../components/ErrorToast';

const ErrorContext = createContext();

export const useError = () => {
  const context = useContext(ErrorContext);
  if (!context) {
    throw new Error('useError must be used within an ErrorProvider');
  }
  return context;
};

export const ErrorProvider = ({ children }) => {
  const [error, setError] = useState(null);

  const showError = useCallback((message) => {
    // Handle different error formats
    let errorMessage = message;

    if (typeof message === 'object' && message !== null) {
      // API error object
      if (message.response?.data?.error) {
        errorMessage = message.response.data.error;
      } else if (message.message) {
        errorMessage = message.message;
      } else {
        errorMessage = 'An unexpected error occurred';
      }
    }

    // Network error
    if (errorMessage.includes('Network Error') || errorMessage.includes('Failed to fetch')) {
      errorMessage = 'Connection error. Please check your internet connection and try again.';
    }

    setError(errorMessage);
  }, []);

  const clearError = useCallback(() => {
    setError(null);
  }, []);

  return (
    <ErrorContext.Provider value={{ showError, clearError }}>
      {children}
      {error && (
        <ErrorToast 
          message={error} 
          onDismiss={clearError}
          autoDismiss={true}
          duration={5000}
        />
      )}
    </ErrorContext.Provider>
  );
};
