## ADDED Requirements

### Requirement: User registration page
The system SHALL provide registration page for new users.

#### Scenario: Registration form fields
- **WHEN** user navigates to registration page
- **THEN** form displays inputs for email, username, and password

#### Scenario: Registration form validation
- **WHEN** user submits form with invalid data
- **THEN** form shows validation errors without making API call

#### Scenario: Successful registration redirect
- **WHEN** registration succeeds
- **THEN** user is redirected to login page with success message

#### Scenario: Registration error display
- **WHEN** registration fails (duplicate email/username)
- **THEN** error message is displayed near relevant field

### Requirement: User login page
The system SHALL provide login page for authentication.

#### Scenario: Login form fields
- **WHEN** user navigates to login page
- **THEN** form displays inputs for email and password

#### Scenario: Successful login redirect
- **WHEN** login succeeds
- **THEN** user is redirected to dashboard and JWT token is stored

#### Scenario: Login error display
- **WHEN** login fails
- **THEN** error message "Invalid credentials" is displayed

#### Scenario: Link to registration
- **WHEN** user is on login page
- **THEN** link to registration page is visible

### Requirement: JWT token storage
The system SHALL store JWT token in browser localStorage.

#### Scenario: Token saved on login
- **WHEN** user logs in successfully
- **THEN** JWT token is saved to localStorage

#### Scenario: Token included in API requests
- **WHEN** making authenticated API request
- **THEN** token is included in Authorization header as "Bearer <token>"

#### Scenario: Token cleared on logout
- **WHEN** user logs out
- **THEN** JWT token is removed from localStorage

### Requirement: Protected routes
The system SHALL restrict access to authenticated pages.

#### Scenario: Unauthenticated user redirected
- **WHEN** unauthenticated user navigates to protected route
- **THEN** user is redirected to login page

#### Scenario: Authenticated user accesses protected route
- **WHEN** authenticated user navigates to protected route
- **THEN** page content is displayed normally

### Requirement: Navigation menu
The system SHALL provide navigation for authenticated users.

#### Scenario: Navigation items visible
- **WHEN** user is authenticated
- **THEN** navigation shows links to Dashboard, New Game, Profile, and Logout

#### Scenario: Current page highlighted
- **WHEN** user is on a page
- **THEN** corresponding navigation item is visually highlighted

#### Scenario: Logout button
- **WHEN** user clicks logout
- **THEN** token is cleared and user is redirected to login page

### Requirement: Dashboard page
The system SHALL provide dashboard showing user statistics.

#### Scenario: Statistics displayed
- **WHEN** user navigates to dashboard
- **THEN** page shows total_games, wins, losses, draws, and win_rate

#### Scenario: Win rate formatting
- **WHEN** displaying win rate
- **THEN** value is shown as percentage with 1 decimal place (e.g., "65.5%")

#### Scenario: New game button
- **WHEN** user is on dashboard
- **THEN** prominent "New Game" button is visible

#### Scenario: Recent games list
- **WHEN** user has played games
- **THEN** dashboard shows list of recent games with result and difficulty

### Requirement: Game board component
The system SHALL provide interactive 3x3 game board.

#### Scenario: Board grid layout
- **WHEN** game page loads
- **THEN** board displays 3x3 grid of clickable cells

#### Scenario: Empty cells clickable
- **WHEN** cell is empty and it's user's turn
- **THEN** cell shows hover effect and is clickable

#### Scenario: Occupied cells not clickable
- **WHEN** cell contains X or O
- **THEN** cell is not clickable and shows no hover effect

#### Scenario: Cell click makes move
- **WHEN** user clicks empty cell during their turn
- **THEN** X appears in that cell and API request is made

#### Scenario: AI move displayed
- **WHEN** AI makes move
- **THEN** O appears in AI's chosen cell

#### Scenario: Winning line highlighted
- **WHEN** game ends with winner
- **THEN** three winning cells are visually highlighted

### Requirement: Difficulty selector
The system SHALL allow users to select AI difficulty before game.

#### Scenario: Difficulty options
- **WHEN** user starts new game
- **THEN** dropdown or buttons show options: Easy, Medium, Hard, Impossible

#### Scenario: Default difficulty
- **WHEN** difficulty selector appears
- **THEN** Medium is selected by default

#### Scenario: Difficulty locked during game
- **WHEN** game is in progress
- **THEN** difficulty cannot be changed

### Requirement: Game status display
The system SHALL show current game status to user.

#### Scenario: Turn indicator
- **WHEN** game is in progress
- **THEN** text indicates "Your turn" or "AI is thinking..."

#### Scenario: Win message
- **WHEN** user wins
- **THEN** message displays "You won!"

#### Scenario: Loss message
- **WHEN** AI wins
- **THEN** message displays "AI wins!"

#### Scenario: Draw message
- **WHEN** game ends in draw
- **THEN** message displays "It's a draw!"

#### Scenario: New game button after completion
- **WHEN** game ends
- **THEN** "Play Again" button is displayed

