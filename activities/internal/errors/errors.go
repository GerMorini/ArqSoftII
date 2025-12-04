package errors

import "errors"

// Repository errors
var (
	ErrActivityNotFound      = errors.New("activity not found")
	ErrActivityFull          = errors.New("activity is full")
	ErrUserAlreadyInscribed  = errors.New("user already inscribed")
	ErrUserNotInscribed      = errors.New("user not inscribed in activity")
	ErrNoFieldsToUpdate      = errors.New("no fields to update")
	ErrInvalidUserID         = errors.New("invalid user id format")
	ErrInvalidIDFormat       = errors.New("invalid ID format")
	ErrActivityAlreadyExists = errors.New("activity with the same ID already exists")
)

// Service validation errors
var (
	ErrValidation                = errors.New("validation error")
	ErrTitleRequired             = errors.New("titulo is required and cannot be empty")
	ErrInstructorRequired        = errors.New("instructor is required and cannot be empty")
	ErrTimeRequired              = errors.New("hora_inicio and hora_fin are required and cannot be empty")
	ErrCapacityRequired          = errors.New("cupo is required and cannot be empty")
	ErrCapacityNegative          = errors.New("cupo cannot be negative")
	ErrDayRequired               = errors.New("dia is required and cannot be empty")
	ErrInvalidDay                = errors.New("dia must be a valid day of the week (e.g., Lunes, Martes, etc.)")
	ErrActivityDoesNotExist      = errors.New("activity does not exist")
	ErrCapacityLessThanInscribed = errors.New("cupo cannot be less than the number of inscribed users")
	ErrInscritosExceedCapacity   = errors.New("number of inscritos cannot exceed capacity")
)

// Service operation errors
var (
	ErrPublishEventFailed            = errors.New("failed to publish event")
	ErrRollbackFailed                = errors.New("rollback failed")
	ErrCreatingActivityInRepository  = errors.New("error creating activity in repository")
	ErrGettingActivityFromRepository = errors.New("error getting activity from repository")
)
