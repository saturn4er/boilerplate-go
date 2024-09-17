package config

import (
	"fmt"
)

type Module struct {
	Types Types `yaml:"types"`
}

func (m *Module) Merge(module *Module) error {
	return m.Types.Merge(&module.Types)
}

func (m *Module) Init(config *Config, moduleName string) error {
	if err := m.Types.Init(config, moduleName); err != nil {
		return fmt.Errorf("init types: %w", err)
	}

	return nil
}
