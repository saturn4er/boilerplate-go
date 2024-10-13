package testservice

import (
	uuid "github.com/google/uuid"

	filter "github.com/saturn4er/boilerplate-go/lib/filter"
)

type SomeOneOf interface {
	isSomeOneOf()
	// user code 'SomeOneOf methods'
	// end user code 'SomeOneOf methods'
}

func (*OneOfValue1) isSomeOneOf() {}
func (*OneOfValue2) isSomeOneOf() {}

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

type PasswordRecoveryEventData interface {
	isPasswordRecoveryEventData()
	// user code 'PasswordRecoveryEventData methods'
	// end user code 'PasswordRecoveryEventData methods'
}

func (*PasswordRecoveryRequestedEventData) isPasswordRecoveryEventData() {}
func (*PasswordRecoveryCompletedEventData) isPasswordRecoveryEventData() {}

func copyPasswordRecoveryEventData(val PasswordRecoveryEventData) PasswordRecoveryEventData {
	if val == nil {
		return nil
	}

	switch val := val.(type) {
	case *PasswordRecoveryRequestedEventData:
		valCopy := val.Copy()
		return &valCopy
	case *PasswordRecoveryCompletedEventData:
		valCopy := val.Copy()
		return &valCopy
	}
	panic("called copyPasswordRecoveryEventData with invalid type")
}

type SomeModelField byte

const (
	SomeModelFieldID SomeModelField = iota + 1
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
	ID  filter.Filter[uuid.UUID]
	Or  []*SomeModelFilter
	And []*SomeModelFilter
}

