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
}

func (a *AST) children() (children []Node) {
	// TODO implement me.
	panic("implement me")
}
