package etx

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/repr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestExpr_Parsing(t *testing.T) {
	type args struct {
		Input string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    *Expr
	}{
		{
			name: "Ternary",
			args: args{
				Input: "1 ? 2 : 3",
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Conditional{
					Condition:    testBuildExprTree[*LogicalOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
					ConditionOp:  "?",
					Left:         testBuildExprTree[*LogicalOr](t, &Value{Number: &Number{big.NewFloat(2)}}),
					ConditionSep: ":",
					Right:        testBuildExprTree[*LogicalOr](t, &Value{Number: &Number{big.NewFloat(3)}}),
				},
			),
		},
		{
			name: "Invalid Ternary",
			args: args{
				Input: "1 ? 2 ; 3",
			},
			wantErr: true,
			want:    nil,
		},

		{
			name: "Logical OR - no spaces",
			args: args{
				Input: `1||2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&LogicalOr{
					Left:  testBuildExprTree[*LogicalAnd](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "||",
					Right: testBuildExprTree[*LogicalOr](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},
		{
			name: "Logical OR - spaces",
			args: args{
				Input: `1 || 2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&LogicalOr{
					Left:  testBuildExprTree[*LogicalAnd](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "||",
					Right: testBuildExprTree[*LogicalOr](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},

		{
			name: "Logical AND - no spaces",
			args: args{
				Input: `1&&2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&LogicalAnd{
					Left:  testBuildExprTree[*BitwiseOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "&&",
					Right: testBuildExprTree[*LogicalAnd](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},
		{
			name: "Logical AND - spaces",
			args: args{
				Input: `1 && 2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&LogicalAnd{
					Left:  testBuildExprTree[*BitwiseOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "&&",
					Right: testBuildExprTree[*LogicalAnd](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},

		{
			name: "Bitwise OR - no spaces",
			args: args{
				Input: `1|2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&BitwiseOr{
					Left:  testBuildExprTree[*BitwiseXor](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "|",
					Right: testBuildExprTree[*BitwiseOr](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},
		{
			name: "Bitwise OR - spaces",
			args: args{
				Input: `1 | 2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&BitwiseOr{
					Left:  testBuildExprTree[*BitwiseXor](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "|",
					Right: testBuildExprTree[*BitwiseOr](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},

		{
			name: "Bitwise XOR - no spaces",
			args: args{
				Input: `1^2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&BitwiseXor{
					Left:  testBuildExprTree[*BitwiseAnd](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "^",
					Right: testBuildExprTree[*BitwiseXor](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},
		{
			name: "Bitwise XOR - spaces",
			args: args{
				Input: `1 ^ 2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&BitwiseXor{
					Left:  testBuildExprTree[*BitwiseAnd](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "^",
					Right: testBuildExprTree[*BitwiseXor](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},

		{
			name: "Bitwise AND - no spaces",
			args: args{
				Input: `1&2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&BitwiseAnd{
					Left:  testBuildExprTree[*Equality](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "&",
					Right: testBuildExprTree[*BitwiseAnd](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},
		{
			name: "Bitwise AND - spaces",
			args: args{
				Input: `1 & 2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&BitwiseAnd{
					Left:  testBuildExprTree[*Equality](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "&",
					Right: testBuildExprTree[*BitwiseAnd](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},

		{
			name: "Equality - Equal - no spaces",
			args: args{
				Input: `1==2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Equality{
					Left:  testBuildExprTree[*Relational](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "==",
					Right: testBuildExprTree[*Equality](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},
		{
			name: "Equality - Equal - spaces",
			args: args{
				Input: `1 == 2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Equality{
					Left:  testBuildExprTree[*Relational](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "==",
					Right: testBuildExprTree[*Equality](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},
		{
			name: "Equality - Not Equal - no spaces",
			args: args{
				Input: `1!=2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Equality{
					Left:  testBuildExprTree[*Relational](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "!=",
					Right: testBuildExprTree[*Equality](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},
		{
			name: "Equality - Not Equal - spaces",
			args: args{
				Input: `1 != 2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Equality{
					Left:  testBuildExprTree[*Relational](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "!=",
					Right: testBuildExprTree[*Equality](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},

		{
			name: "Relational - More - no spaces",
			args: args{
				Input: `1>2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Relational{
					Left:  testBuildExprTree[*Shift](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    ">",
					Right: testBuildExprTree[*Relational](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},
		{
			name: "Relational - More - spaces",
			args: args{
				Input: `1 > 2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Relational{
					Left:  testBuildExprTree[*Shift](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    ">",
					Right: testBuildExprTree[*Relational](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},
		{
			name: "Relational - Less - no spaces",
			args: args{
				Input: `1<2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Relational{
					Left:  testBuildExprTree[*Shift](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "<",
					Right: testBuildExprTree[*Relational](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},
		{
			name: "Relational - Less - spaces",
			args: args{
				Input: `1 < 2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Relational{
					Left:  testBuildExprTree[*Shift](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "<",
					Right: testBuildExprTree[*Relational](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},
		{
			name: "Relational - More Or Equal - no spaces",
			args: args{
				Input: `1>=2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Relational{
					Left:  testBuildExprTree[*Shift](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    ">=",
					Right: testBuildExprTree[*Relational](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},
		{
			name: "Relational - More Or Equal - spaces",
			args: args{
				Input: `1 >= 2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Relational{
					Left:  testBuildExprTree[*Shift](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    ">=",
					Right: testBuildExprTree[*Relational](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},
		{
			name: "Relational - Less Or Equal - no spaces",
			args: args{
				Input: `1<=2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Relational{
					Left:  testBuildExprTree[*Shift](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "<=",
					Right: testBuildExprTree[*Relational](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},
		{
			name: "Relational - Less Or Equal - spaces",
			args: args{
				Input: `1 <= 2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Relational{
					Left:  testBuildExprTree[*Shift](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "<=",
					Right: testBuildExprTree[*Relational](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},

		{
			name: "Shift - Left - no spaces",
			args: args{
				Input: `1<<2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Shift{
					Left:  testBuildExprTree[*Additive](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "<<",
					Right: testBuildExprTree[*Shift](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},
		{
			name: "Shift - Left - spaces",
			args: args{
				Input: `1 << 2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Shift{
					Left:  testBuildExprTree[*Additive](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "<<",
					Right: testBuildExprTree[*Shift](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},
		{
			name: "Shift - Right - no spaces",
			args: args{
				Input: `1>>2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Shift{
					Left:  testBuildExprTree[*Additive](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    ">>",
					Right: testBuildExprTree[*Shift](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},
		{
			name: "Shift - Right - spaces",
			args: args{
				Input: `1 >> 2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Shift{
					Left:  testBuildExprTree[*Additive](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    ">>",
					Right: testBuildExprTree[*Shift](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},

		{
			name: "Additive - Plus - no spaces",
			args: args{
				Input: `1+2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Additive{
					Left:  testBuildExprTree[*Multiplicative](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "+",
					Right: testBuildExprTree[*Additive](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},
		{
			name: "Additive - Plus - spaces",
			args: args{
				Input: `1 + 2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Additive{
					Left:  testBuildExprTree[*Multiplicative](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "+",
					Right: testBuildExprTree[*Additive](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},
		{
			name: "Additive - Minus - no spaces",
			args: args{
				Input: `1-2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Additive{
					Left:  testBuildExprTree[*Multiplicative](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "-",
					Right: testBuildExprTree[*Additive](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},
		{
			name: "Additive - Minus - spaces",
			args: args{
				Input: `1 - 2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Additive{
					Left:  testBuildExprTree[*Multiplicative](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "-",
					Right: testBuildExprTree[*Additive](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},

		{
			name: "Multiplicative - Division - no spaces",
			args: args{
				Input: `1/2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Multiplicative{
					Left:  testBuildExprTree[*Unary](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "/",
					Right: testBuildExprTree[*Multiplicative](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},
		{
			name: "Multiplicative - Division - spaces",
			args: args{
				Input: `1 / 2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Multiplicative{
					Left:  testBuildExprTree[*Unary](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "/",
					Right: testBuildExprTree[*Multiplicative](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},
		{
			name: "Multiplicative - Multiplication - no spaces",
			args: args{
				Input: `1*2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Multiplicative{
					Left:  testBuildExprTree[*Unary](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "*",
					Right: testBuildExprTree[*Multiplicative](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},
		{
			name: "Multiplicative - Multiplication - spaces",
			args: args{
				Input: `1 * 2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Multiplicative{
					Left:  testBuildExprTree[*Unary](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "*",
					Right: testBuildExprTree[*Multiplicative](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},
		{
			name: "Multiplicative - Modulo - no spaces",
			args: args{
				Input: `1%2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Multiplicative{
					Left:  testBuildExprTree[*Unary](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "%",
					Right: testBuildExprTree[*Multiplicative](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},
		{
			name: "Multiplicative - Modulo - spaces",
			args: args{
				Input: `1 % 2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Multiplicative{
					Left:  testBuildExprTree[*Unary](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    "%",
					Right: testBuildExprTree[*Multiplicative](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},

		{
			name: "Unary - Bitwise NOT - no spaces",
			args: args{
				Input: `~1`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Unary{
					Op:    "~",
					Unary: testBuildExprTree[*Unary](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			),
		},
		{
			name: "Unary - Bitwise NOT - spaces",
			args: args{
				Input: `~ 1`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Unary{
					Op:    "~",
					Unary: testBuildExprTree[*Unary](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			),
		},
		{
			name: "Unary - Logical NOT - no spaces",
			args: args{
				Input: `!1`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Unary{
					Op:    "!",
					Unary: testBuildExprTree[*Unary](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			),
		},
		{
			name: "Unary - Logical NOT - spaces",
			args: args{
				Input: `! 1`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Unary{
					Op:    "!",
					Unary: testBuildExprTree[*Unary](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			),
		},
		{
			name: "Unary - Minus - no spaces",
			args: args{
				Input: `-1`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Unary{
					Op:    "-",
					Unary: testBuildExprTree[*Unary](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			),
		},
		{
			name: "Unary - Minus - spaces",
			args: args{
				Input: `- 1`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Unary{
					Op:    "-",
					Unary: testBuildExprTree[*Unary](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			),
		},

		{
			name: "Postfix - no spaces",
			args: args{
				Input: `1[2]`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Postfix{
					Left:  testBuildExprTree[*Primary](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Right: testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},
		{
			name: "Postfix - spaces",
			args: args{
				Input: `1 [ 2 ]`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Postfix{
					Left:  testBuildExprTree[*Primary](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Right: testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			),
		},

		{
			name: "Primary - Value",
			args: args{
				Input: `1`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Primary{
					Value: testBuildExprTree[*Value](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			),
		},
		{
			name: "Primary - Sub - no spaces",
			args: args{
				Input: `(1)`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Primary{
					SubExpression: testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			),
		},
		{
			name: "Primary - Sub - spaces",
			args: args{
				Input: `( 1 )`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Primary{
					SubExpression: testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			),
		},

		{
			name: "Add 3",
			args: args{
				Input: `1 + 2 + 3`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&Additive{
					Left: testBuildExprTree[*Multiplicative](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:   "+",
					Right: &Additive{
						Left:  testBuildExprTree[*Multiplicative](t, &Value{Number: &Number{big.NewFloat(2)}}),
						Op:    "+",
						Right: testBuildExprTree[*Additive](t, &Value{Number: &Number{big.NewFloat(3)}}),
					},
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type Exp struct {
				Expr *Expr `parser:"@@*"`
			}

			l := lexer.MustStateful(lexRules(), lexer.InitialState("Expr"))
			parser := participle.MustBuild(&Exp{}, participle.Lexer(l))

			res := &Exp{}
			err := parser.ParseString("", tt.args.Input, res)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			repr.Println(res)

			assert.Equal(t, tt.want, res.Expr)
		})
	}
}
