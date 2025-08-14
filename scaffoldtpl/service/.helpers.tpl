{{- define "filterType" }}
    {{- $filterPkg := import "github.com/saturn4er/boilerplate-go/lib/filter"}}
    {{- if .IsSlice -}}
        {{$filterPkg.Ref "ArrayFilter"}}[{{.ElemType.Ref}}]
    {{- else -}}
        {{$filterPkg.Ref "Filter"}}[{{.Ref}}]
    {{- end -}}
{{- end }}

{{- define "orderType" -}}
    {{- $orderPkg := import "github.com/saturn4er/boilerplate-go/lib/order" -}}
    {{$orderPkg.Ref "Order"}}[{{.}}]
{{- end -}}

{{- define "orderDirectionType" -}}
    {{- $orderPkg := import "github.com/saturn4er/boilerplate-go/lib/order" -}}
    type OrderDirection = {{$orderPkg.Ref "Direction"}}
{{- end -}}

{{- define "copy_value" }}
    {{- $input := index . 0 -}}
    {{- $goType := index . 1 -}}
    {{- $output := index . 2 -}}
    {{- $varNamesGenerator := index . 3 -}}
    {{- if isModuleOneOf $goType}}
      {{- $oneOf := getModuleOneOf $goType.Type }}
      {{$output}} = copy{{$oneOf.Name}}({{$input}})
    {{- else if or (isModuleModel $goType) (isCommonModel $goType) }}
      {{$output}} = {{$input}}.Copy() // model
    {{- else if or (isModuleEnum $goType) (isCommonEnum $goType)}}
      {{$output}} =  {{$input}}// enum
    {{- else if $goType.IsPtr }}
      if {{$input}} != nil{
        {{- $tmpVar := $varNamesGenerator.Var "tmp" -}}
        var {{$tmpVar}} {{$goType.ElemType.Ref}}
        {{- template "copy_value" (list (print "(*" $input ")") $goType.ElemType $tmpVar $varNamesGenerator) }}
        {{$output}} = &{{$tmpVar}}
      }
    {{- else if $goType.IsSlice}}
      {{- $tmpVar := $varNamesGenerator.Var "tmp"}}
      {{- $itemVar := $varNamesGenerator.Var "i" }}
      {{$tmpVar}} := make({{$goType.Ref}}, 0, len({{$input}}))
      for _, {{$itemVar}} := range {{$input}} {
      {{- $itemCopyVar := $varNamesGenerator.Var "itemCopy" -}}
        var {{$itemCopyVar}} {{$goType.ElemType.Ref}}
        {{- template "copy_value" (list $itemVar $goType.ElemType $itemCopyVar $varNamesGenerator) }}
        {{$tmpVar}} = append({{$tmpVar}}, {{$itemCopyVar}})
      }
      {{$output}} = {{$tmpVar}}
    {{- else if $goType.IsMap }}
      {{- $tmpVar := $varNamesGenerator.Var "tmp" }}
      {{$tmpVar}} := make({{$goType.Ref}})
      for k, v := range {{$input}} {
      {{- $keyCopyVar := $varNamesGenerator.Var "keyCopy" }}
      {{- $valueCopyVar := $varNamesGenerator.Var "valueCopy" }}
        var {{$keyCopyVar}} {{$goType.KeyType.Ref}}
        var {{$valueCopyVar}} {{$goType.ElemType.Ref}}
        {{- template "copy_value" (list "k" $goType.KeyType $keyCopyVar $varNamesGenerator) }}
        {{- template "copy_value" (list "v" $goType.ElemType $valueCopyVar $varNamesGenerator) }}
        {{$tmpVar}}[{{$keyCopyVar}}] = {{$valueCopyVar}}
      }
      {{$output}} = {{$tmpVar}}
    {{- else }}
      {{$output}} = {{$input}}
    {{- end }}
{{- end }}


{{- define "equals_value" }}
    {{- $aVal := index . 0 -}}
    {{- $goType := index . 1 -}}
    {{- $bVal := index . 2 -}}
    {{- $varNamesGenerator := index . 3 -}}
    {{- if isModuleOneOf $goType}}
        {{- $oneOf := getModuleOneOf $goType.Type }}
        if !{{$aVal}}.{{$oneOf.Name}}Equals({{$bVal}}){
          return false
        }
    {{- else if or (isModuleModel $goType) (isCommonModel $goType) }}
        if !{{$aVal}}.Equals(&{{$bVal}}){
          return false
        }
    {{- else if or (isModuleEnum $goType) (isCommonEnum $goType)}}
        if {{$aVal}} != {{$bVal}}{
          return false
        }
    {{- else if $goType.IsPtr }}
      if ({{$aVal}} == nil) != ({{$bVal}} == nil) {
        return false
      }
      if {{$aVal}} != nil && {{$bVal}} != nil {
        {{- template "equals_value" (list (print "(*" $aVal ")") $goType.ElemType (print "(*" $bVal ")") $varNamesGenerator) }}
      }
    {{- else if $goType.IsSlice}}
      if len({{$aVal}}) != len({{$bVal}}){
        return false
      }
      {{- $iVar := $varNamesGenerator.Var "i" }}
      for {{$iVar}} := range {{$aVal}} {
        {{- $itemAVal := print $aVal "[" $iVar "]" }}
        {{- $itemBVal := print $bVal "[" $iVar "]" }}
        {{- template "equals_value" (list $itemAVal $goType.ElemType $itemBVal $varNamesGenerator) }}
      }
    {{- else if and (eq $goType.Package "net") (eq $goType.Type "IP") }}
      if !{{$aVal}}.Equal({{$bVal}}){
        return false
      }
    {{- else if $goType.IsMap }}
      // map comparision
      if len({{$aVal}}) != len({{$bVal}}){
        return false
      }
      {{- $kVar := $varNamesGenerator.Var "k" }}
      for {{$kVar}} := range {{$aVal}} {
          {{- $valA := $varNamesGenerator.Var "valA" }}
          {{- $valB := $varNamesGenerator.Var "valB" }}
          {{$valB}}, ok := {{$bVal}}[{{$kVar}}]
          if !ok {
            return false
          }
          {{$valA}} := {{$aVal}}[{{$kVar}}]
          {{- template "equals_value" (list $valA $goType.ElemType $valB $varNamesGenerator) }}
      }
    {{- else }}
        if {{$aVal}} != {{$bVal}}{
          return false
        }
    {{- end }}
{{- end }}
