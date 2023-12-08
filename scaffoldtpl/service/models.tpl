{
"file_path": "{{.Module}}/{{.Module}}/gen.models.go",
"package_name": "{{.Module}}",
"package_path": "{{.Config.RootPackageName}}/{{.Module}}/{{.Module}}",
"condition": "len(Config.Modules[Module].Value.Types.Models) > 0"
}
<><><>
{{ $module := (index $.Config.Modules $.Module).Value}}

{{- range $oneOf := $module.Types.OneOfs}}
  {{- template "oneOfType" $oneOf }}
{{- end }}

{{- range $model := $module.Types.Models }}
  {{- if not $model.DoNotPersists }}
    {{- template "service._model_fields_enum" $model }}
    {{- template "modelFilterType" $model }}
  {{- end}}

  {{- template "modelType" $model }}
  {{- template "modelCopyHelper" $model }}
{{- end }}
