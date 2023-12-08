package dbutil

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"

	"github.com/saturn4er/boilerplate-go/lib/filter"
)

type ColumnArrayFilter[T any] struct {
	Column string
	Filter filter.ArrayFilter[T]
}

func (c ColumnArrayFilter[T]) buildExpression() (goqu.Expression, error) {
	if c.Filter == nil {
		return nil, nil
	}

	return ArrayFilterExpression[T, any](c.Filter, c.Column, nil)
}

type MappedColumnArrayFilter[T, V any] struct {
	Column string
	Filter filter.ArrayFilter[T]
	Mapper func(T) (V, error)
}

func (c MappedColumnArrayFilter[T, V]) buildExpression() (goqu.Expression, error) {
	if c.Filter == nil {
		return nil, nil
	}

	return ArrayFilterExpression[T, V](c.Filter, c.Column, c.Mapper)
}

func ArrayFilterExpression[T, V any](value filter.ArrayFilter[T], column string, mapper func(T) (V, error)) (goqu.Expression, error) {
	if typedValue, ok := value.(*filter.ArrayContainsFilter[T]); ok {
		return arrayContainsFilterGormCondition(typedValue, column, mapper)
	}

	return nil, errors.Errorf("unsupported Filter type: %T", value)
}

func arrayContainsFilterGormCondition[T, V any](containsFilter *filter.ArrayContainsFilter[T], column string, mapper func(T) (V, error)) (goqu.Expression, error) { //nolint:lll
	if mapper == nil {
		return goqu.L(fmt.Sprintf("%s @> ?", column), containsFilter.Values), nil
	}

	values := make([]interface{}, 0, len(containsFilter.Values))

	for _, val := range containsFilter.Values {
		mappedValue, err := mapper(val)
		if err != nil {
			return nil, err
		}

		values = append(values, mappedValue)
	}

	return goqu.L(fmt.Sprintf("%s @> ?", column), values), nil
}
