package teststorage

import (
	json "encoding/json"
	fmt "fmt"

	uuid "github.com/google/uuid"
	errors "github.com/pkg/errors"

	testservice "github.com/saturn4er/boilerplate-go/test/test/testservice"
	// user code 'imports'
	// end user code 'imports'
)

type dbSomeModel struct {
	ID                 uuid.UUID                             `gorm:"column:id;primaryKey"`
	Name               string                                `gorm:"column:name;type:text;"`
	Description        *string                               `gorm:"column:description;type:text;"`
	ModelField         jsonSomeOtherModel                    `gorm:"column:model_field;"`
	ModelPtrField      *jsonSomeOtherModel                   `gorm:"column:model_ptr_field;"`
	OneOfField         *jsonSomeOneOf                        `gorm:"column:one_of_field;"`
	OneOfPtrField      *jsonSomeOneOf                        `gorm:"column:one_of_ptr_field;"`
	EnumField          string                                `gorm:"column:enum_field;type:text;"`
	EnumPtrField       *string                               `gorm:"column:enum_ptr_field;type:text;"`
	AnyField           *string                               `gorm:"column:any_field;type:text;"`
	AnyPtrField        *string                               `gorm:"column:any_ptr_field;type:text;"`
	MapModelField      mapValue[string, jsonSomeOtherModel]  `gorm:"column:map_model_field;"`
	MapModelPtrField   mapValue[string, *jsonSomeOtherModel] `gorm:"column:map_model_ptr_field;"`
	MapOneOfField      mapValue[string, *jsonSomeOneOf]      `gorm:"column:map_one_of_field;"`
	MapOneOfPtrField   mapValue[string, *jsonSomeOneOf]      `gorm:"column:map_one_of_ptr_field;"`
	MapEnumField       mapValue[string, string]              `gorm:"column:map_enum_field;"`
	MapEnumPtrField    mapValue[string, *string]             `gorm:"column:map_enum_ptr_field;"`
	MapAnyField        mapValue[string, *string]             `gorm:"column:map_any_field;"`
	MapAnyPtrField     mapValue[string, *string]             `gorm:"column:map_any_ptr_field;"`
	ModelSliceField    sliceValue[jsonSomeOtherModel]        `gorm:"column:model_slice_field;"`
	ModelPtrSliceField sliceValue[*jsonSomeOtherModel]       `gorm:"column:model_ptr_slice_field;"`
	OneOfSliceField    sliceValue[*jsonSomeOneOf]            `gorm:"column:one_of_slice_field;"`
	OneOfPtrSliceField sliceValue[*jsonSomeOneOf]            `gorm:"column:one_of_ptr_slice_field;"`
	SliceEnumField     stringSliceValue                      `gorm:"column:slice_enum_field;"`
	SliceEnumPtrField  sliceValue[*string]                   `gorm:"column:slice_enum_ptr_field;"`
	SliceAnyField      sliceValue[*string]                   `gorm:"column:slice_any_field;"`
	SliceAnyPtrField   sliceValue[*string]                   `gorm:"column:slice_any_ptr_field;"`
}

func convertSomeModelToDB(src *testservice.SomeModel) (*dbSomeModel, error) {
	result := &dbSomeModel{}
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
	return result, nil
}

