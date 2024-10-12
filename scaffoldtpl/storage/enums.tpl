{
"file_path": "{{.Module}}/{{.Module}}storage/gen.enums.go",
"package_name": "{{.Module}}storage",
"package_path": "{{.Config.RootPackageName}}/{{.Module}}/{{.Module}}storage",
"condition": "len(Config.Modules[Module].Value.Types.Enums) > 0"
}
<><><>

{{ $dbutil := import "github.com/saturn4er/boilerplate-go/lib/dbutil" }}
{{ $clause := import "gorm.io/gorm/clause" }}
{{ $fmtPkg := import "fmt" }}
{{ $servicePkg :=  import (print $.Config.RootPackageName "/" $.Module "/" $.Module "service") }}

{{ $module := (index $.Config.Modules $.Module).Value}}

{{- define "storage.enum.values" }}
  const (
    {{- range $value := .Values }}
        {{$.Name | lCamelCase}}{{$value}} = "{{$value | snakeCase}}"
    {{- end}}
  )
{{- end }}
{{- define "storage.enum.to_db_converter"}}
  {{ $fmtPkg := import "fmt" }}
  {{- $enumGoType := goType . -}}
  {{- $enumPkg := $enumGoType.PackageImport -}}
  {{ $argName := print (.Name| lCamelCase) "Value" }}
  func {{template "storage.func.enum_to_internal" .Name}}({{$argName}} {{$enumGoType.Ref}}) (string, error){
    result, ok := map[{{$enumGoType.Ref}}]string{
    {{- range $value := .Values}}
        {{$enumPkg.Ref (print $.Name $value)}}: {{$.Name | lCamelCase}}{{$value}},
    {{- end }}
    }[{{$argName}}]
    if !ok {
      return "", {{$fmtPkg.Ref "Errorf"}}("unknown {{.Name}} value: %d", {{$argName}})
    }
    return result, nil
  }
{{- end}}

{{- define "storage.enum.from_db_converter"}}
    {{ $fmtPkg := import "fmt" }}
    {{- $enumGoType := goType . }}
    {{- $enumPkg := $enumGoType.PackageImport }}
    {{- $argName := print (.Name| lCamelCase) "Value" }}
    func {{template "storage.func.enum_to_service" .Name}}({{$argName}} string) ({{$enumGoType.Ref}}, error){
    result, ok := map[string]{{$enumGoType.Ref}}{
    {{- range $value := .Values}}
        {{$.Name | lCamelCase}}{{$value}}: {{$enumPkg.Ref (print $.Name $value)}},
    {{- end }}
    }[{{$argName}}]
    if !ok {
    return 0, {{$fmtPkg.Ref "Errorf"}}("unknown {{.Name}} db value: %s", {{$argName}})
    }
    return result, nil
    }
{{- end}}

{{- range $enum := $module.Types.Enums }}
{{- template "storage.enum.values" $enum }}
{{- template "storage.enum.to_db_converter" $enum }}
{{- template "storage.enum.from_db_converter" $enum }}
{{- end }}

{{- range $enum := .Config.Types.Enums }}
    {{- template "storage.enum.values" $enum }}
    {{- template "storage.enum.to_db_converter" $enum }}
    {{- template "storage.enum.from_db_converter" $enum }}
{{- end}}
