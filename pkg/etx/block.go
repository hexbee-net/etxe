package etx

import (
	"strings"
)

// Block represents an optionally labeled block.
type Block struct {
	ASTNode

	Comments         []string     `parser:"(@Comment [ NewLine ])*"                  json:"comments,omitempty"`
	Name             string       `parser:"@Ident"                                   json:"name"`
	Labels           []string     `parser:"((String @Char StringEnd) | @Ident)* '{'" json:"labels,omitempty"`
	Body             []*BlockItem `parser:"(NewLine @@)*"                            json:"body"`
	TrailingComments []string     `parser:"(@Comment [ NewLine ])* [ NewLine ] '}' " json:"trailing_comments,omitempty"`
}

func (n *Block) Clone() *Block {
	if n == nil {
		return nil
	}

	return &Block{
		ASTNode:          n.ASTNode.Clone(),
		Comments:         cloneStrings(n.Comments),
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
	if n.Name == "" {
		return ""
	}

	var sb strings.Builder

	sb.WriteString(n.Name)

	for _, item := range n.Labels {
		mustFprintf(&sb, ` "%v"`, item)
	}

	if len(n.Body) != 0 {
		sb.WriteString(" {\n")

		for _, item := range n.Body {
			sb.WriteString(indent(item.String(), indentationChar))
		}

		sb.WriteString("\n}")
	} else {
		sb.WriteString(" {}")
	}

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
		return ""
	}
}
