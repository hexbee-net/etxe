package etx

import (
	"regexp"

	"github.com/alecthomas/participle/v2/lexer"
)

const (
	OpAssign            = `=`
	OpBitwiseAnd        = `&`
	OpBitwiseNot        = `~`
	OpBitwiseOr         = `|`
	OpBitwiseShiftLeft  = `<<`
	OpBitwiseShiftRight = `>>`
	OpBitwiseXOr        = `^`
	OpColon             = `:`
	OpCondition         = `?`
	OpDivision          = `/`
	OpEqual             = `==`
	OpLBracket          = `[`
	OpLParen            = `(`
	OpLess              = `<`
	OpLessOrEqual       = `<=`
	OpLogicalAnd        = `&&`
	OpLogicalNot        = `!`
	OpLogicalOr         = `||`
	OpMinus             = `-`
	OpModulo            = `%`
	OpMore              = `>`
	OpMoreOrEqual       = `>=`
	OpMultiplication    = `*`
	OpNotEqual          = `!=`
	OpPlus              = `+`
	OpRBracket          = `]`
	OpRParen            = `)`
	OpLambda            = `=>`
	OpLambdaDef         = `->`
	OpLBrace            = `{`
	OpRBrace            = `}`
	OpComma             = `,`
	OpDot               = `.`

	TokenWhitespace = "whitespace"
)

const (
	lexerCore       = "Core"
	lexerRoot       = "Root"
	lexerString     = "String"
	lexerStringExpr = "StringExpr"
	lexerUnicode    = "Unicode"
	lexerHeredoc    = "Heredoc"
)

