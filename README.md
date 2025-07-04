# lisvg

A tool for generating SVG diagrams from S-expression descriptions. It uses a built-in layout engine and generates clean SVG output.

## Features

- S-expression based diagram description language
- Built-in simple tree layout algorithm
- Clean SVG output with embedded CSS
- Command-line interface with no external dependencies
- Validation of node IDs and edge references
- Support for custom node and edge styling

## Installation

### Prerequisites

- Go 1.19 or later

### Build

```bash
go build -o lisvg
```

## Usage

### Basic Usage

```bash
# Compile a diagram file
./lisvg compile sample.sxd

# Specify output file
./lisvg compile sample.sxd -o output.svg

# Verbose output
./lisvg compile sample.sxd -v

# Read from stdin, write to stdout
cat sample.sxd | ./lisvg compile
```

### S-expression Format

```lisp
(diagram
  (size 800 400)                 ; Canvas size (width height)
  (node-style :shape "rect")     ; Default node attributes
  (edge-style :stroke "#555")    ; Default edge attributes

  (nodes                         ; Node definitions
    (id "A" :label "Start")
    (id "B" :label "Process")
    (id "C" :label "End" :shape "ellipse"))

  (edges                         ; Edge definitions
    ("A" "B" :label "next")
    ("B" "C")))
```

### Supported Node Shapes

- `rect`, `rectangle`, `box`
- `ellipse`, `circle`, `oval`
- `diamond`, `rhombus`

### Supported Edge Styles

- `solid`, `dashed`, `dotted`, `bold`
- `invis`, `invisible`

## Examples

See `sample.sxd` for a complete example.

## Architecture

The implementation follows a clean pipeline:

1. **Lexer/Parser**: Converts S-expressions to AST
2. **Validator**: Checks node ID uniqueness and edge references
3. **Layout Engine**: Uses built-in algorithms for node positioning
4. **SVG Generator**: Creates final SVG output

### Layout Algorithm

The current implementation uses a simple hierarchical tree layout:
- Nodes are arranged in levels based on edge dependencies
- Root nodes (no incoming edges) are placed at the top
- Each level is spaced vertically with nodes centered horizontally
- Edges are drawn as straight lines between node centers

## License

MIT License