package teststorage

import (
	fmt "fmt"

	testservice "github.com/saturn4er/boilerplate-go/test/test/testservice"
)

const (
	someEnumValue1 = "value1"
	someEnumValue2 = "value2"
)

func convertSomeEnumToDB(someEnumValue testservice.SomeEnum) (string, error) {
	result, ok := map[testservice.SomeEnum]string{
		testservice.SomeEnumValue1: someEnumValue1,
		testservice.SomeEnumValue2: someEnumValue2,
	}[someEnumValue]
	if !ok {
		return "", fmt.Errorf("unknown SomeEnum value: %d", someEnumValue)
	}
	return result, nil
}

func convertSomeEnumFromDB(someEnumValue string) (testservice.SomeEnum, error) {
	result, ok := map[string]testservice.SomeEnum{
		someEnumValue1: testservice.SomeEnumValue1,
		someEnumValue2: testservice.SomeEnumValue2,
	}[someEnumValue]
	if !ok {
		return 0, fmt.Errorf("unknown SomeEnum db value: %s", someEnumValue)
	}
	return result, nil
}

const (
	roleUser      = "user"
	roleAdmin     = "admin"
	roleAnonymous = "anonymous"
)

func convertRoleToDB(roleValue testservice.Role) (string, error) {
	result, ok := map[testservice.Role]string{
		testservice.RoleUser:      roleUser,
		testservice.RoleAdmin:     roleAdmin,
		testservice.RoleAnonymous: roleAnonymous,
	}[roleValue]
	if !ok {
		return "", fmt.Errorf("unknown Role value: %d", roleValue)
	}
	return result, nil
}

func convertRoleFromDB(roleValue string) (testservice.Role, error) {
	result, ok := map[string]testservice.Role{
		roleUser:      testservice.RoleUser,
		roleAdmin:     testservice.RoleAdmin,
		roleAnonymous: testservice.RoleAnonymous,
	}[roleValue]
	if !ok {
		return 0, fmt.Errorf("unknown Role db value: %s", roleValue)
	}
	return result, nil
}

const (
	oidcProviderGoogle = "google"
)

func convertOIDCProviderToDB(oidcProviderValue testservice.OIDCProvider) (string, error) {
	result, ok := map[testservice.OIDCProvider]string{
		testservice.OIDCProviderGoogle: oidcProviderGoogle,
	}[oidcProviderValue]
	if !ok {
		return "", fmt.Errorf("unknown OIDCProvider value: %d", oidcProviderValue)
	}
	return result, nil
}

func convertOIDCProviderFromDB(oidcProviderValue string) (testservice.OIDCProvider, error) {
	result, ok := map[string]testservice.OIDCProvider{
		oidcProviderGoogle: testservice.OIDCProviderGoogle,
	}[oidcProviderValue]
	if !ok {
		return 0, fmt.Errorf("unknown OIDCProvider db value: %s", oidcProviderValue)
	}
	return result, nil
}
