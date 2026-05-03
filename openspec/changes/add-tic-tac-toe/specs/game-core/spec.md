## ADDED Requirements

### Requirement: Game board representation
The system SHALL represent the tic-tac-toe board as a 3x3 grid with positions numbered 0-8.

#### Scenario: Empty board initialization
- **WHEN** new game is created
- **THEN** all 9 positions are empty (null)

#### Scenario: Position mapping
- **WHEN** referencing board positions
- **THEN** positions are numbered 0-8 from top-left to bottom-right (0=top-left, 4=center, 8=bottom-right)

### Requirement: Player assignment
The system SHALL assign X to the human player and O to the AI opponent.

#### Scenario: User plays as X
- **WHEN** new game starts
- **THEN** user is assigned player X and makes the first move

#### Scenario: AI plays as O
- **WHEN** new game starts
- **THEN** AI is assigned player O and moves second

### Requirement: Move validation
The system SHALL validate all moves before updating board state.

#### Scenario: Valid move accepted
- **WHEN** player selects empty position during their turn
- **THEN** system places player's mark at that position

#### Scenario: Occupied position rejected
- **WHEN** player selects position already containing X or O
- **THEN** system rejects move and returns error "Position already occupied"

#### Scenario: Move during opponent turn rejected
- **WHEN** player attempts move during AI's turn
- **THEN** system rejects move and returns error "Not your turn"

#### Scenario: Move after game ended rejected
- **WHEN** player attempts move after game has ended (win/loss/draw)
- **THEN** system rejects move and returns error "Game already completed"

#### Scenario: Invalid position rejected
- **WHEN** player selects position outside range 0-8
- **THEN** system rejects move and returns error "Invalid position"

### Requirement: Win condition detection
The system SHALL detect win conditions after each move.

#### Scenario: Horizontal win detected
- **WHEN** player marks three consecutive positions in same row (0-1-2, 3-4-5, or 6-7-8)
- **THEN** system declares that player as winner and ends game

#### Scenario: Vertical win detected
- **WHEN** player marks three consecutive positions in same column (0-3-6, 1-4-7, or 2-5-8)
- **THEN** system declares that player as winner and ends game

#### Scenario: Diagonal win detected (top-left to bottom-right)
- **WHEN** player marks positions 0-4-8
- **THEN** system declares that player as winner and ends game

#### Scenario: Diagonal win detected (top-right to bottom-left)
- **WHEN** player marks positions 2-4-6
- **THEN** system declares that player as winner and ends game

#### Scenario: No win condition continues game
- **WHEN** move does not complete three in a row
- **THEN** system continues game and switches turn to next player

### Requirement: Draw condition detection
The system SHALL detect draw conditions when board is full with no winner.

#### Scenario: Draw when board full
- **WHEN** all 9 positions are filled and no player has three in a row
- **THEN** system declares game as draw and ends game

#### Scenario: Board not full continues game
- **WHEN** empty positions remain and no winner exists
- **THEN** system continues game

### Requirement: AI difficulty - Easy
The system SHALL implement Easy AI that selects random valid moves.

#### Scenario: Easy AI makes random move
- **WHEN** Easy AI's turn arrives
- **THEN** AI randomly selects one empty position and places O

#### Scenario: Easy AI only chooses valid positions
- **WHEN** Easy AI calculates move
- **THEN** AI only considers empty positions (never occupied ones)

### Requirement: AI difficulty - Medium
The system SHALL implement Medium AI that uses optimal strategy 50% of the time and random moves 50% of the time.

#### Scenario: Medium AI uses minimax sometimes
- **WHEN** Medium AI's turn arrives
- **THEN** AI has 50% probability of using minimax algorithm for optimal move

#### Scenario: Medium AI uses random sometimes
- **WHEN** Medium AI's turn arrives
- **THEN** AI has 50% probability of selecting random valid move

### Requirement: AI difficulty - Hard
The system SHALL implement Hard AI using minimax algorithm with alpha-beta pruning limited to depth 4.

#### Scenario: Hard AI uses minimax with depth limit
- **WHEN** Hard AI calculates move
- **THEN** AI uses minimax algorithm with alpha-beta pruning searching up to depth 4

#### Scenario: Hard AI makes near-optimal moves
- **WHEN** Hard AI's turn arrives
- **THEN** AI selects move that maximizes its winning probability within depth limit

#### Scenario: Hard AI is beatable
- **WHEN** user plays optimally against Hard AI
- **THEN** user can achieve draw or win due to depth limitation

### Requirement: AI difficulty - Impossible
The system SHALL implement Impossible AI using full minimax algorithm with alpha-beta pruning.

#### Scenario: Impossible AI uses full minimax
- **WHEN** Impossible AI calculates move
- **THEN** AI uses minimax algorithm with alpha-beta pruning searching entire game tree

#### Scenario: Impossible AI plays optimally
- **WHEN** Impossible AI's turn arrives
- **THEN** AI selects the mathematically optimal move

#### Scenario: Impossible AI never loses
- **WHEN** user plays against Impossible AI
- **THEN** best possible outcome for user is a draw (AI cannot lose)

### Requirement: AI response time
The system SHALL ensure AI moves complete within reasonable time.

#### Scenario: AI move completes quickly
- **WHEN** AI calculates move at any difficulty level
- **THEN** calculation completes within 1 second

### Requirement: Turn management
The system SHALL track whose turn it is and enforce alternating turns.

#### Scenario: User starts first
- **WHEN** new game begins
- **THEN** user (X) has the first turn

#### Scenario: Turns alternate
- **WHEN** player makes valid move
- **THEN** turn switches to opponent (AI)

#### Scenario: Turn persists until valid move
- **WHEN** player attempts invalid move
- **THEN** turn remains with same player

### Requirement: Game state persistence
The system SHALL maintain complete game state throughout the game.

#### Scenario: Board state stored
- **WHEN** move is made
- **THEN** updated board state is saved to database

#### Scenario: Current turn stored
- **WHEN** turn switches
- **THEN** current turn indicator is saved to database

#### Scenario: Game result stored
- **WHEN** game ends
- **THEN** result (win/loss/draw) is saved to database

### Requirement: Move history tracking
The system SHALL record every move made during the game.

#### Scenario: Move recorded with details
- **WHEN** player or AI makes move
- **THEN** system saves position, player (X/O), and move number to game_moves table

#### Scenario: Moves ordered sequentially
- **WHEN** retrieving game history
- **THEN** moves are returned in chronological order by move_number

#### Scenario: Move timestamp recorded
- **WHEN** move is made
- **THEN** system records timestamp of when move occurred
