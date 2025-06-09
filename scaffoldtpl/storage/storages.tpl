{
"file_path": "{{.Module}}/{{.Module}}storage/gen.storages.go",
"package_name": "{{.Module}}storage",
"package_path": "{{.Config.RootPackageName}}/{{.Module}}/{{.Module}}storage",
"condition": "len(Config.Modules[Module].Value.Types.Models) > 0"
}
<><><>

{{ $dbutilPkg := import "github.com/saturn4er/boilerplate-go/lib/dbutil" }}
{{ $loggingPkg := import "github.com/go-pnp/go-pnp/logging" }}
{{ $txoutboxPkg := import "github.com/saturn4er/boilerplate-go/lib/txoutbox" }}
{{ $idempotencyPkg := import "github.com/saturn4er/boilerplate-go/lib/idempotency" }}
{{ $gormPkg := import "gorm.io/gorm" }}
{{ $clausePkg := import "gorm.io/gorm/clause" }}
{{ $errorsPkg := import "github.com/pkg/errors" }}
{{ $pgconnPkg := import "github.com/jackc/pgx/v5/pgconn"}}
{{ $xxhashPkg := import "github.com/cespare/xxhash"}}
{{ $contextPkg := import "context" }}
{{ $servicePkg :=  import (print $.Config.RootPackageName "/" $.Module "/" $.Module "service") (print $.Module "svc") }}
{{ $module := (index $.Config.Modules $.Module).Value}}

type Storages struct {
db *{{$gormPkg.Ref "DB"}}
logger *{{$loggingPkg.Ref "Logger"}}
}

var _ {{$servicePkg.Ref "Storage"}} = &Storages{}

{{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold.ConfigModel*/ -}}
{{- range $model := $module.Types.Models }}
    {{- if $model.DoNotPersists }}
        {{continue}}
    {{- end }}
    {{- if eq $model.StorageType "tx_outbox" }}
    func (s Storages) {{$model.PluralName}}() {{$servicePkg.Ref (print $model.PluralName "Outbox")}} {
      return New{{$model.PluralName}}Outbox(s.db)
    }
    {{- else }}
    func (s Storages) {{$model.PluralName}}() {{$servicePkg.Ref (print $model.PluralName "Storage")}} {
      return New{{$model.PluralName}}Storage(s.db, s.logger)
    }
    {{- end }}
{{- end }}

func (s Storages) IdempotencyKeys() {{$idempotencyPkg.Ref "Storage"}} {
  return {{$idempotencyPkg.Ref "GormStorage"}}{
    DB: s.db,
  }
}

{{- range $model := $module.Types.Models }}
  {{- if $model.AdvisoryLock }}
    func (s Storages) {{$model.PluralName}}AdvisoryLock(ctx {{$contextPkg.Ref "Context"}}, lockID {{$model.AdvisoryLockType}}) error {
      offset :=  {{$xxhashPkg.Ref "Sum64String"}}("{{$model.PluralName}}")
      {{- if eq $model.AdvisoryLockType "uuid" }}
        result := s.db.WithContext(ctx).Exec("SELECT pg_advisory_xact_lock(?, ?)", offset, {{$xxhashPkg.Ref "Sum64String"}}(lockID.String()))
      {{- else if eq $model.AdvisoryLockType "string" }}
        result := s.db.WithContext(ctx).Exec("SELECT pg_advisory_xact_lock(?, ?)", offset, {{$xxhashPkg.Ref "Sum64String"}}(lockID))
      {{- else }}
        result := s.db.WithContext(ctx).Exec("SELECT pg_advisory_xact_lock(?, ?)", offset, lockID)
      {{- end }}
      if result.Error != nil {
        return result.Error
      }
      return nil
    }
  {{- end }}
{{- end }}

func (s Storages) ExecuteInTransaction(ctx {{$contextPkg.Ref "Context"}}, cb func(ctx {{$contextPkg.Ref "Context"}}, tx {{$servicePkg.Ref "Storage"}}) error) error {
return s.db.Transaction(func(tx *gorm.DB) error {
return cb(ctx, &Storages{tx, s.logger})
})
}

func NewStorages(db *{{$gormPkg.Ref "DB"}}, logger *{{$loggingPkg.Ref "Logger"}}) *Storages {
return &Storages{db: db, logger: logger}
}

