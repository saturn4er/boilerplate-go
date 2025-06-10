{
"file_path": "{{.Module}}/{{.Module}}service/gen.storages.go",
"package_name": "{{.Module}}service",
"package_path": "{{.Config.RootPackageName}}/{{.Module}}/{{.Module}}",
"condition": "len(Config.Modules[Module].Value.Types.Models) > 0"
}
<><><>
{{ $dbutil := import "github.com/saturn4er/boilerplate-go/lib/dbutil" }}
{{ $idempotencyPkg := import "github.com/saturn4er/boilerplate-go/lib/idempotency" }}
{{ $txoutboxPkg := import "github.com/saturn4er/boilerplate-go/lib/txoutbox" }}
{{ $contextPkg := import "context" }}
{{ $module := (index $.Config.Modules $.Module).Value}}

type Storage interface {
{{- range $index, $model := $module.Types.Models }}
    {{- if $model.DoNotPersists }}
        {{- continue }}
    {{- end }}
    {{- if eq $model.StorageType "tx_outbox" }}
        {{$model.PluralName}}() {{$model.PluralName}}Outbox
    {{- else }}
        {{$model.PluralName}}() {{$model.PluralName}}Storage
    {{- end }}
{{- end }}
IdempotencyKeys() {{$idempotencyPkg.Ref "Storage"}}
ExecuteInTransaction(ctx {{$contextPkg.Ref "Context"}}, cb func(ctx {{$contextPkg.Ref "Context"}}, tx Storage) error) error
WithAdvisoryLock(ctx {{$contextPkg.Ref "Context"}}, scope string, lockID any) error 
}

{{- range $model := $module.Types.Models }}
    {{- if not $model.DoNotPersists }}
        {{- if eq $model.StorageType "tx_outbox" }}
          type {{$model.PluralName}}Outbox {{$txoutboxPkg.Ref "Outbox"}}[{{$model.Name}}]
        {{- else }}
            {{- if $model.HasCustomDBMethods }}
              type {{$model.PluralName}}Storage interface {
                {{$dbutil.Ref "EntityStorage"}}[{{$model.Name}}, {{$model.Name}}Filter]
                {{ userCodeBlock (printf "%s metods" $model.Name) }}
              }
              {{ userCodeBlock (printf "%s definitions" $model.Name) }}
            {{- else }}
              type {{$model.PluralName}}Storage {{$dbutil.Ref "EntityStorage"}}[{{$model.Name}}, {{$model.Name}}Filter]
            {{- end }}
        {{- end }}
    {{- end }}
{{- end }}
