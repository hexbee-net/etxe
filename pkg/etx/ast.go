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
	Pos lexer.Position `parser:"" json:"-"`

	Items            []*Item  `parser:"@@*" json:"items,omitempty"`
	TrailingComments []string `parser:"@Comment*" json:"trailing_comments,omitempty"`
	Schema           bool     `parser:"" json:"schema,omitempty"`
}

func (a *AST) children() (children []Node) {
	// TODO implement me
	panic("implement me")
}

// Item at the top-level of a file.
type Item struct {
	Pos             lexer.Position `parser:"" json:"-"`
	Parent          Node           `parser:"" json:"-"`
	RecursiveSchema bool           `parser:"" json:"-"`

	Block       *Block     `parser:"(   @@  " json:"block,omitempty"`
	Attribute   *Attribute `parser:"  | @@  " json:"attribute,omitempty"`
	Declaration *Decl      `parser:"  | @@  " json:"declaration,omitempty"`
	Function    *Func      `parser:"  | @@  " json:"function,omitempty"`
	Type        *Type      `parser:"  | @@ )" json:"type,omitempty"`
}

// Block represents an optionally labeled block.
type Block struct {
	Pos    lexer.Position `parser:"" json:"-"`
	Parent Node           `parser:"" json:"-"`

	Comments []string `parser:"@Comment*" json:"comments,omitempty"`

	Name   string       `parser:"@Ident" json:"name"`
	Labels []string     `parser:"@( Ident | String )*" json:"labels,omitempty"`
	Body   []*BlockItem `parser:"'{' @@*" json:"body"`

	TrailingComments []string `parser:"@Comment* '}'" json:"trailing_comments,omitempty"`

	// The block can be repeated. This is surfaced in schemas.
	Repeated bool `parser:"" json:"repeated,omitempty"`
}

// BlockItem in a block.
type BlockItem struct {
	Pos             lexer.Position `parser:"" json:"-"`
	Parent          Node           `parser:"" json:"-"`
	RecursiveSchema bool           `parser:"" json:"-"`

	Block     *Block     `parser:"(   @@" json:"block,omitempty"`
	Attribute *Attribute `parser:"  | @@" json:"attribute,omitempty"`
}

// Attribute is a key+value attribute.
type Attribute struct {
	Pos    lexer.Position `parser:"" json:"-"`
	Parent Node           `parser:"" json:"-"`

	Comments []string `parser:"@Comment*" json:"comments,omitempty"`

	Key   string `parser:"@Ident ['='" json:"key"`
	Value *Value `parser:"@@]" json:"value"`

	// This will be populated during unmarshalling.
	Default *Value `parser:"" json:"default,omitempty"`

	// This will be parsed from the enum tag and will be helping the validation during unmarshalling.
	Enum []*Value `parser:"" json:"enum,omitempty"`

	// Set for schemas when the attribute is optional.
	Optional bool `parser:"" json:"optional,omitempty"`
}

// Decl is an `input`, `output`, `const` or `val` short form declaration.
type Decl struct {
	Pos    lexer.Position `parser:"" json:"-"`
	Parent Node           `parser:"" json:"-"`

	Comments []string `parser:"@Comment*" json:"comments,omitempty"`

	DeclType string `parser:"@(Input | Output | Const | Val) Whitespace" json:"decl_type"`
	Label    string `parser:"@Ident" json:"label"`
	Type     string `parser:"(Whitespace? ':' Whitespace? @Ident)?" json:"type"` // TODO: ParameterType
	Value    *Value `parser:"(Whitespace? '=' Whitespace? @@)?" json:"value"`
}

type ParameterType struct {
	Ident string         `parser:"(   @Ident  " json:"ident"`
	Func  *FuncSignature `parser:"  | @@     )" json:"func"`
}

type FuncSignature struct {
	Parameters []*ParameterType `parser:"'(' (@@ ( Whitespace? ',' Whitespace? @@ )*)? ')' Whitespace? LambdaDef" json:"parameters,omitempty"`
	Return     *ParameterType   `parser:"Whitespace? @@" json:"return,omitempty"`
}

type Lambda struct {
	Parameters []*LambdaParameter `parser:"'(' (@@ ( Whitespace? ',' Whitespace? @@ )*)? ')' Whitespace? Lambda Whitespace?" json:"parameters,omitempty"`
	Expr       *Expr              `parser:"@@"`
}

type LambdaParameter struct {
	Label string         `parser:"@Ident" json:"label"`
	Type  *ParameterType `parser:"(Whitespace? ':' Whitespace? @@)?" json:"type"`
}

type Type struct {
}
