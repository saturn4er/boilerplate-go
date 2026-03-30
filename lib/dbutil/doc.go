// Package dbutil provides generic database utilities built on top of GORM.
//
// The core abstraction is [EntityStorage], a generic interface for CRUD operations
// on any entity type with typed filters:
//
//	type EntityStorage[Entity, Filter any] interface {
//	    Create(ctx context.Context, model *Entity) (*Entity, error)
//	    BatchCreate(ctx context.Context, models []*Entity) ([]*Entity, error)
//	    First(ctx context.Context, filter *Filter, options ...Option[SelectOptions]) (*Entity, error)
//	    Find(ctx context.Context, filter *Filter, options ...Option[SelectOptions]) ([]*Entity, error)
//	    Update(ctx context.Context, model *Entity) (*Entity, error)
//	    Save(ctx context.Context, model *Entity) (*Entity, error)
//	    Delete(ctx context.Context, filter *Filter) error
//	    Count(ctx context.Context, filter *Filter) (int, error)
//	    WithAdvisoryLock(ctx context.Context, lockID int64) error
//	}
//
// [GormEntityStorage] is the concrete GORM implementation that handles bidirectional
// conversion between external (domain) and internal (database) types, filter expression
// building, and PostgreSQL advisory locking.
//
// The package also provides filter expression builders that translate typed filter
// structs into GORM WHERE clauses:
//
//   - [ColumnFilter]: Maps a single filter value to a database column.
//   - [MappedColumnArrayFilter]: Handles array fields with element transformation.
//   - [BuildFilterExpression]: Composes multiple filter builders into a single GORM expression.
//
// Select options support pagination, ordering, row locking (FOR UPDATE), and field selection:
//
//	users, err := storage.Find(ctx, &UserFilter{
//	    Email: filter.Equals("alice@example.com"),
//	}, dbutil.WithPagination(&pagination.Pagination{Page: 1, PerPage: 10}),
//	   dbutil.WithOrder(UserFieldName, dbutil.OrderDirAsc),
//	   dbutil.WithForUpdate())
package dbutil
