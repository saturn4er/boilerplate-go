package teststorage

import (
	driver "database/sql/driver"
	json "encoding/json"
	fmt "fmt"
	time "time"

	uuid "github.com/google/uuid"
	errors "github.com/pkg/errors"

	testservice "github.com/saturn4er/boilerplate-go/test/test/testservice"
	// user code 'imports'
	// end user code 'imports'
)

type jsonSomeModel struct {
	ID                 uuid.UUID                             `json:"id"`
	Name               string                                `json:"name"`
	Description        *string                               `json:"description"`
	ModelField         jsonSomeOtherModel                    `json:"model_field"`
	ModelPtrField      *jsonSomeOtherModel                   `json:"model_ptr_field"`
	OneOfField         *jsonSomeOneOf                        `json:"one_of_field"`
	OneOfPtrField      *jsonSomeOneOf                        `json:"one_of_ptr_field"`
	EnumField          string                                `json:"enum_field"`
	EnumPtrField       *string                               `json:"enum_ptr_field"`
	AnyField           *string                               `json:"any_field"`
	AnyPtrField        *string                               `json:"any_ptr_field"`
	MapModelField      mapValue[string, jsonSomeOtherModel]  `json:"map_model_field"`
	MapModelPtrField   mapValue[string, *jsonSomeOtherModel] `json:"map_model_ptr_field"`
	MapOneOfField      mapValue[string, *jsonSomeOneOf]      `json:"map_one_of_field"`
	MapOneOfPtrField   mapValue[string, *jsonSomeOneOf]      `json:"map_one_of_ptr_field"`
	MapEnumField       mapValue[string, string]              `json:"map_enum_field"`
	MapEnumPtrField    mapValue[string, *string]             `json:"map_enum_ptr_field"`
	MapAnyField        mapValue[string, *string]             `json:"map_any_field"`
	MapAnyPtrField     mapValue[string, *string]             `json:"map_any_ptr_field"`
	ModelSliceField    sliceValue[jsonSomeOtherModel]        `json:"model_slice_field"`
	ModelPtrSliceField sliceValue[*jsonSomeOtherModel]       `json:"model_ptr_slice_field"`
	OneOfSliceField    sliceValue[*jsonSomeOneOf]            `json:"one_of_slice_field"`
	OneOfPtrSliceField sliceValue[*jsonSomeOneOf]            `json:"one_of_ptr_slice_field"`
	SliceEnumField     stringSliceValue                      `json:"slice_enum_field"`
	SliceEnumPtrField  sliceValue[*string]                   `json:"slice_enum_ptr_field"`
	SliceAnyField      sliceValue[*string]                   `json:"slice_any_field"`
	SliceAnyPtrField   sliceValue[*string]                   `json:"slice_any_ptr_field"`
	UpdatedAt          time.Time                             `json:"updated_at"`
}

func (s *jsonSomeModel) Scan(value any) error {
	return json.Unmarshal(value.([]byte), s)
}

