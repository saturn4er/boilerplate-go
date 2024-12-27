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

	"github.com/saturn4er/boilerplate-go/scaffold/config"
)

const uuidPackage = "github.com/google/uuid"
const decimalPackage = "github.com/shopspring/decimal"

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

func (g goType) IsSimple() bool {
	switch g.Type {
	case "string", "int", "int16", "int32", "int64", "uint", "uint16", "uint32", "uint64", "float32", "float64", "bool":
		return true
	}

	return false
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
				Type:          "stringSliceValue",
				ElemType:      elemDBAlternative,
				Package:       g.codeGenerator.storagePackagePath(),
			}
		}

		return &goType{
			codeGenerator: g.codeGenerator,
			Type:          "sliceValue",
			TypeParameters: []goType{
				*elemDBAlternative,
			},
			ElemType: elemDBAlternative,
			Package:  g.codeGenerator.storagePackagePath(),
		}
	case g.IsMap:
		return &goType{
			codeGenerator: g.codeGenerator,
			Type:          "mapValue",
			TypeParameters: []goType{
				*g.KeyType.DBAlternative(),
				*g.ElemType.DBAlternative(),
			},
			KeyType:  g.KeyType.DBAlternative(),
			ElemType: g.ElemType.DBAlternative(),
			Package:  g.codeGenerator.storagePackagePath(),
		}

	case g.IsPtr:
		if g.codeGenerator.isModuleOneOf(g.ElemType) {
			return g.ElemType.DBAlternative() // one of db alternative is already a pointer
		}
		if strings.ToLower(g.ElemType.Type) == "any" {
			return g.ElemType.DBAlternative() // any is already a pointer
		}
		return &goType{
			codeGenerator: g.codeGenerator,
			IsPtr:         true,
			ElemType:      g.ElemType.DBAlternative(),
		}
	case g.codeGenerator.isCommonModel(&g),
		g.codeGenerator.isModuleModel(&g):
		return g.InPackage(g.codeGenerator.storagePackagePath()).WithName("json" + g.Type)
	case g.codeGenerator.isModuleOneOf(&g):
		return g.InPackage(g.codeGenerator.storagePackagePath()).WithName("json" + g.Type).Ptr()
	case g.Type == "any":
		return &goType{
			codeGenerator: g.codeGenerator,
			IsPtr:         true,
			ElemType: &goType{
				codeGenerator: g.codeGenerator,
				Type:          "string",
			},
		}
	default:
		if g.codeGenerator.isCommonEnum(&g) ||
			g.codeGenerator.isModuleEnum(&g) {
			return &goType{
				codeGenerator: g.codeGenerator,
				Type:          "string",
			}
		}

		return &g
	}
}

func (g *goType) copy() *goType {
	if g == nil {
		return nil
	}
	return &goType{
		codeGenerator: g.codeGenerator,
		Package:       g.Package,
		Type:          g.Type,
		ElemType:      g.ElemType.copy(),
		KeyType:       g.KeyType.copy(),
		WithTimezone:  g.WithTimezone,
		IsPtr:         g.IsPtr,
		IsOneOf:       g.IsOneOf,
		IsSlice:       g.IsSlice,
		IsMap:         g.IsMap,
		TypeParameters: lo.Map(g.TypeParameters, func(t goType, index int) goType {
			return lo.FromPtr(t.copy())
		}),
	}
}

func (g goType) InPackage(pkg string) *goType {
	result := g.copy()
	result.Package = pkg

	return result
}

func (g goType) InLocalPackage() *goType {
	return g.InPackage(g.codeGenerator.packagePath)
}

