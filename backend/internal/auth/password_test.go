package auth

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := "TestPassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash == "" {
		t.Error("Expected hash to be non-empty")
	}

	if hash == password {
		t.Error("Hash should not equal plain password")
	}

	// Verify it's a valid bcrypt hash
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		t.Errorf("Generated hash is not valid bcrypt: %v", err)
	}
}

func TestHashPasswordCost(t *testing.T) {
	password := "TestPassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// Extract cost from hash
	cost, err := bcrypt.Cost([]byte(hash))
	if err != nil {
		t.Fatalf("Failed to get bcrypt cost: %v", err)
	}

	if cost != 12 {
		t.Errorf("Expected bcrypt cost 12, got %d", cost)
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "CorrectPassword123"
	hash, _ := HashPassword(password)

	tests := []struct {
		name        string
		password    string
		shouldError bool
	}{
		{
			name:        "Correct password",
			password:    "CorrectPassword123",
			shouldError: false,
		},
		{
			name:        "Wrong password",
			password:    "WrongPassword123",
			shouldError: true,
		},
		{
			name:        "Empty password",
			password:    "",
			shouldError: true,
		},
		{
			name:        "Case sensitive",
			password:    "correctpassword123",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyPassword(hash, tt.password)
			if tt.shouldError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestHashPasswordDifferentHashes(t *testing.T) {
	password := "SamePassword123"

	hash1, _ := HashPassword(password)
	hash2, _ := HashPassword(password)

	// Bcrypt should generate different hashes for same password (due to salt)
	if hash1 == hash2 {
		t.Error("Expected different hashes for same password (bcrypt uses random salt)")
	}

	// But both should verify correctly
	if err := VerifyPassword(hash1, password); err != nil {
		t.Errorf("First hash failed verification: %v", err)
	}
	if err := VerifyPassword(hash2, password); err != nil {
		t.Errorf("Second hash failed verification: %v", err)
	}
}
