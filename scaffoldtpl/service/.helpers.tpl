{{- define "filterType" }}
    {{- $filterPkg := import "github.com/saturn4er/boilerplate-go/lib/filter"}}
    {{- if .IsSlice -}}
        {{$filterPkg.Ref "ArrayFilter"}}[{{.ElemType.Ref}}]
    {{- else -}}
        {{$filterPkg.Ref "Filter"}}[{{.Ref}}]
    {{- end -}}
{{- end }}

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
    {{- else }}
      {{$output}} = {{$input}}
    {{- end }}
{{- end }}
