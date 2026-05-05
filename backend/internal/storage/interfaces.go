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

// Ensure UserRepository implements the interface
var _ UserRepositoryInterface = (*UserRepository)(nil)
