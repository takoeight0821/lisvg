package main

import (
	"strings"
	"testing"
)

// TestSVGGenerator tests basic SVG generator functionality
func TestSVGGenerator(t *testing.T) {
	generator := NewSVGGenerator()
	if generator == nil {
		t.Errorf("NewSVGGenerator returned nil")
	}
}

// TestSVGGenerationBasic tests basic SVG generation
func TestSVGGenerationBasic(t *testing.T) {
	generator := NewSVGGenerator()
	
	layout := &Layout{
		Width:  400,
		Height: 300,
		Nodes: map[string]LayoutNode{
			"A": {
				ID:     "A",
				X:      100,
				Y:      50,
				Width:  80,
				Height: 40,
				Label:  "Node A",
				Shape:  "rect",
			},
		},
		Edges: []LayoutEdge{},
	}
	
	diagram := &Diagram{
		NodeStyle: map[string]string{},
		EdgeStyle: map[string]string{},
	}
	
	svg := generator.Generate(layout, diagram)
	
	// Check basic SVG structure
	if !strings.Contains(svg, `<?xml version="1.0" encoding="UTF-8"?>`) {
		t.Errorf("SVG should contain XML declaration")
	}
	
	if !strings.Contains(svg, `<svg xmlns="http://www.w3.org/2000/svg"`) {
		t.Errorf("SVG should contain proper namespace")
	}
	
	if !strings.Contains(svg, `</svg>`) {
		t.Errorf("SVG should be properly closed")
	}
	
	// Check that the node is rendered
	if !strings.Contains(svg, `<rect`) {
		t.Errorf("SVG should contain rect element for node")
	}
	
	if !strings.Contains(svg, "Node A") {
		t.Errorf("SVG should contain node label")
	}
}

// TestSVGNodeShapes tests different node shape rendering
func TestSVGNodeShapes(t *testing.T) {
	generator := NewSVGGenerator()
	
	tests := []struct {
		shape    string
		expected string
	}{
		{"rect", "<rect"},
		{"ellipse", "<ellipse"},
		{"diamond", "<polygon"},
		{"circle", "<ellipse"},
	}
	
	for _, tt := range tests {
		t.Run(tt.shape, func(t *testing.T) {
			layout := &Layout{
				Width:  200,
				Height: 200,
				Nodes: map[string]LayoutNode{
					"test": {
						ID:     "test",
						X:      100,
						Y:      100,
						Width:  80,
						Height: 40,
						Label:  "Test",
						Shape:  tt.shape,
					},
				},
				Edges: []LayoutEdge{},
			}
			
			diagram := &Diagram{
				NodeStyle: map[string]string{},
				EdgeStyle: map[string]string{},
			}
			
			svg := generator.Generate(layout, diagram)
			
			if !strings.Contains(svg, tt.expected) {
				t.Errorf("Expected SVG to contain '%s' for shape '%s'", tt.expected, tt.shape)
			}
		})
	}
}

// TestSVGEdgeGeneration tests edge rendering
func TestSVGEdgeGeneration(t *testing.T) {
	generator := NewSVGGenerator()
	
	layout := &Layout{
		Width:  400,
		Height: 300,
		Nodes: map[string]LayoutNode{
			"A": {ID: "A", X: 100, Y: 50, Width: 80, Height: 40, Label: "A", Shape: "rect"},
			"B": {ID: "B", X: 200, Y: 150, Width: 80, Height: 40, Label: "B", Shape: "rect"},
		},
		Edges: []LayoutEdge{
			{
				From: "A",
				To:   "B",
				Points: []Point{
					{X: 100, Y: 70},
					{X: 200, Y: 130},
				},
				Label: "connection",
				X:     150,
				Y:     100,
			},
		},
	}
	
	diagram := &Diagram{
		NodeStyle: map[string]string{},
		EdgeStyle: map[string]string{},
	}
	
	svg := generator.Generate(layout, diagram)
	
	// Check for path element
	if !strings.Contains(svg, `<path`) {
		t.Errorf("SVG should contain path element for edge")
	}
	
	// Check for arrowhead marker
	if !strings.Contains(svg, `marker-end="url(#arrowhead)"`) {
		t.Errorf("SVG should contain arrowhead marker")
	}
	
	// Check for edge label
	if !strings.Contains(svg, "connection") {
		t.Errorf("SVG should contain edge label")
	}
}

