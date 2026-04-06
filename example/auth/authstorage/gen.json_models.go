package authstorage

import (
	driver "database/sql/driver"
	json "encoding/json"

	uuid "github.com/google/uuid"

	authservice "github.com/saturn4er/boilerplate-go/example/auth/authservice"
	// user code 'imports'
	// end user code 'imports'
)

type jsonUser struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
	Role  string    `json:"role"`
}

func (u *jsonUser) Scan(value any) error {
	return json.Unmarshal(value.([]byte), u)
}

func (u jsonUser) Value() (driver.Value, error) {
	return json.Marshal(u)
}

func convertUserToJsonModel(src *authservice.User) (*jsonUser, error) {
	result := &jsonUser{}
	result.ID = src.ID
	result.Email = src.Email
	result.Name = src.Name
	tmp3, err := convertRoleToDB(src.Role)
	if err != nil {
		return nil, err
	}
	result.Role = tmp3
	return result, nil
}

func convertUserFromJsonModel(src *jsonUser) (*authservice.User, error) {
	result := &authservice.User{}
	result.ID = src.ID
	result.Email = src.Email
	result.Name = src.Name
	tmp7, err := convertRoleFromDB(src.Role)
	if err != nil {
		return nil, err
	}
	result.Role = tmp7
	return result, nil
}
