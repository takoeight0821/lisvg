package main

import (
	"fmt"
	"strings"
	"testing"
)

// TestHelpers contains utility functions for testing

// CreateTestDiagram creates a simple test diagram for use in tests
func CreateTestDiagram() *Diagram {
	return &Diagram{
		Width:  800,
		Height: 400,
		NodeStyle: map[string]string{
			"shape": "rect",
			"fill":  "#ffffff",
		},
		EdgeStyle: map[string]string{
			"stroke": "#000000",
		},
		Nodes: []Node{
			{
				ID:    "A",
				Label: "Node A",
				Attributes: map[string]string{
					"shape": "rect",
				},
			},
			{
				ID:    "B",
				Label: "Node B",
				Attributes: map[string]string{
					"shape": "ellipse",
				},
			},
		},
		Edges: []Edge{
			{
				From:  "A",
				To:    "B",
				Label: "connects",
				Attributes: map[string]string{
					"style": "solid",
				},
			},
		},
	}
}

// CreateLinearDiagram creates a linear chain diagram for testing
func CreateLinearDiagram(nodeCount int) *Diagram {
	nodes := make([]Node, nodeCount)
	edges := make([]Edge, nodeCount-1)

	for i := 0; i < nodeCount; i++ {
		nodes[i] = Node{
			ID:         fmt.Sprintf("node%d", i),
			Label:      fmt.Sprintf("Node %d", i),
			Attributes: map[string]string{"shape": "rect"},
		}
	}

	for i := 0; i < nodeCount-1; i++ {
		edges[i] = Edge{
			From:       fmt.Sprintf("node%d", i),
			To:         fmt.Sprintf("node%d", i+1),
			Label:      fmt.Sprintf("edge%d", i),
			Attributes: map[string]string{"style": "solid"},
		}
	}

	return &Diagram{
		Width:     800,
		Height:    400,
		NodeStyle: map[string]string{"shape": "rect"},
		EdgeStyle: map[string]string{"stroke": "#000"},
		Nodes:     nodes,
		Edges:     edges,
	}
}

// CreateBranchingDiagram creates a branching diagram for testing
func CreateBranchingDiagram() *Diagram {
	return &Diagram{
		Width:  1000,
		Height: 600,
		NodeStyle: map[string]string{
			"shape": "rect",
		},
		EdgeStyle: map[string]string{
			"stroke": "#333",
		},
		Nodes: []Node{
			{ID: "root", Label: "Root Node", Attributes: map[string]string{"shape": "ellipse"}},
			{ID: "left", Label: "Left Branch", Attributes: map[string]string{"shape": "rect"}},
			{ID: "right", Label: "Right Branch", Attributes: map[string]string{"shape": "rect"}},
			{ID: "left_child", Label: "Left Child", Attributes: map[string]string{"shape": "diamond"}},
			{ID: "right_child", Label: "Right Child", Attributes: map[string]string{"shape": "diamond"}},
			{ID: "merge", Label: "Merge Point", Attributes: map[string]string{"shape": "ellipse"}},
		},
		Edges: []Edge{
			{From: "root", To: "left", Label: "left path"},
			{From: "root", To: "right", Label: "right path"},
			{From: "left", To: "left_child", Label: "left child"},
			{From: "right", To: "right_child", Label: "right child"},
			{From: "left_child", To: "merge", Label: "merge left"},
			{From: "right_child", To: "merge", Label: "merge right"},
		},
	}
}

// CreateCyclicDiagram creates a diagram with cycles for testing
func CreateCyclicDiagram() *Diagram {
	return &Diagram{
		Width:     600,
		Height:    400,
		NodeStyle: map[string]string{"shape": "rect"},
		EdgeStyle: map[string]string{"stroke": "#000"},
		Nodes: []Node{
			{ID: "A", Label: "Node A", Attributes: map[string]string{}},
			{ID: "B", Label: "Node B", Attributes: map[string]string{}},
			{ID: "C", Label: "Node C", Attributes: map[string]string{}},
		},
		Edges: []Edge{
			{From: "A", To: "B", Label: "A to B"},
			{From: "B", To: "C", Label: "B to C"},
			{From: "C", To: "A", Label: "C to A"}, // Creates cycle
		},
	}
}

