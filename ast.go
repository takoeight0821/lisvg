package main

import (
	"fmt"
	"strconv"
)

// Diagram represents the root AST node
type Diagram struct {
	Width     int
	Height    int
	NodeStyle map[string]string
	EdgeStyle map[string]string
	Nodes     []Node
	Edges     []Edge
}

// Node represents a diagram node
type Node struct {
	ID         string
	Label      string
	Attributes map[string]string
}

// Edge represents a diagram edge
type Edge struct {
	From       string
	To         string
	Label      string
	Attributes map[string]string
}

// Token represents a lexical token
type Token struct {
	Type  TokenType
	Value string
}

type TokenType int

const (
	TokenLParen TokenType = iota
	TokenRParen
	TokenAtom
	TokenKeyword
	TokenString
	TokenEOF
)

// Lexer tokenizes S-expressions
type Lexer struct {
	input string
	pos   int
	ch    byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.pos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.pos]
	}
	l.pos++
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipComment() {
	if l.ch == ';' {
		// Skip until end of line
		for l.ch != '\n' && l.ch != '\r' && l.ch != 0 {
			l.readChar()
		}
	}
}

func (l *Lexer) readString() string {
	position := l.pos
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}
	return l.input[position-1 : l.pos-1]
}

func (l *Lexer) readAtom() string {
	position := l.pos - 1
	for l.ch != ' ' && l.ch != '\t' && l.ch != '\n' && l.ch != '\r' && l.ch != '(' && l.ch != ')' && l.ch != 0 {
		l.readChar()
	}
	return l.input[position : l.pos-1]
}

func (l *Lexer) NextToken() Token {
	for {
		l.skipWhitespace()

		// Skip comments
		if l.ch == ';' {
			l.skipComment()
			continue
		}
		break
	}

	switch l.ch {
	case '(':
		l.readChar()
		return Token{TokenLParen, "("}
	case ')':
		l.readChar()
		return Token{TokenRParen, ")"}
	case '"':
		l.readChar()
		value := l.readString()
		l.readChar() // skip closing quote
		return Token{TokenString, value}
	case 0:
		return Token{TokenEOF, ""}
	default:
		if l.ch == ':' {
			l.readChar()
			atom := l.readAtom()
			return Token{TokenKeyword, atom}
		}
		atom := l.readAtom()
		return Token{TokenAtom, atom}
	}
}

// Parser parses S-expressions into AST
type Parser struct {
	lexer *Lexer
	cur   Token
	peek  Token
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{lexer: lexer}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.cur = p.peek
	p.peek = p.lexer.NextToken()
}

func (p *Parser) ParseDiagram() (*Diagram, error) {
	if p.cur.Type != TokenLParen {
		return nil, fmt.Errorf("expected '(', got %s", p.cur.Value)
	}
	p.nextToken()

	if p.cur.Type != TokenAtom || p.cur.Value != "diagram" {
		return nil, fmt.Errorf("expected 'diagram', got %s", p.cur.Value)
	}
	p.nextToken()

	diagram := &Diagram{
		Width:     800,
		Height:    400,
		NodeStyle: make(map[string]string),
		EdgeStyle: make(map[string]string),
		Nodes:     []Node{},
		Edges:     []Edge{},
	}

	for p.cur.Type != TokenRParen && p.cur.Type != TokenEOF {
		if p.cur.Type != TokenLParen {
			return nil, fmt.Errorf("expected '(', got %s", p.cur.Value)
		}
		p.nextToken()

		switch p.cur.Value {
		case "size":
			if err := p.parseSize(diagram); err != nil {
				return nil, err
			}
		case "node-style":
			if err := p.parseNodeStyle(diagram); err != nil {
				return nil, err
			}
		case "edge-style":
			if err := p.parseEdgeStyle(diagram); err != nil {
				return nil, err
			}
		case "nodes":
			if err := p.parseNodes(diagram); err != nil {
				return nil, err
			}
		case "edges":
			if err := p.parseEdges(diagram); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unknown directive: %s", p.cur.Value)
		}
	}

	if p.cur.Type != TokenRParen {
		return nil, fmt.Errorf("expected ')', got %s", p.cur.Value)
	}

	return diagram, nil
}

func (p *Parser) parseSize(diagram *Diagram) error {
	p.nextToken() // consume 'size'

	if p.cur.Type != TokenAtom {
		return fmt.Errorf("expected width, got %s", p.cur.Value)
	}
	width, err := strconv.Atoi(p.cur.Value)
	if err != nil {
		return fmt.Errorf("invalid width: %s", p.cur.Value)
	}
	diagram.Width = width
	p.nextToken()

	if p.cur.Type != TokenAtom {
		return fmt.Errorf("expected height, got %s", p.cur.Value)
	}
	height, err := strconv.Atoi(p.cur.Value)
	if err != nil {
		return fmt.Errorf("invalid height: %s", p.cur.Value)
	}
	diagram.Height = height
	p.nextToken()

	if p.cur.Type != TokenRParen {
		return fmt.Errorf("expected ')', got %s", p.cur.Value)
	}
	p.nextToken()
	return nil
}

