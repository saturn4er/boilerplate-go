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
