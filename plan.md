# je - JSON Editor CLI Specification

## Synopsis
```bash
je <file> <assignments...> [options]
```

Edit JSON files in-place using HTTPie-style key=value syntax.

## Assignment Syntax

### Basic Types
- `key=value` - Set string
- `key:=value` - Set raw JSON (number, boolean, null, array, object)
- `key@file` - Set value from file contents
- `key:@file` - Set raw JSON from file

### Path Notation
- `user.name=gary` - Nested object
- `users.0.name=gary` - Array index
- `users.[].active:=true` - Array map (all elements)
- `config.ports[]=8080` - Array append
- `tags[]="new"` - Array append string

### Special Operations
- `user.age:=null` - Set null
- `user.phone=` - Set empty string
- `user.phone:=` - Delete key
- `metadata:={}` - Empty object
- `items:=[]` - Empty array

## Examples

```bash
# Simple edits
je config.json port:=3000 host=localhost debug:=true

# Nested paths
je user.json profile.name=gary profile.age:=30

# Arrays
je data.json users.0.active:=false results[]=newitem

# From files
je config.json ssl.cert@cert.pem settings:@config.json

# Multiple files (sequentially)
je *.json --each version=1.2.0

# Pipe support
echo '{}' | je - name=test | je - age:=25
```

## Options

```bash
-i, --in-place          Edit file in place (default)
-o, --output <file>     Write to different file
-p, --pretty            Pretty print output
-c, --compact           Compact output
-r, --raw               Output raw values (no JSON encoding)
-e, --each              Apply to multiple files independently
-n, --dry-run          Show changes without writing
-d, --diff             Show diff of changes
-q, --quiet            Suppress non-error output
--create               Create file if doesn't exist
--merge                Merge instead of overwrite arrays/objects
--json5                Parse/write JSON5
```

## Behavior

### Creation Rules
- Intermediate objects created automatically
- `user.profile.name=x` on `{}` creates `{"user":{"profile":{"name":"x"}}}`
- Array indices must be sequential or use `[]` append

### Type Coercion
- Overwriting value changes type implicitly
- `{"age":"30"}` + `age:=30` → `{"age":30}`
- Arrays can only be set whole or appended

### Error Handling
- Missing file: error unless `--create`
- Invalid JSON: error with line number
- Invalid path: error with suggestion
- Type conflicts: error (can't index string as array)

### Atomicity
- Write to temp file, rename on success
- Original preserved on error
- File permissions maintained

## Edge Cases

```bash
# Escape special chars in keys
je file.json 'user\.name=gary'          # Key with dot
je file.json 'map["key"]=value'         # Key with brackets
je file.json 'name=value with spaces'   # Value with spaces

# Numeric keys
je file.json items.0=first items.1=second

# Unicode
je file.json name=🦀 emoji:=true

# Large values
je config.json data:=@large.json       # Multi-MB embeds

# stdin/stdout
je - name=test < input.json > output.json
cat file.json | je - name=test
```

## Implementation Notes

### Parser Phases
1. Tokenize arguments into (path, operator, value) tuples
2. Parse paths into segment arrays
3. Validate value types based on operator
4. Build operation list

### Update Algorithm
```
for each operation:
  if path exists:
    update value
  else if parent exists:
    create key/index
  else:
    create parent path recursively
```

### Performance Targets
- < 50ms for 1MB file
- Streaming for files > 100MB
- Memory usage ≤ 2x file size

### Dependencies
- `tidwall/gjson` - Path queries
- `tidwall/sjson` - Path updates  
- `spf13/cobra` - CLI framework
- No external runtime requirements

## Future Extensions
- `--schema <file>` - Validate against JSON Schema
- `--select <query>` - Update only matching elements
- `--transform <expr>` - Apply expression to values
- Plugin system for custom operators

**NEXT**: Start with MVP covering basic string/number assignment and simple paths, then add array operations and file inputs.