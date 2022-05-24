package etx

import (
	"strings"

	"github.com/alecthomas/repr"
)

// Block represents an optionally labeled block.
type Block struct {
	ASTNode

	Name             string       `parser:"@Ident"               json:"name"`
	Labels           []string     `parser:"@( Ident | String )*" json:"labels,omitempty"` // TODO: this is not gonna work
	Body             []*BlockItem `parser:"'{' @@*       "       json:"body"`
	TrailingComments []string     `parser:"@Comment* '}' "       json:"trailing_comments,omitempty"`
}

func (n *Block) Clone() *Block {
	if n == nil {
		return nil
	}

	return &Block{
		ASTNode:          n.ASTNode.Clone(),
		Name:             n.Name,
		Labels:           cloneStrings(n.Labels),
		Body:             cloneCollection(n.Body),
		TrailingComments: cloneStrings(n.TrailingComments),
	}
}

func (n *Block) Children() (children []Node) {
	for _, item := range n.Body {
		children = append(children, item)
	}

	return
}

func (n Block) String() string {
	var sb strings.Builder

	sb.WriteString(n.Name)

	for _, item := range n.Labels {
		mustFprintf(&sb, ` "%v"`, item)
	}

	sb.WriteString(" {\n")

	for _, item := range n.Body {
		sb.WriteString(indent(item.String(), indentationChar))
	}

	sb.WriteString("\n}")

	return sb.String()
}

// /////////////////////////////////////

// BlockItem in a block.
type BlockItem struct {
	ASTNode

	Block     *Block     `parser:"(   @@  " json:"block,omitempty"`
	Attribute *Attribute `parser:"  | @@ )" json:"attribute,omitempty"`
}

func (n *BlockItem) Clone() *BlockItem {
	if n == nil {
		return nil
	}

	return &BlockItem{
		ASTNode:   n.ASTNode.Clone(),
		Block:     n.Block.Clone(),
		Attribute: n.Attribute.Clone(),
	}
}

func (n *BlockItem) Children() (children []Node) {
	if n.Block != nil {
		children = append(children, n.Block)
	}

	if n.Attribute != nil {
		children = append(children, n.Attribute)
	}

	return
}

func (n BlockItem) String() string {
	switch {
	case n.Block != nil:
		return n.Block.String()
	case n.Attribute != nil:
		return n.Attribute.String()
	default:
		panic(repr.String(n, repr.Hide(Position{})))
	}
}
