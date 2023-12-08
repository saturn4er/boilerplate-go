package config

import (
	"fmt"
	"sort"

	"github.com/samber/lo"
)

type indexed[T comparable, V any] struct {
	Index T
	Value V
}

type OneOfValue struct {
	ModelName string
	Model     *Model `yaml:"-"`
}

func (c *OneOfValue) UnmarshalYAML(unmarshal func(any) error) error {
	var strValue string
	if err := unmarshal(&strValue); err == nil {
		c.ModelName = strValue

		return nil
	}

	type resultType OneOfValue

	var result resultType

	if err := unmarshal(&result); err != nil {
		return err
	}

	*c = OneOfValue(result)

	return nil
}

func (c *OneOfValue) Init(config *Config, moduleName string) error {
	var found bool

	for _, model := range config.Modules[moduleName].Value.Types.Models {
		modelCopy := model

		if model.Name == c.ModelName {
			if model.ID == 0 {
				return fmt.Errorf("one of value model '%s' must have id", model.Name)
			}

			c.Model = &modelCopy
			found = true

			break
		}
	}

	if !found {
		return fmt.Errorf("model %s not found", c.ModelName)
	}

	return nil
}

type OneOf struct {
	Name   string       `yaml:"name"`
	Values []OneOfValue `yaml:"values"`
}

func (c *OneOf) SortedValues() []indexed[uint, OneOfValue] {
	result := lo.Map(c.Values, func(value OneOfValue, _ int) indexed[uint, OneOfValue] {
		return indexed[uint, OneOfValue]{Index: value.Model.ID, Value: value}
	})

	sort.Slice(result, func(i, j int) bool {
		return result[i].Index < result[j].Index
	})

	return result
}

func (c *OneOf) Init(config *Config, moduleName string) error {
	for i, value := range c.Values {
		if err := value.Init(config, moduleName); err != nil {
			return fmt.Errorf("init one of value: %w", err)
		}

		c.Values[i] = value
	}

	return nil
}
