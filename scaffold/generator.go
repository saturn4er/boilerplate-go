package scaffold

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/expr-lang/expr"
	"github.com/go-pnp/go-pnp/pkg/optionutil"
	"go.uber.org/multierr"

	"github.com/saturn4er/boilerplate-go/scaffold/config"
	"github.com/saturn4er/boilerplate-go/scaffoldtpl"
)

type generatorOptions struct {
	OutputDir string
}

type generator struct {
	config    *config.Config
	options   *generatorOptions
	templates []generatorTemplate
	envs      map[string]string
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

	for _, tpl := range g.templates {
		tplVars := map[string]interface{}{
			"Module": moduleName,
			"Config": g.config,
			"Env":    g.envs,
		}
		filePath := bytes.NewBuffer(nil)

		if err := tpl.FilePathTemplate.Execute(filePath, tplVars); err != nil {
			return fmt.Errorf("execute file path template: %w", err)
		}

		packageName := bytes.NewBuffer(nil)
		if err := tpl.PackageNameTemplate.Execute(packageName, tplVars); err != nil {
			return fmt.Errorf("execute package name template: %w", err)
		}

		packagePath := bytes.NewBuffer(nil)
		if err := tpl.PackagePathTemplate.Execute(packagePath, tplVars); err != nil {
			return fmt.Errorf("execute package path template: %w", err)
		}

		conditionResult := true

		if tpl.Condition != nil {
			programResult, err := expr.Run(tpl.Condition, tplVars)
			if err != nil {
				return fmt.Errorf("run template condition: %w", err)
			}

			conditionResultBool, ok := programResult.(bool)
			if !ok {
				return fmt.Errorf("template condition must return bool")
			}

			conditionResult = conditionResultBool
		}

		if !conditionResult {
			if err := os.RemoveAll(filePath.String()); err != nil {
				log.Printf("remove file: %v\n", err)
			}
		} else {
			if err := g.generateFile(
				moduleName,
				packageName.String(),
				packagePath.String(),
				filePath.String(),
				tpl.TemplatePath,
				tpl.FileTemplate,
				tpl.HelperTemplates,
				tplVars,
			); err != nil {
				return fmt.Errorf("generate file: %w", err)
			}
		}
	}

	return nil
}

func (g *generator) generateFile(
	module, packageName, packagePath, filePath, templateName, template string,
	helperTemplates []string,
	data any,
) (rerr error) {
	log.Printf("Generating file: %v\n", filePath)

	userCodeBlocks, err := g.getFileUserCodeBlocks(filePath)
	if err != nil {
		return fmt.Errorf("get file user code blocks: %w", err)
	}

	generator, err := newCodeGenerator(
		templateName,
		g.config.GoImportsLocal,
		module,
		packageName,
		packagePath,
		g.config,
		userCodeBlocks,
		append([]string{template}, helperTemplates...),
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
					fmt.Errorf("write file content: %w", writeErr),
				)
			}
		}

		return fmt.Errorf("generate file: %w", err)
	}

	if _, err := file.WriteString(fileContent); err != nil {
		return fmt.Errorf("write file content: %w", err)
	}

	return nil
}

func (g *generator) getFileUserCodeBlocks(path string) (map[string]string, error) {
	// user code is started with `// user code '<block name>'` and ended with `// end user code '<block name>'`
	fileContent, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("read file: %w", err)
	}
	leftFileContent := string(fileContent)
	result := make(map[string]string)
	for {
		// searching for code block name start
		codeBlockIdx := strings.Index(leftFileContent, "// user code '")
		if codeBlockIdx == -1 {
			break
		}

		leftFileContent = leftFileContent[codeBlockIdx+len("// user code '"):]
		// searching for code block name
		codeBlockIdx = strings.Index(leftFileContent, "'")
		if codeBlockIdx == -1 {
			break
		}

		codeBlockName := leftFileContent[:codeBlockIdx]

		newLineIdx := strings.Index(leftFileContent, "\n")
		if newLineIdx == -1 {
			break
		}
		leftFileContent = leftFileContent[newLineIdx+1:]

		// searching for code block end
		codeBlockEndIdx := strings.Index(leftFileContent, fmt.Sprintf("// end user code '%s'", codeBlockName))
		if codeBlockEndIdx == -1 {
			break
		}

		result[codeBlockName] = leftFileContent[:codeBlockEndIdx]
		leftFileContent = leftFileContent[codeBlockEndIdx:]
	}

	return result, nil
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

	templates, err := loadTemplatesFromDir(scaffoldtpl.FS, ".")
	if err != nil {
		return fmt.Errorf("load templates: %w", err)
	}

	modulesGenerator.templates = templates

	return modulesGenerator.generate()
}
