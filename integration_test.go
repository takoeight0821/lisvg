package main

import (
	"os"
	"strings"
	"testing"
)

// TestIntegrationBasicPipeline tests the complete compilation pipeline
func TestIntegrationBasicPipeline(t *testing.T) {
	input := `(diagram
		(size 800 400)
		(nodes
			(id "A" :label "Start")
			(id "B" :label "End"))
		(edges
			("A" "B" :label "process")))`

	// Parse
	lexer := NewLexer(input)
	parser := NewParser(lexer)
	diagram, err := parser.ParseDiagram()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Validate
	validator := NewValidator()
	err = validator.Validate(diagram)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	// Layout
	layouter := NewSimpleLayouter()
	layout, err := layouter.LayoutDiagram(diagram)
	if err != nil {
		t.Fatalf("Layout failed: %v", err)
	}

	// Generate SVG
	svgGenerator := NewSVGGenerator()
	svg := svgGenerator.GenerateWithCustomStyles(layout, diagram)

	// Verify final output
	if !strings.Contains(svg, `<svg`) {
		t.Errorf("Final output should contain SVG")
	}
	if !strings.Contains(svg, "Start") {
		t.Errorf("Final output should contain node labels")
	}
	if !strings.Contains(svg, "process") {
		t.Errorf("Final output should contain edge labels")
	}
}

// TestIntegrationSampleFile tests with the actual sample file
func TestIntegrationSampleFile(t *testing.T) {
	// Read sample file
	content, err := os.ReadFile("sample.sxd")
	if err != nil {
		t.Skipf("sample.sxd not found, skipping integration test: %v", err)
	}

	// Complete pipeline
	lexer := NewLexer(string(content))
	parser := NewParser(lexer)
	diagram, err := parser.ParseDiagram()
	if err != nil {
		t.Fatalf("Failed to parse sample.sxd: %v", err)
	}

	validator := NewValidator()
	err = validator.Validate(diagram)
	if err != nil {
		t.Fatalf("Sample diagram validation failed: %v", err)
	}

	layouter := NewSimpleLayouter()
	layout, err := layouter.LayoutDiagram(diagram)
	if err != nil {
		t.Fatalf("Sample diagram layout failed: %v", err)
	}

	svgGenerator := NewSVGGenerator()
	svg := svgGenerator.GenerateWithCustomStyles(layout, diagram)

	// Check results
	if len(svg) == 0 {
		t.Errorf("Sample diagram should generate non-empty SVG")
	}

	// Check for Japanese characters
	if !strings.Contains(svg, "開始") || !strings.Contains(svg, "処理") || !strings.Contains(svg, "終了") {
		t.Errorf("Sample diagram should contain Japanese labels")
	}
}

