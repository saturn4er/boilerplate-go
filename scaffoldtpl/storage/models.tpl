{
"file_path": "{{.Module}}/{{.Module}}storage/gen.models.go",
"package_name": "{{.Module}}storage",
"package_path": "{{.Config.RootPackageName}}/{{.Module}}/{{.Module}}storage",
"condition": "len(Config.Modules[Module].Value.Types.Models) > 0"
}
<><><>

{{- $dbutil := import "github.com/saturn4er/boilerplate-go/lib/dbutil" }}
{{- $jsonPkg := import "encoding/json" }}
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


{{- range $module.Types.OneOfs }}
{{$goType := goType .}}
{{$dbTypeName := print "db" .Name}}
{{$receiverName := .Name | receiverName }}
type {{$dbTypeName}} struct {
  Val any `json:"value"`
  OneOfType   string `json:"@type"`
  OneOfTypeID uint `json:"@type_id"`
}
func ({{$receiverName}}  *{{$dbTypeName}}) UnmarshalJSON(bytes []byte) error {
  tmp := struct {
    OneOfTypeID uint   `json:"@type_id"`
    OneOfType   string `json:"@type"`
  }{}
  if err := {{$jsonPkg.Ref "Unmarshal"}}(bytes, &tmp); err != nil {
    return {{$fmtPkg.Ref "Errorf"}}("unmarshal OneOfType: %w", err)
  }

  switch tmp.OneOfTypeID {
    {{- range $value := .SortedValues }}
      {{- $dbType := (goType $value.Value.Model).InLocalPackage.WithName (print "db" (goType $value.Value.Model).Type) }}
    case {{$value.Index}}:
      var value struct {
      Value {{$dbType.Ref}}  `json:"value"`
      }
      if err := {{$jsonPkg.Ref "Unmarshal"}}(bytes, &value); err != nil {
      return err
      }
      {{$receiverName}}.Val = &value.Value
    {{- end }}
  }
  return nil
}
func ({{$receiverName}} *{{$dbTypeName}}) Scan(value any) error {
  return {{$jsonPkg.Ref "Unmarshal"}}(value.([]byte), {{$receiverName}})
}

func ({{$receiverName}} {{$dbTypeName}}) Value() (driver.Value, error) {
  return {{$jsonPkg.Ref "Marshal"}}({{$receiverName}})
}

