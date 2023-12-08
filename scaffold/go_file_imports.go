package scaffold

import (
	"strconv"
	"strings"
)

type codeGeneratorImport struct {
	generator *fileGenerator

	ImportPath string
	Alias      string
}

func (c codeGeneratorImport) Ref(val string) string {
	for _, imprt := range c.generator.imports {
		if imprt.ImportPath == c.ImportPath {
			return imprt.Alias + "." + val
		}
	}

	pathParts := strings.Split(c.ImportPath, "/")
	if c.Alias == "" {
		c.Alias = pathParts[len(pathParts)-1]
	}

	// check if alias is already used and add a number suffix if it is
	for i := 0; ; i++ {
		aliasUsed := false
		suffix := ""

		if i > 0 {
			suffix = strconv.FormatInt(int64(i), 10)
		}

		for _, imprt := range c.generator.imports {
			if imprt.Alias == c.Alias+suffix {
				aliasUsed = true

				break
			}
		}

		if !aliasUsed {
			c.Alias += suffix

			break
		}
	}

	c.generator.imports = append(c.generator.imports, c)

	return c.Alias + "." + val
}
