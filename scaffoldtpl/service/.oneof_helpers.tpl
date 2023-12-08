{{- define "oneOfType" }}
    {{- /*gotype: github.com/saturn4er/boilerplate-go/scaffold/config.OneOf */}}
  type {{.Name}} interface{
  is{{.Name}}()
  }
  {{- range $value := .SortedValues}}
    func (*{{$value.Value.ModelName}}) is{{$.Name}}() {}
  {{- end}}
{{- end }}
