package config

import "fmt"

type Types struct {
	Enums  []Enum  `yaml:"enums"`
	Models []Model `yaml:"models"`
	OneOfs []OneOf `yaml:"one_ofs"`
}

func (t *Types) Merge(types *Types) error {
	t.Enums = append(t.Enums, types.Enums...)
	t.Models = append(t.Models, types.Models...)
	t.OneOfs = append(t.OneOfs, types.OneOfs...)

	// error in case of duplicated names
	names := make(map[string]struct{})
	for _, enum := range types.Enums {
		if _, ok := names[enum.Name]; ok {
			return fmt.Errorf("duplicate name: %s", enum.Name)
		}
		names[enum.Name] = struct{}{}
	}
	for _, model := range types.Models {
		if _, ok := names[model.Name]; ok {
			return fmt.Errorf("duplicate name: %s", model.Name)
		}
		names[model.Name] = struct{}{}
	}
	for _, oneOf := range types.OneOfs {
		if _, ok := names[oneOf.Name]; ok {
			return fmt.Errorf("duplicate name: %s", oneOf.Name)
		}
		names[oneOf.Name] = struct{}{}
	}
	return nil
}

func (t *Types) GetTypeByName(name string) (any, bool) {
	if model, ok := t.GetModelByName(name); ok {
		return model, true
	}

	if enum, ok := t.GetEnumByName(name); ok {
		return enum, true
	}

	if oneOf, ok := t.GetOneOfByName(name); ok {
		return oneOf, true
	}

	return nil, false
}

func (t *Types) GetModelByName(name string) (Model, bool) {
	for _, model := range t.Models {
		if model.Name == name {
			return model, true
		}
	}

	return Model{}, false
}

func (t *Types) GetEnumByName(name string) (Enum, bool) {
	for _, enum := range t.Enums {
		if enum.Name == name {
			return enum, true
		}
	}

	return Enum{}, false
}

func (t *Types) GetOneOfByName(name string) (OneOf, bool) {
	for _, oneOf := range t.OneOfs {
		if oneOf.Name == name {
			return oneOf, true
		}
	}

	return OneOf{}, false
}

func (t *Types) Init(config *Config, moduleName string) error {
	for i, enum := range t.Enums {
		if err := enum.Init(config, moduleName); err != nil {
			return fmt.Errorf("init enum %s: %w", enum.Name, err)
		}

		t.Enums[i] = enum
	}

	for i, model := range t.Models {
		if err := model.Init(config, moduleName); err != nil {
			return fmt.Errorf("init model %s: %w", model.Name, err)
		}

		t.Models[i] = model
	}

	for i, oneOf := range t.OneOfs {
		if err := oneOf.Init(config, moduleName); err != nil {
			return fmt.Errorf("init one_of %s: %w", oneOf.Name, err)
		}

		t.OneOfs[i] = oneOf
	}

	return nil
}
