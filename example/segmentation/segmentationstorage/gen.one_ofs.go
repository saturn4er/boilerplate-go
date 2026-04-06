package segmentationstorage

import (
	driver "database/sql/driver"
	json "encoding/json"
	fmt "fmt"

	errors "github.com/pkg/errors"

	segmentationservice "github.com/saturn4er/boilerplate-go/example/segmentation/segmentationservice"
	// user code 'imports'
	// end user code 'imports'
)

type jsonSomeOneOf struct {
	Val         any    `json:"value"`
	OneOfType   string `json:"@type"`
	OneOfTypeID uint   `json:"@type_id"`
}

func (s *jsonSomeOneOf) UnmarshalJSON(bytes []byte) error {
	tmp := struct {
		OneOfTypeID uint   `json:"@type_id"`
		OneOfType   string `json:"@type"`
	}{}
	if err := json.Unmarshal(bytes, &tmp); err != nil {
		return fmt.Errorf("unmarshal OneOfType: %w", err)
	}

	switch tmp.OneOfTypeID {
	case 1:
		var value struct {
			Value jsonOneOfValue1 `json:"value"`
		}
		if err := json.Unmarshal(bytes, &value); err != nil {
			return err
		}
		s.Val = &value.Value
	case 2:
		var value struct {
			Value jsonOneOfValue2 `json:"value"`
		}
		if err := json.Unmarshal(bytes, &value); err != nil {
			return err
		}
		s.Val = &value.Value
	}
	return nil
}
func (s *jsonSomeOneOf) Scan(value any) error {
	return json.Unmarshal(value.([]byte), s)
}

func (s jsonSomeOneOf) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func convertSomeOneOfToDB(val segmentationservice.SomeOneOf) (*jsonSomeOneOf, error) {
	if val == nil {
		return nil, nil
	}
	result := &jsonSomeOneOf{}
	switch v := val.(type) {
	case *segmentationservice.OneOfValue1:
		if v != nil {
			tmp, err := convertOneOfValue1ToJsonModel(v)
			if err != nil {
				return nil, errors.Wrap(err, "convert OneOfValue1 to db")
			}
			result.Val = tmp
		} else {
			result.Val = nil
		}
		result.OneOfType = "OneOfValue1"
		result.OneOfTypeID = 1

		return result, nil
	case *segmentationservice.OneOfValue2:
		if v != nil {
			tmp, err := convertOneOfValue2ToJsonModel(v)
			if err != nil {
				return nil, errors.Wrap(err, "convert OneOfValue2 to db")
			}
			result.Val = tmp
		} else {
			result.Val = nil
		}
		result.OneOfType = "OneOfValue2"
		result.OneOfTypeID = 2

		return result, nil
	}
	return nil, fmt.Errorf("invalid SomeOneOf value type: %T", val)
}

func convertSomeOneOfFromDB(val *jsonSomeOneOf) (segmentationservice.SomeOneOf, error) {
	if val == nil {
		return nil, nil
	}

	switch v := (*val).Val.(type) {
	case *jsonOneOfValue1:
		v1, err := convertOneOfValue1FromJsonModel(v)
		if err != nil {
			return nil, fmt.Errorf("convert OneOfValue1 from db: %w", err)
		}

		return v1, nil
	case *jsonOneOfValue2:
		v1, err := convertOneOfValue2FromJsonModel(v)
		if err != nil {
			return nil, fmt.Errorf("convert OneOfValue2 from db: %w", err)
		}

		return v1, nil
	default:
		return nil, fmt.Errorf("invalid SomeOneOf value type: %T", *val)
	}

	panic("implement me")
}
