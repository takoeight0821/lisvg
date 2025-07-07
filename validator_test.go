package main

import (
	"strings"
	"testing"
)

// TestValidatorNodeIDUniqueness tests node ID uniqueness validation
func TestValidatorNodeIDUniqueness(t *testing.T) {
	tests := []struct {
		name        string
		diagram     *Diagram
		expectErr   bool
		errContains string
	}{
		{
			name: "unique node IDs",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "A", Label: "Node A"},
					{ID: "B", Label: "Node B"},
					{ID: "C", Label: "Node C"},
				},
			},
			expectErr: false,
		},
		{
			name: "duplicate node IDs",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "A", Label: "Node A1"},
					{ID: "B", Label: "Node B"},
					{ID: "A", Label: "Node A2"},
				},
			},
			expectErr:   true,
			errContains: "duplicate node ID: A",
		},
		{
			name: "empty node ID",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "", Label: "Empty ID"},
					{ID: "B", Label: "Node B"},
				},
			},
			expectErr:   true,
			errContains: "node ID cannot be empty",
		},
		{
			name: "multiple duplicates",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "A", Label: "Node A1"},
					{ID: "B", Label: "Node B1"},
					{ID: "A", Label: "Node A2"},
					{ID: "B", Label: "Node B2"},
				},
			},
			expectErr:   true,
			errContains: "duplicate node ID",
		},
		{
			name: "no nodes",
			diagram: &Diagram{
				Nodes: []Node{},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewValidator()
			err := validator.Validate(tt.diagram)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error, but validation passed")
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}
}

// TestValidatorEdgeReferences tests edge reference validation
func TestValidatorEdgeReferences(t *testing.T) {
	tests := []struct {
		name        string
		diagram     *Diagram
		expectErr   bool
		errContains string
	}{
		{
			name: "valid edge references",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "A", Label: "Node A"},
					{ID: "B", Label: "Node B"},
				},
				Edges: []Edge{
					{From: "A", To: "B"},
				},
			},
			expectErr: false,
		},
		{
			name: "invalid from reference",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "A", Label: "Node A"},
					{ID: "B", Label: "Node B"},
				},
				Edges: []Edge{
					{From: "X", To: "B"},
				},
			},
			expectErr:   true,
			errContains: "'from' node 'X' does not exist",
		},
		{
			name: "invalid to reference",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "A", Label: "Node A"},
					{ID: "B", Label: "Node B"},
				},
				Edges: []Edge{
					{From: "A", To: "Y"},
				},
			},
			expectErr:   true,
			errContains: "'to' node 'Y' does not exist",
		},
		{
			name: "empty from node",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "A", Label: "Node A"},
				},
				Edges: []Edge{
					{From: "", To: "A"},
				},
			},
			expectErr:   true,
			errContains: "'from' node ID cannot be empty",
		},
		{
			name: "empty to node",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "A", Label: "Node A"},
				},
				Edges: []Edge{
					{From: "A", To: ""},
				},
			},
			expectErr:   true,
			errContains: "'to' node ID cannot be empty",
		},
		{
			name: "multiple invalid edges",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "A", Label: "Node A"},
				},
				Edges: []Edge{
					{From: "X", To: "Y"},
					{From: "A", To: "Z"},
				},
			},
			expectErr:   true,
			errContains: "does not exist",
		},
		{
			name: "self-reference edge",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "A", Label: "Node A"},
				},
				Edges: []Edge{
					{From: "A", To: "A"},
				},
			},
			expectErr: false, // Self-references should be valid
		},
		{
			name: "no edges",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "A", Label: "Node A"},
				},
				Edges: []Edge{},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewValidator()
			err := validator.Validate(tt.diagram)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error, but validation passed")
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}
}

