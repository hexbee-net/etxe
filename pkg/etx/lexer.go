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
	TokenOpLambda            = `=>`
	TokenOpLambdaDef         = `->`
	TokenOpLBrace            = `{`
	TokenOpRBrace            = `}`
	TokenOpComma             = `,`

	TokenKeywordInput   = `input`
	TokenKeywordOutput  = `output`
	TokenKeywordConst   = `const`
	TokenKeywordVal     = `val`
	TokenKeywordType    = `type`
	TokenKeywordEnum    = `enum`
	TokenKeywordIf      = `if`
	TokenKeywordElse    = `else`
	TokenKeywordSwitch  = `switch`
	TokenKeywordCase    = `case`
	TokenKeywordDefault = `default`
	TokenKeywordReturn  = `return`
)

const (
	lexerCore       = "Core"
	lexerRoot       = "Root"
	lexerString     = "String"
	lexerStringExpr = "StringExpr"
	lexerUnicode    = "Unicode"
	lexerHeredoc    = "Heredoc"
	lexerExpr       = "Expr"
	lexerFunc       = "Func"
)

func lexRules() lexer.Rules {
	return lexer.Rules{
		lexerRoot: {
			{Name: "Input", Pattern: `\b(` + TokenKeywordInput + `)\b`},
			{Name: "Output", Pattern: `\b(` + TokenKeywordOutput + `)\b`},
			{Name: "Const", Pattern: `\b(` + TokenKeywordConst + `)\b`},
			{Name: "Val", Pattern: `\b(` + TokenKeywordVal + `)\b`},
			{Name: "Type", Pattern: `\b(` + TokenKeywordType + `)\b`},
			{Name: "Func", Pattern: `\b(def)\b`, Action: lexer.Push(lexerFunc)},
			{Name: "Punctuation", Pattern: `[][{}=:,]`},

			lexer.Include(lexerCore),
		},
		lexerCore: {
			{Name: "Ident", Pattern: `\b[[:alpha:]]\w*(-\w+)*\b`},
			{Name: "Number", Pattern: `[-+]?(0[xX][0-9a-fA-F_]+|0[bB][01_]*|0[oO][0-7_]*|[0-9_]*\.?[0-9_]+([eE][-+]?[0-9_]+)?)`},
			{Name: "Heredoc", Pattern: `<<[-]?(\w+\b)`, Action: lexer.Push(lexerHeredoc)},
			{Name: "String", Pattern: `(["'])`, Action: lexer.Push(lexerString)},
			{Name: "Comment", Pattern: `(?:(?:\/\/|#).*?$)|\/\*.*?\*\/`},
			{Name: `Whitespace`, Pattern: `\s+`},
			{Name: `NewLine`, Pattern: `(\n|\n\r)+`},

			{Name: "Lambda", Pattern: regexp.QuoteMeta(TokenOpLambda), Action: lexer.Push(lexerStringExpr)},
			{Name: "LambdaDef", Pattern: regexp.QuoteMeta(TokenOpLambdaDef)},
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
			{Name: "End", Pattern: `\n\s*\b\1\b`, Action: lexer.Pop()},
			{Name: "EOL", Pattern: `\n`},
			{Name: "Body", Pattern: `[^\n]+`},
		},
		lexerStringExpr: {
			{Name: "ExprEnd", Pattern: `}`, Action: lexer.Pop()},

			lexer.Include(lexerExpr),
		},

		lexerExpr: {
			{Name: "If", Pattern: `\b(` + TokenKeywordIf + `)\b`},
			{Name: "Else", Pattern: `\b(` + TokenKeywordElse + `)\b`},
			{Name: "Switch", Pattern: `\b(` + TokenKeywordSwitch + `)\b`},
			{Name: "Case", Pattern: `\b(` + TokenKeywordCase + `)\b`},
			{Name: "Default", Pattern: `\b(` + TokenKeywordDefault + `)\b`},

			{Name: `BlockStart`, Pattern: regexp.QuoteMeta(TokenOpLBrace)},
			{Name: `BlockEnd`, Pattern: regexp.QuoteMeta(TokenOpRBrace)},

			{Name: `OpComma`, Pattern: regexp.QuoteMeta(TokenOpComma)},

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

			lexer.Include(lexerCore),

			lexer.Return(),
		},

		lexerFunc: {
			{Name: "BodyStart", Pattern: `{`},
			{Name: "BodyEnd", Pattern: `}`, Action: lexer.Pop()},

			{Name: "FuncPunctuation", Pattern: `[(),:]`},

			{Name: "Const", Pattern: `\b(` + TokenKeywordConst + `)\b`},
			{Name: "Val", Pattern: `\b(` + TokenKeywordVal + `)\b`},
			{Name: "Return", Pattern: `\b(` + TokenKeywordReturn + `)\b`},

			lexer.Include(lexerCore),
		},
	}
}