func convertSomeModelFromDB(src *dbSomeModel) (*testservice.SomeModel, error) {
	result := &testservice.SomeModel{}
	result.ID = src.ID
	result.Name = src.Name
	if src.Description == nil {
		result.Description = nil
	} else {
		result.Description = toPtr(fromPtr(src.Description))
	}
	tmp51, err := convertSomeOtherModelFromJsonModel(toPtr(src.ModelField))
	if err != nil {
		return nil, err
	}

	result.ModelField = fromPtr(tmp51)
	if src.ModelPtrField != nil {
		tmp52, err := convertSomeOtherModelFromJsonModel(src.ModelPtrField)
		if err != nil {
			return nil, err
		}
		result.ModelPtrField = tmp52
	} else {
		result.ModelPtrField = nil
	}
	tmp53, err := convertSomeOneOfFromDB(src.OneOfField)
	if err != nil {
		return nil, fmt.Errorf("convert SomeOneOf to service type: %w", err)
	}
	result.OneOfField = tmp53
	if src.OneOfPtrField != nil {
		tmp54, err := convertSomeOneOfFromDB(src.OneOfPtrField)
		if err != nil {
			return nil, fmt.Errorf("convert SomeOneOf to service type: %w", err)
		}
		result.OneOfPtrField = toPtr(tmp54)
	} else {
		result.OneOfPtrField = nil
	}
	tmp55, err := convertSomeEnumFromDB(src.EnumField)
	if err != nil {
		return nil, err
	}
	result.EnumField = tmp55
	if src.EnumPtrField == nil {
		result.EnumPtrField = nil
	} else {
		tmp57, err := convertSomeEnumFromDB(fromPtr(src.EnumPtrField))
		if err != nil {
			return nil, err
		}
		result.EnumPtrField = toPtr(tmp57)
	}
	if src.AnyField != nil {
		var tmp58 any
		if err := json.Unmarshal([]byte(*src.AnyField), &tmp58); err != nil {
			return nil, err
		}
		result.AnyField = tmp58
	} else {
		result.AnyField = nil
	}
	if src.AnyPtrField == nil {
		result.AnyPtrField = nil
	} else {
		var tmp60 any
		if err := json.Unmarshal([]byte(fromPtr(src.AnyPtrField)), &tmp60); err != nil {
			return nil, err
		}
		result.AnyPtrField = toPtr(tmp60)
	}
	tmp61 := make(map[string]testservice.SomeOtherModel, len(src.MapModelField))
	for k8, v8 := range src.MapModelField {
		tmp62, err := convertSomeOtherModelFromJsonModel(toPtr(v8))
		if err != nil {
			return nil, err
		}

		tmp61[k8] = fromPtr(tmp62)
	}
	result.MapModelField = tmp61
	tmp63 := make(map[string]*testservice.SomeOtherModel, len(src.MapModelPtrField))
	for k9, v9 := range src.MapModelPtrField {
		if v9 != nil {
			tmp64, err := convertSomeOtherModelFromJsonModel(v9)
			if err != nil {
				return nil, err
			}
			tmp63[k9] = tmp64
		} else {
			tmp63[k9] = nil
		}
	}
	result.MapModelPtrField = tmp63
	tmp65 := make(map[string]testservice.SomeOneOf, len(src.MapOneOfField))
	for k10, v10 := range src.MapOneOfField {
		tmp66, err := convertSomeOneOfFromDB(v10)
		if err != nil {
			return nil, fmt.Errorf("convert SomeOneOf to service type: %w", err)
		}
		tmp65[k10] = tmp66
	}
	result.MapOneOfField = tmp65
	tmp67 := make(map[string]*testservice.SomeOneOf, len(src.MapOneOfPtrField))
	for k11, v11 := range src.MapOneOfPtrField {
		if v11 != nil {
			tmp68, err := convertSomeOneOfFromDB(v11)
			if err != nil {
				return nil, fmt.Errorf("convert SomeOneOf to service type: %w", err)
			}
			tmp67[k11] = toPtr(tmp68)
		} else {
			tmp67[k11] = nil
		}
	}
	result.MapOneOfPtrField = tmp67
	tmp69 := make(map[string]testservice.SomeEnum, len(src.MapEnumField))
	for k12, v12 := range src.MapEnumField {
		tmp70, err := convertSomeEnumFromDB(v12)
		if err != nil {
			return nil, err
		}
		tmp69[k12] = tmp70
	}
	result.MapEnumField = tmp69
	tmp71 := make(map[string]*testservice.SomeEnum, len(src.MapEnumPtrField))
	for k13, v13 := range src.MapEnumPtrField {
		if v13 == nil {
			tmp71[k13] = nil
		} else {
			tmp73, err := convertSomeEnumFromDB(fromPtr(v13))
			if err != nil {
				return nil, err
			}
			tmp71[k13] = toPtr(tmp73)
		}
	}
	result.MapEnumPtrField = tmp71
	tmp74 := make(map[string]any, len(src.MapAnyField))
	for k14, v14 := range src.MapAnyField {
		if v14 != nil {
			var tmp75 any
			if err := json.Unmarshal([]byte(*v14), &tmp75); err != nil {
				return nil, err
			}
			tmp74[k14] = tmp75
		} else {
			tmp74[k14] = nil
		}
	}
	result.MapAnyField = tmp74
	tmp76 := make(map[string]*any, len(src.MapAnyPtrField))
	for k15, v15 := range src.MapAnyPtrField {
		if v15 == nil {
			tmp76[k15] = nil
		} else {
			var tmp78 any
			if err := json.Unmarshal([]byte(fromPtr(v15)), &tmp78); err != nil {
				return nil, err
			}
			tmp76[k15] = toPtr(tmp78)
		}
	}
	result.MapAnyPtrField = tmp76
	tmp79 := make([]testservice.SomeOtherModel, 0, len(src.ModelSliceField))
	for _, el := range src.ModelSliceField {

		tmp80, err := convertSomeOtherModelFromJsonModel(toPtr(el))
		if err != nil {
			return nil, err
		}

		tmp79 = append(tmp79, fromPtr(tmp80))
	}
	result.ModelSliceField = tmp79
	tmp81 := make([]*testservice.SomeOtherModel, 0, len(src.ModelPtrSliceField))
	for _, el := range src.ModelPtrSliceField {

		if el != nil {
			tmp82, err := convertSomeOtherModelFromJsonModel(el)
			if err != nil {
				return nil, err
			}
			tmp81 = append(tmp81, tmp82)
		} else {
			tmp81 = append(tmp81, nil)
		}
	}
	result.ModelPtrSliceField = tmp81
	tmp83 := make([]testservice.SomeOneOf, 0, len(src.OneOfSliceField))
	for _, el := range src.OneOfSliceField {

		tmp84, err := convertSomeOneOfFromDB(el)
		if err != nil {
			return nil, fmt.Errorf("convert SomeOneOf to service type: %w", err)
		}
		tmp83 = append(tmp83, tmp84)
	}
	result.OneOfSliceField = tmp83
	tmp85 := make([]*testservice.SomeOneOf, 0, len(src.OneOfPtrSliceField))
	for _, el := range src.OneOfPtrSliceField {

		if el != nil {
			tmp86, err := convertSomeOneOfFromDB(el)
			if err != nil {
				return nil, fmt.Errorf("convert SomeOneOf to service type: %w", err)
			}
			tmp85 = append(tmp85, toPtr(tmp86))
		} else {
			tmp85 = append(tmp85, nil)
		}
	}
	result.OneOfPtrSliceField = tmp85
	tmp87 := make([]testservice.SomeEnum, 0, len(src.SliceEnumField))
	for _, el := range src.SliceEnumField {

		tmp88, err := convertSomeEnumFromDB(el)
		if err != nil {
			return nil, err
		}
		tmp87 = append(tmp87, tmp88)
	}
	result.SliceEnumField = tmp87
	tmp89 := make([]*testservice.SomeEnum, 0, len(src.SliceEnumPtrField))
	for _, el := range src.SliceEnumPtrField {

		if el == nil {
			tmp89 = append(tmp89, nil)
		} else {
			tmp91, err := convertSomeEnumFromDB(fromPtr(el))
			if err != nil {
				return nil, err
			}
			tmp89 = append(tmp89, toPtr(tmp91))
		}
	}
	result.SliceEnumPtrField = tmp89
	tmp92 := make([]any, 0, len(src.SliceAnyField))
	for _, el := range src.SliceAnyField {

		if el != nil {
			var tmp93 any
			if err := json.Unmarshal([]byte(*el), &tmp93); err != nil {
				return nil, err
			}
			tmp92 = append(tmp92, tmp93)
		} else {
			tmp92 = append(tmp92, nil)
		}
	}
	result.SliceAnyField = tmp92
	tmp94 := make([]*any, 0, len(src.SliceAnyPtrField))
	for _, el := range src.SliceAnyPtrField {

		if el == nil {
			tmp94 = append(tmp94, nil)
		} else {
			var tmp96 any
			if err := json.Unmarshal([]byte(fromPtr(el)), &tmp96); err != nil {
				return nil, err
			}
			tmp94 = append(tmp94, toPtr(tmp96))
		}
	}
	result.SliceAnyPtrField = tmp94
	return result, nil
}
func (a dbSomeModel) TableName() string {
	return "some_models"
}

