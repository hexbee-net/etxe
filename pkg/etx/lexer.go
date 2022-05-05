package etx

import (
	"github.com/alecthomas/participle/v2/lexer"
)

const (
	TokenOpBitwiseAnd        = `&`
	TokenOpBitwiseNot        = `~`
	TokenOpBitwiseOr         = `\|`
	TokenOpBitwiseShiftLeft  = `<<`
	TokenOpBitwiseShiftRight = `>>`
	TokenOpBitwiseXOr        = `\^`
	TokenOpColon             = `:`
	TokenOpCondition         = `\?`
	TokenOpDivision          = `\/`
	TokenOpEqual             = `==`
	TokenOpLBracket          = `\[`
	TokenOpLParen            = `\(`
	TokenOpLess              = `<`
	TokenOpLessOrEqual       = `<=`
	TokenOpLogicalAnd        = `&&`
	TokenOpLogicalNot        = `!`
	TokenOpLogicalOr         = `\|\|`
	TokenOpMinus             = `-`
	TokenOpModulo            = `%`
	TokenOpMore              = `>`
	TokenOpMoreOrEqual       = `>=`
	TokenOpMultiplication    = `\*`
	TokenOpNotEqual          = `!=`
	TokenOpPlus              = `\+`
	TokenOpRBracket          = `\]`
	TokenOpRParen            = `\)`
)

func lexRules() lexer.Rules {
	return lexer.Rules{
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

			{Name: `OpEqual`, Pattern: TokenOpEqual},
			{Name: `OpNotEqual`, Pattern: TokenOpNotEqual},
			{Name: `OpLogicalAnd`, Pattern: TokenOpLogicalAnd},
			{Name: `OpLogicalOr`, Pattern: TokenOpLogicalOr},
			{Name: `OpBitwiseShiftLeft`, Pattern: TokenOpBitwiseShiftLeft},
			{Name: `OpBitwiseShiftRight`, Pattern: TokenOpBitwiseShiftRight},
			{Name: `OpLogicalNot`, Pattern: TokenOpLogicalNot},
			{Name: `OpBitwiseNot`, Pattern: TokenOpBitwiseNot},
			{Name: `OpBitwiseAnd`, Pattern: TokenOpBitwiseAnd},
			{Name: `OpBitwiseOr`, Pattern: TokenOpBitwiseOr},
			{Name: `OpBitwiseXOr`, Pattern: TokenOpBitwiseXOr},
			{Name: `OpMultiplication`, Pattern: TokenOpMultiplication},
			{Name: `OpDivision`, Pattern: TokenOpDivision},
			{Name: `OpModulo`, Pattern: TokenOpModulo},
			{Name: `OpPlus`, Pattern: TokenOpPlus},
			{Name: `OpMinus`, Pattern: TokenOpMinus},
			{Name: `OpLessOrEqual`, Pattern: TokenOpLessOrEqual},
			{Name: `OpMoreOrEqual`, Pattern: TokenOpMoreOrEqual},
			{Name: `OpLess`, Pattern: TokenOpLess},
			{Name: `OpMore`, Pattern: TokenOpMore},
			{Name: `OpCondition`, Pattern: TokenOpCondition},
			{Name: `OpColon`, Pattern: TokenOpColon},
			{Name: `OpLParen`, Pattern: TokenOpLParen},
			{Name: `OpRParen`, Pattern: TokenOpRParen},
			{Name: `OpLBracket`, Pattern: TokenOpLBracket},
			{Name: `OpRBracket`, Pattern: TokenOpRBracket},

			lexer.Include("Root"),
		},
	}
}
