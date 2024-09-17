package config

type (
	EnumHelpers struct {
		isTrue     bool
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
	Enum struct {
		Name    string      `yaml:"name"`
		Package string      `yaml:"package"`
		Values  []string    `yaml:"values"`
		Helpers EnumHelpers `yaml:"helpers"`
	}
)

func (e *EnumHelpers) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var boolVal bool
	if err := unmarshal(&boolVal); err == nil {
		if !boolVal {
			return nil
		}
		e.isTrue = true
		e.IsValid = true
		e.Is = true
		e.Stringer = true
		e.Validate = true

		return nil
	}

	type plain EnumHelpers
	if err := unmarshal((*plain)(e)); err != nil {
		return err
	}

	return nil
}
func (e *EnumHelpers) Init(config *Config, moduleName, enumName string) error {
	if !e.isTrue {
		return nil
	}
	e.AllValues.VarName = "All" + enumName

	return nil
}
func (e *Enum) Init(cfg *Config, moduleName string) error {
	return e.Helpers.Init(cfg, moduleName, e.Name)
}
