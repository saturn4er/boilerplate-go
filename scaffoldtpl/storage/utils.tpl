{
"file_path": "{{.Module}}/{{.Module}}storage/gen.utils.go",
"package_name": "{{.Module}}storage",
"package_path": "{{.Config.RootPackageName}}/{{.Module}}/{{.Module}}storage",
"condition": "len(Config.Modules[Module].Value.Types.Models) > 0"
}
<><><>
{{- $driverPkg := import "database/sql/driver"}}
{{- $jsonPkg := import "encoding/json" }}
{{- $pq := import "github.com/lib/pq" }}

type mapValue[C comparable, B any] map[string]B

func (m mapValue[C, B]) Value() ({{$driverPkg.Ref "Value"}}, error) {
return {{$jsonPkg.Ref "Marshal"}}(m)
}

func (m *mapValue[C, B]) Scan(src interface{}) error {
return {{$jsonPkg.Ref "Unmarshal"}}(src.([]byte), m)
}

type sliceValue[B any] []B

func (s sliceValue[B]) Value() ({{$driverPkg.Ref "Value"}}, error) {
result := make({{$pq.Ref "StringArray"}}, 0, len(s))
for _, item := range s {
tmp, err := {{$jsonPkg.Ref "Marshal"}}(item)
if err != nil {
return nil, err
}
result = append(result, string(tmp))
}

return result.Value()
}

func (s *sliceValue[B]) Scan(src interface{}) error {
result := make({{$pq.Ref "StringArray"}}, 0)
err := result.Scan(src)
if err != nil {
return err
}

*s = make([]B, 0, len(result))

for _, item := range result {
var tmp B
err = {{$jsonPkg.Ref "Unmarshal"}}([]byte(item), &tmp)
if err != nil {
return err
}
*s = append(*s, tmp)
}

return nil
}

func fromPtr[T any](ptr *T) T {
  return *ptr
}

func toPtr[T any](val T) *T {
  return &val
}
