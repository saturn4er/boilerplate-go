package authservice

import (
	strconv "strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	// user code 'imports'
	// end user code 'imports'
)

type Role byte

const (
	RoleAdmin Role = iota + 1
	RoleUser
)

// user code 'Role methods'
// end user code 'Role methods'
var AllRoles = []Role{RoleAdmin, RoleUser}

func (r Role) IsValid() bool {
	return r > 0 && r < 3
}

var SuperUserRoles = []Role{
	RoleAdmin,
}

func (r Role) IsSuperUser() bool {
	if r < 1 || r > 3 {
		return false
	}
	return []bool{false, true, false}[r]
}
func (r Role) IsAdmin() bool {
	return r == RoleAdmin
}
func (r Role) IsUser() bool {
	return r == RoleUser
}
func (r Role) Validate() error {
	return validation.Validate(r, validation.In(RoleAdmin, RoleUser))
}
func AllRolesFn() []Role {
	return []Role{RoleAdmin, RoleUser}
}
func (r Role) String() string {
	const names = "AdminUser"

	var indexes = [...]int32{0, 5, 9}
	if r < 1 || r > 2 {
		return "Role(" + strconv.FormatInt(int64(r), 10) + ")"
	}

	return names[indexes[r-1]:indexes[r]]
}
