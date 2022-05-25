package etx

import (
	"fmt"
	"strings"
)

type Lambda struct {
	ASTNode

	Parameters []*LambdaParameter `parser:"'(' [ @@ (',' @@)* ] ')' OpLambda" json:"parameters"`
	Expr       Expr               `parser:"@@"                                json:"expr"`
}

func (n *Lambda) Clone() *Lambda {
	if n == nil {
		return nil
	}

	return &Lambda{
		ASTNode:    n.ASTNode.Clone(),
		Parameters: cloneCollection(n.Parameters),
		Expr:       *n.Expr.Clone(),
	}
}

func (n *Lambda) Children() (children []Node) {
	for _, item := range n.Parameters {
		children = append(children, item)
	}

	children = append(children, &n.Expr)

	return
}

func (n Lambda) String() string {
	params := make([]string, 0, len(n.Parameters))
	for _, p := range n.Parameters {
		params = append(params, p.String())
	}

	return fmt.Sprintf("(%s) %s %s", strings.Join(params, ", "), OpLambda, n.Expr.String())
}

// /////////////////////////////////////

type LambdaParameter struct {
	ASTNode

	Label string         `parser:"@Ident"     json:"label"`
	Type  *ParameterType `parser:"[ ':' @@ ]" json:"type,omitempty"`
}

func (n *LambdaParameter) Clone() *LambdaParameter {
	if n == nil {
		return nil
	}

	return &LambdaParameter{
		ASTNode: n.ASTNode.Clone(),
		Label:   n.Label,
		Type:    n.Type.Clone(),
	}
}

func (n *LambdaParameter) Children() (children []Node) {
	if n.Type != nil {
		children = append(children, n.Type)
	}

	return
}

func (n LambdaParameter) String() string {
	if n.Type != nil {
		return fmt.Sprintf("%v: %v", n.Label, n.Type)
	}

	return n.Label
}
