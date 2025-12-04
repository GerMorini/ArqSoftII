package services

import (
	"activities/internal/errors"
)

// Re-export errors from the central errors package for backward compatibility
var (
	ErrValidation                    = errors.ErrValidation
	ErrTitleRequired                 = errors.ErrTitleRequired
	ErrInstructorRequired            = errors.ErrInstructorRequired
	ErrTimeRequired                  = errors.ErrTimeRequired
	ErrCapacityRequired              = errors.ErrCapacityRequired
	ErrCapacityNegative              = errors.ErrCapacityNegative
	ErrDayRequired                   = errors.ErrDayRequired
	ErrInvalidDay                    = errors.ErrInvalidDay
	ErrActivityDoesNotExist          = errors.ErrActivityDoesNotExist
	ErrCapacityLessThanInscribed     = errors.ErrCapacityLessThanInscribed
	ErrInscritosExceedCapacity       = errors.ErrInscritosExceedCapacity
	ErrPublishEventFailed            = errors.ErrPublishEventFailed
	ErrRollbackFailed                = errors.ErrRollbackFailed
	ErrCreatingActivityInRepository  = errors.ErrCreatingActivityInRepository
	ErrGettingActivityFromRepository = errors.ErrGettingActivityFromRepository
)
