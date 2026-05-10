import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import './Register.css';

const Register = () => {
  const [formData, setFormData] = useState({
    email: '',
    username: '',
    password: '',
  });
  const [errors, setErrors] = useState({});
  const [success, setSuccess] = useState(false);
  const { register, loading, error: authError } = useAuth();
  const navigate = useNavigate();

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
    // Clear error for this field
    if (errors[name]) {
      setErrors(prev => ({
        ...prev,
        [name]: ''
      }));
    }
  };

  const validateEmail = (email) => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  };

  const validatePassword = (password) => {
    return password.length >= 8;
  };

  const getPasswordStrength = (password) => {
    if (password.length === 0) return null;
    if (password.length < 6) return 'weak';
    if (password.length < 8) return 'medium';
    
    let strength = 0;
    if (password.length >= 8) strength++;
    if (/[a-z]/.test(password)) strength++;
    if (/[A-Z]/.test(password)) strength++;
    if (/[0-9]/.test(password)) strength++;
    if (/[^a-zA-Z0-9]/.test(password)) strength++;

    if (strength >= 4) return 'strong';
    return 'medium';
  };

  const validateForm = () => {
    const newErrors = {};

    if (!formData.email) {
      newErrors.email = 'Email is required';
    } else if (!validateEmail(formData.email)) {
      newErrors.email = 'Invalid email format';
    }

    if (!formData.username) {
      newErrors.username = 'Username is required';
    } else if (formData.username.length < 3) {
      newErrors.username = 'Username must be at least 3 characters';
    } else if (formData.username.length > 20) {
      newErrors.username = 'Username must be less than 20 characters';
    }

    if (!formData.password) {
      newErrors.password = 'Password is required';
    } else if (!validatePassword(formData.password)) {
      newErrors.password = 'Password must be at least 8 characters';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    try {
      await register(formData);
      setSuccess(true);
      setTimeout(() => {
        navigate('/login', { state: { message: 'Registration successful! Please login.' } });
      }, 2000);
    } catch (err) {
      // Error is handled by AuthContext
      console.error('Registration failed:', err);
    }
  };

  const passwordStrength = getPasswordStrength(formData.password);

  return (
    <div className="register-container">
      <div className="register-card card">
        <h1>Create Account</h1>
        <p className="register-subtitle">Join us and start playing!</p>

        {success ? (
          <div className="success-message">
            Registration successful! Redirecting to login...
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="register-form">
            <div className="form-group">
              <label htmlFor="email" className="form-label">
                Email
              </label>
              <input
                type="email"
                id="email"
                name="email"
                value={formData.email}
                onChange={handleChange}
                className={`form-input ${errors.email ? 'error' : ''}`}
                disabled={loading}
                autoComplete="email"
              />
              {errors.email && (
                <span className="form-error">{errors.email}</span>
              )}
            </div>

            <div className="form-group">
              <label htmlFor="username" className="form-label">
                Username
              </label>
              <input
                type="text"
                id="username"
                name="username"
                value={formData.username}
                onChange={handleChange}
                className={`form-input ${errors.username ? 'error' : ''}`}
                disabled={loading}
                autoComplete="username"
              />
              {errors.username && (
                <span className="form-error">{errors.username}</span>
              )}
              {formData.username && !errors.username && (
                <span className="form-success">✓ Valid username</span>
              )}
            </div>

            <div className="form-group">
              <label htmlFor="password" className="form-label">
                Password
              </label>
              <input
                type="password"
                id="password"
                name="password"
                value={formData.password}
                onChange={handleChange}
                className={`form-input ${errors.password ? 'error' : ''}`}
                disabled={loading}
                autoComplete="new-password"
              />
              {errors.password && (
                <span className="form-error">{errors.password}</span>
              )}
              {passwordStrength && !errors.password && (
                <div className="password-strength">
                  <div className={`strength-bar ${passwordStrength}`}>
                    <div className="strength-fill"></div>
                  </div>
                  <span className={`strength-text ${passwordStrength}`}>
                    {passwordStrength.charAt(0).toUpperCase() + passwordStrength.slice(1)} password
                  </span>
                </div>
              )}
            </div>

            {authError && (
              <div className="auth-error">
                {authError}
              </div>
            )}

            <button
              type="submit"
              className="btn btn-primary"
              disabled={loading}
            >
              {loading ? 'Creating account...' : 'Register'}
            </button>

            <p className="register-link">
              Already have an account?{' '}
              <Link to="/login">Login here</Link>
            </p>
          </form>
        )}
      </div>
    </div>
  );
};

export default Register;
