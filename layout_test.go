package main

import (
	"fmt"
	"math"
	"testing"
)

// TestSimpleLayouter tests the basic functionality of SimpleLayouter
func TestSimpleLayouter(t *testing.T) {
	layouter := NewSimpleLayouter()
	
	if layouter == nil {
		t.Errorf("NewSimpleLayouter returned nil")
	}
	
	// Test default values
	if layouter.NodeWidth != 100.0 {
		t.Errorf("Expected NodeWidth 100.0, got %f", layouter.NodeWidth)
	}
	if layouter.NodeHeight != 50.0 {
		t.Errorf("Expected NodeHeight 50.0, got %f", layouter.NodeHeight)
	}
	if layouter.HorizontalGap != 80.0 {
		t.Errorf("Expected HorizontalGap 80.0, got %f", layouter.HorizontalGap)
	}
	if layouter.VerticalGap != 80.0 {
		t.Errorf("Expected VerticalGap 80.0, got %f", layouter.VerticalGap)
	}
}

// TestLayoutEmptyDiagram tests layout with empty diagram
func TestLayoutEmptyDiagram(t *testing.T) {
	layouter := NewSimpleLayouter()
	diagram := &Diagram{
		Width:  800,
		Height: 400,
		Nodes:  []Node{},
		Edges:  []Edge{},
	}
	
	layout, err := layouter.LayoutDiagram(diagram)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	if layout.Width != 800 || layout.Height != 400 {
		t.Errorf("Expected size 800x400, got %fx%f", layout.Width, layout.Height)
	}
	
	if len(layout.Nodes) != 0 {
		t.Errorf("Expected 0 nodes, got %d", len(layout.Nodes))
	}
	
	if len(layout.Edges) != 0 {
		t.Errorf("Expected 0 edges, got %d", len(layout.Edges))
	}
}

// TestLayoutSingleNode tests layout with a single node
func TestLayoutSingleNode(t *testing.T) {
	layouter := NewSimpleLayouter()
	diagram := &Diagram{
		Nodes: []Node{
			{ID: "A", Label: "Single Node"},
		},
		Edges: []Edge{},
	}
	
	layout, err := layouter.LayoutDiagram(diagram)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	if len(layout.Nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(layout.Nodes))
	}
	
	node := layout.Nodes["A"]
	if node.ID != "A" {
		t.Errorf("Expected node ID 'A', got '%s'", node.ID)
	}
	if node.Label != "Single Node" {
		t.Errorf("Expected label 'Single Node', got '%s'", node.Label)
	}
	if node.Width != 100.0 || node.Height != 50.0 {
		t.Errorf("Expected size 100x50, got %fx%f", node.Width, node.Height)
	}
}

// TestLayoutLinearChain tests layout with a linear chain of nodes
func TestLayoutLinearChain(t *testing.T) {
	layouter := NewSimpleLayouter()
	diagram := &Diagram{
		Nodes: []Node{
			{ID: "A", Label: "Node A"},
			{ID: "B", Label: "Node B"},
			{ID: "C", Label: "Node C"},
		},
		Edges: []Edge{
			{From: "A", To: "B"},
			{From: "B", To: "C"},
		},
	}
	
	layout, err := layouter.LayoutDiagram(diagram)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	if len(layout.Nodes) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(layout.Nodes))
	}
	
	if len(layout.Edges) != 2 {
		t.Errorf("Expected 2 edges, got %d", len(layout.Edges))
	}
	
	// Check that nodes are arranged vertically
	nodeA := layout.Nodes["A"]
	nodeB := layout.Nodes["B"]
	nodeC := layout.Nodes["C"]
	
	if nodeA.Y >= nodeB.Y {
		t.Errorf("Expected A to be above B, but A.Y=%f, B.Y=%f", nodeA.Y, nodeB.Y)
	}
	if nodeB.Y >= nodeC.Y {
		t.Errorf("Expected B to be above C, but B.Y=%f, C.Y=%f", nodeB.Y, nodeC.Y)
	}
	
	// Check that nodes are aligned horizontally
	if nodeA.X != nodeB.X || nodeB.X != nodeC.X {
		t.Errorf("Expected nodes to be horizontally aligned, got A.X=%f, B.X=%f, C.X=%f", 
			nodeA.X, nodeB.X, nodeC.X)
	}
}