func (s jsonSomeModel) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func convertSomeModelToJsonModel(src *testservice.SomeModel) (*jsonSomeModel, error) {
	result := &jsonSomeModel{}
	result.ID = src.ID
	result.Name = src.Name
	if src.Description == nil {
		result.Description = nil
	} else {
		result.Description = toPtr(fromPtr(src.Description))
	}
	tmp4, err := convertSomeOtherModelToJsonModel(toPtr(src.ModelField))
	if err != nil {
		return nil, errors.Wrap(err, "convert SomeOtherModel to db")
	}
	result.ModelField = *tmp4
	if src.ModelPtrField != nil {
		tmp5, err := convertSomeOtherModelToJsonModel(src.ModelPtrField)
		if err != nil {
			return nil, errors.Wrap(err, "convert SomeOtherModel to db")
		}
		result.ModelPtrField = tmp5
	} else {
		result.ModelPtrField = nil
	}
	tmp6, err := convertSomeOneOfToDB(src.OneOfField)
	if err != nil {
		return nil, err
	}
	result.OneOfField = tmp6
	if src.OneOfPtrField != nil {
		tmp7, err := convertSomeOneOfToDB(*src.OneOfPtrField)
		if err != nil {
			return nil, err
		}
		result.OneOfPtrField = tmp7
	} else {
		result.OneOfPtrField = nil
	}
	tmp8, err := convertSomeEnumToDB(src.EnumField)
	if err != nil {
		return nil, err
	}
	result.EnumField = tmp8
	if src.EnumPtrField == nil {
		result.EnumPtrField = nil
	} else {
		tmp10, err := convertSomeEnumToDB(fromPtr(src.EnumPtrField))
		if err != nil {
			return nil, err
		}
		result.EnumPtrField = toPtr(tmp10)
	}
	if src.AnyField != nil {
		tmp11, err := json.Marshal(src.AnyField)
		if err != nil {
			return nil, err
		}

		marshaledValue := string(tmp11)
		result.AnyField = toPtr(marshaledValue)
	} else {
		result.AnyField = nil
	}
	if src.AnyPtrField != nil && fromPtr(src.AnyPtrField) != nil {
		tmp12, err := json.Marshal(*src.AnyPtrField)
		if err != nil {
			return nil, err
		}

		marshaledValue1 := string(tmp12)
		result.AnyPtrField = toPtr(marshaledValue1)
	} else {
		result.AnyPtrField = nil
	}
	tmp13 := make(mapValue[string, jsonSomeOtherModel], len(src.MapModelField))
	for k, v := range src.MapModelField {
		tmp14, err := convertSomeOtherModelToJsonModel(toPtr(v))
		if err != nil {
			return nil, errors.Wrap(err, "convert SomeOtherModel to db")
		}
		tmp13[k] = *tmp14
	}
	result.MapModelField = tmp13
	tmp15 := make(mapValue[string, *jsonSomeOtherModel], len(src.MapModelPtrField))
	for k1, v1 := range src.MapModelPtrField {
		if v1 != nil {
			tmp16, err := convertSomeOtherModelToJsonModel(v1)
			if err != nil {
				return nil, errors.Wrap(err, "convert SomeOtherModel to db")
			}
			tmp15[k1] = tmp16
		} else {
			tmp15[k1] = nil
		}
	}
	result.MapModelPtrField = tmp15
	tmp17 := make(mapValue[string, *jsonSomeOneOf], len(src.MapOneOfField))
	for k2, v2 := range src.MapOneOfField {
		tmp18, err := convertSomeOneOfToDB(v2)
		if err != nil {
			return nil, err
		}
		tmp17[k2] = tmp18
	}
	result.MapOneOfField = tmp17
	tmp19 := make(mapValue[string, *jsonSomeOneOf], len(src.MapOneOfPtrField))
	for k3, v3 := range src.MapOneOfPtrField {
		if v3 != nil {
			tmp20, err := convertSomeOneOfToDB(*v3)
			if err != nil {
				return nil, err
			}
			tmp19[k3] = tmp20
		} else {
			tmp19[k3] = nil
		}
	}
	result.MapOneOfPtrField = tmp19
	tmp21 := make(mapValue[string, string], len(src.MapEnumField))
	for k4, v4 := range src.MapEnumField {
		tmp22, err := convertSomeEnumToDB(v4)
		if err != nil {
			return nil, err
		}
		tmp21[k4] = tmp22
	}
	result.MapEnumField = tmp21
	tmp23 := make(mapValue[string, *string], len(src.MapEnumPtrField))
	for k5, v5 := range src.MapEnumPtrField {
		if v5 == nil {
			tmp23[k5] = nil
		} else {
			tmp25, err := convertSomeEnumToDB(fromPtr(v5))
			if err != nil {
				return nil, err
			}
			tmp23[k5] = toPtr(tmp25)
		}
	}
	result.MapEnumPtrField = tmp23
	tmp26 := make(mapValue[string, *string], len(src.MapAnyField))
	for k6, v6 := range src.MapAnyField {
		if v6 != nil {
			tmp27, err := json.Marshal(v6)
			if err != nil {
				return nil, err
			}

			marshaledValue2 := string(tmp27)
			tmp26[k6] = toPtr(marshaledValue2)
		} else {
			tmp26[k6] = nil
		}
	}
	result.MapAnyField = tmp26
	tmp28 := make(mapValue[string, *string], len(src.MapAnyPtrField))
	for k7, v7 := range src.MapAnyPtrField {
		if v7 != nil && fromPtr(v7) != nil {
			tmp29, err := json.Marshal(*v7)
			if err != nil {
				return nil, err
			}

			marshaledValue3 := string(tmp29)
			tmp28[k7] = toPtr(marshaledValue3)
		} else {
			tmp28[k7] = nil
		}
	}
	result.MapAnyPtrField = tmp28
	tmp30 := make(sliceValue[jsonSomeOtherModel], 0, len(src.ModelSliceField))
	for _, el := range src.ModelSliceField {
		tmp31, err := convertSomeOtherModelToJsonModel(toPtr(el))
		if err != nil {
			return nil, errors.Wrap(err, "convert SomeOtherModel to db")
		}
		tmp30 = append(tmp30, *tmp31)
	}
	result.ModelSliceField = tmp30
	tmp32 := make(sliceValue[*jsonSomeOtherModel], 0, len(src.ModelPtrSliceField))
	for _, el := range src.ModelPtrSliceField {
		if el != nil {
			tmp33, err := convertSomeOtherModelToJsonModel(el)
			if err != nil {
				return nil, errors.Wrap(err, "convert SomeOtherModel to db")
			}
			tmp32 = append(tmp32, tmp33)
		} else {
			tmp32 = append(tmp32, nil)
		}
	}
	result.ModelPtrSliceField = tmp32
	tmp34 := make(sliceValue[*jsonSomeOneOf], 0, len(src.OneOfSliceField))
	for _, el := range src.OneOfSliceField {
		tmp35, err := convertSomeOneOfToDB(el)
		if err != nil {
			return nil, err
		}
		tmp34 = append(tmp34, tmp35)
	}
	result.OneOfSliceField = tmp34
	tmp36 := make(sliceValue[*jsonSomeOneOf], 0, len(src.OneOfPtrSliceField))
	for _, el := range src.OneOfPtrSliceField {
		if el != nil {
			tmp37, err := convertSomeOneOfToDB(*el)
			if err != nil {
				return nil, err
			}
			tmp36 = append(tmp36, tmp37)
		} else {
			tmp36 = append(tmp36, nil)
		}
	}
	result.OneOfPtrSliceField = tmp36
	tmp38 := make(stringSliceValue, 0, len(src.SliceEnumField))
	for _, el := range src.SliceEnumField {
		tmp39, err := convertSomeEnumToDB(el)
		if err != nil {
			return nil, err
		}
		tmp38 = append(tmp38, tmp39)
	}
	result.SliceEnumField = tmp38
	tmp40 := make(sliceValue[*string], 0, len(src.SliceEnumPtrField))
	for _, el := range src.SliceEnumPtrField {
		if el == nil {
			tmp40 = append(tmp40, nil)
		} else {
			tmp42, err := convertSomeEnumToDB(fromPtr(el))
			if err != nil {
				return nil, err
			}
			tmp40 = append(tmp40, toPtr(tmp42))
		}
	}
	result.SliceEnumPtrField = tmp40
	tmp43 := make(sliceValue[*string], 0, len(src.SliceAnyField))
	for _, el := range src.SliceAnyField {
		if el != nil {
			tmp44, err := json.Marshal(el)
			if err != nil {
				return nil, err
			}

			marshaledValue4 := string(tmp44)
			tmp43 = append(tmp43, toPtr(marshaledValue4))
		} else {
			tmp43 = append(tmp43, nil)
		}
	}
	result.SliceAnyField = tmp43
	tmp45 := make(sliceValue[*string], 0, len(src.SliceAnyPtrField))
	for _, el := range src.SliceAnyPtrField {
		if el != nil && fromPtr(el) != nil {
			tmp46, err := json.Marshal(*el)
			if err != nil {
				return nil, err
			}

			marshaledValue5 := string(tmp46)
			tmp45 = append(tmp45, toPtr(marshaledValue5))
		} else {
			tmp45 = append(tmp45, nil)
		}
	}
	result.SliceAnyPtrField = tmp45
	result.UpdatedAt = (src.UpdatedAt).UTC()
	return result, nil
}

