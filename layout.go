package main

import (
	"math"
)

// LayoutNode represents a node with layout information
type LayoutNode struct {
	ID     string
	X      float64
	Y      float64
	Width  float64
	Height float64
	Label  string
	Shape  string
	Style  string
	Color  string
}

// LayoutEdge represents an edge with layout information
type LayoutEdge struct {
	From   string
	To     string
	Points []Point
	Label  string
	X      float64
	Y      float64
}

// Point represents a 2D coordinate
type Point struct {
	X float64
	Y float64
}

// Layout represents the complete layout information
type Layout struct {
	Width  float64
	Height float64
	Nodes  map[string]LayoutNode
	Edges  []LayoutEdge
}

// SimpleLayouter implements a basic tree layout algorithm
type SimpleLayouter struct {
	NodeWidth    float64
	NodeHeight   float64
	HorizontalGap float64
	VerticalGap   float64
}

func NewSimpleLayouter() *SimpleLayouter {
	return &SimpleLayouter{
		NodeWidth:    100.0,
		NodeHeight:   50.0,
		HorizontalGap: 80.0,
		VerticalGap:   80.0,
	}
}

// LayoutDiagram creates layout for the given diagram
func (l *SimpleLayouter) LayoutDiagram(diagram *Diagram) (*Layout, error) {
	if len(diagram.Nodes) == 0 {
		return &Layout{
			Width:  float64(diagram.Width),
			Height: float64(diagram.Height),
			Nodes:  make(map[string]LayoutNode),
			Edges:  []LayoutEdge{},
		}, nil
	}

	// Build dependency graph
	graph := l.buildDependencyGraph(diagram)
	
	// Calculate layout levels
	levels := l.calculateLevels(graph, diagram.Nodes)
	
	// Position nodes
	layout := l.positionNodes(levels, diagram)
	
	// Add edges
	l.addEdges(layout, diagram)
	
	// Calculate final canvas size
	l.calculateCanvasSize(layout)
	
	return layout, nil
}

// buildDependencyGraph builds a graph of node dependencies from edges
func (l *SimpleLayouter) buildDependencyGraph(diagram *Diagram) map[string][]string {
	graph := make(map[string][]string)
	
	// Initialize all nodes
	for _, node := range diagram.Nodes {
		graph[node.ID] = []string{}
	}
	
	// Add edges (from -> to relationships)
	for _, edge := range diagram.Edges {
		if children, exists := graph[edge.From]; exists {
			graph[edge.From] = append(children, edge.To)
		} else {
			graph[edge.From] = []string{edge.To}
		}
	}
	
	return graph
}

// calculateLevels assigns nodes to levels based on dependencies
func (l *SimpleLayouter) calculateLevels(graph map[string][]string, nodes []Node) [][]string {
	levels := [][]string{}
	visited := make(map[string]bool)
	inDegree := make(map[string]int)
	
	// Calculate in-degrees
	for _, node := range nodes {
		inDegree[node.ID] = 0
	}
	
	for _, children := range graph {
		for _, child := range children {
			inDegree[child]++
		}
	}
	
	// Find root nodes (nodes with no incoming edges)
	var currentLevel []string
	for _, node := range nodes {
		if inDegree[node.ID] == 0 {
			currentLevel = append(currentLevel, node.ID)
		}
	}
	
	// If no root nodes, pick the first node
	if len(currentLevel) == 0 && len(nodes) > 0 {
		currentLevel = append(currentLevel, nodes[0].ID)
	}
	
	// Process levels using topological sort
	for len(currentLevel) > 0 {
		levels = append(levels, currentLevel)
		nextLevel := []string{}
		
		for _, nodeID := range currentLevel {
			visited[nodeID] = true
			
			// Process children
			for _, child := range graph[nodeID] {
				if !visited[child] {
					inDegree[child]--
					if inDegree[child] == 0 {
						nextLevel = append(nextLevel, child)
					}
				}
			}
		}
		
		currentLevel = nextLevel
	}
	
	// Add any remaining nodes (in case of cycles)
	var remainingNodes []string
	for _, node := range nodes {
		if !visited[node.ID] {
			remainingNodes = append(remainingNodes, node.ID)
		}
	}
	
	if len(remainingNodes) > 0 {
		levels = append(levels, remainingNodes)
	}
	
	return levels
}

