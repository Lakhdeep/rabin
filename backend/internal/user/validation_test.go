package user

import "testing"

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		shouldError bool
	}{
		{"Valid email", "user@example.com", false},
		{"Valid email with subdomain", "user@mail.example.com", false},
		{"Valid email with plus", "user+tag@example.com", false},
		{"Valid email with dash", "user-name@example.com", false},
		{"Valid email with numbers", "user123@example.com", false},
		{"Empty email", "", true},
		{"Missing @", "userexample.com", true},
		{"Missing domain", "user@", true},
		{"Missing username", "@example.com", true},
		{"Missing TLD", "user@example", true},
		{"Invalid characters", "user name@example.com", true},
		{"Multiple @", "user@@example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if tt.shouldError && err == nil {
				t.Errorf("Expected error for email: %s", tt.email)
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error for email: %s, got: %v", tt.email, err)
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		shouldError bool
		errorMsg    string
	}{
		{"Valid username", "validuser123", false, ""},
		{"Valid 3 chars", "abc", false, ""},
		{"Valid 20 chars", "12345678901234567890", false, ""},
		{"Empty username", "", true, "username is required"},
		{"Too short", "ab", true, "at least 3 characters"},
		{"Too long", "123456789012345678901", true, "at most 20 characters"},
		{"With spaces", "user name", true, "only contain letters and numbers"},
		{"With special chars", "user@name", true, "only contain letters and numbers"},
		{"With dash", "user-name", true, "only contain letters and numbers"},
		{"With underscore", "user_name", true, "only contain letters and numbers"},
		{"Only numbers", "123456", false, ""},
		{"Only letters", "username", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.username)
			if tt.shouldError && err == nil {
				t.Errorf("Expected error for username: %s", tt.username)
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error for username: %s, got: %v", tt.username, err)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		shouldError bool
		errorMsg    string
	}{
		{"Valid password", "Password123", false, ""},
		{"Valid 8 chars", "Pass123A", false, ""},
		{"Valid long password", "VeryLongPassword123WithManyChars", false, ""},
		{"Empty password", "", true, "password is required"},
		{"Too short", "Pass12", true, "at least 8 characters"},
		{"No uppercase", "password123", true, "at least one uppercase letter"},
		{"No lowercase", "PASSWORD123", true, "at least one lowercase letter"},
		{"No number", "PasswordABC", true, "at least one number"},
		{"Only lowercase", "abcdefgh", true, ""},
		{"Only uppercase", "ABCDEFGH", true, ""},
		{"Only numbers", "12345678", true, ""},
		{"With special chars", "Pass@123", false, ""}, // Special chars allowed but not required
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if tt.shouldError && err == nil {
				t.Errorf("Expected error for password: %s", tt.password)
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error for password: %s, got: %v", tt.password, err)
			}
		})
	}
}

func TestValidateRegistration(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		username    string
		password    string
		shouldError bool
	}{
		{
			name:        "All valid",
			email:       "user@example.com",
			username:    "validuser",
			password:    "Password123",
			shouldError: false,
		},
		{
			name:        "Invalid email",
			email:       "invalid-email",
			username:    "validuser",
			password:    "Password123",
			shouldError: true,
		},
		{
			name:        "Invalid username",
			email:       "user@example.com",
			username:    "ab",
			password:    "Password123",
			shouldError: true,
		},
		{
			name:        "Invalid password",
			email:       "user@example.com",
			username:    "validuser",
			password:    "weak",
			shouldError: true,
		},
		{
			name:        "All invalid",
			email:       "",
			username:    "",
			password:    "",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRegistration(tt.email, tt.username, tt.password)
			if tt.shouldError && err == nil {
				t.Error("Expected validation error")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
		})
	}
}
