package etx

import (
	"fmt"
	"strings"
)

type Lambda struct {
	ASTNode

	Comment    *Comment           `parser:"[ @@ ]"                            json:"comment,omitempty"`
	Parameters []*LambdaParameter `parser:"'(' [ @@ (',' @@)* ] ')' OpLambda" json:"parameters"`
	Expr       Expr               `parser:"@@"                                json:"expr"`
}

func (n *Lambda) Clone() *Lambda {
	if n == nil {
		return nil
	}

	return &Lambda{
		ASTNode:    n.ASTNode.Clone(),
		Comment:    n.Comment.Clone(),
		Parameters: cloneCollection(n.Parameters),
		Expr:       *n.Expr.Clone(),
	}
}

func (n *Lambda) Children() (children []Node) {
	if n.Comment != nil {
		children = append(children, n.Comment)
	}

	for _, item := range n.Parameters {
		children = append(children, item)
	}

	children = append(children, &n.Expr)

	return
}

func (n Lambda) FormattedString() string {
	var sb strings.Builder

	if n.Comment != nil {
		sb.WriteString(n.Comment.FormattedString())
	}

	params := make([]string, 0, len(n.Parameters))
	for _, p := range n.Parameters {
		params = append(params, p.FormattedString())
	}

	mustFprintf(&sb, "(%s) %s %s", strings.Join(params, ", "), OpLambda, n.Expr.FormattedString())

	return sb.String()
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

func (n LambdaParameter) FormattedString() string {
	if n.Type != nil {
		return fmt.Sprintf("%s: %s", n.Label, n.Type.FormattedString())
	}

	return n.Label
}
