## ADDED Requirements

### Requirement: Docker Compose orchestration
The system SHALL use docker-compose to orchestrate all services.

#### Scenario: Three services defined
- **WHEN** docker-compose.yml is loaded
- **THEN** configuration includes backend, frontend, and postgres services

#### Scenario: Services start with single command
- **WHEN** user runs docker-compose up
- **THEN** all three services start and connect to shared network

#### Scenario: Services stop gracefully
- **WHEN** user runs docker-compose down
- **THEN** all services stop gracefully without data loss

### Requirement: Backend Docker container
The system SHALL containerize Golang backend application.

#### Scenario: Multi-stage build
- **WHEN** backend Dockerfile builds
- **THEN** first stage compiles Go binary, second stage creates minimal runtime image

#### Scenario: Alpine base image
- **WHEN** backend container runs
- **THEN** runtime image uses alpine:latest for minimal size

#### Scenario: Backend port exposed
- **WHEN** backend container starts
- **THEN** port 8080 is exposed and mapped to host

#### Scenario: Environment variables configured
- **WHEN** backend container starts
- **THEN** DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, and JWT_SECRET are available from environment

### Requirement: Frontend Docker container
The system SHALL containerize React frontend application.

#### Scenario: Multi-stage build for frontend
- **WHEN** frontend Dockerfile builds
- **THEN** first stage runs npm build, second stage serves with nginx

#### Scenario: Nginx serves static files
- **WHEN** frontend container runs
- **THEN** nginx serves built React files from /usr/share/nginx/html

#### Scenario: Frontend port exposed
- **WHEN** frontend container starts
- **THEN** port 80 is exposed and mapped to host port 80

#### Scenario: SPA fallback routing
- **WHEN** user navigates to non-root URL
- **THEN** nginx routes all requests to index.html for React Router

### Requirement: PostgreSQL container
The system SHALL run PostgreSQL database in container.

#### Scenario: PostgreSQL version
- **WHEN** postgres container starts
- **THEN** uses postgres:15-alpine image

#### Scenario: Database port exposed
- **WHEN** postgres container starts
- **THEN** port 5432 is exposed for backend connection

#### Scenario: Database credentials from environment
- **WHEN** postgres container initializes
- **THEN** uses POSTGRES_USER, POSTGRES_PASSWORD, and POSTGRES_DB from environment

#### Scenario: Data persistence with volume
- **WHEN** postgres container stops
- **THEN** data persists in named volume and survives container restart

### Requirement: Named volume for database
The system SHALL use named volume for PostgreSQL data.

#### Scenario: Volume created automatically
- **WHEN** docker-compose up runs first time
- **THEN** named volume is created for postgres data

#### Scenario: Data survives container deletion
- **WHEN** postgres container is removed and recreated
- **THEN** existing data is preserved in volume

#### Scenario: Volume independent of containers
- **WHEN** listing volumes
- **THEN** postgres volume exists independently of container state

### Requirement: Shared Docker network
The system SHALL create shared network for inter-service communication.

#### Scenario: Services on same network
- **WHEN** containers start
- **THEN** backend, frontend, and postgres are on same Docker network

#### Scenario: Service name DNS resolution
- **WHEN** backend connects to database
- **THEN** uses service name "postgres" as hostname

#### Scenario: Network isolated from host
- **WHEN** services communicate
- **THEN** traffic stays within Docker network unless explicitly exposed

### Requirement: Environment variable configuration
The system SHALL support configuration via environment files.

#### Scenario: .env file loaded
- **WHEN** docker-compose starts
- **THEN** variables from .env file are available to services

#### Scenario: .env.example provided
- **WHEN** repository is cloned
- **THEN** .env.example file shows required variables with placeholder values

#### Scenario: Secrets not in version control
- **WHEN** committing code
- **THEN** .env file is in .gitignore (only .env.example is committed)

### Requirement: Backend health check
The system SHALL implement health check for backend container.

#### Scenario: Health check endpoint called
- **WHEN** Docker runs health check
- **THEN** makes HTTP request to /api/v1/health

#### Scenario: Healthy status
- **WHEN** backend is running and database connected
- **THEN** health check returns success and container marked healthy

#### Scenario: Unhealthy status
- **WHEN** backend cannot connect to database
- **THEN** health check fails and container marked unhealthy

#### Scenario: Health check interval
- **WHEN** container is running
- **THEN** health check runs every 30 seconds

### Requirement: PostgreSQL health check
The system SHALL implement health check for postgres container.

