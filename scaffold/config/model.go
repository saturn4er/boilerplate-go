package config

import (
	"fmt"
	"strings"

	"github.com/stoewer/go-strcase"
)

type ModelStorageType string

const (
	ModelStorageTypeTxOutbox ModelStorageType = "tx_outbox"
)

type ConfigModelAdmin struct {
	Customizable bool `yaml:"customizable"`
}

type Model struct {
	ID                 uint               `yaml:"id"`
	PKIndexName        string             `yaml:"pk_index_name"`
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

func (c *Model) FirstPKField() ConfigModelField {
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

func (c *Model) Init(config *Config, moduleName string) error {
	if !c.DoNotPersists && c.StorageType == "" && c.PKIndexName == "" {
		return fmt.Errorf("model %s has no pk_index_name specified", c.Name)
	}

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
	Type          Type                  `yaml:"type"`
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
