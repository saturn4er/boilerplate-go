package segmentationstorage

import (
	uuid "github.com/google/uuid"
	clause "gorm.io/gorm/clause"

	segmentationservice "github.com/saturn4er/boilerplate-go/example/segmentation/segmentationservice"
	dbutil "github.com/saturn4er/boilerplate-go/lib/dbutil"
	// user code 'imports'
	// end user code 'imports'
)

type filterOptions struct {
	columnPrefix string
}

func withFilterColumnPrefix(prefix string) func(*filterOptions) {
	return func(f *filterOptions) {
		f.columnPrefix = prefix
	}
}
func buildUserTagFilterExpr(filter *segmentationservice.UserTagFilter, options ...func(*filterOptions)) (clause.Expression, error) {
	if filter == nil {
		return nil, nil
	}

	opts := &filterOptions{}
	for _, opt := range options {
		opt(opts)
	}

	return dbutil.BuildFilterExpression(
		dbutil.ColumnFilter[uuid.UUID]{
			Column: opts.columnPrefix + "id",
			Filter: filter.ID,
		},
		dbutil.ColumnFilter[uuid.UUID]{
			Column: opts.columnPrefix + "user_id",
			Filter: filter.UserID,
		},
		dbutil.ColumnFilter[string]{
			Column: opts.columnPrefix + "key",
			Filter: filter.Key,
		},
		dbutil.ExpressionBuilderFunc(func() (clause.Expression, error) {
			if filter.Or == nil {
				return nil, nil
			}
			exprs := make([]clause.Expression, 0, len(filter.Or))
			for _, orFilter := range filter.Or {
				expr, err := buildUserTagFilterExpr(orFilter)
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
				expr, err := buildUserTagFilterExpr(andFilter)
				if err != nil {
					return nil, err
				}
				exprs = append(exprs, expr)
			}
			return clause.And(exprs...), nil
		}),
	)
}
func buildSetUserTagCommandFilterExpr(filter *segmentationservice.SetUserTagCommandFilter, options ...func(*filterOptions)) (clause.Expression, error) {
	if filter == nil {
		return nil, nil
	}

	opts := &filterOptions{}
	for _, opt := range options {
		opt(opts)
	}

	return dbutil.BuildFilterExpression(
		dbutil.ExpressionBuilderFunc(func() (clause.Expression, error) {
			if filter.Or == nil {
				return nil, nil
			}
			exprs := make([]clause.Expression, 0, len(filter.Or))
			for _, orFilter := range filter.Or {
				expr, err := buildSetUserTagCommandFilterExpr(orFilter)
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
				expr, err := buildSetUserTagCommandFilterExpr(andFilter)
				if err != nil {
					return nil, err
				}
				exprs = append(exprs, expr)
			}
			return clause.And(exprs...), nil
		}),
	)
}

func buildSomeModelFilterExpr(filter *segmentationservice.SomeModelFilter, options ...func(*filterOptions)) (clause.Expression, error) {
	if filter == nil {
		return nil, nil
	}

	opts := &filterOptions{}
	for _, opt := range options {
		opt(opts)
	}

	return dbutil.BuildFilterExpression(
		dbutil.ColumnFilter[uuid.UUID]{
			Column: opts.columnPrefix + "id",
			Filter: filter.ID,
		},
		dbutil.ColumnFilter[string]{
			Column: opts.columnPrefix + "name",
			Filter: filter.Name,
		},
		dbutil.ColumnFilter[*string]{
			Column: opts.columnPrefix + "description",
			Filter: filter.Description,
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
func buildSomeOtherModelFilterExpr(filter *segmentationservice.SomeOtherModelFilter, options ...func(*filterOptions)) (clause.Expression, error) {
	if filter == nil {
		return nil, nil
	}

	opts := &filterOptions{}
	for _, opt := range options {
		opt(opts)
	}

	return dbutil.BuildFilterExpression(
		dbutil.ColumnFilter[uuid.UUID]{
			Column: opts.columnPrefix + "id",
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
