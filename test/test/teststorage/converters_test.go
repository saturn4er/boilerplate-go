package teststorage

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/saturn4er/boilerplate-go/test/test/testservice"
)

func TestConvertSomeModelToJsonModel(t *testing.T) {
	model := &testservice.SomeModel{
		ID: uuid.MustParse("ebb58ec0-931d-4318-977d-2179e8141c54"),
		ModelField: testservice.SomeOtherModel{
			ID: uuid.MustParse("b230154b-dfbb-484e-af4d-58edf690638d"),
		},
		ModelPtrField: &testservice.SomeOtherModel{
			ID: uuid.MustParse("0d83279d-502c-4d16-a7ef-9d0a8bb2ff61"),
		},
		OneOfField:    &testservice.OneOfValue1{Value: "asd3"},
		OneOfPtrField: toPtr[testservice.SomeOneOf](&testservice.OneOfValue2{Value: "qwe"}),
		EnumField:     testservice.SomeEnumValue1,
		EnumPtrField:  toPtr(testservice.SomeEnumValue2),
		AnyField:      map[string]any{"key": "value"},
		AnyPtrField:   toPtr[any]([]any{"value"}),
		MapModelField: map[string]testservice.SomeOtherModel{"key": {ID: uuid.MustParse("0a313dcf-dbe3-4da6-b37f-d81456fcda7a")}},
		MapModelPtrField: map[string]*testservice.SomeOtherModel{
			"some":    {ID: uuid.MustParse("f587dfba-0247-4678-9902-c22d9e58c99b")},
			"values":  {ID: uuid.MustParse("8ca81f7e-407e-4c25-bb4c-d95c0fe38314")},
			"and_nil": nil,
		},
		MapOneOfField: map[string]testservice.SomeOneOf{
			"key":       &testservice.OneOfValue1{Value: "123"},
			"and_other": &testservice.OneOfValue2{Value: "sad"},
			"and_nil":   nil,
		},
		MapOneOfPtrField: map[string]*testservice.SomeOneOf{
			"key": toPtr[testservice.SomeOneOf](&testservice.OneOfValue1{Value: "123"}),
			"and": toPtr[testservice.SomeOneOf](&testservice.OneOfValue2{Value: "dgsd"}),
			"nil": nil,
		},
		MapEnumField: map[string]testservice.SomeEnum{
			"val1": testservice.SomeEnumValue1,
			"val2": testservice.SomeEnumValue1,
		},
		MapEnumPtrField: map[string]*testservice.SomeEnum{
			"val1": toPtr(testservice.SomeEnumValue1),
			"val2": toPtr(testservice.SomeEnumValue2),
			"val3": nil,
		},
		MapAnyField:    map[string]any{"key": "value"},
		MapAnyPtrField: map[string]*any{"key": toPtr[any]("value")},
		ModelSliceField: []testservice.SomeOtherModel{
			{ID: uuid.MustParse("b230154b-dfbb-484e-af4d-58edf690638d")},
			{ID: uuid.MustParse("ab27cfcb-7fc1-4ec6-9bab-1088eb024aee")},
		},
		ModelPtrSliceField: []*testservice.SomeOtherModel{
			{ID: uuid.MustParse("b230154b-dfbb-484e-af4d-58edf690638d")},
			nil,
			{ID: uuid.MustParse("b4af1c7e-cfca-4cbd-80c5-c9232fb39a7c")},
		},
		OneOfSliceField:    []testservice.SomeOneOf{&testservice.OneOfValue1{Value: "123"}, &testservice.OneOfValue2{Value: "456"}},
		OneOfPtrSliceField: []*testservice.SomeOneOf{toPtr[testservice.SomeOneOf](&testservice.OneOfValue1{Value: "123"}), nil, toPtr[testservice.SomeOneOf](&testservice.OneOfValue2{Value: "zxc"})},
		SliceEnumField:     []testservice.SomeEnum{testservice.SomeEnumValue1, testservice.SomeEnumValue2},
		SliceEnumPtrField:  []*testservice.SomeEnum{toPtr(testservice.SomeEnumValue1), nil, toPtr(testservice.SomeEnumValue2)},
		SliceAnyField:      []any{"value", "hello", nil},
		SliceAnyPtrField:   []*any{toPtr[any]("value"), toPtr[any]("hello"), nil},
	}
	dbModel, err := convertSomeModelToDB(model)
	require.NoError(t, err)
	require.Equal(t, &dbSomeModel{
		ID: uuid.MustParse("ebb58ec0-931d-4318-977d-2179e8141c54"),
		ModelField: jsonSomeOtherModel{
			ID: uuid.MustParse("b230154b-dfbb-484e-af4d-58edf690638d"),
		},
		ModelPtrField: &jsonSomeOtherModel{
			ID: uuid.MustParse("0d83279d-502c-4d16-a7ef-9d0a8bb2ff61"),
		},
		OneOfField: &jsonSomeOneOf{
			OneOfType:   "OneOfValue1",
			OneOfTypeID: 1,
			Val: &jsonOneOfValue1{
				ValueVal: "asd3",
			},
		},
		OneOfPtrField: &jsonSomeOneOf{
			OneOfType:   "OneOfValue2",
			OneOfTypeID: 2,
			Val: &jsonOneOfValue2{
				ValueVal: "qwe",
			},
		},
		EnumField:     someEnumValue1,
		EnumPtrField:  toPtr(someEnumValue2),
		AnyField:      toPtr("{\"key\":\"value\"}"),
		AnyPtrField:   toPtr("[\"value\"]"),
		MapModelField: map[string]jsonSomeOtherModel{"key": {ID: uuid.MustParse("0a313dcf-dbe3-4da6-b37f-d81456fcda7a")}},
		MapModelPtrField: map[string]*jsonSomeOtherModel{
			"some":    {ID: uuid.MustParse("f587dfba-0247-4678-9902-c22d9e58c99b")},
			"values":  {ID: uuid.MustParse("8ca81f7e-407e-4c25-bb4c-d95c0fe38314")},
			"and_nil": nil,
		},
		MapOneOfField: map[string]*jsonSomeOneOf{
			"key":       {OneOfType: "OneOfValue1", OneOfTypeID: 1, Val: &jsonOneOfValue1{ValueVal: "123"}},
			"and_other": {OneOfType: "OneOfValue2", OneOfTypeID: 2, Val: &jsonOneOfValue2{ValueVal: "sad"}},
			"and_nil":   nil,
		},
		MapOneOfPtrField: map[string]*jsonSomeOneOf{
			"key": {
				OneOfType:   "OneOfValue1",
				OneOfTypeID: 1,
				Val: &jsonOneOfValue1{
					ValueVal: "123",
				},
			},
			"and": {
				OneOfType:   "OneOfValue2",
				OneOfTypeID: 2,
				Val: &jsonOneOfValue2{
					ValueVal: "dgsd",
				},
			},
			"nil": nil,
		},
		MapEnumField: map[string]string{
			"val1": someEnumValue1,
			"val2": someEnumValue1,
		},
		MapEnumPtrField: map[string]*string{
			"val1": toPtr(someEnumValue1),
			"val2": toPtr(someEnumValue2),
			"val3": nil,
		},
		MapAnyField:    mapValue[string, *string]{"key": toPtr("\"value\"")},
		MapAnyPtrField: mapValue[string, *string]{"key": toPtr("\"value\"")},
		ModelSliceField: []jsonSomeOtherModel{
			{ID: uuid.MustParse("b230154b-dfbb-484e-af4d-58edf690638d")},
			{ID: uuid.MustParse("ab27cfcb-7fc1-4ec6-9bab-1088eb024aee")},
		},
		ModelPtrSliceField: []*jsonSomeOtherModel{
			{ID: uuid.MustParse("b230154b-dfbb-484e-af4d-58edf690638d")},
			nil,
			{ID: uuid.MustParse("b4af1c7e-cfca-4cbd-80c5-c9232fb39a7c")},
		},
		OneOfSliceField: sliceValue[*jsonSomeOneOf]{&jsonSomeOneOf{
			OneOfType:   "OneOfValue1",
			OneOfTypeID: 1,
			Val: &jsonOneOfValue1{
				ValueVal: "123",
			},
		}, &jsonSomeOneOf{OneOfType: "OneOfValue2", OneOfTypeID: 2, Val: &jsonOneOfValue2{ValueVal: "456"}}},
		OneOfPtrSliceField: []*jsonSomeOneOf{
			{OneOfType: "OneOfValue1", OneOfTypeID: 1, Val: &jsonOneOfValue1{ValueVal: "123"}},
			nil,
			{OneOfType: "OneOfValue2", OneOfTypeID: 2, Val: &jsonOneOfValue2{ValueVal: "zxc"}},
		},
		SliceEnumField:    []string{someEnumValue1, someEnumValue2},
		SliceEnumPtrField: []*string{toPtr(someEnumValue1), nil, toPtr(someEnumValue2)},
		SliceAnyField:     sliceValue[*string]{toPtr("\"value\""), toPtr("\"hello\""), nil},
		SliceAnyPtrField:  []*string{toPtr("\"value\""), toPtr("\"hello\""), nil},
	}, dbModel)
}

