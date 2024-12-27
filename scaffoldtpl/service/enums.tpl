{
"file_path": "{{.Module}}/{{.Module}}service/gen.enums.go",
"package_name": "{{.Module}}service",
"package_path": "{{.Config.RootPackageName}}/{{.Module}}/{{.Module}}",
"condition": "len(Config.Modules[Module].Value.Types.Enums) > 0"
}
<><><>

{{ $module := (index $.Config.Modules $.Module).Value}}
{{- $validationPkg := import "github.com/go-ozzo/ozzo-validation/v4" "validation"}}
{{- $strconvPkg := import "strconv"}}

{{- range $enum := $module.Types.Enums }}
  type {{$enum.Name}} byte
  const (
  {{$enum.Name}}{{index $enum.Values 0}} {{$enum.Name}} = iota+1
  {{- range $value := (slice $enum.Values 1)}}
  {{$enum.Name}}{{$value}}
  {{- end}}
  )
  {{ userCodeBlock (printf "%s methods" $enum.Name) }}
  {{- if not .Helpers}}
  {{- continue }}
  {{- end }}
  {{- $receiverName := slice $enum.Name 0 1 | lCamelCase}}
  {{- if .Helpers.AllValues.VarName }}
    var {{.Helpers.AllValues.VarName}} = []{{$enum.Name}}{{"{"}}{{- range $value := $enum.Values -}}{{$enum.Name}}{{$value}}, {{- end -}}{{"}"}}
  {{- end }}
  {{- if .Helpers.IsValid }}
      func ({{$receiverName}} {{$enum.Name}}) IsValid() bool {
        return {{$receiverName}} > 0 && {{$receiverName}} < {{addInts (len $enum.Values) 1}}
      }
  {{- end }}
  {{- range $isCategory := .Helpers.IsCategory }}
    var {{$isCategory.Name}}{{plural $enum.Name}} = []{{$enum.Name}}{
        {{- range $isCategoryValue := $isCategory.Values }}
            {{$enum.Name}}{{$isCategoryValue}},
        {{- end }}
    }
    func ({{$receiverName}} {{$enum.Name}}) Is{{$isCategory.Name}}() bool {
      if {{$receiverName}} < 1 || {{$receiverName}} > {{addInts (len $enum.Values) 1}} {
        return false
      }
      return []bool{false,
        {{- range $value := $enum.Values -}}
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
  {{- if .Helpers.Is }}
    {{- range $value := $enum.Values }}
      func ({{$receiverName}} {{$enum.Name}}) Is{{$value}}() bool {
        return {{$receiverName}} == {{$enum.Name}}{{$value}}
      }
    {{- end}}
  {{- end }}
  {{- if .Helpers.Validate }}
    func ({{$receiverName}} {{$enum.Name}}) Validate() error {
    return {{$validationPkg.Ref "Validate"}}({{$receiverName}}, {{$validationPkg.Ref "In"}}({{- range $value := $enum.Values -}}{{$enum.Name}}{{$value}}, {{- end -}}))
    }
  {{- end }}
  {{- if .Helpers.AllValues.FuncName }}
    func {{.Helpers.AllValues.FuncName}}() []{{$enum.Name}} {
      return []{{$enum.Name}}{{"{"}}{{- range $value := $enum.Values -}}{{$enum.Name}}{{$value}}, {{- end -}}{{"}"}}
    }
  {{- end }}
  {{- if .Helpers.Stringer }}
    func ({{$receiverName}} {{$enum.Name}}) String() string {
    const names = "{{range $i, $v := .Values}}{{if $i}}{{end}}{{$v}}{{end}}"
    {{$totalLength := 0}}
    var indexes = [...]int32{{"{0, "}} {{range $i, $v := .Values}}{{- $totalLength = addInts $totalLength (len $v) -}}{{$totalLength}}, {{- end -}}{{"}"}}
    if {{$receiverName}} < 1 || {{$receiverName}} > {{len $enum.Values}} {
    return "{{$enum.Name}}("+{{$strconvPkg.Ref "FormatInt"}}(int64({{$receiverName}}), 10)+")"
    }

    return names[indexes[{{$receiverName}}-1]:indexes[{{$receiverName}}]]
    }
  {{- end}}
{{- end }}
