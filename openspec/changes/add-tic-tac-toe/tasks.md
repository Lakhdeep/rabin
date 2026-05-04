## 1. Project Setup

- [x] 1.1 Create monorepo directory structure (backend/ and frontend/ directories)
- [x] 1.2 Initialize Golang module in backend/ with go mod init
- [x] 1.3 Add backend dependencies (gin, jwt-go, bcrypt, lib/pq, golang-migrate)
- [x] 1.4 Initialize React app in frontend/ using Vite
- [x] 1.5 Add frontend dependencies (react-router-dom, axios)
- [x] 1.6 Create docker-compose.yml with postgres service
- [x] 1.7 Set up golang-migrate for database migrations
- [x] 1.8 Create .env.example file with required environment variables
- [x] 1.9 Add .env to .gitignore
- [x] 1.10 Create initial database migration for users table
- [x] 1.11 Create migration for games table
- [x] 1.12 Create migration for game_moves table

## 2. Backend - Project Structure

- [x] 2.1 Create cmd/server/main.go as application entry point
- [x] 2.2 Create internal/storage/database.go for database connection
- [x] 2.3 Create internal/storage/migrations.go for migration runner
- [x] 2.4 Create pkg/config/config.go for environment configuration
- [x] 2.5 Create pkg/logger/logger.go for structured logging
- [x] 2.6 Set up backend directory structure (api/v1, auth, game, user packages)

## 3. Backend - Database Layer

- [x] 3.1 Implement database connection pool with PostgreSQL
- [x] 3.2 Create internal/storage/user_repository.go for user CRUD operations
- [x] 3.3 Create internal/storage/game_repository.go for game CRUD operations
- [x] 3.4 Implement user creation with unique email/username constraints
- [x] 3.5 Implement user retrieval by email and by ID
- [x] 3.6 Implement game creation and state persistence
- [x] 3.7 Implement game move recording in game_moves table
- [x] 3.8 Implement user statistics updates (wins, losses, draws)
- [x] 3.9 Implement game history retrieval for user
- [x] 3.10 Add database query logging

## 4. Backend - Authentication

- [x] 4.1 Create internal/auth/password.go with bcrypt hashing (cost 12)
- [x] 4.2 Create internal/auth/jwt.go for token generation and validation
- [x] 4.3 Implement JWT token generation with HS256 algorithm
- [x] 4.4 Implement JWT token validation function
- [x] 4.5 Create internal/auth/middleware.go for JWT middleware
- [x] 4.6 Implement middleware to extract user from token
- [x] 4.7 Create internal/user/models.go with User struct
- [x] 4.8 Create internal/user/validation.go for email, username, password validation
- [x] 4.9 Implement registration validation logic (email format, username length, password strength)
- [x] 4.10 Write unit tests for password hashing and verification
- [x] 4.11 Write unit tests for JWT generation and validation
- [x] 4.12 Write unit tests for input validation functions

## 5. Backend - Authentication API

- [ ] 5.1 Create internal/api/v1/auth_handler.go
- [ ] 5.2 Implement POST /api/v1/auth/register handler
- [ ] 5.3 Implement POST /api/v1/auth/login handler
- [ ] 5.4 Implement GET /api/v1/auth/me handler (protected)
- [ ] 5.5 Add error handling for duplicate email/username
- [ ] 5.6 Add error handling for invalid credentials
- [ ] 5.7 Return appropriate HTTP status codes (201, 400, 401)
- [ ] 5.8 Write integration tests for registration endpoint
- [ ] 5.9 Write integration tests for login endpoint
- [ ] 5.10 Write integration tests for get current user endpoint

## 6. Backend - Game Logic

- [ ] 6.1 Create internal/game/board.go with board representation
- [ ] 6.2 Implement board initialization (empty 3x3 grid)
- [ ] 6.3 Implement move validation (position empty, in range, game not ended)
- [ ] 6.4 Implement win condition detection (rows, columns, diagonals)
- [ ] 6.5 Implement draw condition detection (board full, no winner)
- [ ] 6.6 Create internal/game/models.go with Game struct
- [ ] 6.7 Implement turn management (X starts, alternating turns)
- [ ] 6.8 Implement game state transitions (active, won, lost, draw)
- [ ] 6.9 Write unit tests for board initialization
- [ ] 6.10 Write unit tests for move validation
- [ ] 6.11 Write unit tests for win detection (all 8 winning positions)
- [ ] 6.12 Write unit tests for draw detection