// TestLayoutBranching tests layout with branching structure
func TestLayoutBranching(t *testing.T) {
	layouter := NewSimpleLayouter()
	diagram := &Diagram{
		Nodes: []Node{
			{ID: "root", Label: "Root"},
			{ID: "left", Label: "Left"},
			{ID: "right", Label: "Right"},
			{ID: "bottom", Label: "Bottom"},
		},
		Edges: []Edge{
			{From: "root", To: "left"},
			{From: "root", To: "right"},
			{From: "left", To: "bottom"},
			{From: "right", To: "bottom"},
		},
	}
	
	layout, err := layouter.LayoutDiagram(diagram)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	if len(layout.Nodes) != 4 {
		t.Errorf("Expected 4 nodes, got %d", len(layout.Nodes))
	}
	
	root := layout.Nodes["root"]
	left := layout.Nodes["left"]
	right := layout.Nodes["right"]
	bottom := layout.Nodes["bottom"]
	
	// Root should be at the top
	if root.Y >= left.Y || root.Y >= right.Y {
		t.Errorf("Expected root to be above left and right")
	}
	
	// Left and right should be on the same level
	if left.Y != right.Y {
		t.Errorf("Expected left and right to be on same level, got left.Y=%f, right.Y=%f", 
			left.Y, right.Y)
	}
	
	// Bottom should be below left and right
	if bottom.Y <= left.Y || bottom.Y <= right.Y {
		t.Errorf("Expected bottom to be below left and right")
	}
}

// TestLayoutNodeShapes tests different node shapes
func TestLayoutNodeShapes(t *testing.T) {
	layouter := NewSimpleLayouter()
	diagram := &Diagram{
		Nodes: []Node{
			{ID: "rect", Label: "Rectangle", Attributes: map[string]string{"shape": "rect"}},
			{ID: "ellipse", Label: "Ellipse", Attributes: map[string]string{"shape": "ellipse"}},
			{ID: "diamond", Label: "Diamond", Attributes: map[string]string{"shape": "diamond"}},
		},
		Edges: []Edge{},
	}
	
	layout, err := layouter.LayoutDiagram(diagram)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	rectNode := layout.Nodes["rect"]
	ellipseNode := layout.Nodes["ellipse"]
	diamondNode := layout.Nodes["diamond"]
	
	if rectNode.Shape != "rect" {
		t.Errorf("Expected rect shape, got %s", rectNode.Shape)
	}
	if ellipseNode.Shape != "ellipse" {
		t.Errorf("Expected ellipse shape, got %s", ellipseNode.Shape)
	}
	if diamondNode.Shape != "diamond" {
		t.Errorf("Expected diamond shape, got %s", diamondNode.Shape)
	}
}

// TestLayoutEdgeGeneration tests edge generation
func TestLayoutEdgeGeneration(t *testing.T) {
	layouter := NewSimpleLayouter()
	diagram := &Diagram{
		Nodes: []Node{
			{ID: "A", Label: "Node A"},
			{ID: "B", Label: "Node B"},
		},
		Edges: []Edge{
			{From: "A", To: "B", Label: "connects"},
		},
	}
	
	layout, err := layouter.LayoutDiagram(diagram)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	if len(layout.Edges) != 1 {
		t.Errorf("Expected 1 edge, got %d", len(layout.Edges))
	}
	
	edge := layout.Edges[0]
	if edge.From != "A" || edge.To != "B" {
		t.Errorf("Expected edge from A to B, got from %s to %s", edge.From, edge.To)
	}
	
	if edge.Label != "connects" {
		t.Errorf("Expected label 'connects', got '%s'", edge.Label)
	}
	
	if len(edge.Points) != 2 {
		t.Errorf("Expected 2 points for edge, got %d", len(edge.Points))
	}
	
	// Check that edge points are reasonable
	nodeA := layout.Nodes["A"]
	nodeB := layout.Nodes["B"]
	
	if math.Abs(edge.Points[0].X-nodeA.X) > 0.1 {
		t.Errorf("Expected edge start X to match node A X")
	}
	if math.Abs(edge.Points[1].X-nodeB.X) > 0.1 {
		t.Errorf("Expected edge end X to match node B X")
	}
}

