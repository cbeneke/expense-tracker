# Expense Tracker Application

A full-stack application for tracking personal expenses and managing budgets. Built with React (frontend) and Go (backend).

## Features

- **User Authentication**
  - Email/password-based authentication
  - Secure password hashing
  - JWT-based session management

- **Budget Management**
  - Create and manage monthly budgets by budget
  - Support for budget rollover or reset each month
  - Track budget overruns and remaining amounts

- **Expense Tracking**
  - Add, edit, and delete expenses
  - Categorize expenses
  - Automatic date tracking with manual override
  - Optional expense descriptions

- **Reporting & Analytics**
  - Monthly overview of budgets vs. expenses
  - Budget-wise breakdown
  - Historical expense data
  - Export data in CSV or JSON formats

## Technology Stack

### Frontend
- React.js
- React Router for navigation
- Tailwind CSS for styling
- Axios for API communication
- Chart.js for data visualization

### Backend
- Go (Golang)
- Gin web framework
- GORM for database operations
- JWT for authentication
- PostgreSQL database

## Project Structure

```
.
├── frontend/                # React frontend application
│   ├── src/
│   │   ├── components/     # Reusable React components
│   │   ├── pages/         # Page components
│   │   ├── services/      # API services
│   │   └── utils/         # Utility functions
│   └── Dockerfile         # Frontend Docker configuration
│
├── backend/                # Go backend application
│   ├── cmd/
│   │   └── server/        # Application entry point
│   ├── internal/
│   │   ├── api/          # API handlers and routes
│   │   ├── auth/         # Authentication logic
│   │   ├── models/       # Database models
│   │   └── database/     # Database configuration
│   └── Dockerfile        # Backend Docker configuration
│
└── docker-compose.yml     # Docker compose configuration
```

## Getting Started

### Prerequisites
- Docker and Docker Compose
- Go 1.23 or later (for local development)
- Node.js 20 or later (for local development)
- PostgreSQL 16 or later

### Running with Docker-Compose

1. Clone the repository:
```bash
git clone https://github.com/cbeneke/expense-tracker.git
cd expense-tracker
```

2. Create a `.env` file in the root directory (optional, to override defaults):
```env
DB_USER=myuser
DB_PASSWORD=mypassword
DB_NAME=expense_tracker
JWT_SECRET=my-secret-key
```

3. Start the application:
```bash
docker-compose up -d --build
```

The application will be available at:
- Frontend: http://localhost:80
- Backend API: http://localhost:8080

View the logs:
```bash
docker-compose logs -f
```

### Local Development

#### Backend
```bash
cd backend
go mod download
go run cmd/server/main.go
```

#### Frontend
```bash
cd frontend
npm install
npm start
```

## API Documentation

### Authentication Endpoints
- `POST /auth/signup` - Create a new account
- `POST /auth/login` - Login with email and password
- `POST /auth/logout` - Logout current user

### Budget Endpoints
- `GET /budgets` - Get all budgets
- `POST /budgets` - Create a new budget
- `PUT /budgets/:id` - Update a budget
- `DELETE /budgets/:id` - Delete a budget
- `GET /budgets/overview` - Get budget overview

### Expense Endpoints
- `GET /expenses` - Get all expenses
- `POST /expenses` - Create a new expense
- `PUT /expenses/:id` - Update an expense
- `DELETE /expenses/:id` - Delete an expense

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io)
- [React](https://reactjs.org)
- [Tailwind CSS](https://tailwindcss.com)
