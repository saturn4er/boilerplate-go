package testservice

type SomeEnum byte

const (
	SomeEnumValue1 SomeEnum = iota + 1
	SomeEnumValue2
)

type Role byte

const (
	RoleUser Role = iota + 1
	RoleAdmin
	RoleAnonymous
)

type OIDCProvider byte

const (
	OIDCProviderGoogle OIDCProvider = iota + 1
)
