package scaffold

import (
	"plugin"

	"github.com/go-pnp/go-pnp/pkg/optionutil"
)

func WithOutputDir(dir string) optionutil.Option[generatorOptions] {
	return func(options *generatorOptions) {
		options.OutputDir = dir
	}
}

func WithPlugins(plugins ...plugin.Plugin) optionutil.Option[generatorOptions] {
	return func(options *generatorOptions) {
		options.Plugins = append(options.Plugins, plugins...)
	}
}
