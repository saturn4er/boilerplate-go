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
type ModelTypeParameter struct {
	Name       string `yaml:"name"`
	Constraint string `yaml:"constraint"`
}
type QualifiedFunc struct {
	Package string
	Func    string
}

type UniqueIndex struct {
	ConstraintName string   `yaml:"constraint_name"`
	Fields         []string `yaml:"fields"`
}

// ErrorName returns the joined field names used for error constant naming.
// e.g. ["Name", "OrgID"] → "NameOrgID"
func (u UniqueIndex) ErrorName() string {
	return strings.Join(u.Fields, "")
}

type Model struct {
	ID                uint                 `yaml:"id"`
	Admin             ConfigModelAdmin     `yaml:"admin"`
	Package           string               `yaml:"package"`
	StorageType       ModelStorageType     `yaml:"storage_type"`
	TypeParameters    []ModelTypeParameter `yaml:"type_parameters"`
	Name              string               `yaml:"name"`
	Fields            []ModelField         `yaml:"fields"`
	PluralName        string               `yaml:"plural_name"`
	DoNotPersists     bool                 `yaml:"do_not_persists"`
	TableName         string               `yaml:"table_name"`
	NoLocalOutbox     bool                 `yaml:"no_local_outbox"`
	UniqueIndexes     []UniqueIndex        `yaml:"unique_indexes"`
	MessageBuilder    *QualifiedFunc       `yaml:"-"`
	RawMessageBuilder string               `yaml:"message_builder"`
}

func (c *Model) FirstPKField() ModelField {
	var idField *ModelField
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
	for i, field := range c.Fields {
		if err := field.Init(config, moduleName); err != nil {
			return err
		}

		c.Fields[i] = field
	}

	if c.PluralName == "" {
		c.PluralName = c.Name + "s"
	}

	for i, idx := range c.UniqueIndexes {
		if idx.ConstraintName == "" {
			return fmt.Errorf("model %s: unique_indexes[%d] must have a non-empty constraint_name", c.Name, i)
		}
		if len(idx.Fields) == 0 {
			return fmt.Errorf("model %s: unique_indexes[%d] must have at least one field", c.Name, i)
		}
	}

	if c.RawMessageBuilder != "" {
		idx := strings.LastIndex(c.RawMessageBuilder, ".")
		if idx < 0 {
			return fmt.Errorf("model %s: message_builder %q must be a fully qualified function (package.FuncName)", c.Name, c.RawMessageBuilder)
		}

		c.MessageBuilder = &QualifiedFunc{
			Package: c.RawMessageBuilder[:idx],
			Func:    c.RawMessageBuilder[idx+1:],
		}
	}

	return nil
}

type ConfigModelFieldAdmin struct {
	HideForList bool   `yaml:"hide_for_list"`
	Hide        bool   `yaml:"hide"`
	Creatable   bool   `yaml:"creatable"`
	Editable    bool   `yaml:"editable"`
	LinkTo      string `yaml:"link_to"`
}

func (s *ConfigModelFieldAdmin) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain ConfigModelFieldAdmin
	if err := unmarshal((*plain)(s)); err != nil {
		return err
	}

	return nil
}

type ModelField struct {
	Name          string                `yaml:"name"`
	DBName        string                `yaml:"database_name"`
	JSONName      string                `yaml:"json_name"`
	Type          Type                  `yaml:"type"`
	Filterable    bool                  `yaml:"filterable"`
	DoNotPersists bool                  `yaml:"do_not_persists"`
	PrimaryKey    bool                  `yaml:"primary_key"`
	Admin         ConfigModelFieldAdmin `yaml:"admin"`
}

func (c *ModelField) UnmarshalYAML(unmarshal func(interface{}) error) error {
	c.Admin = ConfigModelFieldAdmin{
		Editable:  true,
		Creatable: true,
	}

	type plain ModelField
	if err := unmarshal((*plain)(c)); err != nil {
		return fmt.Errorf("unmarshal model: %w", err)
	}

	return nil
}

func (c *ModelField) Init(config *Config, moduleName string) error {
	if c.DBName == "" {
		c.DBName = strcase.SnakeCase(c.Name)
	}

	if c.JSONName == "" {
		c.JSONName = strcase.SnakeCase(c.Name)
	}

	if err := c.Type.Init(config, moduleName); err != nil {
		return fmt.Errorf("init field %s type: %w", c.Name, err)
	}

	return nil
}
