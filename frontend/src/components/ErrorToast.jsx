import { useEffect } from 'react';
import './ErrorToast.css';

const ErrorToast = ({ message, onDismiss, autoDismiss = true, duration = 5000 }) => {
  useEffect(() => {
    if (autoDismiss && onDismiss) {
      const timer = setTimeout(() => {
        onDismiss();
      }, duration);

      return () => clearTimeout(timer);
    }
  }, [autoDismiss, duration, onDismiss]);

  if (!message) return null;

  return (
    <div className="error-toast">
      <div className="error-toast-content">
        <div className="error-toast-icon">⚠️</div>
        <div className="error-toast-message">{message}</div>
        {onDismiss && (
          <button 
            className="error-toast-close" 
            onClick={onDismiss}
            aria-label="Dismiss error"
          >
            ×
          </button>
        )}
      </div>
    </div>
  );
};

export default ErrorToast;
