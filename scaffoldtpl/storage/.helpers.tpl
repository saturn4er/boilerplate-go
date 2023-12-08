{{- define "storage.storage.func.db_enum_to_enum" -}}
  convert{{.}}ToDB{{.}}
{{- end -}}

{{- define "storage.func.db_enum_to_enum" -}}
  convertDB{{.}}To{{.}}
{{- end -}}


{{- define "storage.func.model_to_db_model" -}}
  convert{{.Name}}ToDB{{.Name}}
{{- end -}}

{{- define "storage.func.db_model_to_model" -}}
  convertDB{{.}}To{{.}}
{{- end -}}

{{- define "storage.type.db_model" -}}
  db{{.Name | camelCase}}
{{- end -}}

{{- define "storage.const.model_table_field" -}}
    {{.model.Name | lCamelCase}}Field{{.field.Name|camelCase}}
{{- end -}}

{{- define "storage.const.enum_value" -}}
    {{.enum.Name | lCamelCase}}{{.value | camelCase}}
{{- end -}}

{{- define "storage.const.model_table_name" -}}
    {{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold/config.Model*/ -}}
    {{.Name | lCamelCase}}TableName
{{- end -}}
{{- define "storage._convert_value_to_db_value"}}
    {{- $jsonPkg := import "encoding/json" -}}
    {{- $fmtPkg := import "fmt" -}}
    {{- $input := index . 0 -}}
    {{- $inputGoType := index . 1 -}}
    {{- $output := index . 2 -}}
    {{- $outputGoType := index . 3 -}}
    {{- $varNamesGenerator := index . 4 -}}
    {{- $tmpVar := $varNamesGenerator.Var "tmp" -}}
    {{- if isModuleOneOf $inputGoType }}
        {{- $oneOf := getModuleOneOf $inputGoType.Type }}
        {
        if {{$input}} == nil {
        {{$output}} = nil
        } else {
        {{$tmpVar}} := struct {
        Value {{$inputGoType.Ref}} `json:"value"`
        OneOfType   string `json:"@type"`
        OneOfTypeID uint `json:"@type_id"`
        }{
        Value: {{$input}},
        }
        switch {{$input}}.(type) {
        {{- range $value := $oneOf.SortedValues }}
          case *{{ (goType $value.Value).Ref }}:
          {{$tmpVar}}.OneOfType = "{{ (goType $value.Value).Ref }}"
          {{$tmpVar}}.OneOfTypeID = {{$value.Index}}
        {{- end }}
        default:
        return nil, {{$fmtPkg.Ref "Errorf"}}("invalid {{$oneOf.Name}} value type: %T", {{$input}})
        }
        {{$dataVar := $varNamesGenerator.Var "tmp"}}
        {{$dataVar}}, err := {{$jsonPkg.Ref "Marshal"}}({{$tmpVar}})
        if err != nil {
        return nil, err
        }
        {{- $tmpVar1 := $varNamesGenerator.Var "tmp" }}
        {{$tmpVar1}} := string({{$dataVar}})
        {{$output}} = &{{$tmpVar1}}
        }
        }
    {{- else if or (isModuleEnum $inputGoType) (isCommonEnum $inputGoType) }}
      {
      {{$tmpVar}}, err := {{- template "storage.storage.func.db_enum_to_enum" $inputGoType.Type }}({{$input}})
      if err != nil {
      return nil, err
      }
      {{$output}} = {{$tmpVar}}
      }
    {{- else if or (isModuleModel $inputGoType) (isCommonModel $inputGoType) (eq $inputGoType.Type "any") }}
      {
      {{$tmpVar}}, err := {{$jsonPkg.Ref "Marshal"}}({{$input}})
      if err !=nil{
      return nil, err
      }
      {{$output}} = string({{$tmpVar}})
      }
    {{- else if $inputGoType.IsPtr }}
      if {{$input}} != nil {
      var {{$tmpVar}} {{$outputGoType.ElemType.Ref}}
      {{- template "storage._convert_value_to_db_value" list (print "*" $input) $inputGoType.ElemType $tmpVar $outputGoType.ElemType $varNamesGenerator}}
      {{$output}} = &{{$tmpVar}}
      } else {
      {{$output}} = nil
      }
    {{- else if $inputGoType.IsSlice}}
      {
      {{$tmpVar}} := make({{$outputGoType.Ref}}, 0, len({{$input}}))
      for _, el := range {{$input}} {
      var res {{$outputGoType.ElemType.Ref}}
      {{- template "storage._convert_value_to_db_value" list "el" $inputGoType.ElemType "res" $outputGoType.ElemType $varNamesGenerator}}
      {{$tmpVar}} = append({{$tmpVar}}, res)
      }
      {{$output}} = {{$tmpVar}}
      }
    {{- else if and (eq $inputGoType.Type "Time") (eq $inputGoType.Package "time") (not (index $inputGoType.Metadata "with_timezone"))}}
        {{$output}} = ({{$input}}).UTC() // {{$inputGoType.Metadata}}
    {{- else}}
        {{$output}} = {{$input}}
    {{- end }}
{{- end }}

