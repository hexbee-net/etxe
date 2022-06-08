package etx

import (
	"strings"
)

type Ident struct {
	ASTNode

	Parts []string `parser:"@Ident [('.' @Ident)*]" json:"parts"`
}

// Clone the AST node.
func (i *Ident) Clone() *Ident {
	if i == nil {
		return nil
	}

	return &Ident{
		ASTNode: i.ASTNode.Clone(),
		Parts:   cloneStrings(i.Parts),
	}
}

func (i *Ident) Children() (children []Node) {
	return
}

func (i Ident) FormattedString() string {
	return strings.Join(i.Parts, ".")
}
