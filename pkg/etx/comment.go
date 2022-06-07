package etx

import "strings"

// Comment is comment block.
type Comment struct {
	ASTNode

	SingleLine []string `parser:"(   (@SingleLineComment [ NewLine ])+" json:"single_line,omitempty"`
	Multiline  string   `parser:"  | @(MultilineComment [ NewLine ])   )" json:"multiline,omitempty"`
}

func (c *Comment) Clone() *Comment {
	if c == nil {
		return nil
	}

	return &Comment{
		ASTNode:    c.ASTNode.Clone(),
		SingleLine: cloneStrings(c.SingleLine),
		Multiline:  c.Multiline,
	}
}

func (c *Comment) Children() (children []Node) {
	return
}

func (c Comment) String() string {
	if c.Multiline != "" {
		return c.Multiline
	}

	var sb strings.Builder
	for _, item := range c.SingleLine {
		sb.WriteString(item)
		sb.WriteString("\n")
	}

	return sb.String()
}
