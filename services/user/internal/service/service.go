package service

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailExists        = errors.New("email is already taken")
	ErrUserNotFound       = errors.New("user not found")
	ErrProfileNotFound    = errors.New("profile not found")
)
