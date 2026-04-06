package segmentationservice

import (
	uuid "github.com/google/uuid"

	filter "github.com/saturn4er/boilerplate-go/lib/filter"
	order "github.com/saturn4er/boilerplate-go/lib/order"
	// user code 'imports'
	// end user code 'imports'
)

type SomeOneOf interface {
	isSomeOneOf()
	SomeOneOfEquals(SomeOneOf) bool
	// user code 'SomeOneOf methods'
	// end user code 'SomeOneOf methods'
}

func (*OneOfValue1) isSomeOneOf() {}
func (o *OneOfValue1) SomeOneOfEquals(to SomeOneOf) bool {
	if (o == nil) != (to == nil) {
		return false
	}
	if o == nil && to == nil {
		return true
	}

	toTyped, ok := to.(*OneOfValue1)
	if !ok {
		return false
	}

	return o.Equals(toTyped)
}
func (*OneOfValue2) isSomeOneOf() {}
func (o *OneOfValue2) SomeOneOfEquals(to SomeOneOf) bool {
	if (o == nil) != (to == nil) {
		return false
	}
	if o == nil && to == nil {
		return true
	}

	toTyped, ok := to.(*OneOfValue2)
	if !ok {
		return false
	}

	return o.Equals(toTyped)
}

func copySomeOneOf(val SomeOneOf) SomeOneOf {
	if val == nil {
		return nil
	}

	switch val := val.(type) {
	case *OneOfValue1:
		valCopy := val.Copy()
		return &valCopy
	case *OneOfValue2:
		valCopy := val.Copy()
		return &valCopy
	}
	panic("called copySomeOneOf with invalid type")
}

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

type SomeModelField byte

const (
	SomeModelFieldID SomeModelField = iota + 1
	SomeModelFieldName
	SomeModelFieldDescription
	SomeModelFieldModelField
	SomeModelFieldModelPtrField
	SomeModelFieldOneOfField
	SomeModelFieldOneOfPtrField
	SomeModelFieldEnumField
	SomeModelFieldEnumPtrField
	SomeModelFieldAnyField
	SomeModelFieldAnyPtrField
	SomeModelFieldMapModelField
	SomeModelFieldMapModelPtrField
	SomeModelFieldMapOneOfField
	SomeModelFieldMapOneOfPtrField
	SomeModelFieldMapEnumField
	SomeModelFieldMapEnumPtrField
	SomeModelFieldMapAnyField
	SomeModelFieldMapAnyPtrField
	SomeModelFieldModelSliceField
	SomeModelFieldModelPtrSliceField
	SomeModelFieldOneOfSliceField
	SomeModelFieldOneOfPtrSliceField
	SomeModelFieldSliceEnumField
	SomeModelFieldSliceEnumPtrField
	SomeModelFieldSliceAnyField
	SomeModelFieldSliceAnyPtrField
)

type SomeModelFilter struct {
	ID          filter.Filter[uuid.UUID]
	Name        filter.Filter[string]
	Description filter.Filter[*string]
	Or          []*SomeModelFilter
	And         []*SomeModelFilter
}
type SomeModelOrder order.Order[SomeModelField]

type SomeModel struct {
	ID                 uuid.UUID
	Name               string
	Description        *string
	ModelField         SomeOtherModel
	ModelPtrField      *SomeOtherModel
	OneOfField         SomeOneOf
	OneOfPtrField      *SomeOneOf
	EnumField          SomeEnum
	EnumPtrField       *SomeEnum
	AnyField           any
	AnyPtrField        *any
	MapModelField      map[string]SomeOtherModel
	MapModelPtrField   map[string]*SomeOtherModel
	MapOneOfField      map[string]SomeOneOf
	MapOneOfPtrField   map[string]*SomeOneOf
	MapEnumField       map[string]SomeEnum
	MapEnumPtrField    map[string]*SomeEnum
	MapAnyField        map[string]any
	MapAnyPtrField     map[string]*any
	ModelSliceField    []SomeOtherModel
	ModelPtrSliceField []*SomeOtherModel
	OneOfSliceField    []SomeOneOf
	OneOfPtrSliceField []*SomeOneOf
	SliceEnumField     []SomeEnum
	SliceEnumPtrField  []*SomeEnum
	SliceAnyField      []any
	SliceAnyPtrField   []*any
}

// user code 'SomeModel methods'
// end user code 'SomeModel methods'

