package teststorage

import (
	driver "database/sql/driver"
	json "encoding/json"
	"net"

	uuid "github.com/google/uuid"
	pq "github.com/lib/pq"
	errors "github.com/pkg/errors"
	// user code 'imports'
	// end user code 'imports'
)

type mapValue[C comparable, B any] map[string]B

func (m mapValue[C, B]) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *mapValue[C, B]) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), m)
}

type ipValue net.IP

func (i *ipValue) Value() (driver.Value, error) {
	return (*net.IP)(i).String(), nil
}

func (i *ipValue) Scan(src any) error {
	switch src := src.(type) {
	case string:
		*i = ipValue(net.ParseIP(src))
	case []byte:
		*i = ipValue(net.ParseIP(string(src)))
	default:
		return errors.Errorf("can't parse ipValue from: %T", src)
	}
	return nil
}

type stringSliceValue []string

func (s stringSliceValue) Value() (driver.Value, error) {
	result := make(pq.StringArray, 0, len(s))
	for _, item := range s {
		result = append(result, item)
	}

	return result.Value()
}

func (s *stringSliceValue) Scan(src interface{}) error {
	result := make(pq.StringArray, 0)
	err := result.Scan(src)
	if err != nil {
		return err
	}

	*s = stringSliceValue(result)

	return nil
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

type uuidSliceValue []uuid.UUID

func (s uuidSliceValue) Value() (driver.Value, error) {
	result := make(pq.StringArray, len(s))
	for i, u := range s {
		result[i] = u.String()
	}
	return result.Value()
}

func (s *uuidSliceValue) Scan(src interface{}) error {
	var result pq.StringArray
	if err := result.Scan(src); err != nil {
		return err
	}

	*s = make(uuidSliceValue, len(result))
	for i, str := range result {
		parsedUUID, err := uuid.Parse(str)
		if err != nil {
			return err
		}
		(*s)[i] = parsedUUID
	}

	return nil
}

func fromPtr[T any](ptr *T) T {
	return *ptr
}

func toPtr[T any](val T) *T {
	return &val
}
