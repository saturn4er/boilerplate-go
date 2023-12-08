{
"file_path": "{{.Module}}/{{.Module}}service/gen.models.go",
"package_name": "{{.Module}}service",
"package_path": "{{.Config.RootPackageName}}/{{.Module}}/{{.Module}}service",
"condition": "len(Config.Modules[Module].Value.Types.Models) > 0"
}
<><><>
{{- $module := (index $.Config.Modules $.Module).Value }}
{{- $fmtPkg := import "fmt" }}

{{- range $oneOf := $module.Types.OneOfs}}
    type {{$oneOf.Name}} interface{
      is{{$oneOf.Name}}()
    }
    {{- range $value := $oneOf.SortedValues}}
    func (*{{$value.Value.ModelName}}) is{{$oneOf.Name}}() {}
    {{- end}}

    func copy{{$oneOf.Name}}(val {{$oneOf.Name}}) {{$oneOf.Name}} {
      if val== nil {
        return nil
      }

      switch val := val.(type) {
        {{- range $value := $oneOf.SortedValues}}
          case *{{$value.Value.ModelName}}:
          valCopy := val.Copy()
          return &valCopy
        {{- end}}
      }
      panic("called copy{{$oneOf.Name}} with invalid type")e
    }
{{- end }}
{{- range $model := $module.Types.Models }}
  {{- if not $model.DoNotPersists }}
    type {{$model.Name}}Field byte
    const (
    {{$model.Name}}Field{{(index $model.Fields 0).Name}} {{$model.Name}}Field = iota + 1
    {{- range $field := (slice $model.Fields 1) }}
          {{$model.Name}}Field{{$field.Name}}
    {{- end }}
    )
    type {{$model.Name}}Filter struct {
    {{- range $field := $model.Fields }}
        {{- if $field.Filterable }}
            {{$field.Name}} {{ template "filterType" (goType $field.Type)}}
        {{- end }}
    {{- end }}
        Or []*{{$model.Name}}Filter
        And []*{{$model.Name}}Filter
    }
  {{- end}}

  type {{$model.Name}} struct {
  {{- range $field := $model.Fields }}
      {{$field.Name}} {{(goType $field.Type).Ref}}
  {{- end }}
  }
  {{- $receiverName := slice $model.Name 0 1 | lCamelCase}}
  func ({{$receiverName}} {{$model.Name}}) Copy() {{$model.Name}} {
    var result {{$model.Name}}
    {{- $varNamesGenerator := varNamesGenerator }}
    {{- range $field := $model.Fields }}
        {{- $fieldGoType := goType $field.Type -}}
        {{- $input := print $receiverName "." $field.Name}}
        {{- $output := print "result." $field.Name}}
        {{- template "copy_value" (list $input $fieldGoType $output $varNamesGenerator) }}
    {{- end }}

    return result
  }
{{- end }}
