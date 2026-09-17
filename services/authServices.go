// Package services: authServices.go is a placeholder for JWT-based
// authentication (token generation/validation), deliberately deferred.
// See docs/BACKLOG.md for the full implementation plan:
//   - GenerateToken(userID string) (string, error)
//   - ValidateToken(tokenString string) (userID string, err error)
// plus the JWT middleware, login handler, and user-scoped authorization
// that depend on it.
package services
