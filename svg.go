package main

import (
	"fmt"
	"strings"
)

// SVGGenerator generates SVG from layout information
type SVGGenerator struct{}

func NewSVGGenerator() *SVGGenerator {
	return &SVGGenerator{}
}

// Generate creates SVG content from layout
func (s *SVGGenerator) Generate(layout *Layout, diagram *Diagram) string {
	var sb strings.Builder

	// Calculate viewBox with some padding
	padding := 20.0
	viewWidth := layout.Width + padding*2
	viewHeight := layout.Height + padding*2

	// SVG header
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">`,
		viewWidth, viewHeight, viewWidth, viewHeight))
	sb.WriteString("\n")

	// Add CSS styles
	sb.WriteString(s.generateCSS())

	// Create definitions for arrowheads
	sb.WriteString(s.generateDefs())

	// Transform coordinate system (Graphviz uses bottom-left origin, SVG uses top-left)
	sb.WriteString(fmt.Sprintf(`<g transform="translate(%.0f, %.0f) scale(1, -1)">`, padding, viewHeight-padding))
	sb.WriteString("\n")

	// Generate edges first (so they appear behind nodes)
	for _, edge := range layout.Edges {
		sb.WriteString(s.generateEdge(edge))
	}

	// Generate nodes
	for _, node := range layout.Nodes {
		sb.WriteString(s.generateNode(node))
	}

	// Close transform group
	sb.WriteString("</g>\n")

	// SVG footer
	sb.WriteString("</svg>\n")

	return sb.String()
}

// generateCSS generates CSS styles for the SVG
func (s *SVGGenerator) generateCSS() string {
	return `
  <style>
    .node {
      fill: #ffffff;
      stroke: #000000;
      stroke-width: 1;
    }
    .node.rect {
      fill: #f0f0f0;
    }
    .node.ellipse {
      fill: #e0e0e0;
    }
    .node.diamond {
      fill: #d0d0d0;
    }
    .node-label {
      font-family: Arial, sans-serif;
      font-size: 12px;
      text-anchor: middle;
      dominant-baseline: middle;
      fill: #000000;
      pointer-events: none;
    }
    .edge {
      fill: none;
      stroke: #000000;
      stroke-width: 1;
    }
    .edge-label {
      font-family: Arial, sans-serif;
      font-size: 10px;
      text-anchor: middle;
      dominant-baseline: middle;
      fill: #000000;
      pointer-events: none;
    }
  </style>
`
}

// generateDefs generates SVG definitions (arrowheads, etc.)
func (s *SVGGenerator) generateDefs() string {
	return `
  <defs>
    <marker id="arrowhead" markerWidth="10" markerHeight="7" 
            refX="9" refY="3.5" orient="auto">
      <polygon points="0 0, 10 3.5, 0 7" fill="#000000"/>
    </marker>
  </defs>
`
}

// generateNode generates SVG for a single node
func (s *SVGGenerator) generateNode(node LayoutNode) string {
	var sb strings.Builder

	shape := s.getNodeShape(node.Shape)
	class := fmt.Sprintf("node %s", shape)

	switch shape {
	case "rect":
		sb.WriteString(fmt.Sprintf(`  <rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" class="%s"/>`,
			node.X-node.Width/2, node.Y-node.Height/2, node.Width, node.Height, class))
	case "ellipse":
		sb.WriteString(fmt.Sprintf(`  <ellipse cx="%.2f" cy="%.2f" rx="%.2f" ry="%.2f" class="%s"/>`,
			node.X, node.Y, node.Width/2, node.Height/2, class))
	case "diamond":
		// Create diamond shape using polygon
		points := fmt.Sprintf("%.2f,%.2f %.2f,%.2f %.2f,%.2f %.2f,%.2f",
			node.X, node.Y-node.Height/2, // top
			node.X+node.Width/2, node.Y, // right
			node.X, node.Y+node.Height/2, // bottom
			node.X-node.Width/2, node.Y) // left
		sb.WriteString(fmt.Sprintf(`  <polygon points="%s" class="%s"/>`, points, class))
	default:
		// Default to ellipse
		sb.WriteString(fmt.Sprintf(`  <ellipse cx="%.2f" cy="%.2f" rx="%.2f" ry="%.2f" class="%s"/>`,
			node.X, node.Y, node.Width/2, node.Height/2, class))
	}
	sb.WriteString("\n")

	// Add label (flip Y coordinate back for text)
	if node.Label != "" {
		sb.WriteString(fmt.Sprintf(`  <text x="%.2f" y="%.2f" class="node-label" transform="scale(1, -1)">%s</text>`,
			node.X, -node.Y, s.escapeXML(node.Label)))
		sb.WriteString("\n")
	}

	return sb.String()
}

