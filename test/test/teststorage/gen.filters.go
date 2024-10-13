package teststorage

import (
	uuid "github.com/google/uuid"
	clause "gorm.io/gorm/clause"

	dbutil "github.com/saturn4er/boilerplate-go/lib/dbutil"
	testservice "github.com/saturn4er/boilerplate-go/test/test/testservice"
)

func buildSomeModelFilterExpr(filter *testservice.SomeModelFilter) (clause.Expression, error) {
	if filter == nil {
		return nil, nil
	}

	return dbutil.BuildFilterExpression(
		dbutil.ColumnFilter[uuid.UUID]{
			Column: "id",
			Filter: filter.ID,
		},
		dbutil.ExpressionBuilderFunc(func() (clause.Expression, error) {
			if filter.Or == nil {
				return nil, nil
			}
			exprs := make([]clause.Expression, 0, len(filter.Or))
			for _, orFilter := range filter.Or {
				expr, err := buildSomeModelFilterExpr(orFilter)
				if err != nil {
					return nil, err
				}
				exprs = append(exprs, expr)
			}
			return clause.Or(exprs...), nil
		}),
		dbutil.ExpressionBuilderFunc(func() (clause.Expression, error) {
			if filter.And == nil {
				return nil, nil
			}
			exprs := make([]clause.Expression, 0, len(filter.And))
			for _, andFilter := range filter.And {
				expr, err := buildSomeModelFilterExpr(andFilter)
				if err != nil {
					return nil, err
				}
				exprs = append(exprs, expr)
			}
			return clause.And(exprs...), nil
		}),
	)
}
func buildSomeOtherModelFilterExpr(filter *testservice.SomeOtherModelFilter) (clause.Expression, error) {
	if filter == nil {
		return nil, nil
	}

	return dbutil.BuildFilterExpression(
		dbutil.ColumnFilter[uuid.UUID]{
			Column: "id",
			Filter: filter.ID,
		},
		dbutil.ExpressionBuilderFunc(func() (clause.Expression, error) {
			if filter.Or == nil {
				return nil, nil
			}
			exprs := make([]clause.Expression, 0, len(filter.Or))
			for _, orFilter := range filter.Or {
				expr, err := buildSomeOtherModelFilterExpr(orFilter)
				if err != nil {
					return nil, err
				}
				exprs = append(exprs, expr)
			}
			return clause.Or(exprs...), nil
		}),
		dbutil.ExpressionBuilderFunc(func() (clause.Expression, error) {
			if filter.And == nil {
				return nil, nil
			}
			exprs := make([]clause.Expression, 0, len(filter.And))
			for _, andFilter := range filter.And {
				expr, err := buildSomeOtherModelFilterExpr(andFilter)
				if err != nil {
					return nil, err
				}
				exprs = append(exprs, expr)
			}
			return clause.And(exprs...), nil
		}),
	)
}

func buildPasswordRecoveryEventFilterExpr(filter *testservice.PasswordRecoveryEventFilter) (clause.Expression, error) {
	if filter == nil {
		return nil, nil
	}

	return dbutil.BuildFilterExpression(
		dbutil.ExpressionBuilderFunc(func() (clause.Expression, error) {
			if filter.Or == nil {
				return nil, nil
			}
			exprs := make([]clause.Expression, 0, len(filter.Or))
			for _, orFilter := range filter.Or {
				expr, err := buildPasswordRecoveryEventFilterExpr(orFilter)
				if err != nil {
					return nil, err
				}
				exprs = append(exprs, expr)
			}
			return clause.Or(exprs...), nil
		}),
		dbutil.ExpressionBuilderFunc(func() (clause.Expression, error) {
			if filter.And == nil {
				return nil, nil
			}
			exprs := make([]clause.Expression, 0, len(filter.And))
			for _, andFilter := range filter.And {
				expr, err := buildPasswordRecoveryEventFilterExpr(andFilter)
				if err != nil {
					return nil, err
				}
				exprs = append(exprs, expr)
			}
			return clause.And(exprs...), nil
		}),
	)
}
