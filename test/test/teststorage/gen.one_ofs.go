package teststorage

import (
	driver "database/sql/driver"
	json "encoding/json"
	fmt "fmt"

	testservice "github.com/saturn4er/boilerplate-go/test/test/testservice"
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

func convertSomeOneOfToDB(val testservice.SomeOneOf) (*jsonSomeOneOf, error) {
	if val == nil {
		return nil, nil
	}
	result := &jsonSomeOneOf{}
	switch v := val.(type) {
	case *testservice.OneOfValue1:
		if v != nil {

			tmp, err := convertOneOfValue1ToJsonModel(v)
			if err != nil {
				return nil, fmt.Errorf("convert OneOfValue1 to db: %w", err)
			}
			result.Val = tmp
		} else {
			result.Val = nil
		}
		result.OneOfType = "OneOfValue1"
		result.OneOfTypeID = 1

		return result, nil
	case *testservice.OneOfValue2:
		if v != nil {

			tmp, err := convertOneOfValue2ToJsonModel(v)
			if err != nil {
				return nil, fmt.Errorf("convert OneOfValue2 to db: %w", err)
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

func convertSomeOneOfFromDB(val *jsonSomeOneOf) (testservice.SomeOneOf, error) {
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

type jsonPasswordRecoveryEventData struct {
	Val         any    `json:"value"`
	OneOfType   string `json:"@type"`
	OneOfTypeID uint   `json:"@type_id"`
}

func (p *jsonPasswordRecoveryEventData) UnmarshalJSON(bytes []byte) error {
	tmp := struct {
		OneOfTypeID uint   `json:"@type_id"`
		OneOfType   string `json:"@type"`
	}{}
	if err := json.Unmarshal(bytes, &tmp); err != nil {
		return fmt.Errorf("unmarshal OneOfType: %w", err)
	}

	switch tmp.OneOfTypeID {
	case 100001:
		var value struct {
			Value jsonPasswordRecoveryRequestedEventData `json:"value"`
		}
		if err := json.Unmarshal(bytes, &value); err != nil {
			return err
		}
		p.Val = &value.Value
	case 100002:
		var value struct {
			Value jsonPasswordRecoveryCompletedEventData `json:"value"`
		}
		if err := json.Unmarshal(bytes, &value); err != nil {
			return err
		}
		p.Val = &value.Value
	}
	return nil
}
func (p *jsonPasswordRecoveryEventData) Scan(value any) error {
	return json.Unmarshal(value.([]byte), p)
}

func (p jsonPasswordRecoveryEventData) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func convertPasswordRecoveryEventDataToDB(val testservice.PasswordRecoveryEventData) (*jsonPasswordRecoveryEventData, error) {
	if val == nil {
		return nil, nil
	}
	result := &jsonPasswordRecoveryEventData{}
	switch v := val.(type) {
	case *testservice.PasswordRecoveryRequestedEventData:
		if v != nil {

			tmp, err := convertPasswordRecoveryRequestedEventDataToJsonModel(v)
			if err != nil {
				return nil, fmt.Errorf("convert PasswordRecoveryRequestedEventData to db: %w", err)
			}
			result.Val = tmp
		} else {
			result.Val = nil
		}
		result.OneOfType = "PasswordRecoveryRequestedEventData"
		result.OneOfTypeID = 100001

		return result, nil
	case *testservice.PasswordRecoveryCompletedEventData:
		if v != nil {

			tmp, err := convertPasswordRecoveryCompletedEventDataToJsonModel(v)
			if err != nil {
				return nil, fmt.Errorf("convert PasswordRecoveryCompletedEventData to db: %w", err)
			}
			result.Val = tmp
		} else {
			result.Val = nil
		}
		result.OneOfType = "PasswordRecoveryCompletedEventData"
		result.OneOfTypeID = 100002

		return result, nil
	}
	return nil, fmt.Errorf("invalid PasswordRecoveryEventData value type: %T", val)
}

func convertPasswordRecoveryEventDataFromDB(val *jsonPasswordRecoveryEventData) (testservice.PasswordRecoveryEventData, error) {
	if val == nil {
		return nil, nil
	}

	switch v := (*val).Val.(type) {
	case *jsonPasswordRecoveryRequestedEventData:
		v1, err := convertPasswordRecoveryRequestedEventDataFromJsonModel(v)
		if err != nil {
			return nil, fmt.Errorf("convert PasswordRecoveryRequestedEventData from db: %w", err)
		}

		return v1, nil
	case *jsonPasswordRecoveryCompletedEventData:
		v1, err := convertPasswordRecoveryCompletedEventDataFromJsonModel(v)
		if err != nil {
			return nil, fmt.Errorf("convert PasswordRecoveryCompletedEventData from db: %w", err)
		}

		return v1, nil
	default:
		return nil, fmt.Errorf("invalid PasswordRecoveryEventData value type: %T", *val)
	}

	panic("implement me")
}