// TestLayoutCanvasSize tests canvas size calculation
func TestLayoutCanvasSize(t *testing.T) {
	layouter := NewSimpleLayouter()
	diagram := &Diagram{
		Nodes: []Node{
			{ID: "A", Label: "Node A"},
			{ID: "B", Label: "Node B"},
			{ID: "C", Label: "Node C"},
		},
		Edges: []Edge{
			{From: "A", To: "B"},
			{From: "B", To: "C"},
		},
	}
	
	layout, err := layouter.LayoutDiagram(diagram)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	// Canvas should be large enough to contain all nodes with padding
	if layout.Width <= 0 || layout.Height <= 0 {
		t.Errorf("Canvas size should be positive, got %fx%f", layout.Width, layout.Height)
	}
	
	// All nodes should have positive coordinates
	for id, node := range layout.Nodes {
		if node.X < 0 || node.Y < 0 {
			t.Errorf("Node %s has negative coordinates: (%f, %f)", id, node.X, node.Y)
		}
	}
}

// TestLayoutCyclicGraph tests handling of cyclic graphs
func TestLayoutCyclicGraph(t *testing.T) {
	layouter := NewSimpleLayouter()
	diagram := &Diagram{
		Nodes: []Node{
			{ID: "A", Label: "Node A"},
			{ID: "B", Label: "Node B"},
			{ID: "C", Label: "Node C"},
		},
		Edges: []Edge{
			{From: "A", To: "B"},
			{From: "B", To: "C"},
			{From: "C", To: "A"}, // Creates a cycle
		},
	}
	
	layout, err := layouter.LayoutDiagram(diagram)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	// Should still place all nodes
	if len(layout.Nodes) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(layout.Nodes))
	}
	
	// All nodes should have valid positions
	for id, node := range layout.Nodes {
		if node.X < 0 || node.Y < 0 {
			t.Errorf("Node %s has invalid position: (%f, %f)", id, node.X, node.Y)
		}
	}
}

// TestLayoutMultipleRoots tests layout with multiple root nodes
func TestLayoutMultipleRoots(t *testing.T) {
	layouter := NewSimpleLayouter()
	diagram := &Diagram{
		Nodes: []Node{
			{ID: "root1", Label: "Root 1"},
			{ID: "root2", Label: "Root 2"},
			{ID: "child1", Label: "Child 1"},
			{ID: "child2", Label: "Child 2"},
		},
		Edges: []Edge{
			{From: "root1", To: "child1"},
			{From: "root2", To: "child2"},
		},
	}
	
	layout, err := layouter.LayoutDiagram(diagram)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	root1 := layout.Nodes["root1"]
	root2 := layout.Nodes["root2"]
	child1 := layout.Nodes["child1"]
	child2 := layout.Nodes["child2"]
	
	// Both roots should be on the same level (top level)
	if root1.Y != root2.Y {
		t.Errorf("Expected roots on same level, got root1.Y=%f, root2.Y=%f", root1.Y, root2.Y)
	}
	
	// Children should be below roots
	if child1.Y <= root1.Y || child2.Y <= root2.Y {
		t.Errorf("Expected children to be below roots")
	}
}

// TestLayoutIsolatedNodes tests layout with isolated nodes
func TestLayoutIsolatedNodes(t *testing.T) {
	layouter := NewSimpleLayouter()
	diagram := &Diagram{
		Nodes: []Node{
			{ID: "connected1", Label: "Connected 1"},
			{ID: "connected2", Label: "Connected 2"},
			{ID: "isolated", Label: "Isolated"},
		},
		Edges: []Edge{
			{From: "connected1", To: "connected2"},
		},
	}
	
	layout, err := layouter.LayoutDiagram(diagram)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	// All nodes should be placed
	if len(layout.Nodes) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(layout.Nodes))
	}
	
	// Isolated node should also have valid position
	isolated := layout.Nodes["isolated"]
	if isolated.X < 0 || isolated.Y < 0 {
		t.Errorf("Isolated node has invalid position: (%f, %f)", isolated.X, isolated.Y)
	}
}