## 7. Backend - AI Implementation

- [ ] 7.1 Create internal/game/ai/ai.go with AIStrategy interface
- [ ] 7.2 Implement EasyAI struct with random move selection
- [ ] 7.3 Implement MediumAI struct with mixed strategy (50% minimax, 50% random)
- [ ] 7.4 Implement minimax algorithm with alpha-beta pruning
- [ ] 7.5 Implement HardAI struct with depth-limited minimax (depth=4)
- [ ] 7.6 Implement ImpossibleAI struct with full minimax
- [ ] 7.7 Create AI factory function to return AI by difficulty string
- [ ] 7.8 Add timeout protection for AI calculations (1 second max)
- [ ] 7.9 Write unit tests for EasyAI (verifies random but valid moves)
- [ ] 7.10 Write unit tests for minimax algorithm correctness
- [ ] 7.11 Write unit tests for ImpossibleAI (verify unbeatable)
- [ ] 7.12 Write performance tests for AI response time

## 8. Backend - Game API

- [ ] 8.1 Create internal/api/v1/game_handler.go
- [ ] 8.2 Implement POST /api/v1/games handler (create game with difficulty)
- [ ] 8.3 Implement GET /api/v1/games/:id handler (get game state)
- [ ] 8.4 Implement POST /api/v1/games/:id/move handler (make move)
- [ ] 8.5 Implement GET /api/v1/games handler (list user's games)
- [ ] 8.6 Add game ownership validation (user can only access their games)
- [ ] 8.7 Implement AI move calculation and response in move handler
- [ ] 8.8 Update user statistics when game ends (atomic transaction)
- [ ] 8.9 Return updated board state and game result in response
- [ ] 8.10 Add error handling for invalid moves (occupied position, wrong turn)
- [ ] 8.11 Add error handling for invalid difficulty
- [ ] 8.12 Write integration tests for create game endpoint
- [ ] 8.13 Write integration tests for make move endpoint
- [ ] 8.14 Write integration tests for game completion and score updates

## 9. Backend - User Statistics API

- [ ] 9.1 Create internal/api/v1/user_handler.go
- [ ] 9.2 Implement GET /api/v1/users/:id/stats handler
- [ ] 9.3 Calculate and return win_rate (wins / total_games)
- [ ] 9.4 Handle division by zero for users with no games
- [ ] 9.5 Make endpoint public (no authentication required)
- [ ] 9.6 Add error handling for user not found
- [ ] 9.7 Write integration tests for user stats endpoint

## 10. Backend - Server Setup

- [ ] 10.1 Create Gin router in main.go
- [ ] 10.2 Register all /api/v1/auth routes
- [ ] 10.3 Register all /api/v1/games routes with JWT middleware
- [ ] 10.4 Register all /api/v1/users routes
- [ ] 10.5 Add CORS middleware (allow all origins in dev, configurable for prod)
- [ ] 10.6 Add request logging middleware
- [ ] 10.7 Add request timeout middleware (30 seconds)
- [ ] 10.8 Create GET /api/v1/health endpoint with database check
- [ ] 10.9 Implement graceful shutdown on SIGINT/SIGTERM
- [ ] 10.10 Add panic recovery middleware
- [ ] 10.11 Test server starts and connects to database

## 11. Backend - Docker

- [ ] 11.1 Create backend/Dockerfile with multi-stage build
- [ ] 11.2 First stage: build Go binary with golang:1.21-alpine
- [ ] 11.3 Second stage: runtime with alpine:latest
- [ ] 11.4 Copy binary and expose port 8080
- [ ] 11.5 Update docker-compose.yml to include backend service
- [ ] 11.6 Configure backend environment variables in docker-compose
- [ ] 11.7 Add depends_on for postgres service
- [ ] 11.8 Add health check using /api/v1/health endpoint
- [ ] 11.9 Set restart policy to unless-stopped
- [ ] 11.10 Test backend container builds and starts
- [ ] 11.11 Verify backend connects to postgres container
- [ ] 11.12 Test migrations run on container startup

## 12. Frontend - Project Setup

- [ ] 12.1 Configure Vite with React and TypeScript (or JavaScript)
- [ ] 12.2 Install react-router-dom for routing
- [ ] 12.3 Install axios for API requests
- [ ] 12.4 Set up environment variable support (.env files)
- [ ] 12.5 Configure API base URL from environment variable
- [ ] 12.6 Create frontend/src/index.css with global styles
- [ ] 12.7 Set up CSS reset or normalize.css

## 13. Frontend - API Client

- [ ] 13.1 Create src/services/api.js with axios instance
- [ ] 13.2 Configure base URL from environment variable
- [ ] 13.3 Add request interceptor to include JWT token from localStorage
- [ ] 13.4 Add response interceptor for error handling (401, 500)
- [ ] 13.5 Create src/services/auth.js with register, login functions
- [ ] 13.6 Create src/services/game.js with game API functions
- [ ] 13.7 Create src/services/user.js with user stats functions
- [ ] 13.8 Implement token storage helpers (save, get, clear from localStorage)

## 14. Frontend - Authentication Context

- [ ] 14.1 Create src/context/AuthContext.jsx
- [ ] 14.2 Implement AuthProvider with user state
- [ ] 14.3 Implement login function (call API, store token, update state)
- [ ] 14.4 Implement logout function (clear token, clear state)
- [ ] 14.5 Implement register function
- [ ] 14.6 Add loading state for auth operations
- [ ] 14.7 Add error state for auth errors
- [ ] 14.8 Implement useAuth custom hook for accessing context
- [ ] 14.9 Check for existing token on app load

## 15. Frontend - Routing

- [ ] 15.1 Create src/App.jsx with BrowserRouter
- [ ] 15.2 Set up routes for /login, /register, /dashboard, /game, /profile
- [ ] 15.3 Create src/components/ProtectedRoute.jsx wrapper
- [ ] 15.4 Implement redirect to /login for unauthenticated users
- [ ] 15.5 Implement redirect to /dashboard for authenticated users on /login
- [ ] 15.6 Set default route (/) to redirect based on auth status

## 16. Frontend - Registration Page

- [ ] 16.1 Create src/pages/Register.jsx component
- [ ] 16.2 Create form with email, username, password inputs
- [ ] 16.3 Add form validation (email format, username length, password strength)
- [ ] 16.4 Implement real-time validation feedback
- [ ] 16.5 Add password strength indicator
- [ ] 16.6 Implement form submission handler
- [ ] 16.7 Display loading state during registration
- [ ] 16.8 Display error messages from API (duplicate email/username)
- [ ] 16.9 Redirect to login page on success with success message
- [ ] 16.10 Add link to login page for existing users
- [ ] 16.11 Style registration form with CSS

## 17. Frontend - Login Page

- [ ] 17.1 Create src/pages/Login.jsx component
- [ ] 17.2 Create form with email and password inputs
- [ ] 17.3 Add form validation (required fields)
- [ ] 17.4 Implement form submission handler
- [ ] 17.5 Display loading state during login
- [ ] 17.6 Display error message for invalid credentials
- [ ] 17.7 Redirect to dashboard on successful login
- [ ] 17.8 Add link to registration page for new users
- [ ] 17.9 Style login form with CSS

## 18. Frontend - Navigation

- [ ] 18.1 Create src/components/Navigation.jsx component
- [ ] 18.2 Add navigation links (Dashboard, New Game, Profile, Logout)
- [ ] 18.3 Show navigation only when user is authenticated
- [ ] 18.4 Highlight current page in navigation
- [ ] 18.5 Implement logout button handler
- [ ] 18.6 Style navigation bar with CSS
- [ ] 18.7 Make navigation responsive for mobile

## 19. Frontend - Dashboard Page

- [ ] 19.1 Create src/pages/Dashboard.jsx component
- [ ] 19.2 Fetch and display user statistics on mount
- [ ] 19.3 Display total_games, wins, losses, draws
- [ ] 19.4 Calculate and display win rate as percentage
- [ ] 19.5 Add prominent "New Game" button
- [ ] 19.6 Fetch and display recent games list
- [ ] 19.7 Show game result, difficulty, and timestamp for each game
- [ ] 19.8 Add loading state while fetching data
- [ ] 19.9 Handle errors when fetching statistics
- [ ] 19.10 Style dashboard with card layout

## 20. Frontend - Game Context

- [ ] 20.1 Create src/context/GameContext.jsx
- [ ] 20.2 Implement GameProvider with game state (board, difficulty, gameId)
- [ ] 20.3 Implement createGame function (call API, update state)
- [ ] 20.4 Implement makeMove function (call API, update board, handle AI response)
- [ ] 20.5 Implement resetGame function
- [ ] 20.6 Add loading state for game operations
- [ ] 20.7 Add error state for game errors
- [ ] 20.8 Implement useGame custom hook

## 21. Frontend - Game Board Component

- [ ] 21.1 Create src/components/GameBoard.jsx component
- [ ] 21.2 Render 3x3 grid of cells
- [ ] 21.3 Display X or O in occupied cells
- [ ] 21.4 Implement cell click handler
- [ ] 21.5 Disable clicks on occupied cells
- [ ] 21.6 Disable clicks during AI turn (loading state)
- [ ] 21.7 Add hover effect on empty cells when it's user's turn
- [ ] 21.8 Highlight winning line when game ends
- [ ] 21.9 Style game board with CSS (square cells, grid layout)
- [ ] 21.10 Make game board responsive (scale to screen size)

## 22. Frontend - Game Page

- [ ] 22.1 Create src/pages/Game.jsx component
- [ ] 22.2 Add difficulty selector (dropdown or buttons)
- [ ] 22.3 Implement "Start Game" button
- [ ] 22.4 Render GameBoard component
- [ ] 22.5 Display game status (Your turn, AI thinking, You won, etc.)
- [ ] 22.6 Display current turn indicator
- [ ] 22.7 Add "Play Again" button when game ends
- [ ] 22.8 Show loading spinner during AI move
- [ ] 22.9 Display error messages for invalid moves
- [ ] 22.10 Add visual feedback for game completion (win/loss/draw)
- [ ] 22.11 Update user statistics display when game ends
- [ ] 22.12 Style game page layout

## 23. Frontend - Profile Page

- [ ] 23.1 Create src/pages/Profile.jsx component
- [ ] 23.2 Fetch and display user profile information
- [ ] 23.3 Display email, username, account creation date
- [ ] 23.4 Display detailed statistics (total games, W/L/D, win rate)
- [ ] 23.5 Fetch and display complete game history
- [ ] 23.6 Group game history by date or show chronologically
- [ ] 23.7 Add filtering or sorting options for game history
- [ ] 23.8 Add loading state while fetching data
- [ ] 23.9 Handle errors gracefully
- [ ] 23.10 Style profile page with sections

## 24. Frontend - Error Handling

- [ ] 24.1 Create src/components/ErrorToast.jsx component
- [ ] 24.2 Implement toast notification system
- [ ] 24.3 Show error toasts for API failures
- [ ] 24.4 Auto-dismiss toasts after 5 seconds
- [ ] 24.5 Allow manual dismissal of toasts
- [ ] 24.6 Display specific error messages from API
- [ ] 24.7 Show generic error for unexpected failures
- [ ] 24.8 Style error toast with appropriate colors

## 25. Frontend - Loading States

- [ ] 25.1 Create src/components/Spinner.jsx component
- [ ] 25.2 Show spinner during API requests
- [ ] 25.3 Disable form inputs during submission
- [ ] 25.4 Disable game board during AI turn
- [ ] 25.5 Show skeleton loaders for dashboard statistics
- [ ] 25.6 Style loading spinner

## 26. Frontend - Responsive Design

- [ ] 26.1 Add media queries for mobile (<768px)
- [ ] 26.2 Add media queries for tablet (768-1024px)
- [ ] 26.3 Add media queries for desktop (>1024px)
- [ ] 26.4 Make game board responsive (scales to fit screen)
- [ ] 26.5 Make navigation responsive (hamburger menu on mobile)
- [ ] 26.6 Make forms responsive (stack inputs on mobile)
- [ ] 26.7 Test layout on different screen sizes
- [ ] 26.8 Ensure touch targets are 44px minimum for mobile

## 27. Frontend - Accessibility

- [ ] 27.1 Add ARIA labels to game board cells
- [ ] 27.2 Add ARIA labels to form inputs
- [ ] 27.3 Ensure all interactive elements are keyboard accessible
- [ ] 27.4 Add focus visible styles for keyboard navigation
- [ ] 27.5 Ensure color contrast meets WCAG AA standards
- [ ] 27.6 Add alt text for any images
- [ ] 27.7 Test with keyboard navigation (Tab, Enter, Space)
- [ ] 27.8 Test with screen reader (VoiceOver or NVDA)

## 28. Frontend - Docker

- [ ] 28.1 Create frontend/Dockerfile with multi-stage build
- [ ] 28.2 First stage: build React app with node:18-alpine
- [ ] 28.3 Second stage: serve with nginx:alpine
- [ ] 28.4 Copy build output to nginx html directory
- [ ] 28.5 Create nginx.conf for SPA routing (fallback to index.html)
- [ ] 28.6 Update docker-compose.yml to include frontend service
- [ ] 28.7 Expose port 80 and map to host
- [ ] 28.8 Configure API_URL environment variable
- [ ] 28.9 Set restart policy to unless-stopped
- [ ] 28.10 Test frontend container builds and starts
- [ ] 28.11 Verify frontend can communicate with backend container
- [ ] 28.12 Test SPA routing works (refresh on /dashboard doesn't 404)

## 29. Integration Testing

- [ ] 29.1 Test full registration flow (frontend → backend → database)
- [ ] 29.2 Test full login flow (frontend → backend → JWT storage)
- [ ] 29.3 Test game creation with all difficulty levels
- [ ] 29.4 Test complete game flow (moves → win detection → score update)
- [ ] 29.5 Test AI behavior for each difficulty level
- [ ] 29.6 Test invalid move handling (occupied position, wrong turn)
- [ ] 29.7 Test error handling for duplicate email/username
- [ ] 29.8 Test error handling for invalid credentials
- [ ] 29.9 Test protected routes redirect when not authenticated
- [ ] 29.10 Test token expiration handling (401 response)
- [ ] 29.11 Test concurrent games (multiple games for same user)
- [ ] 29.12 Test statistics accuracy after multiple games

## 30. Cross-Browser Testing

- [ ] 30.1 Test in Chrome (latest version)
- [ ] 30.2 Test in Firefox (latest version)
- [ ] 30.3 Test in Safari (latest version)
- [ ] 30.4 Test in Edge (latest version)
- [ ] 30.5 Test on iOS Safari
- [ ] 30.6 Test on Android Chrome
- [ ] 30.7 Fix any browser-specific issues

## 31. Polish & UX

- [ ] 31.1 Add smooth transitions for game state changes
- [ ] 31.2 Add animation for AI move (delay to show AI "thinking")
- [ ] 31.3 Add sound effects for moves (optional, with toggle)
- [ ] 31.4 Add celebration animation for wins
- [ ] 31.5 Improve error messages to be user-friendly
- [ ] 31.6 Add input validation feedback with icons (checkmark/X)
- [ ] 31.7 Add tooltips for difficulty levels explaining strategy
- [ ] 31.8 Improve visual design (colors, typography, spacing)

## 32. Documentation

- [ ] 32.1 Create comprehensive README.md in root directory
- [ ] 32.2 Document prerequisites (Docker, Docker Compose)
- [ ] 32.3 Document local setup instructions (clone, .env, docker-compose up)
- [ ] 32.4 Document API endpoints with examples
- [ ] 32.5 Document environment variables in .env.example
- [ ] 32.6 Add architecture diagram showing services
- [ ] 32.7 Document database schema
- [ ] 32.8 Add troubleshooting section for common issues
- [ ] 32.9 Document how to run tests
- [ ] 32.10 Add code comments in complex logic (AI, game logic)

## 33. Security & Validation

- [ ] 33.1 Verify all SQL queries use parameterized statements (prevent injection)
- [ ] 33.2 Verify password is never logged or returned in responses
- [ ] 33.3 Verify JWT secret is strong and from environment variable
- [ ] 33.4 Verify CORS is properly configured
- [ ] 33.5 Verify input sanitization on frontend
- [ ] 33.6 Verify error messages don't leak sensitive information
- [ ] 33.7 Test with intentionally malicious inputs
- [ ] 33.8 Review code for security vulnerabilities

## 34. Performance & Optimization

- [ ] 34.1 Verify AI move calculation is under 1 second for all difficulties
- [ ] 34.2 Test database query performance with sample data
- [ ] 34.3 Add database indexes on frequently queried columns
- [ ] 34.4 Optimize frontend bundle size (code splitting if needed)
- [ ] 34.5 Enable gzip compression on nginx
- [ ] 34.6 Test with multiple concurrent users (load testing)
- [ ] 34.7 Monitor memory usage of backend container
- [ ] 34.8 Verify no memory leaks in long-running games

## 35. Deployment Preparation

- [ ] 35.1 Add production environment configuration
- [ ] 35.2 Document production deployment steps
- [ ] 35.3 Add health check endpoints for monitoring
- [ ] 35.4 Configure logging for production (log levels, rotation)
- [ ] 35.5 Document backup and restore procedures for database
- [ ] 35.6 Add resource limits to docker-compose (memory, CPU)
- [ ] 35.7 Test deployment process end-to-end locally
- [ ] 35.8 Create deployment checklist
