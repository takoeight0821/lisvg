package main

import (
	"reflect"
	"testing"
)

// TestLexer tests the lexer functionality
func TestLexer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "basic tokens",
			input: "( ) atom :keyword \"string\"",
			expected: []Token{
				{TokenLParen, "("},
				{TokenRParen, ")"},
				{TokenAtom, "atom"},
				{TokenKeyword, "keyword"},
				{TokenString, "string"},
				{TokenEOF, ""},
			},
		},
		{
			name:  "empty input",
			input: "",
			expected: []Token{
				{TokenEOF, ""},
			},
		},
		{
			name:  "whitespace handling",
			input: "  (  \t\n  )  ",
			expected: []Token{
				{TokenLParen, "("},
				{TokenRParen, ")"},
				{TokenEOF, ""},
			},
		},
		{
			name:  "comment handling",
			input: "; this is a comment\n( atom ); another comment\n",
			expected: []Token{
				{TokenLParen, "("},
				{TokenAtom, "atom"},
				{TokenRParen, ")"},
				{TokenEOF, ""},
			},
		},
		{
			name:  "string with spaces",
			input: "\"hello world\"",
			expected: []Token{
				{TokenString, "hello world"},
				{TokenEOF, ""},
			},
		},
		{
			name:  "complex atoms",
			input: "node-style edge-style123 abc_def",
			expected: []Token{
				{TokenAtom, "node-style"},
				{TokenAtom, "edge-style123"},
				{TokenAtom, "abc_def"},
				{TokenEOF, ""},
			},
		},
		{
			name:  "multiple keywords",
			input: ":shape :label :stroke",
			expected: []Token{
				{TokenKeyword, "shape"},
				{TokenKeyword, "label"},
				{TokenKeyword, "stroke"},
				{TokenEOF, ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			var tokens []Token

			for {
				token := lexer.NextToken()
				tokens = append(tokens, token)
				if token.Type == TokenEOF {
					break
				}
			}

			if !reflect.DeepEqual(tokens, tt.expected) {
				t.Errorf("Expected tokens %+v, got %+v", tt.expected, tokens)
			}
		})
	}
}

// TestLexerEdgeCases tests edge cases in lexer
func TestLexerEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(*testing.T, *Lexer)
	}{
		{
			name:  "unclosed string",
			input: "\"unclosed string",
			check: func(t *testing.T, lexer *Lexer) {
				token := lexer.NextToken()
				if token.Type != TokenString {
					t.Errorf("Expected TokenString, got %v", token.Type)
				}
				if token.Value != "unclosed string" {
					t.Errorf("Expected 'unclosed string', got '%s'", token.Value)
				}
			},
		},
		{
			name:  "empty string",
			input: "\"\"",
			check: func(t *testing.T, lexer *Lexer) {
				token := lexer.NextToken()
				if token.Type != TokenString || token.Value != "" {
					t.Errorf("Expected empty string token, got %+v", token)
				}
			},
		},
		{
			name:  "special characters in atoms",
			input: "node-123 _private $special",
			check: func(t *testing.T, lexer *Lexer) {
				expected := []string{"node-123", "_private", "$special"}
				for _, exp := range expected {
					token := lexer.NextToken()
					if token.Type != TokenAtom || token.Value != exp {
						t.Errorf("Expected atom '%s', got %+v", exp, token)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tt.check(t, lexer)
		})
	}
}

// TestParser tests the parser functionality
func TestParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *Diagram
		hasError bool
	}{
		{
			name: "basic diagram",
			input: `(diagram
				(size 800 400)
				(nodes
					(id "A" :label "Node A"))
				(edges
					("A" "B")))`,
			expected: &Diagram{
				Width:     800,
				Height:    400,
				NodeStyle: map[string]string{},
				EdgeStyle: map[string]string{},
				Nodes: []Node{
					{ID: "A", Label: "Node A", Attributes: map[string]string{}},
				},
				Edges: []Edge{
					{From: "A", To: "B", Label: "", Attributes: map[string]string{}},
				},
			},
			hasError: false,
		},
		{
			name: "diagram with styles",
			input: `(diagram
				(node-style :shape "rect" :fill "#fff")
				(edge-style :stroke "#000")
				(nodes
					(id "X" :label "Test")))`,
			expected: &Diagram{
				Width:  800,
				Height: 400,
				NodeStyle: map[string]string{
					"shape": "rect",
					"fill":  "#fff",
				},
				EdgeStyle: map[string]string{
					"stroke": "#000",
				},
				Nodes: []Node{
					{ID: "X", Label: "Test", Attributes: map[string]string{}},
				},
				Edges: []Edge{},
			},
			hasError: false,
		},
		{
			name: "node with attributes",
			input: `(diagram
				(nodes
					(id "N1" :label "Node 1" :shape "ellipse" :color "red")))`,
			expected: &Diagram{
				Width:     800,
				Height:    400,
				NodeStyle: map[string]string{},
				EdgeStyle: map[string]string{},
				Nodes: []Node{
					{
						ID:    "N1",
						Label: "Node 1",
						Attributes: map[string]string{
							"shape": "ellipse",
							"color": "red",
						},
					},
				},
				Edges: []Edge{},
			},
			hasError: false,
		},
		{
			name: "edge with label",
			input: `(diagram
				(nodes
					(id "A")
					(id "B"))
				(edges
					("A" "B" :label "connects")))`,
			expected: &Diagram{
				Width:     800,
				Height:    400,
				NodeStyle: map[string]string{},
				EdgeStyle: map[string]string{},
				Nodes: []Node{
					{ID: "A", Label: "A", Attributes: map[string]string{}},
					{ID: "B", Label: "B", Attributes: map[string]string{}},
				},
				Edges: []Edge{
					{From: "A", To: "B", Label: "connects", Attributes: map[string]string{}},
				},
			},
			hasError: false,
		},
		{
			name: "empty diagram",
			input: `(diagram
				(size 400 200))`,
			expected: &Diagram{
				Width:     400,
				Height:    200,
				NodeStyle: map[string]string{},
				EdgeStyle: map[string]string{},
				Nodes:     []Node{},
				Edges:     []Edge{},
			},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			parser := NewParser(lexer)
			
			diagram, err := parser.ParseDiagram()
			
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error, but parsing succeeded")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if !reflect.DeepEqual(diagram, tt.expected) {
				t.Errorf("Expected diagram %+v, got %+v", tt.expected, diagram)
			}
		})
	}
}

