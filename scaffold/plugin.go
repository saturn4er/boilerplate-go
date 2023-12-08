package scaffold

import (
	"github.com/saturn4er/boilerplate-go/scaffold/config"
)

type Plugin interface {
	// Init is called after config is loaded and before generating
	Init(*config.Config) error
	Name() string
	Templates() []GeneratorTemplate
}
