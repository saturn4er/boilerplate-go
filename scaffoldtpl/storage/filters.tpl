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
type filterOptions struct {
  columnPrefix string
}
func withFilterColumnPrefix(prefix string) func(*filterOptions) {
  return func(f *filterOptions) {
    f.columnPrefix = prefix
  }
}
{{- range $model := $module.Types.Models }}
  {{- if $model.DoNotPersists }}
      {{continue}}
  {{- end }}
  func {{template "storage.func.build_db_filter" $model.Name}}(filter *{{$servicePkg.Ref (print $model.Name "Filter")}}, options ...func(*filterOptions)) ({{ $clause.Ref "Expression"}}, error) {
  if filter == nil {
  return nil, nil
  }

  opts := &filterOptions{}
  for _, opt := range options {
    opt(opts)
  }

  return {{$dbutil.Ref "BuildFilterExpression"}}(
  {{- range $field := $model.Fields }}
      {{- $fieldGoType := goType $field.Type}}
      {{- if $field.Filterable }}
          {{- if $fieldGoType.IsSlice }}
            {{- if isEnum (goType $field.Type.ElemType) }}
              {{$dbutil.Ref "MappedColumnArrayFilter"}}[{{(goType $field.Type.ElemType).Ref}}, string]{
                Column: opts.columnPrefix+"{{$field.DBName}}",
                Filter: filter.{{$field.Name}},
                Mapper: {{ template "storage.func.enum_to_internal" $field.Type.ElemType.Type}},
              },
            {{- else if (goType $field.Type.ElemType).IsPtr }}
              func() {{$dbutil.Ref "ColumnFilter"}}[any] {
                panic("Slice of pointers is not implemented in generator, so we can't handle '{{$field.DBName}}' filter")
              }(),
            {{- else }}
              {{$dbutil.Ref "ColumnArrayFilter"}}[{{(goType $field.Type.ElemType).Ref}}]{
                Column: opts.columnPrefix+"{{$field.DBName}}",
                Filter: filter.{{$field.Name}},
              },
            {{- end }}
          {{- else -}}
            {{- if isEnum $fieldGoType }}
              {{$dbutil.Ref "MappedColumnFilter"}}[{{$fieldGoType.Ref}}, string]{
                Column: opts.columnPrefix+"{{$field.DBName}}",
                Filter: filter.{{$field.Name}},
                Mapper: {{ template "storage.func.enum_to_internal" $field.Type.Type}},
              },
            {{- else if and $fieldGoType.IsPtr (isEnum $fieldGoType.ElemType) }}
                {{$dbutil.Ref "MappedColumnFilter"}}[{{$fieldGoType.Ref}}, *string]{
                Column: opts.columnPrefix+"{{$field.DBName}}",
                Filter: filter.{{$field.Name}},
                Mapper: func(val *{{$fieldGoType.ElemType.Ref}}) (*string, error) {
                  if val == nil{
                    return nil, nil
                  }
                  result, err := {{ template "storage.func.enum_to_internal" $field.Type.ElemType.Type}}(*val)
                  if err != nil {
                    return nil, err
                  }
                  return &result, nil
                  },
                },
            {{- else}}
              {{$dbutil.Ref "ColumnFilter"}}[{{$fieldGoType.Ref}}]{
                Column: opts.columnPrefix+"{{$field.DBName}}",
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
