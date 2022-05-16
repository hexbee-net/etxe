package etx

import "github.com/alecthomas/participle/v2/lexer"

// Position in source file.
type Position = lexer.Position

// Node is the interface implemented by all AST nodes.
type Node interface {
	children() (children []Node)
}

// AST for HCL.
type AST struct {
	Items            []*Item  `parser:"(@@ (NewLine @@)* )*" json:"items,omitempty"`
	TrailingComments []string `parser:"@Comment*"            json:"trailing_comments,omitempty"`
}

func (a *AST) children() (children []Node) {
	// TODO implement me
	panic("implement me")
}

// Item at the top-level of a file.
type Item struct {
	Parent Node `parser:"" json:"-"`

	Attribute *Attribute `parser:"(   @@  " json:"attribute,omitempty"`
	Decl      *Decl      `parser:"  | @@  " json:"decl,omitempty"`
	Func      *Func      `parser:"  | @@  " json:"func,omitempty"`
	Type      *Type      `parser:"  | @@  " json:"type,omitempty"`
	Block     *Block     `parser:"  | @@ )" json:"block,omitempty"`
}

type ParameterType struct {
	Ident *Ident         `parser:"(   @@   " json:"ident"`
	Func  *FuncSignature `parser:"  | @@  )" json:"func"`
}

type FuncSignature struct {
	Parameters []*ParameterType `parser:"'(' (@@ ( Whitespace? ',' Whitespace? @@ )*)? ')' Whitespace? LambdaDef" json:"parameters,omitempty"`
	Return     *ParameterType   `parser:"Whitespace? @@"                                                          json:"return,omitempty"`
}