// CreateInvalidDiagram creates a diagram with validation errors for testing
func CreateInvalidDiagram() *Diagram {
	return &Diagram{
		Width:     400,
		Height:    300,
		NodeStyle: map[string]string{},
		EdgeStyle: map[string]string{},
		Nodes: []Node{
			{ID: "A", Label: "Node A", Attributes: map[string]string{}},
			{ID: "A", Label: "Duplicate A", Attributes: map[string]string{}},                    // Duplicate ID
			{ID: "", Label: "Empty ID", Attributes: map[string]string{}},                        // Empty ID
			{ID: "B", Label: "Node B", Attributes: map[string]string{"shape": "invalid-shape"}}, // Invalid shape
		},
		Edges: []Edge{
			{From: "A", To: "C", Label: "A to C"},    // C doesn't exist
			{From: "", To: "B", Label: "Empty from"}, // Empty from
			{From: "B", To: "", Label: "Empty to"},   // Empty to
			{From: "X", To: "Y", Label: "Both invalid", Attributes: map[string]string{"style": "invalid-style"}}, // Both nodes don't exist, invalid style
		},
	}
}

// ParseTestInput parses S-expression input for testing
func ParseTestInput(t *testing.T, input string) *Diagram {
	t.Helper()

	lexer := NewLexer(input)
	parser := NewParser(lexer)
	diagram, err := parser.ParseDiagram()
	if err != nil {
		t.Fatalf("Failed to parse test input: %v", err)
	}
	return diagram
}

// ValidateTestDiagram validates a diagram and handles test failures
func ValidateTestDiagram(t *testing.T, diagram *Diagram) {
	t.Helper()

	validator := NewValidator()
	err := validator.Validate(diagram)
	if err != nil {
		t.Fatalf("Diagram validation failed: %v", err)
	}
}

// LayoutTestDiagram creates layout for a diagram and handles test failures
func LayoutTestDiagram(t *testing.T, diagram *Diagram) *Layout {
	t.Helper()

	layouter := NewSimpleLayouter()
	layout, err := layouter.LayoutDiagram(diagram)
	if err != nil {
		t.Fatalf("Diagram layout failed: %v", err)
	}
	return layout
}

// GenerateTestSVG generates SVG for a layout and handles test failures
func GenerateTestSVG(t *testing.T, layout *Layout, diagram *Diagram) string {
	t.Helper()

	generator := NewSVGGenerator()
	return generator.GenerateWithCustomStyles(layout, diagram)
}

// CompletePipeline runs the complete compilation pipeline for testing
func CompletePipeline(t *testing.T, input string) string {
	t.Helper()

	// Parse
	diagram := ParseTestInput(t, input)

	// Validate
	ValidateTestDiagram(t, diagram)

	// Layout
	layout := LayoutTestDiagram(t, diagram)

	// Generate SVG
	return GenerateTestSVG(t, layout, diagram)
}

// AssertSVGContains checks that SVG contains expected content
func AssertSVGContains(t *testing.T, svg string, expected ...string) {
	t.Helper()

	for _, exp := range expected {
		if !strings.Contains(svg, exp) {
			t.Errorf("SVG should contain '%s'", exp)
		}
	}
}

// AssertSVGNotContains checks that SVG does not contain specified content
func AssertSVGNotContains(t *testing.T, svg string, notExpected ...string) {
	t.Helper()

	for _, notExp := range notExpected {
		if strings.Contains(svg, notExp) {
			t.Errorf("SVG should not contain '%s'", notExp)
		}
	}
}

// AssertValidSVGStructure checks basic SVG structure
func AssertValidSVGStructure(t *testing.T, svg string) {
	t.Helper()

	requiredElements := []string{
		`<?xml version="1.0" encoding="UTF-8"?>`,
		`<svg xmlns="http://www.w3.org/2000/svg"`,
		`<style>`,
		`<defs>`,
		`</svg>`,
	}

	AssertSVGContains(t, svg, requiredElements...)
}

// AssertNodeCount checks that layout has expected number of nodes
func AssertNodeCount(t *testing.T, layout *Layout, expected int) {
	t.Helper()

	if len(layout.Nodes) != expected {
		t.Errorf("Expected %d nodes, got %d", expected, len(layout.Nodes))
	}
}

// AssertEdgeCount checks that layout has expected number of edges
func AssertEdgeCount(t *testing.T, layout *Layout, expected int) {
	t.Helper()

	if len(layout.Edges) != expected {
		t.Errorf("Expected %d edges, got %d", expected, len(layout.Edges))
	}
}

