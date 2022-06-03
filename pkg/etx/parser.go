package etx

import (
	"fmt"
	"io"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

const (
	parserLookahead = 50
)

func parser() *participle.Parser {
	return participle.MustBuild(&AST{},
		participle.Lexer(lexer.MustStateful(lexRules())),
		participle.UseLookahead(parserLookahead))
}

// Parse ETX from an io.Reader.
func Parse(r io.Reader) (*AST, error) {
	etx := &AST{}
	if err := parser().Parse("", r, etx); err != nil {
		return nil, fmt.Errorf("parsing failed: %w", err)
	}

	return etx, AddParentRefs(etx)
}

// ParseString parses ETX from a string.
func ParseString(str string) (*AST, error) {
	etx := &AST{}
	if err := parser().ParseString("", str, etx); err != nil {
		return nil, fmt.Errorf("parsing failed: %w", err)
	}

	return etx, AddParentRefs(etx)
}

// ParseBytes parses ETX from bytes.
func ParseBytes(data []byte) (*AST, error) {
	etx := &AST{}
	if err := parser().ParseBytes("", data, etx); err != nil {
		return nil, fmt.Errorf("parsing failed: %w", err)
	}

	return etx, AddParentRefs(etx)
}
