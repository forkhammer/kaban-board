package domain

import "fmt"

type ValidationError struct {
	Message string
	Err     error
}

func (e *ValidationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

func NewValidationError(message string, err ...error) *ValidationError {
	var wrapped error
	if len(err) > 0 {
		wrapped = err[0]
	}
	return &ValidationError{
		Message: message,
		Err:     wrapped,
	}
}

type NotFoundError struct {
	Entity string
	ID     string
	Err    error
}

func (e *NotFoundError) Error() string {
	if e.ID != "" {
		if e.Err != nil {
			return fmt.Sprintf("%s with id '%v' not found: %v", e.Entity, e.ID, e.Err)
		}
		return fmt.Sprintf("%s with id '%v' not found", e.Entity, e.ID)
	}
	if e.Err != nil {
		return fmt.Sprintf("%s not found: %v", e.Entity, e.Err)
	}
	return fmt.Sprintf("%s not found", e.Entity)
}

func (e *NotFoundError) Unwrap() error {
	return e.Err
}

func NewNotFoundError(entity string, id string, err error) *NotFoundError {
	return &NotFoundError{
		Entity: entity,
		ID:     id,
		Err:    err,
	}
}

type ConflictError struct {
	Entity string
	Data   any
	Err    error
}

func (e *ConflictError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s conflict: %v", e.Entity, e.Err)
	}
	return fmt.Sprintf("%s conflict", e.Entity)
}

func (e *ConflictError) Unwrap() error {
	return e.Err
}

func NewConflictError(entity string, data any, err error) *ConflictError {
	return &ConflictError{
		Entity: entity,
		Data:   data,
		Err:    err,
	}
}
