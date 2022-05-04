package etx

import (
	"fmt"
	"strings"
)

type String struct {
	Fragments []*Fragment `parser:"String @@* StringEnd"`
}

func (s *String) String() string {
	var sb strings.Builder
	for _, f := range s.Fragments {
		sb.WriteString(f.String())
	}
	return sb.String()
}

type Fragment struct {
	Escaped   string `parser:"(  @Escaped"`
	Unicode   string `parser:" | Unicode@(UnicodeLong|UnicodeShort)"`
	Expr      *Expr  `parser:" | Expr @@ ExprEnd"`
	Directive *Expr  `parser:" | Directive @@ ExprEnd"`
	Text      string `parser:" | @(Char|Quote|NonExpr))"`
}

func (f *Fragment) String() string {
	if f.Escaped != "" {
		return f.Escaped
	}
	if f.Unicode != "" {
		return fmt.Sprintf("\\u%s", f.Unicode)
	}
	if f.Expr != nil {
		return fmt.Sprintf("${%s}", f.Expr)
	}
	if f.Directive != nil {
		return fmt.Sprintf("%%{%s}", f.Directive)
	}
	if f.Text != "" {
		return f.Text
	}

	return ""
}
