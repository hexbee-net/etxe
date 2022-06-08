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
	Node() *ASTNode
}

// /////////////////////////////////////

// ASTNode is a node in the AST.
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

func (n *ASTNode) Node() *ASTNode {
	return n
}

// /////////////////////////////////////

// AST for ETX files.
type AST struct {
	Items []*RootItem `parser:"@@*" json:"items,omitempty"`
}

func (n *AST) Clone() *AST {
	if n == nil {
		return nil
	}

	out := &AST{
		Items: cloneCollection(n.Items),
	}

	return out
}

func (n *AST) Children() (children []Node) {
	for _, item := range n.Items {
		children = append(children, item)
	}

	return
}

func (n AST) FormattedString() string {
	if len(n.Items) == 0 {
		return ""
	}

	items := make([]string, 0, len(n.Items))
	for _, item := range n.Items {
		items = append(items, item.FormattedString())
	}

	return strings.Join(items, "\n\n") + "\n"
}

// UpdateParentRefs recursively updates the AST nodes parent references.
func (n *AST) UpdateParentRefs() {
	for _, c := range n.Children() {
		updateParentRefs(nil, c)
	}
}

// /////////////////////////////////////

// RootItem at the top-level of a file.
type RootItem struct {
	ASTNode

	EmptyLine string     `parser:"(   @NewLine+     " json:"empty_line,omitempty"`
	Decl      *Decl      `parser:"  | @@ [NewLine]  " json:"decl,omitempty"`
	Func      *Func      `parser:"  | @@ [NewLine]  " json:"func,omitempty"`
	Type      *Type      `parser:"  | @@ [NewLine]  " json:"type,omitempty"`
	Block     *Block     `parser:"  | @@ [NewLine]  " json:"block,omitempty"`
	Attribute *Attribute `parser:"  | @@ [NewLine]  " json:"attribute,omitempty"`
	Comment   *Comment   `parser:"  | @@           )" json:"comment,omitempty"`
}

func (n *RootItem) Clone() *RootItem {
	if n == nil {
		return nil
	}

	return &RootItem{
		ASTNode:   n.ASTNode.Clone(),
		Decl:      n.Decl.Clone(),
		Func:      n.Func.Clone(),
		Type:      n.Type.Clone(),
		Block:     n.Block.Clone(),
		Attribute: n.Attribute.Clone(),
		Comment:   n.Comment.Clone(),
		EmptyLine: n.EmptyLine,
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
	case n.Comment != nil:
		children = append(children, n.Comment)

	}

	return
}

func (n RootItem) FormattedString() string {
	switch {
	case n.Decl != nil:
		return n.Decl.FormattedString()
	case n.Func != nil:
		return n.Func.FormattedString()
	case n.Type != nil:
		return n.Type.FormattedString()
	case n.Block != nil:
		return n.Block.FormattedString()
	case n.Attribute != nil:
		return n.Attribute.FormattedString()
	case n.Comment != nil:
		return n.Comment.FormattedString()
	case n.EmptyLine != "":
		return n.EmptyLine
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

func (n ParameterType) FormattedString() string {
	switch {
	case n.Ident != nil:
		return n.Ident.FormattedString()
	case n.Func != nil:
		return n.Func.FormattedString()
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

func (n FuncSignature) FormattedString() string {
	var sb strings.Builder

	params := make([]string, 0, len(n.Parameters))
	for _, item := range n.Parameters {
		params = append(params, item.FormattedString())
	}

	mustFprintf(&sb, "(%s) %s %s", strings.Join(params, ", "), OpLambdaDef, n.Return.FormattedString())

	return sb.String()
}
