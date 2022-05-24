package etx

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/repr"
)

// Position in source file.
type Position = lexer.Position

// Node is the interface implemented by all AST nodes.
type Node interface {
	Children() (children []Node)
}

// /////////////////////////////////////

type ASTNode struct {
	Pos    Position `parser:"" json:"-"`
	Parent Node     `parser:"" json:"-"`
}

func (n ASTNode) Clone() ASTNode {
	out := ASTNode{
		Pos: n.Pos,
	}

	return out
}

// /////////////////////////////////////

type CommentNode struct {
	Comments []string `parser:"@Comment*" json:"comments,omitempty"`
}

func (n CommentNode) Clone() CommentNode {
	return CommentNode{
		Comments: cloneStrings(n.Comments),
	}
}

// /////////////////////////////////////

// AST for ETX files.
type AST struct {
	Items            []*RootItem `parser:"(@@ (NewLine @@)* )*" json:"items,omitempty"`
	TrailingComments []string    `parser:"@Comment*"            json:"trailing_comments,omitempty"`
}

func (n *AST) Clone() *AST {
	if n == nil {
		return nil
	}

	out := &AST{
		Items:            cloneCollection(n.Items),
		TrailingComments: cloneStrings(n.TrailingComments),
	}

	return out
}

func (n *AST) Children() (children []Node) {
	for _, item := range n.Items {
		children = append(children, item)
	}

	return
}

func (n AST) String() string {
	var sb strings.Builder

	// TODO: group declarations by kind?
	for _, item := range n.Items {
		mustFprintf(&sb, "%s\n\n", item)
	}

	return sb.String()
}

// /////////////////////////////////////

// RootItem at the top-level of a file.
type RootItem struct {
	ASTNode

	Attribute *Attribute `parser:"(   @@  " json:"attribute,omitempty"`
	Decl      *Decl      `parser:"  | @@  " json:"decl,omitempty"`
	Func      *Func      `parser:"  | @@  " json:"func,omitempty"`
	Type      *Type      `parser:"  | @@  " json:"type,omitempty"`
	Block     *Block     `parser:"  | @@ )" json:"block,omitempty"`
}

func (n *RootItem) Clone() *RootItem {
	if n == nil {
		return nil
	}

	return &RootItem{
		ASTNode:   n.ASTNode.Clone(),
		Attribute: n.Attribute.Clone(),
		Decl:      n.Decl.Clone(),
		Func:      n.Func.Clone(),
		Type:      n.Type.Clone(),
		Block:     n.Block.Clone(),
	}
}

func (n *RootItem) Children() (children []Node) {
	switch {
	case n.Attribute != nil:
		children = append(children, n.Attribute)
	case n.Decl != nil:
		children = append(children, n.Decl)
	case n.Func != nil:
		children = append(children, n.Func)
	case n.Type != nil:
		children = append(children, n.Type)
	case n.Block != nil:
		children = append(children, n.Block)
	}

	return
}

func (n RootItem) String() string {
	switch {
	case n.Attribute != nil:
		return n.Attribute.String()
	case n.Decl != nil:
		return n.Decl.String()
	case n.Func != nil:
		return n.Func.String()
	case n.Type != nil:
		return n.Type.String()
	case n.Block != nil:
		return n.Block.String()
	default:
		panic(repr.String(n, repr.Hide(Position{})))
	}
}

// /////////////////////////////////////

type ParameterType struct {
	ASTNode

	Ident *Ident         `parser:"(   @@   " json:"ident"`
	Func  *FuncSignature `parser:"  | @@  )" json:"func"`
}

func (n *ParameterType) Clone() *ParameterType {
	if n == nil {
		return nil
	}

	return &ParameterType{
		ASTNode: n.ASTNode.Clone(),
		Ident:   n.Ident.Clone(),
		Func:    n.Func.Clone(),
	}
}

func (n *ParameterType) Children() (children []Node) {
	if n.Ident != nil {
		children = append(children, n.Ident)
	}

	if n.Func != nil {
		children = append(children, n.Func)
	}

	return
}

func (n ParameterType) String() string {
	switch {
	case n.Ident != nil:
		return n.Ident.String()
	case n.Func != nil:
		return n.Func.String()
	default:
		panic(repr.String(n, repr.Hide(Position{})))
	}
}

// /////////////////////////////////////

type FuncSignature struct {
	ASTNode

	Parameters []*ParameterType `parser:"'(' ( @@ (','  @@)* )? ')' OpLambdaDef" json:"parameters,omitempty"`
	Return     ParameterType    `parser:"@@"                                     json:"return,omitempty"`
}

func (n *FuncSignature) Clone() *FuncSignature {
	if n == nil {
		return nil
	}

	return &FuncSignature{
		ASTNode:    n.ASTNode.Clone(),
		Parameters: cloneCollection(n.Parameters),
		Return:     *n.Return.Clone(),
	}
}

func (n *FuncSignature) Children() (children []Node) {
	for _, item := range n.Parameters {
		children = append(children, item)
	}

	children = append(children, &n.Return)

	return
}

func (n FuncSignature) String() string {
	var sb strings.Builder

	sb.WriteString("(")

	params := make([]string, 0, len(n.Parameters))
	for _, item := range n.Parameters {
		params = append(params, item.String())
	}

	mustFprintf(&sb, "(%v) %v %v", strings.Join(params, ", "), OpLambdaDef, n.Return)

	return sb.String()
}
