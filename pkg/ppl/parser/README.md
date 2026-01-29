# PPL Parser

This package contains the ANTLR4 grammar and parser implementation for Quidditch's Piped Processing Language (PPL).

## Overview

The parser converts PPL query strings into Abstract Syntax Trees (AST) for further processing. It consists of:

- **ANTLR4 Grammar Files** (`PPLLexer.g4`, `PPLParser.g4`) - Define the syntax
- **Parser Wrapper** (`parser.go`) - Go interface for parsing queries
- **AST Builder** (`ast_builder.go`) - Converts ANTLR4 parse trees to AST
- **Error Handling** (`error_listener.go`) - Enhanced error reporting

## Prerequisites

### Install ANTLR4

#### Option 1: Download JAR (Recommended)
```bash
cd pkg/ppl/parser
make  # Automatically downloads ANTLR4 JAR and generates code
```

#### Option 2: Install via package manager

**macOS:**
```bash
brew install antlr
```

**Linux (apt):**
```bash
sudo apt-get install antlr4
```

**Arch Linux:**
```bash
sudo pacman -S antlr4
```

### Install ANTLR4 Go Runtime

```bash
go get github.com/antlr4-go/antlr/v4
```

Or run:
```bash
make install-runtime
```

## Building

Generate the parser code from ANTLR4 grammar:

```bash
cd pkg/ppl/parser
make
```

This will:
1. Download ANTLR4 JAR if not present
2. Generate Go code in `generated/` directory
3. Create lexer and parser files

The generated files are:
- `generated/ppllexer_lexer.go` - Tokenizer
- `generated/pplparser_parser.go` - Parser
- `generated/pplparser_base_listener.go` - Base listener
- `generated/pplparser_listener.go` - Listener interface

## Usage

### Basic Parsing

```go
package main

import (
    "fmt"
    "log"

    "github.com/quidditch/quidditch/pkg/ppl/parser"
)

func main() {
    // Create parser
    p := parser.NewParser()

    // Parse query
    query := "source=logs | where status = 200 | stats count() by endpoint"
    ast, err := p.Parse(query)
    if err != nil {
        log.Fatal(err)
    }

    // Use AST
    fmt.Printf("Parsed query with %d commands\n", len(ast.Commands))
}
```

### Syntax Validation

```go
// Just validate syntax without building AST
err := parser.ValidateSyntax("source=logs | where status = 200")
if err != nil {
    fmt.Println("Syntax error:", err)
}
```

### Error Handling

```go
p := parser.NewParser()
_, err := p.Parse("invalid | query |")
if err != nil {
    // Enhanced error message with line/column
    fmt.Println(err)
    // Output: syntax error at line 1, column 15: unexpected '|'
}
```

## Supported PPL Features (Tier 0)

### Commands

- **search/source** - Data source selection: `source=logs`
- **where** - Filtering: `where status = 200`
- **fields** - Field selection: `fields timestamp, message` or `fields - internal_id`
- **stats** - Aggregation: `stats count(), avg(time) by status`
- **sort** - Ordering: `sort timestamp desc`
- **head** - Limit results: `head 100`
- **describe** - Show schema: `describe logs`
- **showdatasources** - List available sources
- **explain** - Show query execution plan

### Expressions

#### Logical Operators
- `AND`, `OR`, `NOT`
- Example: `where status = 200 AND method = 'GET'`

#### Comparison Operators
- `=`, `!=`, `<`, `>`, `<=`, `>=`
- `LIKE` - Pattern matching
- `IN` - List membership
- Example: `where status IN (200, 201, 204)`

#### Arithmetic Operators
- `+`, `-`, `*`, `/`, `%`
- Example: `where response_time > avg_time * 1.5`

#### Aggregation Functions
- `count()`, `sum()`, `avg()`, `min()`, `max()`
- With DISTINCT: `count(DISTINCT user_id)`

#### Case Expressions
```ppl
stats count() as total by
  case
    when status < 300 then 'success'
    when status < 500 then 'client_error'
    else 'server_error'
  end
```

### Field References

- Simple: `status`, `response_time`
- Nested: `user.address.city`
- Array indexing: `tags[0]`

### Literals

