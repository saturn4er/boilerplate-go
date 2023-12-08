{
"file_path": "{{.Module}}/{{.Module}}/gen.errors.go",
"package_name": "{{.Module}}",
"package_path": "{{.Config.RootPackageName}}/{{.Module}}/{{.Module}}",
"condition": "len(Config.Modules[Module].Value.Types.Models) > 0"
}
<><><>
{{- $module := (index $.Config.Modules $.Module).Value -}}
{{- $fmtPkg := import "fmt" -}}

type NotFoundError string
func (e NotFoundError) Error() string {
    return {{$fmtPkg.Ref "Sprintf"}}("%s not found", string(e))
}

type AlreadyExistsError string
func (e AlreadyExistsError) Error() string {
    return {{$fmtPkg.Ref "Sprintf"}}("%s already exists", string(e))
}

var (
{{- range $model := $module.Types.Models }}
    Err{{ $model.Name }}NotFound = NotFoundError("{{ $model.Name }}")
    Err{{ $model.Name }}AlreadyExists = AlreadyExistsError("{{ $model.Name }}")
    {{ userBlock (printf "%s related errors" $model.Name) }}

{{- end }}
)
