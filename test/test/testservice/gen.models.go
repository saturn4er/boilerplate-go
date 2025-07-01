package testservice

import (
	uuid "github.com/google/uuid"

	filter "github.com/saturn4er/boilerplate-go/lib/filter"
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

type PasswordRecoveryEventData interface {
	isPasswordRecoveryEventData()
	PasswordRecoveryEventDataEquals(PasswordRecoveryEventData) bool
	// user code 'PasswordRecoveryEventData methods'
	// end user code 'PasswordRecoveryEventData methods'
}

func (*PasswordRecoveryRequestedEventData) isPasswordRecoveryEventData() {}
func (p *PasswordRecoveryRequestedEventData) PasswordRecoveryEventDataEquals(to PasswordRecoveryEventData) bool {
	if (p == nil) != (to == nil) {
		return false
	}
	if p == nil && to == nil {
		return true
	}

	toTyped, ok := to.(*PasswordRecoveryRequestedEventData)
	if !ok {
		return false
	}

	return p.Equals(toTyped)
}
func (*PasswordRecoveryCompletedEventData) isPasswordRecoveryEventData() {}
func (p *PasswordRecoveryCompletedEventData) PasswordRecoveryEventDataEquals(to PasswordRecoveryEventData) bool {
	if (p == nil) != (to == nil) {
		return false
	}
	if p == nil && to == nil {
		return true
	}

	toTyped, ok := to.(*PasswordRecoveryCompletedEventData)
	if !ok {
		return false
	}

	return p.Equals(toTyped)
}

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

func (s *SomeModel) Copy() SomeModel {
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
	for i := range s.ModelSliceField {
		if !s.ModelSliceField[i].Equals(&to.ModelSliceField[i]) {
			return false
		}
	}
	if len(s.ModelPtrSliceField) != len(to.ModelPtrSliceField) {
		return false
	}
	for i1 := range s.ModelPtrSliceField {
		if (s.ModelPtrSliceField[i1] == nil) != (to.ModelPtrSliceField[i1] == nil) {
			return false
		}
		if s.ModelPtrSliceField[i1] != nil && to.ModelPtrSliceField[i1] != nil {
			if !(*s.ModelPtrSliceField[i1]).Equals(&(*to.ModelPtrSliceField[i1])) {
				return false
			}
		}
	}
	if len(s.OneOfSliceField) != len(to.OneOfSliceField) {
		return false
	}
	for i2 := range s.OneOfSliceField {
		if !s.OneOfSliceField[i2].SomeOneOfEquals(to.OneOfSliceField[i2]) {
			return false
		}
	}
	if len(s.OneOfPtrSliceField) != len(to.OneOfPtrSliceField) {
		return false
	}
	for i3 := range s.OneOfPtrSliceField {
		if (s.OneOfPtrSliceField[i3] == nil) != (to.OneOfPtrSliceField[i3] == nil) {
			return false
		}
		if s.OneOfPtrSliceField[i3] != nil && to.OneOfPtrSliceField[i3] != nil {
			if !(*s.OneOfPtrSliceField[i3]).SomeOneOfEquals((*to.OneOfPtrSliceField[i3])) {
				return false
			}
		}
	}
	if len(s.SliceEnumField) != len(to.SliceEnumField) {
		return false
	}
	for i4 := range s.SliceEnumField {
		if s.SliceEnumField[i4] != to.SliceEnumField[i4] {
			return false
		}
	}
	if len(s.SliceEnumPtrField) != len(to.SliceEnumPtrField) {
		return false
	}
	for i5 := range s.SliceEnumPtrField {
		if (s.SliceEnumPtrField[i5] == nil) != (to.SliceEnumPtrField[i5] == nil) {
			return false
		}
		if s.SliceEnumPtrField[i5] != nil && to.SliceEnumPtrField[i5] != nil {
			if (*s.SliceEnumPtrField[i5]) != (*to.SliceEnumPtrField[i5]) {
				return false
			}
		}
	}
	if len(s.SliceAnyField) != len(to.SliceAnyField) {
		return false
	}
	for i6 := range s.SliceAnyField {
		if s.SliceAnyField[i6] != to.SliceAnyField[i6] {
			return false
		}
	}
	if len(s.SliceAnyPtrField) != len(to.SliceAnyPtrField) {
		return false
	}
	for i7 := range s.SliceAnyPtrField {
		if (s.SliceAnyPtrField[i7] == nil) != (to.SliceAnyPtrField[i7] == nil) {
			return false
		}
		if s.SliceAnyPtrField[i7] != nil && to.SliceAnyPtrField[i7] != nil {
			if (*s.SliceAnyPtrField[i7]) != (*to.SliceAnyPtrField[i7]) {
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

func (p *PasswordRecoveryEvent) Copy() PasswordRecoveryEvent {
	var result PasswordRecoveryEvent
	result.ID = p.ID
	result.Data = copyPasswordRecoveryEventData(p.Data)
	result.IdempotencyKey = p.IdempotencyKey

	return result
}
func (p *PasswordRecoveryEvent) Equals(to *PasswordRecoveryEvent) bool {
	if (p == nil) != (to == nil) {
		return false
	}
	if p == nil && to == nil {
		return true
	}
	if p.ID != to.ID {
		return false
	}
	if !p.Data.PasswordRecoveryEventDataEquals(to.Data) {
		return false
	}
	if p.IdempotencyKey != to.IdempotencyKey {
		return false
	}

	return true
}

type PasswordRecoveryRequestedEventDataField byte

const (
	PasswordRecoveryRequestedEventDataFieldEmail PasswordRecoveryRequestedEventDataField = iota + 1
	PasswordRecoveryRequestedEventDataFieldUserID
	PasswordRecoveryRequestedEventDataFieldVerificationCode
	PasswordRecoveryRequestedEventDataFieldNestedData
)

type PasswordRecoveryRequestedEventDataFilter struct {
	UserID filter.Filter[uuid.UUID]
	Or     []*PasswordRecoveryRequestedEventDataFilter
	And    []*PasswordRecoveryRequestedEventDataFilter
}

type PasswordRecoveryRequestedEventData struct {
	Email            string
	UserID           uuid.UUID
	VerificationCode string
	NestedData       PasswordRecoveryEventData
}

// user code 'PasswordRecoveryRequestedEventData methods'
// end user code 'PasswordRecoveryRequestedEventData methods'

func (p *PasswordRecoveryRequestedEventData) Copy() PasswordRecoveryRequestedEventData {
	var result PasswordRecoveryRequestedEventData
	result.Email = p.Email
	result.UserID = p.UserID
	result.VerificationCode = p.VerificationCode
	result.NestedData = copyPasswordRecoveryEventData(p.NestedData)

	return result
}
func (p *PasswordRecoveryRequestedEventData) Equals(to *PasswordRecoveryRequestedEventData) bool {
	if (p == nil) != (to == nil) {
		return false
	}
	if p == nil && to == nil {
		return true
	}
	if p.Email != to.Email {
		return false
	}
	if p.UserID != to.UserID {
		return false
	}
	if p.VerificationCode != to.VerificationCode {
		return false
	}
	if !p.NestedData.PasswordRecoveryEventDataEquals(to.NestedData) {
		return false
	}

	return true
}

type PasswordRecoveryCompletedEventDataField byte

const (
	PasswordRecoveryCompletedEventDataFieldEmail PasswordRecoveryCompletedEventDataField = iota + 1
	PasswordRecoveryCompletedEventDataFieldUserID
)

type PasswordRecoveryCompletedEventDataFilter struct {
	UserID filter.Filter[uuid.UUID]
	Or     []*PasswordRecoveryCompletedEventDataFilter
	And    []*PasswordRecoveryCompletedEventDataFilter
}

type PasswordRecoveryCompletedEventData struct {
	Email  string
	UserID uuid.UUID
}

// user code 'PasswordRecoveryCompletedEventData methods'
// end user code 'PasswordRecoveryCompletedEventData methods'

func (p *PasswordRecoveryCompletedEventData) Copy() PasswordRecoveryCompletedEventData {
	var result PasswordRecoveryCompletedEventData
	result.Email = p.Email
	result.UserID = p.UserID

	return result
}
func (p *PasswordRecoveryCompletedEventData) Equals(to *PasswordRecoveryCompletedEventData) bool {
	if (p == nil) != (to == nil) {
		return false
	}
	if p == nil && to == nil {
		return true
	}
	if p.Email != to.Email {
		return false
	}
	if p.UserID != to.UserID {
		return false
	}

	return true
}
