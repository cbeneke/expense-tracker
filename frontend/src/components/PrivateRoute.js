import React from 'react';
import { Navigate } from 'react-router-dom';
import { useTokenValidation } from '../hooks/useTokenValidation';

const PrivateRoute = ({ children }) => {
  useTokenValidation();

  const token = localStorage.getItem('token');
  
  if (!token) {
    return <Navigate to="/login" />;
  }

  return children;
};

export default PrivateRoute; 