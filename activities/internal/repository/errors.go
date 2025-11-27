package repository

import "errors"

var (
	ErrActivityNotFound     = errors.New("activity not found")
	ErrActivityFull         = errors.New("activity is full")
	ErrUserAlreadyInscribed = errors.New("user already inscribed")
	ErrUserNotInscribed     = errors.New("user not inscribed in activity")
)
