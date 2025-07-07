package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "lisvg",
		Short: "lisvg - S-expression to SVG diagram generator",
		Long: `lisvg is a tool for generating SVG diagrams from S-expression descriptions.
It uses a built-in layout engine and generates clean SVG output.`,
	}

	var compileCmd = &cobra.Command{
		Use:   "compile [input.sxd]",
		Short: "Compile S-expression diagram to SVG",
		Long: `Compile an S-expression diagram file (.sxd) to SVG format.
The input file should contain a valid S-expression diagram description.`,
		Args: cobra.MaximumNArgs(1),
		RunE: compileCommand,
	}

	var outputFile string
	var verbose bool

	compileCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output SVG file (default: replace .sxd with .svg)")
	compileCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	rootCmd.AddCommand(compileCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func compileCommand(cmd *cobra.Command, args []string) error {
	// Get input file
	var inputFile string
	if len(args) > 0 {
		inputFile = args[0]
	} else {
		// Read from stdin
		inputFile = "-"
	}

	// Get output file
	outputFile, _ := cmd.Flags().GetString("output")
	if outputFile == "" {
		if inputFile == "-" {
			outputFile = "-" // stdout
		} else {
			outputFile = replaceExtension(inputFile, ".svg")
		}
	}

	verbose, _ := cmd.Flags().GetBool("verbose")

	// Compile the diagram
	if err := compileDiagram(inputFile, outputFile, verbose); err != nil {
		return fmt.Errorf("compilation failed: %w", err)
	}

	if verbose {
		fmt.Printf("Successfully compiled %s to %s\n", inputFile, outputFile)
	}

	return nil
}

func compileDiagram(inputFile, outputFile string, verbose bool) error {
	// Read input
	var input []byte
	var err error

	if inputFile == "-" {
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read from stdin: %w", err)
		}
	} else {
		input, err = os.ReadFile(inputFile)
		if err != nil {
			return fmt.Errorf("failed to read input file: %w", err)
		}
	}

	if verbose {
		fmt.Printf("Parsing S-expression from %s...\n", inputFile)
	}

	// Parse S-expression
	lexer := NewLexer(string(input))
	parser := NewParser(lexer)
	diagram, err := parser.ParseDiagram()
	if err != nil {
		return fmt.Errorf("parsing failed: %w", err)
	}

	if verbose {
		fmt.Printf("Parsed diagram with %d nodes and %d edges\n", len(diagram.Nodes), len(diagram.Edges))
	}

	// Validate AST
	validator := NewValidator()
	if err := validator.Validate(diagram); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if verbose {
		fmt.Println("Validation passed")
	}

	// Generate layout
	layouter := NewSimpleLayouter()
	layout, err := layouter.LayoutDiagram(diagram)
	if err != nil {
		return fmt.Errorf("layout failed: %w", err)
	}

	if verbose {
		fmt.Printf("Layout completed: %.0fx%.0f\n", layout.Width, layout.Height)
	}

	// Generate SVG
	svgGenerator := NewSVGGenerator()
	svgContent := svgGenerator.GenerateWithCustomStyles(layout, diagram)

	if verbose {
		fmt.Println("Generated SVG content")
	}

	// Write output
	if outputFile == "-" {
		fmt.Print(svgContent)
	} else {
		if err := os.WriteFile(outputFile, []byte(svgContent), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	return nil
}

func replaceExtension(filename, newExt string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return filename + newExt
	}
	return strings.TrimSuffix(filename, ext) + newExt
}
