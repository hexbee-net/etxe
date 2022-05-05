package etx

import (
	"github.com/alecthomas/participle/v2/lexer"
)

func lex() *lexer.StatefulDefinition {
	return lexer.MustStateful(lexer.Rules{
		"Root": {
			{Name: "Ident", Pattern: `\b[[:alpha:]]\w*(-\w+)*\b`},
			{Name: "Number", Pattern: `[-+]?(0[xX][0-9a-fA-F_]+|0[bB][01_]*|0[oO][0-7_]*|[0-9_]*\.?[0-9_]+([eE][-+]?[0-9_]+)?)`},
			{Name: "Heredoc", Pattern: `<<[-]?(\w+\b)`, Action: lexer.Push("Heredoc")},
			{Name: "String", Pattern: `(["'])`, Action: lexer.Push("String")},
			{Name: "Punctuation", Pattern: `[][{}=:,]`},
			{Name: "Comment", Pattern: `(?:(?://|#)[^\n]*)|/\*.*?\*/`},
			{Name: `Whitespace`, Pattern: `\s+`},
		},
		"String": {
			{Name: "Unicode", Pattern: `\\u`, Action: lexer.Push("Unicode")},
			{Name: "Escaped", Pattern: `\\.`},
			{Name: "StringEnd", Pattern: `\1`, Action: lexer.Pop()},
			{Name: "Quote", Pattern: `["']`},
			{Name: "NonExpr", Pattern: `(\$\${|%%{)`},
			{Name: "Expr", Pattern: `\${`, Action: lexer.Push("Expr")},
			{Name: "Directive", Pattern: `%{`, Action: lexer.Push("Expr")},
			{Name: "Char", Pattern: `[^$%"'\\]+`},
		},
		"Unicode": {
			{Name: "UnicodeLong", Pattern: `[0-9a-fA-F]{8}`, Action: lexer.Pop()},
			{Name: "UnicodeShort", Pattern: `[0-9a-fA-F]{4}`, Action: lexer.Pop()},
		},
		"Heredoc": {
			{Name: "End", Pattern: `\n\s*\b\1\b`, Action: lexer.Pop()},
			{Name: "EOL", Pattern: `\n`},
			{Name: "Body", Pattern: `[^\n]+`},
		},
		"Expr": {
			{Name: "ExprEnd", Pattern: `}`, Action: lexer.Pop()},
			lexer.Include("Root"),
			{Name: `Operator`, Pattern: `[-+/*%]`},
		},
	})
}
