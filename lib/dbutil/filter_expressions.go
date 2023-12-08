package dbutil

import (
	"github.com/pkg/errors"
	"gorm.io/gorm/clause"

	"github.com/saturn4er/boilerplate-go/lib/filter"
)

func FilterExpression[T, V any](value filter.Filter[T], column string, mapper func(T) (V, error)) (clause.Expression, error) {
	switch typedValue := value.(type) {
	case *filter.AndFilter[T]:
		return andFilterGormCondition(typedValue, column, mapper)
	case *filter.EqualsFilter[T]:
		return equalsFilterGormCondition(typedValue, column, mapper)
	case *filter.GreaterFilter[T]:
		return greaterFilterGormCondition(typedValue, column, mapper)
	case *filter.GreaterOrEqualsFilter[T]:
		return greaterOrEqualsFilterGormCondition(typedValue, column, mapper)
	case *filter.LessFilter[T]:
		return lessFilterGormCondition(typedValue, column, mapper)
	case *filter.LessOrEqualsFilter[T]:
		return lessOrEqualsFilterGormCondition(typedValue, column, mapper)
	case *filter.InFilter[T]:
		return inFilterGormCondition(typedValue, column, mapper)
	case *filter.NotEqualsFilter[T]:
		return notEqualsFilterGormCondition(typedValue, column, mapper)
	case *filter.NotInFilter[T]:
		return notInGormCondition(typedValue, column, mapper)
	case *filter.OrFilter[T]:
		return orFilterGormCondition(typedValue, column, mapper)
	}

	return nil, errors.Errorf("unsupported Filter type: %T", value)
}

func andFilterGormCondition[T, V any](andFilter *filter.AndFilter[T], column string, mapper func(T) (V, error)) (clause.Expression, error) {
	expressions := make([]clause.Expression, 0, len(andFilter.Filters))

	for _, f := range andFilter.Filters {
		expr, err := FilterExpression(f, column, mapper)
		if err != nil {
			return nil, err
		}

		expressions = append(expressions, expr)
	}

	return clause.And(expressions...), nil
}

func equalsFilterGormCondition[T, V any](equalsFilter *filter.EqualsFilter[T], column string, mapper func(T) (V, error)) (clause.Expression, error) { //nolint:lll
	if mapper == nil {
		return clause.Eq{Column: column, Value: equalsFilter.Value}, nil
	}

	mappedValue, err := mapper(equalsFilter.Value)
	if err != nil {
		return nil, err
	}

	return clause.Eq{Column: column, Value: mappedValue}, nil
}

func greaterFilterGormCondition[T, V any](greaterFilter *filter.GreaterFilter[T], column string, mapper func(T) (V, error)) (clause.Expression, error) { //nolint:lll
	if mapper == nil {
		return clause.Gt{Column: column, Value: greaterFilter.Value}, nil
	}

	mappedValue, err := mapper(greaterFilter.Value)
	if err != nil {
		return nil, err
	}

	return clause.Gt{Column: column, Value: mappedValue}, nil
}

func greaterOrEqualsFilterGormCondition[T, V any](greaterFilter *filter.GreaterOrEqualsFilter[T], column string, mapper func(T) (V, error)) (clause.Expression, error) { //nolint:lll
	if mapper == nil {
		return clause.Gte{Column: column, Value: greaterFilter.Value}, nil
	}

	mappedValue, err := mapper(greaterFilter.Value)
	if err != nil {
		return nil, err
	}

	return clause.Gte{Column: column, Value: mappedValue}, nil
}

func lessFilterGormCondition[T, V any](lessFilter *filter.LessFilter[T], column string, mapper func(T) (V, error)) (clause.Expression, error) {
	if mapper == nil {
		return clause.Lt{Column: column, Value: lessFilter.Value}, nil
	}

	mappedValue, err := mapper(lessFilter.Value)
	if err != nil {
		return nil, err
	}

	return clause.Lt{Column: column, Value: mappedValue}, nil
}

func lessOrEqualsFilterGormCondition[T, V any](lessFilter *filter.LessOrEqualsFilter[T], column string, mapper func(T) (V, error)) (clause.Expression, error) { //nolint:lll
	if mapper == nil {
		return clause.Lte{Column: column, Value: lessFilter.Value}, nil
	}

	mappedValue, err := mapper(lessFilter.Value)
	if err != nil {
		return nil, err
	}

	return clause.Lte{Column: column, Value: mappedValue}, nil
}

func inFilterGormCondition[T, V any](inFilter *filter.InFilter[T], column string, mapper func(T) (V, error)) (clause.Expression, error) {
	if mapper == nil {
		return clause.Eq{Column: column, Value: inFilter.Values}, nil
	}

	values := make([]interface{}, 0, len(inFilter.Values))

	for _, val := range inFilter.Values {
		mappedValue, err := mapper(val)
		if err != nil {
			return nil, err
		}

		values = append(values, mappedValue)
	}

	return clause.Eq{Column: column, Value: values}, nil
}

func notEqualsFilterGormCondition[T, V any](notEqualsFilter *filter.NotEqualsFilter[T], column string, mapper func(T) (V, error)) (clause.Expression, error) { //nolint:lll
	if mapper == nil {
		return clause.Neq{Column: column, Value: notEqualsFilter.Value}, nil
	}

	mappedValue, err := mapper(notEqualsFilter.Value)
	if err != nil {
		return nil, err
	}

	return clause.Neq{Column: column, Value: mappedValue}, nil
}

func notInGormCondition[T, V any](notInFilter *filter.NotInFilter[T], column string, mapper func(T) (V, error)) (clause.Expression, error) {
	if mapper == nil {
		return clause.Neq{Column: column, Value: notInFilter.Values}, nil
	}

	values := make([]interface{}, 0, len(notInFilter.Values))

	for _, el := range notInFilter.Values {
		mappedValue, err := mapper(el)
		if err != nil {
			return nil, err
		}

		values = append(values, mappedValue)
	}

	return clause.Neq{Column: column, Value: values}, nil
}

func orFilterGormCondition[T, V any](orFilter *filter.OrFilter[T], column string, mapper func(T) (V, error)) (clause.Expression, error) {
	expressions := make([]clause.Expression, 0, len(orFilter.Filters))

	for _, el := range orFilter.Filters {
		expr, err := FilterExpression(el, column, mapper)
		if err != nil {
			return nil, err
		}

		expressions = append(expressions, expr)
	}

	return clause.Or(expressions...), nil
}

type expressionBuilder interface {
	buildExpression() (clause.Expression, error)
}
type ColumnFilter[T any] struct {
	Column string
	Filter filter.Filter[T]
}

func (c ColumnFilter[T]) buildExpression() (clause.Expression, error) {
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

func (c MappedColumnFilter[T, V]) buildExpression() (clause.Expression, error) {
	if c.Filter == nil {
		return nil, nil
	}

	return FilterExpression(c.Filter, c.Column, c.Mapper)
}

func BuildFilterExpression(builders ...expressionBuilder) (clause.Expression, error) {
	expressions := make([]clause.Expression, 0, len(builders))

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

	return clause.And(expressions...), nil
}

type ExpressionBuilderFunc func() (clause.Expression, error)

func (e ExpressionBuilderFunc) buildExpression() (clause.Expression, error) {
	return e()
}
