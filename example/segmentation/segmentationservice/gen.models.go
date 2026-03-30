package segmentationservice

import (
	uuid "github.com/google/uuid"

	filter "github.com/saturn4er/boilerplate-go/lib/filter"
	order "github.com/saturn4er/boilerplate-go/lib/order"
	// user code 'imports'
	// end user code 'imports'
)

type UserTagField byte

const (
	UserTagFieldID UserTagField = iota + 1
	UserTagFieldUserID
	UserTagFieldKey
	UserTagFieldValue
)

type UserTagFilter struct {
	ID     filter.Filter[uuid.UUID]
	UserID filter.Filter[uuid.UUID]
	Key    filter.Filter[string]
	Or     []*UserTagFilter
	And    []*UserTagFilter
}
type UserTagOrder order.Order[UserTagField]

type UserTag struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Key    string
	Value  string
}

// user code 'UserTag methods'
// end user code 'UserTag methods'

func (u *UserTag) Copy() UserTag {
	var result UserTag
	result.ID = u.ID
	result.UserID = u.UserID
	result.Key = u.Key
	result.Value = u.Value

	return result
}
func (u *UserTag) Equals(to *UserTag) bool {
	if (u == nil) != (to == nil) {
		return false
	}
	if u == nil && to == nil {
		return true
	}
	if u.ID != to.ID {
		return false
	}
	if u.UserID != to.UserID {
		return false
	}
	if u.Key != to.Key {
		return false
	}
	if u.Value != to.Value {
		return false
	}

	return true
}

type SetUserTagCommandField byte

const (
	SetUserTagCommandFieldID SetUserTagCommandField = iota + 1
	SetUserTagCommandFieldData
	SetUserTagCommandFieldIdempotencyKey
)

type SetUserTagCommandFilter struct {
	Or  []*SetUserTagCommandFilter
	And []*SetUserTagCommandFilter
}
type SetUserTagCommandOrder order.Order[SetUserTagCommandField]

type SetUserTagCommand struct {
	ID             uuid.UUID
	Data           SetUserTagCommandData
	IdempotencyKey string
}

// user code 'SetUserTagCommand methods'
// end user code 'SetUserTagCommand methods'

func (s *SetUserTagCommand) Copy() SetUserTagCommand {
	var result SetUserTagCommand
	result.ID = s.ID
	result.Data = s.Data.Copy() // model
	result.IdempotencyKey = s.IdempotencyKey

	return result
}
func (s *SetUserTagCommand) Equals(to *SetUserTagCommand) bool {
	if (s == nil) != (to == nil) {
		return false
	}
	if s == nil && to == nil {
		return true
	}
	if s.ID != to.ID {
		return false
	}
	if !s.Data.Equals(&to.Data) {
		return false
	}
	if s.IdempotencyKey != to.IdempotencyKey {
		return false
	}

	return true
}

type SetUserTagCommandDataField byte

const (
	SetUserTagCommandDataFieldUserID SetUserTagCommandDataField = iota + 1
	SetUserTagCommandDataFieldKey
	SetUserTagCommandDataFieldValue
)

type SetUserTagCommandDataFilter struct {
	Or  []*SetUserTagCommandDataFilter
	And []*SetUserTagCommandDataFilter
}
type SetUserTagCommandDataOrder order.Order[SetUserTagCommandDataField]

type SetUserTagCommandData struct {
	UserID uuid.UUID
	Key    string
	Value  string
}

// user code 'SetUserTagCommandData methods'
// end user code 'SetUserTagCommandData methods'

func (s *SetUserTagCommandData) Copy() SetUserTagCommandData {
	var result SetUserTagCommandData
	result.UserID = s.UserID
	result.Key = s.Key
	result.Value = s.Value

	return result
}
func (s *SetUserTagCommandData) Equals(to *SetUserTagCommandData) bool {
	if (s == nil) != (to == nil) {
		return false
	}
	if s == nil && to == nil {
		return true
	}
	if s.UserID != to.UserID {
		return false
	}
	if s.Key != to.Key {
		return false
	}
	if s.Value != to.Value {
		return false
	}

	return true
}
