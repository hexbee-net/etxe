package etx

import (
	"strings"
)

type Ident struct {
	ASTNode

	Parts []string `parser:"@Ident [('.' @Ident)*]" json:"parts"`
}

func (i Ident) String() string {
	return strings.Join(i.Parts, ".")
}

func (i *Ident) Children() (children []Node) {
	return
}

// Clone the AST node.
func (i *Ident) Clone() *Ident {
	if i == nil {
		return nil
	}

	return &Ident{
		Parts: i.Parts,
	}
}