// TestParserErrors tests error handling in parser
func TestParserErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "missing diagram keyword",
			input: "(not-diagram)",
		},
		{
			name:  "invalid size arguments",
			input: "(diagram (size invalid))",
		},
		{
			name:  "malformed node",
			input: "(diagram (nodes (invalid)))",
		},
		{
			name:  "incomplete edge",
			input: "(diagram (edges (\"A\")))",
		},
		{
			name:  "missing closing paren",
			input: "(diagram (size 800 400)",
		},
		{
			name:  "unknown directive",
			input: "(diagram (unknown-directive))",
		},
		{
			name:  "invalid size values",
			input: "(diagram (size abc def))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			parser := NewParser(lexer)
			
			_, err := parser.ParseDiagram()
			
			if err == nil {
				t.Errorf("Expected error for input: %s", tt.input)
			}
		})
	}
}

// TestParserComplexCases tests complex parsing scenarios
func TestParserComplexCases(t *testing.T) {
	t.Run("japanese characters", func(t *testing.T) {
		input := `(diagram
			(nodes
				(id "開始" :label "開始ノード")
				(id "終了" :label "終了ノード"))
			(edges
				("開始" "終了" :label "処理")))`
		
		lexer := NewLexer(input)
		parser := NewParser(lexer)
		
		diagram, err := parser.ParseDiagram()
		if err != nil {
			t.Errorf("Unexpected error with Japanese characters: %v", err)
		}
		
		if len(diagram.Nodes) != 2 {
			t.Errorf("Expected 2 nodes, got %d", len(diagram.Nodes))
		}
		
		if diagram.Nodes[0].Label != "開始ノード" {
			t.Errorf("Expected Japanese label, got %s", diagram.Nodes[0].Label)
		}
	})

	t.Run("mixed quotes and atoms", func(t *testing.T) {
		input := `(diagram
			(nodes
				(id atom-id :label "String Label")
				(id "string-id" :label atom-label)))`
		
		lexer := NewLexer(input)
		parser := NewParser(lexer)
		
		diagram, err := parser.ParseDiagram()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		
		if len(diagram.Nodes) != 2 {
			t.Errorf("Expected 2 nodes, got %d", len(diagram.Nodes))
		}
	})

	t.Run("nested comments", func(t *testing.T) {
		input := `; Top level comment
		(diagram
			; Size comment
			(size 800 400)
			; Nodes comment
			(nodes
				(id "A" :label "Node A")) ; Inline comment
			; Edges comment
			(edges
				("A" "B"))) ; Final comment`
		
		lexer := NewLexer(input)
		parser := NewParser(lexer)
		
		diagram, err := parser.ParseDiagram()
		if err != nil {
			t.Errorf("Unexpected error with comments: %v", err)
		}
		
		if diagram.Width != 800 || diagram.Height != 400 {
			t.Errorf("Expected size 800x400, got %dx%d", diagram.Width, diagram.Height)
		}
	})
}

// BenchmarkLexer benchmarks lexer performance
func BenchmarkLexer(b *testing.B) {
	input := `(diagram
		(size 800 400)
		(node-style :shape "rect")
		(edge-style :stroke "#555")
		(nodes
			(id "A" :label "Start")
			(id "B" :label "Process")
			(id "C" :label "End"))
		(edges
			("A" "B" :label "next")
			("B" "C")))`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lexer := NewLexer(input)
		for {
			token := lexer.NextToken()
			if token.Type == TokenEOF {
				break
			}
		}
	}
}

// BenchmarkParser benchmarks parser performance
func BenchmarkParser(b *testing.B) {
	input := `(diagram
		(size 800 400)
		(node-style :shape "rect")
		(edge-style :stroke "#555")
		(nodes
			(id "A" :label "Start")
			(id "B" :label "Process")
			(id "C" :label "End"))
		(edges
			("A" "B" :label "next")
			("B" "C")))`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lexer := NewLexer(input)
		parser := NewParser(lexer)
		_, err := parser.ParseDiagram()
		if err != nil {
			b.Fatalf("Parse error: %v", err)
		}
	}
}