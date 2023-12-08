package scaffoldtpl

import "embed"

//go:embed all:storage/*.tpl
//sgo:embed all:admin/*.tpl
//go:embed all:service/*.tpl
var FS embed.FS
