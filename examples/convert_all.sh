#!/bin/bash

# Convert all .sxd files in the examples directory to SVG

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if lisvg binary exists
if [ ! -f "../lisvg" ]; then
    echo -e "${RED}Error: lisvg binary not found. Please build it first with 'go build -o lisvg'${NC}"
    exit 1
fi

echo "Converting .sxd files to SVG..."
echo "=============================="

# Convert each .sxd file
for sxd_file in *.sxd; do
    if [ -f "$sxd_file" ]; then
        # Get filename without extension
        base_name="${sxd_file%.sxd}"
        svg_file="${base_name}.svg"
        
        echo -n "Converting $sxd_file -> $svg_file ... "
        
        # Run the conversion
        if ../lisvg compile "$sxd_file" -o "$svg_file" 2>/dev/null; then
            echo -e "${GREEN}✓${NC}"
        else
            echo -e "${RED}✗${NC}"
            echo "  Error converting $sxd_file"
        fi
    fi
done

echo "=============================="
echo "Conversion complete!"

# List generated SVG files
echo ""
echo "Generated SVG files:"
ls -la *.svg 2>/dev/null || echo "No SVG files generated"