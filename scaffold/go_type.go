package scaffold

import (
	"fmt"

	"github.com/samber/lo"
)

type GoType struct {
	codeGenerator  *fileGenerator
	Package        string  // package name (ex: "github.com/google/uuid")
	Type           string  // type name (ex: "string", "int", "Time" for time.Time)
	ElemType       *GoType // type of elements if slices, map values, or pointer underlying type
	KeyType        *GoType // type of map key if this is a map
	WithTimezone   bool
	IsPtr          bool
	IsOneOf        bool
	IsSlice        bool
	IsMap          bool
	TypeParameters []GoType
	Metadata       map[string]any
}

func (g GoType) setMetadata(key string, value any) {
	if g.Metadata == nil {
		g.Metadata = make(map[string]any)
	}
	g.Metadata[key] = value
}

func (g GoType) GoAdminForm() string {
	formImport := g.codeGenerator.packageImport("github.com/GoAdminGroup/go-admin/template/types/form")
	switch {
	case g.IsSlice:
		return formImport.Ref("Code")
	case g.IsPtr:
		return g.ElemType.GoAdminForm()
	case g.IsMap:
		return formImport.Ref("Code")
	case g.Type == "string": //nolint:goconst
		return formImport.Ref("Text")
	case lo.Contains([]string{"int", "int16", "int32", "int64", "uint", "uint16", "uint32", "uint64"}, g.Type):
		return formImport.Ref("Text")
	case lo.Contains([]string{"float32", "float64"}, g.Type):
		return formImport.Ref("Text")
	case g.Type == "bool":
		return formImport.Ref("Text")
	case g.codeGenerator.isModuleEnum(&g) || g.codeGenerator.isCommonEnum(&g):
		return formImport.Ref("SelectSingle")
	case g.Package == "time" && g.Type == "Time":
		return formImport.Ref("Datetime")
	case g.Package == uuidPackage && g.Type == "UUID":
		return formImport.Ref("Text")
	case g.codeGenerator.isModuleOneOf(&g) || g.codeGenerator.isCommonOneOf(&g) || g.codeGenerator.isModuleModel(&g) || g.codeGenerator.isCommonModel(&g):
		return formImport.Ref("Code")
	}

	return formImport.Ref("Text")
}
func (g GoType) GoAdminType() string {
	dbImport := g.codeGenerator.packageImport("github.com/GoAdminGroup/go-admin/modules/db")
	switch {
	case g.IsSlice:
		return dbImport.Ref("Text")
	case g.IsPtr:
		return g.ElemType.GoAdminType()
	case g.IsMap:
		return dbImport.Ref("JSON")
	case g.Type == "string": //nolint:goconst
		return dbImport.Ref("Text")
	case g.Type == "StringArray" && g.Package == "github.com/lib/pq":
		return "Text"
	case lo.Contains([]string{"int", "int16", "int32", "int64", "uint", "uint16", "uint32", "uint64"}, g.Type):
		return dbImport.Ref("Int")
	case lo.Contains([]string{"float32", "float64"}, g.Type):
		return dbImport.Ref("Float")
	case g.Type == "bool":
		return dbImport.Ref("Bool")
	case g.codeGenerator.isModuleEnum(&g) || g.codeGenerator.isCommonEnum(&g):
		return dbImport.Ref("Enum")
	case g.Package == "time" && g.Type == "Time":
		return dbImport.Ref("Timestamp")
	case g.Package == uuidPackage && g.Type == "UUID":
		return dbImport.Ref("UUID")
	case g.codeGenerator.isModuleOneOf(&g) || g.codeGenerator.isCommonOneOf(&g) || g.codeGenerator.isModuleModel(&g) || g.codeGenerator.isCommonModel(&g):
		return dbImport.Ref("JSON")
	}

	return dbImport.Ref("Text")
}
func (g GoType) GormType() string {
	switch {
	case g.IsSlice:
		return g.ElemType.GormType() + "[]"
	case g.IsPtr:
		return g.ElemType.GormType()
	case g.IsMap:
		return "jsonb"
	case g.Type == "string": //nolint:goconst
		return "text"
	case g.Type == "StringArray" && g.Package == "github.com/lib/pq":
		return "text[]"
	}

	return ""
}

func (g GoType) PackageImport() *codeGeneratorImport {
	return g.codeGenerator.packageImport(g.Package)
}

func (g GoType) DBAlternative() *GoType {
	switch {
	case g.IsSlice:
		elemDBAlternative := g.ElemType.DBAlternative()
		if elemDBAlternative.Type == "string" {
			return &GoType{
				codeGenerator: g.codeGenerator,
				ElemType: &GoType{
					codeGenerator: g.codeGenerator,
					Type:          "string",
				},
				Package: "github.com/lib/pq",
				Type:    "StringArray",
			}
		}
		panic(fmt.Sprintf("storing '%v' is not supported", g.Ref()))
	case g.IsMap:
		return &GoType{
			codeGenerator: g.codeGenerator,
			IsMap:         true,
			KeyType:       g.KeyType.DBAlternative(),
			ElemType:      g.ElemType.DBAlternative(),
			Metadata:      g.Metadata,
		}
	case g.IsPtr:
		return &GoType{
			codeGenerator: g.codeGenerator,
			IsPtr:         true,
			ElemType:      g.ElemType.DBAlternative(),
		}
	default:
		if g.codeGenerator.isModuleOneOf(&g) {
			return &GoType{
				codeGenerator: g.codeGenerator,
				IsPtr:         true,
				ElemType: &GoType{
					codeGenerator: g.codeGenerator,
					Type:          "string",
				},
			}
		}
		if g.codeGenerator.isCommonModel(&g) ||
			g.codeGenerator.isModuleModel(&g) ||
			g.codeGenerator.isCommonEnum(&g) ||
			g.codeGenerator.isModuleEnum(&g) ||
			g.Type == "any" {
			return &GoType{
				codeGenerator: g.codeGenerator,
				Type:          "string",
			}
		}

		return &g
	}
}

func (g GoType) InLocalPackage() *GoType {
	g.Package = g.codeGenerator.packagePath

	return &g
}

func (g GoType) WithName(name string) *GoType {
	g.Type = name

	return &g
}

func (g GoType) Ptr() *GoType {
	return &GoType{
		codeGenerator: g.codeGenerator,
		ElemType:      &g,
		IsPtr:         true,
	}
}

func (g GoType) Ref() string {
	if g.IsSlice {
		return "[]" + g.ElemType.Ref()
	}

	if g.IsMap {
		return "map[" + g.KeyType.Ref() + "]" + g.ElemType.Ref()
	}

	if g.IsPtr {
		return "*" + g.ElemType.Ref()
	}

	result := g.Type

	if g.Package != "" && g.Package != g.codeGenerator.packagePath {
		result = g.codeGenerator.packageImport(g.Package).Ref(g.Type)
	}

	if len(g.TypeParameters) > 0 {
		result += "["
		for i, param := range g.TypeParameters {
			if i > 0 {
				result += ","
			}
			result += param.Ref()
		}
		result += "]"
	}

	return result
}
