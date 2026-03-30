// Package filter provides composable, type-safe query filter primitives.
//
// Filters are used with generated storage code to build database queries.
// The [Filter] interface is parameterized by the value type, ensuring
// compile-time type safety.
//
// # Basic Filters
//
//   - [Equals]: Exact match (column = value)
//   - [NotEquals]: Negated match (column != value)
//   - [In]: Set membership (column IN (values...))
//   - [NotIn]: Negated set membership (column NOT IN (values...))
//   - [IsNull]: Null check (column IS NULL)
//   - [IsNotNull]: Non-null check (column IS NOT NULL)
//
// # Comparison Filters
//
//   - [Greater]: column > value
//   - [Less]: column < value
//   - [GreaterOrEquals]: column >= value
//   - [LessOrEquals]: column <= value
//
// # String Filters
//
//   - [Contains]: Substring match (column LIKE '%value%')
//   - [HasPrefix]: Prefix match (column LIKE 'value%')
//   - [HasSuffix]: Suffix match (column LIKE '%value')
//
// # Array Filters (PostgreSQL)
//
//   - [ArrayContains]: Array contains all specified elements
//   - [ArrayContainsAny]: Array contains at least one element
//   - [ArrayDoesNotContainAny]: Array contains none of the elements
//   - [ArrayIsEmpty]: Array has no elements
//   - [ArrayIsNotEmpty]: Array has at least one element
//
// # Logical Composition
//
// Filters compose through And/Or fields on generated filter structs:
//
//	filter := &UserFilter{
//	    Or: []*UserFilter{
//	        {Email: filter.Equals("alice@example.com")},
//	        {Email: filter.Equals("bob@example.com")},
//	    },
//	}
package filter
