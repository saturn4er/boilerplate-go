package scaffold

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/gertd/go-pluralize"
	"github.com/stoewer/go-strcase"
	"golang.org/x/tools/imports"

	"github.com/saturn4er/boilerplate-go/scaffold/config"
)

const uuidPackage = "github.com/google/uuid"
const decimalPackage = "github.com/shopspring/decimal"

type fileGeneratorTemplates []fileGeneratorTemplate

type fileGeneratorTemplate struct {
	TemplateName    string
	Template        string
	HelperTemplates []string
}

type fileGenerator struct {
	templateName   string
	goimportsLocal string
	module         string
	packageName    string
	packagePath    string
	imports        []codeGeneratorImport
	template       *template.Template
	config         *config.Config
	plugins        []Plugin
	userBlocks     map[string]string
}

func (g *fileGenerator) Generate(data any) (string, error) {
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
func (g *fileGenerator) goTypeByConfigType(typ config.Type) (result *GoType) {
	result = &GoType{
		codeGenerator: g,
	}

	if typ.IsPrimitive() {
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
		result.setMetadata("with_timezone", true)
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

// GoType resolves *GoType from any config type
func (g *fileGenerator) goType(val any) *GoType {
	switch val := val.(type) {
	case config.Type:
		return g.goTypeByConfigType(val)
	case *config.Type:
		return g.goTypeByConfigType(*val)
	case config.Enum:
		moduleConfig := g.config.Modules[g.module].Value
		if _, ok := moduleConfig.Types.GetTypeByName(val.Name); ok {
			return &GoType{
				codeGenerator: g,
				Package:       g.servicePackagePath(),
				Type:          val.Name,
			}
		}

		enumType, ok := g.config.Types.GetEnumByName(val.Name)
		if ok {
			return &GoType{
				codeGenerator: g,
				Package:       enumType.Package,
				Type:          enumType.Name,
			}
		}

		return &GoType{
			codeGenerator: g,
			Package:       val.Package,
			Type:          val.Name,
		}
	case *config.Enum:
		return g.goType(*val)
	case *config.OneOfValue:
		return g.goType(*val)
	case config.OneOfValue:
		return &GoType{
			codeGenerator: g,
			Package:       g.servicePackagePath(),
			Type:          val.ModelName,
		}
	case *config.Model:
		return g.goType(*val)
	case config.Model:
		return &GoType{
			codeGenerator: g,
			Package:       g.servicePackagePath(),
			Type:          val.Name,
		}
	}

	panic(fmt.Sprintf("unknown type %T", val))
}

func (g *fileGenerator) servicePackagePath() string {
	return fmt.Sprintf("%s/%s/%s", g.config.RootPackageName, g.module, g.module)
}

func (g *fileGenerator) packageImport(name string, alias ...string) *codeGeneratorImport {
	result := &codeGeneratorImport{
		generator:  g,
		ImportPath: name,
	}
	if len(alias) > 0 {
		result.Alias = alias[0]
	}

	return result
}

func (g *fileGenerator) isEnum(typ *GoType) bool {
	return g.getModuleEnum(typ.Type) != nil || g.getCommonEnum(typ.Type) != nil
}
func (g *fileGenerator) getEnum(typ *GoType) *config.Enum {
	if enum := g.getModuleEnum(typ.Type); enum != nil {
		return enum
	}

	return g.getCommonEnum(typ.Type)
}

func (g *fileGenerator) isModuleEnum(typ *GoType) bool {
	return g.getModuleEnum(typ.Type) != nil
}

func (g *fileGenerator) getModuleEnum(name string) *config.Enum {
	moduleTypes := g.config.Modules[g.module].Value.Types

	result, ok := moduleTypes.GetEnumByName(name)
	if !ok {
		return nil
	}

	return &result
}

func (g *fileGenerator) getModuleOneOf(name string) *config.OneOf {
	moduleTypes := g.config.Modules[g.module].Value.Types

	result, ok := moduleTypes.GetOneOfByName(name)
	if !ok {
		return nil
	}

	return &result
}

func (g *fileGenerator) isModuleOneOf(typ *GoType) bool {
	return g.getModuleOneOf(typ.Type) != nil
}

func (g *fileGenerator) getModuleModel(name string) *config.Model {
	moduleTypes := g.config.Modules[g.module].Value.Types

	result, ok := moduleTypes.GetModelByName(name)
	if !ok {
		return nil
	}

	return &result
}

func (g *fileGenerator) isModuleModel(typ *GoType) bool {
	return g.getModuleModel(typ.Type) != nil
}

func (g *fileGenerator) isCommonEnum(typ *GoType) bool {
	return g.getCommonEnum(typ.Type) != nil
}

func (g *fileGenerator) getCommonEnum(name string) *config.Enum {
	result, ok := g.config.Types.GetEnumByName(name)
	if !ok {
		return nil
	}

	return &result
}

func (g *fileGenerator) getCommonOneOf(name string) *config.OneOf {
	result, ok := g.config.Types.GetOneOfByName(name)
	if !ok {
		return nil
	}

	return &result
}

func (g *fileGenerator) isCommonOneOf(typ *GoType) bool {
	return g.getCommonOneOf(typ.Type) != nil
}

func (g *fileGenerator) getCommonModel(name string) *config.Model {
	result, ok := g.config.Types.GetModelByName(name)
	if !ok {
		return nil
	}

	return &result
}

func (g *fileGenerator) isCommonModel(typ *GoType) bool {
	return g.getCommonModel(typ.Type) != nil
}

func (g *fileGenerator) templateFunctions() template.FuncMap {
	funcMap := sprig.TxtFuncMap()
	customFuncs := template.FuncMap{
		"goType":         g.goType,
		"isModuleOneOf":  g.isModuleOneOf,
		"getModuleOneOf": g.getModuleOneOf,
		"isEnum":         g.isEnum,
		"getEnum":        g.getEnum,
		"isModuleEnum":   g.isModuleEnum,
		"getModuleModel": g.getModuleModel,
		"isModuleModel":  g.isModuleModel,
		"isCommonEnum":   g.isCommonEnum,
		"getCommonEnum":  g.getCommonEnum,
		"isCommonOneOf":  g.isCommonOneOf,
		"getCommonOneOf": g.getCommonOneOf,
		"isCommonModel":  g.isCommonModel,
		"getCommonModel": g.getCommonModel,
		"import":         g.packageImport,
		"lCamelCase":     strcase.LowerCamelCase,
		"camelCase":      strcase.UpperCamelCase,
		"snakeCase":      strcase.SnakeCase,
		"userBlock": func(name string) string {
			result := "// USER CODE: '" + name + "'\n"
			result += g.userBlocks[name]
			result += "\n// END USER CODE: '" + name + "'"
			return result
		},
		"receiverName": func(typeName string) (string, error) {
			if len(typeName) < 1 {
				return "", errors.New("type name is too short")
			}
			return strings.ToLower(typeName)[0:1], nil
		},
		"plural": pluralize.NewClient().Plural,
		"varNamesGenerator": func() *varNamesGenerator {
			return &varNamesGenerator{
				usedNames: make(map[string]bool),
			}
		},
		"include": func(name string, context any) (string, error) {
			result := bytes.NewBuffer(nil)
			if err := g.template.ExecuteTemplate(result, name, context); err != nil {
				return "", err
			}
			return result.String(), nil
		},
	}
	for k, v := range customFuncs {
		funcMap[k] = v
	}

	return funcMap
}

func newCodeGenerator(
	templateName string,
	importsLocal string,
	module, packageName, packagePath string,
	config *config.Config,
	templates []string,
	plugins []Plugin,
	userBlocks map[string]string,
) (*fileGenerator, error) {
	generator := &fileGenerator{
		templateName:   templateName,
		goimportsLocal: importsLocal,
		module:         module,
		packageName:    packageName,
		packagePath:    packagePath,
		config:         config,
		plugins:        plugins,
		userBlocks:     userBlocks,
	}

	textTemplate := template.New(templateName).Funcs(generator.templateFunctions())

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
