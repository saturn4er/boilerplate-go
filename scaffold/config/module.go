package config

import (
	"fmt"
	"strings"
)

type ProducedEvent struct {
	ModuleName string
	Event      Model
}

type Module struct {
	Types          Types           `yaml:"types"`
	Produces       []string        `yaml:"produces"`
	ProducedEvents []ProducedEvent `yaml:"-"`
}

func (m *Module) Merge(module *Module) error {
	if err := m.Types.Merge(&module.Types); err != nil {
		return err
	}
	m.Produces = append(m.Produces, module.Produces...)
	return nil
}

func (m *Module) Init(config *Config, moduleName string) error {
	if err := m.Types.Init(config, moduleName); err != nil {
		return fmt.Errorf("init types: %w", err)
	}

	return nil
}

func (m *Module) resolveProduces(config *Config, moduleName string) error {
	for _, ref := range m.Produces {
		parts := strings.SplitN(ref, ".", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid produces reference %q: expected format \"module.EventName\"", ref)
		}

		srcModuleName, eventName := parts[0], parts[1]

		if srcModuleName == moduleName {
			return fmt.Errorf("produces reference %q: cannot reference own module", ref)
		}

		srcModule, ok := config.Modules[srcModuleName]
		if !ok {
			return fmt.Errorf("produces reference %q: module %q not found", ref, srcModuleName)
		}

		event, ok := srcModule.Value.Types.GetModelByName(eventName)
		if !ok {
			return fmt.Errorf("produces reference %q: event %q not found in module %q", ref, eventName, srcModuleName)
		}

		if event.StorageType != ModelStorageTypeTxOutbox {
			return fmt.Errorf("produces reference %q: model %q is not a tx_outbox event", ref, eventName)
		}

		m.ProducedEvents = append(m.ProducedEvents, ProducedEvent{
			ModuleName: srcModuleName,
			Event:      event,
		})
	}

	return nil
}
