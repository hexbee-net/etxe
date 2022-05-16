package etx

import "github.com/alecthomas/participle/v2/lexer"

// TODO: Since `val`s are immutable, do we really need `const`?
// The one use case would be that const shouldn't accept an expr, only a value.

// Decl is an `input`, `output`, `const` or `val` short form declaration.
type Decl struct {
	Pos    lexer.Position `parser:"" json:"-"`
	Parent Node           `parser:"" json:"-"`

	Comments []string `parser:"@Comment*" json:"comments,omitempty"`

	DeclType string         `parser:"@(Input | Output | Const | Val) Whitespace " json:"decl_type"`
	Label    string         `parser:"@Ident                                     " json:"label"`
	Type     *ParameterType `parser:"(Whitespace? ':' Whitespace? @@)?          " json:"type"`
	Value    *Expr          `parser:"(Whitespace? '=' Whitespace? @@)?          " json:"value"`
}
