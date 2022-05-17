package etx

import "github.com/alecthomas/participle/v2/lexer"

// Attribute is a key+value attribute.
type Attribute struct {
	Pos    lexer.Position `parser:"" json:"-"`
	Parent Node           `parser:"" json:"-"`

	Comments []string `parser:"@Comment*" json:"comments,omitempty"`
	Key      string   `parser:"@Ident"    json:"key"`
	Value    *Expr    `parser:"['=' @@ ]" json:"value,omitempty"`
}
