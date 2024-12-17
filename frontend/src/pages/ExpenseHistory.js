import React, { useState, useEffect, useCallback } from 'react';
import { format } from 'date-fns';
import { expenses } from '../services/api';
import { DataGrid } from '@mui/x-data-grid';
import { Paper, Typography, Fab } from '@mui/material';
import { Add as AddIcon } from '@mui/icons-material';
import AddExpenseDialog from '../components/AddExpenseDialog';

function ExpenseHistory() {
  const [expenseList, setExpenseList] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [openAddExpense, setOpenAddExpense] = useState(false);

  const fetchExpenses = useCallback(async () => {
    try {
      const response = await expenses.getAll();
      setExpenseList(response.data);
      setLoading(false);
    } catch (err) {
      setError(err.message);
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchExpenses();
  }, [fetchExpenses]);

  const columns = [
    {
      field: 'date',
      headerName: 'Date',
      flex: 1,
      valueFormatter: (params) => format(new Date(params.value), 'MMM d, yyyy'),
    },
    {
      field: 'budget',
      headerName: 'Budget',
      flex: 1,
      renderCell: (params) => {
        if (params.row.budget_id === null) {
          return 'No Budget';
        }
        if (params.row.budget) {
          return params.row.budget.name;
        }
        return (
          <Typography color="text.disabled">
            (Deleted)
          </Typography>
        );
      },
    },
    {
      field: 'amount',
      headerName: 'Amount',
      flex: 1,
      valueFormatter: (params) => `$${params.value.toFixed(2)}`,
    },
    {
      field: 'description',
      headerName: 'Description',
      flex: 2,
    },
  ];

  const handleExpenseAdded = () => {
    fetchExpenses();
  };

  if (error) return <Typography color="error">Error: {error}</Typography>;

  return (
    <Paper sx={{ p: 3 }}>
      <Typography variant="h5" gutterBottom>
        Expense History
      </Typography>
      <DataGrid
        rows={expenseList}
        columns={columns}
        loading={loading}
        autoHeight
        pageSizeOptions={[5, 10, 25]}
        initialState={{
          pagination: { paginationModel: { pageSize: 25 } },
        }}
        sx={{ mt: 2 }}
      />
      <Fab
        color="primary"
        sx={{ position: 'fixed', bottom: 24, right: 24 }}
        onClick={() => setOpenAddExpense(true)}
      >
        <AddIcon />
      </Fab>
      <AddExpenseDialog
        open={openAddExpense}
        onClose={() => setOpenAddExpense(false)}
        onExpenseAdded={handleExpenseAdded}
      />
    </Paper>
  );
}

export default ExpenseHistory; 