func convertSomeModelFromJsonModel(src *jsonSomeModel) (*testservice.SomeModel, error) {
	result := &testservice.SomeModel{}
	result.ID = src.ID
	result.Name = src.Name
	if src.Description == nil {
		result.Description = nil
	} else {
		result.Description = toPtr(fromPtr(src.Description))
	}
	tmp52, err := convertSomeOtherModelFromJsonModel(toPtr(src.ModelField))
	if err != nil {
		return nil, err
	}

	result.ModelField = fromPtr(tmp52)
	if src.ModelPtrField != nil {
		tmp53, err := convertSomeOtherModelFromJsonModel(src.ModelPtrField)
		if err != nil {
			return nil, err
		}
		result.ModelPtrField = tmp53
	} else {
		result.ModelPtrField = nil
	}
	tmp54, err := convertSomeOneOfFromDB(src.OneOfField)
	if err != nil {
		return nil, fmt.Errorf("convert SomeOneOf to service type: %w", err)
	}
	result.OneOfField = tmp54
	if src.OneOfPtrField != nil {
		tmp55, err := convertSomeOneOfFromDB(src.OneOfPtrField)
		if err != nil {
			return nil, fmt.Errorf("convert SomeOneOf to service type: %w", err)
		}
		result.OneOfPtrField = toPtr(tmp55)
	} else {
		result.OneOfPtrField = nil
	}
	tmp56, err := convertSomeEnumFromDB(src.EnumField)
	if err != nil {
		return nil, err
	}
	result.EnumField = tmp56
	if src.EnumPtrField == nil {
		result.EnumPtrField = nil
	} else {
		tmp58, err := convertSomeEnumFromDB(fromPtr(src.EnumPtrField))
		if err != nil {
			return nil, err
		}
		result.EnumPtrField = toPtr(tmp58)
	}
	if src.AnyField != nil {
		var tmp59 any
		if err := json.Unmarshal([]byte(*src.AnyField), &tmp59); err != nil {
			return nil, err
		}
		result.AnyField = tmp59
	} else {
		result.AnyField = nil
	}
	if src.AnyPtrField == nil {
		result.AnyPtrField = nil
	} else {
		var tmp61 any
		if err := json.Unmarshal([]byte(fromPtr(src.AnyPtrField)), &tmp61); err != nil {
			return nil, err
		}
		result.AnyPtrField = toPtr(tmp61)
	}
	tmp62 := make(map[string]testservice.SomeOtherModel, len(src.MapModelField))
	for k8, v8 := range src.MapModelField {
		tmp63, err := convertSomeOtherModelFromJsonModel(toPtr(v8))
		if err != nil {
			return nil, err
		}

		tmp62[k8] = fromPtr(tmp63)
	}
	result.MapModelField = tmp62
	tmp64 := make(map[string]*testservice.SomeOtherModel, len(src.MapModelPtrField))
	for k9, v9 := range src.MapModelPtrField {
		if v9 != nil {
			tmp65, err := convertSomeOtherModelFromJsonModel(v9)
			if err != nil {
				return nil, err
			}
			tmp64[k9] = tmp65
		} else {
			tmp64[k9] = nil
		}
	}
	result.MapModelPtrField = tmp64
	tmp66 := make(map[string]testservice.SomeOneOf, len(src.MapOneOfField))
	for k10, v10 := range src.MapOneOfField {
		tmp67, err := convertSomeOneOfFromDB(v10)
		if err != nil {
			return nil, fmt.Errorf("convert SomeOneOf to service type: %w", err)
		}
		tmp66[k10] = tmp67
	}
	result.MapOneOfField = tmp66
	tmp68 := make(map[string]*testservice.SomeOneOf, len(src.MapOneOfPtrField))
	for k11, v11 := range src.MapOneOfPtrField {
		if v11 != nil {
			tmp69, err := convertSomeOneOfFromDB(v11)
			if err != nil {
				return nil, fmt.Errorf("convert SomeOneOf to service type: %w", err)
			}
			tmp68[k11] = toPtr(tmp69)
		} else {
			tmp68[k11] = nil
		}
	}
	result.MapOneOfPtrField = tmp68
	tmp70 := make(map[string]testservice.SomeEnum, len(src.MapEnumField))
	for k12, v12 := range src.MapEnumField {
		tmp71, err := convertSomeEnumFromDB(v12)
		if err != nil {
			return nil, err
		}
		tmp70[k12] = tmp71
	}
	result.MapEnumField = tmp70
	tmp72 := make(map[string]*testservice.SomeEnum, len(src.MapEnumPtrField))
	for k13, v13 := range src.MapEnumPtrField {
		if v13 == nil {
			tmp72[k13] = nil
		} else {
			tmp74, err := convertSomeEnumFromDB(fromPtr(v13))
			if err != nil {
				return nil, err
			}
			tmp72[k13] = toPtr(tmp74)
		}
	}
	result.MapEnumPtrField = tmp72
	tmp75 := make(map[string]any, len(src.MapAnyField))
	for k14, v14 := range src.MapAnyField {
		if v14 != nil {
			var tmp76 any
			if err := json.Unmarshal([]byte(*v14), &tmp76); err != nil {
				return nil, err
			}
			tmp75[k14] = tmp76
		} else {
			tmp75[k14] = nil
		}
	}
	result.MapAnyField = tmp75
	tmp77 := make(map[string]*any, len(src.MapAnyPtrField))
	for k15, v15 := range src.MapAnyPtrField {
		if v15 == nil {
			tmp77[k15] = nil
		} else {
			var tmp79 any
			if err := json.Unmarshal([]byte(fromPtr(v15)), &tmp79); err != nil {
				return nil, err
			}
			tmp77[k15] = toPtr(tmp79)
		}
	}
	result.MapAnyPtrField = tmp77
	tmp80 := make([]testservice.SomeOtherModel, 0, len(src.ModelSliceField))
	for _, el := range src.ModelSliceField {

		tmp81, err := convertSomeOtherModelFromJsonModel(toPtr(el))
		if err != nil {
			return nil, err
		}

		tmp80 = append(tmp80, fromPtr(tmp81))
	}
	result.ModelSliceField = tmp80
	tmp82 := make([]*testservice.SomeOtherModel, 0, len(src.ModelPtrSliceField))
	for _, el := range src.ModelPtrSliceField {

		if el != nil {
			tmp83, err := convertSomeOtherModelFromJsonModel(el)
			if err != nil {
				return nil, err
			}
			tmp82 = append(tmp82, tmp83)
		} else {
			tmp82 = append(tmp82, nil)
		}
	}
	result.ModelPtrSliceField = tmp82
	tmp84 := make([]testservice.SomeOneOf, 0, len(src.OneOfSliceField))
	for _, el := range src.OneOfSliceField {

		tmp85, err := convertSomeOneOfFromDB(el)
		if err != nil {
			return nil, fmt.Errorf("convert SomeOneOf to service type: %w", err)
		}
		tmp84 = append(tmp84, tmp85)
	}
	result.OneOfSliceField = tmp84
	tmp86 := make([]*testservice.SomeOneOf, 0, len(src.OneOfPtrSliceField))
	for _, el := range src.OneOfPtrSliceField {

		if el != nil {
			tmp87, err := convertSomeOneOfFromDB(el)
			if err != nil {
				return nil, fmt.Errorf("convert SomeOneOf to service type: %w", err)
			}
			tmp86 = append(tmp86, toPtr(tmp87))
		} else {
			tmp86 = append(tmp86, nil)
		}
	}
	result.OneOfPtrSliceField = tmp86
	tmp88 := make([]testservice.SomeEnum, 0, len(src.SliceEnumField))
	for _, el := range src.SliceEnumField {

		tmp89, err := convertSomeEnumFromDB(el)
		if err != nil {
			return nil, err
		}
		tmp88 = append(tmp88, tmp89)
	}
	result.SliceEnumField = tmp88
	tmp90 := make([]*testservice.SomeEnum, 0, len(src.SliceEnumPtrField))
	for _, el := range src.SliceEnumPtrField {

		if el == nil {
			tmp90 = append(tmp90, nil)
		} else {
			tmp92, err := convertSomeEnumFromDB(fromPtr(el))
			if err != nil {
				return nil, err
			}
			tmp90 = append(tmp90, toPtr(tmp92))
		}
	}
	result.SliceEnumPtrField = tmp90
	tmp93 := make([]any, 0, len(src.SliceAnyField))
	for _, el := range src.SliceAnyField {

		if el != nil {
			var tmp94 any
			if err := json.Unmarshal([]byte(*el), &tmp94); err != nil {
				return nil, err
			}
			tmp93 = append(tmp93, tmp94)
		} else {
			tmp93 = append(tmp93, nil)
		}
	}
	result.SliceAnyField = tmp93
	tmp95 := make([]*any, 0, len(src.SliceAnyPtrField))
	for _, el := range src.SliceAnyPtrField {

		if el == nil {
			tmp95 = append(tmp95, nil)
		} else {
			var tmp97 any
			if err := json.Unmarshal([]byte(fromPtr(el)), &tmp97); err != nil {
				return nil, err
			}
			tmp95 = append(tmp95, toPtr(tmp97))
		}
	}
	result.SliceAnyPtrField = tmp95
	result.UpdatedAt = src.UpdatedAt
	return result, nil
}