func (g goType) WithName(name string) *goType {
	result := g.copy()
	result.Type = name

	return result
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
	config         *config.Config
	userCodeBlocks map[string]string
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
func (g *codeGenerator) goTypeByConfigType(typ config.Type) (result *goType) {
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
	for _, typeParameter := range typ.TypeParameters {
		resolvedTypeParameter := g.goType(typeParameter)
		if resolvedTypeParameter == nil {
			panic("resolvedTypeParameter is nil")
		}
		result.TypeParameters = append(result.TypeParameters, *resolvedTypeParameter)
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
			case config.Enum:
				result.Package = typ.Package
				result.Type = typ.Name
			case config.Model:
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
	case config.Type:
		return g.goTypeByConfigType(val)
	case *config.Type:
		return g.goTypeByConfigType(*val)
	case config.Enum:
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
	case *config.Enum:
		return g.goType(*val)
	case *config.OneOf:
		return g.goType(*val)
	case config.OneOf:
		return &goType{
			codeGenerator: g,
			Package:       g.servicePackagePath(),
			Type:          val.Name,
		}
	case *config.OneOfValue:
		return g.goType(*val)
	case config.OneOfValue:
		return &goType{
			codeGenerator: g,
			Package:       g.servicePackagePath(),
			Type:          val.ModelName,
		}
	case *config.Model:
		return g.goType(*val)
	case config.Model:
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

func (g *codeGenerator) storagePackagePath() string {
	return fmt.Sprintf("%s/%s/%sstorage", g.config.RootPackageName, g.module, g.module)
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
func (g *codeGenerator) getEnum(typ *goType) *config.Enum {
	if enum := g.getModuleEnum(typ.Type); enum != nil {
		return enum
	}

	return g.getCommonEnum(typ.Type)
}

func (g *codeGenerator) getModel(typ *goType) *config.Model {
	if model := g.getModuleModel(typ.Type); model != nil {
		return model
	}

	return g.getCommonModel(typ.Type)
}

func (g *codeGenerator) isModuleEnum(typ *goType) bool {
	return g.getModuleEnum(typ.Type) != nil
}

func (g *codeGenerator) getModuleEnum(name string) *config.Enum {
	moduleTypes := g.config.Modules[g.module].Value.Types

	result, ok := moduleTypes.GetEnumByName(name)
	if !ok {
		return nil
	}

	return &result
}

func (g *codeGenerator) getModuleOneOf(name string) *config.OneOf {
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

func (g *codeGenerator) getModuleModel(name string) *config.Model {
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

func (g *codeGenerator) getCommonEnum(name string) *config.Enum {
	result, ok := g.config.Types.GetEnumByName(name)
	if !ok {
		return nil
	}

	return &result
}

func (g *codeGenerator) getCommonOneOf(name string) *config.OneOf {
	result, ok := g.config.Types.GetOneOfByName(name)
	if !ok {
		return nil
	}

	return &result
}

func (g *codeGenerator) isCommonOneOf(typ *goType) bool {
	return g.getCommonOneOf(typ.Type) != nil
}

func (g *codeGenerator) getCommonModel(name string) *config.Model {
	result, ok := g.config.Types.GetModelByName(name)
	if !ok {
		return nil
	}

	return &result
}

func (g *codeGenerator) isCommonModel(typ *goType) bool {
	return g.getCommonModel(typ.Type) != nil
}

func (g *codeGenerator) userCodeBlock(name string) string {
	result := fmt.Sprintf("// user code '%s'\n", name)
	if code, ok := g.userCodeBlocks[name]; ok {
		result += code
	}
	result += fmt.Sprintf("// end user code '%s'", name)

	return result
}

func (g *codeGenerator) include(tpl string, context any) (string, error) {
	buf := bytes.NewBuffer(nil)
	if err := g.template.ExecuteTemplate(buf, tpl, context); err != nil {
		return "", fmt.Errorf("execute template %s: %w", tpl, err)
	}

	return buf.String(), nil
}

func newCodeGenerator(
	templateName string,
	importsLocal string,
	module, packageName, packagePath string,
	config *config.Config,
	userCodeBlocks map[string]string,
	templates []string,
) (*codeGenerator, error) {
	generator := &codeGenerator{
		templateName:   templateName,
		goimportsLocal: importsLocal,
		module:         module,
		packageName:    packageName,
		packagePath:    packagePath,
		userCodeBlocks: userCodeBlocks,
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
		"userCodeBlock": generator.userCodeBlock,
		"addInts": func(values ...int) int {
			result := 0
			for _, value := range values {
				result += value
			}

			return result
		},
		"include": generator.include,
		"list": func(vals ...interface{}) []interface{} {
			return vals
		},
		"receiverName": func(name string) string {
			return strings.ToLower(name[:1])
		},
		"setToVarFn": func(varName string) func(val string) string {
			return func(val string) string {
				return varName + " = " + val
			}
		},
		"setToNewVarFn": func(varName string) func(val string) string {
			return func(val string) string {
				return varName + " := " + val
			}
		},
		"appendFn": func(slice string) func(val string) string {
			return func(val string) string {
				return slice + " = append(" + slice + ", " + val + ")"
			}
		},
		"takePtrFn": func() func(val string) string {
			return func(val string) string {
				return "toPtr(" + val + ")"
			}
		},
		"derefFn": func() func(val string) string {
			return func(val string) string {
				return "fromPtr(" + val + ")"
			}
		},
		"putToMapFn": func(mapName, key string) func(val string) string {
			return func(val string) string {
				return mapName + "[" + key + "] = " + val
			}
		},
		"chainFn": func(fns ...func(string) string) func(string) string {
			return func(val string) string {
				for _, fn := range fns {
					val = fn(val)
				}

				return val
			}
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