// TestIntegrationComplexDiagram tests with a complex diagram
func TestIntegrationComplexDiagram(t *testing.T) {
	input := `(diagram
		(size 1000 600)
		(node-style :shape "rect" :fill "#f0f0f0")
		(edge-style :stroke "#333" :stroke-width "2")
		(nodes
			(id "start" :label "Start Process" :shape "ellipse")
			(id "input" :label "Get Input" :shape "rect")
			(id "validate" :label "Validate Data" :shape "diamond")
			(id "process1" :label "Process A" :shape "rect")
			(id "process2" :label "Process B" :shape "rect")
			(id "merge" :label "Merge Results" :shape "rect")
			(id "output" :label "Generate Output" :shape "rect")
			(id "end" :label "End Process" :shape "ellipse"))
		(edges
			("start" "input" :label "begin")
			("input" "validate" :label "data")
			("validate" "process1" :label "valid")
			("validate" "process2" :label "also valid")
			("process1" "merge" :label "result A")
			("process2" "merge" :label "result B")
			("merge" "output" :label "combined")
			("output" "end" :label "complete")))`

	// Complete pipeline
	lexer := NewLexer(input)
	parser := NewParser(lexer)
	diagram, err := parser.ParseDiagram()
	if err != nil {
		t.Fatalf("Complex diagram parse failed: %v", err)
	}

	if len(diagram.Nodes) != 8 {
		t.Errorf("Expected 8 nodes, got %d", len(diagram.Nodes))
	}
	if len(diagram.Edges) != 8 {
		t.Errorf("Expected 8 edges, got %d", len(diagram.Edges))
	}

	validator := NewValidator()
	err = validator.Validate(diagram)
	if err != nil {
		t.Fatalf("Complex diagram validation failed: %v", err)
	}

	layouter := NewSimpleLayouter()
	layout, err := layouter.LayoutDiagram(diagram)
	if err != nil {
		t.Fatalf("Complex diagram layout failed: %v", err)
	}

	// Check layout results
	if len(layout.Nodes) != 8 {
		t.Errorf("Layout should contain 8 nodes, got %d", len(layout.Nodes))
	}
	if len(layout.Edges) != 8 {
		t.Errorf("Layout should contain 8 edges, got %d", len(layout.Edges))
	}

	svgGenerator := NewSVGGenerator()
	svg := svgGenerator.GenerateWithCustomStyles(layout, diagram)

	// Verify complex SVG output
	if !strings.Contains(svg, `<ellipse`) {
		t.Errorf("Complex diagram should contain ellipses")
	}
	if !strings.Contains(svg, `<rect`) {
		t.Errorf("Complex diagram should contain rectangles")
	}
	if !strings.Contains(svg, `<polygon`) {
		t.Errorf("Complex diagram should contain polygons (diamonds)")
	}

	// Check for custom styles
	if !strings.Contains(svg, "#f0f0f0") {
		t.Errorf("Complex diagram should contain custom fill color")
	}
	if !strings.Contains(svg, "#333") {
		t.Errorf("Complex diagram should contain custom stroke color")
	}
}

// TestIntegrationErrorHandling tests error handling throughout pipeline
func TestIntegrationErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectStage string // Which stage should fail
	}{
		{
			name:        "parse error",
			input:       "(invalid syntax",
			expectStage: "parse",
		},
		{
			name: "validation error",
			input: `(diagram
				(nodes
					(id "A" :label "Node A")
					(id "A" :label "Duplicate"))
				(edges
					("A" "B")))`, // B doesn't exist
			expectStage: "validate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			parser := NewParser(lexer)
			diagram, err := parser.ParseDiagram()

			if tt.expectStage == "parse" {
				if err == nil {
					t.Errorf("Expected parse error, but parsing succeeded")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected parse error: %v", err)
			}

			validator := NewValidator()
			err = validator.Validate(diagram)

			if tt.expectStage == "validate" {
				if err == nil {
					t.Errorf("Expected validation error, but validation succeeded")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected validation error: %v", err)
			}
		})
	}
}

// TestIntegrationEmptyDiagram tests empty diagram handling
func TestIntegrationEmptyDiagram(t *testing.T) {
	input := `(diagram
		(size 400 200))`

	// Complete pipeline
	lexer := NewLexer(input)
	parser := NewParser(lexer)
	diagram, err := parser.ParseDiagram()
	if err != nil {
		t.Fatalf("Empty diagram parse failed: %v", err)
	}

	validator := NewValidator()
	err = validator.Validate(diagram)
	if err != nil {
		t.Fatalf("Empty diagram validation failed: %v", err)
	}

	layouter := NewSimpleLayouter()
	layout, err := layouter.LayoutDiagram(diagram)
	if err != nil {
		t.Fatalf("Empty diagram layout failed: %v", err)
	}

	if layout.Width != 400 || layout.Height != 200 {
		t.Errorf("Expected size 400x200, got %fx%f", layout.Width, layout.Height)
	}

	svgGenerator := NewSVGGenerator()
	svg := svgGenerator.Generate(layout, diagram)

	// Should still generate valid SVG
	if !strings.Contains(svg, `<svg`) {
		t.Errorf("Empty diagram should generate valid SVG")
	}
	if !strings.Contains(svg, `viewBox="0 0 440 240"`) {
		t.Errorf("Empty diagram should have correct viewBox")
	}
}

