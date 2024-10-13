package teststorage

import (
	driver "database/sql/driver"
	json "encoding/json"

	pq "github.com/lib/pq"
)

type mapValue[C comparable, B any] map[string]B

func (m mapValue[C, B]) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *mapValue[C, B]) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), m)
}

type sliceValue[B any] []B

func (s sliceValue[B]) Value() (driver.Value, error) {
	result := make(pq.StringArray, 0, len(s))
	for _, item := range s {
		tmp, err := json.Marshal(item)
		if err != nil {
			return nil, err
		}
		result = append(result, string(tmp))
	}

	return result.Value()
}

func (s *sliceValue[B]) Scan(src interface{}) error {
	result := make(pq.StringArray, 0)
	err := result.Scan(src)
	if err != nil {
		return err
	}

	*s = make([]B, 0, len(result))

	for _, item := range result {
		var tmp B
		err = json.Unmarshal([]byte(item), &tmp)
		if err != nil {
			return err
		}
		*s = append(*s, tmp)
	}

	return nil
}

func fromPtr[T any](ptr *T) T {
	return *ptr
}

func toPtr[T any](val T) *T {
	return &val
}