// TestValidatorAttributes tests attribute validation
func TestValidatorAttributes(t *testing.T) {
	tests := []struct {
		name        string
		diagram     *Diagram
		expectErr   bool
		errContains string
	}{
		{
			name: "valid node shapes",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "A", Attributes: map[string]string{"shape": "rect"}},
					{ID: "B", Attributes: map[string]string{"shape": "ellipse"}},
					{ID: "C", Attributes: map[string]string{"shape": "diamond"}},
				},
			},
			expectErr: false,
		},
		{
			name: "invalid node shape",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "A", Attributes: map[string]string{"shape": "triangle"}},
				},
			},
			expectErr: false, // triangle is valid now per isValidShape function
		},
		{
			name: "invalid node shape - completely wrong",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "A", Attributes: map[string]string{"shape": "invalid-shape"}},
				},
			},
			expectErr:   true,
			errContains: "invalid shape 'invalid-shape'",
		},
		{
			name: "valid edge styles",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "A"}, {ID: "B"},
				},
				Edges: []Edge{
					{From: "A", To: "B", Attributes: map[string]string{"style": "solid"}},
				},
			},
			expectErr: false,
		},
		{
			name: "invalid edge style",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "A"}, {ID: "B"},
				},
				Edges: []Edge{
					{From: "A", To: "B", Attributes: map[string]string{"style": "invalid-style"}},
				},
			},
			expectErr:   true,
			errContains: "invalid style 'invalid-style'",
		},
		{
			name: "no attributes",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "A", Attributes: map[string]string{}},
				},
				Edges: []Edge{
					{From: "A", To: "A", Attributes: map[string]string{}},
				},
			},
			expectErr: false,
		},
		{
			name: "mixed valid and invalid",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "A", Attributes: map[string]string{"shape": "rect"}},
					{ID: "B", Attributes: map[string]string{"shape": "bad-shape"}},
				},
				Edges: []Edge{
					{From: "A", To: "B", Attributes: map[string]string{"style": "solid"}},
				},
			},
			expectErr:   true,
			errContains: "invalid shape 'bad-shape'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewValidator()
			err := validator.Validate(tt.diagram)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error, but validation passed")
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}
}

// TestValidatorGetErrors tests error collection
func TestValidatorGetErrors(t *testing.T) {
	diagram := &Diagram{
		Nodes: []Node{
			{ID: "A", Attributes: map[string]string{"shape": "invalid"}},
			{ID: "A", Label: "Duplicate"}, // Duplicate ID
			{ID: "", Label: "Empty ID"},   // Empty ID
		},
		Edges: []Edge{
			{From: "A", To: "B"}, // B doesn't exist
			{From: "X", To: "Y"}, // Both don't exist
		},
	}

	validator := NewValidator()
	err := validator.Validate(diagram)

	if err == nil {
		t.Errorf("Expected validation to fail")
		return
	}

	errors := validator.GetErrors()
	if len(errors) == 0 {
		t.Errorf("Expected validation errors to be collected")
	}

	// Check that we have multiple error types
	hasNodeError := false
	hasEdgeError := false
	for _, e := range errors {
		if e.NodeID != "" {
			hasNodeError = true
		}
		if e.EdgeID != "" {
			hasEdgeError = true
		}
	}

	if !hasNodeError {
		t.Errorf("Expected at least one node error")
	}
	if !hasEdgeError {
		t.Errorf("Expected at least one edge error")
	}
}

// TestValidatorComplexScenarios tests complex validation scenarios
func TestValidatorComplexScenarios(t *testing.T) {
	t.Run("large valid diagram", func(t *testing.T) {
		nodes := make([]Node, 100)
		edges := make([]Edge, 99)

		for i := 0; i < 100; i++ {
			nodes[i] = Node{
				ID:         string(rune('A' + i)),
				Label:      "Node " + string(rune('A'+i)),
				Attributes: map[string]string{"shape": "rect"},
			}
		}

		for i := 0; i < 99; i++ {
			edges[i] = Edge{
				From:       string(rune('A' + i)),
				To:         string(rune('A' + i + 1)),
				Attributes: map[string]string{"style": "solid"},
			}
		}

		diagram := &Diagram{
			Nodes: nodes,
			Edges: edges,
		}

		validator := NewValidator()
		err := validator.Validate(diagram)

		if err != nil {
			t.Errorf("Expected large valid diagram to pass validation, got: %v", err)
		}
	})

	t.Run("empty diagram", func(t *testing.T) {
		diagram := &Diagram{
			Nodes: []Node{},
			Edges: []Edge{},
		}

		validator := NewValidator()
		err := validator.Validate(diagram)

		if err != nil {
			t.Errorf("Expected empty diagram to pass validation, got: %v", err)
		}
	})

	t.Run("multiple validation passes", func(t *testing.T) {
		diagram := &Diagram{
			Nodes: []Node{
				{ID: "A", Label: "Node A"},
			},
		}

		validator := NewValidator()

		// First validation
		err1 := validator.Validate(diagram)
		if err1 != nil {
			t.Errorf("First validation failed: %v", err1)
		}

		// Second validation should also pass
		err2 := validator.Validate(diagram)
		if err2 != nil {
			t.Errorf("Second validation failed: %v", err2)
		}

		// Errors should be reset between validations
		errors := validator.GetErrors()
		if len(errors) != 0 {
			t.Errorf("Expected no errors after successful validation, got %d", len(errors))
		}
	})
}