type jsonSomeOtherModel struct {
	ID uuid.UUID `json:"id"`
}

func (s *jsonSomeOtherModel) Scan(value any) error {
	return json.Unmarshal(value.([]byte), s)
}

func (s jsonSomeOtherModel) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func convertSomeOtherModelToJsonModel(src *testservice.SomeOtherModel) (*jsonSomeOtherModel, error) {
	result := &jsonSomeOtherModel{}
	result.ID = src.ID
	return result, nil
}

func convertSomeOtherModelFromJsonModel(src *jsonSomeOtherModel) (*testservice.SomeOtherModel, error) {
	result := &testservice.SomeOtherModel{}
	result.ID = src.ID
	return result, nil
}

type jsonOneOfValue1 struct {
	ValueVal string `json:"value"`
}

func (o *jsonOneOfValue1) Scan(value any) error {
	return json.Unmarshal(value.([]byte), o)
}

func (o jsonOneOfValue1) Value() (driver.Value, error) {
	return json.Marshal(o)
}

func convertOneOfValue1ToJsonModel(src *testservice.OneOfValue1) (*jsonOneOfValue1, error) {
	result := &jsonOneOfValue1{}
	result.ValueVal = src.Value
	return result, nil
}

func convertOneOfValue1FromJsonModel(src *jsonOneOfValue1) (*testservice.OneOfValue1, error) {
	result := &testservice.OneOfValue1{}
	result.Value = src.ValueVal
	return result, nil
}