// TestBuildDependencyGraph tests dependency graph construction
func TestBuildDependencyGraph(t *testing.T) {
	layouter := NewSimpleLayouter()
	diagram := &Diagram{
		Nodes: []Node{
			{ID: "A"}, {ID: "B"}, {ID: "C"},
		},
		Edges: []Edge{
			{From: "A", To: "B"},
			{From: "A", To: "C"},
		},
	}
	
	graph := layouter.buildDependencyGraph(diagram)
	
	if len(graph) != 3 {
		t.Errorf("Expected 3 nodes in graph, got %d", len(graph))
	}
	
	if len(graph["A"]) != 2 {
		t.Errorf("Expected A to have 2 children, got %d", len(graph["A"]))
	}
	
	if len(graph["B"]) != 0 {
		t.Errorf("Expected B to have 0 children, got %d", len(graph["B"]))
	}
	
	if len(graph["C"]) != 0 {
		t.Errorf("Expected C to have 0 children, got %d", len(graph["C"]))
	}
}

// TestCalculateLevels tests level calculation
func TestCalculateLevels(t *testing.T) {
	layouter := NewSimpleLayouter()
	nodes := []Node{
		{ID: "A"}, {ID: "B"}, {ID: "C"},
	}
	graph := map[string][]string{
		"A": {"B"},
		"B": {"C"},
		"C": {},
	}
	
	levels := layouter.calculateLevels(graph, nodes)
	
	if len(levels) != 3 {
		t.Errorf("Expected 3 levels, got %d", len(levels))
	}
	
	// Check level 0 contains A
	if len(levels[0]) != 1 || levels[0][0] != "A" {
		t.Errorf("Expected level 0 to contain A, got %v", levels[0])
	}
	
	// Check level 1 contains B
	if len(levels[1]) != 1 || levels[1][0] != "B" {
		t.Errorf("Expected level 1 to contain B, got %v", levels[1])
	}
	
	// Check level 2 contains C
	if len(levels[2]) != 1 || levels[2][0] != "C" {
		t.Errorf("Expected level 2 to contain C, got %v", levels[2])
	}
}

// BenchmarkLayout benchmarks layout performance
func BenchmarkLayout(b *testing.B) {
	layouter := NewSimpleLayouter()
	
	// Create a moderately complex diagram
	nodes := make([]Node, 20)
	edges := make([]Edge, 19)
	
	for i := 0; i < 20; i++ {
		nodes[i] = Node{
			ID:    string(rune('A' + i)),
			Label: "Node " + string(rune('A' + i)),
		}
	}
	
	for i := 0; i < 19; i++ {
		edges[i] = Edge{
			From: string(rune('A' + i)),
			To:   string(rune('A' + i + 1)),
		}
	}
	
	diagram := &Diagram{
		Nodes: nodes,
		Edges: edges,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := layouter.LayoutDiagram(diagram)
		if err != nil {
			b.Fatalf("Layout error: %v", err)
		}
	}
}

// BenchmarkLargeLayout benchmarks layout with larger graphs
func BenchmarkLargeLayout(b *testing.B) {
	layouter := NewSimpleLayouter()
	
	// Create a large diagram
	nodes := make([]Node, 100)
	edges := make([]Edge, 99)
	
	for i := 0; i < 100; i++ {
		nodes[i] = Node{
			ID:    fmt.Sprintf("%c%d", rune('A' + i%26), i),
			Label: fmt.Sprintf("Node %d", i),
		}
	}
	
	for i := 0; i < 99; i++ {
		from := fmt.Sprintf("%c%d", rune('A' + i%26), i)
		to := fmt.Sprintf("%c%d", rune('A' + (i+1)%26), i+1)
		edges[i] = Edge{
			From: from,
			To:   to,
		}
	}
	
	diagram := &Diagram{
		Nodes: nodes,
		Edges: edges,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := layouter.LayoutDiagram(diagram)
		if err != nil {
			b.Fatalf("Layout error: %v", err)
		}
	}
}