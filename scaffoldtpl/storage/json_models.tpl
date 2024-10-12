{
"file_path": "{{.Module}}/{{.Module}}storage/gen.json_models.go",
"package_name": "{{.Module}}storage",
"package_path": "{{.Config.RootPackageName}}/{{.Module}}/{{.Module}}storage",
"condition": "len(Config.Modules[Module].Value.Types.Models) > 0"
}
<><><>

{{- $dbutil := import "github.com/saturn4er/boilerplate-go/lib/dbutil" }}
{{- $jsonPkg := import "encoding/json" }}
{{- $driverPkg := import "database/sql/driver"}}
{{- $fmtPkg := import "fmt" }}
{{ $clause := import "gorm.io/gorm/clause" }}
{{ $servicePkg :=  import (print $.Config.RootPackageName "/" $.Module "/" $.Module "service") }}
{{ $module := (index $.Config.Modules $.Module).Value}}


{{- define "fieldGormTag" -}}
{{- $fieldType := (goType .Type).DBAlternative }}
{{- if or $fieldType.GormType .PrimaryKey -}}
  `gorm:"{{- if $fieldType.GormType -}}type:{{- $fieldType.GormType -}};{{- end -}}{{if .PrimaryKey}}primaryKey{{end}}"`
{{- else -}}{{- end -}}
{{- end -}}


{{- range $model := $module.Types.Models }}
    {{- $modelGoType := goType $model }}

    {{- $servicePtrType := $modelGoType.Ptr -}}
    {{- $servicePtrTypeRef := $servicePtrType.Ref -}}
    {{- $jsonType := $modelGoType.InLocalPackage.WithName (print "json" $modelGoType.Type)  -}}
    {{- $jsonTypeName := $jsonType.Type }}
    {{- $varNamesGenerator := varNamesGenerator }}
    {{- $receiverName := $model.Name | receiverName }}

    type {{$jsonType.Ref}} struct {
    {{- range $field := $model.Fields}}

        {{- if not $field.DoNotPersists }}
            {{- $fieldType := (goType $field.Type).DBAlternative }}
            {{ template "storage.field.json_model" $field }} {{(goType $field.Type).DBAlternative.Ref}} `json:"{{$field.Name | snakeCase}}"`
        {{- end }}
    {{- end }}
    }

    func ({{$receiverName}} *{{$jsonTypeName}}) Scan(value any) error {
    return {{$jsonPkg.Ref "Unmarshal"}}(value.([]byte), {{$receiverName}})
    }

    func ({{$receiverName}} {{$jsonTypeName}}) Value() ({{$driverPkg.Ref "Value"}}, error) {
    return {{$jsonPkg.Ref "Marshal"}}({{$receiverName}})
    }

    {{ $toArgName := "src" }}
    func {{ template "storage.func.json_model_to_internal" $model.Name}}({{$toArgName}} {{$servicePtrTypeRef}}) (*{{$jsonType.Ref}}, error){
      result := &{{$jsonType.Ref}}{}
      {{- range $field := $model.Fields }}
        {{- if $field.DoNotPersists }} {{- continue }} {{- end }}
        {{- $input := print $toArgName "." $field.Name -}}
        {{- $inputGoType := goType $field.Type -}}
        {{- $dbFieldName := include "storage.field.json_model" $field -}}
        {{- $output := print "result." $dbFieldName -}}
        {{- $outputGoType := (goType $field.Type).DBAlternative -}}

        {{ template "storage.block.convert_value_to_internal" (list $input $inputGoType $output $outputGoType $varNamesGenerator) }}
      {{- end }}
    return result, nil
    }

    {{ $fromArgName := "src" }}
    func {{template "storage.func.json_model_to_service" $model.Name}}({{$fromArgName}} *{{$jsonType.Ref}}) ({{$servicePtrTypeRef}}, error){
    result := &{{$modelGoType.Ref}}{}
    {{- range $field := $model.Fields }}
        //{{$field.Name}}
        {{- if $field.DoNotPersists }} {{- continue }} {{- end }}
        {{- $dbFieldName := include "storage.field.json_model" $field -}}
        {{- $input := print $fromArgName "." $dbFieldName -}}
        {{- $inputGoType := (goType $field.Type).DBAlternative -}}
        {{- $output := print "result." $field.Name -}}
        {{- $outputGoType := goType $field.Type -}}
        {{ template "storage.block.convert_value_to_service" (list $input $inputGoType $output $outputGoType $varNamesGenerator) }}
    {{- end }}
    return result, nil
    }
{{- end}}
