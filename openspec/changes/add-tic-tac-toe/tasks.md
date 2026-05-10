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

- [x] 5.1 Create internal/api/v1/auth_handler.go
- [x] 5.2 Implement POST /api/v1/auth/register handler
- [x] 5.3 Implement POST /api/v1/auth/login handler
- [x] 5.4 Implement GET /api/v1/auth/me handler (protected)
- [x] 5.5 Add error handling for duplicate email/username
- [x] 5.6 Add error handling for invalid credentials
- [x] 5.7 Return appropriate HTTP status codes (201, 400, 401)
- [x] 5.8 Write integration tests for registration endpoint
- [x] 5.9 Write integration tests for login endpoint
- [x] 5.10 Write integration tests for get current user endpoint

## 6. Backend - Game Logic

- [x] 6.1 Create internal/game/board.go with board representation
- [x] 6.2 Implement board initialization (empty 3x3 grid)
- [x] 6.3 Implement move validation (position empty, in range, game not ended)
- [x] 6.4 Implement win condition detection (rows, columns, diagonals)
- [x] 6.5 Implement draw condition detection (board full, no winner)
- [x] 6.6 Create internal/game/models.go with Game struct
- [x] 6.7 Implement turn management (X starts, alternating turns)
- [x] 6.8 Implement game state transitions (active, won, lost, draw)
- [x] 6.9 Write unit tests for board initialization
- [x] 6.10 Write unit tests for move validation
- [x] 6.11 Write unit tests for win detection (all 8 winning positions)
- [x] 6.12 Write unit tests for draw detection

## 7. Backend - AI Implementation

- [x] 7.1 Create internal/game/ai/ai.go with AIStrategy interface
- [x] 7.2 Implement EasyAI struct with random move selection
- [x] 7.3 Implement MediumAI struct with mixed strategy (50% minimax, 50% random)
- [x] 7.4 Implement minimax algorithm with alpha-beta pruning
- [x] 7.5 Implement HardAI struct with depth-limited minimax (depth=4)
- [x] 7.6 Implement ImpossibleAI struct with full minimax
- [x] 7.7 Create AI factory function to return AI by difficulty string
- [x] 7.8 Add timeout protection for AI calculations (1 second max)
- [x] 7.9 Write unit tests for EasyAI (verifies random but valid moves)
- [x] 7.10 Write unit tests for minimax algorithm correctness
- [x] 7.11 Write unit tests for ImpossibleAI (verify unbeatable)
- [x] 7.12 Write performance tests for AI response time

## 8. Backend - Game API

