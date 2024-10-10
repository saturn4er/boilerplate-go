{{- define "storage.func.enum_to_internal" -}}
  convert{{.}}ToDB{{.}}
{{- end -}}

{{- define "storage.func.enum_to_service" -}}
  convertDB{{.}}To{{.}}
{{- end -}}

{{- define "storage.func.one_of_to_internal" -}}
  {{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold.ConfigOneOf*/ -}}
    convert{{.Name}}ToDB{{.Name}}
{{- end -}}

{{- define "storage.func.one_of_to_string_ptr" -}}
    {{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold.ConfigOneOf*/ -}}
    convert{{.Name}}ToStrPtr
{{- end -}}

{{- define "storage.func.one_of_to_service" -}}
    {{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold.ConfigOneOf*/ -}}
    convertDB{{.Name}}To{{.Name}}
{{- end -}}

{{- define "storage.func.one_of_string_ptr_to_service" -}}
    {{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold.ConfigOneOf*/ -}}
    convertStrPtrTo{{.Name}}
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
        {{- $oneOf := getModuleOneOf $inputGoType.Type }}
        {{- if and $outputGoType.IsPtr (eq $outputGoType.ElemType.Type "string") }}
            // convert oneof to string ptr
            {{$tmpVar}}, err := {{template "storage.func.one_of_to_string_ptr" $oneOf}}({{$input}})
            if err != nil {
            return nil, err
            }
            {{$output}} = {{$tmpVar}}
        {{- else if eq $outputGoType.Type "string" }}
            // convert oneof to string
            {{$tmpVar}}, err := {{template "storage.func.one_of_to_string_ptr" $oneOf}}({{$input}})
            if err != nil {
            return nil, err
            }
            if {{$tmpVar}} == nil {
              {{$output}} = "null"
            }else{
              {{$output}} = *{{$tmpVar}}
            }
        {{- else }}
            // oneof to db
            {{$tmpVar}}, err := {{template "storage.func.one_of_to_internal" $oneOf}}({{$input}})
            if err != nil {
            return nil, err
            }
            {{$output}} = {{$tmpVar}}
        {{- end }}
    {{- else if or (isModuleEnum $inputGoType) (isCommonEnum $inputGoType) }}
      // enum to db
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
        {{if not $inputGoType.IsPtr}}
            {{$input = print "&" $input}}
        {{- end }}

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
    {{- else if $inputGoType.IsMap }}
      // map to db
    {{- else if or (eq $inputGoType.Type "any") }}
      // any to db
      {
          if {{$input}} != nil {
            {{$tmpVar}}, err := {{$jsonPkg.Ref "Marshal"}}({{$input}})
            if err !=nil{
            return nil, err
            }
            {{$tmpVar2 := $varNamesGenerator.Var "marshaledValue"}}
            {{$tmpVar2}} := string({{$tmpVar}})
            {{$output}} = &{{$tmpVar2}}
          }
      }
    {{- else if and $inputGoType.IsPtr $outputGoType.IsPtr }}
      {
        // ptr to db
        if {{$input}} != nil {
        var {{$tmpVar}} {{$outputGoType.ElemType.Ref}}
        {{- template "storage.block.convert_value_to_internal" list (print "*" $input) $inputGoType.ElemType $tmpVar $outputGoType.ElemType $varNamesGenerator}}
        {{$output}} = &{{$tmpVar}}
        } else {
        {{$output}} = nil
        }
      }
    {{- else if $inputGoType.IsSlice}}
      // slice to db
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
        // time to db
        {{$output}} = ({{$input}}).UTC()
    {{- else }}
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
        {{- $oneOfType := getModuleOneOf $outputGoType.Type }}
        {
          // one-of from db
          {{- if and $inputGoType.IsPtr (eq $inputGoType.ElemType.Type "string")  }}
            {{$tmpVar}}, err := {{template "storage.func.one_of_string_ptr_to_service" $oneOfType}}({{$input}})
            if err != nil{
              return nil, {{$fmtPkg.Ref "Errorf"}}("convert {{$oneOfType.Name}} to service type: %w", err)
            }
            {{$output}} = {{$tmpVar}}
          {{- else if eq $inputGoType.Type "string" }}
            {{$tmpVar}}, err := {{template "storage.func.one_of_string_ptr_to_service" $oneOfType}}(&{{$input}})
            if err != nil{
              return nil, {{$fmtPkg.Ref "Errorf"}}("convert {{$oneOfType.Name}} to service type: %w", err)
            }
            {{$output}} = {{$tmpVar}}
          {{- else }}
              {{$tmpVar}}, err := {{template "storage.func.one_of_to_service" $outputGoType.Type}}({{$tmpVar}})
              if err != nil{
                return nil, {{$fmtPkg.Ref "Errorf"}}("convert {{$oneOfType.Name}} to service type: %w", err)
              }
              {{$output}} = {{$tmpVar}}
          {{- end }}
        }
    {{- else if or (isModuleEnum $outputGoType) (isCommonEnum $outputGoType) }}
      {
      // enum from db
      {{$tmpVar}}, err := {{- template "storage.func.enum_to_service" $outputGoType.Type }}({{$input}})
      if err != nil {
      return nil, err
      }
      {{$output}} = {{$tmpVar}}
      }
    {{- else if or (isModuleModel $outputGoType) (isCommonModel $outputGoType) (eq $outputGoType.Type "any")}}
      {
        {{- if $inputGoType.IsPtr}}
          // model/any ptr from db
          if {{$input}} != nil {
            var {{$tmpVar}} {{$outputGoType.Ref}}
            if err := {{$jsonPkg.Ref "Unmarshal"}}([]byte(*{{$input}}), &{{$tmpVar}}); err !=nil{
            return nil, err
            }
            {{$output}} = {{$tmpVar}}
          }
        {{- else}}
          // model/any from db
          var {{$tmpVar}} {{$outputGoType.Ref}}
          if err := {{$jsonPkg.Ref "Unmarshal"}}([]byte({{$input}}), &{{$tmpVar}}); err !=nil{
          return nil, err
          }
          {{$output}} = {{$tmpVar}}
        {{- end }}
      }
    {{- else if $outputGoType.IsMap}}
      {
        // map from db
        var val map[{{$outputGoType.ElemType.DBAlternative.Ref}}]{{$outputGoType.ElemType.DBAlternative.Ref}}
      }
    {{- else if and $outputGoType.IsPtr $inputGoType.IsPtr }}
      // ptr from map
      if {{$input}} != nil {
      var {{$tmpVar}} {{$outputGoType.ElemType.Ref}}
      {{- template "storage.block.convert_value_to_service" list (print "*" $input) $inputGoType.ElemType $tmpVar $outputGoType.ElemType $varNamesGenerator}}
      {{$output}} = &{{$tmpVar}}
      }
    {{- else if $outputGoType.IsSlice}}
      {
          // slice from db
      {{$tmpVar}} := make({{$outputGoType.Ref}}, 0, len({{$input}}))
      for _, el := range {{$input}} {
      var res {{$outputGoType.ElemType.Ref}}
      {{- template "storage.block.convert_value_to_service" list "el" $inputGoType.ElemType "res" $outputGoType.ElemType $varNamesGenerator}}
      {{$tmpVar}} = append({{$tmpVar}}, res)
      }
      {{$output}} = {{$tmpVar}}
      }
    {{- else}}
        // default from db
        {{$output}} = {{$input}}
    {{- end }}
{{- end }}

