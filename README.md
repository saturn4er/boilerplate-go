# boilerplate-go

A Go code generation tool that produces service and storage layer boilerplate from YAML specifications. Define your models, enums, and events in YAML, and the tool generates type-safe CRUD operations, filters, database models, transactional outbox support, and admin UI integrations.

## Features

- **Service layer generation** — domain models, storage interfaces, enums with validation helpers, error types
- **Storage layer generation** — GORM-backed implementations, bidirectional model converters, filter expression builders
- **Transactional outbox** — event models with atomic persistence and async message delivery via [Watermill](https://github.com/ThreeDotsLabs/watermill)
- **Cross-module events** — modules can produce events owned by other modules
- **Typed filters** — composable query filters (`Equals`, `In`, `Contains`, `ArrayContains`, logical `And`/`Or`, etc.)
- **Pagination and ordering** — generic pagination and ordering types
- **Idempotency** — built-in duplicate message detection
- **GoAdmin integration** — auto-generated admin table configurations (can be disabled)
- **User code preservation** — custom code blocks survive regeneration

## Installation

```bash
go install github.com/saturn4er/boilerplate-go/cmd/boilerplate-go@latest
```

## Quick Start

### 1. Create a configuration file

**gen.yaml**
```yaml
root_package_name: "github.com/yourorg/myapp"
goimports_local: "github.com/yourorg/myapp"

modules:
  auth:
    imports:
      - auth.module.yaml
```

### 2. Define your module

**auth.module.yaml**
```yaml
types:
  enums:
    - name: Role
      values: [Admin, User]
  models:
    - name: User
      table_name: users
      fields:
        - { name: ID, type: uuid, primary_key: true, filterable: true }
        - { name: Email, type: string, filterable: true }
        - { name: Name, type: string }
        - { name: Role, type: Role }
```

### 3. Generate code

```bash
boilerplate-go gen.yaml
```

Or with `go generate`:

```go
//go:generate boilerplate-go ./gen.yaml
```

This produces:

```
auth/
├── authservice/
│   ├── gen.models.go       # User struct, Role enum, UserFilter
│   ├── gen.enums.go        # Role enum helpers (Validate, String, etc.)
│   ├── gen.storages.go     # Storage interface with Users() accessor
│   └── gen.errors.go       # Domain error types
└── authstorage/
    ├── gen.storages.go     # GORM storage implementation
    ├── gen.table_models.go # DB models with GORM tags, converters
    ├── gen.filters.go      # Filter expression builders
    ├── gen.enums.go        # Enum DB converters
    ├── gen.errors.go       # DB error wrapping
    └── gen.utils.go        # Helpers
```

### 4. Use the generated code

```go
// Initialize
storage := authstorage.NewStorages(db, logger, processors)
svc := &authservice.Service{Storage: storage}

// Create
user, err := svc.Storage.Users().Create(ctx, &authservice.User{
    ID:    uuid.New(),
    Email: "alice@example.com",
    Name:  "Alice",
    Role:  authservice.RoleUser,
})

// Query with typed filters
user, err := svc.Storage.Users().First(ctx, &authservice.UserFilter{
    Email: filter.Equals("alice@example.com"),
})

// Transactions
err := svc.Storage.ExecuteInTransaction(ctx, func(ctx context.Context, tx authservice.Storage) error {
    // all operations within tx are atomic
    return nil
})
```

## Configuration Reference

### Root Config (`gen.yaml`)

| Field | Description |
|---|---|
| `root_package_name` | Go module path for generated packages |
| `goimports_local` | Local import prefix for goimports grouping |
| `disable_admin` | Disable GoAdmin table generator files (default: `false`) |
| `modules` | Map of module name to module config |
| `types` | Global shared types (enums, models, one_ofs) |

### Module Config

Modules can be defined inline or imported from separate YAML files.

```yaml
modules:
  auth:
    imports:
      - auth.module.yaml      # Import types from file
    produces:
      - other_module.EventName # Cross-module event production
```

### Model Fields

| Field | Type | Description |
|---|---|---|
| `name` | string | Field name (PascalCase) |
| `type` | string | Field type (see [Types](#types)) |
| `primary_key` | bool | Mark as primary key |
| `filterable` | bool | Generate filter support for this field |
| `db_name` | string | Override database column name |

### Model Options

| Field | Type | Description |
|---|---|---|
| `name` | string | Model name (PascalCase) |
| `table_name` | string | Database table name |
| `storage_type` | string | `"tx_outbox"` for event models |
| `do_not_persists` | bool | Transient model, not stored in DB |
| `no_local_outbox` | bool | Don't generate outbox in owning module |
| `message_builder` | string | Custom function to build outbox messages |
| `unique_indexes` | list | Unique index constraints for conflict error handling |
| `fields` | list | Model fields |

### Types

**Primitives:** `bool`, `string`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`, `time`, `uuid`, `decimal`, `json_raw`

**Modifiers:** `*Type` (pointer), `[]Type` (slice), `map[KeyType]ValueType` (map)

**Custom:** Reference enums, models, or one-ofs by name.

### Enums

```yaml
types:
  enums:
    - name: Status
      values: [Active, Inactive, Suspended]
```

Generates: type definition, `Validate()`, `String()`, `All*()` list, string-to-enum parser.

### OneOf (Sum Types)

```yaml
types:
  one_ofs:
    - name: PayloadType
      values:
        - { model: TextPayload, id: 1 }
        - { model: ImagePayload, id: 2 }
```

Generates an interface-based polymorphic type with type-safe matching.

### Unique Index Error Handling

Define unique indexes on models to generate per-constraint conflict errors:

```yaml
types:
  models:
    - name: User
      table_name: users
      unique_indexes:
        - constraint_name: users_email_key
          fields: [Email]
        - constraint_name: users_name_org_id_key
          fields: [Name, OrgID]
      fields:
        - { name: ID, type: uuid, primary_key: true }
        - { name: Email, type: string }
        - { name: Name, type: string }
        - { name: OrgID, type: uuid }
```

This generates a `ConflictError` type in the service layer and specific error variables per constraint:

```go
// Generated error variables
var ErrUserEmailAlreadyExists = &ConflictError{Entity: "User", Fields: []string{"Email"}}
var ErrUserNameOrgIDAlreadyExists = &ConflictError{Entity: "User", Fields: []string{"Name", "OrgID"}}

// Storage layer automatically maps PostgreSQL constraint violations to these errors
```

The generic `ErrUserAlreadyExists` is still returned for unrecognized constraints.

### Transactional Outbox Events

```yaml
types:
  models:
    - name: SetUserTagCommand
      storage_type: tx_outbox
      message_builder: "pkg/path.BuildSetUserTagCommandMessage"
      fields:
        - { name: ID, type: uuid, primary_key: true }
        - { name: Data, type: SetUserTagCommandData }
        - { name: IdempotencyKey, type: string }
```

Events are persisted atomically within the same database transaction as your domain operations, then delivered asynchronously via Watermill.

## Cross-Module Event Production

Modules can produce events owned by other modules using the `produces` field:

```yaml
# auth.module.yaml
produces:
  - segmentation.SetUserTagCommand
```

This generates a `SetUserTagCommands()` accessor on the auth module's `Storage` interface, allowing the auth service to send segmentation events within its transactions.

## User Code Blocks

Generated files support preserved code blocks that survive regeneration:

```go
// user code 'imports'
import "my/custom/package"
// end user code 'imports'
```

Any code between these markers is preserved when re-running the generator.

## CLI Usage

```
boilerplate-go [flags] <config.yaml>

Flags:
  -dotenv string    Path to .env file
  -module string    Generate only a specific module
  -version          Show version
```

## Library Packages

The `lib/` directory provides runtime packages used by generated code:

| Package | Description |
|---|---|
| `lib/dbutil` | Generic `EntityStorage` interface and GORM implementation |
| `lib/filter` | Composable typed query filters |
| `lib/pagination` | Generic `Pagination` and `PaginatedResult[T]` types |
| `lib/order` | Generic `Order[FieldType]` with direction support |
| `lib/idempotency` | Duplicate message detection with PostgreSQL storage |
| `lib/txoutbox` | Transactional outbox pattern with message processing |
| `lib/txoutbox/txoutboxwatermill` | Watermill bridge for outbox message delivery |

## Architecture

The generator follows hexagonal architecture principles:

```
Module/
├── moduleservice/    # Domain layer (public API)
│   ├── Models        # Domain types, enums
│   ├── Interfaces    # Storage contracts
│   └── Errors        # Domain errors
└── modulestorage/    # Infrastructure layer
    ├── GORM models   # Internal DB representations
    ├── Converters    # Domain <-> DB model conversion
    ├── Filters       # Query expression builders
    └── Storage impl  # Interface implementations
```

## Example

See the [`example/`](example/) directory for a complete two-module application demonstrating:

- User CRUD with typed filters
- Cross-module event production (auth -> segmentation)
- Transactional outbox with custom message builders
- Service/storage layer separation

## Schema Support

JSON Schema files are provided for IDE autocompletion and validation of configuration files.

Add this comment to the top of your YAML files:

**gen.yaml:**
```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/saturn4er/boilerplate-go/master/gen.schema.json
```

**\*.module.yaml:**
```yaml
# yaml-language-server: $schema=https://raw.githubusercontent.com/saturn4er/boilerplate-go/master/module.schema.json
```

This works with VS Code (via the [YAML extension](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml)) and JetBrains IDEs.

## Requirements

- Go 1.21+
- PostgreSQL (for generated storage implementations)
- GORM v2

## License

See [LICENSE](LICENSE) for details.