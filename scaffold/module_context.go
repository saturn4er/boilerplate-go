package scaffold

import "github.com/saturn4er/boilerplate-go/scaffold/config"

type ModuleContext struct {
	Module       string
	Config       *config.Config
	ModuleConfig *config.Config
	Env          map[string]string
}
