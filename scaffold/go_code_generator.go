package scaffold

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/gertd/go-pluralize"
	"github.com/samber/lo"
	"github.com/stoewer/go-strcase"
	"golang.org/x/tools/imports"
)

const uuidPackage = "github.com/google/uuid"
const decimalPackage = "github.com/shopspring/decimal"

type indexed[T comparable, V any] struct {
	Index T
	Value V
}
type goType struct {
	codeGenerator  *codeGenerator
	Package        string
	Type           string
	ElemType       *goType
	KeyType        *goType
	WithTimezone   bool
	IsPtr          bool
	IsOneOf        bool
	IsSlice        bool
	IsMap          bool
	TypeParameters []goType
}

func (g goType) GoAdminForm() string {
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
func (g goType) GoAdminType() string {
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
func (g goType) GormType() string {
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

func (g goType) PackageImport() *codeGeneratorImport {
	return g.codeGenerator.packageImport(g.Package)
}

func (g goType) DBAlternative() *goType {
	switch {
	case g.IsSlice:
		elemDBAlternative := g.ElemType.DBAlternative()
		if elemDBAlternative.Type == "string" {
			return &goType{
				codeGenerator: g.codeGenerator,
				ElemType: &goType{
					codeGenerator: g.codeGenerator,
					Type:          "string",
				},
				Package: "github.com/lib/pq",
				Type:    "StringArray",
			}
		}
		panic(fmt.Sprintf("storing '%v' is not supported", g.Ref()))
	case g.IsMap:
		return &goType{
			codeGenerator: g.codeGenerator,
			IsMap:         true,
			KeyType:       g.KeyType.DBAlternative(),
			ElemType:      g.ElemType.DBAlternative(),
		}
	case g.IsPtr:
		return &goType{
			codeGenerator: g.codeGenerator,
			IsPtr:         true,
			ElemType:      g.ElemType.DBAlternative(),
		}
	default:
		if g.codeGenerator.isModuleOneOf(&g) {
			return &goType{
				codeGenerator: g.codeGenerator,
				IsPtr:         true,
				ElemType: &goType{
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
			return &goType{
				codeGenerator: g.codeGenerator,
				Type:          "string",
			}
		}

		return &g
	}
}

func (g goType) InLocalPackage() *goType {
	g.Package = g.codeGenerator.packagePath

	return &g
}

func (g goType) WithName(name string) *goType {
	g.Type = name

	return &g
}

func (g goType) Ptr() *goType {
	return &goType{
		codeGenerator: g.codeGenerator,
		ElemType:      &g,
		IsPtr:         true,
	}
}

func (g goType) Ref() string {
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

type codeGeneratorImport struct {
	generator *codeGenerator

	ImportPath string
	Alias      string
}

func (c codeGeneratorImport) Ref(val string) string {
	for _, imprt := range c.generator.imports {
		if imprt.ImportPath == c.ImportPath {
			return imprt.Alias + "." + val
		}
	}

	pathParts := strings.Split(c.ImportPath, "/")
	if c.Alias == "" {
		c.Alias = pathParts[len(pathParts)-1]
	}

	// check if alias is already used and add a number suffix if it is
	for i := 0; ; i++ {
		aliasUsed := false
		suffix := ""

		if i > 0 {
			suffix = strconv.FormatInt(int64(i), 10)
		}

		for _, imprt := range c.generator.imports {
			if imprt.Alias == c.Alias+suffix {
				aliasUsed = true

				break
			}
		}

		if !aliasUsed {
			c.Alias += suffix

			break
		}
	}

	c.generator.imports = append(c.generator.imports, c)

	return c.Alias + "." + val
}

type codeGenerator struct {
	templateName   string
	goimportsLocal string
	module         string
	packageName    string
	packagePath    string
	imports        []codeGeneratorImport
	template       *template.Template
	config         *Config
}

func (g *codeGenerator) Generate(data any) (string, error) {
	body := bytes.NewBuffer(nil)
	if err := g.template.ExecuteTemplate(body, g.templateName, data); err != nil {
		return "", err
	}

	head := bytes.NewBuffer(nil)
	head.WriteString("package " + g.packageName + "\n")

	if len(g.imports) > 0 {
		head.WriteString("import (\n")

		for _, imprt := range g.imports {
			head.WriteString("\t" + imprt.Alias + " \"" + imprt.ImportPath + "\"\n")
		}

		head.WriteString(")\n")
	}

	result := head.String() + "\n" + body.String()
	imports.LocalPrefix = g.goimportsLocal
	// execute goimports
	formattedResult, err := imports.Process("name", []byte(result), &imports.Options{
		TabIndent:  true,
		TabWidth:   8,
		Comments:   true,
		Fragment:   true,
		FormatOnly: false,
	})
	if err != nil {
		return result, err
	}

	return string(formattedResult), nil
}

//nolint:funlen
func (g *codeGenerator) goTypeByConfigType(typ ConfigType) (result *goType) {
	result = &goType{
		codeGenerator: g,
	}
	isPrimitive := func(typ string) bool {
		switch typ {
		case "bool", "string", "byte", "rune",
			"int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64",
			"uintptr", "float32", "float64", "complex64", "complex128":
			return true
		}

		return false
	}

	if isPrimitive(typ.Type) {
		result.Type = typ.Type

		return result
	}

	if typ.ElemType != nil {
		result.ElemType = g.goType(*typ.ElemType)
	}

	if typ.KeyType != nil {
		result.KeyType = g.goType(*typ.KeyType)
	}

	switch typ.Type {
	case "ptr":
		result.IsPtr = true
		if result.ElemType == nil {
			panic("ptr must have an element type")
		}
	case "map":
		result.IsMap = true
		if result.KeyType == nil {
			panic("map must have a key type")
		}

		if result.ElemType == nil {
			panic("map must have an element type")
		}
	case "slice":
		result.IsSlice = true
		if result.ElemType == nil {
			panic("slice must have a key type")
		}
	case "decimal":
		result.Package = decimalPackage
		result.Type = "Decimal"
	case "uuid":
		result.Package = uuidPackage
		result.Type = "UUID"
	case "timestamp":
		result.Package = "time"
		result.Type = "Time"
	case "timestamptz":
		result.Package = "time"
		result.Type = "Time"
		result.WithTimezone = true
	case "duration":
		result.Package = "time"
		result.Type = "Duration"
	default:
		// check if this is module type
		moduleConfig := g.config.Modules[g.module].Value
		if _, ok := moduleConfig.Types.GetTypeByName(typ.Type); ok {
			result.Package = g.servicePackagePath()
			result.Type = typ.Type

			return result
		}

		commonType, ok := g.config.Types.GetTypeByName(typ.Type)
		if ok {
			switch typ := commonType.(type) {
			case ConfigEnum:
				result.Package = typ.Package
				result.Type = typ.Name
			case ConfigModel:
				result.Package = typ.Package
				result.Type = typ.Name
			}

			return result
		}

		result.Package = g.packagePath
		if typ.Package != "" {
			result.Package = typ.Package
		}

		result.Type = typ.Type
	}

	return result
}

func (g *codeGenerator) goType(val any) *goType {
	switch val := val.(type) {
	case ConfigType:
		return g.goTypeByConfigType(val)
	case *ConfigType:
		return g.goTypeByConfigType(*val)
	case ConfigEnum:
		moduleConfig := g.config.Modules[g.module].Value
		if _, ok := moduleConfig.Types.GetTypeByName(val.Name); ok {
			return &goType{
				codeGenerator: g,
				Package:       g.servicePackagePath(),
				Type:          val.Name,
			}
		}

		enumType, ok := g.config.Types.GetEnumByName(val.Name)
		if ok {
			return &goType{
				codeGenerator: g,
				Package:       enumType.Package,
				Type:          enumType.Name,
			}
		}

		return &goType{
			codeGenerator: g,
			Package:       val.Package,
			Type:          val.Name,
		}
	case *ConfigEnum:
		return g.goType(*val)
	case *ConfigOneOfValue:
		return g.goType(*val)
	case ConfigOneOfValue:
		return &goType{
			codeGenerator: g,
			Package:       g.servicePackagePath(),
			Type:          val.ModelName,
		}
	case *ConfigModel:
		return g.goType(*val)
	case ConfigModel:
		return &goType{
			codeGenerator: g,
			Package:       g.servicePackagePath(),
			Type:          val.Name,
		}
	}

	panic(fmt.Sprintf("unknown type %T", val))
}

func (g *codeGenerator) servicePackagePath() string {
	return fmt.Sprintf("%s/%s/%sservice", g.config.RootPackageName, g.module, g.module)
}

func (g *codeGenerator) packageImport(name string, alias ...string) *codeGeneratorImport {
	result := &codeGeneratorImport{
		generator:  g,
		ImportPath: name,
	}
	if len(alias) > 0 {
		result.Alias = alias[0]
	}

	return result
}

func (g *codeGenerator) isEnum(typ *goType) bool {
	return g.getModuleEnum(typ.Type) != nil || g.getCommonEnum(typ.Type) != nil
}
func (g *codeGenerator) getEnum(typ *goType) *ConfigEnum {
	if enum := g.getModuleEnum(typ.Type); enum != nil {
		return enum
	}

	return g.getCommonEnum(typ.Type)
}

func (g *codeGenerator) getModel(typ *goType) *ConfigModel {
	if model := g.getModuleModel(typ.Type); model != nil {
		return model
	}

	return g.getCommonModel(typ.Type)
}

func (g *codeGenerator) isModuleEnum(typ *goType) bool {
	return g.getModuleEnum(typ.Type) != nil
}

func (g *codeGenerator) getModuleEnum(name string) *ConfigEnum {
	moduleTypes := g.config.Modules[g.module].Value.Types

	result, ok := moduleTypes.GetEnumByName(name)
	if !ok {
		return nil
	}

	return &result
}

func (g *codeGenerator) getModuleOneOf(name string) *ConfigOneOf {
	moduleTypes := g.config.Modules[g.module].Value.Types

	result, ok := moduleTypes.GetOneOfByName(name)
	if !ok {
		return nil
	}

	return &result
}

func (g *codeGenerator) isModuleOneOf(typ *goType) bool {
	return g.getModuleOneOf(typ.Type) != nil
}

func (g *codeGenerator) getModuleModel(name string) *ConfigModel {
	moduleTypes := g.config.Modules[g.module].Value.Types

	result, ok := moduleTypes.GetModelByName(name)
	if !ok {
		return nil
	}

	return &result
}

func (g *codeGenerator) isModuleModel(typ *goType) bool {
	return g.getModuleModel(typ.Type) != nil
}

func (g *codeGenerator) isCommonEnum(typ *goType) bool {
	return g.getCommonEnum(typ.Type) != nil
}

func (g *codeGenerator) getCommonEnum(name string) *ConfigEnum {
	result, ok := g.config.Types.GetEnumByName(name)
	if !ok {
		return nil
	}

	return &result
}

func (g *codeGenerator) getCommonOneOf(name string) *ConfigOneOf {
	result, ok := g.config.Types.GetOneOfByName(name)
	if !ok {
		return nil
	}

	return &result
}

func (g *codeGenerator) isCommonOneOf(typ *goType) bool {
	return g.getCommonOneOf(typ.Type) != nil
}

func (g *codeGenerator) getCommonModel(name string) *ConfigModel {
	result, ok := g.config.Types.GetModelByName(name)
	if !ok {
		return nil
	}

	return &result
}

func (g *codeGenerator) isCommonModel(typ *goType) bool {
	return g.getCommonModel(typ.Type) != nil
}

func newCodeGenerator(
	templateName string,
	importsLocal string,
	module, packageName, packagePath string,
	config *Config,
	templates []string,
) (*codeGenerator, error) {
	generator := &codeGenerator{
		templateName:   templateName,
		goimportsLocal: importsLocal,
		module:         module,
		packageName:    packageName,
		packagePath:    packagePath,
		config:         config,
	}
	textTemplate := template.New(templateName).Funcs(template.FuncMap{
		"goType":         generator.goType,
		"isModuleOneOf":  generator.isModuleOneOf,
		"getModuleOneOf": generator.getModuleOneOf,
		"isEnum":         generator.isEnum,
		"getEnum":        generator.getEnum,
		"getModel":       generator.getModel,
		"isModuleEnum":   generator.isModuleEnum,
		"getModuleModel": generator.getModuleModel,
		"isModuleModel":  generator.isModuleModel,
		"isCommonEnum":   generator.isCommonEnum,
		"getCommonEnum":  generator.getCommonEnum,
		"isCommonOneOf":  generator.isCommonOneOf,
		"getCommonOneOf": generator.getCommonOneOf,
		"isCommonModel":  generator.isCommonModel,
		"getCommonModel": generator.getCommonModel,
		"import":         generator.packageImport,
		"uCamelCase":     strcase.UpperCamelCase,
		"lCamelCase":     strcase.LowerCamelCase,
		"snakeCase":      strcase.SnakeCase,
		"replace":        strings.ReplaceAll,
		"plural":         pluralize.NewClient().Plural,
		"varNamesGenerator": func() *varNamesGenerator {
			return &varNamesGenerator{
				usedNames: make(map[string]bool),
			}
		},
		"addInts": func(values ...int) int {
			result := 0
			for _, value := range values {
				result += value
			}

			return result
		},
		"list": func(vals ...interface{}) []interface{} {
			return vals
		},
		"receiverName": func(name string) string {
			return strings.ToLower(name[:1])
		},
	})

	var err error

	for _, template := range templates {
		textTemplate, err = textTemplate.Parse(template)
		if err != nil {
			return nil, fmt.Errorf("parse template: %w", err)
		}
	}

	generator.template = textTemplate

	return generator, nil
}
