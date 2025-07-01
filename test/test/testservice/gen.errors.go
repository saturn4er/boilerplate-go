package testservice

import (
	fmt "fmt"
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

const (
	ErrSomeModelNotFound      = NotFoundError("SomeModel")
	ErrSomeModelAlreadyExists = AlreadyExistsError("SomeModel")
)
const (
	ErrSomeOtherModelNotFound      = NotFoundError("SomeOtherModel")
	ErrSomeOtherModelAlreadyExists = AlreadyExistsError("SomeOtherModel")
)
const (
	ErrOneOfValue1NotFound      = NotFoundError("OneOfValue1")
	ErrOneOfValue1AlreadyExists = AlreadyExistsError("OneOfValue1")
)
const (
	ErrOneOfValue2NotFound      = NotFoundError("OneOfValue2")
	ErrOneOfValue2AlreadyExists = AlreadyExistsError("OneOfValue2")
)
const (
	ErrPasswordRecoveryEventNotFound      = NotFoundError("PasswordRecoveryEvent")
	ErrPasswordRecoveryEventAlreadyExists = AlreadyExistsError("PasswordRecoveryEvent")
)
const (
	ErrPasswordRecoveryRequestedEventDataNotFound      = NotFoundError("PasswordRecoveryRequestedEventData")
	ErrPasswordRecoveryRequestedEventDataAlreadyExists = AlreadyExistsError("PasswordRecoveryRequestedEventData")
)
const (
	ErrPasswordRecoveryCompletedEventDataNotFound      = NotFoundError("PasswordRecoveryCompletedEventData")
	ErrPasswordRecoveryCompletedEventDataAlreadyExists = AlreadyExistsError("PasswordRecoveryCompletedEventData")
)
