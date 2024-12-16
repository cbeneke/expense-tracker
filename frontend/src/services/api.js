import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

const api = axios.create({
    baseURL: `${API_URL}/api`,
    headers: {
        'Content-Type': 'application/json',
    },
});

const authApi = axios.create({
    baseURL: `${API_URL}/auth`,
    headers: {
        'Content-Type': 'application/json',
    },
});

const addAuthHeader = (config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
};

api.interceptors.request.use(addAuthHeader);
authApi.interceptors.request.use(addAuthHeader);

export const auth = {
    login: (credentials) => authApi.post('/login', credentials),
    register: (userData) => authApi.post('/signup', userData),
    validate: () => authApi.get('/validate'),
    logout: () => {
        localStorage.removeItem('token');
        window.location.href = '/login';
    },
};

export const expenses = {
    getAll: () => api.get('/expenses'),
    getById: (id) => api.get(`/expenses/${id}`),
    create: (expense) => api.post('/expenses', {
        ...expense,
        budget_id: expense.budget ? parseInt(expense.budget) : null,
    }),
    update: (id, expense) => api.put(`/expenses/${id}`, expense),
    delete: (id) => api.delete(`/expenses/${id}`),
};

export const budgets = {
    getAll: () => api.get('/budgets'),
    getById: (id) => api.get(`/budgets/${id}`),
    create: (budget) => api.post('/budgets', budget),
    update: (id, budget) => api.put(`/budgets/${id}`, budget),
    delete: (id) => api.delete(`/budgets/${id}`),
};

export const dashboard = {
    getSummary: async () => {
        const currentMonth = new Date().toISOString().slice(0, 7); // YYYY-MM format

        const [budgetsRes, expensesRes] = await Promise.all([
            api.get('/budgets', { params: { month: currentMonth } }),
            api.get('/expenses', { params: { month: currentMonth } })
        ]);

        const budgets = budgetsRes.data;
        const expenses = expensesRes.data;
                
        const monthlyBudget = budgets
            .reduce((sum, b) => sum + b.amount, 0);

        const totalExpenses = expenses
            .reduce((sum, e) => sum + e.amount, 0);

        return {
            data: {
                monthly_budget: monthlyBudget,
                total_expenses: totalExpenses,
                remaining_budget: monthlyBudget - totalExpenses,
                current_month: currentMonth,
                recent_expenses: expenses.slice(0, 5)
            }
        };
    },
};

api.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            localStorage.removeItem('token');
            window.location.href = '/login';
        }
        return Promise.reject(error);
    }
);

export default api; 