// TestSVGCoordinateTransform tests coordinate transformation
func TestSVGCoordinateTransform(t *testing.T) {
	generator := NewSVGGenerator()
	
	layout := &Layout{
		Width:  200,
		Height: 200,
		Nodes: map[string]LayoutNode{
			"test": {
				ID:     "test",
				X:      100,
				Y:      100,
				Width:  80,
				Height: 40,
				Label:  "Test",
				Shape:  "rect",
			},
		},
		Edges: []LayoutEdge{},
	}
	
	diagram := &Diagram{
		NodeStyle: map[string]string{},
		EdgeStyle: map[string]string{},
	}
	
	svg := generator.Generate(layout, diagram)
	
	// Check for coordinate transformation
	if !strings.Contains(svg, `transform="translate(20, 220) scale(1, -1)"`) {
		t.Errorf("SVG should contain coordinate transformation")
	}
	
	// Check for text transformation (to flip text back)
	if !strings.Contains(svg, `transform="scale(1, -1)"`) {
		t.Errorf("SVG should contain text transformation")
	}
}

// TestSVGCSS tests CSS generation
func TestSVGCSS(t *testing.T) {
	generator := NewSVGGenerator()
	
	layout := &Layout{
		Width:  200,
		Height: 200,
		Nodes:  map[string]LayoutNode{},
		Edges:  []LayoutEdge{},
	}
	
	diagram := &Diagram{
		NodeStyle: map[string]string{},
		EdgeStyle: map[string]string{},
	}
	
	svg := generator.Generate(layout, diagram)
	
	// Check for CSS style block
	if !strings.Contains(svg, `<style>`) {
		t.Errorf("SVG should contain style block")
	}
	
	if !strings.Contains(svg, `.node {`) {
		t.Errorf("SVG should contain node styles")
	}
	
	if !strings.Contains(svg, `.edge {`) {
		t.Errorf("SVG should contain edge styles")
	}
	
	if !strings.Contains(svg, `.node-label {`) {
		t.Errorf("SVG should contain node label styles")
	}
	
	if !strings.Contains(svg, `.edge-label {`) {
		t.Errorf("SVG should contain edge label styles")
	}
}

// TestSVGCustomStyles tests custom style application
func TestSVGCustomStyles(t *testing.T) {
	generator := NewSVGGenerator()
	
	layout := &Layout{
		Width:  200,
		Height: 200,
		Nodes:  map[string]LayoutNode{},
		Edges:  []LayoutEdge{},
	}
	
	diagram := &Diagram{
		NodeStyle: map[string]string{
			"fill":   "#ff0000",
			"stroke": "#0000ff",
		},
		EdgeStyle: map[string]string{
			"stroke":       "#00ff00",
			"stroke-width": "2",
		},
	}
	
	svg := generator.GenerateWithCustomStyles(layout, diagram)
	
	// Check for custom node styles
	if !strings.Contains(svg, "fill: #ff0000") {
		t.Errorf("SVG should contain custom node fill color")
	}
	
	if !strings.Contains(svg, "stroke: #0000ff") {
		t.Errorf("SVG should contain custom node stroke color")
	}
	
	// Check for custom edge styles
	if !strings.Contains(svg, "stroke: #00ff00") {
		t.Errorf("SVG should contain custom edge stroke color")
	}
	
	if !strings.Contains(svg, "stroke-width: 2") {
		t.Errorf("SVG should contain custom edge stroke width")
	}
}

// TestSVGViewBox tests viewBox calculation
func TestSVGViewBox(t *testing.T) {
	generator := NewSVGGenerator()
	
	layout := &Layout{
		Width:  400,
		Height: 300,
		Nodes:  map[string]LayoutNode{},
		Edges:  []LayoutEdge{},
	}
	
	diagram := &Diagram{
		NodeStyle: map[string]string{},
		EdgeStyle: map[string]string{},
	}
	
	svg := generator.Generate(layout, diagram)
	
	// Check viewBox with padding (400+40 = 440, 300+40 = 340)
	if !strings.Contains(svg, `viewBox="0 0 440 340"`) {
		t.Errorf("SVG should contain correct viewBox")
	}
	
	// Check width and height
	if !strings.Contains(svg, `width="440" height="340"`) {
		t.Errorf("SVG should contain correct width and height")
	}
}