func lexRules() lexer.Rules {
	return lexer.Rules{
		lexerRoot: {
			{Name: `Input`, Pattern: `\b(input)\b`},
			{Name: `Output`, Pattern: `\b(output)\b`},
			{Name: `Const`, Pattern: `\b(const)\b`},
			{Name: `Val`, Pattern: `\b(val)\b`},
			{Name: `Type`, Pattern: `\b(type)\b`},
			{Name: `Enum`, Pattern: `\b(enum)\b`},
			{Name: `Object`, Pattern: `\b(object)\b`},
			{Name: `Func`, Pattern: `\b(def)\b`},
			{Name: `Return`, Pattern: `\b(return)\b`},

			lexer.Include(lexerCore),
		},
		lexerCore: {
			{Name: "SingleLineComment", Pattern: `(?:\/\/|#)[^\n]*`},
			{Name: "MultilineComment", Pattern: `\/\*(.|\n)*?\*\/`},

			{Name: "OpLambda", Pattern: regexp.QuoteMeta(OpLambda)},
			{Name: "OpLambdaDef", Pattern: regexp.QuoteMeta(OpLambdaDef)},

			{Name: `OpLParen`, Pattern: regexp.QuoteMeta(OpLParen)},
			{Name: `OpRParen`, Pattern: regexp.QuoteMeta(OpRParen)},

			{Name: `If`, Pattern: `\b(if)\b`},
			{Name: `Else`, Pattern: `\b(else)\b`},
			{Name: `Switch`, Pattern: `\b(switch)\b`},
			{Name: `Case`, Pattern: `\b(case)\b`},

			{Name: `BlockStart`, Pattern: regexp.QuoteMeta(OpLBrace)},
			{Name: `BlockEnd`, Pattern: regexp.QuoteMeta(OpRBrace)},

			{Name: "Ident", Pattern: `\b[[:alpha:]]\w*(-\w+)*\b`},
			{Name: "Number", Pattern: `(0[xX][0-9a-fA-F_]+|0[bB][01_]*|0[oO][0-7_]*|[0-9_]*\.?[0-9_]+([eE][-+]?[0-9_]+)?)`},

			{Name: "Heredoc", Pattern: `<<[-]?(\w+)\n`, Action: lexer.Push(lexerHeredoc)},

			{Name: `OpComma`, Pattern: regexp.QuoteMeta(OpComma)},
			{Name: `OpEqual`, Pattern: regexp.QuoteMeta(OpEqual)},
			{Name: `OpNotEqual`, Pattern: regexp.QuoteMeta(OpNotEqual)},
			{Name: `OpLogicalAnd`, Pattern: regexp.QuoteMeta(OpLogicalAnd)},
			{Name: `OpLogicalOr`, Pattern: regexp.QuoteMeta(OpLogicalOr)},
			{Name: `OpBitwiseShiftLeft`, Pattern: regexp.QuoteMeta(OpBitwiseShiftLeft)},
			{Name: `OpBitwiseShiftRight`, Pattern: regexp.QuoteMeta(OpBitwiseShiftRight)},
			{Name: `OpLogicalNot`, Pattern: regexp.QuoteMeta(OpLogicalNot)},
			{Name: `OpBitwiseNot`, Pattern: regexp.QuoteMeta(OpBitwiseNot)},
			{Name: `OpBitwiseAnd`, Pattern: regexp.QuoteMeta(OpBitwiseAnd)},
			{Name: `OpBitwiseOr`, Pattern: regexp.QuoteMeta(OpBitwiseOr)},
			{Name: `OpBitwiseXOr`, Pattern: regexp.QuoteMeta(OpBitwiseXOr)},
			{Name: `OpMultiplication`, Pattern: regexp.QuoteMeta(OpMultiplication)},
			{Name: `OpDivision`, Pattern: regexp.QuoteMeta(OpDivision)},
			{Name: `OpModulo`, Pattern: regexp.QuoteMeta(OpModulo)},
			{Name: `OpPlus`, Pattern: regexp.QuoteMeta(OpPlus)},
			{Name: `OpMinus`, Pattern: regexp.QuoteMeta(OpMinus)},
			{Name: `OpLessOrEqual`, Pattern: regexp.QuoteMeta(OpLessOrEqual)},
			{Name: `OpMoreOrEqual`, Pattern: regexp.QuoteMeta(OpMoreOrEqual)},
			{Name: `OpLess`, Pattern: regexp.QuoteMeta(OpLess)},
			{Name: `OpMore`, Pattern: regexp.QuoteMeta(OpMore)},
			{Name: `OpAssign`, Pattern: regexp.QuoteMeta(OpAssign)},
			{Name: `OpCondition`, Pattern: regexp.QuoteMeta(OpCondition)},
			{Name: `OpColon`, Pattern: regexp.QuoteMeta(OpColon)},
			{Name: `OpLParen`, Pattern: regexp.QuoteMeta(OpLParen)},
			{Name: `OpRParen`, Pattern: regexp.QuoteMeta(OpRParen)},
			{Name: `OpLBracket`, Pattern: regexp.QuoteMeta(OpLBracket)},
			{Name: `OpRBracket`, Pattern: regexp.QuoteMeta(OpRBracket)},

			{Name: "String", Pattern: `(["'])`, Action: lexer.Push(lexerString)},
			{Name: "Dot", Pattern: regexp.QuoteMeta(OpDot)},
			{Name: "LF", Pattern: `[\r\n]`},
			{Name: TokenWhitespace, Pattern: `[\t ]+`},
		},

		lexerString: {
			{Name: "Unicode", Pattern: `\\u`, Action: lexer.Push(lexerUnicode)},
			{Name: "Escaped", Pattern: `\\.`},
			{Name: "StringEnd", Pattern: `\1`, Action: lexer.Pop()},
			{Name: "Quote", Pattern: `["']`},
			{Name: "NonExpr", Pattern: `(\$\${|%%{)`},
			{Name: "Expr", Pattern: `\${`, Action: lexer.Push(lexerStringExpr)},
			{Name: "Directive", Pattern: `%{`, Action: lexer.Push(lexerStringExpr)},
			{Name: "Char", Pattern: `[^$%"'\\]+`},
		},
		lexerUnicode: {
			{Name: "UnicodeLong", Pattern: `[0-9a-fA-F]{8}`, Action: lexer.Pop()},
			{Name: "UnicodeShort", Pattern: `[0-9a-fA-F]{4}`, Action: lexer.Pop()},
		},
		lexerHeredoc: {
			{Name: "HeredocEnd", Pattern: `^\1`, Action: lexer.Pop()},
			{Name: "EOL", Pattern: `\n`},
			{Name: "NonExpr", Pattern: `(\$\${|%%{)`},
			{Name: "Expr", Pattern: `\${`, Action: lexer.Push(lexerStringExpr)},
			{Name: "Directive", Pattern: `%{`, Action: lexer.Push(lexerStringExpr)},
			{Name: "Body", Pattern: `[^\n$%]+`},
		},
		lexerStringExpr: {
			{Name: "ExprEnd", Pattern: `}`, Action: lexer.Pop()},

			lexer.Include(lexerCore),
		},
	}
}
