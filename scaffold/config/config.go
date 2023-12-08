package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	RootPackageName string                         `yaml:"root_package_name"`
	GoImportsLocal  string                         `yaml:"goimports_local"`
	Types           Types                          `yaml:"types"`
	Modules         map[string]Importable[*Module] `yaml:"modules"`
	Module          string                         `yaml:"-"`
}

func (c *Config) EachModel(considerCommon bool, fn func(module *Model)) {
	for _, module := range c.Modules {
		for _, typ := range module.Value.Types.Models {
			fn(&typ)
		}
	}
	if !considerCommon {
		return
	}

	for _, module := range c.Types.Models {
		fn(&module)
	}
}

func (c *Config) EachFieldType(considerCommon bool, fn func(typ *Type)) {
	c.EachModel(considerCommon, func(module *Model) {
		for _, typ := range module.Fields {
			fn(&typ.Type)
		}
	})
}

func (c *Config) EachFieldTypeRecursive(considerCommon bool, fn func(typ *Type)) {
	var execType func(typ *Type)
	execType = func(typ *Type) {
		if typ == nil {
			return
		}
		fn(typ)
		execType(typ.KeyType)
		execType(typ.ElemType)
	}

	c.EachModel(considerCommon, func(module *Model) {
		for _, typ := range module.Fields {
			execType(&typ.Type)
		}
	})
}

func (c *Config) Init() error {
	for moduleName, module := range c.Modules {
		if c.Module != "" && moduleName != c.Module {
			continue
		}

		if err := module.Value.Init(c, moduleName); err != nil {
			return fmt.Errorf("init module %s: %w", moduleName, err)
		}
	}

	return nil
}

func Load(path string) (*Config, error) {
	configContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(configContent, &config); err != nil {
		return nil, fmt.Errorf("unmarshal config file: %w", err)
	}

	if err := config.Init(); err != nil {
		return nil, fmt.Errorf("init config: %w", err)
	}

	return &config, nil
}
