{
"file_path": "{{.Module}}/{{.Module}}storage/gen.storages.go",
"package_name": "{{.Module}}storage",
"package_path": "{{.Config.RootPackageName}}/{{.Module}}/{{.Module}}storage",
"condition": "len(Config.Modules[Module].Value.Types.Models) > 0"
}
<><><>

{{ $dbutilPkg := import "github.com/saturn4er/boilerplate-go/lib/dbutil" }}
{{ $loggingPkg := import "github.com/go-pnp/go-pnp/logging" }}
{{ $txoutboxPkg := import "github.com/saturn4er/boilerplate-go/lib/txoutbox" }}
{{ $idempotencyPkg := import "github.com/saturn4er/boilerplate-go/lib/idempotency" }}
{{ $optionUtilPkg := import "github.com/go-pnp/go-pnp/pkg/optionutil"}}
{{ $fmtPkg := import "fmt" }}
{{ $gormPkg := import "gorm.io/gorm" }}
{{ $sqlPkg := import "database/sql" }}
{{ $sqlxPkg := import "github.com/jmoiron/sqlx" }}
{{ $goquPkg := import "github.com/doug-martin/goqu/v9" "goqu"}}
{{ $errorsPkg := import "errors" }}
{{ $servicePkg :=  import (print $.Config.RootPackageName "/" $.Module "/" $.Module) (print $.Module "svc") }}
{{ $module := (index $.Config.Modules $.Module).Value}}

type Storages struct {
db *{{$sqlxPkg.Ref "DB"}}
connection {{$dbutilPkg.Ref "Connection"}}
dialect {{$goquPkg.Ref "DialectWrapper"}}
isTransaction bool
logger *{{$loggingPkg.Ref "Logger"}}
}

var _ {{$servicePkg.Ref "Storage"}} = &Storages{}

{{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold.ConfigModel*/ -}}
{{- range $model := $module.Types.Models }}
    {{- if $model.DoNotPersists }}
        {{continue}}
    {{- end }}
    {{- if eq $model.StorageType "tx_outbox" }}
        {{/*    func (s Storages) {{$model.PluralName}}() {{$servicePkg.Ref (print $model.PluralName "Outbox")}} {*/}}
        {{/*      return New{{$model.PluralName}}Outbox(s.connection)*/}}
        {{/*    }*/}}
    {{- end }}
{{- end }}

{{/*func (s Storages) IdempotencyKeys() {{$idempotencyPkg.Ref "Storage"}} {*/}}
{{/*  return {{$idempotencyPkg.Ref "GormStorage"}}{*/}}
{{/*    Connection: s.connection,*/}}
{{/*  }*/}}
{{/*}*/}}

func (s Storages) ExecuteInTransaction(ctx context.Context, cb func(ctx context.Context, tx {{$servicePkg.Ref "Storage"}}) error, txOptions ...{{$dbutilPkg.Ref "TxOption"}}) (rerr error) {
sqlTxOptions := &{{$sqlPkg.Ref "TxOptions"}}{}
for _, opt := range txOptions {
if err := opt(sqlTxOptions); err != nil {
return err
}
}
tx, err := s.db.BeginTxx(ctx, sqlTxOptions)
if err != nil {
return {{$fmtPkg.Ref "Errorf"}}("begin transaction: %w", err)
}

defer func() {
if rerr == nil {
return
}
if rollbackErr := tx.Rollback(); rollbackErr != nil {
rerr = {{$errorsPkg.Ref "Join"}}(rerr, {{$fmtPkg.Ref "Errorf"}}("rollback: %w", rollbackErr))
}
}()

if err := cb(ctx, &Storages{db: s.db, dialect: s.dialect, connection: tx, isTransaction: true, logger: s.logger}); err != nil {
return {{$fmtPkg.Ref "Errorf"}}("execute transaction callback: %w", err)
}

if err := tx.Commit(); err != nil {
return {{$fmtPkg.Ref "Errorf"}}("commit transaction: %w", err)
}

return nil
}

func NewStorages(db *{{$sqlxPkg.Ref "DB"}}, logger *{{$loggingPkg.Ref "Logger"}}) *Storages {
return &Storages{db: db, logger: logger}
}

