{{- define "service.type.enum" -}}
    {{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold/config.Enum*/ -}}
    {{.Name}}
{{- end -}}

{{- define "service.type.model" -}}
    {{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold/config.Model*/ -}}
    {{.Name}}
{{- end -}}

{{- define "service.type.model_field" -}}
  {{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold/config.Model*/ -}}
  {{.Name }}Field
{{- end -}}

{{- define "service.type.model_filter" -}}
    {{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold/config.Model*/ -}}
    {{.Name }}Filter
{{- end }}

{{- define "service.const.model_field" -}}
    {{- template "service.type.model_field" .model -}}{{.field.Name}}
{{- end -}}

{{- define "service.const.enum_value" -}}
    {{- template "service.type.enum" .enum -}}{{.value}}
{{- end -}}

{{- define "service._copy_value" }}
    {{- $input := index . 0 -}}
    {{- $goType := index . 1 -}}
    {{- $output := index . 2 -}}
    {{- $varNamesGenerator := index . 3 -}}
    {{- if isModuleOneOf $goType}}
      {{- $oneOf := getModuleOneOf $goType.Type }}
      if {{$input}} != nil {
        switch val := {{$input}}.(type) {
        {{- range $value := $oneOf.SortedValues}}
          case *{{$value.Value.ModelName}}:
          valCopy := val.Copy()
          {{$output}} = &valCopy
        {{- end}}
        }
      }
    {{- else if or (isModuleModel $goType) (isCommonModel $goType) }}
      {{$output}} = {{$input}}.Copy() // model
    {{- else if or (isModuleEnum $goType) (isCommonEnum $goType)}}
      {{$output}} =  {{$input}}// enum
    {{- else if $goType.IsPtr }}
      if {{$input}} != nil{
        {{- $tmpVar := $varNamesGenerator.Var "tmp" -}}
        var {{$tmpVar}} {{$goType.ElemType.Ref}}
        {{- template "service._copy_value" (list (print "(*" $input ")") $goType.ElemType $tmpVar $varNamesGenerator) }}
        {{$output}} = &{{$tmpVar}}
      }
    {{- else if $goType.IsSlice}}
      {{- $tmpVar := $varNamesGenerator.Var "tmp"}}
      {{- $itemVar := $varNamesGenerator.Var "i" }}
      {{$tmpVar}} := make({{$goType.Ref}}, 0, len({{$input}}))
      for _, {{$itemVar}} := range {{$input}} {
      {{- $itemCopyVar := $varNamesGenerator.Var "itemCopy" -}}
        var {{$itemCopyVar}} {{$goType.ElemType.Ref}}
        {{- template "service._copy_value" (list $itemVar $goType.ElemType $itemCopyVar $varNamesGenerator) }}
        {{$tmpVar}} = append({{$tmpVar}}, {{$itemCopyVar}})
      }
      {{$output}} = {{$tmpVar}}
    {{- else }}
      {{$output}} = {{$input}}
    {{- end }}
{{- end }}
