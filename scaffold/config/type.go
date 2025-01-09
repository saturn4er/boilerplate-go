package config

import (
	"errors"
	"fmt"
	"strings"
)

const (
	TypeSlice = "slice"
	TypeMap   = "map"
	TypePtr   = "ptr"
)

type Type struct {
	Package        string `yaml:"package"`
	Type           string `yaml:"type"`
	TypeParameters []Type `yaml:"type_parameters"`
	ElemType       *Type  `yaml:"elem"`
	KeyType        *Type  `yaml:"key"`
}

func (c *Type) IsPrimitive() bool {
	switch c.Type {
	case "bool", "string", "byte", "rune",
		"int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"uintptr", "float32", "float64", "complex64", "complex128":
		return true
	}

	return false
}

func (c *Type) UnmarshalYAML(unmarshal func(any) error) error {
	var strValue string
	// if value is a string, parse it from string to type
	if err := unmarshal(&strValue); err == nil {
		*c = *newTypeFromString(strValue)

		return nil
	}

	type resultType Type

	var result resultType

	if err := unmarshal(&result); err != nil {
		return fmt.Errorf("unmarshal type: %w", err)
	}

	*c = Type(result)

	return nil
}

func (c *Type) Init(_ *Config, _ string) error {
	if c.Type == TypeSlice && c.ElemType == nil {
		return errors.New("slice must have an element type")
	}

	if c.Type == TypeMap && c.KeyType == nil && c.ElemType == nil {
		return errors.New("map must have a key and elem type")
	}
	if c.IsPrimitive() {
		if c.Package != "" {
			return errors.New("primitive type must not have a package specified")
		}
		if c.ElemType != nil {
			return errors.New("primitive type must not have an element type specified")
		}
		if c.KeyType != nil {
			return errors.New("primitive type must not have a key type specified")
		}
	}

	return nil
}

func newTypeFromString(val string) *Type {
	switch {
	case strings.HasPrefix(val, "*"):
		return &Type{
			Type:     TypePtr,
			ElemType: newTypeFromString(val[1:]),
		}
	case strings.HasPrefix(val, "[]"):
		return &Type{
			Type:     TypeSlice,
			ElemType: newTypeFromString(val[2:]),
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

		return &Type{
			Type:     TypeMap,
			KeyType:  newTypeFromString(val[4 : 4+closingBracketIndex]),
			ElemType: newTypeFromString(val[4+closingBracketIndex+1:]),
		}
	}

	return &Type{
		Type: val,
	}
}
