# Tic-Tac-Toe Web Application

A full-stack tic-tac-toe game with user authentication, AI opponents, and score tracking.

## Tech Stack

- **Backend:** Golang (Gin framework)
- **Frontend:** React (Vite)
- **Database:** PostgreSQL
- **Authentication:** JWT tokens
- **Deployment:** Docker

## Current Status

✅ **Completed Features:**
- User authentication (register, login)
- Password hashing with bcrypt
- JWT token generation and validation
- Input validation (email, username, password)
- Database connection with PostgreSQL
- RESTful API endpoints
- CORS support
- Request logging

🚧 **In Progress:**
- Game logic and AI opponent
- Frontend React application
- Docker containerization

## Quick Start

### Prerequisites

- Docker (for PostgreSQL)
- Go 1.21+ (for backend)
- Node.js 18+ (for frontend, coming soon)

### 1. Start the Server

```bash
# Easy way - use the startup script
./start-server.sh

# Or manually:
# 1. Start PostgreSQL
docker-compose up -d postgres

# 2. Run the backend
cd backend
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

### 2. Test the API

```bash
# Automated tests
./test-api.sh

# Or manually test endpoints
curl http://localhost:8080/api/v1/health
```

See [TESTING.md](./TESTING.md) for detailed API testing guide.

## API Endpoints

### Authentication

- **POST** `/api/v1/auth/register` - Register new user
- **POST** `/api/v1/auth/login` - Login and get JWT token
- **GET** `/api/v1/auth/me` - Get current user (requires JWT)

### Health

- **GET** `/api/v1/health` - Health check with database status

## Example Usage

### Register a User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "username": "alice",
    "password": "Password123"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "Password123"
  }'
```

Response includes a JWT token:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": { "id": 1, "email": "alice@example.com", "username": "alice" }
}
```

### Get Current User

```bash
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## Project Structure

```
.
├── backend/
│   ├── cmd/server/          # Application entry point
│   ├── internal/
│   │   ├── api/v1/          # API handlers
│   │   ├── auth/            # Authentication logic
│   │   ├── storage/         # Database layer
│   │   └── user/            # User models
│   ├── pkg/
│   │   ├── config/          # Configuration
│   │   └── logger/          # Logging
│   └── migrations/          # Database migrations
├── frontend/                # React frontend (coming soon)
├── docker-compose.yml       # Docker services
└── .env.example            # Environment variables template
```

## Environment Variables

Copy `.env.example` to `.env` and configure:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=tictactoe
DB_PASSWORD=tictactoe123
DB_NAME=tictactoe

# JWT Secret
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Server
PORT=8080
```

## Development

### Run Tests

```bash
# Backend tests
cd backend
go test ./... -v

# Auth package tests
go test ./internal/auth/... -v

# API handler tests
go test ./internal/api/v1/... -v
```

### Database Access

```bash
# Connect to PostgreSQL
docker-compose exec postgres psql -U tictactoe -d tictactoe

# View users
SELECT * FROM users;
```

## Features Coming Soon

- 🎮 Tic-tac-toe game logic
- 🤖 AI opponent (Easy, Medium, Hard, Impossible)
- 📊 Score tracking and leaderboards
- ⚛️ React frontend
- 🐳 Complete Docker setup
- 🚀 Production deployment guide

## Documentation

- [TESTING.md](./TESTING.md) - API testing guide with examples
- [.env.example](./.env.example) - Environment configuration template

## License

MIT

## Contributing

This is a learning project. Feel free to explore and experiment!