type jsonOneOfValue2 struct {
	ValueVal string `json:"value"`
}

func (o *jsonOneOfValue2) Scan(value any) error {
	return json.Unmarshal(value.([]byte), o)
}

func (o jsonOneOfValue2) Value() (driver.Value, error) {
	return json.Marshal(o)
}

func convertOneOfValue2ToJsonModel(src *testservice.OneOfValue2) (*jsonOneOfValue2, error) {
	result := &jsonOneOfValue2{}
	result.ValueVal = src.Value
	return result, nil
}

func convertOneOfValue2FromJsonModel(src *jsonOneOfValue2) (*testservice.OneOfValue2, error) {
	result := &testservice.OneOfValue2{}
	result.Value = src.ValueVal
	return result, nil
}

type jsonPasswordRecoveryEvent struct {
	ID             uuid.UUID                      `json:"id"`
	Data           *jsonPasswordRecoveryEventData `json:"data"`
	IdempotencyKey string                         `json:"idempotency_key"`
}

func (p *jsonPasswordRecoveryEvent) Scan(value any) error {
	return json.Unmarshal(value.([]byte), p)
}

func (p jsonPasswordRecoveryEvent) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func convertPasswordRecoveryEventToJsonModel(src *testservice.PasswordRecoveryEvent) (*jsonPasswordRecoveryEvent, error) {
	result := &jsonPasswordRecoveryEvent{}
	result.ID = src.ID
	tmp1, err := convertPasswordRecoveryEventDataToDB(src.Data)
	if err != nil {
		return nil, err
	}
	result.Data = tmp1
	result.IdempotencyKey = src.IdempotencyKey
	return result, nil
}

