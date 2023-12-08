{{- define "modelType" }}
    {{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold/config.Model */}}
    type {{.Name}} struct {
    {{- range $field := .Fields }}
        {{$field.Name}} {{(goType $field.Type).Ref}}
    {{- end }}
    }
{{- end }}

{{- define "service._model_fields_enum" }}
    {{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold/config.Model */}}
    type {{template "service.type.model_field" .}} byte
    const (
    {{- template "service.const.model_field" (dict "model" $ "field" (first .Fields))}} {{template "service.type.model_field" $}} = iota + 1
    {{- range $field := rest .Fields  }}
        {{template "service.const.model_field" (dict "model" $ "field" $field)}}
    {{- end }}
    )
{{- end }}


{{- define "service._filter_type" }}
    {{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold.GoType*/}}
    {{- $filterPkg := import "github.com/saturn4er/boilerplate-go/lib/filter"}}
    {{- if .IsSlice -}}
        {{$filterPkg.Ref "ArrayFilter"}}[{{.ElemType.Ref}}]
    {{- else -}}
        {{$filterPkg.Ref "Filter"}}[{{.Ref}}]
    {{- end -}}
{{- end }}

{{- define "modelFilterType" }}
  {{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold/config.Model */}}
  type {{ template "service.type.model_filter" . }} struct {
  {{- range $field := .Fields }}
      {{- if $field.Filterable }}
          {{$field.Name}} {{ template "service._filter_type" (goType $field.Type)}}
      {{- end }}
  {{- end }}
  Or []*{{ template "service.type.model_filter" $ }}
  And []*{{ template "service.type.model_filter" $ }}
  }
{{- end }}

{{- define "modelCopyHelper" }}
    {{- $typeName := (include "service.type.model" .) }}
    {{- $receiverName := $typeName | receiverName }}
    func ({{$receiverName}} {{$typeName}}) Copy() {{$typeName}} {
    var result {{$typeName}}
    {{- $varNamesGenerator := varNamesGenerator }}
    {{- range $field := .Fields }}
        {{- $fieldGoType := goType $field.Type -}}
        {{- $input := print $receiverName "." $field.Name}}
        {{- $output := print "result." $field.Name}}
        {{- template "service._copy_value" (list $input $fieldGoType $output $varNamesGenerator) }}
    {{- end }}

    return result
    }
{{- end }}
