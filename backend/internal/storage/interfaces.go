package storage

// UserRepositoryInterface defines the interface for user repository operations
type UserRepositoryInterface interface {
	Create(email, username, passwordHash string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByID(id int) (*User, error)
	Update(user *User) error
	Delete(id int) error
	UpdateStatistics(userID int, result string) error
}

// GameRepositoryInterface defines the interface for game repository operations
type GameRepositoryInterface interface {
	Create(userID int, difficulty string, boardState []byte) (*Game, error)
	GetByID(id int) (*Game, error)
	Update(game *Game) error
	ListByUserID(userID int) ([]Game, error)
	AddMove(gameID, position int, player string, moveNumber int) error
	GetMoves(gameID int) ([]GameMove, error)
}

// Ensure repositories implement their interfaces
var _ UserRepositoryInterface = (*UserRepository)(nil)
var _ GameRepositoryInterface = (*GameRepository)(nil)
