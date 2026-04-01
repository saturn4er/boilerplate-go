// Package boilerplatego is a code generation tool that produces Go service and storage
// layer boilerplate from YAML specifications.
//
// It reads YAML configuration files describing models, enums, one-of types, and events,
// then generates type-safe CRUD operations, composable filters, GORM-backed database models,
// transactional outbox support, and GoAdmin integrations.
//
// # Architecture
//
// The generator produces two packages per module:
//
//   - moduleservice: The domain layer containing public types, storage interfaces, and enums.
//   - modulestorage: The infrastructure layer with GORM implementations, model converters,
//     and filter expression builders.
//
// # Usage
//
// Define your models in YAML:
//
//	types:
//	  models:
//	    - name: User
//	      table_name: users
//	      fields:
//	        - { name: ID, type: uuid, primary_key: true, filterable: true }
//	        - { name: Email, type: string, filterable: true }
//	        - { name: Name, type: string }
//
// Run the generator:
//
//	boilerplate-go gen.yaml
//
// Or use go generate:
//
//	//go:generate boilerplate-go ./gen.yaml
//
// # Runtime Libraries
//
// The lib/ directory provides runtime packages used by generated code:
//
//   - [github.com/saturn4er/boilerplate-go/lib/dbutil]: Generic EntityStorage interface and GORM implementation.
//   - [github.com/saturn4er/boilerplate-go/lib/filter]: Composable typed query filters (Equals, In, Contains, etc.).
//
// # Unique Index Error Handling
//
// Models with unique_indexes generate per-constraint ConflictError variables,
// allowing callers to distinguish which unique constraint was violated.
//   - [github.com/saturn4er/boilerplate-go/lib/pagination]: Generic pagination types.
//   - [github.com/saturn4er/boilerplate-go/lib/order]: Generic ordering with direction support.
//   - [github.com/saturn4er/boilerplate-go/lib/idempotency]: Duplicate message detection.
//   - [github.com/saturn4er/boilerplate-go/lib/txoutbox]: Transactional outbox pattern for reliable event delivery.
package boilerplatego
