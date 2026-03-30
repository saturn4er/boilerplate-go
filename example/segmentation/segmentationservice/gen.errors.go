package segmentationservice

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
	ErrUserTagNotFound      = NotFoundError("UserTag")
	ErrUserTagAlreadyExists = AlreadyExistsError("UserTag")
)
const (
	ErrSetUserTagCommandNotFound      = NotFoundError("SetUserTagCommand")
	ErrSetUserTagCommandAlreadyExists = AlreadyExistsError("SetUserTagCommand")
)
const (
	ErrSetUserTagCommandDataNotFound      = NotFoundError("SetUserTagCommandData")
	ErrSetUserTagCommandDataAlreadyExists = AlreadyExistsError("SetUserTagCommandData")
)
