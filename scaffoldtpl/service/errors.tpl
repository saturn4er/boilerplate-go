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

{{- $hasUniqueIndexes := false }}
{{- range $model := $module.Types.Models }}
  {{- if $model.UniqueIndexes }}
    {{- $hasUniqueIndexes = true }}
  {{- end }}
{{- end }}

{{- if $hasUniqueIndexes }}
{{- $stringsPkg := import "strings" }}

type ConflictError struct {
  Entity string
  Fields []string
}

func (e *ConflictError) Error() string {
  return {{$fmtPkg.Ref "Sprintf"}}("%s conflict on %s", e.Entity, {{$stringsPkg.Ref "Join"}}(e.Fields, ", "))
}

func (e *ConflictError) Is(target error) bool {
  if ae, ok := target.(AlreadyExistsError); ok {
    return string(ae) == e.Entity
  }
  return false
}
{{- end }}

{{- range $model := $module.Types.Models}}
const(
  Err{{ $model.Name }}NotFound = NotFoundError("{{ $model.Name }}")
  Err{{ $model.Name }}AlreadyExists = AlreadyExistsError("{{ $model.Name }}")
)
{{- if $model.UniqueIndexes }}
var (
  {{- range $idx := $model.UniqueIndexes }}
  Err{{ $model.Name }}{{ $idx.ErrorName }}AlreadyExists = &ConflictError{Entity: "{{ $model.Name }}", Fields: []string{ {{- range $i, $f := $idx.Fields }}{{if $i}}, {{end}}"{{$f}}"{{- end -}} }}
  {{- end }}
)
{{- end }}
{{- end }}