func TestConvertSomeModelFromJsonModel(t *testing.T) {
	serviceModel, err := convertSomeModelFromDB(&dbSomeModel{
		ID: uuid.MustParse("ebb58ec0-931d-4318-977d-2179e8141c54"),
		ModelField: jsonSomeOtherModel{
			ID: uuid.MustParse("b230154b-dfbb-484e-af4d-58edf690638d"),
		},
		ModelPtrField: &jsonSomeOtherModel{
			ID: uuid.MustParse("0d83279d-502c-4d16-a7ef-9d0a8bb2ff61"),
		},
		OneOfField: &jsonSomeOneOf{
			OneOfType:   "OneOfValue1",
			OneOfTypeID: 1,
			Val: &jsonOneOfValue1{
				ValueVal: "asd3",
			},
		},
		OneOfPtrField: &jsonSomeOneOf{
			OneOfType:   "OneOfValue2",
			OneOfTypeID: 2,
			Val: &jsonOneOfValue2{
				ValueVal: "qwe",
			},
		},
		EnumField:     someEnumValue1,
		EnumPtrField:  toPtr(someEnumValue2),
		AnyField:      toPtr("{\"key\":\"value\"}"),
		AnyPtrField:   toPtr("[\"value\"]"),
		MapModelField: map[string]jsonSomeOtherModel{"key": {ID: uuid.MustParse("0a313dcf-dbe3-4da6-b37f-d81456fcda7a")}},
		MapModelPtrField: map[string]*jsonSomeOtherModel{
			"some":    {ID: uuid.MustParse("f587dfba-0247-4678-9902-c22d9e58c99b")},
			"values":  {ID: uuid.MustParse("8ca81f7e-407e-4c25-bb4c-d95c0fe38314")},
			"and_nil": nil,
		},
		MapOneOfField: map[string]*jsonSomeOneOf{
			"key":       {OneOfType: "OneOfValue1", OneOfTypeID: 1, Val: &jsonOneOfValue1{ValueVal: "123"}},
			"and_other": {OneOfType: "OneOfValue2", OneOfTypeID: 2, Val: &jsonOneOfValue2{ValueVal: "sad"}},
			"and_nil":   nil,
		},
		MapOneOfPtrField: map[string]*jsonSomeOneOf{
			"key": {
				OneOfType:   "OneOfValue1",
				OneOfTypeID: 1,
				Val: &jsonOneOfValue1{
					ValueVal: "123",
				},
			},
			"and": {
				OneOfType:   "OneOfValue2",
				OneOfTypeID: 2,
				Val: &jsonOneOfValue2{
					ValueVal: "dgsd",
				},
			},
			"nil": nil,
		},
		MapEnumField: map[string]string{
			"val1": someEnumValue1,
			"val2": someEnumValue1,
		},
		MapEnumPtrField: map[string]*string{
			"val1": toPtr(someEnumValue1),
			"val2": toPtr(someEnumValue2),
			"val3": nil,
		},
		MapAnyField:    mapValue[string, *string]{"key": toPtr("\"value\"")},
		MapAnyPtrField: mapValue[string, *string]{"key": toPtr("\"value\"")},
		ModelSliceField: []jsonSomeOtherModel{
			{ID: uuid.MustParse("b230154b-dfbb-484e-af4d-58edf690638d")},
			{ID: uuid.MustParse("ab27cfcb-7fc1-4ec6-9bab-1088eb024aee")},
		},
		ModelPtrSliceField: []*jsonSomeOtherModel{
			{ID: uuid.MustParse("b230154b-dfbb-484e-af4d-58edf690638d")},
			nil,
			{ID: uuid.MustParse("b4af1c7e-cfca-4cbd-80c5-c9232fb39a7c")},
		},
		OneOfSliceField: sliceValue[*jsonSomeOneOf]{&jsonSomeOneOf{
			OneOfType:   "OneOfValue1",
			OneOfTypeID: 1,
			Val: &jsonOneOfValue1{
				ValueVal: "123",
			},
		}, &jsonSomeOneOf{OneOfType: "OneOfValue2", OneOfTypeID: 2, Val: &jsonOneOfValue2{ValueVal: "456"}}},
		OneOfPtrSliceField: []*jsonSomeOneOf{
			{OneOfType: "OneOfValue1", OneOfTypeID: 1, Val: &jsonOneOfValue1{ValueVal: "123"}},
			nil,
			{OneOfType: "OneOfValue2", OneOfTypeID: 2, Val: &jsonOneOfValue2{ValueVal: "zxc"}},
		},
		SliceEnumField:    []string{someEnumValue1, someEnumValue2},
		SliceEnumPtrField: []*string{toPtr(someEnumValue1), nil, toPtr(someEnumValue2)},
		SliceAnyField:     sliceValue[*string]{toPtr("\"value\""), toPtr("\"hello\""), nil},
		SliceAnyPtrField:  []*string{toPtr("\"value\""), toPtr("\"hello\""), nil},
	})
	require.NoError(t, err)

	require.Equal(t, &testservice.SomeModel{
		ID: uuid.MustParse("ebb58ec0-931d-4318-977d-2179e8141c54"),
		ModelField: testservice.SomeOtherModel{
			ID: uuid.MustParse("b230154b-dfbb-484e-af4d-58edf690638d"),
		},
		ModelPtrField: &testservice.SomeOtherModel{
			ID: uuid.MustParse("0d83279d-502c-4d16-a7ef-9d0a8bb2ff61"),
		},
		OneOfField:    &testservice.OneOfValue1{Value: "asd3"},
		OneOfPtrField: toPtr[testservice.SomeOneOf](&testservice.OneOfValue2{Value: "qwe"}),
		EnumField:     testservice.SomeEnumValue1,
		EnumPtrField:  toPtr(testservice.SomeEnumValue2),
		AnyField:      map[string]any{"key": "value"},
		AnyPtrField:   toPtr[any]([]any{"value"}),
		MapModelField: map[string]testservice.SomeOtherModel{"key": {ID: uuid.MustParse("0a313dcf-dbe3-4da6-b37f-d81456fcda7a")}},
		MapModelPtrField: map[string]*testservice.SomeOtherModel{
			"some":    {ID: uuid.MustParse("f587dfba-0247-4678-9902-c22d9e58c99b")},
			"values":  {ID: uuid.MustParse("8ca81f7e-407e-4c25-bb4c-d95c0fe38314")},
			"and_nil": nil,
		},
		MapOneOfField: map[string]testservice.SomeOneOf{
			"key":       &testservice.OneOfValue1{Value: "123"},
			"and_other": &testservice.OneOfValue2{Value: "sad"},
			"and_nil":   nil,
		},
		MapOneOfPtrField: map[string]*testservice.SomeOneOf{
			"key": toPtr[testservice.SomeOneOf](&testservice.OneOfValue1{Value: "123"}),
			"and": toPtr[testservice.SomeOneOf](&testservice.OneOfValue2{Value: "dgsd"}),
			"nil": nil,
		},
		MapEnumField: map[string]testservice.SomeEnum{
			"val1": testservice.SomeEnumValue1,
			"val2": testservice.SomeEnumValue1,
		},
		MapEnumPtrField: map[string]*testservice.SomeEnum{
			"val1": toPtr(testservice.SomeEnumValue1),
			"val2": toPtr(testservice.SomeEnumValue2),
			"val3": nil,
		},
		MapAnyField:    map[string]any{"key": "value"},
		MapAnyPtrField: map[string]*any{"key": toPtr[any]("value")},
		ModelSliceField: []testservice.SomeOtherModel{
			{ID: uuid.MustParse("b230154b-dfbb-484e-af4d-58edf690638d")},
			{ID: uuid.MustParse("ab27cfcb-7fc1-4ec6-9bab-1088eb024aee")},
		},
		ModelPtrSliceField: []*testservice.SomeOtherModel{
			{ID: uuid.MustParse("b230154b-dfbb-484e-af4d-58edf690638d")},
			nil,
			{ID: uuid.MustParse("b4af1c7e-cfca-4cbd-80c5-c9232fb39a7c")},
		},
		OneOfSliceField:    []testservice.SomeOneOf{&testservice.OneOfValue1{Value: "123"}, &testservice.OneOfValue2{Value: "456"}},
		OneOfPtrSliceField: []*testservice.SomeOneOf{toPtr[testservice.SomeOneOf](&testservice.OneOfValue1{Value: "123"}), nil, toPtr[testservice.SomeOneOf](&testservice.OneOfValue2{Value: "zxc"})},
		SliceEnumField:     []testservice.SomeEnum{testservice.SomeEnumValue1, testservice.SomeEnumValue2},
		SliceEnumPtrField:  []*testservice.SomeEnum{toPtr(testservice.SomeEnumValue1), nil, toPtr(testservice.SomeEnumValue2)},
		SliceAnyField:      []any{"value", "hello", nil},
		SliceAnyPtrField:   []*any{toPtr[any]("value"), toPtr[any]("hello"), nil},
	}, serviceModel)
}