// positionNodes positions nodes based on their levels
func (l *SimpleLayouter) positionNodes(levels [][]string, diagram *Diagram) *Layout {
	layout := &Layout{
		Nodes: make(map[string]LayoutNode),
		Edges: []LayoutEdge{},
	}
	
	// Create a map for quick node lookup
	nodeMap := make(map[string]Node)
	for _, node := range diagram.Nodes {
		nodeMap[node.ID] = node
	}
	
	// Position nodes level by level
	for levelIndex, level := range levels {
		y := float64(levelIndex) * (l.NodeHeight + l.VerticalGap) + l.NodeHeight/2
		
		// Center nodes horizontally in each level
		totalWidth := float64(len(level))*l.NodeWidth + float64(len(level)-1)*l.HorizontalGap
		startX := -totalWidth/2 + l.NodeWidth/2
		
		for nodeIndex, nodeID := range level {
			x := startX + float64(nodeIndex)*(l.NodeWidth+l.HorizontalGap)
			
			// Get original node for attributes
			originalNode := nodeMap[nodeID]
			
			// Determine node shape
			shape := l.getNodeShape(originalNode)
			
			// Create layout node
			layoutNode := LayoutNode{
				ID:     nodeID,
				X:      x,
				Y:      y,
				Width:  l.NodeWidth,
				Height: l.NodeHeight,
				Label:  originalNode.Label,
				Shape:  shape,
				Style:  "solid",
				Color:  "black",
			}
			
			layout.Nodes[nodeID] = layoutNode
		}
	}
	
	return layout
}

// getNodeShape determines the shape of a node
func (l *SimpleLayouter) getNodeShape(node Node) string {
	if shape, exists := node.Attributes["shape"]; exists {
		return shape
	}
	return "rect" // default shape
}

// addEdges adds edges to the layout
func (l *SimpleLayouter) addEdges(layout *Layout, diagram *Diagram) {
	for _, edge := range diagram.Edges {
		fromNode, fromExists := layout.Nodes[edge.From]
		toNode, toExists := layout.Nodes[edge.To]
		
		if !fromExists || !toExists {
			continue // Skip edges to non-existent nodes
		}
		
		// Calculate edge points (simple straight line)
		points := []Point{
			{X: fromNode.X, Y: fromNode.Y + fromNode.Height/2},
			{X: toNode.X, Y: toNode.Y - toNode.Height/2},
		}
		
		// Calculate label position (midpoint)
		labelX := (fromNode.X + toNode.X) / 2
		labelY := (fromNode.Y + toNode.Y) / 2
		
		layoutEdge := LayoutEdge{
			From:   edge.From,
			To:     edge.To,
			Points: points,
			Label:  edge.Label,
			X:      labelX,
			Y:      labelY,
		}
		
		layout.Edges = append(layout.Edges, layoutEdge)
	}
}

// calculateCanvasSize calculates the final canvas size based on node positions
func (l *SimpleLayouter) calculateCanvasSize(layout *Layout) {
	if len(layout.Nodes) == 0 {
		layout.Width = 800
		layout.Height = 400
		return
	}
	
	minX, maxX := math.Inf(1), math.Inf(-1)
	minY, maxY := math.Inf(1), math.Inf(-1)
	
	for _, node := range layout.Nodes {
		nodeMinX := node.X - node.Width/2
		nodeMaxX := node.X + node.Width/2
		nodeMinY := node.Y - node.Height/2
		nodeMaxY := node.Y + node.Height/2
		
		minX = math.Min(minX, nodeMinX)
		maxX = math.Max(maxX, nodeMaxX)
		minY = math.Min(minY, nodeMinY)
		maxY = math.Max(maxY, nodeMaxY)
	}
	
	// Add padding
	padding := 40.0
	layout.Width = maxX - minX + padding*2
	layout.Height = maxY - minY + padding*2
	
	// Translate all nodes to ensure positive coordinates
	offsetX := -minX + padding
	offsetY := -minY + padding
	
	// Update node positions
	for id, node := range layout.Nodes {
		node.X += offsetX
		node.Y += offsetY
		layout.Nodes[id] = node
	}
	
	// Update edge positions
	for i, edge := range layout.Edges {
		for j, point := range edge.Points {
			point.X += offsetX
			point.Y += offsetY
			layout.Edges[i].Points[j] = point
		}
		layout.Edges[i].X += offsetX
		layout.Edges[i].Y += offsetY
	}
}