func convertPasswordRecoveryEventFromJsonModel(src *jsonPasswordRecoveryEvent) (*testservice.PasswordRecoveryEvent, error) {
	result := &testservice.PasswordRecoveryEvent{}
	result.ID = src.ID
	tmp4, err := convertPasswordRecoveryEventDataFromDB(src.Data)
	if err != nil {
		return nil, fmt.Errorf("convert PasswordRecoveryEventData to service type: %w", err)
	}
	result.Data = tmp4
	result.IdempotencyKey = src.IdempotencyKey
	return result, nil
}

type jsonPasswordRecoveryRequestedEventData struct {
	Email            string                         `json:"email"`
	UserID           uuid.UUID                      `json:"user_id"`
	VerificationCode string                         `json:"verification_code"`
	NestedData       *jsonPasswordRecoveryEventData `json:"nested_data"`
}

func (p *jsonPasswordRecoveryRequestedEventData) Scan(value any) error {
	return json.Unmarshal(value.([]byte), p)
}

func (p jsonPasswordRecoveryRequestedEventData) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func convertPasswordRecoveryRequestedEventDataToJsonModel(src *testservice.PasswordRecoveryRequestedEventData) (*jsonPasswordRecoveryRequestedEventData, error) {
	result := &jsonPasswordRecoveryRequestedEventData{}
	result.Email = src.Email
	result.UserID = src.UserID
	result.VerificationCode = src.VerificationCode
	tmp3, err := convertPasswordRecoveryEventDataToDB(src.NestedData)
	if err != nil {
		return nil, err
	}
	result.NestedData = tmp3
	return result, nil
}

