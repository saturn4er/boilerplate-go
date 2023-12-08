{{- define "storage.func.enum_to_internal" -}}
  convert{{.}}ToDB{{.}}
{{- end -}}

{{- define "storage.func.enum_to_service" -}}
  convertDB{{.}}To{{.}}
{{- end -}}

{{- define "storage.func.one_of_to_internal" -}}
  {{/*gotype: github.com/saturn4er/boilerplate-go/scaffold.ConfigOneOf*/}}
    convert{{.Name}}ToDB{{.Name}}
{{- end -}}

{{- define "storage.func.one_of_to_service" -}}
    {{/*gotype: github.com/saturn4er/boilerplate-go/scaffold.ConfigOneOf*/}}
    convertDB{{.Name}}To{{.Name}}
{{- end -}}

{{- define "storage.func.build_db_filter" -}}
  build{{(print . "Filter")}}Expr
{{- end}}

{{- define "storage.func.model_to_internal" -}}
  convert{{.}}ToDB{{.}}
{{- end -}}
{{- define "storage.func.model_to_service" -}}
  convertDB{{.}}To{{.}}
{{- end -}}

{{- define "storage.func.errors_wrapper" -}}
  wrap{{.}}Error
{{- end -}}


{{- define "storage.block.convert_value_to_internal"}}
    {{- $jsonPkg := import "encoding/json" -}}
    {{- $fmtPkg := import "fmt" -}}
    {{- $input := index . 0 -}}
    {{- $inputGoType := index . 1 -}}
    {{- $output := index . 2 -}}
    {{- $outputGoType := index . 3 -}}
    {{- $varNamesGenerator := index . 4 -}}
    {{- $tmpVar := $varNamesGenerator.Var "tmp" -}}
    {{- if isModuleOneOf $inputGoType }}
        // OneOf
        {{- $oneOf := getModuleOneOf $inputGoType.Type }}
        {
          if {{$input}} == nil {
            {{$output}} = nil
          } else {
            switch v := {{$input}}.(type) {
            {{- range $value := $oneOf.SortedValues }}
              {{- $dbType := (goType $value.Value.Model).InLocalPackage.WithName (print "db" (goType $value.Value.Model).Type) }}
              case *{{ (goType $value.Value).Ref }}:
              {{$valueOutput := print $output ".Val"}}
              {{- template "storage.block.convert_value_to_internal" list "v" (goType $value.Value) $valueOutput $dbType $varNamesGenerator}}

              {{$output}}.OneOfType = "{{ (goType $value.Value).Ref }}"
              {{$output}}.OneOfTypeID = {{$value.Index}}
            {{- end }}
            default:
            return nil, {{$fmtPkg.Ref "Errorf"}}("invalid {{$oneOf.Name}} value type: %T", {{$input}})
            }
          }
        }
    {{- else if or (isModuleEnum $inputGoType) (isCommonEnum $inputGoType) }}
      // Enum
      {
      {{$tmpVar}}, err := {{- template "storage.func.enum_to_internal" $inputGoType.Type }}({{$input}})
      if err != nil {
      return nil, err
      }
      {{$output}} = {{$tmpVar}}
      }
    {{- else if or (isModuleModel $inputGoType) (isCommonModel $inputGoType)}}
      // Model
      {
      {{ $model := getModel $inputGoType }}
      {{$tmpVar}}, err := {{template "storage.func.model_to_internal" $model.Name}}({{$input}})
      if err!=nil {
        return nil, fmt.Errorf("convert {{$model.Name}} to db: %w", err)
      }

      {{$jsonData := $varNamesGenerator.Var "jsonData"}}

      {{$jsonData}}, err := {{$jsonPkg.Ref "Marshal"}}({{$tmpVar}})
      if err !=nil{
      return nil, err
      }
      {{$output}} = string({{$jsonData}})
      }
    {{- else if eq $inputGoType.Type "any" }}
      // Any
      {
      {{$tmpVar}}, err := {{$jsonPkg.Ref "Marshal"}}({{$input}})
      if err !=nil{
      return nil, err
      }
      {{$output}} = string({{$tmpVar}})
      }
    {{- else if $inputGoType.IsPtr }}
      // Ptr
      if {{$input}} != nil {
      var {{$tmpVar}} {{$outputGoType.ElemType.Ref}}
      {{- template "storage.block.convert_value_to_internal" list (print "*" $input) $inputGoType.ElemType $tmpVar $outputGoType.ElemType $varNamesGenerator}}
      {{$output}} = &{{$tmpVar}}
      } else {
      {{$output}} = nil
      }
    {{- else if $inputGoType.IsSlice}}
      // Slice
      {
      {{$tmpVar}} := make({{$outputGoType.Ref}}, 0, len({{$input}}))
      for _, el := range {{$input}} {
      var res {{$outputGoType.ElemType.Ref}}
      {{- template "storage.block.convert_value_to_internal" list "el" $inputGoType.ElemType "res" $outputGoType.ElemType $varNamesGenerator}}
      {{$tmpVar}} = append({{$tmpVar}}, res)
      }
      {{$output}} = {{$tmpVar}}
      }
    {{- else if and (eq $inputGoType.Type "Time") (eq $inputGoType.Package "time") (not $inputGoType.WithTimezone)}}
        {{$output}} = ({{$input}}).UTC()
    {{- else}}
        {{$output}} = {{$input}}
    {{- end }}
{{- end }}