// TestIntegrationLargeDiagram tests performance with large diagrams
func TestIntegrationLargeDiagram(t *testing.T) {
	// Generate a large diagram programmatically
	var sb strings.Builder
	sb.WriteString("(diagram (size 2000 1000) (nodes ")

	// Create 50 nodes
	for i := 0; i < 50; i++ {
		sb.WriteString("(id \"node")
		sb.WriteString(string(rune('0' + i%10)))
		sb.WriteString(string(rune('0' + i/10)))
		sb.WriteString("\" :label \"Node ")
		sb.WriteString(string(rune('0' + i%10)))
		sb.WriteString(string(rune('0' + i/10)))
		sb.WriteString("\") ")
	}

	sb.WriteString(") (edges ")

	// Create 49 edges in a chain
	for i := 0; i < 49; i++ {
		sb.WriteString("(\"node")
		sb.WriteString(string(rune('0' + i%10)))
		sb.WriteString(string(rune('0' + i/10)))
		sb.WriteString("\" \"node")
		sb.WriteString(string(rune('0' + (i+1)%10)))
		sb.WriteString(string(rune('0' + (i+1)/10)))
		sb.WriteString("\") ")
	}

	sb.WriteString("))")

	input := sb.String()

	// Complete pipeline
	lexer := NewLexer(input)
	parser := NewParser(lexer)
	diagram, err := parser.ParseDiagram()
	if err != nil {
		t.Fatalf("Large diagram parse failed: %v", err)
	}

	if len(diagram.Nodes) != 50 {
		t.Errorf("Expected 50 nodes, got %d", len(diagram.Nodes))
	}

	validator := NewValidator()
	err = validator.Validate(diagram)
	if err != nil {
		t.Fatalf("Large diagram validation failed: %v", err)
	}

	layouter := NewSimpleLayouter()
	layout, err := layouter.LayoutDiagram(diagram)
	if err != nil {
		t.Fatalf("Large diagram layout failed: %v", err)
	}

	svgGenerator := NewSVGGenerator()
	svg := svgGenerator.Generate(layout, diagram)

	// Should generate reasonable output
	if len(svg) < 1000 {
		t.Errorf("Large diagram should generate substantial SVG output")
	}

	// Check for correct number of elements  
	// Count actual shape elements (rect is default)
	rectCount := strings.Count(svg, `<rect`)
	if rectCount != 50 {
		t.Errorf("Expected 50 rect elements, got %d", rectCount)
	}

	edgeCount := strings.Count(svg, "<path")
	if edgeCount != 49 {
		t.Errorf("Expected 49 edge elements, got %d", edgeCount)
	}
}

// TestIntegrationUnicodeHandling tests Unicode/international character handling
func TestIntegrationUnicodeHandling(t *testing.T) {
	input := `(diagram
		(nodes
			(id "日本" :label "日本のノード")
			(id "中国" :label "中文节点")
			(id "한국" :label "한국 노드")
			(id "русский" :label "Русский узел")
			(id "العربية" :label "العقدة العربية"))
		(edges
			("日本" "中国" :label "アジア")
			("中国" "한국" :label "亚洲")
			("한국" "русский" :label "국제")
			("русский" "العربية" :label "мир")))`

	// Complete pipeline
	lexer := NewLexer(input)
	parser := NewParser(lexer)
	diagram, err := parser.ParseDiagram()
	if err != nil {
		t.Fatalf("Unicode diagram parse failed: %v", err)
	}

	validator := NewValidator()
	err = validator.Validate(diagram)
	if err != nil {
		t.Fatalf("Unicode diagram validation failed: %v", err)
	}

	layouter := NewSimpleLayouter()
	layout, err := layouter.LayoutDiagram(diagram)
	if err != nil {
		t.Fatalf("Unicode diagram layout failed: %v", err)
	}

	svgGenerator := NewSVGGenerator()
	svg := svgGenerator.Generate(layout, diagram)

	// Check for Unicode content
	unicodeStrings := []string{
		"日本のノード", "中文节点", "한국 노드", "Русский узел", "العقدة العربية",
		"アジア", "亚洲", "국제", "мир",
	}

	for _, str := range unicodeStrings {
		if !strings.Contains(svg, str) {
			t.Errorf("Unicode diagram should contain '%s'", str)
		}
	}
}

