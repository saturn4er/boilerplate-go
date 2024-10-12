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
{{ $errorsPkg := import "github.com/pkg/errors" }}
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
          func New{{$model.PluralName}}StorageBase(db *{{$gormPkg.Ref "DB"}}, logger *{{$loggingPkg.Ref "Logger"}}) {{$servicePkg.Ref (print $model.PluralName "StorageBase")}} {
        {{- else }}
          func New{{$model.PluralName}}Storage(db *{{$gormPkg.Ref "DB"}}, logger *{{$loggingPkg.Ref "Logger"}}) {{$servicePkg.Ref (print $model.PluralName "Storage")}} {
        {{- end }}
        return {{$dbutilPkg.Ref "GormEntityStorage"}}[{{$servicePkg.Ref $model.Name}}, {{$dbTypeRef}}, {{$servicePkg.Ref (print $model.Name "Filter")}}]{
        Logger: logger,
        DB: db,
        DBErrorsWrapper:       {{template "storage.func.errors_wrapper" $model.Name}},
        ConvertToInternal:     {{template "storage.func.table_model_to_internal" $model.Name}},
        ConvertToExternal:     {{template "storage.func.table_model_to_service" $model.Name}},
        BuildFilterExpression: {{template "storage.func.build_db_filter" $model.Name}},
        FieldMapping:          map[any]clause.Column{
        {{- range $field := $model.Fields }}
            {{$servicePkg.Ref (print $model.Name "Field" $field.Name)}}: {Name: "{{$field.DBName}}"},
        {{- end}}
        },
        }
        }
    {{- end }}
{{- end }}




