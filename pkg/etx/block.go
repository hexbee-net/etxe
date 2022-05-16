package etx

import "github.com/alecthomas/participle/v2/lexer"

// Block represents an optionally labeled block.
type Block struct {
	Pos    lexer.Position `parser:"" json:"-"`
	Parent Node           `parser:"" json:"-"`

	Comments         []string     `parser:"@Comment*"            json:"comments,omitempty"`
	Name             string       `parser:"@Ident"               json:"name"`
	Labels           []string     `parser:"@( Ident | String )*" json:"labels,omitempty"`
	Body             []*BlockItem `parser:"'{' @@*       "       json:"body"`
	TrailingComments []string     `parser:"@Comment* '}' "       json:"trailing_comments,omitempty"`
}

// BlockItem in a block.
type BlockItem struct {
	Pos    lexer.Position `parser:"" json:"-"`
	Parent Node           `parser:"" json:"-"`

	Block     *Block     `parser:"(   @@  " json:"block,omitempty"`
	Attribute *Attribute `parser:"  | @@ )" json:"attribute,omitempty"`
}
