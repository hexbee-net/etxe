package etx

import (
	"strings"
)

// TODO: Since `val`s are immutable, do we really need `const`?
// The one use case would be that const shouldn't accept an expr, only a value.

// Decl is an `input`, `output`, `const` or `val` short form declaration.
type Decl struct {
	ASTNode

	DeclType string         `parser:"@(input | Output | Const | Val)" json:"decl_type"`
	Label    string         `parser:"@Ident"                          json:"label"`
	Type     *ParameterType `parser:"[':' @@]"                        json:"type,omitempty"`
	Value    *Expr          `parser:"['=' @@]"                        json:"value,omitempty"`
}

func (n *Decl) Clone() *Decl {
	if n == nil {
		return nil
	}

	return &Decl{
		ASTNode:  n.ASTNode.Clone(),
		DeclType: n.DeclType,
		Label:    n.Label,
		Type:     n.Type.Clone(),
		Value:    n.Value.Clone(),
	}
}

func (n *Decl) Children() (children []Node) {
	if n.Type != nil {
		children = append(children, n.Type)
	}

	if n.Value != nil {
		children = append(children, n.Value)
	}

	return
}

func (n Decl) String() string {
	var sb strings.Builder

	mustFprintf(&sb, "%v %v", n.DeclType, n.Label)

	if n.Type != nil {
		mustFprintf(&sb, ": %v", n.Type)
	}

	if n.Value != nil {
		mustFprintf(&sb, "= %v", n.Value)
	}

	return sb.String()
}
