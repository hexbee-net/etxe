package etx

import (
	"github.com/alecthomas/participle/v2/lexer"
)

func lex() *lexer.StatefulDefinition {
	return lexer.MustStateful(lexer.Rules{
		"Root": {
			// {"Ident", `\b[[:alpha:]]\w*(-\w+)*\b`, nil},
			// {"Number", `^[-+]?\d*\.?\d+([eE][-+]?\d+)?`, nil},
			// {"Heredoc", `<<[-]?(\w+\b)`, stateful.Push("Heredoc")},
			{Name: "String", Pattern: `(["'])`, Action: lexer.Push("String")},
			// {"Punct", `[][{}=:,]`, nil},
			// {"Comment", `(?:(?://|#)[^\n]*)|/\*.*?\*/`, nil},
			// {"whitespace", `\s+`, nil},
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
		// "Heredoc": {
		// 	{"End", `\n\s*\b\1\b`, stateful.Pop()},
		// 	{"EOL", `\n`, nil},
		// 	{"Body", `[^\n]+`, nil},
		// },
		"Expr": {
			lexer.Include("Root"),
			{Name: `Whitespace`, Pattern: `\s+`},
			{Name: `Oper`, Pattern: `[-+/*%]`},
			{Name: "Ident", Pattern: `\w+`},
			{Name: "ExprEnd", Pattern: `}`, Action: lexer.Pop()}},
	})
}
