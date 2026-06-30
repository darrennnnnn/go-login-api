package user

import "errors"

var (
	ErrEmailInUse     = errors.New("Email already in use")
	ErrUsernameInUse  = errors.New("Username already in use")
	ErrUserNotFound   = errors.New("User not found")
)