{{- range $model := $module.Types.Models }}
    {{- if $model.DoNotPersists }}
        {{continue}}
    {{- end }}
    {{- if eq $model.StorageType "tx_outbox" }}
        func New{{$model.PluralName}}Outbox(db *{{$gormPkg.Ref "DB"}}) {{$servicePkg.Ref (print $model.PluralName "Outbox")}} {
          return {{$txoutboxPkg.Ref "GormStorage"}}[{{$servicePkg.Ref $model.Name}}]{
            DB: db,
            BuildMessage:     build{{$model.Name}}Message,
          }
        }
    {{- else }}
        {{ $modelGoType := goType $model }}
        {{ $dbType := $modelGoType.InLocalPackage.WithName (print "db" $modelGoType.Type) }}
        {{ $dbTypeRef := $dbType.Ref }}

        {{- if $model.HasCustomDBMethods }}
          type {{$model.PluralName}}Storage struct {
          {{$dbutilPkg.Ref "GormEntityStorage"}}[{{$servicePkg.Ref $model.Name}}, {{$dbTypeRef}}, {{$servicePkg.Ref (print $model.Name "Filter")}}]
          }
          {{ userCodeBlock (printf "%s custom methods" $model.Name) }}
          func New{{$model.PluralName}}Storage(db *{{$gormPkg.Ref "DB"}}, logger *{{$loggingPkg.Ref "Logger"}}) {{$servicePkg.Ref (print $model.PluralName "Storage")}} {
            return &{{$model.PluralName}}Storage{
              GormEntityStorage: {{$dbutilPkg.Ref "GormEntityStorage"}}[{{$servicePkg.Ref $model.Name}}, {{$dbTypeRef}}, {{$servicePkg.Ref (print $model.Name "Filter")}}]{
                Logger: logger,
                DB: db,
                DBErrorsWrapper:       wrap{{$model.Name}}QueryError,
                ConvertToInternal:     {{template "storage.func.table_model_to_internal" $model.Name}},
                ConvertToExternal:     {{template "storage.func.table_model_to_service" $model.Name}},
                BuildFilterExpression: func(filter *{{$servicePkg.Ref (print $model.Name "Filter")}}) ({{ $clausePkg.Ref "Expression"}}, error) {
                  return {{template "storage.func.build_db_filter" $model.Name}}(filter)
                },
                FieldMapping:          map[any]{{$clausePkg.Ref "Column"}}{
                  {{- range $field := $model.Fields }}
                      {{$servicePkg.Ref (print $model.Name "Field" $field.Name)}}: {Name: "{{$field.DBName}}"},
                  {{- end }}
                },
              },
              {{ userCodeBlock (printf "%s custom metods" $model.Name) }}
            }
          }
        {{- else }}
          func New{{$model.PluralName}}Storage(db *{{$gormPkg.Ref "DB"}}, logger *{{$loggingPkg.Ref "Logger"}}) {{$servicePkg.Ref (print $model.PluralName "Storage")}} {
            return {{$dbutilPkg.Ref "GormEntityStorage"}}[{{$servicePkg.Ref $model.Name}}, {{$dbTypeRef}}, {{$servicePkg.Ref (print $model.Name "Filter")}}]{
              Logger: logger,
              DB: db,
              DBErrorsWrapper:       wrap{{$model.Name}}QueryError,
              ConvertToInternal:     {{template "storage.func.table_model_to_internal" $model.Name}},
              ConvertToExternal:     {{template "storage.func.table_model_to_service" $model.Name}},
              BuildFilterExpression: func(filter *{{$servicePkg.Ref (print $model.Name "Filter")}}) ({{ $clausePkg.Ref "Expression"}}, error) {
                return {{template "storage.func.build_db_filter" $model.Name}}(filter)
              },
              FieldMapping:          map[any]{{$clausePkg.Ref "Column"}}{
                {{- range $field := $model.Fields }}
                    {{$servicePkg.Ref (print $model.Name "Field" $field.Name)}}: {Name: "{{$field.DBName}}"},
                {{- end }}
              },
            }
          }
        {{- end }}
    {{- end }}
{{- end }}




