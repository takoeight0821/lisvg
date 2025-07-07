package main

import (
	"fmt"
	"strings"
)

// ValidatorError represents a validation error
type ValidatorError struct {
	Message string
	NodeID  string
	EdgeID  string
}

func (e ValidatorError) Error() string {
	return e.Message
}

// Validator validates AST diagrams
type Validator struct {
	errors []ValidatorError
}

func NewValidator() *Validator {
	return &Validator{
		errors: []ValidatorError{},
	}
}

// Validate validates the diagram AST
func (v *Validator) Validate(diagram *Diagram) error {
	v.errors = []ValidatorError{}

	// Validate node ID uniqueness
	v.validateNodeIDUniqueness(diagram)

	// Validate edge references
	v.validateEdgeReferences(diagram)

	// Validate node and edge attributes
	v.validateAttributes(diagram)

	if len(v.errors) > 0 {
		return fmt.Errorf("validation failed with %d errors: %s", len(v.errors), v.formatErrors())
	}

	return nil
}

func (v *Validator) validateNodeIDUniqueness(diagram *Diagram) {
	nodeIDs := make(map[string]bool)
	duplicates := []string{}

	for _, node := range diagram.Nodes {
		if node.ID == "" {
			v.errors = append(v.errors, ValidatorError{
				Message: "node ID cannot be empty",
				NodeID:  node.ID,
			})
			continue
		}

		if nodeIDs[node.ID] {
			duplicates = append(duplicates, node.ID)
		}
		nodeIDs[node.ID] = true
	}

	for _, id := range duplicates {
		v.errors = append(v.errors, ValidatorError{
			Message: fmt.Sprintf("duplicate node ID: %s", id),
			NodeID:  id,
		})
	}
}

func (v *Validator) validateEdgeReferences(diagram *Diagram) {
	nodeIDs := make(map[string]bool)
	for _, node := range diagram.Nodes {
		nodeIDs[node.ID] = true
	}

	for i, edge := range diagram.Edges {
		if edge.From == "" {
			v.errors = append(v.errors, ValidatorError{
				Message: fmt.Sprintf("edge %d: 'from' node ID cannot be empty", i),
				EdgeID:  fmt.Sprintf("edge_%d", i),
			})
		} else if !nodeIDs[edge.From] {
			v.errors = append(v.errors, ValidatorError{
				Message: fmt.Sprintf("edge %d: 'from' node '%s' does not exist", i, edge.From),
				EdgeID:  fmt.Sprintf("edge_%d", i),
			})
		}

		if edge.To == "" {
			v.errors = append(v.errors, ValidatorError{
				Message: fmt.Sprintf("edge %d: 'to' node ID cannot be empty", i),
				EdgeID:  fmt.Sprintf("edge_%d", i),
			})
		} else if !nodeIDs[edge.To] {
			v.errors = append(v.errors, ValidatorError{
				Message: fmt.Sprintf("edge %d: 'to' node '%s' does not exist", i, edge.To),
				EdgeID:  fmt.Sprintf("edge_%d", i),
			})
		}
	}
}

func (v *Validator) validateAttributes(diagram *Diagram) {
	// Validate node attributes
	for _, node := range diagram.Nodes {
		if shape, ok := node.Attributes["shape"]; ok {
			if !isValidShape(shape) {
				v.errors = append(v.errors, ValidatorError{
					Message: fmt.Sprintf("node '%s': invalid shape '%s'", node.ID, shape),
					NodeID:  node.ID,
				})
			}
		}
	}

	// Validate edge attributes
	for i, edge := range diagram.Edges {
		if style, ok := edge.Attributes["style"]; ok {
			if !isValidEdgeStyle(style) {
				v.errors = append(v.errors, ValidatorError{
					Message: fmt.Sprintf("edge %d: invalid style '%s'", i, style),
					EdgeID:  fmt.Sprintf("edge_%d", i),
				})
			}
		}
	}
}

func isValidShape(shape string) bool {
	validShapes := []string{
		"rect", "rectangle", "box",
		"ellipse", "circle", "oval",
		"diamond", "rhombus",
		"triangle", "trapezium",
		"polygon", "hexagon", "octagon",
	}

	for _, valid := range validShapes {
		if shape == valid {
			return true
		}
	}
	return false
}

func isValidEdgeStyle(style string) bool {
	validStyles := []string{
		"solid", "dashed", "dotted", "bold",
		"invis", "invisible",
	}

	for _, valid := range validStyles {
		if style == valid {
			return true
		}
	}
	return false
}

func (v *Validator) formatErrors() string {
	var messages []string
	for _, err := range v.errors {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, "; ")
}

// GetErrors returns all validation errors
func (v *Validator) GetErrors() []ValidatorError {
	return v.errors
}