type SomeModel struct {
	ID                 uuid.UUID
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
func (s SomeModel) Copy() SomeModel {
	var result SomeModel
	result.ID = s.ID
	result.ModelField = s.ModelField.Copy() // model
	if s.ModelPtrField != nil {
		var tmp SomeOtherModel
		tmp = (*s.ModelPtrField).Copy() // model
		result.ModelPtrField = &tmp
	}
	result.OneOfField = copySomeOneOf(s.OneOfField)
	if s.OneOfPtrField != nil {
		var tmp1 SomeOneOf
		tmp1 = copySomeOneOf((*s.OneOfPtrField))
		result.OneOfPtrField = &tmp1
	}
	result.EnumField = s.EnumField // enum
	if s.EnumPtrField != nil {
		var tmp2 SomeEnum
		tmp2 = (*s.EnumPtrField) // enum
		result.EnumPtrField = &tmp2
	}
	result.AnyField = s.AnyField
	if s.AnyPtrField != nil {
		var tmp3 any
		tmp3 = (*s.AnyPtrField)
		result.AnyPtrField = &tmp3
	}
	tmp4 := make(map[string]SomeOtherModel)
	for k, v := range s.MapModelField {
		var keyCopy string
		var valueCopy SomeOtherModel
		keyCopy = k
		valueCopy = v.Copy() // model
		tmp4[keyCopy] = valueCopy
	}
	result.MapModelField = tmp4
	tmp5 := make(map[string]*SomeOtherModel)
	for k, v := range s.MapModelPtrField {
		var keyCopy1 string
		var valueCopy1 *SomeOtherModel
		keyCopy1 = k
		if v != nil {
			var tmp6 SomeOtherModel
			tmp6 = (*v).Copy() // model
			valueCopy1 = &tmp6
		}
		tmp5[keyCopy1] = valueCopy1
	}
	result.MapModelPtrField = tmp5
	tmp7 := make(map[string]SomeOneOf)
	for k, v := range s.MapOneOfField {
		var keyCopy2 string
		var valueCopy2 SomeOneOf
		keyCopy2 = k
		valueCopy2 = copySomeOneOf(v)
		tmp7[keyCopy2] = valueCopy2
	}
	result.MapOneOfField = tmp7
	tmp8 := make(map[string]*SomeOneOf)
	for k, v := range s.MapOneOfPtrField {
		var keyCopy3 string
		var valueCopy3 *SomeOneOf
		keyCopy3 = k
		if v != nil {
			var tmp9 SomeOneOf
			tmp9 = copySomeOneOf((*v))
			valueCopy3 = &tmp9
		}
		tmp8[keyCopy3] = valueCopy3
	}
	result.MapOneOfPtrField = tmp8
	tmp10 := make(map[string]SomeEnum)
	for k, v := range s.MapEnumField {
		var keyCopy4 string
		var valueCopy4 SomeEnum
		keyCopy4 = k
		valueCopy4 = v // enum
		tmp10[keyCopy4] = valueCopy4
	}
	result.MapEnumField = tmp10
	tmp11 := make(map[string]*SomeEnum)
	for k, v := range s.MapEnumPtrField {
		var keyCopy5 string
		var valueCopy5 *SomeEnum
		keyCopy5 = k
		if v != nil {
			var tmp12 SomeEnum
			tmp12 = (*v) // enum
			valueCopy5 = &tmp12
		}
		tmp11[keyCopy5] = valueCopy5
	}
	result.MapEnumPtrField = tmp11
	tmp13 := make(map[string]any)
	for k, v := range s.MapAnyField {
		var keyCopy6 string
		var valueCopy6 any
		keyCopy6 = k
		valueCopy6 = v
		tmp13[keyCopy6] = valueCopy6
	}
	result.MapAnyField = tmp13
	tmp14 := make(map[string]*any)
	for k, v := range s.MapAnyPtrField {
		var keyCopy7 string
		var valueCopy7 *any
		keyCopy7 = k
		if v != nil {
			var tmp15 any
			tmp15 = (*v)
			valueCopy7 = &tmp15
		}
		tmp14[keyCopy7] = valueCopy7
	}
	result.MapAnyPtrField = tmp14
	tmp16 := make([]SomeOtherModel, 0, len(s.ModelSliceField))
	for _, i := range s.ModelSliceField {
		var itemCopy SomeOtherModel
		itemCopy = i.Copy() // model
		tmp16 = append(tmp16, itemCopy)
	}
	result.ModelSliceField = tmp16
	tmp17 := make([]*SomeOtherModel, 0, len(s.ModelPtrSliceField))
	for _, i1 := range s.ModelPtrSliceField {
		var itemCopy1 *SomeOtherModel
		if i1 != nil {
			var tmp18 SomeOtherModel
			tmp18 = (*i1).Copy() // model
			itemCopy1 = &tmp18
		}
		tmp17 = append(tmp17, itemCopy1)
	}
	result.ModelPtrSliceField = tmp17
	tmp19 := make([]SomeOneOf, 0, len(s.OneOfSliceField))
	for _, i2 := range s.OneOfSliceField {
		var itemCopy2 SomeOneOf
		itemCopy2 = copySomeOneOf(i2)
		tmp19 = append(tmp19, itemCopy2)
	}
	result.OneOfSliceField = tmp19
	tmp20 := make([]*SomeOneOf, 0, len(s.OneOfPtrSliceField))
	for _, i3 := range s.OneOfPtrSliceField {
		var itemCopy3 *SomeOneOf
		if i3 != nil {
			var tmp21 SomeOneOf
			tmp21 = copySomeOneOf((*i3))
			itemCopy3 = &tmp21
		}
		tmp20 = append(tmp20, itemCopy3)
	}
	result.OneOfPtrSliceField = tmp20
	tmp22 := make([]SomeEnum, 0, len(s.SliceEnumField))
	for _, i4 := range s.SliceEnumField {
		var itemCopy4 SomeEnum
		itemCopy4 = i4 // enum
		tmp22 = append(tmp22, itemCopy4)
	}
	result.SliceEnumField = tmp22
	tmp23 := make([]*SomeEnum, 0, len(s.SliceEnumPtrField))
	for _, i5 := range s.SliceEnumPtrField {
		var itemCopy5 *SomeEnum
		if i5 != nil {
			var tmp24 SomeEnum
			tmp24 = (*i5) // enum
			itemCopy5 = &tmp24
		}
		tmp23 = append(tmp23, itemCopy5)
	}
	result.SliceEnumPtrField = tmp23
	tmp25 := make([]any, 0, len(s.SliceAnyField))
	for _, i6 := range s.SliceAnyField {
		var itemCopy6 any
		itemCopy6 = i6
		tmp25 = append(tmp25, itemCopy6)
	}
	result.SliceAnyField = tmp25
	tmp26 := make([]*any, 0, len(s.SliceAnyPtrField))
	for _, i7 := range s.SliceAnyPtrField {
		var itemCopy7 *any
		if i7 != nil {
			var tmp27 any
			tmp27 = (*i7)
			itemCopy7 = &tmp27
		}
		tmp26 = append(tmp26, itemCopy7)
	}
	result.SliceAnyPtrField = tmp26

	return result
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

type SomeOtherModel struct {
	ID uuid.UUID
}

// user code 'SomeOtherModel methods'
// end user code 'SomeOtherModel methods'
func (s SomeOtherModel) Copy() SomeOtherModel {
	var result SomeOtherModel
	result.ID = s.ID

	return result
}

type OneOfValue1 struct {
	Value string
}

// user code 'OneOfValue1 methods'
// end user code 'OneOfValue1 methods'
func (o OneOfValue1) Copy() OneOfValue1 {
	var result OneOfValue1
	result.Value = o.Value

	return result
}

type OneOfValue2 struct {
	Value string
}

// user code 'OneOfValue2 methods'
// end user code 'OneOfValue2 methods'
func (o OneOfValue2) Copy() OneOfValue2 {
	var result OneOfValue2
	result.Value = o.Value

	return result
}

type PasswordRecoveryEventField byte

const (
	PasswordRecoveryEventFieldID PasswordRecoveryEventField = iota + 1
	PasswordRecoveryEventFieldData
	PasswordRecoveryEventFieldIdempotencyKey
)

type PasswordRecoveryEventFilter struct {
	Or  []*PasswordRecoveryEventFilter
	And []*PasswordRecoveryEventFilter
}

type PasswordRecoveryEvent struct {
	ID             uuid.UUID
	Data           PasswordRecoveryEventData
	IdempotencyKey string
}

// user code 'PasswordRecoveryEvent methods'
// end user code 'PasswordRecoveryEvent methods'
func (p PasswordRecoveryEvent) Copy() PasswordRecoveryEvent {
	var result PasswordRecoveryEvent
	result.ID = p.ID
	result.Data = copyPasswordRecoveryEventData(p.Data)
	result.IdempotencyKey = p.IdempotencyKey

	return result
}

type PasswordRecoveryRequestedEventData struct {
	Email            string
	UserID           uuid.UUID
	VerificationCode string
	NestedData       PasswordRecoveryEventData
}

// user code 'PasswordRecoveryRequestedEventData methods'
// end user code 'PasswordRecoveryRequestedEventData methods'
func (p PasswordRecoveryRequestedEventData) Copy() PasswordRecoveryRequestedEventData {
	var result PasswordRecoveryRequestedEventData
	result.Email = p.Email
	result.UserID = p.UserID
	result.VerificationCode = p.VerificationCode
	result.NestedData = copyPasswordRecoveryEventData(p.NestedData)

	return result
}

type PasswordRecoveryCompletedEventData struct {
	Email  string
	UserID uuid.UUID
}

// user code 'PasswordRecoveryCompletedEventData methods'
// end user code 'PasswordRecoveryCompletedEventData methods'
func (p PasswordRecoveryCompletedEventData) Copy() PasswordRecoveryCompletedEventData {
	var result PasswordRecoveryCompletedEventData
	result.Email = p.Email
	result.UserID = p.UserID

	return result
}
