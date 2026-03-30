package authservice

type Role byte

const (
	RoleAdmin Role = iota + 1
	RoleUser
)

// user code 'Role methods'
// end user code 'Role methods'
