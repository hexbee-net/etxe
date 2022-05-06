package etx

import (
	"math/big"
	"testing"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
					ConditionOp:  TokenOpCondition,
					True:         testBuildExprTree[*LogicalOr](t, &Value{Number: &Number{big.NewFloat(2)}}),
					ConditionSep: TokenOpColon,
					False:        testBuildExprTree[*LogicalOr](t, &Value{Number: &Number{big.NewFloat(3)}}),
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
					Op:    TokenOpLogicalOr,
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
					Op:    TokenOpLogicalOr,
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
					Op:    TokenOpLogicalAnd,
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
					Op:    TokenOpLogicalAnd,
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
					Op:    TokenOpBitwiseOr,
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
					Op:    TokenOpBitwiseOr,
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
					Op:    TokenOpBitwiseXOr,
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
					Op:    TokenOpBitwiseXOr,
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
					Op:    TokenOpBitwiseAnd,
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
					Op:    TokenOpBitwiseAnd,
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
					Op:    TokenOpEqual,
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
					Op:    TokenOpEqual,
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
					Op:    TokenOpNotEqual,
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
					Op:    TokenOpNotEqual,
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
					Op:    TokenOpMore,
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
					Op:    TokenOpMore,
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
					Op:    TokenOpLess,
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
					Op:    TokenOpLess,
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
					Op:    TokenOpMoreOrEqual,
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
					Op:    TokenOpMoreOrEqual,
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
					Op:    TokenOpLessOrEqual,
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
					Op:    TokenOpLessOrEqual,
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
					Op:    TokenOpBitwiseShiftLeft,
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
					Op:    TokenOpBitwiseShiftLeft,
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
					Op:    TokenOpBitwiseShiftRight,
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
					Op:    TokenOpBitwiseShiftRight,
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
					Op:    TokenOpPlus,
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
					Op:    TokenOpPlus,
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
					Op:    TokenOpMinus,
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
					Op:    TokenOpMinus,
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
					Op:    TokenOpDivision,
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
					Op:    TokenOpDivision,
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
					Op:    TokenOpMultiplication,
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
					Op:    TokenOpMultiplication,
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
					Op:    TokenOpModulo,
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
					Op:    TokenOpModulo,
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
					Op:    TokenOpBitwiseNot,
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
					Op:    TokenOpBitwiseNot,
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
					Op:    TokenOpLogicalNot,
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
					Op:    TokenOpLogicalNot,
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
					Op:    TokenOpMinus,
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
					Op:    TokenOpMinus,
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

			assert.Equal(t, tt.want, res.Expr)
		})
	}
}

// /////////////////////////////////////

func TestExpr_String(t *testing.T) {}

