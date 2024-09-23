{
"file_path": "{{.Module}}/{{.Module}}storage/gen.admin_table_generators.go",
"package_name": "{{.Module}}storage",
"package_path": "{{.Config.RootPackageName}}/{{.Module}}/{{.Module}}storage",
"condition": "len(Config.Modules[Module].Value.Types.Models) > 0"
}
<><><>
{{ $module := (index $.Config.Modules $.Module).Value}}
{{ $adminContextPkg :=  import "github.com/GoAdminGroup/go-admin/context" }}
{{ $adminDBPkg := import "github.com/GoAdminGroup/go-admin/modules/db" }}
{{ $adminTypesPkg := import "github.com/GoAdminGroup/go-admin/template/types"}}
{{ $formPkg := import "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"}}
{{ $stringsPkg := import "strings" }}
{{ $urlPkg := import "net/url" }}

{{ $tablePkg :=  import "github.com/GoAdminGroup/go-admin/plugins/admin/modules/table" }}
func NewTableGenerators() {{$tablePkg.Ref "GeneratorList"}} {
return map[string]{{$tablePkg.Ref "Generator"}}{
{{- range $model := $module.Types.Models }}
    {{- if or $model.DoNotPersists (eq $model.StorageType "tx_outbox") }}
        {{continue}}
    {{- else }}
      "{{ $model.Name | snakeCase }}": func(ctx *{{$adminContextPkg.Ref "Context"}}) {{$tablePkg.Ref "Table"}} {
      tableConfig :=  {{$tablePkg.Ref "DefaultConfigWithDriver"}}("postgresql")
      {{- $pkField := $model.FirstPKField }}
      tableConfig.PrimaryKey.Type = {{ (goType $pkField.Type).GoAdminType }}
      tableConfig.PrimaryKey.Name = "{{$pkField.Name | snakeCase}}"
      table := {{$tablePkg.Ref "NewDefaultTable"}}(ctx, tableConfig)
      info := table.GetInfo()
      formList := table.GetForm()
      info.SetTable("{{replace $model.TableName "." "\\\".\\\""}}").SetTitle("{{$model.Name}}").SetDescription("{{$model.Name}}")
      formList.SetTable("{{replace $model.TableName "." "\\\".\\\""}}").SetTitle("{{$model.Name}}").SetDescription("{{$model.Name}}")
      {{- range $field := $model.Fields }}
        {{- $fieldGoType := goType $field.Type}}
        info.AddField("{{$field.Name}}", "{{$field.Name | snakeCase}}", {{$fieldGoType.GoAdminType }})
        info.FieldSortable()
        {{- if $field.Admin.LinkTo }}
          info.FieldDisplay(func(model types.FieldModel) interface{} {
            return template.Default().Link().
              SetContent(template.HTML(model.Value)).
              SetURL("/admin/info/{{$field.Admin.LinkTo}}/detail?__goadmin_detail_pk="+model.Value).
              GetContent()
          })
        {{- end }}
        {{- if $field.Admin.HideForList }}
          info.FieldHideForList()
        {{- end }}
        {{- if $field.Admin.Hide }}
          info.FieldHide()
        {{- end }}
        {{- if $field.Filterable }}
          info.FieldFilterable({{$adminTypesPkg.Ref "FilterType" }}{
            FormType: {{$fieldGoType.GoAdminForm }},
            {{- if isEnum $fieldGoType }}
                {{$enum := getEnum $fieldGoType }}
                Options: {{$adminTypesPkg.Ref "FieldOptions"}}{
                {{- range $value := $enum.Values }}
                  {Value: "{{$value | snakeCase}}", Text: "{{$value}}"},
                {{- end}}
                },
            {{- else if $fieldGoType.IsPtr }}
                {{- if isEnum $fieldGoType.ElemType }}
                    {{$enum := getEnum $fieldGoType.ElemType }}
                    Options: {{$adminTypesPkg.Ref "FieldOptions"}}{
                    {{- range $value := $enum.Values }}
                      {Value: "{{$value | snakeCase}}", Text: "{{$value}}"},
                    {{- end}}
                    },
                {{- end}}
            {{- end }}
            },
          )
        {{- end }}
        formList.AddField("{{$field.Name}}", "{{$field.Name | snakeCase}}", {{$fieldGoType.GoAdminType }}, {{$fieldGoType.GoAdminForm }})
        {{- if not $field.Admin.Editable }}
          formList.FieldDisableWhenUpdate()
        {{- end }}
        {{- if not $field.Admin.Creatable }}
          formList.FieldHideWhenCreate() // {{$field.Admin}}
        {{- end }}
        formList.PreProcessFn = func(values {{$formPkg.Ref "Values"}}) {{$formPkg.Ref "Values"}} {
          for k, v := range values {
          	for i, v := range v {
          		if {{$stringsPkg.Ref "Contains"}}(v, "%") {
          			if newV, err := {{$urlPkg.Ref "QueryUnescape"}}(v); err == nil {
          				values[k][i] = newV
          			}
          		}
          	}
          }
          return values
        }
        {{- if isEnum $fieldGoType }}
          {{$enum := getEnum $fieldGoType }}
          formList.FieldOptions({{$adminTypesPkg.Ref "FieldOptions"}}{
            {{- range $value := $enum.Values }}
                {Value: "{{$value | snakeCase}}", Text: "{{$value}}"},
            {{- end}}
          })
        {{- else if $fieldGoType.IsPtr }}
            {{- if isEnum $fieldGoType.ElemType }}
                {{$enum := getEnum $fieldGoType.ElemType }}
                formList.FieldOptions({{$adminTypesPkg.Ref "FieldOptions"}}{
                {{- range $value := $enum.Values }}
                  {Value: "{{$value | snakeCase}}", Text: "{{$value}}"},
                {{- end}}
                })
            {{- end}}
        {{- end }}
      {{- end }}


{{/*      // set id editable is false.*/}}
{{/*      formList.AddField("Id", "id", db.UUID, form.Default).FieldDefault(uuid.New().String()).FieldDisplayButCanNotEditWhenUpdate()*/}}
{{/*      formList.AddField("Email", "email", db.Varchar, form.Email).FieldDisplayButCanNotEditWhenUpdate()*/}}
{{/*      formList.AddField("First name", "first_name", db.Varchar, form.Text)*/}}
{{/*      formList.AddField("Last name", "last_name", db.Varchar, form.Text)*/}}
{{/*      formList.AddField("Created at", "created_at", db.Timestamp, form.Datetime)*/}}
{{/*      formList.AddField("Updated at", "updated_at", db.Timestamp, form.Datetime)*/}}


      return table
      },
    {{- end }}
{{- end }}
}
}
