// Package scaffold is the core code generation engine for boilerplate-go.
//
// It reads YAML configuration files describing modules with models, enums,
// one-of types, and events, then generates Go source files using text/template.
//
// # Generation Pipeline
//
// The [Generate] function orchestrates the full pipeline:
//
//  1. Load and parse YAML configuration via [config.Load].
//  2. Load embedded Go templates from the scaffoldtpl package.
//  3. For each module, evaluate template conditions and execute applicable templates.
//  4. Extract and preserve user code blocks from existing generated files.
//  5. Format output with goimports and write to disk.
//
// # Templates
//
// Templates are embedded .tpl files with a JSON header specifying the output
// file path, package name, and an optional condition expression. Helper
// templates (prefixed with ".") provide shared functions across all templates.
//
// # User Code Preservation
//
// Generated files support preserved code blocks delimited by special comments:
//
//	// user code 'block_name'
//	... custom code here ...
//	// end user code 'block_name'
//
// These blocks survive regeneration, allowing developers to extend generated
// code without losing changes.
package scaffold
