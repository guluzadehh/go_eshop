package storage

import "errors"

var (
	UserNotFound    = errors.New("user not found")
	UserExists      = errors.New("user already exists")
	ProfileNotFound = errors.New("profile not found")
)
