package etx

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/repr"
)

// Value is a scalar, list or map.
type Value struct {
	Pos    lexer.Position `parser:"" json:"-"`
	Parent Node           `parser:"" json:"-"`

	Null             bool        `parser:"(  @'null'" json:"null,omitempty"`
	Bool             *Bool       `parser:" | @('true' | 'false')" json:"bool,omitempty"`
	Number           *Number     `parser:" | @Number" json:"number,omitempty"`
	Str              String      `parser:" | String @@* StringEnd"`
	Ident            *string     `parser:" | @Ident" json:"ident,omitempty"`
	HeredocDelimiter string      `parser:" | (@Heredoc" json:"heredoc_delimiter,omitempty"`
	Heredoc          *string     `parser:"     @(Body | EOL)* End)" json:"heredoc,omitempty"`
	HaveList         bool        `parser:" | ( @'['" json:"have_list,omitempty"` // Need this to detect empty lists.
	List             []*Value    `parser:"     ( Whitespace? @@ ( Whitespace?',' Whitespace? @@ )* )? Whitespace? ','? Whitespace? ']' )" json:"list,omitempty"`
	HaveMap          bool        `parser:" | ( @'{'" json:"have_map,omitempty"` // Need this to detect empty maps.
	Map              []*MapEntry `parser:"     ( Whitespace? @@ ( Whitespace? ',' Whitespace? @@ )* Whitespace? ','? )? Whitespace? '}' ) )" json:"map,omitempty"`
}

// Clone the AST.
func (v *Value) Clone() *Value {
	if v == nil {
		return nil
	}

	out := &Value{}
	*out = *v

	switch {
	case out.Number != nil:
		out.Number = &Number{}
		out.Number.Float.Copy(v.Number.Float)

	case out.Str != nil:
		out.Str = v.Str.Clone()

	case v.HaveList:
		out.List = make([]*Value, len(v.List))
		for i, value := range v.List {
			out.List[i] = value.Clone()
		}

	case v.HaveMap:
		out.Map = make([]*MapEntry, len(v.Map))
		for i, entry := range out.Map {
			out.Map[i] = entry.Clone()
		}
	}

	return out
}

func (v *Value) children() (children []Node) {
	// TODO: Add string fragments.

	for _, el := range v.List {
		children = append(children, el)
	}
	for _, el := range v.Map {
		children = append(children, el)
	}

	return
}

func (v *Value) String() string {
	switch {
	case v.Null == true:
		return "null"

	case v.Bool != nil:
		return fmt.Sprintf("%v", *v.Bool)

	case v.Number != nil:
		return v.Number.String()

	case v.Str != nil:
		return v.Str.String()

	case v.Ident != nil:
		return fmt.Sprintf("%v", *v.Ident)

	case v.HeredocDelimiter != "":
		heredoc := ""
		if v.Heredoc != nil {
			heredoc = *v.Heredoc
		}

		return fmt.Sprintf("<<%s%s\n%s", v.HeredocDelimiter, heredoc, strings.TrimPrefix(v.HeredocDelimiter, "-"))

	case v.HaveList:
		entries := make([]string, 0, len(v.List))
		for _, e := range v.List {
			entries = append(entries, e.String())
		}

		return fmt.Sprintf("[%s]", strings.Join(entries, ", "))

	case v.HaveMap:
		entries := make([]string, 0, len(v.Map))
		for _, e := range v.Map {
			entries = append(entries, fmt.Sprintf("%s: %s", e.Key, e.Value))
		}

		return fmt.Sprintf("{%s}", strings.Join(entries, ", "))

	default:
		panic(repr.String(v, repr.Hide(lexer.Position{})))
	}
}