// TestSVGXMLEscaping tests XML character escaping
func TestSVGXMLEscaping(t *testing.T) {
	generator := NewSVGGenerator()
	
	layout := &Layout{
		Width:  200,
		Height: 200,
		Nodes: map[string]LayoutNode{
			"test": {
				ID:     "test",
				X:      100,
				Y:      100,
				Width:  80,
				Height: 40,
				Label:  `<tag>"quotes"&ampersand`,
				Shape:  "rect",
			},
		},
		Edges: []LayoutEdge{},
	}
	
	diagram := &Diagram{
		NodeStyle: map[string]string{},
		EdgeStyle: map[string]string{},
	}
	
	svg := generator.Generate(layout, diagram)
	
	// Check for escaped characters
	if !strings.Contains(svg, "&lt;tag&gt;") {
		t.Errorf("SVG should escape < and >")
	}
	
	if !strings.Contains(svg, "&quot;quotes&quot;") {
		t.Errorf("SVG should escape quotes")
	}
	
	if !strings.Contains(svg, "&amp;ampersand") {
		t.Errorf("SVG should escape ampersands")
	}
}

// TestSVGMarkerDefinitions tests SVG marker definitions
func TestSVGMarkerDefinitions(t *testing.T) {
	generator := NewSVGGenerator()
	
	layout := &Layout{
		Width:  200,
		Height: 200,
		Nodes:  map[string]LayoutNode{},
		Edges:  []LayoutEdge{},
	}
	
	diagram := &Diagram{
		NodeStyle: map[string]string{},
		EdgeStyle: map[string]string{},
	}
	
	svg := generator.Generate(layout, diagram)
	
	// Check for defs section
	if !strings.Contains(svg, `<defs>`) {
		t.Errorf("SVG should contain defs section")
	}
	
	// Check for arrowhead marker definition
	if !strings.Contains(svg, `<marker id="arrowhead"`) {
		t.Errorf("SVG should contain arrowhead marker definition")
	}
	
	if !strings.Contains(svg, `<polygon points="0 0, 10 3.5, 0 7"`) {
		t.Errorf("SVG should contain arrowhead polygon")
	}
}

// TestSVGEmptyLayout tests SVG generation with empty layout
func TestSVGEmptyLayout(t *testing.T) {
	generator := NewSVGGenerator()
	
	layout := &Layout{
		Width:  400,
		Height: 300,
		Nodes:  map[string]LayoutNode{},
		Edges:  []LayoutEdge{},
	}
	
	diagram := &Diagram{
		NodeStyle: map[string]string{},
		EdgeStyle: map[string]string{},
	}
	
	svg := generator.Generate(layout, diagram)
	
	// Should still generate valid SVG structure
	if !strings.Contains(svg, `<svg`) {
		t.Errorf("Empty layout should still generate SVG")
	}
	
	if !strings.Contains(svg, `</svg>`) {
		t.Errorf("Empty layout should still close SVG")
	}
	
	// Should not contain any nodes or edges
	if strings.Contains(svg, `<rect`) || strings.Contains(svg, `<ellipse`) {
		t.Errorf("Empty layout should not contain node elements")
	}
	
	if strings.Contains(svg, `<path`) {
		t.Errorf("Empty layout should not contain edge elements")
	}
}

