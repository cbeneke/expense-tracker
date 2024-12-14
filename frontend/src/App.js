import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import Box from '@mui/material/Box';
import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import ExpenseHistory from './pages/ExpenseHistory';
import BudgetManagement from './pages/BudgetManagement';
import PrivateRoute from './components/PrivateRoute';
import Navigation from './components/Navigation';

const theme = createTheme({
  palette: {
    primary: {
      main: '#1976d2',
    },
    secondary: {
      main: '#dc004e',
    },
    background: {
      default: '#f5f5f5',
    },
  },
});

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Router>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route
            path="/"
            element={
              <PrivateRoute>
                <Box sx={{ display: 'flex' }}>
                  <Navigation />
                  <Box
                    component="main"
                    sx={{
                      flexGrow: 1,
                      p: 3,
                      bgcolor: 'background.default',
                      minHeight: '100vh',
                    }}
                  >
                    <Dashboard />
                  </Box>
                </Box>
              </PrivateRoute>
            }
          />
          <Route
            path="/history"
            element={
              <PrivateRoute>
                <Box sx={{ display: 'flex' }}>
                  <Navigation />
                  <Box
                    component="main"
                    sx={{
                      flexGrow: 1,
                      p: 3,
                      bgcolor: 'background.default',
                      minHeight: '100vh',
                    }}
                  >
                    <ExpenseHistory />
                  </Box>
                </Box>
              </PrivateRoute>
            }
          />
          <Route
            path="/budgets"
            element={
              <PrivateRoute>
                <Box sx={{ display: 'flex' }}>
                  <Navigation />
                  <Box
                    component="main"
                    sx={{
                      flexGrow: 1,
                      p: 3,
                      bgcolor: 'background.default',
                      minHeight: '100vh',
                    }}
                  >
                    <BudgetManagement />
                  </Box>
                </Box>
              </PrivateRoute>
            }
          />
        </Routes>
      </Router>
    </ThemeProvider>
  );
}

export default App; 