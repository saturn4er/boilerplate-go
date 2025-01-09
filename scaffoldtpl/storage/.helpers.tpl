{{- define "storage.func.enum_to_internal" -}}
  convert{{.}}ToDB
{{- end -}}

{{- define "storage.func.enum_to_service" -}}
  convert{{.}}FromDB
{{- end -}}

{{- define "storage.func.one_of_to_internal" -}}
    {{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold.ConfigOneOf*/ -}}
    convert{{.Name}}ToDB
{{- end -}}

{{- define "storage.func.one_of_to_string_ptr" -}}
    {{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold.ConfigOneOf*/ -}}
    convert{{.Name}}ToStrPtr
{{- end -}}

{{- define "storage.func.one_of_to_service" -}}
    {{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold.ConfigOneOf*/ -}}
    convert{{.Name}}FromDB
{{- end -}}

{{- define "storage.func.one_of_string_ptr_to_service" -}}
    {{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold.ConfigOneOf*/ -}}
    convertStrPtrTo{{.Name}}
{{- end -}}

{{- define "storage.func.build_db_filter" -}}
  build{{(print . "Filter")}}Expr
{{- end}}

{{- define "storage.func.table_model_to_internal" -}}
  convert{{.}}ToDB
{{- end -}}
{{- define "storage.func.table_model_to_service" -}}
  convert{{.}}FromDB
{{- end -}}

{{- define "storage.func.json_model_to_internal" -}}
  convert{{.}}ToJsonModel
{{- end -}}
{{- define "storage.func.json_model_to_service" -}}
  convert{{.}}FromJsonModel
{{- end -}}

{{- define "storage.func.errors_wrapper" -}}
  wrap{{.}}Error
{{- end -}}

{{- define "storage.field.json_model"}}
    {{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold/config.ModelField*/ -}}
    {{- $name := .Name}}
    {{- if or (eq $name "Scan") (eq $name "Value") -}}
        {{- $name = print $name "Val" }}
    {{- end -}}
    {{- $name -}}
{{- end }}

