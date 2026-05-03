## Context

This is a greenfield full-stack web application with no existing codebase. The goal is to create a production-ready, containerized tic-tac-toe game with user authentication and persistent score tracking. The application will be deployed locally using Docker containers initially, with architecture that supports future cloud deployment and real-time multiplayer features.

Current state: Empty repository with only OpenSpec configuration.

Key constraints:
- Must use Golang for backend (requirement)
- Must use React for frontend (requirement)
- Must use PostgreSQL for database (requirement)
- Must support Docker deployment (requirement)
- Must track wins/losses/draws per user
- Must support 4 AI difficulty levels

## Goals / Non-Goals

**Goals:**
- Create a fully functional tic-tac-toe web application playable in browser
- Implement secure user authentication with JWT tokens
- Provide 4 AI difficulty levels (easy, medium, hard, impossible)
- Track and persist user game statistics (wins, losses, draws)
- Use monorepo structure with clear separation of frontend/backend
- Containerize all services for consistent deployment
- Version API endpoints for backward compatibility (`/api/v1/...`)
- Design extensible architecture that can support PvP multiplayer in the future

**Non-Goals:**
- Real-time multiplayer (PvP) - architecture supports it, but not implementing initially
- User profile pictures or social features
- Tournament/bracket systems
- Mobile native apps (web-first, but responsive)
- Advanced analytics or game replays (basic history only)
- Email verification for registration (can add later)
- Rate limiting or DDoS protection (local deployment initially)
- Production-grade monitoring/observability (basic health checks only)

## Decisions

### 1. Monorepo Structure

**Decision:** Use a monorepo with `backend/` and `frontend/` directories in the same repository.

**Rationale:**
- Simplifies development workflow (single clone, single docker-compose)
- Easier to keep frontend/backend API contracts in sync
- Shared documentation and issue tracking
- Simpler CI/CD setup for initial deployment

**Alternatives considered:**
- Separate repositories: More complex coordination, overkill for this scope
- Microservices architecture: Unnecessary complexity for a single-domain application

### 2. Backend Framework: Gin

**Decision:** Use Gin web framework for Golang backend.

**Rationale:**
- Lightweight and performant (fastest Golang HTTP framework)
- Rich middleware ecosystem (CORS, JWT, logging)
- Easy routing with parameter binding
- Well-documented and widely adopted

**Alternatives considered:**
- Echo: Similar performance, slightly different API
- Standard library net/http: More verbose, would need to build middleware layer
- Fiber: Fast but less mature ecosystem

### 3. Database Schema Design

**Decision:** Three core tables - `users`, `games`, `game_moves`.

**Schema:**
```sql
users (
  id SERIAL PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  username VARCHAR(50) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  total_games INT DEFAULT 0,
  wins INT DEFAULT 0,
  losses INT DEFAULT 0,
  draws INT DEFAULT 0,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
)

games (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id) ON DELETE CASCADE,
  difficulty VARCHAR(20) NOT NULL, -- easy/medium/hard/impossible
  result VARCHAR(10), -- win/loss/draw/null (in progress)
  board_state JSONB, -- current board state as JSON array
  current_turn VARCHAR(1), -- 'X' or 'O'
  created_at TIMESTAMP DEFAULT NOW(),
  completed_at TIMESTAMP
)

game_moves (
  id SERIAL PRIMARY KEY,
  game_id INT REFERENCES games(id) ON DELETE CASCADE,
  position INT NOT NULL, -- 0-8
  player VARCHAR(1) NOT NULL, -- 'X' or 'O'
  move_number INT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
)
```

**Rationale:**
- Normalized design prevents data duplication
- User stats denormalized (total_games, wins, etc.) for fast leaderboard queries
- Game state stored as JSONB for flexibility and fast retrieval
- Move history allows future replay features
- Cascade deletes maintain referential integrity

**Alternatives considered:**
- Single table with embedded game data: Poor query performance for game history
- NoSQL (MongoDB): Relational queries (user stats, leaderboards) fit SQL better

### 4. JWT Authentication

**Decision:** Use JWT tokens with HS256 signing, 24-hour expiration, stored in localStorage.

