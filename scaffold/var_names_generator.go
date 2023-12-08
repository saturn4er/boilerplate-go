package scaffold

import "strconv"

type varNamesGenerator struct {
	usedNames map[string]bool
}

func (v *varNamesGenerator) Var(name string) string {
	if v.usedNames == nil {
		v.usedNames = map[string]bool{}
	}

	if _, exists := v.usedNames[name]; !exists {
		v.usedNames[name] = true

		return name
	}

	for i := 1; ; i++ {
		suffixedName := name + strconv.Itoa(i)
		if _, exists := v.usedNames[suffixedName]; !exists {
			v.usedNames[suffixedName] = true

			return suffixedName
		}
	}
}
