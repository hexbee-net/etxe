package etx

import (
	"strings"
)

// Attribute is a key+value attribute.
type Attribute struct {
	ASTNode

	Comment *Comment `parser:"[ @@ ]"    json:"comment,omitempty"`
	Key     string   `parser:"@Ident"    json:"key"`
	Value   *Expr    `parser:"['=' @@ ]" json:"value,omitempty"`
}

func (n *Attribute) Clone() *Attribute {
	if n == nil {
		return nil
	}

	return &Attribute{
		ASTNode: n.ASTNode.Clone(),
		Comment: n.Comment.Clone(),
		Key:     n.Key,
		Value:   n.Value.Clone(),
	}
}

func (n *Attribute) Children() (children []Node) {
	if n.Comment != nil {
		children = append(children, n.Comment)
	}

	if n.Value != nil {
		children = append(children, n.Value)
	}

	return
}

func (n Attribute) String() string {
	var sb strings.Builder

	if n.Comment != nil {
		sb.WriteString(n.Comment.String())
	}

	if n.Value != nil {
		mustFprintf(&sb, "%v: %v", n.Key, n.Value)
	} else {
		sb.WriteString(n.Key)
	}

	return sb.String()
}
