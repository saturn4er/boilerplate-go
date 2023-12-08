package dbutil

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/pkg/errors"

	"github.com/saturn4er/boilerplate-go/lib/filter"
)

func FilterExpression[T, V any](value filter.Filter[T], column string, mapper func(T) (V, error)) (goqu.Expression, error) {
	switch typedValue := value.(type) {
	case *filter.AndFilter[T]:
		return andFilterGoquCondition(typedValue, column, mapper)
	case *filter.EqualsFilter[T]:
		return equalsFilterGoquCondition(typedValue, column, mapper)
	case *filter.GreaterFilter[T]:
		return greaterFilterGoquCondition(typedValue, column, mapper)
	case *filter.GreaterOrEqualsFilter[T]:
		return greaterOrEqualsFilterGoquCondition(typedValue, column, mapper)
	case *filter.LessFilter[T]:
		return lessFilterGoquCondition(typedValue, column, mapper)
	case *filter.LessOrEqualsFilter[T]:
		return lessOrEqualsFilterGoquCondition(typedValue, column, mapper)
	case *filter.InFilter[T]:
		return inFilterGoquCondition(typedValue, column, mapper)
	case *filter.NotEqualsFilter[T]:
		return notEqualsFilterGoquCondition(typedValue, column, mapper)
	case *filter.NotInFilter[T]:
		return notInGoquCondition(typedValue, column, mapper)
	case *filter.OrFilter[T]:
		return orFilterGoquCondition(typedValue, column, mapper)
	}

	return nil, errors.Errorf("unsupported Filter type: %T", value)
}

func andFilterGoquCondition[T, V any](andFilter *filter.AndFilter[T], column string, mapper func(T) (V, error)) (goqu.Expression, error) {
	expressions := make([]goqu.Expression, 0, len(andFilter.Filters))

	for _, f := range andFilter.Filters {
		expr, err := FilterExpression(f, column, mapper)
		if err != nil {
			return nil, err
		}

		expressions = append(expressions, expr)
	}

	return goqu.And(expressions...), nil
}

func equalsFilterGoquCondition[T, V any](equalsFilter *filter.EqualsFilter[T], column string, mapper func(T) (V, error)) (goqu.Expression, error) { //nolint:lll
	if mapper == nil {
		return goqu.Ex{column: equalsFilter.Value}, nil
	}

	mappedValue, err := mapper(equalsFilter.Value)
	if err != nil {
		return nil, err
	}

	return goqu.Ex{column: mappedValue}, nil
}

func greaterFilterGoquCondition[T, V any](greaterFilter *filter.GreaterFilter[T], column string, mapper func(T) (V, error)) (goqu.Expression, error) { //nolint:lll
	if mapper == nil {
		return goqu.C(column).Gt(greaterFilter.Value), nil
	}

	mappedValue, err := mapper(greaterFilter.Value)
	if err != nil {
		return nil, err
	}

	return goqu.C(column).Gt(mappedValue), nil
}

func greaterOrEqualsFilterGoquCondition[T, V any](greaterFilter *filter.GreaterOrEqualsFilter[T], column string, mapper func(T) (V, error)) (goqu.Expression, error) { //nolint:lll
	if mapper == nil {
		return goqu.C(column).Gte(greaterFilter.Value), nil
	}

	mappedValue, err := mapper(greaterFilter.Value)
	if err != nil {
		return nil, err
	}

	return goqu.C(column).Gte(mappedValue), nil
}

func lessFilterGoquCondition[T, V any](lessFilter *filter.LessFilter[T], column string, mapper func(T) (V, error)) (goqu.Expression, error) {
	if mapper == nil {
		return goqu.C(column).Lt(lessFilter.Value), nil
	}

	mappedValue, err := mapper(lessFilter.Value)
	if err != nil {
		return nil, err
	}

	return goqu.C(column).Lt(mappedValue), nil
}

func lessOrEqualsFilterGoquCondition[T, V any](lessFilter *filter.LessOrEqualsFilter[T], column string, mapper func(T) (V, error)) (goqu.Expression, error) { //nolint:lll
	if mapper == nil {
		return goqu.C(column).Lte(lessFilter.Value), nil
	}

	mappedValue, err := mapper(lessFilter.Value)
	if err != nil {
		return nil, err
	}

	return goqu.C(column).Lte(mappedValue), nil
}

func inFilterGoquCondition[T, V any](inFilter *filter.InFilter[T], column string, mapper func(T) (V, error)) (goqu.Expression, error) {
	if mapper == nil {
		return goqu.C(column).In(inFilter.Values), nil
	}

	values := make([]interface{}, 0, len(inFilter.Values))

	for _, val := range inFilter.Values {
		mappedValue, err := mapper(val)
		if err != nil {
			return nil, err
		}

		values = append(values, mappedValue)
	}

	return goqu.C(column).In(values), nil
}

func notEqualsFilterGoquCondition[T, V any](notEqualsFilter *filter.NotEqualsFilter[T], column string, mapper func(T) (V, error)) (goqu.Expression, error) { //nolint:lll
	if mapper == nil {
		return goqu.C(column).Neq(notEqualsFilter.Value), nil
	}

	mappedValue, err := mapper(notEqualsFilter.Value)
	if err != nil {
		return nil, err
	}

	return goqu.C(column).Neq(mappedValue), nil
}

func notInGoquCondition[T, V any](notInFilter *filter.NotInFilter[T], column string, mapper func(T) (V, error)) (goqu.Expression, error) {
	if mapper == nil {
		return goqu.C(column).NotIn(notInFilter.Values), nil
	}

	values := make([]interface{}, 0, len(notInFilter.Values))

	for _, el := range notInFilter.Values {
		mappedValue, err := mapper(el)
		if err != nil {
			return nil, err
		}

		values = append(values, mappedValue)
	}

	return goqu.C(column).NotIn(values), nil
}

func orFilterGoquCondition[T, V any](orFilter *filter.OrFilter[T], column string, mapper func(T) (V, error)) (goqu.Expression, error) {
	expressions := make([]goqu.Expression, 0, len(orFilter.Filters))

	for _, el := range orFilter.Filters {
		expr, err := FilterExpression(el, column, mapper)
		if err != nil {
			return nil, err
		}

		expressions = append(expressions, expr)
	}

	return goqu.Or(expressions...), nil
}

type expressionBuilder interface {
	buildExpression() (goqu.Expression, error)
}
type ColumnFilter[T any] struct {
	Column string
	Filter filter.Filter[T]
}

func (c ColumnFilter[T]) buildExpression() (goqu.Expression, error) {
	if c.Filter == nil {
		return nil, nil
	}

	return FilterExpression[T, any](c.Filter, c.Column, nil)
}

type MappedColumnFilter[T, V any] struct {
	Column string
	Filter filter.Filter[T]
	Mapper func(T) (V, error)
}

func (c MappedColumnFilter[T, V]) buildExpression() (goqu.Expression, error) {
	if c.Filter == nil {
		return nil, nil
	}

	return FilterExpression(c.Filter, c.Column, c.Mapper)
}

func BuildFilterExpression(builders ...expressionBuilder) (goqu.Expression, error) {
	expressions := make([]goqu.Expression, 0, len(builders))

	for _, builder := range builders {
		expr, err := builder.buildExpression()
		if err != nil {
			return nil, err
		}

		if expr == nil {
			continue
		}

		expressions = append(expressions, expr)
	}

	return goqu.And(expressions...), nil
}

type ExpressionBuilderFunc func() (goqu.Expression, error)

func (e ExpressionBuilderFunc) buildExpression() (goqu.Expression, error) {
	return e()
}
