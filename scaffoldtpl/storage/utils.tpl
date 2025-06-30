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
{{- $errorsPkg := import "github.com/pkg/errors" -}}

type mapValue[C comparable, B any] map[string]B

func (m mapValue[C, B]) Value() ({{$driverPkg.Ref "Value"}}, error) {
return {{$jsonPkg.Ref "Marshal"}}(m)
}

func (m *mapValue[C, B]) Scan(src interface{}) error {
return {{$jsonPkg.Ref "Unmarshal"}}(src.([]byte), m)
}

type ipValue net.IP

func (i *ipValue) Value() (driver.Value, error) {
return (*net.IP)(i).String(), nil
}

func (i *ipValue) Scan(src any) error {
switch src := src.(type) {
case string:
*i = ipValue(net.ParseIP(src))
case []byte:
*i = ipValue(net.ParseIP(string(src)))
default:
return {{$errorsPkg.Ref "Errorf"}}("can't parse ipValue from: %T", src)
}
return nil
}

type stringSliceValue []string
func (s stringSliceValue) Value() ({{$driverPkg.Ref "Value"}}, error) {
result := make({{$pq.Ref "StringArray"}}, 0, len(s))
for _, item := range s {
  result = append(result, item)
}

return result.Value()
}

func (s *stringSliceValue) Scan(src interface{}) error {
result := make({{$pq.Ref "StringArray"}}, 0)
err := result.Scan(src)
if err != nil {
return err
}

*s = stringSliceValue(result)

return nil
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

type uuidSlice []{{$uuidPkg.Ref "UUID"}}

func (s uuidSlice) Value() ({{$driverPkg.Ref "Value"}}, error) {
	result := make({{$pq.Ref "StringArray"}}, len(s))
	for i, u := range s {
		result[i] = u.String()
	}
	return result.Value()
}

func (s *uuidSlice) Scan(src interface{}) error {
	var result {{$pq.Ref "StringArray"}}
	if err := result.Scan(src); err != nil {
		return err
	}

	*s = make(uuidSlice, len(result))
	for i, str := range result {
		parsedUUID, err := {{$uuidPkg.Ref "Parse"}}(str)
		if err != nil {
			return err
		}
		(*s)[i] = parsedUUID
	}

	return nil
}

func fromPtr[T any](ptr *T) T {
  return *ptr
}

func toPtr[T any](val T) *T {
  return &val
}
