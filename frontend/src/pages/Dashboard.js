import React, { useState, useEffect } from 'react';
import {
  Grid,
  Paper,
  Typography,
  Box,
  Card,
  CardContent,
  IconButton,
  Fab,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
} from '@mui/material';
import {
  AccountBalance as AccountBalanceIcon,
  TrendingUp as TrendingUpIcon,
  TrendingDown as TrendingDownIcon,
  MoreVert as MoreVertIcon,
  Add as AddIcon,
} from '@mui/icons-material';
import AddExpenseDialog from '../components/AddExpenseDialog';
import { dashboard } from '../services/api';
import { format } from 'date-fns';

function Dashboard() {
  const [summary, setSummary] = useState({
    monthly_budget: 0,
    total_expenses: 0,
    remaining_budget: 0,
    current_month: '',
    recent_expenses: []
  });
  const [openAddExpense, setOpenAddExpense] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchSummary = async () => {
    try {
      const response = await dashboard.getSummary();
      setSummary(response.data);
      setLoading(false);
    } catch (err) {
      setError(err.message);
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchSummary();
  }, []);

  const handleExpenseAdded = () => {
    fetchSummary();
  };

  if (error) return <Typography color="error">Error: {error}</Typography>;
  if (loading) return <Typography>Loading...</Typography>;

  const StatCard = ({ title, value, icon, color }) => (
    <Card>
      <CardContent>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
          <Box
            sx={{
              backgroundColor: `${color}.lighter`,
              borderRadius: '50%',
              p: 1,
              display: 'flex',
              alignItems: 'center',
            }}
          >
            {icon}
          </Box>
          <IconButton size="small">
            <MoreVertIcon />
          </IconButton>
        </Box>
        <Typography variant="h4" component="div" sx={{ mb: 1 }}>
          ${value.toFixed(2)}
        </Typography>
        <Typography color="text.secondary" variant="body2">
          {title}
        </Typography>
      </CardContent>
    </Card>
  );

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Dashboard - {summary.current_month}
      </Typography>
      <Grid container spacing={3}>
        <Grid item xs={12} md={4}>
          <StatCard
            title="Monthly Budget"
            value={summary.monthly_budget}
            icon={<AccountBalanceIcon sx={{ color: 'primary.main' }} />}
            color="primary"
          />
        </Grid>
        <Grid item xs={12} md={4}>
          <StatCard
            title="Total Expenses"
            value={summary.total_expenses}
            icon={<TrendingDownIcon sx={{ color: 'error.main' }} />}
            color="error"
          />
        </Grid>
        <Grid item xs={12} md={4}>
          <StatCard
            title="Remaining Budget"
            value={summary.remaining_budget}
            icon={<TrendingUpIcon sx={{ color: 'success.main' }} />}
            color="success"
          />
        </Grid>
        <Grid item xs={12}>
          <Paper sx={{ p: 3 }}>
            <Typography variant="h6" gutterBottom>
              Recent Expenses
            </Typography>
            <Box sx={{ mt: 2 }}>
              {summary.recent_expenses?.length > 0 ? (
                <Table size="small">
                  <TableHead>
                    <TableRow>
                      <TableCell>Date</TableCell>
                      <TableCell>Budget</TableCell>
                      <TableCell>Description</TableCell>
                      <TableCell align="right">Amount</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {summary.recent_expenses.slice(0, 5).map((expense) => (
                      <TableRow key={expense.id}>
                        <TableCell>{format(new Date(expense.date), 'MMM d, yyyy')}</TableCell>
                        <TableCell>
                          {expense.budget ? (
                            expense.budget.name
                          ) : expense.budget_id === null ? (
                            'No Budget'
                          ) : (
                            <Typography color="text.disabled" component="span">
                              (Deleted)
                            </Typography>
                          )}
                        </TableCell>
                        <TableCell>{expense.description}</TableCell>
                        <TableCell align="right">${expense.amount.toFixed(2)}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              ) : (
                <Typography color="text.secondary" align="center">
                  No recent transactions
                </Typography>
              )}
            </Box>
          </Paper>
        </Grid>
      </Grid>
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
    </Box>
  );
}

export default Dashboard; 