package etx

import "github.com/alecthomas/participle/lexer"

// MapEntry represents a key+value in a map.
type MapEntry struct {
	Pos    lexer.Position `parser:"" json:"-"`
	Parent Node           `parser:"" json:"-"`

	Comments []string `parser:"@Comment*" json:"comments,omitempty"`
	Key      *Value   `parser:"@@ ':'"    json:"key"`
	Value    *Value   `parser:"@@"        json:"value"`
}

func (e *MapEntry) children() (children []Node) {
	return []Node{e.Key, e.Value}
}

// Clone the AST.
func (e *MapEntry) Clone() *MapEntry {
	if e == nil {
		return nil
	}

	return &MapEntry{
		Pos:      e.Pos,
		Key:      e.Key.Clone(),
		Value:    e.Value.Clone(),
		Comments: cloneStrings(e.Comments),
	}
}
