package helpers

import (
	"testing"
	"time"
)

func TestTokenHelper(t *testing.T) {
	secretKey := "test_secret"
	otherSecret := "wrong_secret"
	now := time.Now().Unix()

	tokenHelper := NewTokenHelper(secretKey)

	tests := []struct {
		name        string
		secret      string
		email       string
		username    string
		createdAt   int64
		expireAfter time.Duration
		modifyToken func(string) string
		wantErr     bool
	}{
		{
			name:        "valid token",
			secret:      secretKey,
			email:       "test@example.com",
			username:    "testuser",
			createdAt:   now,
			expireAfter: 1 * time.Hour,
			wantErr:     false,
		},
		{
			name:        "expired token",
			secret:      secretKey,
			email:       "expired@example.com",
			username:    "olduser",
			createdAt:   now,
			expireAfter: -1 * time.Second,
			wantErr:     true,
		},
		{
			name:        "invalid token format",
			secret:      secretKey,
			modifyToken: func(_ string) string { return "not.a.jwt" },
			wantErr:     true,
		},
		{
			name:        "wrong secret key",
			secret:      otherSecret,
			email:       "wrong@example.com",
			username:    "wronguser",
			createdAt:   now,
			expireAfter: 1 * time.Hour,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate token only if we have email and username (skip for invalid format case)
			var token string
			var err error

			if tt.email != "" {
				// Always generate with the *correct* secret so we control the content
				token, err = tokenHelper.GenerateToken(tt.email, tt.username, tt.createdAt, tt.expireAfter)
				if err != nil {
					t.Fatalf("GenerateToken() error = %v", err)
				}
			}

			// If we have a modifyToken function, apply it
			if tt.modifyToken != nil {
				token = tt.modifyToken(token)
			}

			// Use the provided secret key for validation
			claims, err := NewTokenHelper(tt.secret).ValidateToken(token)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("Expected error but got none (claims: %+v)", claims)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Validate claims
			if claims.Email != tt.email {
				t.Errorf("Email mismatch: got %q, want %q", claims.Email, tt.email)
			}
			if claims.Username != tt.username {
				t.Errorf("Username mismatch: got %q, want %q", claims.Username, tt.username)
			}
			if claims.CreatedAt != tt.createdAt {
				t.Errorf("CreatedAt mismatch: got %d, want %d", claims.CreatedAt, tt.createdAt)
			}
		})
	}
}
