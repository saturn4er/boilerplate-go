{
"file_path": "{{.Module}}/{{.Module}}/gen.enums.go",
"package_name": "{{.Module}}",
"package_path": "{{.Config.RootPackageName}}/{{.Module}}/{{.Module}}",
"condition": "len(Config.Modules[Module].Value.Types.Enums) > 0"
}
<><><>

{{ $module := (index $.Config.Modules $.Module).Value}}

{{- range $enum := $module.Types.Enums }}
  {{ template "service._enum_all_values_helper" $enum }}
  {{ template "service._enum_is_category_helper" $enum }}
  {{ template "service._enum_declaration" $enum }}
  {{ template "service._enum_is_valid_helper" $enum }}
  {{ template "service._enum_is_helper" $enum }}
  {{ template "service._enum_stringer_helper" $enum }}
  {{ template "service._enum_validate_helper" $enum }}
  {{ userBlock $enum.Name }}
{{- end }}