{{- define "storage.block.convert_value_to_service"}}
    {{- $jsonPkg := import "encoding/json" }}
    {{- $fmtPkg := import "fmt" }}
    {{- $input := index . 0 -}}
    {{- $inputGoType := index . 1 -}}
    {{- $output := index . 2 -}}
    {{- $outputGoType := index . 3 -}}
    {{- $varNamesGenerator := index . 4 -}}
    {{- $tmpVar := $varNamesGenerator.Var "tmp" -}}
    {{- if isModuleOneOf $outputGoType }}
        // OneOf From DB
        {{- $oneOf := getModuleOneOf $outputGoType.Type }}
        {
          if {{$input}} != nil {
            switch v := {{$input}}.Val.(type) {
            {{- range $value := $oneOf.SortedValues }}
                {{- $dbType := (goType $value.Value.Model).InLocalPackage.WithName (print "db" (goType $value.Value.Model).Type) }}
                {{- $model := $value.Value.Model }}
                case *{{ $dbType.Ref }}:
                  {{$tmpVar}}, err := {{template "storage.func.model_to_service" $model.Name}}(v)
                  if err!=nil {
                    return nil, {{$fmtPkg.Ref "Errorf"}}("convert {{$model.Name}} from db: %w", err)
                  }
                  {{$output}} = {{$tmpVar}}
            {{- end }}
            default:
            return nil, {{$fmtPkg.Ref "Errorf"}}("invalid {{$oneOf.Name}} value type: %T", {{$input}})
            }


          }
        }
    {{- else if or (isModuleEnum $outputGoType) (isCommonEnum $outputGoType) }}
      {
      {{$tmpVar}}, err := {{- template "storage.func.enum_to_service" $outputGoType.Type }}({{$input}})
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
      {{- template "storage.block.convert_value_to_service" list (print "*" $input) $inputGoType.ElemType $tmpVar $outputGoType.ElemType $varNamesGenerator}}
      {{$output}} = &{{$tmpVar}}
      }
    {{- else if $outputGoType.IsSlice}}
      {
      {{$tmpVar}} := make({{$outputGoType.Ref}}, 0, len({{$input}}))
      for _, el := range {{$input}} {
      var res {{$outputGoType.ElemType.Ref}}
      {{- template "storage.block.convert_value_to_service" list "el" $inputGoType.ElemType "res" $outputGoType.ElemType $varNamesGenerator}}
      {{$tmpVar}} = append({{$tmpVar}}, res)
      }
      {{$output}} = {{$tmpVar}}
      }
    {{- else}}
        {{$output}} = {{$input}}
    {{- end }}
{{- end }}

