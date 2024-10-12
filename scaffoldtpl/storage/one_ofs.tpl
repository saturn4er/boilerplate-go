{
"file_path": "{{.Module}}/{{.Module}}storage/gen.one_ofs.go",
"package_name": "{{.Module}}storage",
"package_path": "{{.Config.RootPackageName}}/{{.Module}}/{{.Module}}storage",
"condition": "len(Config.Modules[Module].Value.Types.OneOfs) > 0"
}
<><><>
{{- $dbutil := import "github.com/saturn4er/boilerplate-go/lib/dbutil" }}
{{- $jsonPkg := import "encoding/json" }}
{{- $driverPkg := import "database/sql/driver"}}
{{- $fmtPkg := import "fmt" }}
{{ $module := (index $.Config.Modules $.Module).Value}}


{{- range $module.Types.OneOfs }}
    {{$goType := goType .}}
    {{$dbTypeName := print "json" .Name}}
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
        {{- $dbType := (goType $value.Value.Model).InLocalPackage.WithName (print "json" (goType $value.Value.Model).Type) }}
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

    func ({{$receiverName}} {{$dbTypeName}}) Value() ({{$driverPkg.Ref "Value"}}, error) {
    return {{$jsonPkg.Ref "Marshal"}}({{$receiverName}})
    }

    func {{template "storage.func.one_of_to_internal" .}}(val {{$goType.Ref}}) (*{{$dbTypeName}}, error) {
    if val == nil {
    return nil, nil
    }
    result := &{{$dbTypeName}}{}
    switch v := val.(type) {
    {{- range $value := .SortedValues }}
        {{- $dbType := (goType $value.Value.Model).InLocalPackage.WithName (print "json" (goType $value.Value.Model).Type) }}
        case *{{ (goType $value.Value).Ref }}:
        {{$valueOutput := print "result.Val"}}
        {{- template "storage.block.convert_value_to_internal" list "v" (goType $value.Value).Ptr $valueOutput $dbType varNamesGenerator}}

        result.OneOfType = "{{ (goType $value.Value).Ref }}"
        result.OneOfTypeID = {{$value.Index}}
    {{- end }}
    }
    return nil, {{$fmtPkg.Ref "Errorf"}}("invalid {{.Name}} value type: %T", val)
    }

    func {{template "storage.func.one_of_to_service" .}}(val *{{$dbTypeName}})  ({{$goType.Ref}}, error) {
    if val == nil {
    return nil, nil
    }

    switch v := (*val).Val.(type) {
    {{- range $value := .SortedValues }}
        {{- $dbType := (goType $value.Value.Model).InLocalPackage.WithName (print "json" (goType $value.Value.Model).Type) }}
        {{- $model := $value.Value.Model }}
        case *{{ $dbType.Ref }}:
        v1, err := {{template "storage.func.json_model_to_service" $model.Name}}(v)
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
