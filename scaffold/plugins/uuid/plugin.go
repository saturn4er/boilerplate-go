package uuid

import (
	"text/template"

	"github.com/saturn4er/boilerplate-go/scaffold"
	"github.com/saturn4er/boilerplate-go/scaffold/config"
)

const (
	UUIDPackage = "github.com/google/uuid"
)

type Plugin struct {
}

var _ scaffold.Plugin = &Plugin{}

func (p Plugin) Init(cfg *config.Config) error {
	cfg.EachFieldTypeRecursive(true, func(typ *config.Type) {
		if typ.Type == "uuid" {
			typ.Package = UUIDPackage
			typ.Type = "UUID"
		}
	})
	return nil
}

func (p Plugin) Name() string {
	return "uuid"
}

func (p Plugin) TemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"isUUID": func(val any) bool {
			switch v := val.(type) {
			case *scaffold.GoType:
				return v.Type == "UUID" && v.Package == UUIDPackage && v.IsSlice == false && v.IsPtr == false && v.IsMap == false
			case scaffold.GoType:
				return v.Type == "UUID" && v.Package == UUIDPackage && v.IsSlice == false && v.IsPtr == false && v.IsMap == false
			case *config.Type:
				return v.Type == "UUID" && v.Package == UUIDPackage
			case config.Type:
				return v.Type == "UUID" && v.Package == UUIDPackage
			}
			return false
		},
	}
}

func New() Plugin {
	return Plugin{}
}
