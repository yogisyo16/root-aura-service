// Package services: authServices.go handles JWT-based authentication —
// issuing tokens on login and validating them on protected requests.
// See docs/BACKLOG.md for the rest of the plan (login handler, JWT
// middleware, user-scoped authorization) that builds on top of this.
package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// tokenTTL is how long an issued token stays valid. Kept short since there
// are no refresh tokens yet — the user just logs in again once it expires.
const tokenTTL = 60 * time.Minute

type AuthServices interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(tokenString string) (userID string, err error)
}

type authService struct {
	secret []byte
}

// NewAuthServices takes the HS256 signing secret as an argument rather than
// reading JWT_SECRET itself, so the env lookup (and fail-fast on a missing
// value) lives in api/main.go, and tests can pass in a throwaway secret.
func NewAuthServices(secret string) AuthServices {
	return &authService{secret: []byte(secret)}
}

// GenerateToken issues a signed HS256 token for userID. Claims are kept to
// the registered minimum: sub (user ID), iat, exp.
func (s *authService) GenerateToken(userID string) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(tokenTTL)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

// ValidateToken verifies the signature and expiry and returns the user ID
// stored in sub.
//
// WithValidMethods pins the algorithm to HS256 so a token can't pick its own
// (the "alg confusion" attack). WithExpirationRequired rejects tokens with
// no exp at all — by default v5 only checks exp when it's present.
func (s *authService) ValidateToken(tokenString string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (any, error) { return s.secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return "", err
	}
	if !token.Valid || claims.Subject == "" {
		return "", errors.New("invalid token")
	}

	return claims.Subject, nil
}