// TestSVGComplexDiagram tests SVG generation with complex diagram
func TestSVGComplexDiagram(t *testing.T) {
	generator := NewSVGGenerator()
	
	layout := &Layout{
		Width:  600,
		Height: 400,
		Nodes: map[string]LayoutNode{
			"start": {ID: "start", X: 100, Y: 50, Width: 80, Height: 40, Label: "Start", Shape: "ellipse"},
			"proc1": {ID: "proc1", X: 100, Y: 150, Width: 80, Height: 40, Label: "Process 1", Shape: "rect"},
			"proc2": {ID: "proc2", X: 250, Y: 150, Width: 80, Height: 40, Label: "Process 2", Shape: "rect"},
			"end":   {ID: "end", X: 175, Y: 250, Width: 80, Height: 40, Label: "End", Shape: "diamond"},
		},
		Edges: []LayoutEdge{
			{From: "start", To: "proc1", Points: []Point{{100, 70}, {100, 130}}, Label: "begin", X: 100, Y: 100},
			{From: "start", To: "proc2", Points: []Point{{100, 70}, {250, 130}}, Label: "alt", X: 175, Y: 100},
			{From: "proc1", To: "end", Points: []Point{{100, 170}, {175, 230}}, X: 137, Y: 200},
			{From: "proc2", To: "end", Points: []Point{{250, 170}, {175, 230}}, X: 212, Y: 200},
		},
	}
	
	diagram := &Diagram{
		NodeStyle: map[string]string{"fill": "#f0f0f0"},
		EdgeStyle: map[string]string{"stroke": "#333"},
	}
	
	svg := generator.GenerateWithCustomStyles(layout, diagram)
	
	// Check that all node types are present
	if !strings.Contains(svg, `<ellipse`) {
		t.Errorf("Complex diagram should contain ellipse")
	}
	
	if !strings.Contains(svg, `<rect`) {
		t.Errorf("Complex diagram should contain rectangles")
	}
	
	if !strings.Contains(svg, `<polygon`) {
		t.Errorf("Complex diagram should contain polygon (diamond)")
	}
	
	// Check that all edges are present
	edgeCount := strings.Count(svg, `<path`)
	if edgeCount != 4 {
		t.Errorf("Expected 4 edges, found %d", edgeCount)
	}
	
	// Check that labels are present
	if !strings.Contains(svg, "Start") || !strings.Contains(svg, "Process 1") {
		t.Errorf("Complex diagram should contain all node labels")
	}
	
	if !strings.Contains(svg, "begin") || !strings.Contains(svg, "alt") {
		t.Errorf("Complex diagram should contain edge labels")
	}
}

// TestGetNodeShape tests node shape mapping
func TestGetNodeShape(t *testing.T) {
	generator := NewSVGGenerator()
	
	tests := []struct {
		input    string
		expected string
	}{
		{"rect", "rect"},
		{"rectangle", "rect"},
		{"box", "rect"},
		{"ellipse", "ellipse"},
		{"oval", "ellipse"},
		{"circle", "ellipse"},
		{"diamond", "diamond"},
		{"rhombus", "diamond"},
		{"unknown", "ellipse"}, // default
		{"", "ellipse"},        // default
	}
	
	for _, tt := range tests {
		result := generator.getNodeShape(tt.input)
		if result != tt.expected {
			t.Errorf("getNodeShape(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

// TestEscapeXML tests XML escaping function
func TestEscapeXML(t *testing.T) {
	generator := NewSVGGenerator()
	
	tests := []struct {
		input    string
		expected string
	}{
		{"normal text", "normal text"},
		{"<tag>", "&lt;tag&gt;"},
		{"\"quotes\"", "&quot;quotes&quot;"},
		{"'apostrophe'", "&apos;apostrophe&apos;"},
		{"&ampersand", "&amp;ampersand"},
		{"<>&\"'", "&lt;&gt;&amp;&quot;&apos;"},
		{"", ""},
	}
	
	for _, tt := range tests {
		result := generator.escapeXML(tt.input)
		if result != tt.expected {
			t.Errorf("escapeXML(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

// BenchmarkSVGGeneration benchmarks SVG generation performance
func BenchmarkSVGGeneration(b *testing.B) {
	generator := NewSVGGenerator()
	
	// Create a moderate-sized layout
	nodes := make(map[string]LayoutNode)
	edges := make([]LayoutEdge, 19)
	
	for i := 0; i < 20; i++ {
		id := string(rune('A' + i))
		nodes[id] = LayoutNode{
			ID:     id,
			X:      float64(i * 50),
			Y:      float64((i % 4) * 100),
			Width:  80,
			Height: 40,
			Label:  "Node " + id,
			Shape:  "rect",
		}
	}
	
	for i := 0; i < 19; i++ {
		from := string(rune('A' + i))
		to := string(rune('A' + i + 1))
		edges[i] = LayoutEdge{
			From: from,
			To:   to,
			Points: []Point{
				{float64(i * 50), float64((i % 4) * 100)},
				{float64((i + 1) * 50), float64(((i + 1) % 4) * 100)},
			},
		}
	}
	
	layout := &Layout{
		Width:  1000,
		Height: 400,
		Nodes:  nodes,
		Edges:  edges,
	}
	
	diagram := &Diagram{
		NodeStyle: map[string]string{"fill": "#fff"},
		EdgeStyle: map[string]string{"stroke": "#000"},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = generator.GenerateWithCustomStyles(layout, diagram)
	}
}