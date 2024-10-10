{
"file_path": "{{.Module}}/{{.Module}}service/gen.errors.go",
"package_name": "{{.Module}}service",
"package_path": "{{.Config.RootPackageName}}/{{.Module}}/{{.Module}}service",
"condition": "len(Config.Modules[Module].Value.Types.Models) > 0"
}
<><><>
{{- $module := (index $.Config.Modules $.Module).Value }}
{{- $fmtPkg := import "fmt" }}

type NotFoundError string
func (n NotFoundError) Error() string {
  return {{$fmtPkg.Ref "Sprintf"}}("%s not found", string(n))
}

type AlreadyExistsError string
func (a AlreadyExistsError) Error() string {
  return {{$fmtPkg.Ref "Sprintf"}}("%s already exists", string(a))
}



{{- range $model := $module.Types.Models}}
const(
  Err{{ $model.Name }}NotFound = NotFoundError("{{ $model.Name }}")
  Err{{ $model.Name }}AlreadyExists = AlreadyExistsError("{{ $model.Name }}")
)
{{- end }}
