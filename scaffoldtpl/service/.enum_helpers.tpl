{{- define "service._enum_declaration" }}
    {{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold/config.Enum */}}
    type {{ template "service.type.enum" . }} byte
    const (
    {{ template "service.const.enum_value" (dict "enum" . "value" (first .Values)) }} {{ template "service.type.enum" . }} = iota + 1
    {{- range $value := rest .Values }}
        {{ template "service.const.enum_value" (dict "enum" $ "value" $value) }}
    {{- end }}
    )
{{- end }}

{{- define "service._enum_is_valid_helper"}}
    {{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold/config.Enum */}}
    {{- $typeName := (include "service.type.enum" .) }}
    {{- $receiverName := $typeName | receiverName }}
    func ({{$receiverName}} {{$typeName}}) IsValid() bool {
    return {{$receiverName}} > 0 && {{$receiverName}} < {{add (len .Values) 1}}
    }
{{- end }}

{{- define "service._enum_all_values_helper"}}
    {{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold/config.Enum */}}
    {{- $typeName := (include "service.type.enum" .) }}
    func All{{$typeName}}() []{{$typeName}} {
    return []{{$typeName}}{
    {{- range $value := .Values -}}
        {{- template "service.const.enum_value" (dict "enum" $ "value" $value) }},
    {{- end -}}
    }
    }
{{- end }}

{{- define "service._enum_is_helper"}}
    {{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold/config.Enum */}}
    {{- $typeName := (include "service.type.enum" .) }}
    {{- $receiverName := $typeName | receiverName }}
    {{- if .Helpers.Is }}
        {{- range $value := .Values }}
          func ({{$receiverName}} {{$typeName}}) Is{{$value}}() bool {
          return {{$receiverName}} == {{- template "service.const.enum_value" (dict "enum" $ "value" $value) }}
          }
        {{- end}}
    {{- end }}
{{- end }}

{{- define "service._enum_stringer_helper"}}
    {{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold/config.Enum */}}
    {{- $strconvPkg := import "strconv"}}
    {{- $typeName := (include "service.type.enum" .) }}
    {{- $receiverName := $typeName | receiverName }}
    {{- if .Helpers.Stringer }}
      func ({{$receiverName}} {{$typeName}}) String() string {
      const names = "{{range $i, $v := .Values}}{{if $i}}{{end}}{{$v}}{{end}}"
      {{$totalLength := 0}}
      var indexes = [...]int32{{"{0, "}} {{range $i, $v := .Values}}{{- $totalLength = add $totalLength (len $v) -}}{{$totalLength}}, {{- end -}}{{"}"}}
      if {{$receiverName}} < 1 || {{$receiverName}} > {{len .Values}} {
      return "{{$typeName}}("+{{$strconvPkg.Ref "FormatInt"}}(int64({{$receiverName}}), 10)+")"
      }

      return names[indexes[{{$receiverName}}-1]:indexes[{{$receiverName}}]]
      }
    {{- end}}
{{- end }}


{{- define "service._enum_is_category_helper"}}
    {{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold/config.Enum */}}
    {{- $typeName := (include "service.type.enum" .) }}
    {{- $receiverName := $typeName | receiverName }}
    {{- range $isCategory := .Helpers.IsCategory }}
      var {{$isCategory.Name}}{{plural $.Name}} = []{{$typeName}}{
      {{- range $isCategoryValue := $isCategory.Values }}
          {{- template "service.const.enum_value" (dict "enum" $ "value" $isCategoryValue) }},
      {{- end }}
      }
      func ({{$receiverName}} {{$typeName}}) Is{{$isCategory.Name}}() bool {
      if {{$receiverName}} < 1 || {{$receiverName}} > {{add (len .Values) 1}} {
      return false
      }
      return [...]bool{false,
      {{- range $value := $.Values -}}
          {{- $isTrue := false -}}
          {{- range $isCategoryValue := $isCategory.Values }}
              {{- if eq $value $isCategoryValue  }}
                  {{- $isTrue = true -}}
              {{- end }}
          {{- end -}}
          {{$isTrue}},
      {{- end -}}}[{{$receiverName}}]
      }
    {{- end }}
{{- end }}

{{- define "service._enum_validate_helper"}}
    {{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold/config.Enum */}}
    {{- $typeName := (include "service.type.enum" .) }}
    {{- $receiverName := $typeName | receiverName }}
    {{- $validationPkg := import "github.com/go-ozzo/ozzo-validation/v4" "validation"}}
    {{- if .Helpers.Validate }}
      func ({{$receiverName}} {{$typeName}}) Validate() error {
      return {{$validationPkg.Ref "Validate"}}({{$receiverName}}, {{$validationPkg.Ref "In"}}(
      {{- range $value := .Values -}}
          {{- template "service.const.enum_value" (dict "enum" $ "value" $value) }}
      {{- end -}}
      ))
      }
    {{- end }}
{{- end }}
