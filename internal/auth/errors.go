package auth

import "errors"

var (
	ErrUserNotFound  = errors.New("User not found")
	ErrUnauthorized  = errors.New("Unauthorized")
	ErrTokenNotFound = errors.New("Token not found, please try to login again")
	ErrTokenRevoked  = errors.New("Token is revoked")
	ErrTokenExpired  = errors.New("Token is expired")
)

type TokenValidator interface {
	ValidateAccessToken(tokenID string) error
}
