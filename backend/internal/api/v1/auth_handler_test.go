package v1

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rabin/tictactoe/internal/auth"
	"github.com/rabin/tictactoe/internal/storage"
)

// Mock UserRepository for testing
type mockUserRepo struct {
	users map[string]*storage.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users: make(map[string]*storage.User),
	}
}

func (m *mockUserRepo) Create(email, username, passwordHash string) (*storage.User, error) {
	// Check for duplicates
	for _, u := range m.users {
		if u.Email == email {
			return nil, sql.ErrNoRows // Simulate duplicate email error
		}
		if u.Username == username {
			return nil, sql.ErrNoRows // Simulate duplicate username error
		}
	}

	user := &storage.User{
		ID:           len(m.users) + 1,
		Email:        email,
		Username:     username,
		PasswordHash: passwordHash,
	}
	m.users[email] = user
	return user, nil
}

func (m *mockUserRepo) GetByEmail(email string) (*storage.User, error) {
	if user, ok := m.users[email]; ok {
		return user, nil
	}
	return nil, sql.ErrNoRows
}

func (m *mockUserRepo) GetByID(id int) (*storage.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (m *mockUserRepo) Update(user *storage.User) error {
	return nil
}

func (m *mockUserRepo) Delete(id int) error {
	return nil
}

func (m *mockUserRepo) UpdateStatistics(userID int, result string) error {
	return nil
}

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name: "Valid registration",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"username": "testuser",
				"password": "Password123",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["email"] != "test@example.com" {
					t.Errorf("Expected email test@example.com, got %v", resp["email"])
				}
				if resp["username"] != "testuser" {
					t.Errorf("Expected username testuser, got %v", resp["username"])
				}
			},
		},
		{
			name: "Invalid email",
			requestBody: map[string]string{
				"email":    "invalid-email",
				"username": "testuser",
				"password": "Password123",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				// May be VALIDATION_ERROR or INVALID_REQUEST depending on Gin binding
				code, _ := resp["code"].(string)
				if code != "VALIDATION_ERROR" && code != "INVALID_REQUEST" {
					t.Errorf("Expected VALIDATION_ERROR or INVALID_REQUEST code, got %s", code)
				}
			},
		},
		{
			name: "Username too short",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"username": "ab",
				"password": "Password123",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				// May be VALIDATION_ERROR or INVALID_REQUEST depending on Gin binding
				code, _ := resp["code"].(string)
				if code != "VALIDATION_ERROR" && code != "INVALID_REQUEST" {
					t.Errorf("Expected VALIDATION_ERROR or INVALID_REQUEST code, got %s", code)
				}
			},
		},
		{
			name: "Weak password",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"username": "testuser",
				"password": "weak",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				// May be VALIDATION_ERROR or INVALID_REQUEST depending on Gin binding
				code, _ := resp["code"].(string)
				if code != "VALIDATION_ERROR" && code != "INVALID_REQUEST" {
					t.Errorf("Expected VALIDATION_ERROR or INVALID_REQUEST code, got %s", code)
				}
			},
		},
		{
			name: "Missing fields",
			requestBody: map[string]string{
				"email": "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := newMockUserRepo()
			jwtService := auth.NewJWTService("test-secret")
			handler := NewAuthHandler(mockRepo, jwtService)

			router := gin.New()
			router.POST("/register", handler.Register)

			// Create request
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Execute
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse != nil && w.Code == tt.expectedStatus {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				tt.checkResponse(t, response)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup user in mock repo
	mockRepo := newMockUserRepo()
	hashedPassword, _ := auth.HashPassword("Password123")
	mockRepo.users["test@example.com"] = &storage.User{
		ID:           1,
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: hashedPassword,
	}

	jwtService := auth.NewJWTService("test-secret")
	handler := NewAuthHandler(mockRepo, jwtService)

	tests := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name: "Valid login",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"password": "Password123",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["token"] == nil || resp["token"] == "" {
					t.Error("Expected token in response")
				}
				if user, ok := resp["user"].(map[string]interface{}); ok {
					if user["email"] != "test@example.com" {
						t.Errorf("Expected email test@example.com")
					}
				}
			},
		},
		{
			name: "Invalid email",
			requestBody: map[string]string{
				"email":    "nonexistent@example.com",
				"password": "Password123",
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["code"] != "INVALID_CREDENTIALS" {
					t.Error("Expected INVALID_CREDENTIALS code")
				}
			},
		},
		{
			name: "Invalid password",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"password": "WrongPassword123",
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["code"] != "INVALID_CREDENTIALS" {
					t.Error("Expected INVALID_CREDENTIALS code")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/login", handler.Login)

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse != nil {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				tt.checkResponse(t, response)
			}
		})
	}
}

func TestGetCurrentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	mockRepo := newMockUserRepo()
	mockRepo.users["test@example.com"] = &storage.User{
		ID:       1,
		Email:    "test@example.com",
		Username: "testuser",
	}

	jwtService := auth.NewJWTService("test-secret")
	handler := NewAuthHandler(mockRepo, jwtService)

	// Generate valid token
	token, _ := jwtService.GenerateToken(1, "test@example.com")

	tests := []struct {
		name           string
		token          string
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:           "Valid token",
			token:          token,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				if resp["email"] != "test@example.com" {
					t.Error("Expected user email in response")
				}
			},
		},
		{
			name:           "Missing token",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid token",
			token:          "invalid-token",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/me", auth.AuthMiddleware(jwtService), handler.GetCurrentUser)

			req, _ := http.NewRequest(http.MethodGet, "/me", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse != nil && w.Code == tt.expectedStatus {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				tt.checkResponse(t, response)
			}
		})
	}
}
