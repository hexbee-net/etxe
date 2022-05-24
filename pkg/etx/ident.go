package etx

import (
	"strings"
)

type Ident struct {
	Parts []string `parser:"@Ident [('.' @Ident)*]"`
}

func (i Ident) String() string {
	return strings.Join(i.Parts, ".")
}

func (i *Ident) Children() (children []Node) {
	return
}

// Clone the AST.
func (i *Ident) Clone() *Ident {
	if i == nil {
		return nil
	}

	return &Ident{
		Parts: i.Parts,
	}
}
