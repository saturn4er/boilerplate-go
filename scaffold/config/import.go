package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Importable[T interface{ Merge(T) error }] struct {
	Value T
}

func (c *Importable[T]) UnmarshalYAML(unmarshal func(any) error) error {
	var configWithImports struct {
		Imports []string `yaml:"imports"`
	}

	if err := unmarshal(&configWithImports); err != nil {
		return fmt.Errorf("unmarshal import config: %w", err)
	}

	importValues := make([]T, 0, len(configWithImports.Imports))
	for _, importPath := range configWithImports.Imports {
		importContent, err := os.ReadFile(importPath)
		if err != nil {
			return fmt.Errorf("read import file '%s': %w", importValues, err)
		}
		var importValue T
		if err := yaml.Unmarshal(importContent, &importValue); err != nil {
			return fmt.Errorf("unmarshal import value '%s': %w", importPath, err)
		}

		importValues = append(importValues, importValue)
	}

	if err := unmarshal(&c.Value); err != nil {
		return err
	}
	for i, importValue := range importValues {
		if err := c.Value.Merge(importValue); err != nil {
			return fmt.Errorf("merge value with import '%s': %w", configWithImports.Imports[i], err)
		}
	}

	return nil
}
