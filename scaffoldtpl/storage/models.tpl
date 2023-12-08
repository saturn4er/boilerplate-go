{
"file_path": "{{.Module}}/{{.Module}}storage/gen.models.go",
"package_name": "{{.Module}}storage",
"package_path": "{{.Config.RootPackageName}}/{{.Module}}/{{.Module}}storage",
"condition": "len(Config.Modules[Module].Value.Types.Models) > 0"
}
<><><>

{{ $dbutil := import "github.com/saturn4er/boilerplate-go/lib/dbutil" }}
{{ $clause := import "gorm.io/gorm/clause" }}
{{ $servicePkg :=  import (print $.Config.RootPackageName "/" $.Module "/" $.Module) }}
{{ $module := (index $.Config.Modules $.Module).Value}}


{{- define "fieldGormTag" -}}
{{- $fieldType := (goType .Type).DBAlternative }}
{{- if or $fieldType.GormType .PrimaryKey -}}
  `gorm:"{{- if $fieldType.GormType -}}type:{{- $fieldType.GormType -}};{{- end -}}{{if .PrimaryKey}}primaryKey;autoIncrement:false{{end}}"`
{{- else -}}{{- end -}}
{{- end -}}

{{- range $model := $module.Types.Models }}
    {{- if $model.DoNotPersists }} {{- continue}} {{- end }}
    {{- $modelGoType := goType $model }}

    {{- $servicePtrType := $modelGoType.Ptr -}}
    {{- $servicePtrTypeRef := $servicePtrType.Ref -}}
    {{- $dbType := $modelGoType.InLocalPackage.WithName (include "storage.type.db_model" $model)  -}}
    {{- $dbTypeRef := $dbType.Ref -}}
    {{- $dbPtrType := $dbType.Ptr -}}
    {{- $dbPtrTypeRef := print $dbPtrType.Ref -}}
    {{- $varNamesGenerator := varNamesGenerator }}

    type {{$dbTypeRef}} struct {
    {{- range $field := $model.Fields}}
        {{- if not $field.DoNotPersists }}
            {{- $fieldType := (goType $field.Type).DBAlternative }}
            {{$field.Name}} {{(goType $field.Type).DBAlternative.Ref}} {{ template "fieldGormTag" $field }}
        {{- end }}
    {{- end }}
    }
    {{- if $model.DoNotPersists }}
        {{continue}}
    {{- end }}
    {{- if $model.TableName }}
      func (a {{$dbTypeRef}}) TableName() string {
      return "{{$model.TableName}}"
      }
    {{- end }}

    {{ $toArgName := print $.Module $model.Name }}
    func {{ template "storage.func.model_to_db_model" $model }}({{$toArgName}} {{$servicePtrTypeRef}}) ({{$dbPtrTypeRef}}, error){
    result := &{{$dbTypeRef}}{}
    {{- range $field := $model.Fields }}
        {{- if $field.DoNotPersists }} {{- continue }} {{- end }}
        {{- $input := print $toArgName "." $field.Name -}}
        {{- $inputGoType := goType $field.Type -}}
        {{- $output := print "result." $field.Name -}}
        {{- $outputGoType := (goType $field.Type).DBAlternative -}}

        {{ template "storage._convert_value_to_db_value" (list $input $inputGoType $output $outputGoType $varNamesGenerator) }}
    {{- end }}
    return result, nil
    }

    {{ $fromArgName := print "db" $model.Name }}
    func {{template "storage.func.db_model_to_model" $model.Name}}({{$fromArgName}} {{$dbPtrTypeRef}}) ({{$servicePtrTypeRef}}, error){
    result := &{{$modelGoType.Ref}}{}
    {{- range $field := $model.Fields }}
        {{- if $field.DoNotPersists }} {{- continue }} {{- end }}
        {{- $input := print $fromArgName "." $field.Name -}}
        {{- $inputGoType := (goType $field.Type).DBAlternative -}}
        {{- $output := print "result." $field.Name -}}
        {{- $outputGoType := goType $field.Type -}}
        {{ template "storage._convert_db_value_to_value" (list $input $inputGoType $output $outputGoType $varNamesGenerator) }}
    {{- end }}
    return result, nil
    }
{{- end}}
