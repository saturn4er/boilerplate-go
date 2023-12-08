package dbutil

import (
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/stretchr/testify/require"

	"github.com/saturn4er/boilerplate-go/lib/filter"
)

func TestOr(t *testing.T) {
	expr, err := BuildFilterExpression(
		ColumnFilter[string]{
			Column: "column_a",
			Filter: filter.Equals("value_a"),
		},
		ColumnFilter[int]{
			Column: "column_b",
			Filter: filter.Greater(10),
		},
	)
	require.NoError(t, err)
	query, params, err := goqu.Select("*").Where(expr).Prepared(true).From("table").ToSQL()
	require.NoError(t, err)
	require.Equal(t, `SELECT * FROM "table" WHERE (("column_a" = ?) AND ("column_b" > ?))`, query)
	require.Equal(t, []interface{}{"value_a", int64(10)}, params)
}
