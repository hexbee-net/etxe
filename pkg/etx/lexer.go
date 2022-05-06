package etx

import (
	"regexp"

	"github.com/alecthomas/participle/v2/lexer"
)

const (
	TokenOpBitwiseAnd        = `&`
	TokenOpBitwiseNot        = `~`
	TokenOpBitwiseOr         = `|`
	TokenOpBitwiseShiftLeft  = `<<`
	TokenOpBitwiseShiftRight = `>>`
	TokenOpBitwiseXOr        = `^`
	TokenOpColon             = `:`
	TokenOpCondition         = `?`
	TokenOpDivision          = `/`
	TokenOpEqual             = `==`
	TokenOpLBracket          = `[`
	TokenOpLParen            = `(`
	TokenOpLess              = `<`
	TokenOpLessOrEqual       = `<=`
	TokenOpLogicalAnd        = `&&`
	TokenOpLogicalNot        = `!`
	TokenOpLogicalOr         = `||`
	TokenOpMinus             = `-`
	TokenOpModulo            = `%`
	TokenOpMore              = `>`
	TokenOpMoreOrEqual       = `>=`
	TokenOpMultiplication    = `*`
	TokenOpNotEqual          = `!=`
	TokenOpPlus              = `+`
	TokenOpRBracket          = `]`
	TokenOpRParen            = `)`
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

			{Name: `OpEqual`, Pattern: regexp.QuoteMeta(TokenOpEqual)},
			{Name: `OpNotEqual`, Pattern: regexp.QuoteMeta(TokenOpNotEqual)},
			{Name: `OpLogicalAnd`, Pattern: regexp.QuoteMeta(TokenOpLogicalAnd)},
			{Name: `OpLogicalOr`, Pattern: regexp.QuoteMeta(TokenOpLogicalOr)},
			{Name: `OpBitwiseShiftLeft`, Pattern: regexp.QuoteMeta(TokenOpBitwiseShiftLeft)},
			{Name: `OpBitwiseShiftRight`, Pattern: regexp.QuoteMeta(TokenOpBitwiseShiftRight)},
			{Name: `OpLogicalNot`, Pattern: regexp.QuoteMeta(TokenOpLogicalNot)},
			{Name: `OpBitwiseNot`, Pattern: regexp.QuoteMeta(TokenOpBitwiseNot)},
			{Name: `OpBitwiseAnd`, Pattern: regexp.QuoteMeta(TokenOpBitwiseAnd)},
			{Name: `OpBitwiseOr`, Pattern: regexp.QuoteMeta(TokenOpBitwiseOr)},
			{Name: `OpBitwiseXOr`, Pattern: regexp.QuoteMeta(TokenOpBitwiseXOr)},
			{Name: `OpMultiplication`, Pattern: regexp.QuoteMeta(TokenOpMultiplication)},
			{Name: `OpDivision`, Pattern: regexp.QuoteMeta(TokenOpDivision)},
			{Name: `OpModulo`, Pattern: regexp.QuoteMeta(TokenOpModulo)},
			{Name: `OpPlus`, Pattern: regexp.QuoteMeta(TokenOpPlus)},
			{Name: `OpMinus`, Pattern: regexp.QuoteMeta(TokenOpMinus)},
			{Name: `OpLessOrEqual`, Pattern: regexp.QuoteMeta(TokenOpLessOrEqual)},
			{Name: `OpMoreOrEqual`, Pattern: regexp.QuoteMeta(TokenOpMoreOrEqual)},
			{Name: `OpLess`, Pattern: regexp.QuoteMeta(TokenOpLess)},
			{Name: `OpMore`, Pattern: regexp.QuoteMeta(TokenOpMore)},
			{Name: `OpCondition`, Pattern: regexp.QuoteMeta(TokenOpCondition)},
			{Name: `OpColon`, Pattern: regexp.QuoteMeta(TokenOpColon)},
			{Name: `OpLParen`, Pattern: regexp.QuoteMeta(TokenOpLParen)},
			{Name: `OpRParen`, Pattern: regexp.QuoteMeta(TokenOpRParen)},
			{Name: `OpLBracket`, Pattern: regexp.QuoteMeta(TokenOpLBracket)},
			{Name: `OpRBracket`, Pattern: regexp.QuoteMeta(TokenOpRBracket)},

			lexer.Include("Root"),
		},
	}
}