func (p *Parser) parseNodeStyle(diagram *Diagram) error {
	p.nextToken() // consume 'node-style'

	for p.cur.Type != TokenRParen {
		if p.cur.Type != TokenKeyword {
			return fmt.Errorf("expected keyword, got %s", p.cur.Value)
		}
		key := p.cur.Value
		p.nextToken()

		if p.cur.Type != TokenString && p.cur.Type != TokenAtom {
			return fmt.Errorf("expected value, got %s", p.cur.Value)
		}
		value := p.cur.Value
		diagram.NodeStyle[key] = value
		p.nextToken()
	}
	p.nextToken() // consume ')'
	return nil
}

func (p *Parser) parseEdgeStyle(diagram *Diagram) error {
	p.nextToken() // consume 'edge-style'

	for p.cur.Type != TokenRParen {
		if p.cur.Type != TokenKeyword {
			return fmt.Errorf("expected keyword, got %s", p.cur.Value)
		}
		key := p.cur.Value
		p.nextToken()

		if p.cur.Type != TokenString && p.cur.Type != TokenAtom {
			return fmt.Errorf("expected value, got %s", p.cur.Value)
		}
		value := p.cur.Value
		diagram.EdgeStyle[key] = value
		p.nextToken()
	}
	p.nextToken() // consume ')'
	return nil
}

func (p *Parser) parseNodes(diagram *Diagram) error {
	p.nextToken() // consume 'nodes'

	for p.cur.Type != TokenRParen {
		if p.cur.Type != TokenLParen {
			return fmt.Errorf("expected '(', got %s", p.cur.Value)
		}
		p.nextToken()

		if p.cur.Type != TokenAtom || p.cur.Value != "id" {
			return fmt.Errorf("expected 'id', got %s", p.cur.Value)
		}
		p.nextToken()

		if p.cur.Type != TokenString && p.cur.Type != TokenAtom {
			return fmt.Errorf("expected node id, got %s", p.cur.Value)
		}
		nodeID := p.cur.Value
		p.nextToken()

		node := Node{
			ID:         nodeID,
			Label:      nodeID,
			Attributes: make(map[string]string),
		}

		// Parse attributes
		for p.cur.Type != TokenRParen {
			if p.cur.Type != TokenKeyword {
				return fmt.Errorf("expected keyword, got %s", p.cur.Value)
			}
			key := p.cur.Value
			p.nextToken()

			if p.cur.Type != TokenString && p.cur.Type != TokenAtom {
				return fmt.Errorf("expected value, got %s", p.cur.Value)
			}
			value := p.cur.Value

			if key == "label" {
				node.Label = value
			} else {
				node.Attributes[key] = value
			}
			p.nextToken()
		}

		diagram.Nodes = append(diagram.Nodes, node)
		p.nextToken() // consume ')'
	}
	p.nextToken() // consume ')'
	return nil
}

func (p *Parser) parseEdges(diagram *Diagram) error {
	p.nextToken() // consume 'edges'

	for p.cur.Type != TokenRParen {
		if p.cur.Type != TokenLParen {
			return fmt.Errorf("expected '(', got %s", p.cur.Value)
		}
		p.nextToken()

		if p.cur.Type != TokenString && p.cur.Type != TokenAtom {
			return fmt.Errorf("expected from node, got %s", p.cur.Value)
		}
		from := p.cur.Value
		p.nextToken()

		if p.cur.Type != TokenString && p.cur.Type != TokenAtom {
			return fmt.Errorf("expected to node, got %s", p.cur.Value)
		}
		to := p.cur.Value
		p.nextToken()

		edge := Edge{
			From:       from,
			To:         to,
			Label:      "",
			Attributes: make(map[string]string),
		}

		// Parse attributes
		for p.cur.Type != TokenRParen {
			if p.cur.Type != TokenKeyword {
				return fmt.Errorf("expected keyword, got %s", p.cur.Value)
			}
			key := p.cur.Value
			p.nextToken()

			if p.cur.Type != TokenString && p.cur.Type != TokenAtom {
				return fmt.Errorf("expected value, got %s", p.cur.Value)
			}
			value := p.cur.Value

			if key == "label" {
				edge.Label = value
			} else {
				edge.Attributes[key] = value
			}
			p.nextToken()
		}

		diagram.Edges = append(diagram.Edges, edge)
		p.nextToken() // consume ')'
	}
	p.nextToken() // consume ')'
	return nil
}
