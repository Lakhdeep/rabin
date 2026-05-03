## ADDED Requirements

### Requirement: User registration with email and username
The system SHALL allow users to register with a unique email address and username. Both email and username MUST be unique across all users.

#### Scenario: Successful registration
- **WHEN** user submits valid registration form with email, username, and password
- **THEN** system creates new user account and returns success confirmation

#### Scenario: Duplicate email
- **WHEN** user attempts to register with an email that already exists
- **THEN** system rejects registration and returns error "Email already registered"

#### Scenario: Duplicate username
- **WHEN** user attempts to register with a username that already exists
- **THEN** system rejects registration and returns error "Username already taken"

#### Scenario: Invalid email format
- **WHEN** user submits registration with malformed email address
- **THEN** system rejects registration and returns error "Invalid email format"

### Requirement: Username validation
The system SHALL enforce username constraints to ensure consistency and prevent abuse.

#### Scenario: Valid username
- **WHEN** user submits username with 3-20 alphanumeric characters
- **THEN** system accepts the username

#### Scenario: Username too short
- **WHEN** user submits username with fewer than 3 characters
- **THEN** system rejects registration and returns error "Username must be at least 3 characters"

#### Scenario: Username too long
- **WHEN** user submits username with more than 20 characters
- **THEN** system rejects registration and returns error "Username must be at most 20 characters"

#### Scenario: Username with invalid characters
- **WHEN** user submits username containing special characters or spaces
- **THEN** system rejects registration and returns error "Username can only contain letters and numbers"

### Requirement: Password security requirements
The system SHALL enforce password strength requirements to protect user accounts.

#### Scenario: Valid password
- **WHEN** user submits password with at least 8 characters including 1 uppercase, 1 lowercase, and 1 number
- **THEN** system accepts the password

#### Scenario: Password too short
- **WHEN** user submits password with fewer than 8 characters
- **THEN** system rejects registration and returns error "Password must be at least 8 characters"

#### Scenario: Password missing uppercase
- **WHEN** user submits password without uppercase letter
- **THEN** system rejects registration and returns error "Password must contain at least one uppercase letter"

#### Scenario: Password missing lowercase
- **WHEN** user submits password without lowercase letter
- **THEN** system rejects registration and returns error "Password must contain at least one lowercase letter"

#### Scenario: Password missing number
- **WHEN** user submits password without numeric digit
- **THEN** system rejects registration and returns error "Password must contain at least one number"

### Requirement: Password hashing with bcrypt
The system SHALL hash all passwords using bcrypt with cost factor 12 before storing in database.

#### Scenario: Password stored securely
- **WHEN** user registers with valid credentials
- **THEN** system stores bcrypt hash of password, not plaintext

#### Scenario: Password verification
- **WHEN** user logs in with correct password
- **THEN** system verifies password against stored bcrypt hash

### Requirement: User login with credentials
The system SHALL allow registered users to login with email and password.

#### Scenario: Successful login
- **WHEN** user submits correct email and password
- **THEN** system returns JWT token and user profile information

#### Scenario: Invalid email
- **WHEN** user submits email that does not exist in database
- **THEN** system rejects login and returns error "Invalid credentials"

#### Scenario: Incorrect password
- **WHEN** user submits correct email but wrong password
- **THEN** system rejects login and returns error "Invalid credentials"

#### Scenario: Empty credentials
- **WHEN** user submits empty email or password
- **THEN** system rejects login and returns error "Email and password are required"

### Requirement: JWT token generation
The system SHALL generate JWT tokens for authenticated users using HS256 algorithm.

#### Scenario: Token contains user information
- **WHEN** user successfully logs in
- **THEN** system generates JWT token containing user_id, email, and expiration timestamp

#### Scenario: Token expiration
- **WHEN** JWT token is generated
- **THEN** token expires after 24 hours from creation time

#### Scenario: Token signed with secret
- **WHEN** JWT token is generated
- **THEN** token is signed with server's secret key using HS256 algorithm

### Requirement: JWT token validation
The system SHALL validate JWT tokens on protected API endpoints.

#### Scenario: Valid token grants access
- **WHEN** user includes valid JWT token in Authorization header
- **THEN** system allows access to protected endpoint

#### Scenario: Missing token denied
- **WHEN** user makes request to protected endpoint without token
- **THEN** system returns 401 Unauthorized error

#### Scenario: Expired token denied
- **WHEN** user includes expired JWT token in request
- **THEN** system returns 401 Unauthorized error with message "Token expired"

#### Scenario: Invalid signature denied
- **WHEN** user includes JWT token with invalid signature
- **THEN** system returns 401 Unauthorized error with message "Invalid token"

#### Scenario: Malformed token denied
- **WHEN** user includes malformed JWT token in request
- **THEN** system returns 401 Unauthorized error with message "Invalid token format"

### Requirement: Get current user profile
The system SHALL allow authenticated users to retrieve their own profile information.

#### Scenario: Retrieve own profile
- **WHEN** authenticated user requests their profile
- **THEN** system returns user's email, username, and statistics (total_games, wins, losses, draws)

#### Scenario: Unauthenticated access denied
- **WHEN** unauthenticated user requests profile endpoint
- **THEN** system returns 401 Unauthorized error
