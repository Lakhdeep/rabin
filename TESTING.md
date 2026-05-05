# Testing the Tic-Tac-Toe Authentication API

This guide will help you test the authentication endpoints locally.

## Prerequisites

1. **Docker** (for PostgreSQL database)
2. **Go 1.21+** (for running the backend)

## Quick Start

### 1. Start the Database

```bash
# Start PostgreSQL container
docker-compose up -d postgres

# Verify it's running
docker-compose ps
```

### 2. Configure Environment Variables

Create a `.env` file in the root directory (you can copy from `.env.example`):

```bash
cp .env.example .env
```

The default values should work for local development:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=tictactoe
DB_PASSWORD=tictactoe123
DB_NAME=tictactoe
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
PORT=8080
```

### 3. Run the Server

```bash
cd backend
go run cmd/server/main.go
```

You should see:
```
✓ Server started on port 8080
✓ Health check: http://localhost:8080/api/v1/health
✓ Auth endpoints:
  - POST http://localhost:8080/api/v1/auth/register
  - POST http://localhost:8080/api/v1/auth/login
  - GET  http://localhost:8080/api/v1/auth/me (requires JWT)
```

## Testing the Endpoints

### 1. Health Check

```bash
curl http://localhost:8080/api/v1/health
```

**Expected Response:**
```json
{
  "status": "ok",
  "timestamp": 1714826400
}
```

### 2. Register a New User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "username": "alice",
    "password": "Password123"
  }'
```

**Expected Response (201 Created):**
```json
{
  "id": 1,
  "email": "alice@example.com",
  "username": "alice",
  "created_at": "2026-05-04T11:16:00Z"
}
```

**Error Cases:**

Duplicate email:
```json
{
  "error": "Email already registered",
  "code": "DUPLICATE_EMAIL"
}
```

Duplicate username:
```json
{
  "error": "Username already taken",
  "code": "DUPLICATE_USERNAME"
}
```

Weak password:
```json
{
  "error": "password must be at least 8 characters",
  "code": "VALIDATION_ERROR"
}
```

### 3. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "Password123"
  }'
```

**Expected Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "alice@example.com",
    "username": "alice"
  }
}
```

**Save the token** - you'll need it for authenticated requests!

**Error Case (Invalid Credentials):**
```json
{
  "error": "Invalid credentials",
  "code": "INVALID_CREDENTIALS"
}
```

### 4. Get Current User (Authenticated)

```bash
# Replace YOUR_TOKEN_HERE with the token from login response
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer $TOKEN"
```

**Expected Response (200 OK):**
```json
{
  "id": 1,
  "email": "alice@example.com",
  "username": "alice",
  "total_games": 0,
  "wins": 0,
  "losses": 0,
  "draws": 0,
  "created_at": "2026-05-04T11:16:00Z"
}
```

**Error Case (Missing Token):**
```json
{
  "error": "Missing authorization header",
  "code": "UNAUTHORIZED"
}
```

**Error Case (Invalid Token):**
```json
{
  "error": "Invalid token",
  "code": "UNAUTHORIZED"
}
```

## Using Postman or Thunder Client

### Import Collection

Create a new collection with these requests:

**1. Register**
- Method: POST
- URL: `http://localhost:8080/api/v1/auth/register`
- Headers: `Content-Type: application/json`
- Body (JSON):
  ```json
  {
    "email": "alice@example.com",
    "username": "alice",
    "password": "Password123"
  }
  ```

**2. Login**
- Method: POST
- URL: `http://localhost:8080/api/v1/auth/login`
- Headers: `Content-Type: application/json`
- Body (JSON):
  ```json
  {
    "email": "alice@example.com",
    "password": "Password123"
  }
  ```
- After login, **copy the token** from response

**3. Get Current User**
- Method: GET
- URL: `http://localhost:8080/api/v1/auth/me`
- Headers: 
  - `Content-Type: application/json`
  - `Authorization: Bearer YOUR_TOKEN_HERE`

## Validation Rules

### Email
- Must be valid email format
- Must be unique

### Username
- 3-20 characters
- Only letters and numbers (no spaces or special characters)
- Must be unique

### Password
- Minimum 8 characters
- At least 1 uppercase letter
- At least 1 lowercase letter
- At least 1 number

## Troubleshooting

### Database Connection Failed

```json
{
  "status": "error",
  "message": "Database connection failed"
}
```

**Solution:** Make sure PostgreSQL is running:
```bash
docker-compose up -d postgres
docker-compose ps
```

### Port Already in Use

```
Failed to start server: listen tcp :8080: bind: address already in use
```

**Solution:** Either:
1. Stop the process using port 8080
2. Change the port in `.env`: `PORT=8081`

### Migrations Not Applied

```
Warning: Failed to run migrations: ...
```

**Solution:** 
- This is usually OK - migrations may already be applied
- If tables don't exist, check database connection and migrations folder

## Database Access

To inspect the database directly:

```bash
# Connect to PostgreSQL
docker-compose exec postgres psql -U tictactoe -d tictactoe

# List tables
\dt

# View users
SELECT * FROM users;

# Exit
\q
```

## Next Steps

Once authentication is working, the next phases will add:
- **Game Logic** - Create and play tic-tac-toe games
- **AI Opponent** - Play against AI with multiple difficulty levels
- **Score Tracking** - Track wins, losses, and draws
- **Frontend** - React web interface

## Current Implementation Status

✅ **Completed:**
- User registration with validation
- User login with JWT tokens
- Get current user endpoint
- Password hashing with bcrypt
- Input validation (email, username, password)
- Health check endpoint
- CORS support
- Request logging

🚧 **Coming Soon:**
- Game endpoints
- AI opponent
- Frontend interface
- Docker containerization for backend

---

Enjoy testing the API! 🎮