// TestIsValidShape tests the shape validation function
func TestIsValidShape(t *testing.T) {
	validShapes := []string{
		"rect", "rectangle", "box",
		"ellipse", "circle", "oval",
		"diamond", "rhombus",
		"triangle", "trapezium",
		"polygon", "hexagon", "octagon",
	}

	invalidShapes := []string{
		"star", "invalid", "random", "123", "",
	}

	for _, shape := range validShapes {
		if !isValidShape(shape) {
			t.Errorf("Expected '%s' to be valid shape", shape)
		}
	}

	for _, shape := range invalidShapes {
		if isValidShape(shape) {
			t.Errorf("Expected '%s' to be invalid shape", shape)
		}
	}
}

// TestIsValidEdgeStyle tests the edge style validation function
func TestIsValidEdgeStyle(t *testing.T) {
	validStyles := []string{
		"solid", "dashed", "dotted", "bold",
		"invis", "invisible",
	}

	invalidStyles := []string{
		"wavy", "invalid", "random", "123", "",
	}

	for _, style := range validStyles {
		if !isValidEdgeStyle(style) {
			t.Errorf("Expected '%s' to be valid edge style", style)
		}
	}

	for _, style := range invalidStyles {
		if isValidEdgeStyle(style) {
			t.Errorf("Expected '%s' to be invalid edge style", style)
		}
	}
}

// TestValidatorEdgeCases tests additional edge cases and error scenarios
func TestValidatorEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		diagram     *Diagram
		expectErr   bool
		errContains string
	}{
		{
			name: "node with whitespace-only ID",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "   ", Label: "Whitespace ID"},
				},
			},
			expectErr: false, // Validator doesn't trim whitespace
		},
		{
			name: "edge with whitespace-only from",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "A", Label: "Node A"},
				},
				Edges: []Edge{
					{From: "  ", To: "A"},
				},
			},
			expectErr:   true,
			errContains: "'from' node '  ' does not exist",
		},
		{
			name: "very long node ID",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: strings.Repeat("A", 1000), Label: "Long ID"},
				},
			},
			expectErr: false, // Should be valid
		},
		{
			name: "special characters in node ID",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "node-123_test.v2", Label: "Special chars"},
					{ID: "@node!", Label: "More special"},
				},
			},
			expectErr: false, // Should be valid
		},
		{
			name: "nil attributes map",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "A", Label: "Node A", Attributes: nil},
				},
				Edges: []Edge{
					{From: "A", To: "A", Attributes: nil},
				},
			},
			expectErr: false, // Should handle nil gracefully
		},
		{
			name: "multiple shape attributes",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "A", Attributes: map[string]string{
						"shape": "rect",
						"Shape": "ellipse", // Different case
					}},
				},
			},
			expectErr: false, // Both should be checked
		},
		{
			name: "cyclic edges",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "A"},
					{ID: "B"},
					{ID: "C"},
				},
				Edges: []Edge{
					{From: "A", To: "B"},
					{From: "B", To: "C"},
					{From: "C", To: "A"},
				},
			},
			expectErr: false, // Cycles should be allowed
		},
		{
			name: "duplicate edges",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "A"},
					{ID: "B"},
				},
				Edges: []Edge{
					{From: "A", To: "B"},
					{From: "A", To: "B"}, // Duplicate
				},
			},
			expectErr: false, // Duplicates might be allowed
		},
		{
			name: "edge style case variations",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "A"},
				},
				Edges: []Edge{
					{From: "A", To: "A", Attributes: map[string]string{"style": "SOLID"}},
				},
			},
			expectErr:   true,
			errContains: "invalid style 'SOLID'", // Case sensitive
		},
		{
			name: "empty attribute values",
			diagram: &Diagram{
				Nodes: []Node{
					{ID: "A", Attributes: map[string]string{
						"shape": "",
						"color": "",
					}},
				},
			},
			expectErr:   true,
			errContains: "invalid shape ''",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewValidator()
			err := validator.Validate(tt.diagram)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error, but validation passed")
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}
}

