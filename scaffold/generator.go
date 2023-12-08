package scaffold

import (
	"fmt"
	"log"
	"os"
	"path"
	"plugin"
	"strings"
	"sync"

	"github.com/go-pnp/go-pnp/pkg/optionutil"
	"github.com/samber/lo"
	"go.uber.org/multierr"

	"github.com/saturn4er/boilerplate-go/scaffold/config"
	"github.com/saturn4er/boilerplate-go/scaffoldtpl"
)

type generatorOptions struct {
	OutputDir string
	Plugins   []plugin.Plugin
}

type generator struct {
	config          *config.Config
	options         *generatorOptions
	templates       []GeneratorTemplate
	helperTemplates []string
	envs            map[string]string
	plugins         []Plugin
}

func (g *generator) generate() error {
	var (
		generateErr error
		errMu       sync.RWMutex
		wg          sync.WaitGroup
	)

	for moduleName, module := range g.config.Modules {
		wg.Add(1)
		go func(moduleName string, module config.Importable[*config.Module]) {
			defer wg.Done()

			if err := g.generateModule(moduleName, module); err != nil {
				errMu.Lock()
				generateErr = multierr.Append(generateErr, fmt.Errorf("generate module %s: %w", moduleName, err))
				errMu.Unlock()
			}
		}(moduleName, module)
	}
	wg.Wait()

	return generateErr
}

func (g *generator) generateModule(moduleName string, module config.Importable[*config.Module]) error {
	if g.config.Module != "" && moduleName != g.config.Module {
		return nil
	}

	log.Println("Generating module:", moduleName)
	log.Printf("Module has %d enums", len(module.Value.Types.Enums))
	log.Printf("Module has %d models", len(module.Value.Types.Models))
	templates := append([]GeneratorTemplate{}, g.templates...)
	for _, p := range g.plugins {
		templates = append(templates, p.Templates()...)
	}
	tplContext := ModuleContext{
		Module: moduleName,
		Config: g.config,
		Env:    g.envs,
	}
	fileInfos := make([]*TemplateFileInfo, 0, len(templates))
	fileInfosByFilePath := make(map[string][]*TemplateFileInfo)
	for _, tpl := range templates {
		fileInfo, err := tpl.ResolveTemplateFileInfo(tplContext)
		if err != nil {
			return fmt.Errorf("execute template fields: %w", err)
		}
		fileInfos = append(fileInfos, fileInfo)
		fileInfosByFilePath[fileInfo.FilePath] = append(fileInfosByFilePath[fileInfo.FilePath], fileInfo)
	}

	for filePath, fileInfos := range fileInfosByFilePath {
		var needFile = lo.SomeBy(fileInfos, func(fileInfo *TemplateFileInfo) bool {
			return fileInfo.Condition
		})
		if !needFile {
			if err := os.RemoveAll(filePath); err != nil {
				return fmt.Errorf("remove file %s: %w", filePath, err)
			}
		}
	}

	for _, tpl := range g.templates {
		tplInfo, err := tpl.ResolveTemplateFileInfo(tplContext)
		if err != nil {
			return fmt.Errorf("resolve template file info: %w", err)
		}

		if !tplInfo.Condition {
			if err := os.RemoveAll(tplInfo.FilePath); err != nil {
				log.Printf("remove file: %v\n", err)
			}
		} else {
			if err := g.generateFile(
				moduleName,
				tplInfo.PackageName,
				tplInfo.PackagePath,
				tplInfo.FilePath,
				tpl.TemplatePath,
				tpl.FileTemplate,
				g.helperTemplates,
				tplContext,
			); err != nil {
				return fmt.Errorf("generate file: %w", err)
			}
		}
	}

	return nil
}

func (g *generator) getFileUserCodeBlocks(filePath string) (map[string]string, error) {
	// search for comment like // USER CODE: '{name}' and closing commit like // END USER CODE: '{name}'
	// and return map with name as key and code as value
	var result = make(map[string]string)
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return result, nil
		}
		return nil, fmt.Errorf("read file: %w", err)
	}
	file := string(fileBytes)
	for {
		commentIndex := strings.Index(file, "// USER CODE: '")
		if commentIndex == -1 {
			break
		}
		file = file[commentIndex+len("// USER CODE: '"):]
		nameEndsIndex := strings.Index(file, "'")
		name := file[:nameEndsIndex]
		file = file[nameEndsIndex+1:]
		// find new line after user block start and trim file
		newLineIndex := strings.Index(file, "\n")
		file = file[newLineIndex+1:]

		endsIndex := strings.Index(file, "\n// END USER CODE: '"+name+"'")
		if endsIndex == -1 {
			continue
			return nil, fmt.Errorf("no closing comment for user code block '%s'", name)
		}
		content := file[:endsIndex]
		result[name] = content
		file = file[endsIndex:]
	}

	return result, nil
}

func (g *generator) generateFile(
	module, packageName, packagePath, filePath, templateName, template string,
	helperTemplates []string,
	data any,
) (rerr error) {
	log.Printf("Generating file: %v\n", filePath)

	userBlocks, err := g.getFileUserCodeBlocks(filePath)
	if err != nil {
		return fmt.Errorf("get user code blocks from file %s: %w", filePath, err)
	}

	generator, err := newCodeGenerator(
		templateName,
		g.config.GoImportsLocal,
		module,
		packageName,
		packagePath,
		g.config,
		append([]string{template}, helperTemplates...),
		nil,
		userBlocks,
	)
	if err != nil {
		return fmt.Errorf("new code generator: %w", err)
	}

	if err := os.MkdirAll(path.Dir(filePath), 0o755); err != nil {
		return fmt.Errorf("mkdirall: %w", err)
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("open enums file: %w", err)
	}

	defer func() { rerr = multierr.Append(rerr, file.Close()) }()

	fileContent, err := generator.Generate(data)
	if err != nil {
		if len(fileContent) > 0 {
			if _, writeErr := file.WriteString(fileContent); writeErr != nil {
				return multierr.Append(
					err,
					fmt.Errorf("write file content %s: %w", filePath, writeErr),
				)
			}
		}

		return fmt.Errorf("generate file %s: %w", filePath, err)
	}

	if _, err := file.WriteString(fileContent); err != nil {
		return fmt.Errorf("write file content: %w", err)
	}

	return nil
}

func Generate(config *config.Config, options ...optionutil.Option[generatorOptions]) error {
	modulesGenerator := generator{
		config: config,
		envs:   map[string]string{},
	}

	for _, i2 := range os.Environ() {
		split := strings.SplitN(i2, "=", 2)
		modulesGenerator.envs[split[0]] = split[1]
	}

	modulesGenerator.options = optionutil.ApplyOptions(&generatorOptions{
		OutputDir: "./",
	}, options...)

	templates, err := LoadTemplatesFromDir(scaffoldtpl.FS, ".")
	if err != nil {
		return fmt.Errorf("load templates: %w", err)
	}

	modulesGenerator.templates = templates
	for _, template := range templates {
		modulesGenerator.helperTemplates = append(modulesGenerator.helperTemplates, template.HelperTemplates...)
	}

	return modulesGenerator.generate()
}
