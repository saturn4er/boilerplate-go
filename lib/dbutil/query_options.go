package dbutil

import (
	"github.com/go-pnp/go-pnp/pkg/optionutil"
	"gorm.io/gorm/clause"

	"github.com/saturn4er/boilerplate-go/lib/pagination"
)

type ClausesBuilder func() func() ([]clause.Expression, error)

type SelectOptions struct {
	Pagination *pagination.Pagination
	Order      []Order
	ForUpdate  bool
}

func (s *SelectOptions) BuildExpressions(fieldMapping map[any]clause.Column) ([]clause.Expression, error) {
	var exprs []clause.Expression

	if s.Pagination != nil {
		expr, err := BuildPaginationExpression(s.Pagination)
		if err != nil {
			return nil, err
		}

		exprs = append(exprs, expr)
	}

	if len(s.Order) > 0 {
		expr, err := BuildOrderExpression(s.Order, fieldMapping)
		if err != nil {
			return nil, err
		}

		exprs = append(exprs, expr)
	}

	if s.ForUpdate {
		exprs = append(exprs, BuildForUpdateExpression())
	}

	return exprs, nil
}

func WithPagination(pagination *pagination.Pagination) optionutil.Option[SelectOptions] {
	return func(options *SelectOptions) {
		options.Pagination = pagination
	}
}

func WithOrder(field any, direction OrderDirection) optionutil.Option[SelectOptions] {
	return func(options *SelectOptions) {
		options.Order = append(options.Order, Order{
			Field:     field,
			Direction: direction,
		})
	}
}

func WithForUpdate() optionutil.Option[SelectOptions] {
	return func(options *SelectOptions) {
		options.ForUpdate = true
	}
}

func BuildPaginationExpression(pagination *pagination.Pagination) (clause.Expression, error) {
	return clause.Limit{
		Limit:  &pagination.PerPage,
		Offset: (pagination.Page - 1) * pagination.PerPage,
	}, nil
}

type OrderDirection byte

const (
	OrderDirAsc OrderDirection = iota + 1
	OrderDirDesc
)

type Order struct {
	Field     any
	Direction OrderDirection
}

func BuildOrderExpression(orders []Order, fieldsMapping map[any]clause.Column) (clause.Expression, error) {
	var result clause.OrderBy

	for _, order := range orders {
		column, ok := fieldsMapping[order.Field]
		if !ok {
			return nil, ErrInvalidField
		}

		result.Columns = append(result.Columns, clause.OrderByColumn{
			Column: column,
			Desc:   order.Direction == OrderDirDesc,
		})
	}

	return result, nil
}

func BuildForUpdateExpression() clause.Expression {
	return clause.Locking{
		Strength: "UPDATE",
	}
}
