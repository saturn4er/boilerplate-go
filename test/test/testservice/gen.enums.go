package testservice

type SomeEnum byte

const (
	SomeEnumValue1 SomeEnum = iota + 1
	SomeEnumValue2
)

// user code 'SomeEnum methods'
// end user code 'SomeEnum methods'
type Role byte

const (
	RoleUser Role = iota + 1
	RoleAdmin
	RoleAnonymous
)

// user code 'Role methods'
// end user code 'Role methods'
type OIDCProvider byte

const (
	OIDCProviderGoogle OIDCProvider = iota + 1
)

// user code 'OIDCProvider methods'
// end user code 'OIDCProvider methods'
