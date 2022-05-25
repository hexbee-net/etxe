package etx

import (
	"fmt"
	"strings"
)

type Func struct {
	ASTNode

	Label      string           `parser:"Func @Ident "                               json:"label"`
	Parameters []*FuncParameter `parser:"'(' [ @@ (',' @@)* ] ')'"                   json:"parameters,omitempty"`
	Return     *ParameterType   `parser:"@@?"                                        json:"return,omitempty"` // TODO: let return several values
	Body       []*FuncStatement `parser:"NewLine? '{' ( NewLine? @@ NewLine? )* '}'" json:"body,omitempty"`
}

func (n *Func) Clone() *Func {
	if n == nil {
		return nil
	}

	return &Func{
		ASTNode:    n.ASTNode.Clone(),
		Label:      n.Label,
		Parameters: cloneCollection(n.Parameters),
		Return:     n.Return.Clone(),
		Body:       cloneCollection(n.Body),
	}
}

func (n *Func) Children() (children []Node) {
	for _, item := range n.Parameters {
		children = append(children, item)
	}

	if n.Return != nil {
		children = append(children, n.Return)
	}

	for _, item := range n.Body {
		children = append(children, item)
	}

	return
}

func (n Func) String() string {
	var sb strings.Builder

	params := make([]string, 0, len(n.Parameters))
	for _, p := range n.Parameters {
		params = append(params, p.String())
	}

	mustFprintf(&sb, "def %v(%v)", n.Label, strings.Join(params, ", "))

	if n.Return != nil {
		mustFprintf(&sb, " %v", n.Return)
	}

	if len(n.Body) != 0 {
		sb.WriteString(" {\n")

		for _, b := range n.Body {
			sb.WriteString(indent(b.String(), indentationChar))
			sb.WriteString("\n")
		}

		sb.WriteString("}")
	} else {
		sb.WriteString(" {}")
	}

	return sb.String()
}

// /////////////////////////////////////

type FuncParameter struct {
	ASTNode

	Label string         `parser:"@Ident"   json:"label"`
	Type  *ParameterType `parser:"[':' @@]" json:"type"`
}

func (n *FuncParameter) Clone() *FuncParameter {
	if n == nil {
		return nil
	}

	return &FuncParameter{
		ASTNode: n.ASTNode.Clone(),
		Label:   n.Label,
		Type:    n.Type.Clone(),
	}
}

func (n *FuncParameter) Children() (children []Node) {
	if n.Type != nil {
		children = append(children, n.Type)
	}

	return
}

func (n FuncParameter) String() string {
	if n.Type != nil {
		return fmt.Sprintf("%v: %v", n.Label, n.Type)
	}

	return n.Label
}

// /////////////////////////////////////

type FuncStatement struct {
	ASTNode
	CommentNode

	Decl *FuncDecl `parser:"(   @@  "  json:"decl,omitempty"`
	Expr *Expr     `parser:"  | @@ )"  json:"expr,omitempty"`
}

func (n *FuncStatement) Clone() *FuncStatement {
	if n == nil {
		return nil
	}

	return &FuncStatement{
		ASTNode:     n.ASTNode.Clone(),
		CommentNode: n.CommentNode.Clone(),
		Decl:        n.Decl.Clone(),
		Expr:        n.Expr.Clone(),
	}
}

func (n *FuncStatement) Children() (children []Node) {
	if n.Decl != nil {
		children = append(children, n.Decl)
	}

	if n.Expr != nil {
		children = append(children, n.Expr)
	}

	return
}

func (n FuncStatement) String() string {
	switch {
	case n.Decl != nil:
		return n.Decl.String()
	case n.Expr != nil:
		return n.Expr.String()
	default:
		return ""
	}
}

// /////////////////////////////////////

type FuncDecl struct {
	ASTNode

	DeclType string         `parser:"@(Const | Val)" json:"decl_type"`
	Label    string         `parser:"@Ident"         json:"label"`
	Type     *ParameterType `parser:"[ ':' @@ ]"     json:"type,omitempty"`
	Value    *Expr          `parser:"[ '=' @@     ]" json:"value,omitempty"`
}

func (n *FuncDecl) Clone() *FuncDecl {
	if n == nil {
		return nil
	}

	return &FuncDecl{
		ASTNode:  n.ASTNode.Clone(),
		DeclType: n.DeclType,
		Label:    n.Label,
		Type:     n.Type,
		Value:    n.Value.Clone(),
	}
}

func (n *FuncDecl) Children() (children []Node) {
	if n.Type != nil {
		children = append(children, n.Type)
	}

	if n.Value != nil {
		children = append(children, n.Value)
	}

	return
}

func (n FuncDecl) String() string {
	var sb strings.Builder

	if n.Label == "" {
		return ""
	}

	mustFprintf(&sb, "%v %v", n.DeclType, n.Label)

	if n.Type != nil {
		mustFprintf(&sb, ": %s", n.Type.String())
	}

	if n.Value != nil {
		mustFprintf(&sb, " = %v", n.Value)
	}

	return sb.String()
}