// AssertValidCoordinates checks that all nodes have valid coordinates
func AssertValidCoordinates(t *testing.T, layout *Layout) {
	t.Helper()

	for id, node := range layout.Nodes {
		if node.X < 0 || node.Y < 0 {
			t.Errorf("Node %s has invalid coordinates: (%.2f, %.2f)", id, node.X, node.Y)
		}
		if node.Width <= 0 || node.Height <= 0 {
			t.Errorf("Node %s has invalid dimensions: %.2fx%.2f", id, node.Width, node.Height)
		}
	}
}

// AssertCanvasSize checks that layout has reasonable canvas size
func AssertCanvasSize(t *testing.T, layout *Layout) {
	t.Helper()

	if layout.Width <= 0 || layout.Height <= 0 {
		t.Errorf("Canvas has invalid size: %.2fx%.2f", layout.Width, layout.Height)
	}
}

// AssertNodeShape checks that a specific node has expected shape
func AssertNodeShape(t *testing.T, layout *Layout, nodeID, expectedShape string) {
	t.Helper()

	node, exists := layout.Nodes[nodeID]
	if !exists {
		t.Errorf("Node %s not found in layout", nodeID)
		return
	}

	if node.Shape != expectedShape {
		t.Errorf("Node %s has shape %s, expected %s", nodeID, node.Shape, expectedShape)
	}
}

// AssertNodeLabel checks that a specific node has expected label
func AssertNodeLabel(t *testing.T, layout *Layout, nodeID, expectedLabel string) {
	t.Helper()

	node, exists := layout.Nodes[nodeID]
	if !exists {
		t.Errorf("Node %s not found in layout", nodeID)
		return
	}

	if node.Label != expectedLabel {
		t.Errorf("Node %s has label '%s', expected '%s'", nodeID, node.Label, expectedLabel)
	}
}

// AssertNodesConnected checks that two nodes are connected by an edge
func AssertNodesConnected(t *testing.T, layout *Layout, fromID, toID string) {
	t.Helper()

	for _, edge := range layout.Edges {
		if edge.From == fromID && edge.To == toID {
			return // Found connection
		}
	}

	t.Errorf("No edge found connecting %s to %s", fromID, toID)
}

// AssertEdgeLabel checks that an edge has expected label
func AssertEdgeLabel(t *testing.T, layout *Layout, fromID, toID, expectedLabel string) {
	t.Helper()

	for _, edge := range layout.Edges {
		if edge.From == fromID && edge.To == toID {
			if edge.Label != expectedLabel {
				t.Errorf("Edge %s->%s has label '%s', expected '%s'", fromID, toID, edge.Label, expectedLabel)
			}
			return
		}
	}

	t.Errorf("No edge found connecting %s to %s", fromID, toID)
}

// AssertValidationError checks that validation produces expected error
func AssertValidationError(t *testing.T, diagram *Diagram, expectedErrorSubstring string) {
	t.Helper()

	validator := NewValidator()
	err := validator.Validate(diagram)

	if err == nil {
		t.Errorf("Expected validation error containing '%s', but validation passed", expectedErrorSubstring)
		return
	}

	if !strings.Contains(err.Error(), expectedErrorSubstring) {
		t.Errorf("Expected validation error to contain '%s', got: %v", expectedErrorSubstring, err)
	}
}

// AssertParseError checks that parsing produces expected error
func AssertParseError(t *testing.T, input string, expectedErrorSubstring string) {
	t.Helper()

	lexer := NewLexer(input)
	parser := NewParser(lexer)
	_, err := parser.ParseDiagram()

	if err == nil {
		t.Errorf("Expected parse error containing '%s', but parsing succeeded", expectedErrorSubstring)
		return
	}

	if !strings.Contains(err.Error(), expectedErrorSubstring) {
		t.Errorf("Expected parse error to contain '%s', got: %v", expectedErrorSubstring, err)
	}
}

// CountSVGElements counts occurrences of specific SVG elements
func CountSVGElements(svg string, elementType string) int {
	return strings.Count(svg, "<"+elementType)
}

// CreateTestLayout creates a simple test layout for testing SVG generation
func CreateTestLayout() *Layout {
	return &Layout{
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
			"B": {
				ID:     "B",
				X:      200,
				Y:      150,
				Width:  80,
				Height: 40,
				Label:  "Node B",
				Shape:  "ellipse",
			},
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
}
