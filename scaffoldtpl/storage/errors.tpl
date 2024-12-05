{
"file_path": "{{.Module}}/{{.Module}}storage/gen.errors.go",
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
{{ $contextPkg := import "context" }}
{{ $servicePkg :=  import (print $.Config.RootPackageName "/" $.Module "/" $.Module "service") (print $.Module "svc") }}
{{ $module := (index $.Config.Modules $.Module).Value}}

{{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold.ConfigModel*/ -}}
{{- range $model := $module.Types.Models }}
    {{- if $model.DoNotPersists }}
        {{continue}}
    {{- end }}
    func wrap{{$model.Name}}QueryError(err error) error {
      if err == nil{
        return nil
      }

      if {{$errorsPkg.Ref "Is"}}(err, {{$gormPkg.Ref "ErrRecordNotFound"}}) {
        return {{$servicePkg.Ref (printf "Err%sNotFound" $model.Name)}}
      }

      var pgErr *{{$pgconnPkg.Ref "PgError"}}

      if {{$errorsPkg.Ref "As"}}(err, &pgErr) {
        if pgErr.Code == "23505" {
          return {{$servicePkg.Ref (printf "Err%sAlreadyExists" $model.Name)}}
        }
      }

      return err
    }
{{- end }}