func {{template "storage.func.one_of_to_internal" .}}(val {{$goType.Ref}}) (*{{$dbTypeName}}, error) {
  if val == nil {
    return nil, nil
  }
  result := &{{$dbTypeName}}{}
  switch v := val.(type) {
  {{- range $value := .SortedValues }}
      {{- $dbType := (goType $value.Value.Model).InLocalPackage.WithName (print "db" (goType $value.Value.Model).Type) }}
      case *{{ (goType $value.Value).Ref }}:
      {{$valueOutput := print "result.Val"}}
      {{- template "storage.block.convert_value_to_internal" list "v" (goType $value.Value).Ptr $valueOutput $dbType varNamesGenerator}}

      result.OneOfType = "{{ (goType $value.Value).Ref }}"
      result.OneOfTypeID = {{$value.Index}}
  {{- end }}
  }
  return nil, {{$fmtPkg.Ref "Errorf"}}("invalid {{.Name}} value type: %T", val)
}
func {{template "storage.func.one_of_to_string_ptr" .}}(val {{$goType.Ref}}) (*string, error) {
  if val == nil {
    return nil, nil
  }

  dbVal, err := {{template "storage.func.one_of_to_internal" .}}(val)
  if err != nil {
    return nil, {{$fmtPkg.Ref "Errorf"}}("convert {{.Name}} to internal: %w", err)
  }

  encodedVal, err := {{$jsonPkg.Ref "Marshal"}}(dbVal)
  if err!=nil {
    return nil, {{$fmtPkg.Ref "Errorf"}}("marshal {{$dbTypeName}}: %w", err)
  }

  strVal := string(encodedVal)

  return &strVal, nil
}
func {{template "storage.func.one_of_to_service" .}}(val *{{$dbTypeName}})  ({{$goType.Ref}}, error) {
  if val == nil {
    return nil, nil
  }

  switch v := (*val).Val.(type) {
  {{- range $value := .SortedValues }}
      {{- $dbType := (goType $value.Value.Model).InLocalPackage.WithName (print "db" (goType $value.Value.Model).Type) }}
      {{- $model := $value.Value.Model }}
      case *{{ $dbType.Ref }}:
      v1, err := {{template "storage.func.model_to_service" $model.Name}}(v)
      if err!=nil {
      return nil, {{$fmtPkg.Ref "Errorf"}}("convert {{$model.Name}} from db: %w", err)
      }

      return v1, nil
  {{- end }}
  default:
  return nil, {{$fmtPkg.Ref "Errorf"}}("invalid {{.Name}} value type: %T", *val)
  }

  panic("implement me")
}
func {{template "storage.func.one_of_string_ptr_to_service" .}}(val *string) ({{$goType.Ref}}, error) {
  if val == nil || *val == "null" {
    return nil, nil
  }

  var dbVal {{$dbTypeName}}
  if err := {{$jsonPkg.Ref "Unmarshal"}}([]byte(*val), &dbVal); err != nil {
    return nil, {{$fmtPkg.Ref "Errorf"}}("unmarshal {{$dbTypeName}}: %w", err)
  }

  return {{template "storage.func.one_of_to_service" .}}(&dbVal)
}
{{- end }}


{{- range $model := $module.Types.Models }}
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

    type {{$jsonType.Ref}} struct {
    {{- range $field := $model.Fields}}
        {{- if not $field.DoNotPersists }}
            {{- $fieldType := (goType $field.Type).DBAlternative }}
            {{$field.Name}} {{(goType $field.Type).DBAlternative.Ref}}
        {{- end }}
    {{- end }}
    }

    {{ $toArgName := "src" }}
    func {{ template "storage.func.model_to_internal" $model.Name}}({{$toArgName}} {{$servicePtrTypeRef}}) ({{$dbPtrTypeRef}}, error){
    result := &{{$dbTypeRef}}{}
    {{- range $field := $model.Fields }}
        {{- if $field.DoNotPersists }} {{- continue }} {{- end }}
        {{- $input := print $toArgName "." $field.Name -}}
        {{- $inputGoType := goType $field.Type -}}
        {{- $output := print "result." $field.Name -}}
        {{- $outputGoType := (goType $field.Type).DBAlternative -}}

        {{ template "storage.block.convert_value_to_internal" (list $input $inputGoType $output $outputGoType $varNamesGenerator) }}
    {{- end }}
    return result, nil
    }

    {{ $fromArgName := "src" }}
    func {{template "storage.func.model_to_service" $model.Name}}({{$fromArgName}} {{$dbPtrTypeRef}}) ({{$servicePtrTypeRef}}, error){
    result := &{{$modelGoType.Ref}}{}
    {{- range $field := $model.Fields }}
        {{- if $field.DoNotPersists }} {{- continue }} {{- end }}
        {{- $input := print $fromArgName "." $field.Name -}}
        {{- $inputGoType := (goType $field.Type).DBAlternative -}}
        {{- $output := print "result." $field.Name -}}
        {{- $outputGoType := goType $field.Type -}}
        {{ template "storage.block.convert_value_to_service" (list $input $inputGoType $output $outputGoType $varNamesGenerator) }}
    {{- end }}
    return result, nil
    }
    {{- if $model.DoNotPersists }}
      {{continue}}
    {{- end }}
    {{- if $model.TableName }}
      func (a {{$dbTypeRef}}) TableName() string {
      return "{{$model.TableName}}"
      }
    {{- end }}
{{- end}}
