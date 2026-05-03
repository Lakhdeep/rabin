## ADDED Requirements

### Requirement: User statistics initialization
The system SHALL initialize user statistics to zero when user registers.

#### Scenario: New user has zero stats
- **WHEN** user completes registration
- **THEN** user's total_games, wins, losses, and draws are all set to 0

### Requirement: Win tracking
The system SHALL increment user's win count when they win a game.

#### Scenario: User wins against AI
- **WHEN** user achieves three in a row before AI
- **THEN** system increments user's wins by 1 and total_games by 1

#### Scenario: Win count persists
- **WHEN** user wins multiple games
- **THEN** each win increments the counter independently

### Requirement: Loss tracking
The system SHALL increment user's loss count when they lose a game.

#### Scenario: User loses to AI
- **WHEN** AI achieves three in a row before user
- **THEN** system increments user's losses by 1 and total_games by 1

#### Scenario: Loss count persists
- **WHEN** user loses multiple games
- **THEN** each loss increments the counter independently

### Requirement: Draw tracking
The system SHALL increment user's draw count when game ends in draw.

#### Scenario: Game ends in draw
- **WHEN** board fills completely with no winner
- **THEN** system increments user's draws by 1 and total_games by 1

#### Scenario: Draw count persists
- **WHEN** user plays multiple drawn games
- **THEN** each draw increments the counter independently

### Requirement: Total games tracking
The system SHALL maintain accurate count of total games played.

#### Scenario: Total equals sum of outcomes
- **WHEN** user completes any game (win/loss/draw)
- **THEN** total_games equals (wins + losses + draws)

#### Scenario: Abandoned games not counted
- **WHEN** user starts game but never completes it
- **THEN** total_games is not incremented

### Requirement: Atomic statistics update
The system SHALL update statistics atomically when game ends.

#### Scenario: Statistics updated together
- **WHEN** game ends with result
- **THEN** total_games and outcome counter (wins/losses/draws) are updated in single transaction

#### Scenario: Update failure rolls back
- **WHEN** statistics update fails for any reason
- **THEN** neither total_games nor outcome counter is incremented

### Requirement: Game result persistence
The system SHALL store individual game records with result and difficulty.

#### Scenario: Game saved with result
- **WHEN** game ends
- **THEN** system saves game record with user_id, difficulty, result (win/loss/draw), and timestamp

#### Scenario: Game saved with difficulty level
- **WHEN** game ends
- **THEN** game record includes which difficulty level was played (easy/medium/hard/impossible)

#### Scenario: Completed timestamp recorded
- **WHEN** game ends
- **THEN** system records completed_at timestamp

### Requirement: User statistics retrieval
The system SHALL allow users to view their complete statistics.

#### Scenario: Retrieve statistics
- **WHEN** user requests their statistics
- **THEN** system returns total_games, wins, losses, draws, and win_rate

#### Scenario: Win rate calculation
- **WHEN** user has played at least one game
- **THEN** win_rate equals (wins / total_games) expressed as percentage

#### Scenario: Win rate for new user
- **WHEN** user has played zero games
- **THEN** win_rate is 0.0 or null

### Requirement: Game history retrieval
The system SHALL allow users to view their past games.

#### Scenario: List user's games
- **WHEN** user requests game history
- **THEN** system returns list of games with difficulty, result, and timestamp

#### Scenario: Games ordered by recency
- **WHEN** retrieving game history
- **THEN** games are ordered by completed_at timestamp (most recent first)

#### Scenario: Only completed games shown
- **WHEN** retrieving game history
- **THEN** only games with result (win/loss/draw) are included, not in-progress games

### Requirement: Leaderboard ranking
The system SHALL provide leaderboard of top players.

#### Scenario: Leaderboard sorted by wins
- **WHEN** user views leaderboard
- **THEN** users are ordered by wins (descending), then by win_rate (descending)

#### Scenario: Leaderboard shows statistics
- **WHEN** user views leaderboard
- **THEN** each entry shows username, total_games, wins, losses, draws, and win_rate

#### Scenario: Leaderboard limited to active players
- **WHEN** generating leaderboard
- **THEN** only users with at least one completed game are included

### Requirement: Statistics accuracy
The system SHALL ensure statistics remain accurate and consistent.

#### Scenario: No negative statistics
- **WHEN** accessing user statistics
- **THEN** all counters (wins, losses, draws, total_games) are non-negative integers

#### Scenario: Total games consistency
- **WHEN** accessing user statistics
- **THEN** total_games equals count of user's completed games in games table

#### Scenario: Statistics survive server restart
- **WHEN** server restarts
- **THEN** all user statistics persist unchanged in database
