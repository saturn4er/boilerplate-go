{
"file_path": "{{.Module}}/{{.Module}}storage/gen.filters.go",
"package_name": "{{.Module}}storage",
"package_path": "{{.Config.RootPackageName}}/{{.Module}}/{{.Module}}storage",
"condition": "len(Config.Modules[Module].Value.Types.Models) > 0"
}
<><><>

{{ $module := (index $.Config.Modules $.Module).Value}}
{{ $dbutil := import "github.com/saturn4er/boilerplate-go/lib/dbutil" }}
{{ $clause := import "gorm.io/gorm/clause" }}
{{ $servicePkg :=  import (print $.Config.RootPackageName "/" $.Module "/" $.Module "service") }}

{{- range $model := $module.Types.Models }}
  {{- if $model.DoNotPersists }}
      {{continue}}
  {{- end }}
  func {{template "storage.func.build_db_filter" $model.Name}}(filter *{{$servicePkg.Ref (print $model.Name "Filter")}}) ({{ $clause.Ref "Expression"}}, error) {
  if filter == nil {
  return nil, nil
  }

  return {{$dbutil.Ref "BuildFilterExpression"}}(
  {{- range $field := $model.Fields }}
      {{- $fieldGoType := goType $field.Type}}
      {{- if $field.Filterable }}
          {{- if $fieldGoType.IsSlice }}
            {{- if or (isModuleEnum (goType $field.Type.ElemType)) (isCommonEnum (goType $field.Type.ElemType)) }}
              {{$dbutil.Ref "MappedColumnArrayFilter"}}[{{(goType $field.Type.ElemType).Ref}}, string]{
                Column: "{{$field.DBName}}",
                Filter: filter.{{$field.Name}},
                Mapper: {{ template "storage.func.enum_to_internal" $field.Type.ElemType.Type}},
              },
            {{- else if (goType $field.Type.ElemType).IsPtr }}
              func() {{$dbutil.Ref "ColumnFilter"}}[any] {
                panic("Slice of pointers is not implemented in generator, so we can't handle '{{$field.DBName}}' filter")
              }(),
            {{- else }}
              {{$dbutil.Ref "ColumnArrayFilter"}}[{{(goType $field.Type.ElemType).Ref}}]{
                Column: "{{$field.DBName}}",
                Filter: filter.{{$field.Name}},
              },
            {{- end }}
          {{- else -}}
            {{- if or (isModuleEnum (goType $field.Type)) (isCommonEnum (goType $field.Type)) }}
              {{$dbutil.Ref "MappedColumnFilter"}}[{{$fieldGoType.Ref}}, string]{
                Column: "{{$field.DBName}}",
                Filter: filter.{{$field.Name}},
                Mapper: {{ template "storage.func.enum_to_internal" $field.Type.Type}},
              },
            {{- else}}
              {{$dbutil.Ref "ColumnFilter"}}[{{$fieldGoType.Ref}}]{
                Column: "{{$field.DBName}}",
                Filter: filter.{{$field.Name}},
              },
            {{- end }}
          {{- end }}
      {{- end }}
  {{- end }}
    {{$dbutil.Ref "ExpressionBuilderFunc"}}(func() (clause.Expression, error) {
      if filter.Or == nil {
        return nil, nil
      }
      exprs := make([]clause.Expression, 0, len(filter.Or))
      for _, orFilter := range filter.Or {
        expr, err := {{template "storage.func.build_db_filter" $model.Name}}(orFilter)
        if err != nil {
          return nil, err
        }
        exprs = append(exprs, expr)
      }
      return clause.Or(exprs...), nil
    }),
    {{$dbutil.Ref "ExpressionBuilderFunc"}}(func() (clause.Expression, error) {
      if filter.And == nil {
        return nil, nil
      }
      exprs := make([]clause.Expression, 0, len(filter.And))
      for _, andFilter := range filter.And {
        expr, err := {{template "storage.func.build_db_filter" $model.Name}}(andFilter)
        if err != nil {
          return nil, err
        }
        exprs = append(exprs, expr)
      }
      return clause.And(exprs...), nil
    }),
  )
  }

{{- end }}
