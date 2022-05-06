package etx

import (
	"reflect"
	"testing"
)

func testValPtr[T any](t *testing.T, v T) *T {
	t.Helper()
	return &v
}

func testBuildExprTree[E any](t *testing.T, value interface{}) E {
	t.Helper()

	var build func(interface{}, reflect.Type) interface{}
	build = func(value interface{}, stop reflect.Type) interface{} {
		if reflect.TypeOf(value) == stop {
		}
		if stop != nil && reflect.TypeOf(value) == stop {
			return value
		}

		switch v := value.(type) {
		case *Expr:
			return v
		case *Conditional:
			return &Expr{Left: v}
		case *LogicalOr:
			return build(&Conditional{Condition: v}, stop)
		case *LogicalAnd:
			return build(&LogicalOr{Left: v}, stop)
		case *BitwiseOr:
			return build(&LogicalAnd{Left: v}, stop)
		case *BitwiseXor:
			return build(&BitwiseOr{Left: v}, stop)
		case *BitwiseAnd:
			return build(&BitwiseXor{Left: v}, stop)
		case *Equality:
			return build(&BitwiseAnd{Left: v}, stop)
		case *Relational:
			return build(&Equality{Left: v}, stop)
		case *Shift:
			return build(&Relational{Left: v}, stop)
		case *Additive:
			return build(&Shift{Left: v}, stop)
		case *Multiplicative:
			return build(&Additive{Left: v}, stop)
		case *Unary:
			return build(&Multiplicative{Left: v}, stop)
		case *Postfix:
			return build(&Unary{Postfix: v}, stop)
		case *Primary:
			return build(&Postfix{Left: v}, stop)
		case *Value:
			return build(&Primary{Value: v}, stop)
		default:
			panic("invalid type for expression tree")
		}
	}

	var stopVal E
	return build(value, reflect.TypeOf(stopVal)).(E)
}
