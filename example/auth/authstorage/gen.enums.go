package authstorage

import (
	fmt "fmt"

	authservice "github.com/saturn4er/boilerplate-go/example/auth/authservice"
	// user code 'imports'
	// end user code 'imports'
)

const (
	roleAdmin = "admin"
	roleUser  = "user"
)

func convertRoleToDB(roleValue authservice.Role) (string, error) {
	result, ok := map[authservice.Role]string{
		authservice.RoleAdmin: roleAdmin,
		authservice.RoleUser:  roleUser,
	}[roleValue]
	if !ok {
		return "", fmt.Errorf("unknown Role value: %d", roleValue)
	}
	return result, nil
}

func convertRoleFromDB(roleValue string) (authservice.Role, error) {
	result, ok := map[string]authservice.Role{
		roleAdmin: authservice.RoleAdmin,
		roleUser:  authservice.RoleUser,
	}[roleValue]
	if !ok {
		return 0, fmt.Errorf("unknown Role db value: %s", roleValue)
	}
	return result, nil
}
