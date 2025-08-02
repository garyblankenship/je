# Project Memory

## TASKS
- [x] Build MVP je tool with basic functionality
- [x] Add support for array map operations (users.[].active:=true)
- [x] Add --each flag for multiple file operations
- [x] Add --diff flag to show changes
- [x] Create comprehensive README.md
- [x] Add .gitignore and LICENSE
- [x] Implement all audit recommendations (A- to A+ improvements)
- [x] Refactor large functions into smaller units
- [x] Fix error wrapping to use %w throughout
- [x] Add edge case tests for array operations
- [x] Add linter configuration (.golangci.yml)
- [ ] Add --merge flag for arrays/objects
- [ ] Add JSON5 support
- [ ] Add JSON Schema validation support
- [ ] Optimize for large files (streaming)
- [ ] Publish to GitHub

## REFERENCE  
### Key Files & Patterns
- `cmd/je/main.go` - CLI entry point using Cobra
- `internal/parser/` - Argument parsing (key=value syntax)
- `internal/json/` - File I/O and formatting
- `internal/operations/` - JSON manipulation using gjson/sjson
- `internal/diff/` - Diff display functionality
- `internal/cli/` - CLI helper functions
- Modular design with separate files for array_ops, array_map, json_parser

### Architecture Decisions
- Using tidwall/gjson and tidwall/sjson for path-based JSON operations
- Atomic file writes (temp file + rename) for safety
- Operator precedence: []= checked before = to handle array syntax correctly
- Path segments split by dots, escaped dots handled with backslash

### Cached Research
- Build: `go build -o je ./cmd/je`
- Test: `go test ./...` and `./test_e2e.sh`
- Lint: `golangci-lint run ./...`
- Dependencies: cobra (CLI), gjson/sjson (JSON manipulation), go-diff (diff display)
- Array map syntax: `users.[].active:=true` sets property on all array elements
- Parser must handle [].[property] carefully to extract base path correctly
- --each flag uses filepath.Glob for pattern matching
- --diff uses line-by-line comparison with color output
- Functions refactored to stay under 80 lines and 20 cyclomatic complexity
- Error wrapping uses %w for proper context propagation

### Environment Quirks
- JSON key ordering is not guaranteed, tests must be order-agnostic
- File operations need explicit permission handling
- stdin is represented as "-" filename