func (s *SomeModel) Copy() SomeModel {
	var result SomeModel
	result.ID = s.ID
	result.Name = s.Name
	if s.Description != nil {
		var tmp string
		tmp = (*s.Description)
		result.Description = &tmp
	}
	result.ModelField = s.ModelField.Copy() // model
	if s.ModelPtrField != nil {
		var tmp1 SomeOtherModel
		tmp1 = (*s.ModelPtrField).Copy() // model
		result.ModelPtrField = &tmp1
	}
	result.OneOfField = copySomeOneOf(s.OneOfField)
	if s.OneOfPtrField != nil {
		var tmp2 SomeOneOf
		tmp2 = copySomeOneOf((*s.OneOfPtrField))
		result.OneOfPtrField = &tmp2
	}
	result.EnumField = s.EnumField // enum
	if s.EnumPtrField != nil {
		var tmp3 SomeEnum
		tmp3 = (*s.EnumPtrField) // enum
		result.EnumPtrField = &tmp3
	}
	result.AnyField = s.AnyField
	if s.AnyPtrField != nil {
		var tmp4 any
		tmp4 = (*s.AnyPtrField)
		result.AnyPtrField = &tmp4
	}
	tmp5 := make(map[string]SomeOtherModel)
	for k, v := range s.MapModelField {
		var keyCopy string
		var valueCopy SomeOtherModel
		keyCopy = k
		valueCopy = v.Copy() // model
		tmp5[keyCopy] = valueCopy
	}
	result.MapModelField = tmp5
	tmp6 := make(map[string]*SomeOtherModel)
	for k, v := range s.MapModelPtrField {
		var keyCopy1 string
		var valueCopy1 *SomeOtherModel
		keyCopy1 = k
		if v != nil {
			var tmp7 SomeOtherModel
			tmp7 = (*v).Copy() // model
			valueCopy1 = &tmp7
		}
		tmp6[keyCopy1] = valueCopy1
	}
	result.MapModelPtrField = tmp6
	tmp8 := make(map[string]SomeOneOf)
	for k, v := range s.MapOneOfField {
		var keyCopy2 string
		var valueCopy2 SomeOneOf
		keyCopy2 = k
		valueCopy2 = copySomeOneOf(v)
		tmp8[keyCopy2] = valueCopy2
	}
	result.MapOneOfField = tmp8
	tmp9 := make(map[string]*SomeOneOf)
	for k, v := range s.MapOneOfPtrField {
		var keyCopy3 string
		var valueCopy3 *SomeOneOf
		keyCopy3 = k
		if v != nil {
			var tmp10 SomeOneOf
			tmp10 = copySomeOneOf((*v))
			valueCopy3 = &tmp10
		}
		tmp9[keyCopy3] = valueCopy3
	}
	result.MapOneOfPtrField = tmp9
	tmp11 := make(map[string]SomeEnum)
	for k, v := range s.MapEnumField {
		var keyCopy4 string
		var valueCopy4 SomeEnum
		keyCopy4 = k
		valueCopy4 = v // enum
		tmp11[keyCopy4] = valueCopy4
	}
	result.MapEnumField = tmp11
	tmp12 := make(map[string]*SomeEnum)
	for k, v := range s.MapEnumPtrField {
		var keyCopy5 string
		var valueCopy5 *SomeEnum
		keyCopy5 = k
		if v != nil {
			var tmp13 SomeEnum
			tmp13 = (*v) // enum
			valueCopy5 = &tmp13
		}
		tmp12[keyCopy5] = valueCopy5
	}
	result.MapEnumPtrField = tmp12
	tmp14 := make(map[string]any)
	for k, v := range s.MapAnyField {
		var keyCopy6 string
		var valueCopy6 any
		keyCopy6 = k
		valueCopy6 = v
		tmp14[keyCopy6] = valueCopy6
	}
	result.MapAnyField = tmp14
	tmp15 := make(map[string]*any)
	for k, v := range s.MapAnyPtrField {
		var keyCopy7 string
		var valueCopy7 *any
		keyCopy7 = k
		if v != nil {
			var tmp16 any
			tmp16 = (*v)
			valueCopy7 = &tmp16
		}
		tmp15[keyCopy7] = valueCopy7
	}
	result.MapAnyPtrField = tmp15
	tmp17 := make([]SomeOtherModel, 0, len(s.ModelSliceField))
	for _, i := range s.ModelSliceField {
		var itemCopy SomeOtherModel
		itemCopy = i.Copy() // model
		tmp17 = append(tmp17, itemCopy)
	}
	result.ModelSliceField = tmp17
	tmp18 := make([]*SomeOtherModel, 0, len(s.ModelPtrSliceField))
	for _, i1 := range s.ModelPtrSliceField {
		var itemCopy1 *SomeOtherModel
		if i1 != nil {
			var tmp19 SomeOtherModel
			tmp19 = (*i1).Copy() // model
			itemCopy1 = &tmp19
		}
		tmp18 = append(tmp18, itemCopy1)
	}
	result.ModelPtrSliceField = tmp18
	tmp20 := make([]SomeOneOf, 0, len(s.OneOfSliceField))
	for _, i2 := range s.OneOfSliceField {
		var itemCopy2 SomeOneOf
		itemCopy2 = copySomeOneOf(i2)
		tmp20 = append(tmp20, itemCopy2)
	}
	result.OneOfSliceField = tmp20
	tmp21 := make([]*SomeOneOf, 0, len(s.OneOfPtrSliceField))
	for _, i3 := range s.OneOfPtrSliceField {
		var itemCopy3 *SomeOneOf
		if i3 != nil {
			var tmp22 SomeOneOf
			tmp22 = copySomeOneOf((*i3))
			itemCopy3 = &tmp22
		}
		tmp21 = append(tmp21, itemCopy3)
	}
	result.OneOfPtrSliceField = tmp21
	tmp23 := make([]SomeEnum, 0, len(s.SliceEnumField))
	for _, i4 := range s.SliceEnumField {
		var itemCopy4 SomeEnum
		itemCopy4 = i4 // enum
		tmp23 = append(tmp23, itemCopy4)
	}
	result.SliceEnumField = tmp23
	tmp24 := make([]*SomeEnum, 0, len(s.SliceEnumPtrField))
	for _, i5 := range s.SliceEnumPtrField {
		var itemCopy5 *SomeEnum
		if i5 != nil {
			var tmp25 SomeEnum
			tmp25 = (*i5) // enum
			itemCopy5 = &tmp25
		}
		tmp24 = append(tmp24, itemCopy5)
	}
	result.SliceEnumPtrField = tmp24
	tmp26 := make([]any, 0, len(s.SliceAnyField))
	for _, i6 := range s.SliceAnyField {
		var itemCopy6 any
		itemCopy6 = i6
		tmp26 = append(tmp26, itemCopy6)
	}
	result.SliceAnyField = tmp26
	tmp27 := make([]*any, 0, len(s.SliceAnyPtrField))
	for _, i7 := range s.SliceAnyPtrField {
		var itemCopy7 *any
		if i7 != nil {
			var tmp28 any
			tmp28 = (*i7)
			itemCopy7 = &tmp28
		}
		tmp27 = append(tmp27, itemCopy7)
	}
	result.SliceAnyPtrField = tmp27

	return result
}
func (s *SomeModel) Equals(to *SomeModel) bool {
	if (s == nil) != (to == nil) {
		return false
	}
	if s == nil && to == nil {
		return true
	}
	if s.ID != to.ID {
		return false
	}
	if s.Name != to.Name {
		return false
	}
	if (s.Description == nil) != (to.Description == nil) {
		return false
	}
	if s.Description != nil && to.Description != nil {
		if (*s.Description) != (*to.Description) {
			return false
		}
	}
	if !s.ModelField.Equals(&to.ModelField) {
		return false
	}
	if (s.ModelPtrField == nil) != (to.ModelPtrField == nil) {
		return false
	}
	if s.ModelPtrField != nil && to.ModelPtrField != nil {
		if !(*s.ModelPtrField).Equals(&(*to.ModelPtrField)) {
			return false
		}
	}
	if !s.OneOfField.SomeOneOfEquals(to.OneOfField) {
		return false
	}
	if (s.OneOfPtrField == nil) != (to.OneOfPtrField == nil) {
		return false
	}
	if s.OneOfPtrField != nil && to.OneOfPtrField != nil {
		if !(*s.OneOfPtrField).SomeOneOfEquals((*to.OneOfPtrField)) {
			return false
		}
	}
	if s.EnumField != to.EnumField {
		return false
	}
	if (s.EnumPtrField == nil) != (to.EnumPtrField == nil) {
		return false
	}
	if s.EnumPtrField != nil && to.EnumPtrField != nil {
		if (*s.EnumPtrField) != (*to.EnumPtrField) {
			return false
		}
	}
	if s.AnyField != to.AnyField {
		return false
	}
	if (s.AnyPtrField == nil) != (to.AnyPtrField == nil) {
		return false
	}
	if s.AnyPtrField != nil && to.AnyPtrField != nil {
		if (*s.AnyPtrField) != (*to.AnyPtrField) {
			return false
		}
	}
	// map comparision
	if len(s.MapModelField) != len(to.MapModelField) {
		return false
	}
	for k := range s.MapModelField {
		valB, ok := to.MapModelField[k]
		if !ok {
			return false
		}
		valA := s.MapModelField[k]
		if !valA.Equals(&valB) {
			return false
		}
	}
	// map comparision
	if len(s.MapModelPtrField) != len(to.MapModelPtrField) {
		return false
	}
	for k1 := range s.MapModelPtrField {
		valB1, ok := to.MapModelPtrField[k1]
		if !ok {
			return false
		}
		valA1 := s.MapModelPtrField[k1]
		if (valA1 == nil) != (valB1 == nil) {
			return false
		}
		if valA1 != nil && valB1 != nil {
			if !(*valA1).Equals(&(*valB1)) {
				return false
			}
		}
	}
	// map comparision
	if len(s.MapOneOfField) != len(to.MapOneOfField) {
		return false
	}
	for k2 := range s.MapOneOfField {
		valB2, ok := to.MapOneOfField[k2]
		if !ok {
			return false
		}
		valA2 := s.MapOneOfField[k2]
		if !valA2.SomeOneOfEquals(valB2) {
			return false
		}
	}
	// map comparision
	if len(s.MapOneOfPtrField) != len(to.MapOneOfPtrField) {
		return false
	}
	for k3 := range s.MapOneOfPtrField {
		valB3, ok := to.MapOneOfPtrField[k3]
		if !ok {
			return false
		}
		valA3 := s.MapOneOfPtrField[k3]
		if (valA3 == nil) != (valB3 == nil) {
			return false
		}
		if valA3 != nil && valB3 != nil {
			if !(*valA3).SomeOneOfEquals((*valB3)) {
				return false
			}
		}
	}
	// map comparision
	if len(s.MapEnumField) != len(to.MapEnumField) {
		return false
	}
	for k4 := range s.MapEnumField {
		valB4, ok := to.MapEnumField[k4]
		if !ok {
			return false
		}
		valA4 := s.MapEnumField[k4]
		if valA4 != valB4 {
			return false
		}
	}
	// map comparision
	if len(s.MapEnumPtrField) != len(to.MapEnumPtrField) {
		return false
	}
	for k5 := range s.MapEnumPtrField {
		valB5, ok := to.MapEnumPtrField[k5]
		if !ok {
			return false
		}
		valA5 := s.MapEnumPtrField[k5]
		if (valA5 == nil) != (valB5 == nil) {
			return false
		}
		if valA5 != nil && valB5 != nil {
			if (*valA5) != (*valB5) {
				return false
			}
		}
	}
	// map comparision
	if len(s.MapAnyField) != len(to.MapAnyField) {
		return false
	}
	for k6 := range s.MapAnyField {
		valB6, ok := to.MapAnyField[k6]
		if !ok {
			return false
		}
		valA6 := s.MapAnyField[k6]
		if valA6 != valB6 {
			return false
		}
	}
	// map comparision
	if len(s.MapAnyPtrField) != len(to.MapAnyPtrField) {
		return false
	}
	for k7 := range s.MapAnyPtrField {
		valB7, ok := to.MapAnyPtrField[k7]
		if !ok {
			return false
		}
		valA7 := s.MapAnyPtrField[k7]
		if (valA7 == nil) != (valB7 == nil) {
			return false
		}
		if valA7 != nil && valB7 != nil {
			if (*valA7) != (*valB7) {
				return false
			}
		}
	}
	if len(s.ModelSliceField) != len(to.ModelSliceField) {
		return false
	}
	for idx := range s.ModelSliceField {
		if !s.ModelSliceField[idx].Equals(&to.ModelSliceField[idx]) {
			return false
		}
	}
	if len(s.ModelPtrSliceField) != len(to.ModelPtrSliceField) {
		return false
	}
	for idx1 := range s.ModelPtrSliceField {
		if (s.ModelPtrSliceField[idx1] == nil) != (to.ModelPtrSliceField[idx1] == nil) {
			return false
		}
		if s.ModelPtrSliceField[idx1] != nil && to.ModelPtrSliceField[idx1] != nil {
			if !(*s.ModelPtrSliceField[idx1]).Equals(&(*to.ModelPtrSliceField[idx1])) {
				return false
			}
		}
	}
	if len(s.OneOfSliceField) != len(to.OneOfSliceField) {
		return false
	}
	for idx2 := range s.OneOfSliceField {
		if !s.OneOfSliceField[idx2].SomeOneOfEquals(to.OneOfSliceField[idx2]) {
			return false
		}
	}
	if len(s.OneOfPtrSliceField) != len(to.OneOfPtrSliceField) {
		return false
	}
	for idx3 := range s.OneOfPtrSliceField {
		if (s.OneOfPtrSliceField[idx3] == nil) != (to.OneOfPtrSliceField[idx3] == nil) {
			return false
		}
		if s.OneOfPtrSliceField[idx3] != nil && to.OneOfPtrSliceField[idx3] != nil {
			if !(*s.OneOfPtrSliceField[idx3]).SomeOneOfEquals((*to.OneOfPtrSliceField[idx3])) {
				return false
			}
		}
	}
	if len(s.SliceEnumField) != len(to.SliceEnumField) {
		return false
	}
	for idx4 := range s.SliceEnumField {
		if s.SliceEnumField[idx4] != to.SliceEnumField[idx4] {
			return false
		}
	}
	if len(s.SliceEnumPtrField) != len(to.SliceEnumPtrField) {
		return false
	}
	for idx5 := range s.SliceEnumPtrField {
		if (s.SliceEnumPtrField[idx5] == nil) != (to.SliceEnumPtrField[idx5] == nil) {
			return false
		}
		if s.SliceEnumPtrField[idx5] != nil && to.SliceEnumPtrField[idx5] != nil {
			if (*s.SliceEnumPtrField[idx5]) != (*to.SliceEnumPtrField[idx5]) {
				return false
			}
		}
	}
	if len(s.SliceAnyField) != len(to.SliceAnyField) {
		return false
	}
	for idx6 := range s.SliceAnyField {
		if s.SliceAnyField[idx6] != to.SliceAnyField[idx6] {
			return false
		}
	}
	if len(s.SliceAnyPtrField) != len(to.SliceAnyPtrField) {
		return false
	}
	for idx7 := range s.SliceAnyPtrField {
		if (s.SliceAnyPtrField[idx7] == nil) != (to.SliceAnyPtrField[idx7] == nil) {
			return false
		}
		if s.SliceAnyPtrField[idx7] != nil && to.SliceAnyPtrField[idx7] != nil {
			if (*s.SliceAnyPtrField[idx7]) != (*to.SliceAnyPtrField[idx7]) {
				return false
			}
		}
	}

	return true
}