func convertPasswordRecoveryRequestedEventDataFromJsonModel(src *jsonPasswordRecoveryRequestedEventData) (*testservice.PasswordRecoveryRequestedEventData, error) {
	result := &testservice.PasswordRecoveryRequestedEventData{}
	result.Email = src.Email
	result.UserID = src.UserID
	result.VerificationCode = src.VerificationCode
	tmp7, err := convertPasswordRecoveryEventDataFromDB(src.NestedData)
	if err != nil {
		return nil, fmt.Errorf("convert PasswordRecoveryEventData to service type: %w", err)
	}
	result.NestedData = tmp7
	return result, nil
}

type jsonPasswordRecoveryCompletedEventData struct {
	Email  string    `json:"email"`
	UserID uuid.UUID `json:"user_id"`
}

func (p *jsonPasswordRecoveryCompletedEventData) Scan(value any) error {
	return json.Unmarshal(value.([]byte), p)
}

func (p jsonPasswordRecoveryCompletedEventData) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func convertPasswordRecoveryCompletedEventDataToJsonModel(src *testservice.PasswordRecoveryCompletedEventData) (*jsonPasswordRecoveryCompletedEventData, error) {
	result := &jsonPasswordRecoveryCompletedEventData{}
	result.Email = src.Email
	result.UserID = src.UserID
	return result, nil
}

func convertPasswordRecoveryCompletedEventDataFromJsonModel(src *jsonPasswordRecoveryCompletedEventData) (*testservice.PasswordRecoveryCompletedEventData, error) {
	result := &testservice.PasswordRecoveryCompletedEventData{}
	result.Email = src.Email
	result.UserID = src.UserID
	return result, nil
}
