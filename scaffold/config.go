package scaffold

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/stoewer/go-strcase"
	"gopkg.in/yaml.v3"
)

type importableConfig[T any] struct {
	Value T
}

func (c *importableConfig[T]) UnmarshalYAML(unmarshal func(any) error) error {
	var configWithImport struct {
		Import string `yaml:"import"`
	}

	if err := unmarshal(&configWithImport); err != nil {
		return fmt.Errorf("unmarshal import config: %w", err)
	}

	if configWithImport.Import == "" {
		return unmarshal(&c.Value)
	}

	importContent, err := os.ReadFile(configWithImport.Import)
	if err != nil {
		return fmt.Errorf("read import file '%s': %w", configWithImport.Import, err)
	}

	return yaml.Unmarshal(importContent, &c.Value)
}

type ConfigOneOfValue struct {
	ModelName string
	Model     *ConfigModel `yaml:"-"`
}

func (c *ConfigOneOfValue) UnmarshalYAML(unmarshal func(any) error) error {
	var strValue string
	if err := unmarshal(&strValue); err == nil {
		c.ModelName = strValue

		return nil
	}

	type resultType ConfigOneOfValue

	var result resultType

	if err := unmarshal(&result); err != nil {
		return err
	}

	*c = ConfigOneOfValue(result)

	return nil
}

func (c *ConfigOneOfValue) Init(config *Config, moduleName string) error {
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

type ConfigOneOf struct {
	Name   string             `yaml:"name"`
	Values []ConfigOneOfValue `yaml:"values"`
}

func (c *ConfigOneOf) SortedValues() []indexed[uint, ConfigOneOfValue] {
	result := make([]indexed[uint, ConfigOneOfValue], 0, len(c.Values))
	for _, value := range c.Values {
		result = append(result, indexed[uint, ConfigOneOfValue]{Index: value.Model.ID, Value: value})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Index < result[j].Index
	})

	return result
}

func (c *ConfigOneOf) Init(config *Config, moduleName string) error {
	for i, value := range c.Values {
		if err := value.Init(config, moduleName); err != nil {
			return fmt.Errorf("init one of value: %w", err)
		}

		c.Values[i] = value
	}

	return nil
}

type (
	ConfigEnumHelpers struct {
		IsValid    bool `yaml:"is_valid"`
		Is         bool `yaml:"is"`
		IsCategory []struct {
			Name   string   `yaml:"name"`
			Values []string `yaml:"values"`
		} `yaml:"is_category"`
		Validate  bool `yaml:"validate"`
		Stringer  bool `yaml:"stringer"`
		AllValues struct {
			VarName  string `yaml:"var_name"`
			FuncName string `yaml:"func_name"`
		} `yaml:"all_values"`
	}
	ConfigEnum struct {
		Name    string            `yaml:"name"`
		Package string            `yaml:"package"`
		Values  []string          `yaml:"values"`
		Helpers ConfigEnumHelpers `yaml:"helpers"`
	}
)

func (e *ConfigEnum) Init(*Config, string) error {
	return nil
}

type ConfigType struct {
	Package  string      `yaml:"package"`
	Type     string      `yaml:"type"`
	ElemType *ConfigType `yaml:"elem"`
	KeyType  *ConfigType `yaml:"key"`
}

func (c *ConfigType) UnmarshalYAML(unmarshal func(any) error) error {
	var (
		strValue  string
		parseType func(val string) *ConfigType
	)

	parseType = func(val string) *ConfigType {
		switch {
		case strings.HasPrefix(val, "*"):
			return &ConfigType{
				Type:     "ptr",
				ElemType: parseType(val[1:]),
			}
		case strings.HasPrefix(val, "[]"):
			return &ConfigType{
				Type:     "slice",
				ElemType: parseType(val[2:]),
			}
		case strings.HasPrefix(val, "map["):
			// Search for the closing bracket considering nested maps and slices and other chars
			var (
				bracketCount        = 1
				closingBracketIndex int
			)

			for i, char := range val[4:] { //nolint:varnamelen
				switch char {
				case '[':
					bracketCount++
				case ']':
					bracketCount--
				}

				if bracketCount == 0 {
					closingBracketIndex = i

					break
				}
			}

			return &ConfigType{
				Type:     "map",
				KeyType:  parseType(val[4 : 4+closingBracketIndex]),
				ElemType: parseType(val[4+closingBracketIndex+1:]),
			}
		}

		return &ConfigType{
			Type: val,
		}
	}

	if err := unmarshal(&strValue); err == nil {
		*c = *parseType(strValue)

		return nil
	}

	type resultType ConfigType

	var result resultType

	if err := unmarshal(&result); err != nil {
		return err
	}

	*c = ConfigType(result)

	return nil
}

func (c *ConfigType) Init(_ *Config, _ string) error {
	if c.Type == "slice" && c.ElemType == nil {
		return errors.New("slice must have an element type")
	}

	if c.Type == "map" && c.KeyType == nil && c.ElemType == nil {
		return errors.New("map must have a key and elem type")
	}

	return nil
}

