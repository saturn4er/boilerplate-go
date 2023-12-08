package scaffold

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"text/template"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"

	"github.com/saturn4er/boilerplate-go/scaffold/config"
)

type TemplateFileInfo struct {
	FilePath    string
	PackageName string
	PackagePath string
	Condition   bool
}

type GeneratorTemplate struct {
	FilePathTemplate    *template.Template
	PackageNameTemplate *template.Template
	PackagePathTemplate *template.Template
	Condition           *vm.Program
	TemplatePath        string
	HelperTemplates     []string
	FileTemplate        string
	CustomFuncs         template.FuncMap
}

func (g *GeneratorTemplate) ResolveTemplateFileInfo(ctx ModuleContext) (*TemplateFileInfo, error) {
	result := &TemplateFileInfo{}
	filePathBuffer := bytes.NewBuffer(nil)
	if err := g.FilePathTemplate.Execute(filePathBuffer, ctx); err != nil {
		return nil, fmt.Errorf("execute file path template: %w", err)
	}
	result.FilePath = filePathBuffer.String()

	packageNameBuffer := bytes.NewBuffer(nil)
	if err := g.PackageNameTemplate.Execute(packageNameBuffer, ctx); err != nil {
		return nil, fmt.Errorf("execute package name template: %w", err)
	}
	result.PackageName = packageNameBuffer.String()

	packagePathBuffer := bytes.NewBuffer(nil)
	if err := g.PackagePathTemplate.Execute(packagePathBuffer, ctx); err != nil {
		return nil, fmt.Errorf("execute package path template: %w", err)
	}
	result.PackagePath = packagePathBuffer.String()

	if g.Condition != nil {
		result.Condition = true

		return result, nil
	}

	var programResult interface{}
	programResult, err := expr.Run(g.Condition, ctx)
	if err != nil {
		return nil, fmt.Errorf("run template condition: %w", err)
	}

	var ok bool
	condition, ok := programResult.(bool)
	if !ok {
		return nil, fmt.Errorf("template condition must return bool")
	}

	result.Condition = condition

	return result, nil
}

func LoadTemplatesFromDir(fs embed.FS, dir string) ([]GeneratorTemplate, error) {
	dirEntries, err := fs.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir: %w", err)
	}

	helperTemplates := make([]string, 0, len(dirEntries))

	for _, entry := range dirEntries {
		if entry.IsDir() || entry.Name()[0] != '.' {
			continue
		}

		entryContent, err := fs.ReadFile(path.Join(dir, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("read file: %w", err)
		}

		helperTemplates = append(helperTemplates, string(entryContent))
	}

	result := make([]GeneratorTemplate, 0, len(dirEntries))

	for _, entry := range dirEntries {
		entryPath := path.Join(dir, entry.Name())
		if entry.IsDir() {
			dirTemplates, err := LoadTemplatesFromDir(fs, entryPath)
			if err != nil {
				return nil, fmt.Errorf("load templates from dir '%v': %w", entry.Name(), err)
			}

			result = append(result, dirTemplates...)

			continue
		}

		if entry.Name()[0] == '.' {
			continue
		}

		entryContent, err := fs.ReadFile(entryPath)
		if err != nil {
			return nil, fmt.Errorf("read file '%v': %w", entryPath, err)
		}

		tplContentParts := strings.SplitN(string(entryContent), "<><><>", 2)
		if len(tplContentParts) != 2 {
			return nil, fmt.Errorf("template '%v' does not contain header separated by template with <><><>", entry.Name())
		}

		var header struct {
			FilePath    string `json:"file_path"`
			PackageName string `json:"package_name"`
			PackagePath string `json:"package_path"`
			Condition   string `json:"condition"`
		}

		if err := json.Unmarshal([]byte(tplContentParts[0]), &header); err != nil {
			return nil, fmt.Errorf("unmarshal template header: %w", err)
		}

		if header.FilePath == "" {
			return nil, fmt.Errorf("template header file_path is empty")
		}

		filePathTemplate, err := template.New(entry.Name() + ".file_path").Parse(header.FilePath)
		if err != nil {
			return nil, fmt.Errorf("parse template header file_path: %w", err)
		}

		if header.PackageName == "" {
			return nil, fmt.Errorf("template header package is empty")
		}

		packageNameTemplate, err := template.New(entry.Name() + ".package_name").Parse(header.PackageName)
		if err != nil {
			return nil, fmt.Errorf("parse template header package_name: %w", err)
		}

		if header.PackagePath == "" {
			return nil, fmt.Errorf("template header package_path is empty")
		}

		packagePathTemplate, err := template.New(entry.Name() + ".package_path").Parse(header.PackagePath)
		if err != nil {
			return nil, fmt.Errorf("parse template header package: %w", err)
		}

		var condition *vm.Program
		if header.Condition != "" {
			condition, err = expr.Compile(header.Condition, expr.Env(map[string]interface{}{
				"Module": "",
				"Config": &config.Config{},
			}))
			if err != nil {
				return nil, fmt.Errorf("compile template header condition: %w", err)
			}
		}

		result = append(result, GeneratorTemplate{
			FilePathTemplate:    filePathTemplate,
			PackageNameTemplate: packageNameTemplate,
			PackagePathTemplate: packagePathTemplate,
			Condition:           condition,
			TemplatePath:        entryPath,
			FileTemplate:        tplContentParts[1],
			HelperTemplates:     helperTemplates,
		})
	}

	return result, nil
}
