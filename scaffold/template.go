package scaffold

import (
	"embed"
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"text/template"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

type generatorTemplate struct {
	FilePathTemplate    *template.Template
	PackageNameTemplate *template.Template
	PackagePathTemplate *template.Template
	Condition           *vm.Program
	TemplatePath        string
	HelperTemplates     []string
	FileTemplate        string
}

func loadTemplatesFromDir(fs embed.FS, dir string) ([]generatorTemplate, error) {
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

	result := make([]generatorTemplate, 0, len(dirEntries))

	for _, entry := range dirEntries {
		entryPath := path.Join(dir, entry.Name())
		if entry.IsDir() {
			dirTemplates, err := loadTemplatesFromDir(fs, entryPath)
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
				"Config": &Config{},
			}))
			if err != nil {
				return nil, fmt.Errorf("compile template header condition: %w", err)
			}
		}

		result = append(result, generatorTemplate{
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