// TestValidatorNilHandling tests validator behavior with nil inputs
func TestValidatorNilHandling(t *testing.T) {
	t.Run("nil diagram", func(t *testing.T) {
		validator := NewValidator()

		// Validator doesn't handle nil diagrams, it will panic
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for nil diagram")
			}
		}()

		_ = validator.Validate(nil)
	})

	t.Run("nil nodes slice", func(t *testing.T) {
		diagram := &Diagram{
			Nodes: nil,
			Edges: []Edge{},
		}
		validator := NewValidator()
		err := validator.Validate(diagram)
		if err != nil {
			t.Errorf("Expected nil nodes slice to be valid, got: %v", err)
		}
	})

	t.Run("nil edges slice", func(t *testing.T) {
		diagram := &Diagram{
			Nodes: []Node{{ID: "A"}},
			Edges: nil,
		}
		validator := NewValidator()
		err := validator.Validate(diagram)
		if err != nil {
			t.Errorf("Expected nil edges slice to be valid, got: %v", err)
		}
	})
}

// TestValidatorErrorAccumulation tests that validator collects all errors
func TestValidatorErrorAccumulation(t *testing.T) {
	diagram := &Diagram{
		Nodes: []Node{
			{ID: "", Label: "Empty1"},
			{ID: "", Label: "Empty2"},
			{ID: "A", Attributes: map[string]string{"shape": "invalid1"}},
			{ID: "B", Attributes: map[string]string{"shape": "invalid2"}},
			{ID: "A", Label: "Duplicate"},
		},
		Edges: []Edge{
			{From: "X", To: "Y"},
			{From: "A", To: "Z"},
			{From: "", To: "A"},
			{From: "A", To: "", Attributes: map[string]string{"style": "wrong"}},
		},
	}

	validator := NewValidator()
	err := validator.Validate(diagram)

	if err == nil {
		t.Fatalf("Expected validation to fail")
	}

	errors := validator.GetErrors()
	if len(errors) < 5 {
		t.Errorf("Expected at least 5 errors, got %d", len(errors))
	}

	// Verify error variety
	emptyIDCount := 0
	duplicateCount := 0
	invalidShapeCount := 0
	invalidRefCount := 0

	for _, e := range errors {
		if strings.Contains(e.Message, "empty") {
			emptyIDCount++
		}
		if strings.Contains(e.Message, "duplicate") {
			duplicateCount++
		}
		if strings.Contains(e.Message, "shape") {
			invalidShapeCount++
		}
		if strings.Contains(e.Message, "does not exist") {
			invalidRefCount++
		}
	}

	if emptyIDCount == 0 {
		t.Errorf("Expected empty ID errors")
	}
	if duplicateCount == 0 {
		t.Errorf("Expected duplicate ID errors")
	}
	if invalidShapeCount == 0 {
		t.Errorf("Expected invalid shape errors")
	}
	if invalidRefCount == 0 {
		t.Errorf("Expected invalid reference errors")
	}
}

// BenchmarkValidator benchmarks validator performance
func BenchmarkValidator(b *testing.B) {
	// Create a reasonably complex diagram
	nodes := make([]Node, 50)
	edges := make([]Edge, 49)

	for i := 0; i < 50; i++ {
		nodes[i] = Node{
			ID:         string(rune('A' + i)),
			Label:      "Node " + string(rune('A'+i)),
			Attributes: map[string]string{"shape": "rect"},
		}
	}

	for i := 0; i < 49; i++ {
		edges[i] = Edge{
			From:       string(rune('A' + i)),
			To:         string(rune('A' + i + 1)),
			Attributes: map[string]string{"style": "solid"},
		}
	}

	diagram := &Diagram{
		Nodes: nodes,
		Edges: edges,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator := NewValidator()
		err := validator.Validate(diagram)
		if err != nil {
			b.Fatalf("Validation error: %v", err)
		}
	}
}
