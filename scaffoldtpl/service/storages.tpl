{
"file_path": "{{.Module}}/{{.Module}}/gen.storages.go",
"package_name": "{{.Module}}",
"package_path": "{{.Config.RootPackageName}}/{{.Module}}/{{.Module}}",
"condition": "len(Config.Modules[Module].Value.Types.Models) > 0"
}
<><><>
{{ $dbutilPkg := import "github.com/saturn4er/boilerplate-go/lib/dbutil" }}
{{ $optionutilPkg := import "github.com/go-pnp/go-pnp/pkg/optionutil" }}
{{ $contextPkg := import "context" }}
{{ $module := (index $.Config.Modules $.Module).Value}}

{{- $optionsType := printf "%s[%s]" ($optionutilPkg.Ref "Option") ($dbutilPkg.Ref "SelectOptions") }}
type EntityStorage[Entity, Filter any] interface {
Create(ctx {{$contextPkg.Ref "Context"}}, model *Entity) (*Entity, error)
Count(ctx {{$contextPkg.Ref "Context"}}, filter *Filter) (int, error)
Update(ctx {{$contextPkg.Ref "Context"}}, model *Entity) (*Entity, error)
Save(ctx {{$contextPkg.Ref "Context"}}, model *Entity) (*Entity, error)
First(ctx {{$contextPkg.Ref "Context"}}, filter *Filter, options ...{{$optionsType}}) (*Entity, error)
FirstOrCreate(ctx {{$contextPkg.Ref "Context"}}, filter *Filter, model *Entity, options ...{{$optionsType}}) (*Entity, error)
Find(ctx {{$contextPkg.Ref "Context"}}, filter *Filter, options ...{{$optionsType}}) ([]*Entity, error)
Delete(ctx {{$contextPkg.Ref "Context"}}, filter *Filter) error
}
{{/*type OutboxStorage[Entity any] interface {*/}}
{{/*Send(ctx {{$contextPkg.Ref "Context"}}, model *Entity) error*/}}
{{/*}*/}}

{{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold/config.Module */}}
{{ $idempotencyPkg := import "github.com/saturn4er/boilerplate-go/lib/idempotency" }}

type Storage interface {
{{- range $index, $model := $module.Types.Models }}
    {{- if $model.DoNotPersists }}
        {{- continue }}
    {{- end }}
    {{- if eq $model.StorageType "tx_outbox" }}
{{/*        {{$model.PluralName}}() {{$model.PluralName}}Outbox*/}}
    {{- else }}
        {{$model.PluralName}}() {{$model.PluralName}}Storage
    {{- end }}
{{- end }}
{{/*IdempotencyKeys() {{ $idempotencyPkg.Ref "Storage" }}*/}}
ExecuteInTransaction(ctx {{$contextPkg.Ref "Context"}}, cb func(ctx {{$contextPkg.Ref "Context"}}, tx Storage) error, txOptions ...{{$dbutilPkg.Ref "TxOption"}}) error
}


{{- range $model := $module.Types.Models }}
    {{- if not $model.DoNotPersists }}
        {{- if eq $model.StorageType "tx_outbox" }}
          type {{$model.PluralName}}Outbox OutboxStorage[{{$model.Name}}]
        {{- else }}
              type {{$model.PluralName}}Storage interface {
                EntityStorage[{{$model.Name}}, {{$model.Name}}Filter]
                {{ userBlock (cat $model.Name "Storage") }}
              }
        {{- end }}
    {{- end }}
{{- end }}
