package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken(t *testing.T) {
	jwtService := NewJWTService("test-secret-key")
	userID := 123
	email := "test@example.com"

	token, err := jwtService.GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Error("Expected token to be non-empty")
	}

	// Verify it's a valid JWT format (3 parts separated by dots)
	parts := 0
	for _, char := range token {
		if char == '.' {
			parts++
		}
	}
	if parts != 2 {
		t.Errorf("Expected JWT to have 2 dots (3 parts), got %d dots", parts)
	}
}

func TestValidateToken(t *testing.T) {
	jwtService := NewJWTService("test-secret-key")
	userID := 456
	email := "user@example.com"

	token, _ := jwtService.GenerateToken(userID, email)

	claims, err := jwtService.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, claims.UserID)
	}

	if claims.Email != email {
		t.Errorf("Expected Email %s, got %s", email, claims.Email)
	}
}

func TestValidateTokenInvalidSecret(t *testing.T) {
	// Generate token with one secret
	jwtService1 := NewJWTService("secret-one")
	token, _ := jwtService1.GenerateToken(123, "test@example.com")

	// Try to validate with different secret
	jwtService2 := NewJWTService("secret-two")
	_, err := jwtService2.ValidateToken(token)

	if err == nil {
		t.Error("Expected error when validating token with wrong secret")
	}
}

func TestValidateTokenMalformed(t *testing.T) {
	jwtService := NewJWTService("test-secret")

	tests := []struct {
		name  string
		token string
	}{
		{"Empty token", ""},
		{"Invalid format", "not.a.valid.jwt"},
		{"Random string", "random-string"},
		{"Only header", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := jwtService.ValidateToken(tt.token)
			if err == nil {
				t.Error("Expected error for malformed token")
			}
		})
	}
}

func TestValidateTokenExpired(t *testing.T) {
	jwtService := NewJWTService("test-secret")

	// Create an expired token (manually)
	claims := Claims{
		UserID: 123,
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-25 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtService.secretKey)

	// Try to validate expired token
	_, err := jwtService.ValidateToken(tokenString)
	if err == nil {
		t.Error("Expected error for expired token")
	}
}

func TestTokenExpiration(t *testing.T) {
	jwtService := NewJWTService("test-secret")
	token, _ := jwtService.GenerateToken(123, "test@example.com")

	claims, _ := jwtService.ValidateToken(token)

	// Check that expiration is ~24 hours from now
	expectedExpiry := time.Now().Add(24 * time.Hour)
	actualExpiry := claims.ExpiresAt.Time

	diff := actualExpiry.Sub(expectedExpiry)
	if diff < -1*time.Minute || diff > 1*time.Minute {
		t.Errorf("Expected expiry ~24 hours from now, got %v (diff: %v)", actualExpiry, diff)
	}
}

func TestTokenIssuedAt(t *testing.T) {
	jwtService := NewJWTService("test-secret")
	beforeIssue := time.Now().Add(-1 * time.Second) // Allow 1 second buffer
	token, _ := jwtService.GenerateToken(123, "test@example.com")
	afterIssue := time.Now().Add(1 * time.Second)

	claims, _ := jwtService.ValidateToken(token)

	if claims.IssuedAt.Time.Before(beforeIssue) || claims.IssuedAt.Time.After(afterIssue) {
		t.Errorf("Expected IssuedAt between %v and %v, got %v", beforeIssue, afterIssue, claims.IssuedAt.Time)
	}
}

func TestSigningMethodHS256(t *testing.T) {
	jwtService := NewJWTService("test-secret")
	tokenString, _ := jwtService.GenerateToken(123, "test@example.com")

	// Parse token without validating to check algorithm
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtService.secretKey, nil
	})

	if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		t.Error("Expected signing method to be HMAC")
	} else if method.Alg() != "HS256" {
		t.Errorf("Expected algorithm HS256, got %s", method.Alg())
	}
}
