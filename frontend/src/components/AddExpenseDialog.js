import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  MenuItem,
  Box,
  Alert,
} from '@mui/material';
import { expenses, budgets } from '../services/api';

function AddExpenseDialog({ open, onClose, onExpenseAdded }) {
  const [availableBudgets, setAvailableBudgets] = useState([]);
  const [error, setError] = useState('');
  const [formData, setFormData] = useState({
    amount: '',
    budget: '',
    description: '',
    date: new Date().toISOString().split('T')[0],
  });

  useEffect(() => {
    const fetchBudgets = async () => {
      try {
        const response = await budgets.getAll();
        setAvailableBudgets(response.data);
      } catch (error) {
        console.error('Failed to fetch budgets:', error);
      }
    };

    if (open) {
      fetchBudgets();
    }
  }, [open]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    try {
      await expenses.create({
        ...formData,
        amount: parseFloat(formData.amount),
      });
      onExpenseAdded();
      handleClose();
    } catch (error) {
      setError(error.response?.data?.error || 'Failed to add expense');
    }
  };

  const handleClose = () => {
    setFormData({
      amount: '',
      budget: '',
      description: '',
      date: new Date().toISOString().split('T')[0],
    });
    onClose();
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>Add New Expense</DialogTitle>
      <DialogContent>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}
        <Box component="form" onSubmit={handleSubmit} sx={{ mt: 2 }}>
          <TextField
            fullWidth
            label="Amount"
            type="number"
            value={formData.amount}
            onChange={(e) =>
              setFormData({ ...formData, amount: e.target.value })
            }
            required
            margin="normal"
            inputProps={{ step: "0.01" }}
          />
          <TextField
            fullWidth
            select
            label="Budget"
            value={formData.budget}
            onChange={(e) =>
              setFormData({ ...formData, budget: e.target.value })
            }
            margin="normal"
            helperText="Select from your existing budgets (optional)"
          >
            <MenuItem value="">
              <em>No Budget</em>
            </MenuItem>
            {availableBudgets.map((budget) => (
              <MenuItem key={budget.id} value={budget.id}>
                {budget.category} (Available: ${(budget.amount - (budget.roll_over_amount || 0)).toFixed(2)})
              </MenuItem>
            ))}
          </TextField>
          <TextField
            fullWidth
            label="Description"
            value={formData.description}
            onChange={(e) =>
              setFormData({ ...formData, description: e.target.value })
            }
            required
            margin="normal"
          />
          <TextField
            fullWidth
            type="date"
            label="Date"
            value={formData.date}
            onChange={(e) => setFormData({ ...formData, date: e.target.value })}
            required
            margin="normal"
            InputLabelProps={{ shrink: true }}
          />
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose}>Cancel</Button>
        <Button onClick={handleSubmit} variant="contained">
          Add Expense
        </Button>
      </DialogActions>
    </Dialog>
  );
}

export default AddExpenseDialog; 