- [x] 8.1 Create internal/api/v1/game_handler.go
- [x] 8.2 Implement POST /api/v1/games handler (create game with difficulty)
- [x] 8.3 Implement GET /api/v1/games/:id handler (get game state)
- [x] 8.4 Implement POST /api/v1/games/:id/move handler (make move)
- [x] 8.5 Implement GET /api/v1/games handler (list user's games)
- [x] 8.6 Add game ownership validation (user can only access their games)
- [x] 8.7 Implement AI move calculation and response in move handler
- [x] 8.8 Update user statistics when game ends (atomic transaction)
- [x] 8.9 Return updated board state and game result in response
- [x] 8.10 Add error handling for invalid moves (occupied position, wrong turn)
- [x] 8.11 Add error handling for invalid difficulty
- [x] 8.12 Write integration tests for create game endpoint
- [x] 8.13 Write integration tests for make move endpoint
- [x] 8.14 Write integration tests for game completion and score updates

## 9. Backend - User Statistics API

- [x] 9.1 Create internal/api/v1/user_handler.go
- [x] 9.2 Implement GET /api/v1/users/:id/stats handler
- [x] 9.3 Calculate and return win_rate (wins / total_games)
- [x] 9.4 Handle division by zero for users with no games
- [x] 9.5 Make endpoint public (no authentication required)
- [x] 9.6 Add error handling for user not found
- [x] 9.7 Write integration tests for user stats endpoint

## 10. Backend - Server Setup

- [x] 10.1 Create Gin router in main.go
- [x] 10.2 Register all /api/v1/auth routes
- [x] 10.3 Register all /api/v1/games routes with JWT middleware
- [x] 10.4 Register all /api/v1/users routes
- [x] 10.5 Add CORS middleware (allow all origins in dev, configurable for prod)
- [x] 10.6 Add request logging middleware
- [x] 10.7 Add request timeout middleware (30 seconds)
- [x] 10.8 Create GET /api/v1/health endpoint with database check
- [x] 10.9 Implement graceful shutdown on SIGINT/SIGTERM
- [x] 10.10 Add panic recovery middleware
- [x] 10.11 Test server starts and connects to database

## 11. Backend - Docker

- [x] 11.1 Create backend/Dockerfile with multi-stage build
- [x] 11.2 First stage: build Go binary with golang:1.21-alpine
- [x] 11.3 Second stage: runtime with alpine:latest
- [x] 11.4 Copy binary and expose port 8080
- [x] 11.5 Update docker-compose.yml to include backend service
- [x] 11.6 Configure backend environment variables in docker-compose
- [x] 11.7 Add depends_on for postgres service
- [x] 11.8 Add health check using /api/v1/health endpoint
- [x] 11.9 Set restart policy to unless-stopped
- [x] 11.10 Test backend container builds and starts
- [x] 11.11 Verify backend connects to postgres container
- [x] 11.12 Test migrations run on container startup

## 12. Frontend - Project Setup

- [x] 12.1 Configure Vite with React and TypeScript (or JavaScript)
- [x] 12.2 Install react-router-dom for routing
- [x] 12.3 Install axios for API requests
- [x] 12.4 Set up environment variable support (.env files)
- [x] 12.5 Configure API base URL from environment variable
- [x] 12.6 Create frontend/src/index.css with global styles
- [x] 12.7 Set up CSS reset or normalize.css

## 13. Frontend - API Client

- [x] 13.1 Create src/services/api.js with axios instance
- [x] 13.2 Configure base URL from environment variable
- [x] 13.3 Add request interceptor to include JWT token from localStorage
- [x] 13.4 Add response interceptor for error handling (401, 500)
- [x] 13.5 Create src/services/auth.js with register, login functions
- [x] 13.6 Create src/services/game.js with game API functions
- [x] 13.7 Create src/services/user.js with user stats functions
- [x] 13.8 Implement token storage helpers (save, get, clear from localStorage)

## 14. Frontend - Authentication Context

- [x] 14.1 Create src/context/AuthContext.jsx
- [x] 14.2 Implement AuthProvider with user state
- [x] 14.3 Implement login function (call API, store token, update state)
- [x] 14.4 Implement logout function (clear token, clear state)
- [x] 14.5 Implement register function
- [x] 14.6 Add loading state for auth operations
- [x] 14.7 Add error state for auth errors
- [x] 14.8 Implement useAuth custom hook for accessing context
- [x] 14.9 Check for existing token on app load

## 15. Frontend - Routing

- [x] 15.1 Create src/App.jsx with BrowserRouter
- [x] 15.2 Set up routes for /login, /register, /dashboard, /game, /profile
- [x] 15.3 Create src/components/ProtectedRoute.jsx wrapper
- [x] 15.4 Implement redirect to /login for unauthenticated users
- [x] 15.5 Implement redirect to /dashboard for authenticated users on /login
- [x] 15.6 Set default route (/) to redirect based on auth status

## 16. Frontend - Registration Page

- [x] 16.1 Create src/pages/Register.jsx component
- [x] 16.2 Create form with email, username, password inputs
- [x] 16.3 Add form validation (email format, username length, password strength)
- [x] 16.4 Implement real-time validation feedback
- [x] 16.5 Add password strength indicator
- [x] 16.6 Implement form submission handler
- [x] 16.7 Display loading state during registration
- [x] 16.8 Display error messages from API (duplicate email/username)
- [x] 16.9 Redirect to login page on success with success message
- [x] 16.10 Add link to login page for existing users
- [x] 16.11 Style registration form with CSS

## 17. Frontend - Login Page

- [x] 17.1 Create src/pages/Login.jsx component
- [x] 17.2 Create form with email and password inputs
- [x] 17.3 Add form validation (required fields)
- [x] 17.4 Implement form submission handler
- [x] 17.5 Display loading state during login
- [x] 17.6 Display error message for invalid credentials
- [x] 17.7 Redirect to dashboard on successful login
- [x] 17.8 Add link to registration page for new users
- [x] 17.9 Style login form with CSS

## 18. Frontend - Navigation

- [x] 18.1 Create src/components/Navigation.jsx component
- [x] 18.2 Add navigation links (Dashboard, New Game, Profile, Logout)
- [x] 18.3 Show navigation only when user is authenticated
- [x] 18.4 Highlight current page in navigation
- [x] 18.5 Implement logout button handler
- [x] 18.6 Style navigation bar with CSS
- [x] 18.7 Make navigation responsive for mobile

## 19. Frontend - Dashboard Page

- [x] 19.1 Create src/pages/Dashboard.jsx component
- [x] 19.2 Fetch and display user statistics on mount
- [x] 19.3 Display total_games, wins, losses, draws
- [x] 19.4 Calculate and display win rate as percentage
- [x] 19.5 Add prominent "New Game" button
- [x] 19.6 Fetch and display recent games list
- [x] 19.7 Show game result, difficulty, and timestamp for each game
- [x] 19.8 Add loading state while fetching data
- [x] 19.9 Handle errors when fetching statistics
- [x] 19.10 Style dashboard with card layout

## 20. Frontend - Game Context

- [x] 20.1 Create src/context/GameContext.jsx
- [x] 20.2 Implement GameProvider with game state (board, difficulty, gameId)
- [x] 20.3 Implement createGame function (call API, update state)
- [x] 20.4 Implement makeMove function (call API, update board, handle AI response)
- [x] 20.5 Implement resetGame function
- [x] 20.6 Add loading state for game operations
- [x] 20.7 Add error state for game errors
- [x] 20.8 Implement useGame custom hook

## 21. Frontend - Game Board Component

- [x] 21.1 Create src/components/GameBoard.jsx component
- [x] 21.2 Render 3x3 grid of cells
- [x] 21.3 Display X or O in occupied cells
- [x] 21.4 Implement cell click handler
- [x] 21.5 Disable clicks on occupied cells
- [x] 21.6 Disable clicks during AI turn (loading state)
- [x] 21.7 Add hover effect on empty cells when it's user's turn
- [x] 21.8 Highlight winning line when game ends
- [x] 21.9 Style game board with CSS (square cells, grid layout)
- [x] 21.10 Make game board responsive (scale to screen size)

## 22. Frontend - Game Page

- [x] 22.1 Create src/pages/Game.jsx component
- [x] 22.2 Add difficulty selector (dropdown or buttons)
- [x] 22.3 Implement "Start Game" button
- [x] 22.4 Render GameBoard component
- [x] 22.5 Display game status (Your turn, AI thinking, You won, etc.)
- [x] 22.6 Display current turn indicator
- [x] 22.7 Add "Play Again" button when game ends
- [x] 22.8 Show loading spinner during AI move
- [x] 22.9 Display error messages for invalid moves
- [x] 22.10 Add visual feedback for game completion (win/loss/draw)
- [x] 22.11 Update user statistics display when game ends
- [x] 22.12 Style game page layout

## 23. Frontend - Profile Page

- [x] 23.1 Create src/pages/Profile.jsx component
- [x] 23.2 Fetch and display user profile information
- [x] 23.3 Display email, username, account creation date
- [x] 23.4 Display detailed statistics (total games, W/L/D, win rate)
- [x] 23.5 Fetch and display complete game history
- [x] 23.6 Group game history by date or show chronologically
- [x] 23.7 Add filtering or sorting options for game history
- [x] 23.8 Add loading state while fetching data
- [x] 23.9 Handle errors gracefully
- [x] 23.10 Style profile page with sections

## 24. Frontend - Error Handling

- [x] 24.1 Create src/components/ErrorToast.jsx component
- [x] 24.2 Implement toast notification system
- [x] 24.3 Show error toasts for API failures
- [x] 24.4 Auto-dismiss toasts after 5 seconds
- [x] 24.5 Allow manual dismissal of toasts
- [x] 24.6 Display specific error messages from API
- [x] 24.7 Show generic error for unexpected failures
- [x] 24.8 Style error toast with appropriate colors

## 25. Frontend - Loading States

- [x] 25.1 Create src/components/Spinner.jsx component
- [x] 25.2 Show spinner during API requests
- [x] 25.3 Disable form inputs during submission
- [x] 25.4 Disable game board during AI turn
- [x] 25.5 Show skeleton loaders for dashboard statistics
- [x] 25.6 Style loading spinner

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

- [x] 28.1 Create frontend/Dockerfile with multi-stage build
- [x] 28.2 First stage: build React app with node:18-alpine
- [x] 28.3 Second stage: serve with nginx:alpine
- [x] 28.4 Copy build output to nginx html directory
- [x] 28.5 Create nginx.conf for SPA routing (fallback to index.html)
- [x] 28.6 Update docker-compose.yml to include frontend service
- [x] 28.7 Expose port 80 and map to host
- [x] 28.8 Configure API_URL environment variable
- [x] 28.9 Set restart policy to unless-stopped
- [x] 28.10 Test frontend container builds and starts
- [x] 28.11 Verify frontend can communicate with backend container
- [x] 28.12 Test SPA routing works (refresh on /dashboard doesn't 404)

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
