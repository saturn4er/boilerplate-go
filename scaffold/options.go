package scaffold

import "github.com/go-pnp/go-pnp/pkg/optionutil"

func WithOutputDir(dir string) optionutil.Option[generatorOptions] {
	return func(options *generatorOptions) {
		options.OutputDir = dir
	}
}
