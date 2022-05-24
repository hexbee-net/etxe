package etx

import "fmt"

// Attribute is a key+value attribute.
type Attribute struct {
	ASTNode

	Key   string `parser:"@Ident"    json:"key"`
	Value *Expr  `parser:"['=' @@ ]" json:"value,omitempty"`
}

func (n *Attribute) Clone() *Attribute {
	if n == nil {
		return nil
	}

	return &Attribute{
		ASTNode: n.ASTNode.Clone(),
		Key:     n.Key,
		Value:   n.Value.Clone(),
	}
}

func (n *Attribute) Children() (children []Node) {
	if n.Value != nil {
		children = append(children, n.Value)
	}

	return
}

func (n Attribute) String() string {
	if n.Value != nil {
		return fmt.Sprintf("%v: %v", n.Key, n.Value)
	}

	return n.Key
}
