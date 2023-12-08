package scaffoldtpl

import "embed"

//go:embed all:service/*.tpl
//go:embed all:storage/*.tpl
//go:embed all:admin/*.tpl
var FS embed.FS