{{- define "storage.block.convert_value_to_internal"}}
    {{- $jsonPkg := import "encoding/json" -}}
    {{- $fmtPkg := import "fmt" -}}
    {{- $input := index . 0 -}}
    {{- $inputGoType := index . 1 -}}
    {{- $output := index . 2 -}}
    {{- $outputGoType := index . 3 -}}
    {{- $varNamesGenerator := index . 4 -}}
    {{- $tmpVar := $varNamesGenerator.Var "tmp" -}}
    {{- if or (isModuleOneOf $inputGoType)  }}
      {{- $oneOf := getModuleOneOf $inputGoType.Type }}
      {{$tmpVar}}, err := {{template "storage.func.one_of_to_internal" $oneOf}}({{$input}})
      if err != nil {
      return nil, err
      }
      {{- if $outputGoType.IsPtr}}
          {{call $output $tmpVar}}
      {{- else}}
          {{call $output print "*" $tmpVar}}
      {{- end }}
    {{- else if and $inputGoType.IsPtr (isModuleOneOf $inputGoType.ElemType) }}
        {{- $oneOf := getModuleOneOf $inputGoType.ElemType.Type }}
        if {{$input}} != nil {
        {{$tmpVar}}, err := {{template "storage.func.one_of_to_internal" $oneOf}}(*{{$input}})
        if err != nil {
        return nil, err
        }
        {{call $output $tmpVar}}
        } else {
        {{call $output "nil"}}
        }
    {{- else if or (isModuleEnum $inputGoType) (isCommonEnum $inputGoType) }}
      {{$tmpVar}}, err := {{- template "storage.func.enum_to_internal" $inputGoType.Type }}({{$input}})
      if err != nil {
      return nil, err
      }
      {{call $output $tmpVar}}
    {{- else if or (isModuleModel $inputGoType) (isCommonModel $inputGoType) }}
          {{- $model := getModel $inputGoType }}
          {{$tmpVar}}, err := {{template "storage.func.json_model_to_internal" $model.Name}}(toPtr({{$input}}))
          if err!=nil {
          return nil, fmt.Errorf("convert {{$model.Name}} to db: %w", err)
          }
          {{call $output (print "*" $tmpVar)}}
    {{- else if and $inputGoType.IsPtr (or (isModuleModel $inputGoType.ElemType) (isCommonModel $inputGoType.ElemType)) }}
      if {{$input}} != nil {
        {{- $model := getModel $inputGoType.ElemType}}
        {{$tmpVar}}, err := {{template "storage.func.json_model_to_internal" $model.Name}}({{$input}})
        if err!=nil {
        return nil, fmt.Errorf("convert {{$model.Name}} to db: %w", err)
        }
        {{call $output $tmpVar}}
        } else {
        {{call $output "nil"}}
        }
    {{- else if $inputGoType.IsMap }}
      {{$tmpVar}} := make({{$outputGoType.Ref}}, len({{$input}}))
      {{- $kVar := $varNamesGenerator.Var "k"}}
      {{- $vVar := $varNamesGenerator.Var "v"}}
      for {{$kVar}}, {{$vVar}} := range {{$input}} {
      {{- if $inputGoType.KeyType.IsSimple}}
          {{- $itemOutput := putToMapFn $tmpVar $kVar }}
          {{- template "storage.block.convert_value_to_internal" list $vVar $inputGoType.ElemType $itemOutput $outputGoType.ElemType $varNamesGenerator}}
      {{- else }}
          {{- $kResultVar := $varNamesGenerator.Var "kResult"}}
          {{- $kOutput := setToNewVarFn $kResultVar }}
          {{- template "storage.block.convert_value_to_internal" list $kVar $inputGoType.KeyType $kOutput $outputGoType.KeyType $varNamesGenerator}}
          {{- $vResultVar := $varNamesGenerator.Var "vResult"}}
          {{- $itemOutput := putToMapFn $tmpVar $kResultVar }}
          {{- template "storage.block.convert_value_to_internal" list $vVar $inputGoType.ElemType $itemOutput $outputGoType.ElemType $varNamesGenerator}}
      {{- end }}
      }
      {{call $output $tmpVar}}
    {{- else if eq $inputGoType.Type "any" }}
      if {{$input}} != nil {
      {{$tmpVar}}, err := {{$jsonPkg.Ref "Marshal"}}({{$input}})
      if err !=nil{
      return nil, err
      }
      {{$tmpVar2 := $varNamesGenerator.Var "marshaledValue"}}
      {{$tmpVar2}} := string({{$tmpVar}})
      {{call $output (print "toPtr(" $tmpVar2 ")")}}
      }else{
      {{call $output "nil"}}
      }
    {{- else if and $inputGoType.IsPtr (eq $inputGoType.ElemType.Type "any") }}
      if {{$input}} != nil && fromPtr({{$input}}) != nil {
      {{$tmpVar}}, err := {{$jsonPkg.Ref "Marshal"}}(*{{$input}})
      if err !=nil{
      return nil, err
      }
      {{$tmpVar2 := $varNamesGenerator.Var "marshaledValue"}}
      {{$tmpVar2}} := string({{$tmpVar}})
      {{call $output (print "toPtr(" $tmpVar2 ")")}}
      }else{
      {{call $output "nil"}}
      }
    {{- else if and $inputGoType.IsPtr $outputGoType.IsPtr }}
      if {{$input}} == nil {
      {{call $output "nil"}}
      }else{
      {{- $localOutput := chainFn takePtrFn $output}}
      {{- template "storage.block.convert_value_to_internal" list (print "fromPtr(" $input ")") $inputGoType.ElemType $localOutput $outputGoType.ElemType $varNamesGenerator}}
      }
    {{- else if $inputGoType.IsSlice}}
      {{$tmpVar}} := make({{$outputGoType.Ref}}, 0, len({{$input}}))
      for _, el := range {{$input}} {
      {{- $sliceOutput := appendFn $tmpVar}}
      {{- template "storage.block.convert_value_to_internal" list "el" $inputGoType.ElemType $sliceOutput $outputGoType.ElemType $varNamesGenerator}}
      }
      {{call $output $tmpVar}}
    {{- else if and (eq $inputGoType.Type "IP") (eq $inputGoType.Package "net")}}
        {{call $output (print "(ipValue)(" $input ")")}}
    {{- else if and (eq $inputGoType.Type "Time") (eq $inputGoType.Package "time") (not $inputGoType.WithTimezone)}}
      {{call $output (print "(" $input ").UTC()")}}
    {{- else }}
        {{call $output $input}}
    {{- end }}
{{- end }}