- **Integer**: `42`, `-10`
- **Decimal**: `3.14`, `-0.5`
- **String**: `'hello'`, `"world"`, `` `backticks` ``
- **Boolean**: `true`, `false`
- **Null**: `null`

## Examples

### Simple Query
```ppl
source=logs | where status = 200
```

### Analytics Query
```ppl
source=logs
| where timestamp > '2024-01-01' AND method = 'GET'
| stats count() as requests, avg(response_time) as avg_time by endpoint
| sort requests desc
| head 10
```

### Complex Filter
```ppl
source=logs
| where (status >= 200 AND status < 300) OR status = 304
| where response_time > 1000
| fields timestamp, endpoint, status, response_time
```

### Schema Inspection
```ppl
describe logs
```

```ppl
showdatasources
```

### Query Explanation
```ppl
explain source=logs | where status = 200 | stats count() by endpoint
```

## Testing

Run tests (requires generated code):

```bash
# Generate code first
make

# Run tests
make test
```

Or with go test:
```bash
go test -v ./...
```

## Grammar Development

### Modifying the Grammar

1. Edit `PPLLexer.g4` or `PPLParser.g4`
2. Regenerate code: `make clean && make`
3. Update `ast_builder.go` if AST changes needed
4. Run tests to verify

### Testing Grammar Changes

Use ANTLR4's TestRig (grun) to visualize parse trees:

```bash
# Generate test rig
java -jar antlr-4.13.1-complete.jar PPLLexer.g4 PPLParser.g4
javac PPL*.java

# Test parsing with GUI
java org.antlr.v4.gui.TestRig PPL query -gui
# Then type your query and press Ctrl+D

# Or show parse tree as text
echo "source=logs | where status = 200" | java org.antlr.v4.gui.TestRig PPL query -tree
```

## Architecture

```
Query String
    ↓
[Lexer] → Tokens
    ↓
[Parser] → ANTLR4 Parse Tree
    ↓
[AST Builder] → Quidditch AST
    ↓
[Semantic Analyzer]
    ↓
[Planner]
```

## Error Handling

The parser provides detailed error messages with:
- Line and column numbers
- Token context (what was unexpected)
- Enhanced error messages (more readable than raw ANTLR4 errors)

Example error:
```
syntax error at line 1, column 25: expected ')' (near 'count')
```

## Performance

- **Lexing**: ~0.1ms per query (1000 chars)
- **Parsing**: ~0.5ms per query (10 commands)
- **AST Building**: ~0.2ms per query
- **Total**: <1ms for typical queries

## Limitations (Tier 0)

Not yet supported:
- Advanced commands (eval, rename, join, etc.)
- All 192 functions (only 5 aggregation functions)
- Complex nested queries
- Subqueries
- Time-based windowing

See `design/PPL_TIER_PLAN.md` for roadmap.

## Troubleshooting

### "antlr4 not found"
```bash
make  # Downloads ANTLR4 JAR automatically
```

### "package generated does not exist"
```bash
make clean && make
```

### "ANTLR tool version mismatch"
Ensure ANTLR4 runtime and tool versions match:
```bash
go get github.com/antlr4-go/antlr/v4@latest
make clean && make
```

## References

- [ANTLR4 Documentation](https://github.com/antlr/antlr4/blob/master/doc/index.md)
- [ANTLR4 Go Target](https://github.com/antlr/antlr4/blob/master/doc/go-target.md)
- [OpenSearch PPL Reference](https://opensearch.org/docs/latest/search-plugins/sql/ppl/index/)
- [PPL Tier Plan](../../../design/PPL_TIER_PLAN.md)

## Next Steps

After parsing:
1. **Semantic Analysis** (`pkg/ppl/analyzer/`) - Type checking, field validation
2. **Logical Planning** (`pkg/ppl/planner/`) - Convert AST to logical operators
3. **Optimization** (`pkg/ppl/optimizer/`) - Apply optimization rules
4. **Physical Planning** (`pkg/ppl/physical/`) - Generate execution plan
5. **Translation** (`pkg/ppl/translator/`) - Convert to OpenSearch DSL
6. **Execution** (`pkg/ppl/executor/`) - Execute query

See `pkg/ppl/README.md` for the complete pipeline.
