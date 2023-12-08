package txoutboxdbstate

import "embed"

//go:embed *_pg.sql
var PGMigrationsFS embed.FS
