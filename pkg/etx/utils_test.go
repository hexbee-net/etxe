package etx

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
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
		case *ExprConditional:
			return &Expr{Left: v}
		case *ExprLogicalOr:
			return build(&ExprConditional{Condition: v}, stop)
		case *ExprLogicalAnd:
			return build(&ExprLogicalOr{Left: v}, stop)
		case *ExprBitwiseOr:
			return build(&ExprLogicalAnd{Left: v}, stop)
		case *ExprBitwiseXor:
			return build(&ExprBitwiseOr{Left: v}, stop)
		case *ExprBitwiseAnd:
			return build(&ExprBitwiseXor{Left: v}, stop)
		case *ExprEquality:
			return build(&ExprBitwiseAnd{Left: v}, stop)
		case *ExprRelational:
			return build(&ExprEquality{Left: v}, stop)
		case *ExprShift:
			return build(&ExprRelational{Left: v}, stop)
		case *ExprAdditive:
			return build(&ExprShift{Left: v}, stop)
		case *ExprMultiplicative:
			return build(&ExprAdditive{Left: v}, stop)
		case *ExprUnary:
			return build(&ExprMultiplicative{Left: v}, stop)
		case *ExprPostfix:
			return build(&ExprUnary{Postfix: v}, stop)
		case *ExprPrimary:
			return build(&ExprPostfix{Left: v}, stop)
		case *Value:
			return build(&ExprPrimary{Value: v}, stop)
		default:
			panic("invalid type for expression tree")
		}
	}

	var stopVal E
	return build(value, reflect.TypeOf(stopVal)).(E)
}

func TestIndent(t *testing.T) {
	type args struct {
		s      string
		prefix string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "One char indent",
			args: args{
				s: `
Not tricks, Michael, illusions.
I care deeply for nature.
Well, what do you expect, mother? Get me a vodka rocks.
And a piece of toast.
Marry me.`[1:],
				prefix: "X",
			},
			want: `
XNot tricks, Michael, illusions.
XI care deeply for nature.
XWell, what do you expect, mother? Get me a vodka rocks.
XAnd a piece of toast.
XMarry me.`[1:],
		},
		{
			name: "Multiple chars indent",
			args: args{
				s: `
Not tricks, Michael, illusions.
I care deeply for nature.
Well, what do you expect, mother? Get me a vodka rocks.
And a piece of toast.
Marry me.`[1:],
				prefix: "XXX",
			},
			want: `
XXXNot tricks, Michael, illusions.
XXXI care deeply for nature.
XXXWell, what do you expect, mother? Get me a vodka rocks.
XXXAnd a piece of toast.
XXXMarry me.`[1:],
		},
		{
			name: "Empty lines",
			args: args{
				s: `
Not tricks, Michael, illusions.
I care deeply for nature.

Well, what do you expect, mother? Get me a vodka rocks.

And a piece of toast.
Marry me.`[1:],
				prefix: "XXX",
			},
			want: `
XXXNot tricks, Michael, illusions.
XXXI care deeply for nature.

XXXWell, what do you expect, mother? Get me a vodka rocks.

XXXAnd a piece of toast.
XXXMarry me.`[1:],
		},
		{
			name: "Trailing empty lines",
			args: args{
				s: `
Not tricks, Michael, illusions.
I care deeply for nature.
Well, what do you expect, mother? Get me a vodka rocks.
And a piece of toast.
Marry me.


`[1:],
				prefix: "XXX",
			},
			want: `
XXXNot tricks, Michael, illusions.
XXXI care deeply for nature.
XXXWell, what do you expect, mother? Get me a vodka rocks.
XXXAnd a piece of toast.
XXXMarry me.


`[1:],
		},
		{
			name: "Preceding empty lines",
			args: args{
				s: `


Not tricks, Michael, illusions.
I care deeply for nature.
Well, what do you expect, mother? Get me a vodka rocks.
And a piece of toast.
Marry me.`[1:],
				prefix: "XXX",
			},
			want: `


XXXNot tricks, Michael, illusions.
XXXI care deeply for nature.
XXXWell, what do you expect, mother? Get me a vodka rocks.
XXXAnd a piece of toast.
XXXMarry me.`[1:],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, indent(tt.args.s, tt.args.prefix), "indent(%v, %v)", tt.args.s, tt.args.prefix)
		})
	}
}
