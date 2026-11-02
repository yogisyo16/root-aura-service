package services

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test-secret"

// signWith builds a token by hand so tests can produce tokens
// GenerateToken never would (expired, no exp, other algorithms).
func signWith(t *testing.T, method jwt.SigningMethod, key any, claims jwt.RegisteredClaims) string {
	t.Helper()
	s, err := jwt.NewWithClaims(method, claims).SignedString(key)
	if err != nil {
		t.Fatalf("signing test token: %v", err)
	}
	return s
}

func TestGenerateAndValidateToken(t *testing.T) {
	auth := NewAuthServices(testSecret)

	token, err := auth.GenerateToken("user-123")
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	userID, err := auth.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken: %v", err)
	}
	if userID != "user-123" {
		t.Errorf("userID = %q, want %q", userID, "user-123")
	}
}

func TestValidateTokenRejects(t *testing.T) {
	auth := NewAuthServices(testSecret)
	now := time.Now()

	valid, err := auth.GenerateToken("user-123")
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	tests := []struct {
		name  string
		token string
	}{
		{
			name: "expired",
			token: signWith(t, jwt.SigningMethodHS256, []byte(testSecret), jwt.RegisteredClaims{
				Subject:   "user-123",
				IssuedAt:  jwt.NewNumericDate(now.Add(-2 * time.Hour)),
				ExpiresAt: jwt.NewNumericDate(now.Add(-time.Hour)),
			}),
		},
		{
			name: "wrong secret",
			token: signWith(t, jwt.SigningMethodHS256, []byte("other-secret"), jwt.RegisteredClaims{
				Subject:   "user-123",
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
			}),
		},
		{
			name: "no exp",
			token: signWith(t, jwt.SigningMethodHS256, []byte(testSecret), jwt.RegisteredClaims{
				Subject: "user-123",
			}),
		},
		{
			name: "other algorithm",
			token: signWith(t, jwt.SigningMethodHS512, []byte(testSecret), jwt.RegisteredClaims{
				Subject:   "user-123",
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
			}),
		},
		{
			name: "alg none",
			token: signWith(t, jwt.SigningMethodNone, jwt.UnsafeAllowNoneSignatureType, jwt.RegisteredClaims{
				Subject:   "user-123",
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
			}),
		},
		{
			name:  "tampered signature",
			token: valid[:len(valid)-2] + "xx",
		},
		{
			name:  "garbage",
			token: "not-a-jwt",
		},
		{
			name:  "empty",
			token: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if userID, err := auth.ValidateToken(tt.token); err == nil {
				t.Errorf("ValidateToken accepted token, got userID %q", userID)
			}
		})
	}
}