**Rationale:**
- Stateless authentication (no server-side session storage)
- Scalable (API servers don't need to share session state)
- Industry standard for SPAs
- Easy to implement with existing Golang libraries (golang-jwt)

**Security considerations:**
- JWT secret stored in environment variable
- HTTPS required in production (prevents token interception)
- localStorage vs cookies: Using localStorage for simplicity, aware of XSS risk
- Token refresh not implemented initially (user re-authenticates after 24h)

**Alternatives considered:**
- Session-based auth: Requires session store (Redis), adds complexity
- OAuth2: Overkill for email/password auth, can add later for social login

### 5. AI Implementation Strategy

**Decision:** Implement 4 difficulty levels with different algorithms:

1. **Easy**: Random selection from valid moves
2. **Medium**: 50% chance of optimal move (minimax), 50% random
3. **Hard**: Minimax with alpha-beta pruning, limited to depth 4
4. **Impossible**: Full minimax with alpha-beta pruning (unbeatable)

**Rationale:**
- Easy mode gives beginners a chance to win
- Progressive difficulty provides replay value
- Minimax is proven optimal for tic-tac-toe
- Alpha-beta pruning prevents performance issues
- Depth limiting for "hard" mode creates near-optimal but beatable AI

**Implementation:**
```go
type AIStrategy interface {
    GetMove(board Board) int
}

type EasyAI struct{}
type MediumAI struct{}
type HardAI struct{}
type ImpossibleAI struct{}
```

**Alternatives considered:**
- Neural network: Overkill for deterministic game with small state space
- Rule-based heuristics: Less elegant than minimax, harder to tune difficulty

### 6. API Versioning

**Decision:** All endpoints prefixed with `/api/v1/`.

**Rationale:**
- Supports backward compatibility when introducing breaking changes
- Client code explicitly knows which API version it's using
- Industry best practice for REST APIs
- Minimal overhead (just routing prefix)

**Future strategy:**
- `/api/v2/` for breaking changes
- Maintain v1 for 6-12 months after v2 release
- Version header alternative not chosen (URL versioning more explicit)

### 7. Frontend State Management

**Decision:** Use React Context API for global state (auth, game state).

**Rationale:**
- Built into React, no additional dependencies
- Sufficient for this app's complexity (2-3 global states)
- Avoids Redux boilerplate for simple use cases
- Easy to migrate to Redux/Zustand later if needed

**State structure:**
```javascript
AuthContext: { user, token, login(), logout(), register() }
GameContext: { board, difficulty, gameId, makeMove(), newGame() }
```

**Alternatives considered:**
- Redux: Too much boilerplate for this scope
- Zustand: Simpler than Redux, but Context API is sufficient
- Component state only: Would require prop drilling

### 8. Docker Multi-Stage Builds

**Decision:** Use multi-stage builds for both backend and frontend.

**Backend Dockerfile:**
```dockerfile
# Stage 1: Build
FROM golang:1.21-alpine AS builder
# ... build binary

# Stage 2: Runtime
FROM alpine:latest
COPY --from=builder /app/server .
CMD ["./server"]
```

**Frontend Dockerfile:**
```dockerfile
# Stage 1: Build
FROM node:18-alpine AS builder
# ... npm run build

# Stage 2: Serve
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
```

**Rationale:**
- Smaller production images (no build tools in runtime)
- Faster deployments (smaller image size)
- More secure (reduced attack surface)

### 9. Password Security

**Decision:** Use bcrypt with cost factor 12 for password hashing.

**Rationale:**
- Industry standard for password hashing
- Adaptive (cost can increase over time as hardware improves)
- Built-in salt generation
- Cost 12 balances security and performance (~250ms per hash)

**Alternatives considered:**
- Argon2: More modern, but bcrypt is battle-tested and simpler
- SHA-256: Cryptographic hash, NOT suitable for passwords (too fast)

### 10. CORS Configuration

**Decision:** Allow all origins in development, restrict to specific domains in production.

**Implementation:**
```go
// Development
router.Use(cors.New(cors.Config{
    AllowAllOrigins: true,
}))

// Production (env-based)
router.Use(cors.New(cors.Config{
    AllowOrigins: []string{os.Getenv("FRONTEND_URL")},
}))
```

**Rationale:**
- Flexible development (frontend can run on any port)
- Secure production (prevents unauthorized domain access)
- Environment-based configuration

## Risks / Trade-offs

### Risk: JWT Token Theft (XSS Attack)
**Mitigation:**
- Sanitize all user inputs on frontend
- Use React's built-in XSS protection (JSX escaping)
- Consider httpOnly cookies in future iteration for enhanced security
- HTTPS enforced in production

### Risk: Database Connection Pool Exhaustion
**Mitigation:**
- Set max connections to 25 (reasonable for local deployment)
- Implement connection timeouts (30 seconds)
- Monitor connection usage in logs
- Use connection pooling library (pgx)

### Risk: AI Performance Bottleneck (Impossible Mode)
**Mitigation:**
- Tic-tac-toe has small state space (max 9 moves)
- Minimax completes in <10ms even with full search
- If needed: add move timeout (1 second) and fallback to random
- Alpha-beta pruning reduces search space significantly

### Risk: Concurrent Game State Updates
**Mitigation:**
- Backend enforces turn validation (reject moves when not user's turn)
- Optimistic locking not needed (AI responds synchronously)
- Database transactions ensure atomic updates
- Frontend polls game state if needed (future WebSocket for PvP)

### Risk: Docker Volume Permissions on Different OS
**Mitigation:**
- Use named volumes instead of bind mounts for PostgreSQL data
- Document permission issues in README
- Provide docker-compose examples for macOS/Linux/Windows

### Trade-off: No Token Refresh Mechanism
**Impact:** Users must re-login after 24 hours.
**Rationale:** Simplifies initial implementation. Can add refresh tokens later if needed.
**Acceptance:** Reasonable for a game application (sessions don't need to be long-lived).

### Trade-off: No Real-Time Updates (Polling vs WebSocket)
**Impact:** Frontend shows AI moves only after user refresh or next action.
**Rationale:** AI responds synchronously, so no polling needed. WebSocket adds complexity.
**Future:** When adding PvP, implement WebSocket for real-time board updates.

### Trade-off: Limited Input Validation
**Impact:** Relies on bcrypt and PostgreSQL constraints for security.
**Mitigation:**
- Email format validation (regex)
- Username length/character restrictions (3-20 alphanumeric)
- Password minimum requirements (8+ chars, complexity)
- SQL injection prevented by parameterized queries
- Can add rate limiting later if abuse occurs

## Migration Plan

### Initial Deployment (Local Docker)

1. **Prerequisites:**
   - Docker and Docker Compose installed
   - Ports 80 (frontend), 8080 (backend), 5432 (postgres) available

2. **Deployment Steps:**
   ```bash
   # Clone repository
   git clone <repo-url>
   cd rabin
   
   # Create environment file
   cp .env.example .env
   # Edit .env with JWT_SECRET and DB credentials
   
   # Start all services
   docker-compose up -d
   
   # Run database migrations
   docker-compose exec backend ./migrate up
   
   # Verify health
   curl http://localhost:8080/api/v1/health
   curl http://localhost
   ```

3. **Verification:**
   - Backend health check: `GET /api/v1/health` returns 200
   - Frontend loads at `http://localhost`
   - Register test user successfully
   - Play complete game and verify stats update

4. **Rollback Strategy:**
   ```bash
   # Stop services
   docker-compose down
   
   # If database corruption, restore from volume backup
   docker volume ls
   docker run --rm -v rabin_postgres_data:/data -v $(pwd):/backup alpine tar -czf /backup/db-backup.tar.gz /data
   ```

### Future Production Deployment

(Not implementing initially, but design supports it)

- Deploy to cloud platform (AWS ECS, GCP Cloud Run, DigitalOcean)
- Use managed PostgreSQL (AWS RDS, GCP Cloud SQL)
- Add HTTPS with Let's Encrypt or cloud load balancer
- Implement CI/CD pipeline (GitHub Actions)
- Add environment-specific configs (staging, production)
- Set up monitoring (Prometheus, Grafana) and logging (ELK stack)

## Open Questions

1. **Leaderboard pagination:** Should we paginate leaderboard results? If so, how many users per page?
   - **Recommendation:** Start with top 100, add pagination later if needed

2. **Game cleanup:** Should we delete old games after X months to reduce database size?
   - **Recommendation:** No cleanup initially, add archive strategy if storage becomes issue

3. **Username changes:** Should users be able to change their username after registration?
   - **Recommendation:** No for initial version (simplifies uniqueness constraints)

4. **Forgot password:** Should we implement password reset via email?
   - **Recommendation:** No for initial version (local deployment, no email service)

5. **HTTPS in local development:** Should docker-compose include self-signed certificates?
   - **Recommendation:** No, use HTTP locally. Document HTTPS setup for production.

6. **Frontend build optimization:** Should we implement code splitting and lazy loading?
   - **Recommendation:** Not initially, app is small. Add if bundle size exceeds 500KB.
