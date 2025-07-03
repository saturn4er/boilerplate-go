package dbutil

import (
	"fmt"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gorm.io/gorm/clause"

	"github.com/saturn4er/boilerplate-go/lib/filter"
)

type ColumnArrayFilter[T any] struct {
	Column string
	Filter filter.ArrayFilter[T]
}

func (c ColumnArrayFilter[T]) buildExpression() (clause.Expression, error) {
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

func (c MappedColumnArrayFilter[T, V]) buildExpression() (clause.Expression, error) {
	if c.Filter == nil {
		return nil, nil
	}

	return ArrayFilterExpression[T, V](c.Filter, c.Column, c.Mapper)
}

func ArrayFilterExpression[T, V any](value filter.ArrayFilter[T], column string, mapper func(T) (V, error)) (clause.Expression, error) {
	switch typedValue := value.(type) {
	case *filter.ArrayContainsFilter[T]:
		return arrayContainsFilterGormCondition(typedValue, column, mapper)
	case *filter.ArrayContainsAnyFilter[T]:
		return arrayContainsAnyFilterGormCondition(typedValue, column, mapper)
	default:
		return nil, errors.Errorf("unsupported Filter type: %T", value)
	}
}

func arrayContainsFilterGormCondition[T, V any](containsFilter *filter.ArrayContainsFilter[T], column string, mapper func(T) (V, error)) (clause.Expression, error) { //nolint:lll
	if mapper == nil {
		return clause.Expr{SQL: fmt.Sprintf("%s @> ?", column), Vars: []interface{}{pq.Array(containsFilter.Values)}}, nil
	}

	values := make([]interface{}, 0, len(containsFilter.Values))

	for _, val := range containsFilter.Values {
		mappedValue, err := mapper(val)
		if err != nil {
			return nil, err
		}

		values = append(values, mappedValue)
	}

	return clause.Expr{SQL: fmt.Sprintf("%s @> ?", column), Vars: []interface{}{values}}, nil
}

func arrayContainsAnyFilterGormCondition[T, V any](containsFilter *filter.ArrayContainsAnyFilter[T], column string, mapper func(T) (V, error)) (clause.Expression, error) { //nolint:lll
	if mapper == nil {
		return clause.Expr{SQL: fmt.Sprintf("%s && ?", column), Vars: []interface{}{pq.Array(containsFilter.Values)}}, nil
	}

	values := make([]V, 0, len(containsFilter.Values))
	for _, val := range containsFilter.Values {
		mappedValue, err := mapper(val)
		if err != nil {
			return nil, err
		}
		values = append(values, mappedValue)
	}

	return clause.Expr{SQL: fmt.Sprintf("%s && ?", column), Vars: []interface{}{pq.Array(values)}}, nil
}