// TestIntegrationStyleInheritance tests style inheritance and overrides
func TestIntegrationStyleInheritance(t *testing.T) {
	input := `(diagram
		(node-style :shape "rect" :fill "#ffffff")
		(edge-style :stroke "#000000" :stroke-width "1")
		(nodes
			(id "default" :label "Default Node")
			(id "custom" :label "Custom Node" :shape "ellipse" :fill "#ff0000"))
		(edges
			("default" "custom" :label "connection" :stroke "#0000ff")))`

	// Complete pipeline
	lexer := NewLexer(input)
	parser := NewParser(lexer)
	diagram, err := parser.ParseDiagram()
	if err != nil {
		t.Fatalf("Style inheritance parse failed: %v", err)
	}

	// Check that styles are parsed correctly
	if diagram.NodeStyle["shape"] != "rect" {
		t.Errorf("Expected default node shape 'rect', got '%s'", diagram.NodeStyle["shape"])
	}

	// Check node-specific overrides
	customNode := diagram.Nodes[1] // Assuming order is preserved
	if customNode.Attributes["shape"] != "ellipse" {
		t.Errorf("Expected custom node shape 'ellipse', got '%s'", customNode.Attributes["shape"])
	}

	validator := NewValidator()
	err = validator.Validate(diagram)
	if err != nil {
		t.Fatalf("Style inheritance validation failed: %v", err)
	}

	layouter := NewSimpleLayouter()
	layout, err := layouter.LayoutDiagram(diagram)
	if err != nil {
		t.Fatalf("Style inheritance layout failed: %v", err)
	}

	// Check that layout preserves shapes
	defaultNode := layout.Nodes["default"]
	customNodeLayout := layout.Nodes["custom"]

	if defaultNode.Shape != "rect" {
		t.Errorf("Expected default node layout shape 'rect', got '%s'", defaultNode.Shape)
	}
	if customNodeLayout.Shape != "ellipse" {
		t.Errorf("Expected custom node layout shape 'ellipse', got '%s'", customNodeLayout.Shape)
	}

	svgGenerator := NewSVGGenerator()
	svg := svgGenerator.GenerateWithCustomStyles(layout, diagram)

	// Check that SVG contains appropriate elements
	if !strings.Contains(svg, "<rect") {
		t.Errorf("SVG should contain rect for default node")
	}
	if !strings.Contains(svg, "<ellipse") {
		t.Errorf("SVG should contain ellipse for custom node")
	}
}

// BenchmarkIntegrationFullPipeline benchmarks the complete pipeline
func BenchmarkIntegrationFullPipeline(b *testing.B) {
	input := `(diagram
		(size 800 400)
		(node-style :shape "rect")
		(edge-style :stroke "#555")
		(nodes
			(id "A" :label "Node A")
			(id "B" :label "Node B")
			(id "C" :label "Node C")
			(id "D" :label "Node D"))
		(edges
			("A" "B" :label "edge1")
			("B" "C" :label "edge2")
			("C" "D" :label "edge3")))`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Complete pipeline
		lexer := NewLexer(input)
		parser := NewParser(lexer)
		diagram, err := parser.ParseDiagram()
		if err != nil {
			b.Fatalf("Parse failed: %v", err)
		}

		validator := NewValidator()
		err = validator.Validate(diagram)
		if err != nil {
			b.Fatalf("Validation failed: %v", err)
		}

		layouter := NewSimpleLayouter()
		layout, err := layouter.LayoutDiagram(diagram)
		if err != nil {
			b.Fatalf("Layout failed: %v", err)
		}

		svgGenerator := NewSVGGenerator()
		_ = svgGenerator.GenerateWithCustomStyles(layout, diagram)
	}
}