type ConfigModelFieldAdmin struct {
	HideForList bool   `yaml:"hide_for_list"`
	Hide        bool   `yaml:"hide"`
	Editable    bool   `yaml:"editable"`
	LinkTo      string `yaml:"link_to"`
}

func (s *ConfigModelFieldAdmin) UnmarshalYAML(unmarshal func(interface{}) error) error {
	s.Editable = true

	type plain ConfigModelFieldAdmin
	if err := unmarshal((*plain)(s)); err != nil {
		return err
	}

	return nil
}

type ConfigModelField struct {
	Name          string                `yaml:"name"`
	DBName        string                `yaml:"database_name"`
	Type          ConfigType            `yaml:"type"`
	Filterable    bool                  `yaml:"filterable"`
	DoNotPersists bool                  `yaml:"do_not_persists"`
	PrimaryKey    bool                  `yaml:"primary_key"`
	Admin         ConfigModelFieldAdmin `yaml:"admin"`
}

func (c *ConfigModelField) UnmarshalYAML(unmarshal func(interface{}) error) error {
	c.Admin = ConfigModelFieldAdmin{
		Editable: true,
	}

	type plain ConfigModelField
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}

	return nil
}

func (c *ConfigModelField) Init(config *Config, moduleName string) error {
	if c.DBName == "" {
		c.DBName = strcase.SnakeCase(c.Name)
	}

	if err := c.Type.Init(config, moduleName); err != nil {
		return fmt.Errorf("init field %s type: %w", c.Name, err)
	}

	return nil
}

type ModelStorageType string

const (
	ModelStorageTypeTxOutbox ModelStorageType = "tx_outbox"
)

type ConfigModelAdmin struct {
	Customizable bool `yaml:"customizable"`
}

type ConfigModel struct {
	ID                 uint               `yaml:"id"`
	Admin              ConfigModelAdmin   `yaml:"admin"`
	Package            string             `yaml:"package"`
	StorageType        ModelStorageType   `yaml:"storage_type"`
	Name               string             `yaml:"name"`
	Fields             []ConfigModelField `yaml:"fields"`
	PluralName         string             `yaml:"plural_name"`
	DoNotPersists      bool               `yaml:"do_not_persists"`
	HasCustomDBMethods bool               `yaml:"has_custom_db_methods"`
	TableName          string             `yaml:"table_name"`
}

func (c *ConfigModel) FirstPKField() ConfigModelField {
	var idField *ConfigModelField
	for _, field := range c.Fields {
		if field.PrimaryKey == true {
			return field
		}
		if strings.ToLower(field.Name) == "id" {
			fieldCp := field
			idField = &fieldCp
		}
	}
	if idField != nil {
		return *idField
	}

	panic("no primary key field found")
}

func (c *ConfigModel) Init(config *Config, moduleName string) error {
	for i, field := range c.Fields {
		if err := field.Init(config, moduleName); err != nil {
			return err
		}

		c.Fields[i] = field
	}

	if c.PluralName == "" {
		c.PluralName = c.Name + "s"
	}

	return nil
}

type ConfigTypes struct {
	Enums  []ConfigEnum  `yaml:"enums"`
	Models []ConfigModel `yaml:"models"`
	OneOfs []ConfigOneOf `yaml:"one_ofs"`
}

func (t *ConfigTypes) GetTypeByName(name string) (any, bool) {
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

func (t *ConfigTypes) GetModelByName(name string) (ConfigModel, bool) {
	for _, model := range t.Models {
		if model.Name == name {
			return model, true
		}
	}

	return ConfigModel{}, false
}

func (t *ConfigTypes) GetEnumByName(name string) (ConfigEnum, bool) {
	for _, enum := range t.Enums {
		if enum.Name == name {
			return enum, true
		}
	}

	return ConfigEnum{}, false
}

func (t *ConfigTypes) GetOneOfByName(name string) (ConfigOneOf, bool) {
	for _, oneOf := range t.OneOfs {
		if oneOf.Name == name {
			return oneOf, true
		}
	}

	return ConfigOneOf{}, false
}

func (t *ConfigTypes) Init(config *Config, moduleName string) error {
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

type ConfigModule struct {
	Types ConfigTypes `yaml:"types"`
}

func (m *ConfigModule) Init(config *Config, moduleName string) error {
	if err := m.Types.Init(config, moduleName); err != nil {
		return fmt.Errorf("init types: %w", err)
	}

	return nil
}

type Config struct {
	RootPackageName string                                    `yaml:"root_package_name"`
	GoImportsLocal  string                                    `yaml:"goimports_local"`
	Types           ConfigTypes                               `yaml:"types"`
	Modules         map[string]importableConfig[ConfigModule] `yaml:"modules"`
	Module          string                                    `yaml:"-"`
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

func LoadConfig(path string) (*Config, error) {
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
