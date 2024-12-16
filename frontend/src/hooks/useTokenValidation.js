import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { auth } from '../services/api';

export function useTokenValidation(interval = 60000) { // Default: check every minute
  const navigate = useNavigate();

  useEffect(() => {
    const validateToken = async () => {
      const token = localStorage.getItem('token');
      if (!token) {
        navigate('/login');
        return;
      }

      try {
        await auth.validate();
      } catch (error) {
        console.error('Token validation failed:', error);
        localStorage.removeItem('token');
        navigate('/login');
      }
    };

    // Validate immediately
    validateToken();

    // Set up periodic validation
    const intervalId = setInterval(validateToken, interval);

    return () => clearInterval(intervalId);
  }, [navigate, interval]);
} 