### Requirement: Loading states
The system SHALL provide visual feedback during API requests.

#### Scenario: Loading spinner on request
- **WHEN** API request is in progress
- **THEN** loading spinner or indicator is displayed

#### Scenario: Disabled state during loading
- **WHEN** waiting for API response
- **THEN** form inputs or game board are disabled

#### Scenario: Loading cleared on response
- **WHEN** API request completes (success or error)
- **THEN** loading indicator is hidden

### Requirement: Error handling
The system SHALL display user-friendly error messages.

#### Scenario: Network error message
- **WHEN** API request fails due to network issue
- **THEN** error toast displays "Connection error. Please try again."

#### Scenario: Validation error message
- **WHEN** API returns validation error
- **THEN** specific error message is displayed near relevant field

#### Scenario: Generic error message
- **WHEN** unexpected error occurs
- **THEN** error toast displays "Something went wrong. Please try again."

#### Scenario: Error dismissible
- **WHEN** error message is displayed
- **THEN** user can dismiss it by clicking X or waiting for auto-dismiss

### Requirement: Responsive design
The system SHALL adapt layout for different screen sizes.

#### Scenario: Mobile layout (< 768px)
- **WHEN** viewport width is less than 768px
- **THEN** layout switches to single-column mobile view

#### Scenario: Tablet layout (768-1024px)
- **WHEN** viewport width is between 768px and 1024px
- **THEN** layout adjusts to tablet-friendly spacing

#### Scenario: Desktop layout (> 1024px)
- **WHEN** viewport width exceeds 1024px
- **THEN** full desktop layout with optimal use of space

#### Scenario: Game board responsive
- **WHEN** screen size changes
- **THEN** game board scales proportionally while maintaining square cells

### Requirement: Profile page
The system SHALL provide user profile page.

#### Scenario: Profile information displayed
- **WHEN** user navigates to profile page
- **THEN** page shows email, username, and account creation date

#### Scenario: Detailed statistics
- **WHEN** user views profile
- **THEN** page shows complete statistics breakdown by difficulty level

#### Scenario: Game history
- **WHEN** user views profile
- **THEN** page shows chronological list of all completed games

### Requirement: Accessibility
The system SHALL support keyboard navigation and screen readers.

#### Scenario: Keyboard navigation
- **WHEN** user navigates with Tab key
- **THEN** focus moves through interactive elements in logical order

#### Scenario: ARIA labels for game board
- **WHEN** screen reader user accesses game board
- **THEN** cells have descriptive labels (e.g., "Position 0, empty")

#### Scenario: Focus visible
- **WHEN** element receives keyboard focus
- **THEN** clear visual focus indicator is displayed

### Requirement: API client service
The system SHALL centralize API communication logic.

#### Scenario: Base URL configuration
- **WHEN** making API requests
- **THEN** client uses configured base URL from environment variable

#### Scenario: Token automatically included
- **WHEN** making authenticated request
- **THEN** API client automatically includes JWT from localStorage

#### Scenario: Error handling centralized
- **WHEN** API request fails
- **THEN** client handles common errors (401, 500) consistently

### Requirement: Authentication context
The system SHALL provide global authentication state.

#### Scenario: User state available globally
- **WHEN** any component needs current user information
- **THEN** component can access user from AuthContext

#### Scenario: Login updates context
- **WHEN** user logs in
- **THEN** AuthContext is updated with user information

#### Scenario: Logout clears context
- **WHEN** user logs out
- **THEN** AuthContext user state is cleared

### Requirement: Game context
The system SHALL provide global game state management.

#### Scenario: Game state available to components
- **WHEN** game components need board state
- **THEN** components can access game state from GameContext

#### Scenario: Move updates context
- **WHEN** user makes move
- **THEN** GameContext is updated with new board state

#### Scenario: New game resets context
- **WHEN** user starts new game
- **THEN** GameContext is reset to initial state

### Requirement: Form validation
The system SHALL validate user inputs before submission.

#### Scenario: Real-time validation
- **WHEN** user types in form field
- **THEN** validation runs on blur or after short delay

#### Scenario: Email validation
- **WHEN** user enters email
- **THEN** system validates email format using regex

#### Scenario: Password strength indicator
- **WHEN** user types password
- **THEN** visual indicator shows password strength (weak/medium/strong)

#### Scenario: Submit disabled when invalid
- **WHEN** form has validation errors
- **THEN** submit button is disabled

### Requirement: Color scheme and contrast
The system SHALL use accessible color combinations.

#### Scenario: Text contrast ratio
- **WHEN** text is displayed on background
- **THEN** contrast ratio meets WCAG AA standard (4.5:1 minimum)

#### Scenario: Distinct player colors
- **WHEN** displaying X and O on board
- **THEN** colors are clearly distinguishable for colorblind users

#### Scenario: Error colors
- **WHEN** displaying error messages
- **THEN** red color has sufficient contrast with background
