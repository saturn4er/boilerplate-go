package dbutil

import (
	"testing"

	"gorm.io/gorm/clause"

	"github.com/saturn4er/boilerplate-go/lib/filter"
)

func TestFilterExpression_StringFilters(t *testing.T) {
	tests := []struct {
		name           string
		filter         filter.Filter[string]
		expectedClause clause.Expression
	}{
		{
			name:   "HasPrefix filter",
			filter: filter.HasPrefix[string]("hello"),
			expectedClause: clause.Like{
				Column: "test_column",
				Value:  "hello%",
			},
		},
		{
			name:   "HasSuffix filter",
			filter: filter.HasSuffix[string]("world"),
			expectedClause: clause.Like{
				Column: "test_column",
				Value:  "%world",
			},
		},
		{
			name:   "Contains filter",
			filter: filter.Contains[string]("test"),
			expectedClause: clause.Like{
				Column: "test_column",
				Value:  "%test%",
			},
		},
		{
			name:   "Equals filter",
			filter: filter.Equals("exact"),
			expectedClause: clause.Eq{
				Column: "test_column",
				Value:  "exact",
			},
		},
		{
			name:   "NotEquals filter",
			filter: filter.NotEquals("notthis"),
			expectedClause: clause.Neq{
				Column: "test_column",
				Value:  "notthis",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := FilterExpression[string, any](tt.filter, "test_column", nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if expr == nil {
				t.Fatal("expected expression, got nil")
			}

			// Compare expressions
			if expr != tt.expectedClause {
				t.Errorf("expected %+v, got %+v", tt.expectedClause, expr)
			}
		})
	}
}

func TestFilterExpression_StringPtrFilters(t *testing.T) {
	tests := []struct {
		name           string
		filter         filter.Filter[*string]
		expectedClause clause.Expression
	}{
		{
			name:   "HasPrefix filter for *string",
			filter: filter.HasPrefix[*string]("hello"),
			expectedClause: clause.Like{
				Column: "test_column",
				Value:  "hello%",
			},
		},
		{
			name:   "HasSuffix filter for *string",
			filter: filter.HasSuffix[*string]("world"),
			expectedClause: clause.Like{
				Column: "test_column",
				Value:  "%world",
			},
		},
		{
			name:   "Contains filter for *string",
			filter: filter.Contains[*string]("test"),
			expectedClause: clause.Like{
				Column: "test_column",
				Value:  "%test%",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := FilterExpression[*string, any](tt.filter, "test_column", nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if expr == nil {
				t.Fatal("expected expression, got nil")
			}

			// Compare expressions
			if expr != tt.expectedClause {
				t.Errorf("expected %+v, got %+v", tt.expectedClause, expr)
			}
		})
	}
}

func TestColumnFilter_String(t *testing.T) {
	t.Run("with HasPrefix filter", func(t *testing.T) {
		cf := ColumnFilter[string]{
			Column: "name",
			Filter: filter.HasPrefix[string]("John"),
		}

		expr, err := cf.buildExpression()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := clause.Like{
			Column: "name",
			Value:  "John%",
		}

		if expr != expected {
			t.Errorf("expected %+v, got %+v", expected, expr)
		}
	})

	t.Run("with Contains filter", func(t *testing.T) {
		cf := ColumnFilter[string]{
			Column: "description",
			Filter: filter.Contains[string]("test"),
		}

		expr, err := cf.buildExpression()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := clause.Like{
			Column: "description",
			Value:  "%test%",
		}

		if expr != expected {
			t.Errorf("expected %+v, got %+v", expected, expr)
		}
	})

	t.Run("with nil filter", func(t *testing.T) {
		cf := ColumnFilter[string]{
			Column: "name",
			Filter: nil,
		}

		expr, err := cf.buildExpression()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if expr != nil {
			t.Errorf("expected nil expression, got %+v", expr)
		}
	})
}
