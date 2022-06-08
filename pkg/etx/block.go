package etx

import (
	"strings"
)

// Block represents an optionally labeled block.
type Block struct {
	ASTNode

	Name   string       `parser:"@Ident"                                    json:"name"`
	Labels []string     `parser:"((String @Char StringEnd) | @Ident)* '{' [NewLine] "  json:"labels,omitempty"`
	Body   []*BlockItem `parser:"@@* '}'"            json:"body"`
}

func (n *Block) Clone() *Block {
	if n == nil {
		return nil
	}

	return &Block{
		ASTNode: n.ASTNode.Clone(),
		Name:    n.Name,
		Labels:  cloneStrings(n.Labels),
		Body:    cloneCollection(n.Body),
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

	EmptyLine string     `parser:"(   @NewLine+       " json:"empty_line,omitempty"`
	Block     *Block     `parser:"  | (@@ [NewLine])  " json:"block,omitempty"`
	Attribute *Attribute `parser:"  | (@@ [NewLine])  " json:"attribute,omitempty"`
	Comment   *Comment   `parser:"  | @@             )" json:"comment,omitempty"`
}

func (n *BlockItem) Clone() *BlockItem {
	if n == nil {
		return nil
	}

	return &BlockItem{
		ASTNode:   n.ASTNode.Clone(),
		Block:     n.Block.Clone(),
		Attribute: n.Attribute.Clone(),
		Comment:   n.Comment.Clone(),
		EmptyLine: n.EmptyLine,
	}
}

func (n *BlockItem) Children() (children []Node) {
	if n.Comment != nil {
		children = append(children, n.Comment)
	}

	if n.Block != nil {
		children = append(children, n.Block)
	}

	if n.Attribute != nil {
		children = append(children, n.Attribute)
	}

	return
}

func (n BlockItem) String() string {
	var sb strings.Builder

	if n.Comment != nil {
		sb.WriteString(n.Comment.String())
	}

	switch {
	case n.Block != nil:
		sb.WriteString(n.Block.String())
	case n.Attribute != nil:
		sb.WriteString(n.Attribute.String())
	default:
	}

	return sb.String()
}
