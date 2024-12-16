import React, { useState, useEffect, useCallback } from 'react';
import {
  Box,
  Typography,
  Grid,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  IconButton,
  LinearProgress,
  Card,
  CardContent,
  Tooltip,
} from '@mui/material';
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
} from '@mui/icons-material';
import { budgets } from '../services/api';

function BudgetManagement() {
  const [budgetList, setBudgetList] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [openDialog, setOpenDialog] = useState(false);
  const [selectedBudget, setSelectedBudget] = useState(null);
  const [formData, setFormData] = useState({
    name: '',
    amount: '',
  });

  const fetchBudgets = useCallback(async () => {
    try {
      const response = await budgets.getAll();
      setBudgetList(response.data);
      setLoading(false);
    } catch (err) {
      setError(err.message);
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchBudgets();
  }, [fetchBudgets]);

  const handleOpenDialog = (budget = null) => {
    if (budget) {
      setFormData({
        name: budget.name,
        amount: budget.amount.toString(),
      });
      setSelectedBudget(budget);
    } else {
      setFormData({ name: '', amount: '' });
      setSelectedBudget(null);
    }
    setOpenDialog(true);
  };

  const handleCloseDialog = () => {
    setOpenDialog(false);
    setSelectedBudget(null);
    setFormData({ name: '', amount: '' });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      if (selectedBudget) {
        await budgets.update(selectedBudget.id, {
          ...formData,
          amount: parseFloat(formData.amount),
        });
      } else {
        await budgets.create({
          ...formData,
          amount: parseFloat(formData.amount),
        });
      }
      fetchBudgets();
      handleCloseDialog();
    } catch (err) {
      setError(err.message);
    }
  };

  const handleDelete = async (id) => {
    if (window.confirm('Are you sure you want to delete this budget?')) {
      try {
        await budgets.delete(id);
        fetchBudgets();
      } catch (err) {
        setError(err.message);
      }
    }
  };

  const BudgetCard = ({ budget }) => (
    <Card>
      <CardContent>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
          <Typography variant="h6" component="div">
            {budget.name}
          </Typography>
          <Box>
            <Tooltip title="Edit">
              <IconButton
                size="small"
                onClick={() => handleOpenDialog(budget)}
                sx={{ mr: 1 }}
              >
                <EditIcon />
              </IconButton>
            </Tooltip>
            <Tooltip title="Delete">
              <IconButton
                size="small"
                onClick={() => handleDelete(budget.id)}
                color="error"
              >
                <DeleteIcon />
              </IconButton>
            </Tooltip>
          </Box>
        </Box>
        <Typography variant="h5" color="primary" sx={{ mb: 2 }}>
          ${budget.amount.toFixed(2)}
        </Typography>
        <Box sx={{ width: '100%', mb: 1 }}>
          <LinearProgress
            variant="determinate"
            value={((budget.roll_over_amount || 0) / budget.amount) * 100}
            color={budget.roll_over_amount > budget.amount ? 'error' : 'primary'}
          />
        </Box>
        <Typography variant="body2" color="text.secondary">
          Spent: ${(budget.roll_over_amount || 0).toFixed(2)} / ${budget.amount.toFixed(2)}
        </Typography>
      </CardContent>
    </Card>
  );

  if (error) return <Typography color="error">Error: {error}</Typography>;

  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
        <Typography variant="h4">Budget Management</Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => handleOpenDialog()}
        >
          Add Budget
        </Button>
      </Box>

      {loading ? (
        <LinearProgress />
      ) : (
        <Grid container spacing={3}>
          {budgetList.map((budget) => (
            <Grid item xs={12} md={4} key={budget.id}>
              <BudgetCard budget={budget} />
            </Grid>
          ))}
        </Grid>
      )}

      <Dialog open={openDialog} onClose={handleCloseDialog}>
        <DialogTitle>
          {selectedBudget ? 'Edit Budget' : 'Add New Budget'}
        </DialogTitle>
        <DialogContent>
          <Box component="form" onSubmit={handleSubmit} sx={{ mt: 2 }}>
            <TextField
              fullWidth
              label="Name"
              value={formData.name}
              onChange={(e) =>
                setFormData({ ...formData, name: e.target.value })
              }
              margin="normal"
              required
            />
            <TextField
              fullWidth
              label="Amount"
              type="number"
              value={formData.amount}
              onChange={(e) =>
                setFormData({ ...formData, amount: e.target.value })
              }
              margin="normal"
              required
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseDialog}>Cancel</Button>
          <Button onClick={handleSubmit} variant="contained">
            {selectedBudget ? 'Update' : 'Create'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}

export default BudgetManagement; 