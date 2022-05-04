package etx

import "github.com/alecthomas/participle/lexer"

// Position in source file.
type Position = lexer.Position

// Node is the interface implemented by all AST nodes.
type Node interface {
	children() (children []Node)
}

// AST for HCL.
type AST struct {
	Pos lexer.Position `parser:"" json:"-"`

	// Entries          []*Entry `parser:"@@*" json:"entries,omitempty"`
	// TrailingComments []string `parser:"@Comment*" json:"trailing_comments,omitempty"`
	// Schema           bool     `parser:"" json:"schema,omitempty"`
}

func (a *AST) children() (children []Node) {
	// TODO implement me.
	panic("implement me")
}