#### Scenario: pg_isready check
- **WHEN** Docker runs postgres health check
- **THEN** executes pg_isready command

#### Scenario: Database ready
- **WHEN** postgres is accepting connections
- **THEN** health check returns success

#### Scenario: Database not ready
- **WHEN** postgres is still initializing
- **THEN** health check fails and retries

### Requirement: Container restart policy
The system SHALL configure restart policy for reliability.

#### Scenario: Restart policy set
- **WHEN** container stops unexpectedly
- **THEN** Docker automatically restarts it

#### Scenario: Manual stop respected
- **WHEN** user runs docker-compose down
- **THEN** containers do not auto-restart

#### Scenario: Restart policy is unless-stopped
- **WHEN** configuring services
- **THEN** restart policy is set to unless-stopped for all services

### Requirement: Build optimization
The system SHALL optimize Docker builds for speed and size.

#### Scenario: Layer caching for dependencies
- **WHEN** backend Dockerfile builds
- **THEN** go.mod and go.sum copied before source code to cache dependencies

#### Scenario: Layer caching for npm
- **WHEN** frontend Dockerfile builds
- **THEN** package.json and package-lock.json copied before source to cache node_modules

#### Scenario: Small image sizes
- **WHEN** images are built
- **THEN** backend runtime image is under 50MB, frontend under 100MB

#### Scenario: No build tools in runtime
- **WHEN** runtime images are created
- **THEN** Go compiler and Node.js are not included (only in build stage)

### Requirement: Database migrations on startup
The system SHALL run database migrations when backend starts.

#### Scenario: Migrations run automatically
- **WHEN** backend container starts
- **THEN** migrations are applied before accepting requests

#### Scenario: Migrations idempotent
- **WHEN** migrations run multiple times
- **THEN** already-applied migrations are skipped

#### Scenario: Migration failure prevents startup
- **WHEN** migration fails
- **THEN** backend container exits with error and restarts

### Requirement: Service dependencies
The system SHALL ensure services start in correct order.

#### Scenario: Backend depends on postgres
- **WHEN** docker-compose starts services
- **THEN** postgres starts before backend

#### Scenario: Backend waits for database ready
- **WHEN** backend container starts
- **THEN** waits for postgres health check to pass before connecting

#### Scenario: Frontend independent
- **WHEN** docker-compose starts
- **THEN** frontend can start independently of backend (calls API at runtime)

### Requirement: Local development setup
The system SHALL support easy local development.

#### Scenario: One command setup
- **WHEN** developer clones repository
- **THEN** can start entire stack with docker-compose up

#### Scenario: README with setup instructions
- **WHEN** developer reads README
- **THEN** clear step-by-step instructions for local setup are provided

#### Scenario: Port mappings for localhost
- **WHEN** services are running
- **THEN** frontend accessible at http://localhost, backend at http://localhost:8080

### Requirement: Log aggregation
The system SHALL make container logs easily accessible.

#### Scenario: View all logs
- **WHEN** developer runs docker-compose logs
- **THEN** logs from all services are displayed

#### Scenario: View service-specific logs
- **WHEN** developer runs docker-compose logs backend
- **THEN** only backend logs are displayed

#### Scenario: Follow logs in real-time
- **WHEN** developer runs docker-compose logs -f
- **THEN** logs stream in real-time as events occur

### Requirement: Volume backup capability
The system SHALL support backing up database volume.

#### Scenario: Volume can be backed up
- **WHEN** administrator needs to backup database
- **THEN** can use docker commands to export volume contents

#### Scenario: Volume can be restored
- **WHEN** administrator needs to restore database
- **THEN** can use docker commands to import volume contents

### Requirement: Production-ready configuration
The system SHALL support production deployment with configuration changes.

#### Scenario: Environment-based config
- **WHEN** deploying to production
- **THEN** can override environment variables for production settings

#### Scenario: CORS configuration
- **WHEN** running in production
- **THEN** backend restricts CORS to production frontend URL

#### Scenario: Database credentials secured
- **WHEN** deploying to production
- **THEN** database credentials are loaded from secure secrets management

### Requirement: Resource limits
The system SHALL define resource limits for containers.

#### Scenario: Memory limits
- **WHEN** containers run
- **THEN** each service has defined memory limit to prevent resource exhaustion

#### Scenario: CPU limits
- **WHEN** containers run under load
- **THEN** CPU usage is limited to prevent single service consuming all resources

#### Scenario: Limits configurable
- **WHEN** deploying to different environments
- **THEN** resource limits can be adjusted via environment variables
