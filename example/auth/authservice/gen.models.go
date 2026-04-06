package authservice

import (
	uuid "github.com/google/uuid"

	filter "github.com/saturn4er/boilerplate-go/lib/filter"
	order "github.com/saturn4er/boilerplate-go/lib/order"
	// user code 'imports'
	// end user code 'imports'
)

type UserField byte

const (
	UserFieldID UserField = iota + 1
	UserFieldEmail
	UserFieldName
	UserFieldRole
)

type UserFilter struct {
	ID    filter.Filter[uuid.UUID]
	Email filter.Filter[string]
	Or    []*UserFilter
	And   []*UserFilter
}
type UserOrder order.Order[UserField]

type User struct {
	ID    uuid.UUID
	Email string
	Name  string
	Role  Role
}

// user code 'User methods'
// end user code 'User methods'

func (u *User) Copy() User {
	var result User
	result.ID = u.ID
	result.Email = u.Email
	result.Name = u.Name
	result.Role = u.Role // enum

	return result
}
func (u *User) Equals(to *User) bool {
	if (u == nil) != (to == nil) {
		return false
	}
	if u == nil && to == nil {
		return true
	}
	if u.ID != to.ID {
		return false
	}
	if u.Email != to.Email {
		return false
	}
	if u.Name != to.Name {
		return false
	}
	if u.Role != to.Role {
		return false
	}

	return true
}