{{- define "storage._convert_db_value_to_value"}}
    {{- $jsonPkg := import "encoding/json" }}
    {{- $fmtPkg := import "fmt" }}
    {{- $input := index . 0 -}}
    {{- $inputGoType := index . 1 -}}
    {{- $output := index . 2 -}}
    {{- $outputGoType := index . 3 -}}
    {{- $varNamesGenerator := index . 4 -}}
    {{- $tmpVar := $varNamesGenerator.Var "tmp" -}}
    {{- if isModuleOneOf $outputGoType }}
        {{- $oneOf := getModuleOneOf $outputGoType.Type }}
        {
        if {{$input}} != nil {
        {{- $typeInfoVar := $varNamesGenerator.Var "typeInfo" -}}
        {{$typeInfoVar}} := struct { OneOfTypeID uint `json:"@type_id"` }{}
        if err := {{$jsonPkg.Ref "Unmarshal"}}([]byte(*{{$input}}), &{{$typeInfoVar}}); err!=nil {
        return nil, err
        }
        {{- $valueVar := $varNamesGenerator.Var "value" }}
        switch {{$typeInfoVar}}.OneOfTypeID {
        {{- range $value := $oneOf.SortedValues }}
          case {{$value.Index}}:
          var {{$valueVar}} struct{
          Value {{(goType $value.Value.Model).Ref}} `json:"value"`
          }
          if err := {{$jsonPkg.Ref "Unmarshal"}}([]byte(*{{$input}}), &{{$valueVar}}); err!=nil {
          return nil, err
          }
          {{$output}} = &{{$valueVar}}.Value
        {{- end }}
        default:
        return nil, {{$fmtPkg.Ref "Errorf"}}("invalid {{$oneOf.Name}} value type id: %d", {{$typeInfoVar}}.OneOfTypeID)
        }
        }
        }
    {{- else if or (isModuleEnum $outputGoType) (isCommonEnum $outputGoType) }}
      {
      {{$tmpVar}}, err := {{- template "storage.func.db_enum_to_enum" $outputGoType.Type }}({{$input}})
      if err != nil {
      return nil, err
      }
      {{$output}} = {{$tmpVar}}
      }
    {{- else if or (isModuleModel $outputGoType) (isCommonModel $outputGoType) (eq $outputGoType.Type "any")}}
      {
      var {{$tmpVar}} {{$outputGoType.Ref}}
      if err := {{$jsonPkg.Ref "Unmarshal"}}([]byte({{$input}}), &{{$tmpVar}}); err !=nil{
      return nil, err
      }
      {{$output}} = {{$tmpVar}}
      }
    {{- else if $outputGoType.IsPtr }}
      if {{$input}} != nil {
      var {{$tmpVar}} {{$outputGoType.ElemType.Ref}}
      {{- template "storage._convert_db_value_to_value" list (print "*" $input) $inputGoType.ElemType $tmpVar $outputGoType.ElemType $varNamesGenerator}}
      {{$output}} = &{{$tmpVar}}
      }
    {{- else if $outputGoType.IsSlice}}
      {
      {{$tmpVar}} := make({{$outputGoType.Ref}}, 0, len({{$input}}))
      for _, el := range {{$input}} {
      var res {{$outputGoType.ElemType.Ref}}
      {{- template "storage._convert_db_value_to_value" list "el" $inputGoType.ElemType "res" $outputGoType.ElemType $varNamesGenerator}}
      {{$tmpVar}} = append({{$tmpVar}}, res)
      }
      {{$output}} = {{$tmpVar}}
      }
    {{- else}}
        {{$output}} = {{$input}}
    {{- end }}
{{- end }}

