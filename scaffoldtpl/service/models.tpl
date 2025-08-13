{
"file_path": "{{.Module}}/{{.Module}}service/gen.models.go",
"package_name": "{{.Module}}service",
"package_path": "{{.Config.RootPackageName}}/{{.Module}}/{{.Module}}service",
"condition": "len(Config.Modules[Module].Value.Types.Models) > 0"
}
<><><>
{{- $module := (index $.Config.Modules $.Module).Value }}
{{- $fmtPkg := import "fmt" }}

type OrderDirection = order.Direction

{{- range $oneOf := $module.Types.OneOfs}}
  type {{$oneOf.Name}} interface{
  is{{$oneOf.Name}}()
  {{$oneOf.Name}}Equals({{$oneOf.Name}}) bool
  {{ userCodeBlock (printf "%s methods" $oneOf.Name) }}
  }
  {{- range $value := $oneOf.SortedValues}}
      {{- $receiverName := $value.Value.ModelName | receiverName }}
      func (*{{$value.Value.ModelName}}) is{{$oneOf.Name}}() {}
      func ({{$receiverName}} *{{$value.Value.ModelName}}) {{$oneOf.Name}}Equals(to {{$oneOf.Name}}) bool {
      if ({{$receiverName}} == nil) != (to == nil) {
      return false
      }
      if {{$receiverName}} == nil && to == nil {
      return true
      }

      toTyped, ok := to.(*{{$value.Value.ModelName}})
      if !ok {
      return false
      }

      return {{$receiverName}}.Equals(toTyped)
      }
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
  panic("called copy{{$oneOf.Name}} with invalid type")
  }
{{- end }}
{{- range $model := $module.Types.Models }}
  {{- if len $model.Fields}}
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
    type {{$model.Name}}Order order.Order[{{$model.Name}}Field]
  {{- end }}

    type {{$model.Name}}{{if gt (len $model.TypeParameters) 0}}[{{range $model.TypeParameters}} {{.Name}} {{.Constraint}}{{end}}]{{end}} struct {
    {{- range $field := $model.Fields }}
        {{$field.Name}} {{(goType $field.Type).Ref}}
    {{- end }}
    }
    {{ userCodeBlock (printf "%s methods" $model.Name) }}
    {{- $receiverName := slice $model.Name 0 1 | lCamelCase}}
    {{$methodDefinitionTypeParams := ""}}
    {{- if gt (len $model.TypeParameters) 0 }}
      {{- $methodDefinitionTypeParams = "[" }}
      {{- range $param := $model.TypeParameters }}
        {{- $methodDefinitionTypeParams = print $methodDefinitionTypeParams $param.Name ", "}}
      {{- end }}
      {{- $methodDefinitionTypeParams = print $methodDefinitionTypeParams "]"}}
    {{- end }}
    {{ $methodDefinition := print "func (" $receiverName " *" $model.Name $methodDefinitionTypeParams ")" }}
    {{- if gt (len $model.TypeParameters) 0 }}
      {{- continue }}
    {{- end }}
    {{ $methodDefinition}} Copy() {{$model.Name}} {
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
    func ({{$receiverName}} *{{$model.Name}}) Equals(to *{{$model.Name}}) bool {
    if ({{$receiverName}} == nil) != (to == nil) {
    return false
    }
    if {{$receiverName}} == nil && to == nil {
    return true
    }

    {{- $varNamesGenerator = varNamesGenerator }}
    {{- range $field := $model.Fields }}
        {{- $fieldGoType := goType $field.Type -}}
        {{- $aVal := print $receiverName "." $field.Name}}
        {{- $bVal := print "to." $field.Name}}
        {{- template "equals_value" (list $aVal $fieldGoType $bVal $varNamesGenerator) }}
    {{- end }}

    return true
    }
{{- end }}
