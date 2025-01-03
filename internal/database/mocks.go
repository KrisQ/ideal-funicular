package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// DBInterface defines just the methods we need in our handlers.
// *database.Queries already implements these, so no changes needed there.
type DBInterface interface {
	CreateChirp(ctx context.Context, arg CreateChirpParams) (Chirp, error)
	GetAllChirps(ctx context.Context) ([]Chirp, error)
	GetChirpById(ctx context.Context, id uuid.UUID) (Chirp, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	ResetUsers(ctx context.Context) error
}

// MockDB implements DBInterface, returning stubbed data or errors.
type MockDB struct {
	// You can store fields here that let you define behavior per test.
}

func (m *MockDB) CreateChirp(ctx context.Context, arg CreateChirpParams) (Chirp, error) {
	return Chirp{
		ID:        uuid.New(),
		Body:      arg.Body,
		UserID:    arg.UserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *MockDB) GetAllChirps(ctx context.Context) ([]Chirp, error) {
	// Return some sample data
	return []Chirp{
		{
			ID:        uuid.New(),
			Body:      "Hello World",
			UserID:    uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}, nil
}

// Implement other methods similarly...
func (m *MockDB) GetChirpById(ctx context.Context, id uuid.UUID) (Chirp, error) {
	return Chirp{
		ID:        id,
		Body:      "Test chirp",
		UserID:    uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *MockDB) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	return User{
		ID:             uuid.New(),
		Email:          arg.Email,
		HashedPassword: arg.HashedPassword,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

func (m *MockDB) GetUserByEmail(ctx context.Context, email string) (User, error) {
	return User{
		ID:             uuid.New(),
		Email:          email,
		HashedPassword: "fake_hash",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

func (m *MockDB) ResetUsers(ctx context.Context) error {
	return nil
}
