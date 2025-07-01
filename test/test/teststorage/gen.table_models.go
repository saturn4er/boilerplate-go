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
	tmp1, err := convertSomeOtherModelToJsonModel(toPtr(src.ModelField))
	if err != nil {
		return nil, errors.Wrap(err, "convert SomeOtherModel to db")
	}
	result.ModelField = *tmp1
	if src.ModelPtrField != nil {
		tmp2, err := convertSomeOtherModelToJsonModel(src.ModelPtrField)
		if err != nil {
			return nil, errors.Wrap(err, "convert SomeOtherModel to db")
		}
		result.ModelPtrField = tmp2
	} else {
		result.ModelPtrField = nil
	}
	tmp3, err := convertSomeOneOfToDB(src.OneOfField)
	if err != nil {
		return nil, err
	}
	result.OneOfField = tmp3
	if src.OneOfPtrField != nil {
		tmp4, err := convertSomeOneOfToDB(*src.OneOfPtrField)
		if err != nil {
			return nil, err
		}
		result.OneOfPtrField = tmp4
	} else {
		result.OneOfPtrField = nil
	}
	tmp5, err := convertSomeEnumToDB(src.EnumField)
	if err != nil {
		return nil, err
	}
	result.EnumField = tmp5
	if src.EnumPtrField == nil {
		result.EnumPtrField = nil
	} else {
		tmp7, err := convertSomeEnumToDB(fromPtr(src.EnumPtrField))
		if err != nil {
			return nil, err
		}
		result.EnumPtrField = toPtr(tmp7)
	}
	if src.AnyField != nil {
		tmp8, err := json.Marshal(src.AnyField)
		if err != nil {
			return nil, err
		}

		marshaledValue := string(tmp8)
		result.AnyField = toPtr(marshaledValue)
	} else {
		result.AnyField = nil
	}
	if src.AnyPtrField != nil && fromPtr(src.AnyPtrField) != nil {
		tmp9, err := json.Marshal(*src.AnyPtrField)
		if err != nil {
			return nil, err
		}

		marshaledValue1 := string(tmp9)
		result.AnyPtrField = toPtr(marshaledValue1)
	} else {
		result.AnyPtrField = nil
	}
	tmp10 := make(mapValue[string, jsonSomeOtherModel], len(src.MapModelField))
	for k, v := range src.MapModelField {
		tmp11, err := convertSomeOtherModelToJsonModel(toPtr(v))
		if err != nil {
			return nil, errors.Wrap(err, "convert SomeOtherModel to db")
		}
		tmp10[k] = *tmp11
	}
	result.MapModelField = tmp10
	tmp12 := make(mapValue[string, *jsonSomeOtherModel], len(src.MapModelPtrField))
	for k1, v1 := range src.MapModelPtrField {
		if v1 != nil {
			tmp13, err := convertSomeOtherModelToJsonModel(v1)
			if err != nil {
				return nil, errors.Wrap(err, "convert SomeOtherModel to db")
			}
			tmp12[k1] = tmp13
		} else {
			tmp12[k1] = nil
		}
	}
	result.MapModelPtrField = tmp12
	tmp14 := make(mapValue[string, *jsonSomeOneOf], len(src.MapOneOfField))
	for k2, v2 := range src.MapOneOfField {
		tmp15, err := convertSomeOneOfToDB(v2)
		if err != nil {
			return nil, err
		}
		tmp14[k2] = tmp15
	}
	result.MapOneOfField = tmp14
	tmp16 := make(mapValue[string, *jsonSomeOneOf], len(src.MapOneOfPtrField))
	for k3, v3 := range src.MapOneOfPtrField {
		if v3 != nil {
			tmp17, err := convertSomeOneOfToDB(*v3)
			if err != nil {
				return nil, err
			}
			tmp16[k3] = tmp17
		} else {
			tmp16[k3] = nil
		}
	}
	result.MapOneOfPtrField = tmp16
	tmp18 := make(mapValue[string, string], len(src.MapEnumField))
	for k4, v4 := range src.MapEnumField {
		tmp19, err := convertSomeEnumToDB(v4)
		if err != nil {
			return nil, err
		}
		tmp18[k4] = tmp19
	}
	result.MapEnumField = tmp18
	tmp20 := make(mapValue[string, *string], len(src.MapEnumPtrField))
	for k5, v5 := range src.MapEnumPtrField {
		if v5 == nil {
			tmp20[k5] = nil
		} else {
			tmp22, err := convertSomeEnumToDB(fromPtr(v5))
			if err != nil {
				return nil, err
			}
			tmp20[k5] = toPtr(tmp22)
		}
	}
	result.MapEnumPtrField = tmp20
	tmp23 := make(mapValue[string, *string], len(src.MapAnyField))
	for k6, v6 := range src.MapAnyField {
		if v6 != nil {
			tmp24, err := json.Marshal(v6)
			if err != nil {
				return nil, err
			}

			marshaledValue2 := string(tmp24)
			tmp23[k6] = toPtr(marshaledValue2)
		} else {
			tmp23[k6] = nil
		}
	}
	result.MapAnyField = tmp23
	tmp25 := make(mapValue[string, *string], len(src.MapAnyPtrField))
	for k7, v7 := range src.MapAnyPtrField {
		if v7 != nil && fromPtr(v7) != nil {
			tmp26, err := json.Marshal(*v7)
			if err != nil {
				return nil, err
			}

			marshaledValue3 := string(tmp26)
			tmp25[k7] = toPtr(marshaledValue3)
		} else {
			tmp25[k7] = nil
		}
	}
	result.MapAnyPtrField = tmp25
	tmp27 := make(sliceValue[jsonSomeOtherModel], 0, len(src.ModelSliceField))
	for _, el := range src.ModelSliceField {
		tmp28, err := convertSomeOtherModelToJsonModel(toPtr(el))
		if err != nil {
			return nil, errors.Wrap(err, "convert SomeOtherModel to db")
		}
		tmp27 = append(tmp27, *tmp28)
	}
	result.ModelSliceField = tmp27
	tmp29 := make(sliceValue[*jsonSomeOtherModel], 0, len(src.ModelPtrSliceField))
	for _, el := range src.ModelPtrSliceField {
		if el != nil {
			tmp30, err := convertSomeOtherModelToJsonModel(el)
			if err != nil {
				return nil, errors.Wrap(err, "convert SomeOtherModel to db")
			}
			tmp29 = append(tmp29, tmp30)
		} else {
			tmp29 = append(tmp29, nil)
		}
	}
	result.ModelPtrSliceField = tmp29
	tmp31 := make(sliceValue[*jsonSomeOneOf], 0, len(src.OneOfSliceField))
	for _, el := range src.OneOfSliceField {
		tmp32, err := convertSomeOneOfToDB(el)
		if err != nil {
			return nil, err
		}
		tmp31 = append(tmp31, tmp32)
	}
	result.OneOfSliceField = tmp31
	tmp33 := make(sliceValue[*jsonSomeOneOf], 0, len(src.OneOfPtrSliceField))
	for _, el := range src.OneOfPtrSliceField {
		if el != nil {
			tmp34, err := convertSomeOneOfToDB(*el)
			if err != nil {
				return nil, err
			}
			tmp33 = append(tmp33, tmp34)
		} else {
			tmp33 = append(tmp33, nil)
		}
	}
	result.OneOfPtrSliceField = tmp33
	tmp35 := make(stringSliceValue, 0, len(src.SliceEnumField))
	for _, el := range src.SliceEnumField {
		tmp36, err := convertSomeEnumToDB(el)
		if err != nil {
			return nil, err
		}
		tmp35 = append(tmp35, tmp36)
	}
	result.SliceEnumField = tmp35
	tmp37 := make(sliceValue[*string], 0, len(src.SliceEnumPtrField))
	for _, el := range src.SliceEnumPtrField {
		if el == nil {
			tmp37 = append(tmp37, nil)
		} else {
			tmp39, err := convertSomeEnumToDB(fromPtr(el))
			if err != nil {
				return nil, err
			}
			tmp37 = append(tmp37, toPtr(tmp39))
		}
	}
	result.SliceEnumPtrField = tmp37
	tmp40 := make(sliceValue[*string], 0, len(src.SliceAnyField))
	for _, el := range src.SliceAnyField {
		if el != nil {
			tmp41, err := json.Marshal(el)
			if err != nil {
				return nil, err
			}

			marshaledValue4 := string(tmp41)
			tmp40 = append(tmp40, toPtr(marshaledValue4))
		} else {
			tmp40 = append(tmp40, nil)
		}
	}
	result.SliceAnyField = tmp40
	tmp42 := make(sliceValue[*string], 0, len(src.SliceAnyPtrField))
	for _, el := range src.SliceAnyPtrField {
		if el != nil && fromPtr(el) != nil {
			tmp43, err := json.Marshal(*el)
			if err != nil {
				return nil, err
			}

			marshaledValue5 := string(tmp43)
			tmp42 = append(tmp42, toPtr(marshaledValue5))
		} else {
			tmp42 = append(tmp42, nil)
		}
	}
	result.SliceAnyPtrField = tmp42
	return result, nil
}

