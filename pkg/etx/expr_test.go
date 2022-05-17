package etx

import (
	"math/big"
	"testing"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/repr"
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
			name: "Invocation",
			args: args{
				Input: `foo(bar, baz)`,
			},
			wantErr: false,
			want: &Expr{
				Left: testBuildExprTree[*ExprConditional](t, &ExprInvocation{
					Ident: Ident{
						Parts: []string{"foo"},
					},
					InvocationParams: []InvocationParams{{
						Values: []Expr{
							*testBuildExprTree[*Expr](t, &Value{Ident: &Ident{Parts: []string{"bar"}}}),
							*testBuildExprTree[*Expr](t, &Value{Ident: &Ident{Parts: []string{"baz"}}}),
						},
					}},
				}),
			},
		},
		{
			name: "Dot Invocation",
			args: args{
				Input: `foo.bar(baz, qux)`,
			},
			wantErr: false,
			want: &Expr{
				Left: testBuildExprTree[*ExprConditional](t, &ExprInvocation{
					Ident: Ident{
						Parts: []string{
							"foo",
							"bar",
						},
					},
					InvocationParams: []InvocationParams{{
						[]Expr{
							*testBuildExprTree[*Expr](t, &Value{Ident: &Ident{Parts: []string{"baz"}}}),
							*testBuildExprTree[*Expr](t, &Value{Ident: &Ident{Parts: []string{"qux"}}}),
						},
					}},
				}),
			},
		},
		{
			name: "Monadic invocation",
			args: args{
				Input: `foo(bar)(baz)`,
			},
			wantErr: false,
			want: &Expr{
				Left: testBuildExprTree[*ExprConditional](t, &ExprInvocation{
					Ident: Ident{
						Parts: []string{"foo"},
					},
					InvocationParams: []InvocationParams{
						{
							[]Expr{
								*testBuildExprTree[*Expr](t, &Value{Ident: &Ident{Parts: []string{"bar"}}}),
							},
						},
						{
							[]Expr{
								*testBuildExprTree[*Expr](t, &Value{Ident: &Ident{Parts: []string{"baz"}}}),
							},
						},
					},
				}),
			},
		},

		{
			name: "If with empty body",
			args: args{
				Input: `if foo { }`,
			},
			wantErr: false,
			want: &Expr{
				If: &ExprIf{
					Condition: *testBuildExprTree[*ExprLogicalOr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
					Left:      nil,
					Right:     nil,
				},
			},
		},
		{
			name: "If",
			args: args{
				Input: `if foo { 1 }`,
			},
			wantErr: false,
			want: &Expr{
				If: &ExprIf{
					Condition: *testBuildExprTree[*ExprLogicalOr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
					Left:      testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Right:     nil,
				},
			},
		},
		{
			name: "If/Else",
			args: args{
				Input: `if foo { 1 } else { 2 }`,
			},
			wantErr: false,
			want: &Expr{
				If: &ExprIf{
					Condition: *testBuildExprTree[*ExprLogicalOr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
					Left:      testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Right:     testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			},
		},

		{
			name: "Switch - Empty",
			args: args{
				Input: `
switch foo {}`[1:],
			},
			wantErr: false,
			want: &Expr{
				Switch: &ExprSwitch{
					Selector: *testBuildExprTree[*ExprLogicalOr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
					Cases:    nil,
				},
			},
		},
		{
			name: "Switch - One single case, no default",
			args: args{
				Input: `
switch foo {
	case 1: { 2 }
}`[1:],
			},
			wantErr: false,
			want: &Expr{
				Switch: &ExprSwitch{
					Selector: *testBuildExprTree[*ExprLogicalOr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
					Cases: []*ExprCase{
						{
							Conditions: []ExprLogicalOr{
								*testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
							},
							Expr: testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(2)}}),
						},
					},
				},
			},
		},
		{
			name: "Switch - One multiple case, no default",
			args: args{
				Input: `
switch foo {
	case 1, 2: { 3 }
}`[1:],
			},
			wantErr: false,
			want: &Expr{
				Switch: &ExprSwitch{
					Selector: *testBuildExprTree[*ExprLogicalOr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
					Cases: []*ExprCase{
						{
							Conditions: []ExprLogicalOr{
								*testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
								*testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(2)}}),
							},
							Expr: testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(3)}}),
						},
					},
				},
			},
		},
		{
			name: "Switch - Two single case, no default",
			args: args{
				Input: `
switch foo {
	case 1: { 2 }
	case 3: { 4 }
}`[1:],
			},
			wantErr: false,
			want: &Expr{
				Switch: &ExprSwitch{
					Selector: *testBuildExprTree[*ExprLogicalOr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
					Cases: []*ExprCase{
						{
							Conditions: []ExprLogicalOr{
								*testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
							},
							Expr: testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(2)}}),
						},
						{
							Conditions: []ExprLogicalOr{
								*testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(3)}}),
							},
							Expr: testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(4)}}),
						},
					},
				},
			},
		},
		{
			name: "Switch - Only default",
			args: args{
				Input: `
switch foo {
	default: { 3 }
}`[1:],
			},
			wantErr: false,
			want: &Expr{
				Switch: &ExprSwitch{
					Selector: *testBuildExprTree[*ExprLogicalOr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
					Cases: []*ExprCase{
						{
							Conditions: nil,
							Default:    true,
							Expr:       testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(3)}}),
						},
					},
				},
			},
		},
		{
			name: "Switch - One single case, default",
			args: args{
				Input: `
switch foo {
	case 1: { 2 }
	default: { 3 }
}`[1:],
			},
			wantErr: false,
			want: &Expr{
				Switch: &ExprSwitch{
					Selector: *testBuildExprTree[*ExprLogicalOr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
					Cases: []*ExprCase{
						{
							Conditions: []ExprLogicalOr{
								*testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
							},
							Expr: testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(2)}}),
						},
						{
							Conditions: nil,
							Default:    true,
							Expr:       testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(3)}}),
						},
					},
					// Default: testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(3)}}),
				},
			},
		},

		{
			name: "Ternary",
			args: args{
				Input: "1 ? 2 : 3",
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprConditional{
					Condition:   *testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
					ConditionOp: true,
					TrueExpr:    testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(2)}}),
					FalseExpr:   testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(3)}}),
				},
			),
		},

		{
			name: "Logical OR - no spaces",
			args: args{
				Input: `1||2`,
			},
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprLogicalOr{
					Left:  *testBuildExprTree[*ExprLogicalAnd](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpLogicalOr,
					Right: testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprLogicalOr{
					Left:  *testBuildExprTree[*ExprLogicalAnd](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpLogicalOr,
					Right: testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprLogicalAnd{
					Left:  *testBuildExprTree[*ExprBitwiseOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpLogicalAnd,
					Right: testBuildExprTree[*ExprLogicalAnd](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprLogicalAnd{
					Left:  *testBuildExprTree[*ExprBitwiseOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpLogicalAnd,
					Right: testBuildExprTree[*ExprLogicalAnd](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprBitwiseOr{
					Left:  *testBuildExprTree[*ExprBitwiseXor](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpBitwiseOr,
					Right: testBuildExprTree[*ExprBitwiseOr](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprBitwiseOr{
					Left:  *testBuildExprTree[*ExprBitwiseXor](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpBitwiseOr,
					Right: testBuildExprTree[*ExprBitwiseOr](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprBitwiseXor{
					Left:  *testBuildExprTree[*ExprBitwiseAnd](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpBitwiseXOr,
					Right: testBuildExprTree[*ExprBitwiseXor](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprBitwiseXor{
					Left:  *testBuildExprTree[*ExprBitwiseAnd](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpBitwiseXOr,
					Right: testBuildExprTree[*ExprBitwiseXor](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprBitwiseAnd{
					Left:  *testBuildExprTree[*ExprEquality](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpBitwiseAnd,
					Right: testBuildExprTree[*ExprBitwiseAnd](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprBitwiseAnd{
					Left:  *testBuildExprTree[*ExprEquality](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpBitwiseAnd,
					Right: testBuildExprTree[*ExprBitwiseAnd](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprEquality{
					Left:  *testBuildExprTree[*ExprRelational](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpEqual,
					Right: testBuildExprTree[*ExprEquality](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprEquality{
					Left:  *testBuildExprTree[*ExprRelational](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpEqual,
					Right: testBuildExprTree[*ExprEquality](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprEquality{
					Left:  *testBuildExprTree[*ExprRelational](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpNotEqual,
					Right: testBuildExprTree[*ExprEquality](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprEquality{
					Left:  *testBuildExprTree[*ExprRelational](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpNotEqual,
					Right: testBuildExprTree[*ExprEquality](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprRelational{
					Left:  *testBuildExprTree[*ExprShift](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpMore,
					Right: testBuildExprTree[*ExprRelational](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprRelational{
					Left:  *testBuildExprTree[*ExprShift](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpMore,
					Right: testBuildExprTree[*ExprRelational](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprRelational{
					Left:  *testBuildExprTree[*ExprShift](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpLess,
					Right: testBuildExprTree[*ExprRelational](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprRelational{
					Left:  *testBuildExprTree[*ExprShift](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpLess,
					Right: testBuildExprTree[*ExprRelational](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprRelational{
					Left:  *testBuildExprTree[*ExprShift](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpMoreOrEqual,
					Right: testBuildExprTree[*ExprRelational](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprRelational{
					Left:  *testBuildExprTree[*ExprShift](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpMoreOrEqual,
					Right: testBuildExprTree[*ExprRelational](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprRelational{
					Left:  *testBuildExprTree[*ExprShift](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpLessOrEqual,
					Right: testBuildExprTree[*ExprRelational](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprRelational{
					Left:  *testBuildExprTree[*ExprShift](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpLessOrEqual,
					Right: testBuildExprTree[*ExprRelational](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprShift{
					Left:  *testBuildExprTree[*ExprAdditive](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpBitwiseShiftLeft,
					Right: testBuildExprTree[*ExprShift](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprShift{
					Left:  *testBuildExprTree[*ExprAdditive](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpBitwiseShiftLeft,
					Right: testBuildExprTree[*ExprShift](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprShift{
					Left:  *testBuildExprTree[*ExprAdditive](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpBitwiseShiftRight,
					Right: testBuildExprTree[*ExprShift](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprShift{
					Left:  *testBuildExprTree[*ExprAdditive](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpBitwiseShiftRight,
					Right: testBuildExprTree[*ExprShift](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprAdditive{
					Left:  *testBuildExprTree[*ExprMultiplicative](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpPlus,
					Right: testBuildExprTree[*ExprAdditive](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprAdditive{
					Left:  *testBuildExprTree[*ExprMultiplicative](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpPlus,
					Right: testBuildExprTree[*ExprAdditive](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprAdditive{
					Left:  *testBuildExprTree[*ExprMultiplicative](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpMinus,
					Right: testBuildExprTree[*ExprAdditive](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprAdditive{
					Left:  *testBuildExprTree[*ExprMultiplicative](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpMinus,
					Right: testBuildExprTree[*ExprAdditive](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprMultiplicative{
					Left:  *testBuildExprTree[*ExprUnary](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpDivision,
					Right: testBuildExprTree[*ExprMultiplicative](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprMultiplicative{
					Left:  *testBuildExprTree[*ExprUnary](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpDivision,
					Right: testBuildExprTree[*ExprMultiplicative](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprMultiplicative{
					Left:  *testBuildExprTree[*ExprUnary](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpMultiplication,
					Right: testBuildExprTree[*ExprMultiplicative](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprMultiplicative{
					Left:  *testBuildExprTree[*ExprUnary](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpMultiplication,
					Right: testBuildExprTree[*ExprMultiplicative](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprMultiplicative{
					Left:  *testBuildExprTree[*ExprUnary](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpModulo,
					Right: testBuildExprTree[*ExprMultiplicative](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprMultiplicative{
					Left:  *testBuildExprTree[*ExprUnary](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpModulo,
					Right: testBuildExprTree[*ExprMultiplicative](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
				&ExprUnary{
					Op:    OpBitwiseNot,
					Right: *testBuildExprTree[*ExprPostfix](t, &Value{Number: &Number{big.NewFloat(1)}}),
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
				&ExprUnary{
					Op:    OpBitwiseNot,
					Right: *testBuildExprTree[*ExprPostfix](t, &Value{Number: &Number{big.NewFloat(1)}}),
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
				&ExprUnary{
					Op:    OpLogicalNot,
					Right: *testBuildExprTree[*ExprPostfix](t, &Value{Number: &Number{big.NewFloat(1)}}),
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
				&ExprUnary{
					Op:    OpLogicalNot,
					Right: *testBuildExprTree[*ExprPostfix](t, &Value{Number: &Number{big.NewFloat(1)}}),
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
				&ExprUnary{
					Op:    OpMinus,
					Right: *testBuildExprTree[*ExprPostfix](t, &Value{Number: &Number{big.NewFloat(1)}}),
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
				&ExprUnary{
					Op:    OpMinus,
					Right: *testBuildExprTree[*ExprPostfix](t, &Value{Number: &Number{big.NewFloat(1)}}),
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
				&ExprPostfix{
					Left:  *testBuildExprTree[*ExprPrimary](t, &Value{Number: &Number{big.NewFloat(1)}}),
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
				&ExprPostfix{
					Left:  *testBuildExprTree[*ExprPrimary](t, &Value{Number: &Number{big.NewFloat(1)}}),
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
				&ExprPrimary{
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
				&ExprPrimary{
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
				&ExprPrimary{
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
				&ExprAdditive{
					Left: *testBuildExprTree[*ExprMultiplicative](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:   "+",
					Right: &ExprAdditive{
						Left:  *testBuildExprTree[*ExprMultiplicative](t, &Value{Number: &Number{big.NewFloat(2)}}),
						Op:    "+",
						Right: testBuildExprTree[*ExprAdditive](t, &Value{Number: &Number{big.NewFloat(3)}}),
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

			parser := participle.MustBuild(
				&Exp{},
				participle.Lexer(lexer.MustStateful(lexRules(), lexer.InitialState(lexerExpr))),
				participle.Elide(TokenWhitespace),
			)

			res := &Exp{}
			err := parser.ParseString("", tt.args.Input, res)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if !assert.Equal(t, tt.want, res.Expr) {
				repr.Println(res.Expr)
			}
		})
	}
}

// /////////////////////////////////////

func TestExpr_String(t *testing.T) {
	type args struct {
		Input *Expr
	}
	tests := []struct {
		name        string
		description string
		args        args
		wantPanic   bool
		want        string
	}{
		{
			name:        "nil",
			description: "",
			args: args{
				Input: nil,
			},
			wantPanic: true,
		},
		{
			name:        "Empty",
			description: "",
			args: args{
				Input: &Expr{},
			},
			wantPanic: true,
			want:      "expression not set",
		},
		{
			name:        "Left",
			description: "",
			args: args{
				Input: &Expr{
					Left: testBuildExprTree[*ExprConditional](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			},
			want: "1",
		},
		{
			name:        "If",
			description: "",
			args: args{
				Input: &Expr{
					If: &ExprIf{
						Condition: *testBuildExprTree[*ExprLogicalOr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
						Left:      nil,
						Right:     nil,
					},
				},
			},
			want: `if foo { }`,
		},
		{
			name:        "Switch",
			description: "",
			args: args{
				Input: &Expr{
					Switch: &ExprSwitch{
						Selector: *testBuildExprTree[*ExprLogicalOr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
						Cases:    nil,
					},
				},
			},
			want: `switch foo { }`,
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

func TestExprIf_String(t *testing.T) {
	type args struct {
		Input *ExprIf
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
				Input: &ExprIf{},
			},
			want: "if null { }",
		},
		{
			name: "If with empty body",
			args: args{
				Input: &ExprIf{
					Condition: *testBuildExprTree[*ExprLogicalOr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
					Left:      nil,
					Right:     nil,
				},
			},
			want: `if foo { }`,
		},
		{
			name: "If",
			args: args{
				Input: &ExprIf{
					Condition: *testBuildExprTree[*ExprLogicalOr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
					Left:      testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Right:     nil,
				},
			},
			want: `
if foo {
	1
}`[1:],
		},
		{
			name: "If/Else",
			args: args{
				Input: &ExprIf{
					Condition: *testBuildExprTree[*ExprLogicalOr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
					Left:      testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Right:     testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			},
			want: `
if foo {
	1
} else {
	2
}`[1:],
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

func TestExprSwitch_String(t *testing.T) {
	type args struct {
		Input *ExprSwitch
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
			name: "Empty",
			args: args{
				Input: &ExprSwitch{},
			},
			want: `switch null { }`,
		},
		{
			name: "Only Selector",
			args: args{
				Input: &ExprSwitch{
					Selector: *testBuildExprTree[*ExprLogicalOr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
					Cases:    nil,
				},
			},
			want: `switch foo { }`,
		},
		{
			name: "One single case, no default",
			args: args{
				Input: &ExprSwitch{
					Selector: *testBuildExprTree[*ExprLogicalOr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
					Cases: []*ExprCase{
						{
							Conditions: []ExprLogicalOr{
								*testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
							},
							Expr: testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(2)}}),
						},
					},
				},
			},
			want: `
switch foo {
	case 1: {
		2
	}
}`[1:],
		},
		{
			name: "One multiple case, no default",
			args: args{
				Input: &ExprSwitch{
					Selector: *testBuildExprTree[*ExprLogicalOr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
					Cases: []*ExprCase{
						{
							Conditions: []ExprLogicalOr{
								*testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
								*testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(2)}}),
							},
							Expr: testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(3)}}),
						},
					},
				},
			},
			want: `
switch foo {
	case 1, 2: {
		3
	}
}`[1:],
		},
		{
			name: "Two single case, no default",
			args: args{
				Input: &ExprSwitch{
					Selector: *testBuildExprTree[*ExprLogicalOr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
					Cases: []*ExprCase{
						{
							Conditions: []ExprLogicalOr{
								*testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
							},
							Expr: testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(2)}}),
						},
						{
							Conditions: []ExprLogicalOr{
								*testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(3)}}),
							},
							Expr: testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(4)}}),
						},
					},
				},
			},
			want: `
switch foo {
	case 1: {
		2
	}
	case 3: {
		4
	}
}`[1:],
		},
		{
			name: "Only default",
			args: args{
				Input: &ExprSwitch{
					Selector: *testBuildExprTree[*ExprLogicalOr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
					Cases: []*ExprCase{
						{
							Conditions: nil,
							Default:    true,
							Expr:       testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(3)}}),
						},
					},
				},
			},
			want: `
switch foo {
	default: {
		3
	}
}`[1:],
		},
		{
			name: "One single case, default",
			args: args{
				Input: &ExprSwitch{
					Selector: *testBuildExprTree[*ExprLogicalOr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
					Cases: []*ExprCase{
						{
							Conditions: []ExprLogicalOr{
								*testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
							},
							Expr: testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(2)}}),
						},
						{
							Conditions: nil,
							Default:    true,
							Expr:       testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(3)}}),
						},
					},
					// Default: testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(3)}}),
				},
			},
			want: `
switch foo {
	case 1: {
		2
	}
	default: {
		3
	}
}`[1:],
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

func TestExprCase_String(t *testing.T) {
	type args struct {
		Input *ExprCase
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
			description: "No condition and not default",
			args: args{
				Input: &ExprCase{},
			},
			wantPanic: true,
			want:      "non-default case statement without condition",
		},
		{
			name: "One condition",
			args: args{
				Input: &ExprCase{
					Conditions: []ExprLogicalOr{
						*testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
					},
					Default: false,
					Expr:    testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			},
			want: `
case 1: {
	2
}`[1:],
		},
		{
			name: "Several condition",
			args: args{
				Input: &ExprCase{
					Conditions: []ExprLogicalOr{
						*testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
						*testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(2)}}),
					},
					Default: false,
					Expr:    testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(3)}}),
				},
			},
			want: `
case 1, 2: {
	3
}`[1:],
		},
		{
			name: "Default",
			args: args{
				Input: &ExprCase{
					Conditions: nil,
					Default:    true,
					Expr:       testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(3)}}),
				},
			},
			want: `
default: {
	3
}`[1:],
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

func TestConditional_String(t *testing.T) {
	type args struct {
		Input *ExprConditional
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
			name: "Only Condition",
			args: args{
				Input: &ExprConditional{
					Condition: *testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			},
			want: "1",
		},
		{
			name: "All expressions but not operators",
			args: args{
				Input: &ExprConditional{
					Condition: *testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
					TrueExpr:  testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(2)}}),
					FalseExpr: testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(3)}}),
				},
			},
			want: "1",
		},
		{
			name: "Operator with neither True nor False",
			args: args{
				Input: &ExprConditional{
					Condition:   *testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
					ConditionOp: true,
				},
			},
			want: "1 ? null : null",
		},
		{
			name: "Only True expression",
			args: args{
				Input: &ExprConditional{
					Condition:   *testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
					ConditionOp: true,
					TrueExpr:    testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(2)}}),
				},
			},
			want: "1 ? 2 : null",
		},
		{
			name: "Only False expression",
			args: args{
				Input: &ExprConditional{
					Condition:   *testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
					ConditionOp: true,
					FalseExpr:   testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(3)}}),
				},
			},
			want: "1 ? null : 3",
		},
		{
			name:        "All parts",
			description: "Both sides of the condition must be present.",
			args: args{
				Input: &ExprConditional{
					Condition:   *testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
					ConditionOp: true,
					TrueExpr:    testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(2)}}),
					FalseExpr:   testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(3)}}),
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
		Input *ExprLogicalOr
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
			name: "Empty",
			args: args{
				Input: &ExprLogicalOr{},
			},
			want: TokenNull,
		},
		{
			name: "Left",
			args: args{
				Input: &ExprLogicalOr{
					Left: *testBuildExprTree[*ExprLogicalAnd](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			args: args{
				Input: &ExprLogicalOr{
					Left: ExprLogicalAnd{},
					Op:   OpLogicalOr,
				},
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			args: args{
				Input: &ExprLogicalOr{
					Left:  *testBuildExprTree[*ExprLogicalAnd](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpLogicalOr,
					Right: testBuildExprTree[*ExprLogicalOr](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
		Input *ExprLogicalAnd
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
			name: "Empty",
			args: args{
				Input: &ExprLogicalAnd{},
			},
			want: TokenNull,
		},
		{
			name: "Left",
			args: args{
				Input: &ExprLogicalAnd{
					Left: *testBuildExprTree[*ExprBitwiseOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			args: args{
				Input: &ExprLogicalAnd{
					Left: ExprBitwiseOr{},
					Op:   OpLogicalAnd,
				},
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			args: args{
				Input: &ExprLogicalAnd{
					Left:  *testBuildExprTree[*ExprBitwiseOr](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpLogicalAnd,
					Right: testBuildExprTree[*ExprLogicalAnd](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
		Input *ExprBitwiseOr
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
			name: "Empty",
			args: args{
				Input: &ExprBitwiseOr{},
			},
			want: TokenNull,
		},
		{
			name: "Left",
			args: args{
				Input: &ExprBitwiseOr{
					Left: *testBuildExprTree[*ExprBitwiseXor](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			args: args{
				Input: &ExprBitwiseOr{
					Left: ExprBitwiseXor{},
					Op:   OpBitwiseOr,
				},
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			args: args{
				Input: &ExprBitwiseOr{
					Left:  *testBuildExprTree[*ExprBitwiseXor](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpBitwiseOr,
					Right: testBuildExprTree[*ExprBitwiseOr](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
		Input *ExprBitwiseXor
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
			name: "Empty",
			args: args{
				Input: &ExprBitwiseXor{},
			},
			want: TokenNull,
		},
		{
			name: "Left",
			args: args{
				Input: &ExprBitwiseXor{
					Left: *testBuildExprTree[*ExprBitwiseAnd](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			args: args{
				Input: &ExprBitwiseXor{
					Left: ExprBitwiseAnd{},
					Op:   OpBitwiseXOr,
				},
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			args: args{
				Input: &ExprBitwiseXor{
					Left:  *testBuildExprTree[*ExprBitwiseAnd](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpBitwiseXOr,
					Right: testBuildExprTree[*ExprBitwiseXor](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
		Input *ExprBitwiseAnd
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
			name: "Empty",
			args: args{
				Input: &ExprBitwiseAnd{},
			},
			want: TokenNull,
		},
		{
			name: "Left",
			args: args{
				Input: &ExprBitwiseAnd{
					Left: *testBuildExprTree[*ExprEquality](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			args: args{
				Input: &ExprBitwiseAnd{
					Left: ExprEquality{},
					Op:   OpBitwiseAnd,
				},
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			args: args{
				Input: &ExprBitwiseAnd{
					Left:  *testBuildExprTree[*ExprEquality](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpBitwiseAnd,
					Right: testBuildExprTree[*ExprBitwiseAnd](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
		Input *ExprEquality
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
			name: "Empty",
			args: args{
				Input: &ExprEquality{},
			},
			want: TokenNull,
		},
		{
			name: "Left",
			args: args{
				Input: &ExprEquality{
					Left: *testBuildExprTree[*ExprRelational](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			args: args{
				Input: &ExprEquality{
					Left: ExprRelational{},
					Op:   OpEqual,
				},
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			args: args{
				Input: &ExprEquality{
					Left:  *testBuildExprTree[*ExprRelational](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpEqual,
					Right: testBuildExprTree[*ExprEquality](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
		Input *ExprRelational
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
			name: "Empty",
			args: args{
				Input: &ExprRelational{},
			},
			wantPanic: false,
			want:      TokenNull,
		},
		{
			name: "Left",
			args: args{
				Input: &ExprRelational{
					Left: *testBuildExprTree[*ExprShift](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			args: args{
				Input: &ExprRelational{
					Left: ExprShift{},
					Op:   OpMoreOrEqual,
				},
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			args: args{
				Input: &ExprRelational{
					Left:  *testBuildExprTree[*ExprShift](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpMoreOrEqual,
					Right: testBuildExprTree[*ExprRelational](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
		Input *ExprShift
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
			name: "Empty",
			args: args{
				Input: &ExprShift{},
			},
			want: TokenNull,
		},
		{
			name: "Left",
			args: args{
				Input: &ExprShift{
					Left: *testBuildExprTree[*ExprAdditive](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			args: args{
				Input: &ExprShift{
					Left: ExprAdditive{},
					Op:   OpBitwiseShiftRight,
				},
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			args: args{
				Input: &ExprShift{
					Left:  *testBuildExprTree[*ExprAdditive](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpBitwiseShiftRight,
					Right: testBuildExprTree[*ExprShift](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
		Input *ExprAdditive
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
			name: "Empty",
			args: args{
				Input: &ExprAdditive{},
			},
			want: TokenNull,
		},
		{
			name: "Left",
			args: args{
				Input: &ExprAdditive{
					Left: *testBuildExprTree[*ExprMultiplicative](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			args: args{
				Input: &ExprAdditive{
					Left: ExprMultiplicative{},
					Op:   OpPlus,
				},
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			args: args{
				Input: &ExprAdditive{
					Left:  *testBuildExprTree[*ExprMultiplicative](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpPlus,
					Right: testBuildExprTree[*ExprAdditive](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
		Input *ExprMultiplicative
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
			name: "Empty",
			args: args{
				Input: &ExprMultiplicative{},
			},
			want: TokenNull,
		},
		{
			name: "Left",
			args: args{
				Input: &ExprMultiplicative{
					Left: *testBuildExprTree[*ExprUnary](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			args: args{
				Input: &ExprMultiplicative{
					Left: *testBuildExprTree[*ExprUnary](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
					Op:   OpMultiplication,
				},
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			args: args{
				Input: &ExprMultiplicative{
					Left:  *testBuildExprTree[*ExprUnary](t, &Value{Number: &Number{big.NewFloat(1)}}),
					Op:    OpMultiplication,
					Right: testBuildExprTree[*ExprMultiplicative](t, &Value{Number: &Number{big.NewFloat(2)}}),
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
		Input *ExprUnary
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
				Input: &ExprUnary{},
			},
			want: TokenNull,
		},
		{
			name: "Right",
			args: args{
				Input: &ExprUnary{
					Right: *testBuildExprTree[*ExprPostfix](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
				},
			},
			want: "foo",
		},
		{
			name: "Operator",
			args: args{
				Input: &ExprUnary{
					Op:    OpBitwiseNot,
					Right: *testBuildExprTree[*ExprPostfix](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
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
		Input *ExprPostfix
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
			name: "Empty",
			args: args{
				Input: &ExprPostfix{},
			},
			want: TokenNull,
		},
		{
			name: "Left",
			args: args{
				Input: &ExprPostfix{
					Left: *testBuildExprTree[*ExprPrimary](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
				},
			},
			want: "foo",
		},
		{
			name: "Both",
			args: args{
				Input: &ExprPostfix{
					Left:  *testBuildExprTree[*ExprPrimary](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
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
		Input *ExprPrimary
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
				Input: &ExprPrimary{},
			},
			want: TokenNull,
		},
		{
			name: "Sub Expression",
			args: args{
				Input: &ExprPrimary{
					SubExpression: testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(1)}}),
				},
			},
			want: "1",
		},
		{
			name: "Value",
			args: args{
				Input: &ExprPrimary{
					Value: &Value{Number: &Number{big.NewFloat(1)}},
				},
			},
			want: "1",
		},
		{
			name: "Invocation",
			args: args{
				Input: &ExprPrimary{
					Value: &Value{Number: &Number{big.NewFloat(1)}},
					Invocation: &ExprInvocation{
						Ident: Ident{
							Parts: []string{"foo"},
						},
						InvocationParams: []InvocationParams{{
							[]Expr{
								*testBuildExprTree[*Expr](t, &Value{Ident: &Ident{Parts: []string{"bar"}}}),
								*testBuildExprTree[*Expr](t, &Value{Ident: &Ident{Parts: []string{"baz"}}}),
							},
						}},
					},
				},
			},
			want: "foo(bar, baz)",
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

func TestInvocation_String(t *testing.T) {
	type args struct {
		Input *ExprInvocation
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
				Input: &ExprInvocation{},
			},
			want: "()",
		},
		{
			name: "No parameters",
			args: args{
				Input: &ExprInvocation{
					Ident: Ident{
						Parts: []string{
							"foo",
						},
					},
				},
			},
			want: "foo()",
		},
		{
			name: "One parameter",
			args: args{
				Input: &ExprInvocation{
					Ident: Ident{
						Parts: []string{
							"foo",
						},
					},
					InvocationParams: []InvocationParams{{
						[]Expr{
							*testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(1)}}),
						},
					}},
				},
			},
			want: "foo(1)",
		},
		{
			name: "Several parameters",
			args: args{
				Input: &ExprInvocation{
					Ident: Ident{
						Parts: []string{
							"foo",
						},
					},
					InvocationParams: []InvocationParams{{
						[]Expr{
							*testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(1)}}),
							*testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(2)}}),
						},
					}},
				},
			},
			want: "foo(1, 2)",
		},
		{
			name: "Monadic invocation",
			args: args{
				Input: &ExprInvocation{
					Ident: Ident{
						Parts: []string{
							"foo",
						},
					},
					InvocationParams: []InvocationParams{
						{
							[]Expr{
								*testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(1)}}),
							},
						},
						{
							[]Expr{
								*testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(2)}}),
							},
						},
					},
				},
			},
			want: "foo(1)(2)",
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
					Left: &ExprConditional{},
				},
			},
			want: &Expr{
				Left: &ExprConditional{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.Input.Clone())
		})
	}
}

func TestIf_Clone(t *testing.T) {
	type args struct {
		Input *ExprIf
	}
	tests := []struct {
		name string
		args args
		want *ExprIf
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
				Input: &ExprIf{},
			},
			want: &ExprIf{},
		},

		{
			name: "Condition",
			args: args{
				Input: &ExprIf{
					Condition: ExprLogicalOr{},
				},
			},
			want: &ExprIf{
				Condition: ExprLogicalOr{},
			},
		},
		{
			name: "Left",
			args: args{
				Input: &ExprIf{
					Left: &Expr{},
				},
			},
			want: &ExprIf{
				Left: &Expr{},
			},
		},
		{
			name: "Right",
			args: args{
				Input: &ExprIf{
					Right: &Expr{},
				},
			},
			want: &ExprIf{
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

func TestSwitch_Clone(t *testing.T) {
	type args struct {
		Input *ExprSwitch
	}
	tests := []struct {
		name string
		args args
		want *ExprSwitch
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
				Input: &ExprSwitch{},
			},
			want: &ExprSwitch{},
		},

		{
			name: "Selector",
			args: args{
				Input: &ExprSwitch{
					Selector: ExprLogicalOr{},
				},
			},
			want: &ExprSwitch{
				Selector: ExprLogicalOr{},
			},
		},
		{
			name: "Cases",
			args: args{
				Input: &ExprSwitch{
					Cases: []*ExprCase{{}},
				},
			},
			want: &ExprSwitch{
				Cases: []*ExprCase{{}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.Input.Clone())
		})
	}
}

func TestCase_Clone(t *testing.T) {
	type args struct {
		Input *ExprCase
	}
	tests := []struct {
		name string
		args args
		want *ExprCase
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
				Input: &ExprCase{},
			},
			want: &ExprCase{},
		},

		{
			name: "Conditions",
			args: args{
				Input: &ExprCase{
					Conditions: []ExprLogicalOr{{}},
				},
			},
			want: &ExprCase{
				Conditions: []ExprLogicalOr{{}},
			},
		},
		{
			name: "Default",
			args: args{
				Input: &ExprCase{
					Default: true,
				},
			},
			want: &ExprCase{
				Default: true,
			},
		},
		{
			name: "Expr",
			args: args{
				Input: &ExprCase{
					Expr: &Expr{},
				},
			},
			want: &ExprCase{
				Expr: &Expr{},
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
		Input *ExprConditional
	}
	tests := []struct {
		name string
		args args
		want *ExprConditional
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
				Input: &ExprConditional{},
			},
			want: &ExprConditional{},
		},

		{
			name: "With Data",
			args: args{
				Input: &ExprConditional{
					Condition:   ExprLogicalOr{},
					ConditionOp: true,
					TrueExpr:    &ExprLogicalOr{},
					FalseExpr:   &ExprLogicalOr{},
				},
			},
			want: &ExprConditional{
				Condition:   ExprLogicalOr{},
				ConditionOp: true,
				TrueExpr:    &ExprLogicalOr{},
				FalseExpr:   &ExprLogicalOr{},
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
		Input *ExprLogicalAnd
	}
	tests := []struct {
		name string
		args args
		want *ExprLogicalAnd
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
				Input: &ExprLogicalAnd{},
			},
			want: &ExprLogicalAnd{},
		},

		{
			name: "Left",
			args: args{
				Input: &ExprLogicalAnd{
					Left: ExprBitwiseOr{},
				},
			},
			want: &ExprLogicalAnd{
				Left: ExprBitwiseOr{},
			},
		},
		{
			name: "Both",
			args: args{
				Input: &ExprLogicalAnd{
					Left:  ExprBitwiseOr{},
					Op:    OpLogicalAnd,
					Right: &ExprLogicalAnd{},
				},
			},
			want: &ExprLogicalAnd{
				Left:  ExprBitwiseOr{},
				Op:    OpLogicalAnd,
				Right: &ExprLogicalAnd{},
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
		Input *ExprBitwiseAnd
	}
	tests := []struct {
		name string
		args args
		want *ExprBitwiseAnd
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
				Input: &ExprBitwiseAnd{},
			},
			want: &ExprBitwiseAnd{},
		},

		{
			name: "Left",
			args: args{
				Input: &ExprBitwiseAnd{
					Left: ExprEquality{},
				},
			},
			want: &ExprBitwiseAnd{
				Left: ExprEquality{},
			},
		},
		{
			name: "Both",
			args: args{
				Input: &ExprBitwiseAnd{
					Left:  ExprEquality{},
					Op:    OpBitwiseAnd,
					Right: &ExprBitwiseAnd{},
				},
			},
			want: &ExprBitwiseAnd{
				Left:  ExprEquality{},
				Op:    OpBitwiseAnd,
				Right: &ExprBitwiseAnd{},
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
		Input *ExprEquality
	}
	tests := []struct {
		name string
		args args
		want *ExprEquality
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
				Input: &ExprEquality{},
			},
			want: &ExprEquality{},
		},

		{
			name: "Left",
			args: args{
				Input: &ExprEquality{
					Left: ExprRelational{},
				},
			},
			want: &ExprEquality{
				Left: ExprRelational{},
			},
		},
		{
			name: "Both",
			args: args{
				Input: &ExprEquality{
					Left:  ExprRelational{},
					Op:    OpNotEqual,
					Right: &ExprEquality{},
				},
			},
			want: &ExprEquality{
				Left:  ExprRelational{},
				Op:    OpNotEqual,
				Right: &ExprEquality{},
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
		Input *ExprRelational
	}
	tests := []struct {
		name string
		args args
		want *ExprRelational
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
				Input: &ExprRelational{},
			},
			want: &ExprRelational{},
		},

		{
			name: "Left",
			args: args{
				Input: &ExprRelational{
					Left: ExprShift{},
				},
			},
			want: &ExprRelational{
				Left: ExprShift{},
			},
		},
		{
			name: "Both",
			args: args{
				Input: &ExprRelational{
					Left:  ExprShift{},
					Op:    OpLessOrEqual,
					Right: &ExprRelational{},
				},
			},
			want: &ExprRelational{
				Left:  ExprShift{},
				Op:    OpLessOrEqual,
				Right: &ExprRelational{},
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
		Input *ExprShift
	}
	tests := []struct {
		name string
		args args
		want *ExprShift
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
				Input: &ExprShift{},
			},
			want: &ExprShift{},
		},

		{
			name: "Left",
			args: args{
				Input: &ExprShift{
					Left: ExprAdditive{},
				},
			},
			want: &ExprShift{
				Left: ExprAdditive{},
			},
		},
		{
			name: "Both",
			args: args{
				Input: &ExprShift{
					Left:  ExprAdditive{},
					Op:    OpBitwiseShiftRight,
					Right: &ExprShift{},
				},
			},
			want: &ExprShift{
				Left:  ExprAdditive{},
				Op:    OpBitwiseShiftRight,
				Right: &ExprShift{},
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
		Input *ExprAdditive
	}
	tests := []struct {
		name string
		args args
		want *ExprAdditive
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
				Input: &ExprAdditive{},
			},
			want: &ExprAdditive{},
		},

		{
			name: "Left",
			args: args{
				Input: &ExprAdditive{
					Left: ExprMultiplicative{},
				},
			},
			want: &ExprAdditive{
				Left: ExprMultiplicative{},
			},
		},
		{
			name: "Both",
			args: args{
				Input: &ExprAdditive{
					Left:  ExprMultiplicative{},
					Op:    OpMinus,
					Right: &ExprAdditive{},
				},
			},
			want: &ExprAdditive{
				Left:  ExprMultiplicative{},
				Op:    OpMinus,
				Right: &ExprAdditive{},
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
		Input *ExprMultiplicative
	}
	tests := []struct {
		name string
		args args
		want *ExprMultiplicative
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
				Input: &ExprMultiplicative{},
			},
			want: &ExprMultiplicative{},
		},

		{
			name: "Left",
			args: args{
				Input: &ExprMultiplicative{
					Left: ExprUnary{},
				},
			},
			want: &ExprMultiplicative{
				Left: ExprUnary{},
			},
		},
		{
			name: "Both",
			args: args{
				Input: &ExprMultiplicative{
					Left:  ExprUnary{},
					Op:    OpModulo,
					Right: &ExprMultiplicative{},
				},
			},
			want: &ExprMultiplicative{
				Left:  ExprUnary{},
				Op:    OpModulo,
				Right: &ExprMultiplicative{},
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
		Input *ExprUnary
	}
	tests := []struct {
		name string
		args args
		want *ExprUnary
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
				Input: &ExprUnary{},
			},
			want: &ExprUnary{},
		},

		{
			name: "Operator",
			args: args{
				Input: &ExprUnary{
					Op:    OpMinus,
					Right: ExprPostfix{},
				},
			},
			want: &ExprUnary{
				Op:    OpMinus,
				Right: ExprPostfix{},
			},
		},
		{
			name: "No Operator",
			args: args{
				Input: &ExprUnary{
					Right: ExprPostfix{},
				},
			},
			want: &ExprUnary{
				Right: ExprPostfix{},
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
		Input *ExprPostfix
	}
	tests := []struct {
		name string
		args args
		want *ExprPostfix
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
				Input: &ExprPostfix{},
			},
			want: &ExprPostfix{},
		},

		{
			name: "Left",
			args: args{
				Input: &ExprPostfix{
					Left: ExprPrimary{},
				},
			},
			want: &ExprPostfix{
				Left: ExprPrimary{},
			},
		},
		{
			name: "Right",
			args: args{
				Input: &ExprPostfix{
					Right: &Expr{},
				},
			},
			want: &ExprPostfix{
				Right: &Expr{},
			},
		},
		{
			name: "Left and Right",
			args: args{
				Input: &ExprPostfix{
					Left:  ExprPrimary{},
					Right: &Expr{},
				},
			},
			want: &ExprPostfix{
				Left:  ExprPrimary{},
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
		Input *ExprPrimary
	}
	tests := []struct {
		name string
		args args
		want *ExprPrimary
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
				Input: &ExprPrimary{},
			},
			want: &ExprPrimary{},
		},

		{
			name: "Sub Expression",
			args: args{
				Input: &ExprPrimary{
					SubExpression: &Expr{},
				},
			},
			want: &ExprPrimary{
				SubExpression: &Expr{},
			},
		},
		{
			name: "Invocation",
			args: args{
				Input: &ExprPrimary{
					Invocation: &ExprInvocation{},
				},
			},
			want: &ExprPrimary{
				Invocation: &ExprInvocation{},
			},
		},
		{
			name: "Value",
			args: args{
				Input: &ExprPrimary{
					Value: &Value{},
				},
			},
			want: &ExprPrimary{
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
