package authservice

import (
	fmt "fmt"
	strings "strings"
	// user code 'imports'
	// end user code 'imports'
)

type NotFoundError string

func (n NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", string(n))
}

type AlreadyExistsError string

func (a AlreadyExistsError) Error() string {
	return fmt.Sprintf("%s already exists", string(a))
}

type ConflictError struct {
	Entity string
	Fields []string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("%s conflict on %s", e.Entity, strings.Join(e.Fields, ", "))
}

func (e *ConflictError) Is(target error) bool {
	if ae, ok := target.(AlreadyExistsError); ok {
		return string(ae) == e.Entity
	}
	return false
}

const (
	ErrUserNotFound      = NotFoundError("User")
	ErrUserAlreadyExists = AlreadyExistsError("User")
)

var (
	ErrUserEmailAlreadyExists = &ConflictError{Entity: "User", Fields: []string{"Email"}}
)