// generateEdge generates SVG for a single edge
func (s *SVGGenerator) generateEdge(edge LayoutEdge) string {
	var sb strings.Builder

	if len(edge.Points) < 2 {
		return ""
	}

	// Create path from points
	pathData := fmt.Sprintf("M %.2f %.2f", edge.Points[0].X, edge.Points[0].Y)

	if len(edge.Points) == 2 {
		// Simple line
		pathData += fmt.Sprintf(" L %.2f %.2f", edge.Points[1].X, edge.Points[1].Y)
	} else {
		// Bezier curve - use points as control points
		for i := 1; i < len(edge.Points); i++ {
			if i == 1 {
				pathData += fmt.Sprintf(" Q %.2f %.2f", edge.Points[i].X, edge.Points[i].Y)
			} else if i == len(edge.Points)-1 {
				pathData += fmt.Sprintf(" %.2f %.2f", edge.Points[i].X, edge.Points[i].Y)
			} else {
				pathData += fmt.Sprintf(" T %.2f %.2f", edge.Points[i].X, edge.Points[i].Y)
			}
		}
	}

	sb.WriteString(fmt.Sprintf(`  <path d="%s" class="edge" marker-end="url(#arrowhead)"/>`, pathData))
	sb.WriteString("\n")

	// Add edge label if present
	if edge.Label != "" && edge.X != 0 && edge.Y != 0 {
		sb.WriteString(fmt.Sprintf(`  <text x="%.2f" y="%.2f" class="edge-label" transform="scale(1, -1)">%s</text>`,
			edge.X, -edge.Y, s.escapeXML(edge.Label)))
		sb.WriteString("\n")
	}

	return sb.String()
}

// getNodeShape maps Graphviz shapes to SVG shapes
func (s *SVGGenerator) getNodeShape(shape string) string {
	switch strings.ToLower(shape) {
	case "box", "rect", "rectangle":
		return "rect"
	case "ellipse", "oval":
		return "ellipse"
	case "circle":
		return "ellipse"
	case "diamond", "rhombus":
		return "diamond"
	default:
		return "ellipse"
	}
}

// escapeXML escapes special XML characters
func (s *SVGGenerator) escapeXML(text string) string {
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	text = strings.ReplaceAll(text, "\"", "&quot;")
	text = strings.ReplaceAll(text, "'", "&apos;")
	return text
}

// GenerateWithCustomStyles generates SVG with custom styles from diagram
func (s *SVGGenerator) GenerateWithCustomStyles(layout *Layout, diagram *Diagram) string {
	var sb strings.Builder

	// Calculate viewBox with some padding
	padding := 20.0
	viewWidth := layout.Width + padding*2
	viewHeight := layout.Height + padding*2

	// SVG header
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">`,
		viewWidth, viewHeight, viewWidth, viewHeight))
	sb.WriteString("\n")

	// Add CSS styles with custom overrides
	sb.WriteString(s.generateCustomCSS(diagram))

	// Create definitions for arrowheads
	sb.WriteString(s.generateDefs())

	// Transform coordinate system
	sb.WriteString(fmt.Sprintf(`<g transform="translate(%.0f, %.0f) scale(1, -1)">`, padding, viewHeight-padding))
	sb.WriteString("\n")

	// Generate edges first
	for _, edge := range layout.Edges {
		sb.WriteString(s.generateEdge(edge))
	}

	// Generate nodes
	for _, node := range layout.Nodes {
		sb.WriteString(s.generateNode(node))
	}

	// Close transform group
	sb.WriteString("</g>\n")

	// SVG footer
	sb.WriteString("</svg>\n")

	return sb.String()
}

// generateCustomCSS generates CSS with custom styles from diagram
func (s *SVGGenerator) generateCustomCSS(diagram *Diagram) string {
	var sb strings.Builder

	sb.WriteString("  <style>\n")
	sb.WriteString("    .node {\n")
	sb.WriteString("      fill: #ffffff;\n")
	sb.WriteString("      stroke: #000000;\n")
	sb.WriteString("      stroke-width: 1;\n")

	// Apply custom node styles
	for key, value := range diagram.NodeStyle {
		switch key {
		case "fill", "stroke", "stroke-width":
			sb.WriteString(fmt.Sprintf("      %s: %s;\n", key, value))
		}
	}

	sb.WriteString("    }\n")

	sb.WriteString("    .edge {\n")
	sb.WriteString("      fill: none;\n")
	sb.WriteString("      stroke: #000000;\n")
	sb.WriteString("      stroke-width: 1;\n")

	// Apply custom edge styles
	for key, value := range diagram.EdgeStyle {
		switch key {
		case "stroke", "stroke-width", "stroke-dasharray":
			sb.WriteString(fmt.Sprintf("      %s: %s;\n", key, value))
		}
	}

	sb.WriteString("    }\n")

	sb.WriteString("    .node-label {\n")
	sb.WriteString("      font-family: Arial, sans-serif;\n")
	sb.WriteString("      font-size: 12px;\n")
	sb.WriteString("      text-anchor: middle;\n")
	sb.WriteString("      dominant-baseline: middle;\n")
	sb.WriteString("      fill: #000000;\n")
	sb.WriteString("      pointer-events: none;\n")
	sb.WriteString("    }\n")

	sb.WriteString("    .edge-label {\n")
	sb.WriteString("      font-family: Arial, sans-serif;\n")
	sb.WriteString("      font-size: 10px;\n")
	sb.WriteString("      text-anchor: middle;\n")
	sb.WriteString("      dominant-baseline: middle;\n")
	sb.WriteString("      fill: #000000;\n")
	sb.WriteString("      pointer-events: none;\n")
	sb.WriteString("    }\n")

	sb.WriteString("  </style>\n")

	return sb.String()
}
