package authstorage

import (
	uuid "github.com/google/uuid"

	authservice "github.com/saturn4er/boilerplate-go/example/auth/authservice"
	// user code 'imports'
	// end user code 'imports'
)

type dbUser struct {
	ID    uuid.UUID `gorm:"column:id;primaryKey"`
	Email string    `gorm:"column:email;type:text;"`
	Name  string    `gorm:"column:name;type:text;"`
	Role  string    `gorm:"column:role;type:text;"`
}

func convertUserToDB(src *authservice.User) (*dbUser, error) {
	result := &dbUser{}
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

func convertUserFromDB(src *dbUser) (*authservice.User, error) {
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
func (a dbUser) TableName() string {
	return "users"
}
