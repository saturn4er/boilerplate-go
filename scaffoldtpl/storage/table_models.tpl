{
"file_path": "{{.Module}}/{{.Module}}storage/gen.table_models.go",
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
    {{- $fieldType := (goType .Type).DBAlternative -}}
    `gorm:"column:{{- .DBName -}};{{- if $fieldType.GormType -}}type:{{- $fieldType.GormType -}};{{- end -}}{{if .PrimaryKey}}primaryKey{{end}}"`
{{- end -}}


{{- range $model := $module.Types.Models }}
    {{- if $model.DoNotPersists }}
        {{continue}}
    {{- end }}

    {{- $modelGoType := goType $model }}
    {{- $servicePtrType := $modelGoType.Ptr -}}
    {{- $servicePtrTypeRef := $servicePtrType.Ref -}}
    {{- $dbType := $modelGoType.InLocalPackage.WithName (print "db" $modelGoType.Type)  -}}
    {{- $jsonType := $modelGoType.InLocalPackage.WithName (print "json" $modelGoType.Type)  -}}
    {{- $dbTypeRef := $dbType.Ref -}}
    {{- $dbPtrType := $dbType.Ptr -}}
    {{- $dbPtrTypeRef := print $dbPtrType.Ref -}}
    {{- $varNamesGenerator := varNamesGenerator }}

    type {{$dbType.Ref}} struct {
    {{- range $field := $model.Fields}}
        {{- if not $field.DoNotPersists }}
            {{- $fieldType := (goType $field.Type).DBAlternative }}
            {{$field.Name}} {{(goType $field.Type).DBAlternative.Ref}} {{ template "fieldGormTag" $field }}
        {{- end }}
    {{- end }}
    }

    {{ $toArgName := "src" }}
    func {{ template "storage.func.table_model_to_internal" $model.Name}}({{$toArgName}} {{$servicePtrTypeRef}}) ({{$dbPtrTypeRef}}, error){
    result := &{{$dbTypeRef}}{}
    {{- range $field := $model.Fields }}
        {{- if $field.DoNotPersists }} {{- continue }} {{- end }}
        {{- $input := print $toArgName "." $field.Name -}}
        {{- $inputGoType := goType $field.Type -}}
        {{- $output := setToVarFn (print "result." $field.Name) -}}
        {{- $outputGoType := (goType $field.Type).DBAlternative -}}

        {{ template "storage.block.convert_value_to_internal" (list $input $inputGoType $output $outputGoType $varNamesGenerator) }}
    {{- end }}
    return result, nil
    }

    {{ $fromArgName := "src" }}
    func {{template "storage.func.table_model_to_service" $model.Name}}({{$fromArgName}} {{$dbPtrTypeRef}}) ({{$servicePtrTypeRef}}, error){
    result := &{{$modelGoType.Ref}}{}
    {{- range $field := $model.Fields }}
        {{- if $field.DoNotPersists }} {{- continue }} {{- end }}
        {{- $input := print $fromArgName "." $field.Name -}}
        {{- $inputGoType := (goType $field.Type).DBAlternative -}}
        {{- $output := setToVarFn (print "result." $field.Name) -}}
        {{- $outputGoType := goType $field.Type -}}
        {{ template "storage.block.convert_value_to_service" (list $input $inputGoType $output $outputGoType $varNamesGenerator) }}
    {{- end }}
    return result, nil
    }
    {{- if $model.TableName }}
      func (a {{$dbTypeRef}}) TableName() string {
      return "{{$model.TableName}}"
      }
    {{- end }}
{{- end}}