type dbSomeOtherModel struct {
	ID uuid.UUID `gorm:"column:id;primaryKey"`
}

func convertSomeOtherModelToDB(src *testservice.SomeOtherModel) (*dbSomeOtherModel, error) {
	result := &dbSomeOtherModel{}
	result.ID = src.ID
	return result, nil
}

func convertSomeOtherModelFromDB(src *dbSomeOtherModel) (*testservice.SomeOtherModel, error) {
	result := &testservice.SomeOtherModel{}
	result.ID = src.ID
	return result, nil
}
func (a dbSomeOtherModel) TableName() string {
	return "some_other_models"
}

type dbPasswordRecoveryEvent struct {
	ID             uuid.UUID                      `gorm:"column:id;primaryKey"`
	Data           *jsonPasswordRecoveryEventData `gorm:"column:data;"`
	IdempotencyKey string                         `gorm:"column:idempotency_key;type:text;"`
}

func convertPasswordRecoveryEventToDB(src *testservice.PasswordRecoveryEvent) (*dbPasswordRecoveryEvent, error) {
	result := &dbPasswordRecoveryEvent{}
	result.ID = src.ID
	tmp1, err := convertPasswordRecoveryEventDataToDB(src.Data)
	if err != nil {
		return nil, err
	}
	result.Data = tmp1
	result.IdempotencyKey = src.IdempotencyKey
	return result, nil
}

func convertPasswordRecoveryEventFromDB(src *dbPasswordRecoveryEvent) (*testservice.PasswordRecoveryEvent, error) {
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
