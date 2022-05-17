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

func (i *Ident) children() (children []Node) {
	return []Node{}
}

// Clone the AST.
func (i *Ident) Clone() *Ident {
	return &Ident{
		Parts: i.Parts,
	}
}
