package etx

import (
	"fmt"
	"io"
	"reflect"
	"testing"
)

const (
	indentationChar = "\t"
)

type Cloner[C any] interface {
	Clone() C
}

type FormattedStringer interface {
	FormattedString() string
}

func cloneCollection[T Cloner[T]](src []T) []T {
	if src == nil {
		return nil
	}

	out := make([]T, 0, len(src))
	for _, item := range src {
		out = append(out, item.Clone())
	}

	return out
}

func cloneStrings(strings []string) []string {
	if strings == nil {
		return nil
	}
	out := make([]string, len(strings))
	copy(out, strings)

	return out
}

// indent inserts prefix at the beginning of each non-empty line of s.
func indent(s, prefix string) string {
	return string(indentBytes([]byte(s), []byte(prefix)))
}

// indentBytes inserts prefix at the beginning of each non-empty line of b.
func indentBytes(b, prefix []byte) []byte {
	var res []byte
	bol := true
	for _, c := range b {
		if bol && c != '\n' {
			res = append(res, prefix...)
		}
		res = append(res, c)
		bol = c == '\n'
	}

	return res
}

func mustFprintf(w io.Writer, format string, a ...any) {
	if _, err := fmt.Fprintf(w, format, a...); err != nil {
		panic(err)
	}
}

func BuildTestExprTree[E any](t *testing.T, value interface{}) E {
	t.Helper()

	var build func(interface{}, reflect.Type) interface{}
	build = func(value interface{}, stop reflect.Type) interface{} {
		if (stop != nil && reflect.TypeOf(value) == stop) || (reflect.TypeOf(value) == reflect.TypeOf(&Expr{})) {
			return value
		}

		switch v := value.(type) {
		case *ExprConditional:
			return &Expr{ASTNode: v.ASTNode, Left: v}
		case *ExprLogicalOr:
			return build(&ExprConditional{ASTNode: v.ASTNode, Condition: *v}, stop)
		case *ExprLogicalAnd:
			return build(&ExprLogicalOr{ASTNode: v.ASTNode, Left: *v}, stop)
		case *ExprBitwiseOr:
			return build(&ExprLogicalAnd{ASTNode: v.ASTNode, Left: *v}, stop)
		case *ExprBitwiseXor:
			return build(&ExprBitwiseOr{ASTNode: v.ASTNode, Left: *v}, stop)
		case *ExprBitwiseAnd:
			return build(&ExprBitwiseXor{ASTNode: v.ASTNode, Left: *v}, stop)
		case *ExprEquality:
			return build(&ExprBitwiseAnd{ASTNode: v.ASTNode, Left: *v}, stop)
		case *ExprRelational:
			return build(&ExprEquality{ASTNode: v.ASTNode, Left: *v}, stop)
		case *ExprShift:
			return build(&ExprRelational{ASTNode: v.ASTNode, Left: *v}, stop)
		case *ExprAdditive:
			return build(&ExprShift{ASTNode: v.ASTNode, Left: *v}, stop)
		case *ExprMultiplicative:
			return build(&ExprAdditive{ASTNode: v.ASTNode, Left: *v}, stop)
		case *ExprUnary:
			return build(&ExprMultiplicative{ASTNode: v.ASTNode, Left: *v}, stop)
		case *ExprPostfix:
			return build(&ExprUnary{ASTNode: v.ASTNode, Right: *v}, stop)
		case *ExprPrimary:
			return build(&ExprPostfix{ASTNode: v.ASTNode, Value: *v}, stop)
		case *Ident:
			return build(&ExprPrimary{ASTNode: v.ASTNode, Ident: v}, stop)
		case *Value:
			return build(&ExprPrimary{ASTNode: v.ASTNode, Value: v}, stop)
		default:
			panic("invalid type for expression tree")
		}
	}

	var stopVal E

	return build(value, reflect.TypeOf(stopVal)).(E) //nolint:forcetypeassert // only used in tests
}