func convertSomeModelFromDB(src *dbSomeModel) (*testservice.SomeModel, error) {
	result := &testservice.SomeModel{}
	result.ID = src.ID
	tmp45, err := convertSomeOtherModelFromJsonModel(toPtr(src.ModelField))
	if err != nil {
		return nil, err
	}

	result.ModelField = fromPtr(tmp45)
	if src.ModelPtrField != nil {
		tmp46, err := convertSomeOtherModelFromJsonModel(src.ModelPtrField)
		if err != nil {
			return nil, err
		}
		result.ModelPtrField = tmp46
	} else {
		result.ModelPtrField = nil
	}
	tmp47, err := convertSomeOneOfFromDB(src.OneOfField)
	if err != nil {
		return nil, fmt.Errorf("convert SomeOneOf to service type: %w", err)
	}
	result.OneOfField = tmp47
	if src.OneOfPtrField != nil {
		tmp48, err := convertSomeOneOfFromDB(src.OneOfPtrField)
		if err != nil {
			return nil, fmt.Errorf("convert SomeOneOf to service type: %w", err)
		}
		result.OneOfPtrField = toPtr(tmp48)
	} else {
		result.OneOfPtrField = nil
	}
	tmp49, err := convertSomeEnumFromDB(src.EnumField)
	if err != nil {
		return nil, err
	}
	result.EnumField = tmp49
	if src.EnumPtrField == nil {
		result.EnumPtrField = nil
	} else {
		tmp51, err := convertSomeEnumFromDB(fromPtr(src.EnumPtrField))
		if err != nil {
			return nil, err
		}
		result.EnumPtrField = toPtr(tmp51)
	}
	if src.AnyField != nil {
		var tmp52 any
		if err := json.Unmarshal([]byte(*src.AnyField), &tmp52); err != nil {
			return nil, err
		}
		result.AnyField = tmp52
	} else {
		result.AnyField = nil
	}
	if src.AnyPtrField == nil {
		result.AnyPtrField = nil
	} else {
		var tmp54 any
		if err := json.Unmarshal([]byte(fromPtr(src.AnyPtrField)), &tmp54); err != nil {
			return nil, err
		}
		result.AnyPtrField = toPtr(tmp54)
	}
	tmp55 := make(map[string]testservice.SomeOtherModel, len(src.MapModelField))
	for k8, v8 := range src.MapModelField {
		tmp56, err := convertSomeOtherModelFromJsonModel(toPtr(v8))
		if err != nil {
			return nil, err
		}

		tmp55[k8] = fromPtr(tmp56)
	}
	result.MapModelField = tmp55
	tmp57 := make(map[string]*testservice.SomeOtherModel, len(src.MapModelPtrField))
	for k9, v9 := range src.MapModelPtrField {
		if v9 != nil {
			tmp58, err := convertSomeOtherModelFromJsonModel(v9)
			if err != nil {
				return nil, err
			}
			tmp57[k9] = tmp58
		} else {
			tmp57[k9] = nil
		}
	}
	result.MapModelPtrField = tmp57
	tmp59 := make(map[string]testservice.SomeOneOf, len(src.MapOneOfField))
	for k10, v10 := range src.MapOneOfField {
		tmp60, err := convertSomeOneOfFromDB(v10)
		if err != nil {
			return nil, fmt.Errorf("convert SomeOneOf to service type: %w", err)
		}
		tmp59[k10] = tmp60
	}
	result.MapOneOfField = tmp59
	tmp61 := make(map[string]*testservice.SomeOneOf, len(src.MapOneOfPtrField))
	for k11, v11 := range src.MapOneOfPtrField {
		if v11 != nil {
			tmp62, err := convertSomeOneOfFromDB(v11)
			if err != nil {
				return nil, fmt.Errorf("convert SomeOneOf to service type: %w", err)
			}
			tmp61[k11] = toPtr(tmp62)
		} else {
			tmp61[k11] = nil
		}
	}
	result.MapOneOfPtrField = tmp61
	tmp63 := make(map[string]testservice.SomeEnum, len(src.MapEnumField))
	for k12, v12 := range src.MapEnumField {
		tmp64, err := convertSomeEnumFromDB(v12)
		if err != nil {
			return nil, err
		}
		tmp63[k12] = tmp64
	}
	result.MapEnumField = tmp63
	tmp65 := make(map[string]*testservice.SomeEnum, len(src.MapEnumPtrField))
	for k13, v13 := range src.MapEnumPtrField {
		if v13 == nil {
			tmp65[k13] = nil
		} else {
			tmp67, err := convertSomeEnumFromDB(fromPtr(v13))
			if err != nil {
				return nil, err
			}
			tmp65[k13] = toPtr(tmp67)
		}
	}
	result.MapEnumPtrField = tmp65
	tmp68 := make(map[string]any, len(src.MapAnyField))
	for k14, v14 := range src.MapAnyField {
		if v14 != nil {
			var tmp69 any
			if err := json.Unmarshal([]byte(*v14), &tmp69); err != nil {
				return nil, err
			}
			tmp68[k14] = tmp69
		} else {
			tmp68[k14] = nil
		}
	}
	result.MapAnyField = tmp68
	tmp70 := make(map[string]*any, len(src.MapAnyPtrField))
	for k15, v15 := range src.MapAnyPtrField {
		if v15 == nil {
			tmp70[k15] = nil
		} else {
			var tmp72 any
			if err := json.Unmarshal([]byte(fromPtr(v15)), &tmp72); err != nil {
				return nil, err
			}
			tmp70[k15] = toPtr(tmp72)
		}
	}
	result.MapAnyPtrField = tmp70
	tmp73 := make([]testservice.SomeOtherModel, 0, len(src.ModelSliceField))
	for _, el := range src.ModelSliceField {

		tmp74, err := convertSomeOtherModelFromJsonModel(toPtr(el))
		if err != nil {
			return nil, err
		}

		tmp73 = append(tmp73, fromPtr(tmp74))
	}
	result.ModelSliceField = tmp73
	tmp75 := make([]*testservice.SomeOtherModel, 0, len(src.ModelPtrSliceField))
	for _, el := range src.ModelPtrSliceField {

		if el != nil {
			tmp76, err := convertSomeOtherModelFromJsonModel(el)
			if err != nil {
				return nil, err
			}
			tmp75 = append(tmp75, tmp76)
		} else {
			tmp75 = append(tmp75, nil)
		}
	}
	result.ModelPtrSliceField = tmp75
	tmp77 := make([]testservice.SomeOneOf, 0, len(src.OneOfSliceField))
	for _, el := range src.OneOfSliceField {

		tmp78, err := convertSomeOneOfFromDB(el)
		if err != nil {
			return nil, fmt.Errorf("convert SomeOneOf to service type: %w", err)
		}
		tmp77 = append(tmp77, tmp78)
	}
	result.OneOfSliceField = tmp77
	tmp79 := make([]*testservice.SomeOneOf, 0, len(src.OneOfPtrSliceField))
	for _, el := range src.OneOfPtrSliceField {

		if el != nil {
			tmp80, err := convertSomeOneOfFromDB(el)
			if err != nil {
				return nil, fmt.Errorf("convert SomeOneOf to service type: %w", err)
			}
			tmp79 = append(tmp79, toPtr(tmp80))
		} else {
			tmp79 = append(tmp79, nil)
		}
	}
	result.OneOfPtrSliceField = tmp79
	tmp81 := make([]testservice.SomeEnum, 0, len(src.SliceEnumField))
	for _, el := range src.SliceEnumField {

		tmp82, err := convertSomeEnumFromDB(el)
		if err != nil {
			return nil, err
		}
		tmp81 = append(tmp81, tmp82)
	}
	result.SliceEnumField = tmp81
	tmp83 := make([]*testservice.SomeEnum, 0, len(src.SliceEnumPtrField))
	for _, el := range src.SliceEnumPtrField {

		if el == nil {
			tmp83 = append(tmp83, nil)
		} else {
			tmp85, err := convertSomeEnumFromDB(fromPtr(el))
			if err != nil {
				return nil, err
			}
			tmp83 = append(tmp83, toPtr(tmp85))
		}
	}
	result.SliceEnumPtrField = tmp83
	tmp86 := make([]any, 0, len(src.SliceAnyField))
	for _, el := range src.SliceAnyField {

		if el != nil {
			var tmp87 any
			if err := json.Unmarshal([]byte(*el), &tmp87); err != nil {
				return nil, err
			}
			tmp86 = append(tmp86, tmp87)
		} else {
			tmp86 = append(tmp86, nil)
		}
	}
	result.SliceAnyField = tmp86
	tmp88 := make([]*any, 0, len(src.SliceAnyPtrField))
	for _, el := range src.SliceAnyPtrField {

		if el == nil {
			tmp88 = append(tmp88, nil)
		} else {
			var tmp90 any
			if err := json.Unmarshal([]byte(fromPtr(el)), &tmp90); err != nil {
				return nil, err
			}
			tmp88 = append(tmp88, toPtr(tmp90))
		}
	}
	result.SliceAnyPtrField = tmp88
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
