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
	ast := &AST{}
	if err := parser().Parse("", r, ast); err != nil {
		return nil, fmt.Errorf("parsing failed: %w", err)
	}

	ast.UpdateParentRefs()

	return ast, nil
}

// ParseString parses ETX from a string.
func ParseString(str string) (*AST, error) {
	ast := &AST{}
	if err := parser().ParseString("", str, ast); err != nil {
		return nil, fmt.Errorf("parsing failed: %w", err)
	}

	ast.UpdateParentRefs()

	return ast, nil
}

// ParseBytes parses ETX from bytes.
func ParseBytes(data []byte) (*AST, error) {
	ast := &AST{}
	if err := parser().ParseBytes("", data, ast); err != nil {
		return nil, fmt.Errorf("parsing failed: %w", err)
	}

	ast.UpdateParentRefs()

	return ast, nil
}

// updateParentRefs recursively updates an AST parent references.
func updateParentRefs(parent, node Node) {
	node.Node().Parent = parent

	for _, c := range node.Children() {
		updateParentRefs(node, c)
	}
}