{{- range $model := $module.Types.Models }}
    {{- if $model.DoNotPersists }}
        {{continue}}
    {{- end }}
    {{- if eq $model.StorageType "tx_outbox" }}
        {{/*        func New{{$model.PluralName}}Outbox(db *{{$gormPkg.Ref "DB"}}) {{$servicePkg.Ref (print $model.PluralName "Outbox")}} {*/}}
        {{/*          return {{$txoutboxPkg.Ref "GormStorage"}}[{{$servicePkg.Ref $model.Name}}]{*/}}
        {{/*            DB: db,*/}}
        {{/*            BuildMessage:     build{{$model.Name}}Message,*/}}
        {{/*          }*/}}
        {{/*        }*/}}
    {{- else }}
        {{ $modelGoType := goType $model }}
        {{ $dbType := $modelGoType.InLocalPackage.WithName (include "storage.type.db_model" $model) }}
        {{ $dbTypeRef := $dbType.Ref }}
        {{ $storageTypeName := print $model.PluralName "Storage" }}
        const (
        {{- range $field := $model.Fields }}
            {{- if not $field.DoNotPersists }}
                {{template "storage.const.model_table_field" (dict "model" $model "field" $field)}} = "{{$field.Name | snakeCase}}"
            {{- end }}
        {{- end }}
        )

        const {{template "storage.const.model_table_name" $model}} = "{{$model.TableName | default ($model.Name | snakeCase)}}"
        {{- $tableName := include "storage.const.model_table_name" $model}}

        type {{$storageTypeName}} struct {
        connection {{$dbutilPkg.Ref "Connection"}}
        logger *{{$loggingPkg.Ref "Logger"}}
        dialect {{$goquPkg.Ref "DialectWrapper"}}
        }
        func (s Storages) {{$model.PluralName}}() {{$servicePkg.Ref (print $model.PluralName "Storage")}} {
        return New{{$storageTypeName}}(s.connection, s.logger, s.dialect)
        }
        func New{{$storageTypeName}}(
        connection {{$dbutilPkg.Ref "Connection"}},
        logger *{{$loggingPkg.Ref "Logger"}},
        dialect {{$goquPkg.Ref "DialectWrapper"}},
        ) {{$servicePkg.Ref (print $model.PluralName "Storage")}} {
        return &{{$storageTypeName}}{
        logger: logger,
        connection: connection,
        dialect: dialect,
        }
        }
        {{- $receiverName := receiverName $storageTypeName}}
        {{- $methodStart := printf "func (%s %s) " $receiverName $storageTypeName}}
        {{- $svcModelRef := $servicePkg.Ref (include "service.type.model" $model)}}
        {{- $svcModelFilterRef := $servicePkg.Ref (include "service.type.model_filter" $model)}}
        {{- $selectOptions := printf "%s[%s]" ($optionUtilPkg.Ref "Option") ($dbutilPkg.Ref "SelectOptions")}}

        {{/*            CREATE        */}}
        {{$methodStart}} Create(ctx context.Context, {{$model.Name | lCamelCase}} *{{$svcModelRef}}) (*{{$svcModelRef}}, error) {
        DB{{$model.Name}}, err := {{ template "storage.func.model_to_db_model" $model}}({{$model.Name | lCamelCase}})
        if err != nil {
        return nil, {{$fmtPkg.Ref "Errorf"}}("convert {{$model.Name | lCamelCase}} to db {{$model.Name | lCamelCase}}: %w", err)
        }

        query, params, err := {{$receiverName}}.dialect.
        Insert({{$tableName}}).
        Prepared(true).
        Rows(DB{{$model.Name}}).
        Returning(
        {{- range $field := $model.Fields }}
            {{- if not $field.DoNotPersists }}
                {{template "storage.const.model_table_field" (dict "model" $model "field" $field)}},
            {{- end }}
        {{- end }}
        ).
        ToSQL()
        if err != nil {
        return nil, {{$fmtPkg.Ref "Errorf"}}("build query: %w", err)
        }


        rows, err := {{$receiverName}}.connection.QueryxContext(ctx, query, params...)
        if err != nil {
        return nil, {{$fmtPkg.Ref "Errorf"}}("execute query: %w", {{$receiverName}}.mapError(err))
        }

        for rows.Next() {
        var DB{{$model.Name}} {{$dbTypeRef}}
        err := rows.StructScan(&DB{{$model.Name}})
        if err != nil {
        return nil, {{$fmtPkg.Ref "Errorf"}}("scan row: %w", err)
        }

        result, err := {{ template "storage.func.db_model_to_model" $model.Name}}(&DB{{$model.Name}})
        if err != nil {
        return nil, {{$fmtPkg.Ref "Errorf"}}("convert db {{$model.Name | lCamelCase}} to {{$model.Name | lCamelCase}}: %w", err)
        }
        return result, nil
        }

        return nil, {{$fmtPkg.Ref "Errorf"}}("no rows returned")
        }

        {{/*            COUNT        */}}
        {{$methodStart}} Count(ctx context.Context, filter *{{$svcModelFilterRef}}) (int, error) {
        queryBuilder := {{$receiverName}}.dialect.From({{$tableName}}).Select(goqu.COUNT("*")).Prepared(true)
        if filter != nil {
        expr, err := {{$receiverName}}.filterExpr(filter)
        if err != nil {
        return 0, fmt.Errorf("build filter expression: %w", err)
        }

        queryBuilder = queryBuilder.Where(expr)
        }

        query, params, err := queryBuilder.ToSQL()
        if err != nil {
        return 0, fmt.Errorf("build query: %w", err)
        }

        var count int
        if err = {{$receiverName}}.connection.GetContext(ctx, &count, query, params...); err != nil {
        return 0, fmt.Errorf("execute query: %w", err)
        }

        return count, nil
        }


        {{/*            UPDATE        */}}
        {{$methodStart}} Update(ctx context.Context, model *{{$svcModelRef}}) (*{{$svcModelRef}}, error) {
        queryBuilder := {{$receiverName}}.dialect.Update({{$tableName}}).Prepared(true)
        DB{{$model.Name}}, err := {{ template "storage.func.model_to_db_model" $model}}(model)
        if err != nil {
        return nil, fmt.Errorf("convert UserOIDCAccount to db userOIDCAccount: %w", err)
        }

        query, parameters, err := queryBuilder.
        Returning(
        {{- range $field := $model.Fields }}
            {{- if not $field.DoNotPersists }}
                {{template "storage.const.model_table_field" (dict "model" $model "field" $field)}},
            {{- end }}
        {{- end }}
        ).
        Set({{$goquPkg.Ref "Record"}}{
        {{- range $field := $model.Fields }}
            {{- if not $field.DoNotPersists }}
                {{template "storage.const.model_table_field" (dict "model" $model "field" $field)}}: DB{{$model.Name}}.{{$field.Name}},
            {{- end }}
        {{- end }}
        }).
        ToSQL()
        if err!=nil {
        return nil, fmt.Errorf("build query: %w", err)
        }

        rows, err := {{$receiverName}}.connection.QueryxContext(ctx, query, parameters...)
        if err != nil {
        return nil, fmt.Errorf("execute query: %w", {{$receiverName}}.mapError(err))
        }
        for rows.Next() {
        var DB{{$model.Name}} {{$dbTypeRef}}
        err := rows.StructScan(&DB{{$model.Name}})
        if err != nil {
        return nil, fmt.Errorf("scan row: %w", err)
        }

        result, err := {{ template "storage.func.db_model_to_model" $model.Name}}(&DB{{$model.Name}})
        if err != nil {
        return nil, fmt.Errorf("convert db {{$model.Name | lCamelCase}} to {{$model.Name | lCamelCase}}: %w", err)
        }
        return result, nil
        }

        return nil, fmt.Errorf("no rows returned during update")
        }

        {{/*            SAVE        */}}
        {{$methodStart}} Save(ctx context.Context, model *{{$svcModelRef}}) (*{{$svcModelRef}}, error) {
        queryBuilder := {{$receiverName}}.dialect.Insert({{$tableName}}).Prepared(true)
        DB{{$model.Name}}, err := {{ template "storage.func.model_to_db_model" $model}}(model)
        if err != nil {
        return nil, fmt.Errorf("convert UserOIDCAccount to db userOIDCAccount: %w", err)
        }

        query, parameters, err := queryBuilder.
        Returning(
        {{- range $field := $model.Fields }}
            {{- if not $field.DoNotPersists }}
                {{template "storage.const.model_table_field" (dict "model" $model "field" $field)}},
            {{- end }}
        {{- end }}
        ).
        Rows({{$goquPkg.Ref "Record"}}{
        {{- range $field := $model.Fields }}
            {{- if not $field.DoNotPersists }}
                {{template "storage.const.model_table_field" (dict "model" $model "field" $field)}}: DB{{$model.Name}}.{{$field.Name}},
            {{- end }}
        {{- end }}
        }).
        OnConflict(
        {{$goquPkg.Ref "DoUpdate"}}("{{$model.PKIndexName}}", goqu.Record{
        {{- range $field := $model.Fields }}
            {{- if not $field.DoNotPersists }}
                {{template "storage.const.model_table_field" (dict "model" $model "field" $field)}}: DB{{$model.Name}}.{{$field.Name}},
            {{- end }}
        {{- end }}
        }),
        ).
        ToSQL()
        if err!=nil {
        return nil, fmt.Errorf("build query: %w", err)
        }

        rows, err := {{$receiverName}}.connection.QueryxContext(ctx, query, parameters...)
        if err != nil {
        return nil, fmt.Errorf("execute query: %w", u.mapError(err))
        }
        for rows.Next() {
        var DB{{$model.Name}} {{$dbTypeRef}}
        err := rows.StructScan(&DB{{$model.Name}})
        if err != nil {
        return nil, fmt.Errorf("scan row: %w", err)
        }

        result, err := {{ template "storage.func.db_model_to_model" $model.Name}}(&DB{{$model.Name}})
        if err != nil {
        return nil, fmt.Errorf("convert db {{$model.Name | lCamelCase}} to {{$model.Name | lCamelCase}}: %w", err)
        }
        return result, nil
        }

        return nil, fmt.Errorf("no rows returned during update")
        }

        {{/*            FIRST        */}}
        {{$methodStart}} First(ctx context.Context, filter *{{$svcModelFilterRef}}, options ...{{$selectOptions}}) (*{{$svcModelRef}}, error) {
        queryBuilder := {{$receiverName}}.dialect.Select("*").From({{$tableName}}).Prepared(true)


        // Build the query based on the filter and options
        query, parameters, err := queryBuilder.Where().Limit(1).ToSQL()
        if err != nil {
        return nil, fmt.Errorf("build query: %w", err)
        }

        row := {{$receiverName}}.connection.QueryRowxContext(ctx, query, parameters...)
        var DB{{$model.Name}} {{$dbTypeRef}}
        if err := row.StructScan(&DB{{$model.Name}}); err != nil {
        return nil, fmt.Errorf("scan row: %w", err)
        }

        result, err := {{ template "storage.func.db_model_to_model" $model.Name }}(&DB{{$model.Name}})
        if err != nil {
        return nil, fmt.Errorf("convert db {{$model.Name | lCamelCase}} to {{$model.Name | lCamelCase}}: %w", err)
        }
        return result, nil
        }

        {{/*            FIRST OR CREATE        */}}
        {{$methodStart}} FirstOrCreate(ctx context.Context, filter *{{$svcModelFilterRef}}, model *{{$svcModelRef}}, options ...{{$selectOptions}}) (*{{$svcModelRef}}, error) {
        // Try to find the first matching record
        result, err := u.First(ctx, filter, options...)
        if err == nil {
        return result, nil
        }

        // If not found, create the record
        return u.Save(ctx, model)
        }

        {{/*            FIND        */}}
        {{$methodStart}} Find(ctx context.Context, filter *{{$svcModelFilterRef}}, options ...{{$selectOptions}}) ([]*{{$svcModelRef}}, error) {
        queryBuilder := {{$receiverName}}.dialect.Select("*").From({{$tableName}}).Prepared(true)

        query, parameters, err := queryBuilder.Where().ToSQL()
        if err != nil {
        return nil, fmt.Errorf("build query: %w", err)
        }

        rows, err := {{$receiverName}}.connection.QueryxContext(ctx, query, parameters...)
        if err != nil {
        return nil, fmt.Errorf("execute query: %w", u.mapError(err))
        }

        var results []*{{$svcModelRef}}
        for rows.Next() {
        var DB{{$model.Name}} {{$dbTypeRef}}
        if err := rows.StructScan(&DB{{$model.Name}}); err != nil {
        return nil, fmt.Errorf("scan row: %w", err)
        }

        result, err := {{ template "storage.func.db_model_to_model" $model.Name }}(&DB{{$model.Name}})
        if err != nil {
        return nil, fmt.Errorf("convert db {{$model.Name | lCamelCase}} to {{$model.Name | lCamelCase}}: %w", err)
        }
        results = append(results, result)
        }

        return results, nil
        }

        {{/*            DELETE        */}}
        {{$methodStart}} Delete(ctx context.Context, filter *{{$svcModelFilterRef}}) error {
        queryBuilder := {{$receiverName}}.dialect.Delete({{$tableName}}).Prepared(true)

        query, parameters, err := queryBuilder.Where().ToSQL()
        if err != nil {
        return fmt.Errorf("build query: %w", err)
        }

        _, err = {{$receiverName}}.connection.ExecContext(ctx, query, parameters...)
        if err != nil {
        return fmt.Errorf("execute query: %w", u.mapError(err))
        }

        return nil
        }

        {{$methodStart}} mapError(err error) error {
        if err == nil {
        return nil
        }

        {{ userBlock (printf "%s error mapper" $model.Name) }}

        return err
        }
        {{$methodStart}} mapField(field {{$servicePkg.Ref (printf "%sField" $model.Name)}}) (string, error) {
        switch field {
        {{- range $field := $model.Fields }}
            {{- if not $field.DoNotPersists }}
              case {{$servicePkg.Ref (include "service.const.model_field" (dict "model" $model "field" $field))}}:
              return {{template "storage.const.model_table_field" (dict "model" $model "field" $field)}}, nil
            {{- end }}
        {{- end }}
        }
        return "", fmt.Errorf("unknown field: %s", field)
        }

        {{$methodStart}} filterExpr(filter *{{$svcModelFilterRef}}) ({{ $goquPkg.Ref "Expression"}}, error) {
        if filter == nil {
        return nil, nil
        }

        return {{$dbutilPkg.Ref "BuildFilterExpression"}}(
        {{- range $field := $model.Fields }}
            {{- $fieldGoType := goType $field.Type}}
            {{- if $field.Filterable }}
                {{- if $fieldGoType.IsSlice }}
                    {{- if or (isModuleEnum (goType $field.Type.ElemType)) (isCommonEnum (goType $field.Type.ElemType)) }}
                        {{$dbutilPkg.Ref "MappedColumnArrayFilter"}}[{{(goType $field.Type.ElemType).Ref}}, string]{
                        Column: "{{$field.DBName}}",
                        Filter: filter.{{$field.Name}},
                        Mapper: {{ template "storage.storage.func.db_enum_to_enum" $field.Type.ElemType.Type}},
                        },
                    {{- else if (goType $field.Type.ElemType).IsPtr }}
                      func() {{$dbutilPkg.Ref "ColumnFilter"}}[any] {
                      panic("Slice of pointers is not implemented in generator, so we can't handle '{{$field.DBName}}' filter")
                      }(),
                    {{- else }}
                        {{$dbutilPkg.Ref "ColumnArrayFilter"}}[{{(goType $field.Type.ElemType).Ref}}]{
                        Column: "{{$field.DBName}}",
                        Filter: filter.{{$field.Name}},
                        },
                    {{- end }}
                {{- else -}}
                    {{- if or (isModuleEnum (goType $field.Type)) (isCommonEnum (goType $field.Type)) }}
                        {{$dbutilPkg.Ref "MappedColumnFilter"}}[{{$fieldGoType.Ref}}, string]{
                        Column: "{{$field.DBName}}",
                        Filter: filter.{{$field.Name}},
                        Mapper: {{ template "storage.storage.func.db_enum_to_enum" $field.Type.Type}},
                        },
                    {{- else}}
                        {{$dbutilPkg.Ref "ColumnFilter"}}[{{$fieldGoType.Ref}}]{
                        Column: "{{$field.DBName}}",
                        Filter: filter.{{$field.Name}},
                        },
                    {{- end }}
                {{- end }}
            {{- end }}
        {{- end }}
        {{$dbutilPkg.Ref "ExpressionBuilderFunc"}}(func() ({{$goquPkg.Ref "Expression"}}, error) {
        if filter.Or == nil {
        return nil, nil
        }
        exprs := make([]{{$goquPkg.Ref "Expression"}}, 0, len(filter.Or))
        for _, orFilter := range filter.Or {
        expr, err := {{$receiverName}}.filterExpr(orFilter)
        if err != nil {
        return nil, err
        }
        exprs = append(exprs, expr)
        }
        return {{$goquPkg.Ref "Or"}}(exprs...), nil
        }),
        {{$dbutilPkg.Ref "ExpressionBuilderFunc"}}(func() ({{$goquPkg.Ref "Expression"}}, error) {
        if filter.And == nil {
        return nil, nil
        }
        exprs := make([]{{$goquPkg.Ref "Expression"}}, 0, len(filter.And))
        for _, andFilter := range filter.And {
        expr, err := {{$receiverName}}.filterExpr(andFilter)
        if err != nil {
        return nil, err
        }
        exprs = append(exprs, expr)
        }
        return {{$goquPkg.Ref "And"}}(exprs...), nil
        }),
        )
        }
    {{- end }}
{{- end }}