{{- define "storage.block.convert_value_to_service"}}
    {{- $jsonPkg := import "encoding/json" -}}
    {{- $fmtPkg := import "fmt" -}}
    {{- $input := index . 0 -}}
    {{- $inputGoType := index . 1 -}}
    {{- $output := index . 2 -}}
    {{- $outputGoType := index . 3 -}}
    {{- $varNamesGenerator := index . 4 -}}
    {{- $tmpVar := $varNamesGenerator.Var "tmp" -}}
    {{- if isModuleOneOf $outputGoType }}
        {{- $oneOfType := getModuleOneOf $outputGoType.Type }}
        {{$tmpVar}}, err := {{template "storage.func.one_of_to_service" $oneOfType}}({{$input}})
        if err != nil {
        return nil, {{$fmtPkg.Ref "Errorf"}}("convert {{$oneOfType.Name}} to service type: %w", err)
        }
        {{call $output $tmpVar}}
    {{- else if and $outputGoType.IsPtr (isModuleOneOf $outputGoType.ElemType) }}
        if({{$input}} != nil){
          {{- $oneOfType := getModuleOneOf $outputGoType.ElemType.Type }}
          {{$tmpVar}}, err := {{template "storage.func.one_of_to_service" $oneOfType}}({{$input}})
          if err != nil {
          return nil, {{$fmtPkg.Ref "Errorf"}}("convert {{$oneOfType.Name}} to service type: %w", err)
          }
          {{call $output (print "toPtr(" $tmpVar ")")}}
        } else {
          {{call $output "nil"}}
        }
    {{- else if or (isModuleEnum $outputGoType) (isCommonEnum $outputGoType) }}
      {{$tmpVar}}, err := {{- template "storage.func.enum_to_service" $outputGoType.Type }}({{$input}})
      if err != nil {
      return nil, err
      }
      {{call $output $tmpVar}}
    {{- else if or (isModuleModel $outputGoType) (isCommonModel $outputGoType) }}
        {{$tmpVar}}, err := {{- template "storage.func.json_model_to_service" $outputGoType.Type }}(toPtr({{$input}}))
        if err != nil{
        return nil, err
        }

        {{call $output (print "fromPtr(" $tmpVar ")")}}
    {{- else if and $outputGoType.IsPtr (or (isModuleModel $outputGoType.ElemType) (isCommonModel $outputGoType.ElemType)) }}
        if {{$input}} != nil {
          {{$tmpVar}}, err := {{- template "storage.func.json_model_to_service" $outputGoType.ElemType.Type }}({{$input}})
          if err != nil{
          return nil, err
          }
            {{call $output $tmpVar}}
        } else {
          {{call $output "nil"}}
        }
    {{- else if and $inputGoType.IsPtr $outputGoType.IsPtr }}
      if {{$input}} == nil {
        {{call $output "nil"}}
      }else{
      {{- $localOutput := chainFn takePtrFn $output}}
      {{- template "storage.block.convert_value_to_service" list (print "fromPtr(" $input ")") $inputGoType.ElemType $localOutput $outputGoType.ElemType $varNamesGenerator}}
      }
    {{- else if (eq $outputGoType.Type "any")}}
        {{- if $inputGoType.IsPtr}}
          if {{$input}} != nil {
          var {{$tmpVar}} {{$outputGoType.Ref}}
          if err := {{$jsonPkg.Ref "Unmarshal"}}([]byte(*{{$input}}), &{{$tmpVar}}); err !=nil{
          return nil, err
          }
          {{call $output $tmpVar}}
          }else {
          {{call $output "nil"}}
          }
        {{- else}}
          var {{$tmpVar}} {{$outputGoType.Ref}}
          if err := {{$jsonPkg.Ref "Unmarshal"}}([]byte({{$input}}), &{{$tmpVar}}); err !=nil{
          return nil, err
          }
          {{call $output $tmpVar}}
        {{- end }}
    {{- else if $outputGoType.IsMap}}
      {{$tmpVar}} := make(map[{{$outputGoType.KeyType.Ref}}]{{$outputGoType.ElemType.Ref}}, len({{$input}}))
      {{- $kVar := $varNamesGenerator.Var "k"}}
      {{- $vVar := $varNamesGenerator.Var "v"}}
      for {{$kVar}}, {{$vVar}} := range {{$input}} {
        {{- if $inputGoType.KeyType.IsSimple }}
            {{- $valueOutput := putToMapFn $tmpVar $kVar}}
            {{- template "storage.block.convert_value_to_service" list $vVar $inputGoType.ElemType $valueOutput $outputGoType.ElemType $varNamesGenerator}}
        {{- else }}
          {{- $kResultVar := $varNamesGenerator.Var "kResult"}}
          {{- $keyOutput := setToNewVarFn $kResultVar}}
          {{- template "storage.block.convert_value_to_service" list $kVar $inputGoType.KeyType $keyOutput $outputGoType.KeyType $varNamesGenerator}}
          {{- $vResultVar := $varNamesGenerator.Var "vResult"}}
          {{- $valueOutput := putToMapFn $tmpVar $kResultVar}}
          {{- template "storage.block.convert_value_to_service" list $vVar $inputGoType.ElemType $valueOutput $outputGoType.ElemType $varNamesGenerator}}
        {{- end }}
      }
      {{call $output $tmpVar}}
    {{- else if $outputGoType.IsSlice}}
      {{$tmpVar}} := make({{$outputGoType.Ref}}, 0, len({{$input}}))
      for _, el := range {{$input}} {
      {{$itemOutput := appendFn $tmpVar}}
      {{- template "storage.block.convert_value_to_service" list "el" $inputGoType.ElemType $itemOutput $outputGoType.ElemType $varNamesGenerator}}
      }
      {{call $output $tmpVar}}
    {{- else if and (eq $outputGoType.Type "IP") (eq $outputGoType.Package "net")}}
        {{call $output (print "(net.IP)(" $input ")")}}
    {{- else}}
      {{call $output $input}}
    {{- end }}
{{- end }}

