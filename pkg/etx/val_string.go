package etx

import (
	"fmt"
	"strings"
)

type String []*Fragment

func (s String) String() string {
	var sb strings.Builder
	for _, f := range s {
		sb.WriteString(f.String())
	}

	return sb.String()
}

func (s String) Clone() String {
	if s == nil {
		return nil
	}

	out := make(String, 0, len(s))
	for _, f := range s {
		out = append(out, f.Clone())
	}

	return out
}

type Fragment struct {
	Escaped   string `parser:"(  @Escaped"                           json:"escaped,omitempty"`
	Unicode   string `parser:" | Unicode@(UnicodeLong|UnicodeShort)" json:"unicode,omitempty"`
	Expr      *Expr  `parser:" | Expr @@ ExprEnd"                    json:"expr,omitempty"`
	Directive *Expr  `parser:" | Directive @@ ExprEnd"               json:"directive,omitempty"`
	Text      string `parser:" | @(Char|Quote|NonExpr))"             json:"text,omitempty"`
}

func (f *Fragment) String() string {
	switch {
	case f.Escaped != "":
		return f.Escaped
	case f.Unicode != "":
		return fmt.Sprintf("\\u%s", f.Unicode)
	case f.Expr != nil:
		return fmt.Sprintf("${%s}", f.Expr)
	case f.Directive != nil:
		return fmt.Sprintf("%%{%s}", f.Directive)
	case f.Text != "":
		return f.Text
	default:
		return ""
	}
}

func (f *Fragment) Clone() *Fragment {
	if f == nil {
		return nil
	}

	out := &Fragment{}

	switch {
	case f.Escaped != "":
		out.Escaped = f.Escaped
	case f.Unicode != "":
		out.Unicode = f.Unicode
	case f.Expr != nil:
		out.Expr = f.Expr.Clone()
	case f.Directive != nil:
		out.Directive = f.Directive.Clone()
	case f.Text != "":
		out.Text = f.Text
	}

	return out
}
