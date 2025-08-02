# je - JSON Editor CLI

A command-line tool for editing JSON files in-place using HTTPie-style syntax.

## Installation

```bash
go install github.com/vampire/je/cmd/je@latest
```

Or build from source:

```bash
git clone https://github.com/vampire/je
cd je
go build -o je ./cmd/je
```

## Usage

```bash
je <file> <assignments...> [options]
```

Edit JSON files using intuitive key=value syntax:

```bash
# Simple edits
je config.json port:=3000 host=localhost debug:=true

# Nested paths
je user.json profile.name=gary profile.age:=30

# Arrays
je data.json users.0.active:=false tags[]=newitem

# From files
je config.json ssl.cert@cert.pem settings:@config.json

# Multiple files
je '*.json' --each version=1.2.0

# View changes
je config.json --diff --dry-run port:=8080
```

## Assignment Syntax

### Basic Types
- `key=value` - Set string value
- `key:=value` - Set raw JSON (number, boolean, null, array, object)
- `key@file` - Set value from file contents
- `key:@file` - Set raw JSON from file

### Path Notation
- `user.name=gary` - Nested object access
- `users.0.name=gary` - Array index access
- `users.[].active:=true` - Set property on all array elements
- `config.ports[]=8080` - Append to array
- `tags[]="new"` - Append string to array

### Special Operations
- `user.age:=null` - Set null value
- `user.phone=` - Set empty string
- `user.phone:=` - Delete key
- `metadata:={}` - Empty object
- `items:=[]` - Empty array

## Options

```
-i, --in-place          Edit file in place (default)
-o, --output <file>     Write to different file
-p, --pretty            Pretty print output
-c, --compact           Compact output
-r, --raw               Output raw values (no JSON encoding)
-e, --each              Apply to multiple files independently
-n, --dry-run           Show changes without writing
-d, --diff              Show diff of changes
-q, --quiet             Suppress non-error output
--create                Create file if doesn't exist
--merge                 Merge instead of overwrite arrays/objects
--json5                 Parse/write JSON5
```

## Examples

### Basic Editing

```bash
# Update server configuration
je server.json host=0.0.0.0 port:=8080 ssl.enabled:=true

# Set nested values
je package.json scripts.test="npm test" scripts.build="npm run build"

# Work with arrays
je data.json users.0.admin:=true users.1.admin:=false
```

### Array Operations

```bash
# Append to arrays
je config.json servers[]="server3.example.com" ports[]:=8080

# Update all array elements
je users.json 'users.[].status=active' 'users.[].verified:=true'
```

### File Operations

```bash
# Load certificate from file
je config.json ssl.cert@/path/to/cert.pem ssl.key@/path/to/key.pem

# Load entire config from file
je app.json database:@db-config.json
```

### Batch Operations

```bash
# Update multiple files
je '*.json' --each version=2.0.0 updated:=true

# Preview changes with diff
je config.json --diff --dry-run port:=3000

# Quiet mode for scripts
je data.json --quiet status=processed
```

### Working with stdin/stdout

```bash
# Pipe support
echo '{}' | je - name=test | je - age:=25

# Read from stdin, write to file
cat input.json | je - name=updated > output.json
```

## Advanced Features

### Escaping Special Characters

```bash
# Escape dots in keys
je file.json 'user\.name=gary'

# Keys with brackets
je file.json 'map["key"]=value'

# Values with spaces
je file.json 'message=Hello World'
```

### Complex Data Types

```bash
# Set objects
je config.json 'database:={"host":"localhost","port":5432}'

# Set arrays
je config.json 'tags:=["prod","stable","v1"]'

# Null values
je user.json email:=null
```

## Error Handling

- Missing files result in an error unless `--create` is used
- Invalid JSON causes an error with line number
- Type conflicts (e.g., indexing a string as array) result in clear error messages
- File permissions are preserved during in-place edits
- Atomic writes ensure data safety (temp file + rename)

## Performance

- Optimized for files under 100MB
- Memory usage approximately 2x file size
- Sub-50ms operations for typical config files

## License

MIT