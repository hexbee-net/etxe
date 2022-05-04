package etx

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

const (
	parserLookahead = 50
)

var (
	stripCommentRe = regexp.MustCompile(`^//\s*|^/\*|\*/$`)
)

func parser() *participle.Parser {
	return participle.MustBuild(&AST{},
		participle.Lexer(lex()),
		participle.Map(unquoteString, "String"),
		participle.Map(cleanHeredocStart, "Heredoc"),
		participle.Map(stripComment, "Comment"),
		// We need lookahead to ensure prefixed comments are associated with the right nodes.
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

func unquoteString(token lexer.Token) (lexer.Token, error) {
	if token.Value[0] == '\'' {
		token.Value = "\"" + strings.ReplaceAll(token.Value[1:len(token.Value)-1], "\"", "\\\"") + "\""
	}

	var err error
	token.Value, err = strconv.Unquote(token.Value)
	if err != nil {
		return token, fmt.Errorf("%s: %w", token.Pos, err)
	}

	return token, nil
}

func cleanHeredocStart(token lexer.Token) (lexer.Token, error) {
	token.Value = token.Value[2:]

	return token, nil
}

func stripComment(token lexer.Token) (lexer.Token, error) {
	token.Value = stripCommentRe.ReplaceAllString(token.Value, "")

	return token, nil
}

func cloneStrings(strings []string) []string {
	if strings == nil {
		return nil
	}
	out := make([]string, len(strings))
	copy(out, strings)

	return out
}
