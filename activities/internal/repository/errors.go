package repository

import (
	"activities/internal/errors"
)

// Re-export errors from the central errors package for backward compatibility
var (
	ErrActivityNotFound      = errors.ErrActivityNotFound
	ErrActivityFull          = errors.ErrActivityFull
	ErrUserAlreadyInscribed  = errors.ErrUserAlreadyInscribed
	ErrUserNotInscribed      = errors.ErrUserNotInscribed
	ErrNoFieldsToUpdate      = errors.ErrNoFieldsToUpdate
	ErrInvalidUserID         = errors.ErrInvalidUserID
	ErrInvalidIDFormat       = errors.ErrInvalidIDFormat
	ErrActivityAlreadyExists = errors.ErrActivityAlreadyExists
)