func TestConditional_String(t *testing.T) {
	type args struct {
		Input *Conditional
	}
	tests := []struct {
		name        string
		description string
		args        args
		wantPanic   bool
		want        string
	}{
		{
			name: "Nil",
			args: args{
				Input: nil,
			},
			wantPanic: true,
		},
		{
			name:        "Empty",
			description: "Condition can never be nil",
			args: args{
				Input: &Conditional{},
			},
			wantPanic: true,
			want:      "condition cannot be <nil>",
		},
		{
			name: "No Condition operator",
			args: args{
				Input: &Conditional{
					Condition:    testBuildExprTree[*LogicalOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
					ConditionSep: TokenOpColon,
					True:         testBuildExprTree[*LogicalOr](t, &Value{Number: &Number{big.NewFloat(2)}}),
					False:        testBuildExprTree[*LogicalOr](t, &Value{Number: &Number{big.NewFloat(3)}}),
				},
			},
			wantPanic: true,
			want:      "both operators need to be set",
		},
		{
			name: "No Sep operator",
			args: args{
				Input: &Conditional{
					Condition:   testBuildExprTree[*LogicalOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
					ConditionOp: TokenOpCondition,
					True:        testBuildExprTree[*LogicalOr](t, &Value{Number: &Number{big.NewFloat(2)}}),
					False:       testBuildExprTree[*LogicalOr](t, &Value{Number: &Number{big.NewFloat(3)}}),
				},
			},
			wantPanic: true,
			want:      "both operators need to be set",
		},
		{
			name: "Only Condition",
			args: args{
				Input: &Conditional{
					Condition: testBuildExprTree[*LogicalOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			},
			want: "1",
		},
		{
			name: "All expressions but not operators",
			args: args{
				Input: &Conditional{
					Condition: testBuildExprTree[*LogicalOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
					True:      testBuildExprTree[*LogicalOr](t, &Value{Number: &Number{big.NewFloat(2)}}),
					False:     testBuildExprTree[*LogicalOr](t, &Value{Number: &Number{big.NewFloat(3)}}),
				},
			},
			want: "1",
		},
		{
			name:        "Neither True nor False",
			description: "Both sides of the condition must be present.",
			args: args{
				Input: &Conditional{
					Condition:    testBuildExprTree[*LogicalOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
					ConditionOp:  TokenOpCondition,
					ConditionSep: TokenOpColon,
				},
			},
			wantPanic: true,
			want:      "true and false expressions must be set when operators are set",
		},
		{
			name:        "Only True expression",
			description: "Both sides of the condition must be present.",
			args: args{
				Input: &Conditional{
					Condition:    testBuildExprTree[*LogicalOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
					ConditionOp:  TokenOpCondition,
					ConditionSep: TokenOpColon,
					True:         testBuildExprTree[*LogicalOr](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			},
			wantPanic: true,
			want:      "false expression must be set when operators are set",
		},
		{
			name:        "Only False expression",
			description: "Both sides of the condition must be present.",
			args: args{
				Input: &Conditional{
					Condition:    testBuildExprTree[*LogicalOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
					ConditionOp:  TokenOpCondition,
					ConditionSep: TokenOpColon,
					False:        testBuildExprTree[*LogicalOr](t, &Value{Number: &Number{big.NewFloat(3)}}),
				},
			},
			wantPanic: true,
			want:      "true expression must be set when operators are set",
		},
		{
			name:        "All parts",
			description: "Both sides of the condition must be present.",
			args: args{
				Input: &Conditional{
					Condition:    testBuildExprTree[*LogicalOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
					ConditionOp:  TokenOpCondition,
					ConditionSep: TokenOpColon,
					True:         testBuildExprTree[*LogicalOr](t, &Value{Number: &Number{big.NewFloat(2)}}),
					False:        testBuildExprTree[*LogicalOr](t, &Value{Number: &Number{big.NewFloat(3)}}),
				},
			},
			want: "1 ? 2 : 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				if tt.want != "" {
					assert.PanicsWithValuef(t, tt.want, func() {
						_ = tt.args.Input.String()
					}, tt.description)
				} else {
					assert.Panicsf(t, func() {
						_ = tt.args.Input.String()
					}, tt.description)
				}
			} else {
				assert.Equal(t, tt.want, tt.args.Input.String())
			}
		})
	}
}

func TestLogicalOr_String(t *testing.T) {
	type args struct {
		Input *LogicalOr
	}
	tests := []struct {
		name        string
		description string
		args        args
		wantPanic   bool
		want        string
	}{
		{
			name: "Nil",
			args: args{
				Input: nil,
			},
			wantPanic: true,
		},
		{
			name:        "Empty",
			description: "Left side can never be nil",
			args: args{
				Input: &LogicalOr{},
			},
			wantPanic: true,
			want:      "left side cannot be <nil>",
		},
		{
			name: "Left",
			args: args{
				Input: &LogicalOr{
					Left: testBuildExprTree[*LogicalAnd](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			args: args{
				Input: &LogicalOr{
					Left: &LogicalAnd{},
					Op:   TokenOpLogicalOr,
				},
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			args: args{
				Input: &LogicalOr{
					Left:  testBuildExprTree[*LogicalAnd](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    TokenOpLogicalOr,
					Right: testBuildExprTree[*LogicalOr](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			},
			want: "1 || 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				if tt.want != "" {
					assert.PanicsWithValuef(t, tt.want, func() {
						_ = tt.args.Input.String()
					}, tt.description)
				} else {
					assert.Panicsf(t, func() {
						_ = tt.args.Input.String()
					}, tt.description)
				}
			} else {
				assert.Equal(t, tt.want, tt.args.Input.String())
			}
		})
	}
}

func TestLogicalAnd_String(t *testing.T) {
	type args struct {
		Input *LogicalAnd
	}
	tests := []struct {
		name        string
		description string
		args        args
		wantPanic   bool
		want        string
	}{
		{
			name: "Nil",
			args: args{
				Input: nil,
			},
			wantPanic: true,
		},
		{
			name:        "Empty",
			description: "Left side can never be nil",
			args: args{
				Input: &LogicalAnd{},
			},
			wantPanic: true,
			want:      "left side cannot be <nil>",
		},
		{
			name: "Left",
			args: args{
				Input: &LogicalAnd{
					Left: testBuildExprTree[*BitwiseOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			args: args{
				Input: &LogicalAnd{
					Left: &BitwiseOr{},
					Op:   TokenOpLogicalAnd,
				},
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			args: args{
				Input: &LogicalAnd{
					Left:  testBuildExprTree[*BitwiseOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    TokenOpLogicalAnd,
					Right: testBuildExprTree[*LogicalAnd](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			},
			want: "1 && 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				if tt.want != "" {
					assert.PanicsWithValuef(t, tt.want, func() {
						_ = tt.args.Input.String()
					}, tt.description)
				} else {
					assert.Panicsf(t, func() {
						_ = tt.args.Input.String()
					}, tt.description)
				}
			} else {
				assert.Equal(t, tt.want, tt.args.Input.String())
			}
		})
	}
}

func TestBitwiseOr_String(t *testing.T) {
	type args struct {
		Input *BitwiseOr
	}
	tests := []struct {
		name        string
		description string
		args        args
		wantPanic   bool
		want        string
	}{
		{
			name: "Nil",
			args: args{
				Input: nil,
			},
			wantPanic: true,
		},
		{
			name:        "Empty",
			description: "Left side can never be nil",
			args: args{
				Input: &BitwiseOr{},
			},
			wantPanic: true,
			want:      "left side cannot be <nil>",
		},
		{
			name: "Left",
			args: args{
				Input: &BitwiseOr{
					Left: testBuildExprTree[*BitwiseXor](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			args: args{
				Input: &BitwiseOr{
					Left: &BitwiseXor{},
					Op:   TokenOpBitwiseOr,
				},
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			args: args{
				Input: &BitwiseOr{
					Left:  testBuildExprTree[*BitwiseXor](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    TokenOpBitwiseOr,
					Right: testBuildExprTree[*BitwiseOr](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			},
			want: "1 | 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				if tt.want != "" {
					assert.PanicsWithValuef(t, tt.want, func() {
						_ = tt.args.Input.String()
					}, tt.description)
				} else {
					assert.Panicsf(t, func() {
						_ = tt.args.Input.String()
					}, tt.description)
				}
			} else {
				assert.Equal(t, tt.want, tt.args.Input.String())
			}
		})
	}
}

func TestBitwiseXor_String(t *testing.T) {
	type args struct {
		Input *BitwiseXor
	}
	tests := []struct {
		name        string
		description string
		args        args
		wantPanic   bool
		want        string
	}{
		{
			name: "Nil",
			args: args{
				Input: nil,
			},
			wantPanic: true,
		},
		{
			name:        "Empty",
			description: "Left side can never be nil",
			args: args{
				Input: &BitwiseXor{},
			},
			wantPanic: true,
			want:      "left side cannot be <nil>",
		},
		{
			name: "Left",
			args: args{
				Input: &BitwiseXor{
					Left: testBuildExprTree[*BitwiseAnd](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			args: args{
				Input: &BitwiseXor{
					Left: &BitwiseAnd{},
					Op:   TokenOpBitwiseXOr,
				},
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			args: args{
				Input: &BitwiseXor{
					Left:  testBuildExprTree[*BitwiseAnd](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    TokenOpBitwiseXOr,
					Right: testBuildExprTree[*BitwiseXor](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			},
			want: "1 ^ 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				if tt.want != "" {
					assert.PanicsWithValuef(t, tt.want, func() {
						_ = tt.args.Input.String()
					}, tt.description)
				} else {
					assert.Panicsf(t, func() {
						_ = tt.args.Input.String()
					}, tt.description)
				}
			} else {
				assert.Equal(t, tt.want, tt.args.Input.String())
			}
		})
	}
}

func TestBitwiseAnd_String(t *testing.T) {
	type args struct {
		Input *BitwiseAnd
	}
	tests := []struct {
		name        string
		description string
		args        args
		wantPanic   bool
		want        string
	}{
		{
			name: "Nil",
			args: args{
				Input: nil,
			},
			wantPanic: true,
		},
		{
			name:        "Empty",
			description: "Left side can never be nil",
			args: args{
				Input: &BitwiseAnd{},
			},
			wantPanic: true,
			want:      "left side cannot be <nil>",
		},
		{
			name: "Left",
			args: args{
				Input: &BitwiseAnd{
					Left: testBuildExprTree[*Equality](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			args: args{
				Input: &BitwiseAnd{
					Left: &Equality{},
					Op:   TokenOpBitwiseAnd,
				},
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			args: args{
				Input: &BitwiseAnd{
					Left:  testBuildExprTree[*Equality](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    TokenOpBitwiseAnd,
					Right: testBuildExprTree[*BitwiseAnd](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			},
			want: "1 & 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				if tt.want != "" {
					assert.PanicsWithValuef(t, tt.want, func() {
						_ = tt.args.Input.String()
					}, tt.description)
				} else {
					assert.Panicsf(t, func() {
						_ = tt.args.Input.String()
					}, tt.description)
				}
			} else {
				assert.Equal(t, tt.want, tt.args.Input.String())
			}
		})
	}
}

func TestEquality_String(t *testing.T) {
	type args struct {
		Input *Equality
	}
	tests := []struct {
		name        string
		description string
		args        args
		wantPanic   bool
		want        string
	}{
		{
			name: "Nil",
			args: args{
				Input: nil,
			},
			wantPanic: true,
		},
		{
			name:        "Empty",
			description: "Left side can never be nil",
			args: args{
				Input: &Equality{},
			},
			wantPanic: true,
			want:      "left side cannot be <nil>",
		},
		{
			name: "Left",
			args: args{
				Input: &Equality{
					Left: testBuildExprTree[*Relational](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			args: args{
				Input: &Equality{
					Left: &Relational{},
					Op:   TokenOpEqual,
				},
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			args: args{
				Input: &Equality{
					Left:  testBuildExprTree[*Relational](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    TokenOpEqual,
					Right: testBuildExprTree[*Equality](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			},
			want: "1 == 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				if tt.want != "" {
					assert.PanicsWithValuef(t, tt.want, func() {
						_ = tt.args.Input.String()
					}, tt.description)
				} else {
					assert.Panicsf(t, func() {
						_ = tt.args.Input.String()
					}, tt.description)
				}
			} else {
				assert.Equal(t, tt.want, tt.args.Input.String())
			}
		})
	}
}

func TestRelational_String(t *testing.T) {
	type args struct {
		Input *Relational
	}
	tests := []struct {
		name        string
		description string
		args        args
		wantPanic   bool
		want        string
	}{
		{
			name: "Nil",
			args: args{
				Input: nil,
			},
			wantPanic: true,
		},
		{
			name:        "Empty",
			description: "Left side can never be nil",
			args: args{
				Input: &Relational{},
			},
			wantPanic: true,
			want:      "left side cannot be <nil>",
		},
		{
			name: "Left",
			args: args{
				Input: &Relational{
					Left: testBuildExprTree[*Shift](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			args: args{
				Input: &Relational{
					Left: &Shift{},
					Op:   TokenOpMoreOrEqual,
				},
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			args: args{
				Input: &Relational{
					Left:  testBuildExprTree[*Shift](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    TokenOpMoreOrEqual,
					Right: testBuildExprTree[*Relational](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			},
			want: "1 >= 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				if tt.want != "" {
					assert.PanicsWithValuef(t, tt.want, func() {
						_ = tt.args.Input.String()
					}, tt.description)
				} else {
					assert.Panicsf(t, func() {
						_ = tt.args.Input.String()
					}, tt.description)
				}
			} else {
				assert.Equal(t, tt.want, tt.args.Input.String())
			}
		})
	}
}

func TestShift_String(t *testing.T) {
	type args struct {
		Input *Shift
	}
	tests := []struct {
		name        string
		description string
		args        args
		wantPanic   bool
		want        string
	}{
		{
			name: "Nil",
			args: args{
				Input: nil,
			},
			wantPanic: true,
		},
		{
			name:        "Empty",
			description: "Left side can never be nil",
			args: args{
				Input: &Shift{},
			},
			wantPanic: true,
			want:      "left side cannot be <nil>",
		},
		{
			name: "Left",
			args: args{
				Input: &Shift{
					Left: testBuildExprTree[*Additive](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			args: args{
				Input: &Shift{
					Left: &Additive{},
					Op:   TokenOpBitwiseShiftRight,
				},
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			args: args{
				Input: &Shift{
					Left:  testBuildExprTree[*Additive](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    TokenOpBitwiseShiftRight,
					Right: testBuildExprTree[*Shift](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			},
			want: "1 >> 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				if tt.want != "" {
					assert.PanicsWithValuef(t, tt.want, func() {
						_ = tt.args.Input.String()
					}, tt.description)
				} else {
					assert.Panicsf(t, func() {
						_ = tt.args.Input.String()
					}, tt.description)
				}
			} else {
				assert.Equal(t, tt.want, tt.args.Input.String())
			}
		})
	}
}

func TestAdditive_String(t *testing.T) {
	type args struct {
		Input *Additive
	}
	tests := []struct {
		name        string
		description string
		args        args
		wantPanic   bool
		want        string
	}{
		{
			name: "Nil",
			args: args{
				Input: nil,
			},
			wantPanic: true,
		},
		{
			name:        "Empty",
			description: "Left side can never be nil",
			args: args{
				Input: &Additive{},
			},
			wantPanic: true,
			want:      "left side cannot be <nil>",
		},
		{
			name: "Left",
			args: args{
				Input: &Additive{
					Left: testBuildExprTree[*Multiplicative](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			args: args{
				Input: &Additive{
					Left: &Multiplicative{},
					Op:   TokenOpPlus,
				},
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			args: args{
				Input: &Additive{
					Left:  testBuildExprTree[*Multiplicative](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    TokenOpPlus,
					Right: testBuildExprTree[*Additive](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			},
			want: "1 + 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				if tt.want != "" {
					assert.PanicsWithValuef(t, tt.want, func() {
						_ = tt.args.Input.String()
					}, tt.description)
				} else {
					assert.Panicsf(t, func() {
						_ = tt.args.Input.String()
					}, tt.description)
				}
			} else {
				assert.Equal(t, tt.want, tt.args.Input.String())
			}
		})
	}
}

func TestMultiplicative_String(t *testing.T) {
	type args struct {
		Input *Multiplicative
	}
	tests := []struct {
		name        string
		description string
		args        args
		wantPanic   bool
		want        string
	}{
		{
			name: "Nil",
			args: args{
				Input: nil,
			},
			wantPanic: true,
		},
		{
			name:        "Empty",
			description: "Left side can never be nil",
			args: args{
				Input: &Multiplicative{},
			},
			wantPanic: true,
			want:      "left side cannot be <nil>",
		},
		{
			name: "Left",
			args: args{
				Input: &Multiplicative{
					Left: testBuildExprTree[*Unary](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			args: args{
				Input: &Multiplicative{
					Left: testBuildExprTree[*Unary](t, &Value{Ident: testValPtr(t, "foo")}),
					Op:   TokenOpMultiplication,
				},
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			args: args{
				Input: &Multiplicative{
					Left:  testBuildExprTree[*Unary](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    TokenOpMultiplication,
					Right: testBuildExprTree[*Multiplicative](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			},
			want: "1 * 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				if tt.want != "" {
					assert.PanicsWithValuef(t, tt.want, func() {
						_ = tt.args.Input.String()
					}, tt.description)
				} else {
					assert.Panicsf(t, func() {
						_ = tt.args.Input.String()
					}, tt.description)
				}
			} else {
				assert.Equal(t, tt.want, tt.args.Input.String())
			}
		})
	}
}

func TestUnary_String(t *testing.T) {
	type args struct {
		Input *Unary
	}
	tests := []struct {
		name        string
		description string
		args        args
		want        string
		wantPanic   bool
	}{
		{
			name: "Nil",
			args: args{
				Input: nil,
			},
			wantPanic: true,
		},
		{
			name: "Empty",
			args: args{
				Input: &Unary{},
			},
			wantPanic: true,
		},
		{
			name: "Postfix",
			args: args{
				Input: &Unary{
					Postfix: testBuildExprTree[*Postfix](t, &Value{Ident: testValPtr(t, "foo")}),
				},
			},
			want: "foo",
		},
		{
			name: "Operator",
			args: args{
				Input: &Unary{
					Op:    TokenOpBitwiseNot,
					Unary: testBuildExprTree[*Unary](t, &Value{Ident: testValPtr(t, "foo")}),
				},
			},
			want: "~foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				assert.Panicsf(t, func() {
					_ = tt.args.Input.String()
				}, tt.description)
			} else {
				assert.Equal(t, tt.want, tt.args.Input.String())
			}
		})
	}
}

func TestPostfix_String(t *testing.T) {
	type args struct {
		Input *Postfix
	}
	tests := []struct {
		name        string
		description string
		args        args
		wantPanic   bool
		want        string
	}{
		{
			name: "Nil",
			args: args{
				Input: nil,
			},
			wantPanic: true,
		},
		{
			name:        "Empty",
			description: "It should not be possible to have the left field not initialized",
			args: args{
				Input: &Postfix{},
			},
			wantPanic: true,
		},
		{
			name: "Left",
			args: args{
				Input: &Postfix{
					Left: testBuildExprTree[*Primary](t, &Value{Ident: testValPtr(t, "foo")}),
				},
			},
			want: "foo",
		},
		{
			name: "Both",
			args: args{
				Input: &Postfix{
					Left:  testBuildExprTree[*Primary](t, &Value{Ident: testValPtr(t, "foo")}),
					Right: testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			},
			want: "foo[1]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				assert.Panicsf(t, func() {
					_ = tt.args.Input.String()
				}, tt.description)
			} else {
				assert.Equal(t, tt.want, tt.args.Input.String())
			}
		})
	}
}

func TestPrimary_String(t *testing.T) {
	type args struct {
		Input *Primary
	}
	tests := []struct {
		name        string
		description string
		args        args
		want        string
		wantPanic   bool
	}{
		{
			name: "Nil",
			args: args{
				Input: nil,
			},
			wantPanic: true,
		},
		{
			name: "Empty",
			args: args{
				Input: &Primary{},
			},
			want: "",
		},
		{
			name: "Sub Expression",
			args: args{
				Input: &Primary{
					SubExpression: testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			},
			want: "1",
		},
		{
			name: "Value",
			args: args{
				Input: &Primary{
					Value: &Value{Number: &Number{big.NewFloat(1)}},
				},
			},
			want: "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				assert.Panicsf(t, func() {
					_ = tt.args.Input.String()
				}, tt.description)
			} else {
				assert.Equal(t, tt.want, tt.args.Input.String())
			}
		})
	}
}

// /////////////////////////////////////

func TestExpr_Clone(t *testing.T) {
	type args struct {
		Input *Expr
	}
	tests := []struct {
		name string
		args args
		want *Expr
	}{
		{
			name: "Nil",
			args: args{
				Input: nil,
			},
			want: nil,
		},
		{
			name: "Empty",
			args: args{
				Input: &Expr{},
			},
			want: &Expr{},
		},

		{
			name: "Left",
			args: args{
				Input: &Expr{
					Left: &Conditional{},
				},
			},
			want: &Expr{
				Left: &Conditional{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.Input.Clone())
		})
	}
}

func TestConditional_Clone(t *testing.T) {
	type args struct {
		Input *Conditional
	}
	tests := []struct {
		name string
		args args
		want *Conditional
	}{
		{
			name: "Nil",
			args: args{
				Input: nil,
			},
			want: nil,
		},
		{
			name: "Empty",
			args: args{
				Input: &Conditional{},
			},
			want: &Conditional{},
		},

		{
			name: "With Data",
			args: args{
				Input: &Conditional{
					Condition:    &LogicalOr{},
					ConditionOp:  TokenOpCondition,
					True:         &LogicalOr{},
					ConditionSep: TokenOpColon,
					False:        &LogicalOr{},
				},
			},
			want: &Conditional{
				Condition:    &LogicalOr{},
				ConditionOp:  TokenOpCondition,
				True:         &LogicalOr{},
				ConditionSep: TokenOpColon,
				False:        &LogicalOr{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.Input.Clone())
		})
	}
}

func TestLogicalAnd_Clone(t *testing.T) {
	type args struct {
		Input *LogicalAnd
	}
	tests := []struct {
		name string
		args args
		want *LogicalAnd
	}{
		{
			name: "Nil",
			args: args{
				Input: nil,
			},
			want: nil,
		},
		{
			name: "Empty",
			args: args{
				Input: &LogicalAnd{},
			},
			want: &LogicalAnd{},
		},

		{
			name: "Left",
			args: args{
				Input: &LogicalAnd{
					Left: &BitwiseOr{},
				},
			},
			want: &LogicalAnd{
				Left: &BitwiseOr{},
			},
		},
		{
			name: "Both",
			args: args{
				Input: &LogicalAnd{
					Left:  &BitwiseOr{},
					Op:    TokenOpLogicalAnd,
					Right: &LogicalAnd{},
				},
			},
			want: &LogicalAnd{
				Left:  &BitwiseOr{},
				Op:    TokenOpLogicalAnd,
				Right: &LogicalAnd{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.Input.Clone())
		})
	}
}

func TestBitwiseAnd_Clone(t *testing.T) {
	type args struct {
		Input *BitwiseAnd
	}
	tests := []struct {
		name string
		args args
		want *BitwiseAnd
	}{
		{
			name: "Nil",
			args: args{
				Input: nil,
			},
			want: nil,
		},
		{
			name: "Empty",
			args: args{
				Input: &BitwiseAnd{},
			},
			want: &BitwiseAnd{},
		},

		{
			name: "Left",
			args: args{
				Input: &BitwiseAnd{
					Left: &Equality{},
				},
			},
			want: &BitwiseAnd{
				Left: &Equality{},
			},
		},
		{
			name: "Both",
			args: args{
				Input: &BitwiseAnd{
					Left:  &Equality{},
					Op:    TokenOpBitwiseAnd,
					Right: &BitwiseAnd{},
				},
			},
			want: &BitwiseAnd{
				Left:  &Equality{},
				Op:    TokenOpBitwiseAnd,
				Right: &BitwiseAnd{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.Input.Clone())
		})
	}
}

func TestEquality_Clone(t *testing.T) {
	type args struct {
		Input *Equality
	}
	tests := []struct {
		name string
		args args
		want *Equality
	}{
		{
			name: "Nil",
			args: args{
				Input: nil,
			},
			want: nil,
		},
		{
			name: "Empty",
			args: args{
				Input: &Equality{},
			},
			want: &Equality{},
		},

		{
			name: "Left",
			args: args{
				Input: &Equality{
					Left: &Relational{},
				},
			},
			want: &Equality{
				Left: &Relational{},
			},
		},
		{
			name: "Both",
			args: args{
				Input: &Equality{
					Left:  &Relational{},
					Op:    TokenOpNotEqual,
					Right: &Equality{},
				},
			},
			want: &Equality{
				Left:  &Relational{},
				Op:    TokenOpNotEqual,
				Right: &Equality{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.Input.Clone())
		})
	}
}

func TestRelational_Clone(t *testing.T) {
	type args struct {
		Input *Relational
	}
	tests := []struct {
		name string
		args args
		want *Relational
	}{
		{
			name: "Nil",
			args: args{
				Input: nil,
			},
			want: nil,
		},
		{
			name: "Empty",
			args: args{
				Input: &Relational{},
			},
			want: &Relational{},
		},

		{
			name: "Left",
			args: args{
				Input: &Relational{
					Left: &Shift{},
				},
			},
			want: &Relational{
				Left: &Shift{},
			},
		},
		{
			name: "Both",
			args: args{
				Input: &Relational{
					Left:  &Shift{},
					Op:    TokenOpLessOrEqual,
					Right: &Relational{},
				},
			},
			want: &Relational{
				Left:  &Shift{},
				Op:    TokenOpLessOrEqual,
				Right: &Relational{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.Input.Clone())
		})
	}
}

func TestShift_Clone(t *testing.T) {
	type args struct {
		Input *Shift
	}
	tests := []struct {
		name string
		args args
		want *Shift
	}{
		{
			name: "Nil",
			args: args{
				Input: nil,
			},
			want: nil,
		},
		{
			name: "Empty",
			args: args{
				Input: &Shift{},
			},
			want: &Shift{},
		},

		{
			name: "Left",
			args: args{
				Input: &Shift{
					Left: &Additive{},
				},
			},
			want: &Shift{
				Left: &Additive{},
			},
		},
		{
			name: "Both",
			args: args{
				Input: &Shift{
					Left:  &Additive{},
					Op:    TokenOpBitwiseShiftRight,
					Right: &Shift{},
				},
			},
			want: &Shift{
				Left:  &Additive{},
				Op:    TokenOpBitwiseShiftRight,
				Right: &Shift{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.Input.Clone())
		})
	}
}

func TestAdditive_Clone(t *testing.T) {
	type args struct {
		Input *Additive
	}
	tests := []struct {
		name string
		args args
		want *Additive
	}{
		{
			name: "Nil",
			args: args{
				Input: nil,
			},
			want: nil,
		},
		{
			name: "Empty",
			args: args{
				Input: &Additive{},
			},
			want: &Additive{},
		},

		{
			name: "Left",
			args: args{
				Input: &Additive{
					Left: &Multiplicative{},
				},
			},
			want: &Additive{
				Left: &Multiplicative{},
			},
		},
		{
			name: "Both",
			args: args{
				Input: &Additive{
					Left:  &Multiplicative{},
					Op:    TokenOpMinus,
					Right: &Additive{},
				},
			},
			want: &Additive{
				Left:  &Multiplicative{},
				Op:    TokenOpMinus,
				Right: &Additive{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.Input.Clone())
		})
	}
}

func TestMultiplicative_Clone(t *testing.T) {
	type args struct {
		Input *Multiplicative
	}
	tests := []struct {
		name string
		args args
		want *Multiplicative
	}{
		{
			name: "Nil",
			args: args{
				Input: nil,
			},
			want: nil,
		},
		{
			name: "Empty",
			args: args{
				Input: &Multiplicative{},
			},
			want: &Multiplicative{},
		},

		{
			name: "Left",
			args: args{
				Input: &Multiplicative{
					Left: &Unary{},
				},
			},
			want: &Multiplicative{
				Left: &Unary{},
			},
		},
		{
			name: "Both",
			args: args{
				Input: &Multiplicative{
					Left:  &Unary{},
					Op:    TokenOpModulo,
					Right: &Multiplicative{},
				},
			},
			want: &Multiplicative{
				Left:  &Unary{},
				Op:    TokenOpModulo,
				Right: &Multiplicative{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.Input.Clone())
		})
	}
}

func TestUnary_Clone(t *testing.T) {
	type args struct {
		Input *Unary
	}
	tests := []struct {
		name string
		args args
		want *Unary
	}{
		{
			name: "Nil",
			args: args{
				Input: nil,
			},
			want: nil,
		},
		{
			name: "Empty",
			args: args{
				Input: &Unary{},
			},
			want: &Unary{},
		},

		{
			name: "Unary",
			args: args{
				Input: &Unary{
					Op:    TokenOpMinus,
					Unary: &Unary{},
				},
			},
			want: &Unary{
				Op:    TokenOpMinus,
				Unary: &Unary{},
			},
		},
		{
			name: "Postfix",
			args: args{
				Input: &Unary{
					Postfix: &Postfix{},
				},
			},
			want: &Unary{
				Postfix: &Postfix{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.Input.Clone())
		})
	}
}

func TestPostfix_Clone(t *testing.T) {
	type args struct {
		Input *Postfix
	}
	tests := []struct {
		name string
		args args
		want *Postfix
	}{
		{
			name: "Nil",
			args: args{
				Input: nil,
			},
			want: nil,
		},
		{
			name: "Empty",
			args: args{
				Input: &Postfix{},
			},
			want: &Postfix{},
		},

		{
			name: "Left",
			args: args{
				Input: &Postfix{
					Left: &Primary{},
				},
			},
			want: &Postfix{
				Left: &Primary{},
			},
		},
		{
			name: "Right",
			args: args{
				Input: &Postfix{
					Right: &Expr{},
				},
			},
			want: &Postfix{
				Right: &Expr{},
			},
		},
		{
			name: "Left and Right",
			args: args{
				Input: &Postfix{
					Left:  &Primary{},
					Right: &Expr{},
				},
			},
			want: &Postfix{
				Left:  &Primary{},
				Right: &Expr{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.Input.Clone())
		})
	}
}

func TestPrimary_Clone(t *testing.T) {
	type args struct {
		Input *Primary
	}
	tests := []struct {
		name string
		args args
		want *Primary
	}{
		{
			name: "Nil",
			args: args{
				Input: nil,
			},
			want: nil,
		},
		{
			name: "Empty",
			args: args{
				Input: &Primary{},
			},
			want: &Primary{},
		},

		{
			name: "Sub Expression",
			args: args{
				Input: &Primary{
					SubExpression: &Expr{},
				},
			},
			want: &Primary{
				SubExpression: &Expr{},
			},
		},
		{
			name: "Value",
			args: args{
				Input: &Primary{
					Value: &Value{},
				},
			},
			want: &Primary{
				Value: &Value{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.Input.Clone())
		})
	}
}