type SomeOtherModelField byte

const (
	SomeOtherModelFieldID SomeOtherModelField = iota + 1
)

type SomeOtherModelFilter struct {
	ID  filter.Filter[uuid.UUID]
	Or  []*SomeOtherModelFilter
	And []*SomeOtherModelFilter
}
type SomeOtherModelOrder order.Order[SomeOtherModelField]

type SomeOtherModel struct {
	ID uuid.UUID
}

// user code 'SomeOtherModel methods'
// end user code 'SomeOtherModel methods'

func (s *SomeOtherModel) Copy() SomeOtherModel {
	var result SomeOtherModel
	result.ID = s.ID

	return result
}
func (s *SomeOtherModel) Equals(to *SomeOtherModel) bool {
	if (s == nil) != (to == nil) {
		return false
	}
	if s == nil && to == nil {
		return true
	}
	if s.ID != to.ID {
		return false
	}

	return true
}

type OneOfValue1Field byte

const (
	OneOfValue1FieldValue OneOfValue1Field = iota + 1
)

type OneOfValue1Filter struct {
	Or  []*OneOfValue1Filter
	And []*OneOfValue1Filter
}
type OneOfValue1Order order.Order[OneOfValue1Field]

type OneOfValue1 struct {
	Value string
}

// user code 'OneOfValue1 methods'
// end user code 'OneOfValue1 methods'

func (o *OneOfValue1) Copy() OneOfValue1 {
	var result OneOfValue1
	result.Value = o.Value

	return result
}
func (o *OneOfValue1) Equals(to *OneOfValue1) bool {
	if (o == nil) != (to == nil) {
		return false
	}
	if o == nil && to == nil {
		return true
	}
	if o.Value != to.Value {
		return false
	}

	return true
}

type OneOfValue2Field byte

const (
	OneOfValue2FieldValue OneOfValue2Field = iota + 1
)

type OneOfValue2Filter struct {
	Or  []*OneOfValue2Filter
	And []*OneOfValue2Filter
}
type OneOfValue2Order order.Order[OneOfValue2Field]

type OneOfValue2 struct {
	Value string
}

// user code 'OneOfValue2 methods'
// end user code 'OneOfValue2 methods'

func (o *OneOfValue2) Copy() OneOfValue2 {
	var result OneOfValue2
	result.Value = o.Value

	return result
}
func (o *OneOfValue2) Equals(to *OneOfValue2) bool {
	if (o == nil) != (to == nil) {
		return false
	}
	if o == nil && to == nil {
		return true
	}
	if o.Value != to.Value {
		return false
	}

	return true
}
