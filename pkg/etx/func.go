package etx

import (
	"fmt"
	"strings"
)

type Func struct {
	ASTNode

	Label      string           `parser:"Func @Ident "                                            json:"label"`
	Parameters []*FuncParameter `parser:"'(' [ @@ (',' @@)* ] ')'"                                json:"parameters,omitempty"`
	Return     []*ParameterType `parser:"('(' [ @@ (',' @@)* ] ')' | @@)?"                        json:"return,omitempty"`
	Body       []*FuncStatement `parser:"[ NewLine+ ] '{' [ NewLine+ ] @@ * '}' " json:"body,omitempty"`
}

func (n *Func) Clone() *Func {
	if n == nil {
		return nil
	}

	return &Func{
		ASTNode:    n.ASTNode.Clone(),
		Label:      n.Label,
		Parameters: cloneCollection(n.Parameters),
		Return:     cloneCollection(n.Return),
		Body:       cloneCollection(n.Body),
	}
}

func (n *Func) Children() (children []Node) {
	for _, item := range n.Parameters {
		children = append(children, item)
	}

	for _, item := range n.Return {
		children = append(children, item)
	}

	for _, item := range n.Body {
		children = append(children, item)
	}

	return
}

func (n Func) FormattedString() string {
	var sb strings.Builder

	params := make([]string, 0, len(n.Parameters))
	for _, p := range n.Parameters {
		params = append(params, p.FormattedString())
	}

	mustFprintf(&sb, "def %s(%s)", n.Label, strings.Join(params, ", "))

	switch l := len(n.Return); {
	case l == 1:
		mustFprintf(&sb, " %s", n.Return[0].FormattedString())
	case l > 1:
		rets := make([]string, 0, l)
		for _, item := range n.Return {
			rets = append(rets, item.FormattedString())
		}

		mustFprintf(&sb, " (%s)", strings.Join(rets, ", "))
	}

	if len(n.Body) != 0 {
		sb.WriteString(" {\n")

		for _, b := range n.Body {
			sb.WriteString(indent(b.FormattedString(), indentationChar))
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

func (n FuncParameter) FormattedString() string {
	if n.Type != nil {
		return fmt.Sprintf("%s: %s", n.Label, n.Type.FormattedString())
	}

	return n.Label
}

// /////////////////////////////////////

type FuncStatement struct {
	ASTNode

	EmptyLine string    `parser:"(   @NewLine+     " json:"empty_line,omitempty"`
	Comment   *Comment  `parser:"  | @@            "        json:"comment,omitempty"`
	Decl      *FuncDecl `parser:"  | @@ [NewLine]  "        json:"decl,omitempty"`
	Expr      *Expr     `parser:"  | @@ [NewLine] )"        json:"expr,omitempty"`
}

func (n *FuncStatement) Clone() *FuncStatement {
	if n == nil {
		return nil
	}

	return &FuncStatement{
		ASTNode:   n.ASTNode.Clone(),
		Comment:   n.Comment.Clone(),
		Decl:      n.Decl.Clone(),
		Expr:      n.Expr.Clone(),
		EmptyLine: n.EmptyLine,
	}
}

func (n *FuncStatement) Children() (children []Node) {
	if n.Comment != nil {
		children = append(children, n.Comment)
	}

	if n.Decl != nil {
		children = append(children, n.Decl)
	}

	if n.Expr != nil {
		children = append(children, n.Expr)
	}

	return
}

func (n FuncStatement) FormattedString() string {
	switch {
	case n.Comment != nil:
		return n.Comment.FormattedString()
	case n.Decl != nil:
		return n.Decl.FormattedString()
	case n.Expr != nil:
		return n.Expr.FormattedString()
	case n.EmptyLine != "":
		return n.EmptyLine
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

func (n FuncDecl) FormattedString() string {
	var sb strings.Builder

	if n.Label == "" {
		return ""
	}

	mustFprintf(&sb, "%s %s", n.DeclType, n.Label)

	if n.Type != nil {
		mustFprintf(&sb, ": %s", n.Type.FormattedString())
	}

	if n.Value != nil {
		mustFprintf(&sb, " = %s", n.Value.FormattedString())
	}

	return sb.String()
}
