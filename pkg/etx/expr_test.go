package etx

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpr_Parsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Expr
	}{
		{
			name:    "Invocation - no parameters",
			input:   `foo()`,
			wantErr: false,
			want: &Expr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Left: testBuildExprTree[*ExprConditional](t, &ExprInvocation{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Ident: Ident{
						Parts: []string{"foo"},
					},
					Monads: []*ExprInvocationParams{{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
					}},
				}),
			},
		},
		{
			name:    "Invocation - parameters",
			input:   `foo(bar, baz)`,
			wantErr: false,
			want: &Expr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Left: testBuildExprTree[*ExprConditional](t, &ExprInvocation{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Ident: Ident{
						Parts: []string{"foo"},
					},
					Monads: []*ExprInvocationParams{{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Values: []*Expr{
							testBuildExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
								Ident:   &Ident{Parts: []string{"bar"}},
							}),
							testBuildExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 1, Column: 10}},
								Ident:   &Ident{Parts: []string{"baz"}},
							}),
						},
					}},
				}),
			},
		},
		{
			name:    "Dot Invocation",
			input:   `foo.bar(baz, qux)`,
			wantErr: false,
			want: &Expr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Left: testBuildExprTree[*ExprConditional](t, &ExprInvocation{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Ident: Ident{
						Parts: []string{
							"foo",
							"bar",
						},
					},
					Monads: []*ExprInvocationParams{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
							Values: []*Expr{
								testBuildExprTree[*Expr](t, &Value{
									ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
									Ident:   &Ident{Parts: []string{"baz"}},
								}),
								testBuildExprTree[*Expr](t, &Value{
									ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 1, Column: 14}},
									Ident:   &Ident{Parts: []string{"qux"}},
								}),
							},
						},
					},
				}),
			},
		},
		{
			name:    "Monadic invocation",
			input:   `foo(bar)(baz)`,
			wantErr: false,
			want: &Expr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Left: testBuildExprTree[*ExprConditional](t, &ExprInvocation{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Ident: Ident{
						Parts: []string{"foo"},
					},
					Monads: []*ExprInvocationParams{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
							Values: []*Expr{
								testBuildExprTree[*Expr](t, &Value{
									ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
									Ident:   &Ident{Parts: []string{"bar"}},
								}),
							},
						},
						{
							ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 1, Column: 10}},
							Values: []*Expr{
								testBuildExprTree[*Expr](t, &Value{
									ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 1, Column: 10}},
									Ident:   &Ident{Parts: []string{"baz"}},
								}),
							},
						},
					},
				}),
			},
		},
		{
			name:  "Dot reference on invocation",
			input: `foo().bar`,
			want: &Expr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Left: testBuildExprTree[*ExprConditional](t, &ExprInvocation{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Ident: Ident{
						Parts: []string{
							"foo",
						},
					},
					Monads: []*ExprInvocationParams{
						{ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}}},
					},
					Postfix: testBuildExprTree[*ExprPostfix](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
						Ident:   &Ident{Parts: []string{"bar"}},
					}),
				}),
			},
		},
		{
			name:  "Dot reference invocation on invocation",
			input: `foo().bar()`,
			want: &Expr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Left: testBuildExprTree[*ExprConditional](t, &ExprInvocation{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Ident: Ident{
						Parts: []string{
							"foo",
						},
					},
					Monads: []*ExprInvocationParams{
						{ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}}},
					},
					Postfix: testBuildExprTree[*ExprPostfix](t, &ExprInvocation{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
						Ident: Ident{
							Parts: []string{
								"bar",
							},
						},
						Monads: []*ExprInvocationParams{
							{ASTNode: ASTNode{Pos: Position{Offset: 10, Line: 1, Column: 11}}},
						},
					}),
				}),
			},
		},

		{
			name:    "If with empty body",
			input:   `if foo { }`,
			wantErr: false,
			want: &Expr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				If: &ExprIf{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Condition: *testBuildExprTree[*ExprLogicalOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
						Ident:   &Ident{Parts: []string{"foo"}},
					}),
					Left:  nil,
					Right: nil,
				},
			},
		},
		{
			name:    "If",
			input:   `if foo { 1 }`,
			wantErr: false,
			want: &Expr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				If: &ExprIf{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Condition: *testBuildExprTree[*ExprLogicalOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
						Ident:   &Ident{Parts: []string{"foo"}},
					}),
					Left: testBuildExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 1, Column: 10}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Right: nil,
				},
			},
		},
		{
			name:    "If/Else",
			input:   `if foo { 1 } else { 2 }`,
			wantErr: false,
			want: &Expr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				If: &ExprIf{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Condition: *testBuildExprTree[*ExprLogicalOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
						Ident:   &Ident{Parts: []string{"foo"}},
					}),
					Left: testBuildExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 1, Column: 10}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Right: testBuildExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 20, Line: 1, Column: 21}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			},
		},

		{
			name: "Switch - Empty",
			input: `
switch foo {}`[1:],
			wantErr: false,
			want: &Expr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Switch: &ExprSwitch{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Selector: *testBuildExprTree[*ExprLogicalOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
						Ident:   &Ident{Parts: []string{"foo"}},
					}),
					Cases: nil,
				},
			},
		},
		{
			name: "Switch - One single case, no default",
			input: `
switch foo {
	case 1: { 2 }
}`[1:],
			wantErr: false,
			want: &Expr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Switch: &ExprSwitch{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Selector: *testBuildExprTree[*ExprLogicalOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
						Ident:   &Ident{Parts: []string{"foo"}},
					}),
					Cases: []*ExprCase{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 14, Line: 2, Column: 2}},
							Conditions: []*ExprLogicalOr{
								testBuildExprTree[*ExprLogicalOr](t, &Value{
									ASTNode: ASTNode{Pos: Position{Offset: 19, Line: 2, Column: 7}},
									Number:  &ValueNumber{big.NewFloat(1), "1"},
								}),
							},
							Expr: testBuildExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 24, Line: 2, Column: 12}},
								Number:  &ValueNumber{big.NewFloat(2), "2"},
							}),
						},
					},
				},
			},
		},
		{
			name: "Switch - One multiple case, no default",
			input: `
switch foo {
	case 1, 2: { 3 }
}`[1:],
			wantErr: false,
			want: &Expr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Switch: &ExprSwitch{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Selector: *testBuildExprTree[*ExprLogicalOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
						Ident:   &Ident{Parts: []string{"foo"}},
					}),
					Cases: []*ExprCase{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 14, Line: 2, Column: 2}},
							Conditions: []*ExprLogicalOr{
								testBuildExprTree[*ExprLogicalOr](t, &Value{
									ASTNode: ASTNode{Pos: Position{Offset: 19, Line: 2, Column: 7}},
									Number:  &ValueNumber{big.NewFloat(1), "1"},
								}),
								testBuildExprTree[*ExprLogicalOr](t, &Value{
									ASTNode: ASTNode{Pos: Position{Offset: 22, Line: 2, Column: 10}},
									Number:  &ValueNumber{big.NewFloat(2), "2"},
								}),
							},
							Expr: testBuildExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 27, Line: 2, Column: 15}},
								Number:  &ValueNumber{big.NewFloat(3), "3"},
							}),
						},
					},
				},
			},
		},
		{
			name: "Switch - Two single case, no default",
			input: `
switch foo {
	case 1: { 2 }
	case 3: { 4 }
}`[1:],
			wantErr: false,
			want: &Expr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Switch: &ExprSwitch{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Selector: *testBuildExprTree[*ExprLogicalOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
						Ident:   &Ident{Parts: []string{"foo"}},
					}),
					Cases: []*ExprCase{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 14, Line: 2, Column: 2}},
							Conditions: []*ExprLogicalOr{
								testBuildExprTree[*ExprLogicalOr](t, &Value{
									ASTNode: ASTNode{Pos: Position{Offset: 19, Line: 2, Column: 7}},
									Number:  &ValueNumber{big.NewFloat(1), "1"},
								}),
							},
							Expr: testBuildExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 24, Line: 2, Column: 12}},
								Number:  &ValueNumber{big.NewFloat(2), "2"},
							}),
						},
						{
							ASTNode: ASTNode{Pos: Position{Offset: 29, Line: 3, Column: 2}},
							Conditions: []*ExprLogicalOr{
								testBuildExprTree[*ExprLogicalOr](t, &Value{
									ASTNode: ASTNode{Pos: Position{Offset: 34, Line: 3, Column: 7}},
									Number:  &ValueNumber{big.NewFloat(3), "3"},
								}),
							},
							Expr: testBuildExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 39, Line: 3, Column: 12}},
								Number:  &ValueNumber{big.NewFloat(4), "4"},
							}),
						},
					},
				},
			},
		},
		{
			name: "Switch - Only default",
			input: `
switch foo {
	default: { 3 }
}`[1:],
			wantErr: false,
			want: &Expr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Switch: &ExprSwitch{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Selector: *testBuildExprTree[*ExprLogicalOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
						Ident:   &Ident{Parts: []string{"foo"}},
					}),
					Cases: []*ExprCase{
						{
							ASTNode:    ASTNode{Pos: Position{Offset: 14, Line: 2, Column: 2}},
							Conditions: nil,
							Default:    true,
							Expr: testBuildExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 25, Line: 2, Column: 13}},
								Number:  &ValueNumber{big.NewFloat(3), "3"},
							}),
						},
					},
				},
			},
		},
		{
			name: "Switch - One single case, default",
			input: `
switch foo {
	case 1: { 2 }
	default: { 3 }
}`[1:],
			wantErr: false,
			want: &Expr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Switch: &ExprSwitch{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Selector: *testBuildExprTree[*ExprLogicalOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
						Ident:   &Ident{Parts: []string{"foo"}},
					}),
					Cases: []*ExprCase{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 14, Line: 2, Column: 2}},
							Conditions: []*ExprLogicalOr{
								testBuildExprTree[*ExprLogicalOr](t, &Value{
									ASTNode: ASTNode{Pos: Position{Offset: 19, Line: 2, Column: 7}},
									Number:  &ValueNumber{big.NewFloat(1), "1"},
								}),
							},
							Expr: testBuildExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 24, Line: 2, Column: 12}},
								Number:  &ValueNumber{big.NewFloat(2), "2"},
							}),
						},
						{
							ASTNode:    ASTNode{Pos: Position{Offset: 29, Line: 3, Column: 2}},
							Conditions: nil,
							Default:    true,
							Expr: testBuildExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 40, Line: 3, Column: 13}},
								Number:  &ValueNumber{big.NewFloat(3), "3"},
							}),
						},
					},
					// Default: testBuildExprTree[*Expr](t, &Value{ValueNumber: &ValueNumber{big.NewFloat(3), "3"}}),
				},
			},
		},

		{
			name:    "Ternary",
			input:   "1 ? 2 : 3",
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprConditional{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Condition: *testBuildExprTree[*ExprLogicalOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					ConditionOp: true,
					TrueExpr: testBuildExprTree[*ExprLogicalOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
					FalseExpr: testBuildExprTree[*ExprLogicalOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
						Number:  &ValueNumber{big.NewFloat(3), "3"},
					}),
				},
			),
		},

		{
			name:    "Logical OR - no spaces",
			input:   `1||2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprLogicalOr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprLogicalAnd](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpLogicalOr,
					Right: testBuildExprTree[*ExprLogicalOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},
		{
			name:    "Logical OR - spaces",
			input:   `1 || 2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprLogicalOr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprLogicalAnd](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpLogicalOr,
					Right: testBuildExprTree[*ExprLogicalOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},

		{
			name:    "Logical AND - no spaces",
			input:   `1&&2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprLogicalAnd{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprBitwiseOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpLogicalAnd,
					Right: testBuildExprTree[*ExprLogicalAnd](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},
		{
			name:    "Logical AND - spaces",
			input:   `1 && 2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprLogicalAnd{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprBitwiseOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpLogicalAnd,
					Right: testBuildExprTree[*ExprLogicalAnd](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},

		{
			name:    "Bitwise OR - no spaces",
			input:   `1|2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprBitwiseOr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprBitwiseXor](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpBitwiseOr,
					Right: testBuildExprTree[*ExprBitwiseOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},
		{
			name:    "Bitwise OR - spaces",
			input:   `1 | 2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprBitwiseOr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprBitwiseXor](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpBitwiseOr,
					Right: testBuildExprTree[*ExprBitwiseOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},

		{
			name:    "Bitwise XOR - no spaces",
			input:   `1^2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprBitwiseXor{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprBitwiseAnd](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpBitwiseXOr,
					Right: testBuildExprTree[*ExprBitwiseXor](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},
		{
			name:    "Bitwise XOR - spaces",
			input:   `1 ^ 2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprBitwiseXor{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprBitwiseAnd](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpBitwiseXOr,
					Right: testBuildExprTree[*ExprBitwiseXor](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},

		{
			name:    "Bitwise AND - no spaces",
			input:   `1&2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprBitwiseAnd{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprEquality](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpBitwiseAnd,
					Right: testBuildExprTree[*ExprBitwiseAnd](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},
		{
			name:    "Bitwise AND - spaces",
			input:   `1 & 2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprBitwiseAnd{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprEquality](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpBitwiseAnd,
					Right: testBuildExprTree[*ExprBitwiseAnd](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},

		{
			name:    "Equality - Equal - no spaces",
			input:   `1==2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprEquality{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprRelational](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpEqual,
					Right: testBuildExprTree[*ExprEquality](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},
		{
			name:    "Equality - Equal - spaces",
			input:   `1 == 2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprEquality{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprRelational](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpEqual,
					Right: testBuildExprTree[*ExprEquality](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},
		{
			name:    "Equality - Not Equal - no spaces",
			input:   `1!=2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprEquality{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprRelational](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpNotEqual,
					Right: testBuildExprTree[*ExprEquality](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},
		{
			name:    "Equality - Not Equal - spaces",
			input:   `1 != 2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprEquality{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprRelational](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpNotEqual,
					Right: testBuildExprTree[*ExprEquality](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},

		{
			name:    "Relational - More - no spaces",
			input:   `1>2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprRelational{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprShift](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpMore,
					Right: testBuildExprTree[*ExprRelational](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},
		{
			name:    "Relational - More - spaces",
			input:   `1 > 2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprRelational{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprShift](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpMore,
					Right: testBuildExprTree[*ExprRelational](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},
		{
			name:    "Relational - Less - no spaces",
			input:   `1<2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprRelational{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprShift](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpLess,
					Right: testBuildExprTree[*ExprRelational](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},
		{
			name:    "Relational - Less - spaces",
			input:   `1 < 2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprRelational{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprShift](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpLess,
					Right: testBuildExprTree[*ExprRelational](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},
		{
			name:    "Relational - More Or Equal - no spaces",
			input:   `1>=2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprRelational{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprShift](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpMoreOrEqual,
					Right: testBuildExprTree[*ExprRelational](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},
		{
			name:    "Relational - More Or Equal - spaces",
			input:   `1 >= 2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprRelational{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprShift](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpMoreOrEqual,
					Right: testBuildExprTree[*ExprRelational](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},
		{
			name:    "Relational - Less Or Equal - no spaces",
			input:   `1<=2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprRelational{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprShift](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpLessOrEqual,
					Right: testBuildExprTree[*ExprRelational](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},
		{
			name:    "Relational - Less Or Equal - spaces",
			input:   `1 <= 2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprRelational{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprShift](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpLessOrEqual,
					Right: testBuildExprTree[*ExprRelational](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},

		{
			name:    "Shift - Left - no spaces",
			input:   `1<<2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprShift{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprAdditive](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpBitwiseShiftLeft,
					Right: testBuildExprTree[*ExprShift](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},
		{
			name:    "Shift - Left - spaces",
			input:   `1 << 2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprShift{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprAdditive](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpBitwiseShiftLeft,
					Right: testBuildExprTree[*ExprShift](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},
		{
			name:    "Shift - Right - no spaces",
			input:   `1>>2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprShift{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprAdditive](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpBitwiseShiftRight,
					Right: testBuildExprTree[*ExprShift](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},
		{
			name:    "Shift - Right - spaces",
			input:   `1 >> 2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprShift{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprAdditive](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpBitwiseShiftRight,
					Right: testBuildExprTree[*ExprShift](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},

		{
			name:    "Additive - Plus - no spaces",
			input:   `1+2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprAdditive{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprMultiplicative](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpPlus,
					Right: testBuildExprTree[*ExprAdditive](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},
		{
			name:    "Additive - Plus - spaces",
			input:   `1 + 2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprAdditive{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprMultiplicative](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpPlus,
					Right: testBuildExprTree[*ExprAdditive](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},
		{
			name:    "Additive - Minus - no spaces",
			input:   `1-2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprAdditive{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprMultiplicative](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpMinus,
					Right: testBuildExprTree[*ExprAdditive](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},
		{
			name:    "Additive - Minus - spaces",
			input:   `1 - 2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprAdditive{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprMultiplicative](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpMinus,
					Right: testBuildExprTree[*ExprAdditive](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},

		{
			name:    "Multiplicative - Division - no spaces",
			input:   `1/2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprMultiplicative{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprUnary](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpDivision,
					Right: testBuildExprTree[*ExprMultiplicative](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},
		{
			name:    "Multiplicative - Division - spaces",
			input:   `1 / 2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprMultiplicative{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprUnary](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpDivision,
					Right: testBuildExprTree[*ExprMultiplicative](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},
		{
			name:    "Multiplicative - Multiplication - no spaces",
			input:   `1*2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprMultiplicative{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprUnary](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpMultiplication,
					Right: testBuildExprTree[*ExprMultiplicative](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},
		{
			name:    "Multiplicative - Multiplication - spaces",
			input:   `1 * 2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprMultiplicative{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprUnary](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpMultiplication,
					Right: testBuildExprTree[*ExprMultiplicative](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},
		{
			name:    "Multiplicative - Modulo - no spaces",
			input:   `1%2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprMultiplicative{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprUnary](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpModulo,
					Right: testBuildExprTree[*ExprMultiplicative](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},
		{
			name:    "Multiplicative - Modulo - spaces",
			input:   `1 % 2`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprMultiplicative{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprUnary](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: OpModulo,
					Right: testBuildExprTree[*ExprMultiplicative](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},

		{
			name:    "Unary - Bitwise NOT - no spaces",
			input:   `~1`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprUnary{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Op:      OpBitwiseNot,
					Right: *testBuildExprTree[*ExprPostfix](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
				},
			),
		},
		{
			name:    "Unary - Bitwise NOT - spaces",
			input:   `~ 1`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprUnary{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Op:      OpBitwiseNot,
					Right: *testBuildExprTree[*ExprPostfix](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
				},
			),
		},
		{
			name:    "Unary - Logical NOT - no spaces",
			input:   `!1`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprUnary{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Op:      OpLogicalNot,
					Right: *testBuildExprTree[*ExprPostfix](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
				},
			),
		},
		{
			name:    "Unary - Logical NOT - spaces",
			input:   `! 1`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprUnary{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Op:      OpLogicalNot,
					Right: *testBuildExprTree[*ExprPostfix](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
				},
			),
		},
		{
			name:    "Unary - Minus - no spaces",
			input:   `-1`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprUnary{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Op:      OpMinus,
					Right: *testBuildExprTree[*ExprPostfix](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
				},
			),
		},
		{
			name:    "Unary - Minus - spaces",
			input:   `- 1`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprUnary{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Op:      OpMinus,
					Right: *testBuildExprTree[*ExprPostfix](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
				},
			),
		},
		{
			name:    "Unary - Plus - no spaces",
			input:   `+1`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprUnary{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Op:      OpPlus,
					Right: *testBuildExprTree[*ExprPostfix](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
				},
			),
		},
		{
			name:    "Unary - Plus - spaces",
			input:   `+ 1`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprUnary{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Op:      OpPlus,
					Right: *testBuildExprTree[*ExprPostfix](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
				},
			),
		},

		{
			name:    "Postfix - no spaces",
			input:   `1[2]`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprPostfix{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprPrimary](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Right: testBuildExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},
		{
			name:    "Postfix - spaces",
			input:   `1 [ 2 ]`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprPostfix{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprPrimary](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Right: testBuildExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
			),
		},

		{
			name:    "Primary - Value",
			input:   `1`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprPrimary{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Value: testBuildExprTree[*Value](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
				},
			),
		},
		{
			name:    "Primary - Sub - no spaces",
			input:   `(1)`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprPrimary{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					SubExpression: testBuildExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
				},
			),
		},
		{
			name:    "Primary - Sub - spaces",
			input:   `( 1 )`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprPrimary{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					SubExpression: testBuildExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
				},
			),
		},

		{
			name:    "Add 3",
			input:   `1 + 2 + 3`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprAdditive{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprMultiplicative](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: "+",
					Right: &ExprAdditive{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Left: *testBuildExprTree[*ExprMultiplicative](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
							Number:  &ValueNumber{big.NewFloat(2), "2"},
						}),
						Op: "+",
						Right: testBuildExprTree[*ExprAdditive](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
							Number:  &ValueNumber{big.NewFloat(3), "3"},
						}),
					},
				},
			),
		},
		{
			name:    "Add and Multiply",
			input:   `1 + 2 * 3`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprAdditive{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *testBuildExprTree[*ExprMultiplicative](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					Op: "+",
					Right: testBuildExprTree[*ExprAdditive](t, &ExprMultiplicative{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Left: *testBuildExprTree[*ExprUnary](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
							Number:  &ValueNumber{big.NewFloat(2), "2"},
						}),
						Op: "*",
						Right: testBuildExprTree[*ExprMultiplicative](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
							Number:  &ValueNumber{big.NewFloat(3), "3"},
						}),
					}),
				},
			),
		},
		{
			name:    "Multiply and Add",
			input:   `1 * 2 + 3`,
			wantErr: false,
			want: testBuildExprTree[*Expr](t,
				&ExprAdditive{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: ExprMultiplicative{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Left: *testBuildExprTree[*ExprUnary](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Number:  &ValueNumber{big.NewFloat(1), "1"},
						}),
						Op: "*",
						Right: testBuildExprTree[*ExprMultiplicative](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
							Number:  &ValueNumber{big.NewFloat(2), "2"},
						}),
					},
					Op: "+",
					Right: testBuildExprTree[*ExprAdditive](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
						Number:  &ValueNumber{big.NewFloat(3), "3"},
					}),
				}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, true)
		})
	}
}

// /////////////////////////////////////

func TestExpr_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *Expr
		want  *Expr
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name: "Empty",
			input: &Expr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
			want: &Expr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
		},

		{
			name: "Left",
			input: &Expr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Left:    &ExprConditional{},
			},
			want: &Expr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Left:    &ExprConditional{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*Expr](t, tt.want, tt.input.Clone())
		})
	}
}

func TestIf_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprIf
		want  *ExprIf
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &ExprIf{},
			want:  &ExprIf{},
		},

		{
			name: "Condition",
			input: &ExprIf{
				Condition: ExprLogicalOr{},
			},
			want: &ExprIf{
				Condition: ExprLogicalOr{},
			},
		},
		{
			name: "Left",
			input: &ExprIf{
				Left: &Expr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprIf{
				Left: &Expr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
		},
		{
			name: "Right",
			input: &ExprIf{
				Right: &Expr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprIf{
				Right: &Expr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*ExprIf](t, tt.want, tt.input.Clone())
		})
	}
}

func TestSwitch_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprSwitch
		want  *ExprSwitch
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &ExprSwitch{},
			want:  &ExprSwitch{},
		},

		{
			name: "Selector",
			input: &ExprSwitch{
				Selector: ExprLogicalOr{},
			},
			want: &ExprSwitch{
				Selector: ExprLogicalOr{},
			},
		},
		{
			name: "Cases",
			input: &ExprSwitch{
				Cases: []*ExprCase{{}},
			},
			want: &ExprSwitch{
				Cases: []*ExprCase{{}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*ExprSwitch](t, tt.want, tt.input.Clone())
		})
	}
}

func TestCase_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprCase
		want  *ExprCase
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &ExprCase{},
			want:  &ExprCase{},
		},

		{
			name: "Conditions",
			input: &ExprCase{
				Conditions: []*ExprLogicalOr{{}},
			},
			want: &ExprCase{
				Conditions: []*ExprLogicalOr{{}},
			},
		},
		{
			name: "Default",
			input: &ExprCase{
				Default: true,
			},
			want: &ExprCase{
				Default: true,
			},
		},
		{
			name: "Expr",
			input: &ExprCase{
				Expr: &Expr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprCase{
				Expr: &Expr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*ExprCase](t, tt.want, tt.input.Clone())
		})
	}
}

func TestConditional_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprConditional
		want  *ExprConditional
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &ExprConditional{},
			want:  &ExprConditional{},
		},

		{
			name: "With Data",
			input: &ExprConditional{
				Condition:   ExprLogicalOr{},
				ConditionOp: true,
				TrueExpr:    &ExprLogicalOr{},
				FalseExpr:   &ExprLogicalOr{},
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
			testCloner[*ExprConditional](t, tt.want, tt.input.Clone())
		})
	}
}

func TestLogicalOr_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprLogicalOr
		want  *ExprLogicalOr
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &ExprLogicalOr{},
			want:  &ExprLogicalOr{},
		},

		{
			name: "Left",
			input: &ExprLogicalOr{
				Left: ExprLogicalAnd{},
			},
			want: &ExprLogicalOr{
				Left: ExprLogicalAnd{},
			},
		},
		{
			name: "Op",
			input: &ExprLogicalOr{
				Op: "||",
			},
			want: &ExprLogicalOr{
				Op: "||",
			},
		},
		{
			name: "Right",
			input: &ExprLogicalOr{
				Right: &ExprLogicalOr{},
			},
			want: &ExprLogicalOr{
				Right: &ExprLogicalOr{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*ExprLogicalOr](t, tt.want, tt.input.Clone())
		})
	}
}

func TestLogicalAnd_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprLogicalAnd
		want  *ExprLogicalAnd
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &ExprLogicalAnd{},
			want:  &ExprLogicalAnd{},
		},

		{
			name: "Left",
			input: &ExprLogicalAnd{
				Left: ExprBitwiseOr{},
			},
			want: &ExprLogicalAnd{
				Left: ExprBitwiseOr{},
			},
		},
		{
			name: "Both",
			input: &ExprLogicalAnd{
				Left:  ExprBitwiseOr{},
				Op:    OpLogicalAnd,
				Right: &ExprLogicalAnd{},
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
			testCloner[*ExprLogicalAnd](t, tt.want, tt.input.Clone())
		})
	}
}

func TestBitwiseOr_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprBitwiseOr
		want  *ExprBitwiseOr
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &ExprBitwiseOr{},
			want:  &ExprBitwiseOr{},
		},

		{
			name: "Left",
			input: &ExprBitwiseOr{
				Left: ExprBitwiseXor{},
			},
			want: &ExprBitwiseOr{
				Left: ExprBitwiseXor{},
			},
		},
		{
			name: "Op",
			input: &ExprBitwiseOr{
				Op: "|",
			},
			want: &ExprBitwiseOr{
				Op: "|",
			},
		},
		{
			name: "Right",
			input: &ExprBitwiseOr{
				Right: &ExprBitwiseOr{},
			},
			want: &ExprBitwiseOr{
				Right: &ExprBitwiseOr{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*ExprBitwiseOr](t, tt.want, tt.input.Clone())
		})
	}
}

func TestBitwiseXor_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprBitwiseXor
		want  *ExprBitwiseXor
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &ExprBitwiseXor{},
			want:  &ExprBitwiseXor{},
		},

		{
			name: "Left",
			input: &ExprBitwiseXor{
				Left: ExprBitwiseAnd{},
			},
			want: &ExprBitwiseXor{
				Left: ExprBitwiseAnd{},
			},
		},
		{
			name: "Op",
			input: &ExprBitwiseXor{
				Op: "^",
			},
			want: &ExprBitwiseXor{
				Op: "^",
			},
		},
		{
			name: "Right",
			input: &ExprBitwiseXor{
				Right: &ExprBitwiseXor{},
			},
			want: &ExprBitwiseXor{
				Right: &ExprBitwiseXor{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*ExprBitwiseXor](t, tt.want, tt.input.Clone())
		})
	}
}

func TestBitwiseAnd_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprBitwiseAnd
		want  *ExprBitwiseAnd
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &ExprBitwiseAnd{},
			want:  &ExprBitwiseAnd{},
		},

		{
			name: "Left",
			input: &ExprBitwiseAnd{
				Left: ExprEquality{},
			},
			want: &ExprBitwiseAnd{
				Left: ExprEquality{},
			},
		},
		{
			name: "Both",
			input: &ExprBitwiseAnd{
				Left:  ExprEquality{},
				Op:    OpBitwiseAnd,
				Right: &ExprBitwiseAnd{},
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
			testCloner[*ExprBitwiseAnd](t, tt.want, tt.input.Clone())
		})
	}
}

func TestEquality_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprEquality
		want  *ExprEquality
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &ExprEquality{},
			want:  &ExprEquality{},
		},

		{
			name: "Left",
			input: &ExprEquality{
				Left: ExprRelational{},
			},
			want: &ExprEquality{
				Left: ExprRelational{},
			},
		},
		{
			name: "Both",
			input: &ExprEquality{
				Left:  ExprRelational{},
				Op:    OpNotEqual,
				Right: &ExprEquality{},
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
			testCloner[*ExprEquality](t, tt.want, tt.input.Clone())
		})
	}
}

func TestRelational_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprRelational
		want  *ExprRelational
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &ExprRelational{},
			want:  &ExprRelational{},
		},

		{
			name: "Left",
			input: &ExprRelational{
				Left: ExprShift{},
			},
			want: &ExprRelational{
				Left: ExprShift{},
			},
		},
		{
			name: "Both",
			input: &ExprRelational{
				Left:  ExprShift{},
				Op:    OpLessOrEqual,
				Right: &ExprRelational{},
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
			testCloner[*ExprRelational](t, tt.want, tt.input.Clone())
		})
	}
}

func TestShift_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprShift
		want  *ExprShift
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &ExprShift{},
			want:  &ExprShift{},
		},

		{
			name: "Left",
			input: &ExprShift{
				Left: ExprAdditive{},
			},
			want: &ExprShift{
				Left: ExprAdditive{},
			},
		},
		{
			name: "Both",
			input: &ExprShift{
				Left:  ExprAdditive{},
				Op:    OpBitwiseShiftRight,
				Right: &ExprShift{},
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
			testCloner[*ExprShift](t, tt.want, tt.input.Clone())
		})
	}
}

func TestAdditive_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprAdditive
		want  *ExprAdditive
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &ExprAdditive{},
			want:  &ExprAdditive{},
		},

		{
			name: "Left",
			input: &ExprAdditive{
				Left: ExprMultiplicative{},
			},
			want: &ExprAdditive{
				Left: ExprMultiplicative{},
			},
		},
		{
			name: "Both",
			input: &ExprAdditive{
				Left:  ExprMultiplicative{},
				Op:    OpMinus,
				Right: &ExprAdditive{},
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
			testCloner[*ExprAdditive](t, tt.want, tt.input.Clone())
		})
	}
}

func TestMultiplicative_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprMultiplicative
		want  *ExprMultiplicative
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &ExprMultiplicative{},
			want:  &ExprMultiplicative{},
		},

		{
			name: "Left",
			input: &ExprMultiplicative{
				Left: ExprUnary{},
			},
			want: &ExprMultiplicative{
				Left: ExprUnary{},
			},
		},
		{
			name: "Both",
			input: &ExprMultiplicative{
				Left:  ExprUnary{},
				Op:    OpModulo,
				Right: &ExprMultiplicative{},
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
			testCloner[*ExprMultiplicative](t, tt.want, tt.input.Clone())
		})
	}
}

func TestUnary_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprUnary
		want  *ExprUnary
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &ExprUnary{},
			want:  &ExprUnary{},
		},

		{
			name: "Operator",
			input: &ExprUnary{
				Op:    OpMinus,
				Right: ExprPostfix{},
			},
			want: &ExprUnary{
				Op:    OpMinus,
				Right: ExprPostfix{},
			},
		},
		{
			name: "No Operator",
			input: &ExprUnary{
				Right: ExprPostfix{},
			},
			want: &ExprUnary{
				Right: ExprPostfix{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*ExprUnary](t, tt.want, tt.input.Clone())
		})
	}
}

func TestPostfix_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprPostfix
		want  *ExprPostfix
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &ExprPostfix{},
			want:  &ExprPostfix{},
		},

		{
			name: "Left",
			input: &ExprPostfix{
				Left: ExprPrimary{},
			},
			want: &ExprPostfix{
				Left: ExprPrimary{},
			},
		},
		{
			name: "Right",
			input: &ExprPostfix{
				Right: &Expr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprPostfix{
				Right: &Expr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
		},
		{
			name: "Left and Right",
			input: &ExprPostfix{
				Left: ExprPrimary{},
				Right: &Expr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprPostfix{
				Left: ExprPrimary{},
				Right: &Expr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*ExprPostfix](t, tt.want, tt.input.Clone())
		})
	}
}

func TestPrimary_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		Input *ExprPrimary
		want  *ExprPrimary
	}{
		{
			name:  "Nil",
			Input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			Input: &ExprPrimary{},
			want:  &ExprPrimary{},
		},

		{
			name: "Sub Expression",
			Input: &ExprPrimary{
				SubExpression: &Expr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprPrimary{
				SubExpression: &Expr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
		},
		{
			name: "Invocation",
			Input: &ExprPrimary{
				Invocation: &ExprInvocation{},
			},
			want: &ExprPrimary{
				Invocation: &ExprInvocation{},
			},
		},
		{
			name: "Value",
			Input: &ExprPrimary{
				Value: &Value{},
			},
			want: &ExprPrimary{
				Value: &Value{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*ExprPrimary](t, tt.want, tt.Input.Clone())
		})
	}
}

func TestInvocation_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		Input *ExprInvocation
		want  *ExprInvocation
	}{
		{
			name:  "Nil",
			Input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			Input: &ExprInvocation{},
			want:  &ExprInvocation{},
		},

		{
			name: "Ident",
			Input: &ExprInvocation{
				Ident: Ident{},
			},
			want: &ExprInvocation{
				Ident: Ident{},
			},
		},
		{
			name: "Monads",
			Input: &ExprInvocation{
				Monads: []*ExprInvocationParams{{}},
			},
			want: &ExprInvocation{
				Monads: []*ExprInvocationParams{{}},
			},
		},
		{
			name: "Postfix",
			Input: &ExprInvocation{
				Postfix: &ExprPostfix{},
			},
			want: &ExprInvocation{
				Postfix: &ExprPostfix{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*ExprInvocation](t, tt.want, tt.Input.Clone())
		})
	}
}

func TestInvocationParams_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		Input *ExprInvocationParams
		want  *ExprInvocationParams
	}{
		{
			name:  "Nil",
			Input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			Input: &ExprInvocationParams{},
			want:  &ExprInvocationParams{},
		},
		{
			name: "No Values",
			Input: &ExprInvocationParams{
				Values: []*Expr{},
			},
			want: &ExprInvocationParams{
				Values: []*Expr{},
			},
		},
		{
			name: "Values",
			Input: &ExprInvocationParams{
				Values: []*Expr{
					testBuildExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
				},
			},
			want: &ExprInvocationParams{
				Values: []*Expr{
					testBuildExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*ExprInvocationParams](t, tt.want, tt.Input.Clone())
		})
	}
}

// /////////////////////////////////////

func TestExpr_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *Expr
		want  []Node
	}{
		{
			name:  "Empty",
			input: &Expr{},
			want:  nil,
		},
		{
			name: "Left",
			input: &Expr{
				Left: &ExprConditional{},
			},
			want: []Node{
				&ExprConditional{},
			},
		},
		{
			name: "If",
			input: &Expr{
				If: &ExprIf{},
			},
			want: []Node{
				&ExprIf{},
			},
		},
		{
			name: "Switch",
			input: &Expr{
				Switch: &ExprSwitch{},
			},
			want: []Node{
				&ExprSwitch{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestIf_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprIf
		want  []Node
	}{
		{
			name:  "Empty",
			input: &ExprIf{},
			want: []Node{
				&ExprLogicalOr{},
			},
		},
		{
			name: "Condition",
			input: &ExprIf{
				Condition: ExprLogicalOr{},
			},
			want: []Node{
				&ExprLogicalOr{},
			},
		},
		{
			name: "Left",
			input: &ExprIf{
				Left: &Expr{},
			},
			want: []Node{
				&ExprLogicalOr{},
				&Expr{},
			},
		},
		{
			name: "Right",
			input: &ExprIf{
				Right: &Expr{},
			},
			want: []Node{
				&ExprLogicalOr{},
				&Expr{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestSwitch_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprSwitch
		want  []Node
	}{
		{
			name:  "Empty",
			input: &ExprSwitch{},
			want: []Node{
				&ExprLogicalOr{},
			},
		},
		{
			name: "Selector",
			input: &ExprSwitch{
				Selector: ExprLogicalOr{},
			},
			want: []Node{
				&ExprLogicalOr{},
			},
		},
		{
			name: "Cases",
			input: &ExprSwitch{
				Cases: []*ExprCase{
					{}, {},
				},
			},
			want: []Node{
				&ExprLogicalOr{},
				&ExprCase{},
				&ExprCase{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestCase_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprCase
		want  []Node
	}{
		{
			name:  "Empty",
			input: &ExprCase{},
			want:  nil,
		},
		{
			name: "Conditions",
			input: &ExprCase{
				Conditions: []*ExprLogicalOr{
					{}, {},
				},
			},
			want: []Node{
				&ExprLogicalOr{},
				&ExprLogicalOr{},
			},
		},
		{
			name: "Default",
			input: &ExprCase{
				Default: true,
			},
			want: nil,
		},
		{
			name: "Expr",
			input: &ExprCase{
				Expr: &Expr{},
			},
			want: []Node{
				&Expr{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestConditional_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprConditional
		want  []Node
	}{
		{
			name:  "Empty",
			input: &ExprConditional{},
			want: []Node{
				&ExprLogicalOr{},
			},
		},
		{
			name: "Condition",
			input: &ExprConditional{
				Condition: ExprLogicalOr{},
			},
			want: []Node{
				&ExprLogicalOr{},
			},
		},
		{
			name: "ConditionOp",
			input: &ExprConditional{
				ConditionOp: true,
			},
			want: []Node{
				&ExprLogicalOr{},
			},
		},
		{
			name: "TrueExpr",
			input: &ExprConditional{
				TrueExpr: &ExprLogicalOr{},
			},
			want: []Node{
				&ExprLogicalOr{},
				&ExprLogicalOr{},
			},
		},
		{
			name: "FalseExpr",
			input: &ExprConditional{
				FalseExpr: &ExprLogicalOr{},
			},
			want: []Node{
				&ExprLogicalOr{},
				&ExprLogicalOr{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestLogicalOr_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprLogicalOr
		want  []Node
	}{
		{
			name:  "Empty",
			input: &ExprLogicalOr{},
			want: []Node{
				&ExprLogicalAnd{},
			},
		},
		{
			name: "Left",
			input: &ExprLogicalOr{
				Left: ExprLogicalAnd{},
			},
			want: []Node{
				&ExprLogicalAnd{},
			},
		},
		{
			name: "Op",
			input: &ExprLogicalOr{
				Op: "foo",
			},
			want: []Node{
				&ExprLogicalAnd{},
			},
		},
		{
			name: "Right",
			input: &ExprLogicalOr{
				Right: &ExprLogicalOr{},
			},
			want: []Node{
				&ExprLogicalAnd{},
				&ExprLogicalOr{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestLogicalAnd_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprLogicalAnd
		want  []Node
	}{
		{
			name:  "Empty",
			input: &ExprLogicalAnd{},
			want: []Node{
				&ExprBitwiseOr{},
			},
		},
		{
			name: "Left",
			input: &ExprLogicalAnd{
				Left: ExprBitwiseOr{},
			},
			want: []Node{
				&ExprBitwiseOr{},
			},
		},
		{
			name: "Op",
			input: &ExprLogicalAnd{
				Op: "foo",
			},
			want: []Node{
				&ExprBitwiseOr{},
			},
		},
		{
			name: "Right",
			input: &ExprLogicalAnd{
				Right: &ExprLogicalAnd{},
			},
			want: []Node{
				&ExprBitwiseOr{},
				&ExprLogicalAnd{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestBitwiseOr_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprBitwiseOr
		want  []Node
	}{
		{
			name:  "Empty",
			input: &ExprBitwiseOr{},
			want: []Node{
				&ExprBitwiseXor{},
			},
		},
		{
			name: "Left",
			input: &ExprBitwiseOr{
				Left: ExprBitwiseXor{},
			},
			want: []Node{
				&ExprBitwiseXor{},
			},
		},
		{
			name: "Op",
			input: &ExprBitwiseOr{
				Op: "foo",
			},
			want: []Node{
				&ExprBitwiseXor{},
			},
		},
		{
			name: "Right",
			input: &ExprBitwiseOr{
				Right: &ExprBitwiseOr{},
			},
			want: []Node{
				&ExprBitwiseXor{},
				&ExprBitwiseOr{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestBitwiseXor_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprBitwiseXor
		want  []Node
	}{
		{
			name:  "Empty",
			input: &ExprBitwiseXor{},
			want: []Node{
				&ExprBitwiseAnd{},
			},
		},
		{
			name: "Left",
			input: &ExprBitwiseXor{
				Left: ExprBitwiseAnd{},
			},
			want: []Node{
				&ExprBitwiseAnd{},
			},
		},
		{
			name: "Op",
			input: &ExprBitwiseXor{
				Op: "foo",
			},
			want: []Node{
				&ExprBitwiseAnd{},
			},
		},
		{
			name: "Right",
			input: &ExprBitwiseXor{
				Right: &ExprBitwiseXor{},
			},
			want: []Node{
				&ExprBitwiseAnd{},
				&ExprBitwiseXor{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestBitwiseAnd_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprBitwiseAnd
		want  []Node
	}{
		{
			name:  "Empty",
			input: &ExprBitwiseAnd{},
			want: []Node{
				&ExprEquality{},
			},
		},
		{
			name: "Left",
			input: &ExprBitwiseAnd{
				Left: ExprEquality{},
			},
			want: []Node{
				&ExprEquality{},
			},
		},
		{
			name: "Op",
			input: &ExprBitwiseAnd{
				Op: "foo",
			},
			want: []Node{
				&ExprEquality{},
			},
		},
		{
			name: "Right",
			input: &ExprBitwiseAnd{
				Right: &ExprBitwiseAnd{},
			},
			want: []Node{
				&ExprEquality{},
				&ExprBitwiseAnd{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestEquality_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprEquality
		want  []Node
	}{
		{
			name:  "Empty",
			input: &ExprEquality{},
			want: []Node{
				&ExprRelational{},
			},
		},
		{
			name: "Left",
			input: &ExprEquality{
				Left: ExprRelational{},
			},
			want: []Node{
				&ExprRelational{},
			},
		},
		{
			name: "Op",
			input: &ExprEquality{
				Op: "foo",
			},
			want: []Node{
				&ExprRelational{},
			},
		},
		{
			name: "Right",
			input: &ExprEquality{
				Right: &ExprEquality{},
			},
			want: []Node{
				&ExprRelational{},
				&ExprEquality{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestRelational_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprRelational
		want  []Node
	}{
		{
			name:  "Empty",
			input: &ExprRelational{},
			want: []Node{
				&ExprShift{},
			},
		},
		{
			name: "Left",
			input: &ExprRelational{
				Left: ExprShift{},
			},
			want: []Node{
				&ExprShift{},
			},
		},
		{
			name: "Op",
			input: &ExprRelational{
				Op: "foo",
			},
			want: []Node{
				&ExprShift{},
			},
		},
		{
			name: "Right",
			input: &ExprRelational{
				Right: &ExprRelational{},
			},
			want: []Node{
				&ExprShift{},
				&ExprRelational{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestShift_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprShift
		want  []Node
	}{
		{
			name:  "Empty",
			input: &ExprShift{},
			want: []Node{
				&ExprAdditive{},
			},
		},
		{
			name: "Left",
			input: &ExprShift{
				Left: ExprAdditive{},
			},
			want: []Node{
				&ExprAdditive{},
			},
		},
		{
			name: "Op",
			input: &ExprShift{
				Op: "foo",
			},
			want: []Node{
				&ExprAdditive{},
			},
		},
		{
			name: "Right",
			input: &ExprShift{
				Right: &ExprShift{},
			},
			want: []Node{
				&ExprAdditive{},
				&ExprShift{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestAdditive_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprAdditive
		want  []Node
	}{
		{
			name:  "Empty",
			input: &ExprAdditive{},
			want: []Node{
				&ExprMultiplicative{},
			},
		},
		{
			name: "Left",
			input: &ExprAdditive{
				Left: ExprMultiplicative{},
			},
			want: []Node{
				&ExprMultiplicative{},
			},
		},
		{
			name: "Op",
			input: &ExprAdditive{
				Op: "foo",
			},
			want: []Node{
				&ExprMultiplicative{},
			},
		},
		{
			name: "Right",
			input: &ExprAdditive{
				Right: &ExprAdditive{},
			},
			want: []Node{
				&ExprMultiplicative{},
				&ExprAdditive{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestMultiplicative_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprMultiplicative
		want  []Node
	}{
		{
			name:  "Empty",
			input: &ExprMultiplicative{},
			want: []Node{
				&ExprUnary{},
			},
		},
		{
			name: "Left",
			input: &ExprMultiplicative{
				Left: ExprUnary{},
			},
			want: []Node{
				&ExprUnary{},
			},
		},
		{
			name: "Op",
			input: &ExprMultiplicative{
				Op: "foo",
			},
			want: []Node{
				&ExprUnary{},
			},
		},
		{
			name: "Right",
			input: &ExprMultiplicative{
				Right: &ExprMultiplicative{},
			},
			want: []Node{
				&ExprUnary{},
				&ExprMultiplicative{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestUnary_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprUnary
		want  []Node
	}{
		{
			name:  "Empty",
			input: &ExprUnary{},
			want: []Node{
				&ExprPostfix{},
			},
		},
		{
			name: "Op",
			input: &ExprUnary{
				Op: "foo",
			},
			want: []Node{
				&ExprPostfix{},
			},
		},
		{
			name: "Right",
			input: &ExprUnary{
				Right: ExprPostfix{},
			},
			want: []Node{
				&ExprPostfix{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestPostfix_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprPostfix
		want  []Node
	}{
		{
			name:  "Empty",
			input: &ExprPostfix{},
			want: []Node{
				&ExprPrimary{},
			},
		},
		{
			name: "Left",
			input: &ExprPostfix{
				Left: ExprPrimary{},
			},
			want: []Node{
				&ExprPrimary{},
			},
		},
		{
			name: "Right",
			input: &ExprPostfix{
				Right: &Expr{},
			},
			want: []Node{
				&ExprPrimary{},
				&Expr{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestPrimary_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprPrimary
		want  []Node
	}{
		{
			name:  "Empty",
			input: &ExprPrimary{},
			want:  nil,
		},
		{
			name: "SubExpression",
			input: &ExprPrimary{
				SubExpression: &Expr{},
			},
			want: []Node{
				&Expr{},
			},
		},
		{
			name: "Invocation",
			input: &ExprPrimary{
				Invocation: &ExprInvocation{},
			},
			want: []Node{
				&ExprInvocation{},
			},
		},
		{
			name: "Value",
			input: &ExprPrimary{
				Value: &Value{},
			},
			want: []Node{
				&Value{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestInvocation_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprInvocation
		want  []Node
	}{
		{
			name:  "Empty",
			input: &ExprInvocation{},
			want: []Node{
				&Ident{},
			},
		},
		{
			name: "Ident",
			input: &ExprInvocation{
				Ident: Ident{},
			},
			want: []Node{
				&Ident{},
			},
		},
		{
			name: "Monads",
			input: &ExprInvocation{
				Monads: []*ExprInvocationParams{
					{}, {},
				},
			},
			want: []Node{
				&Ident{},
				&ExprInvocationParams{},
				&ExprInvocationParams{},
			},
		},
		{
			name: "Postfix",
			input: &ExprInvocation{
				Postfix: &ExprPostfix{},
			},
			want: []Node{
				&Ident{},
				&ExprPostfix{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestInvocationParams_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprInvocationParams
		want  []Node
	}{
		{
			name:  "Empty",
			input: &ExprInvocationParams{},
			want:  nil,
		},
		{
			name: "Values",
			input: &ExprInvocationParams{
				Values: []*Expr{
					{}, {},
				},
			},
			want: []Node{
				&Expr{},
				&Expr{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

// /////////////////////////////////////

func TestExpr_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		description string
		input       *Expr
		wantPanic   bool
		want        string
	}{
		{
			name:        "nil",
			description: "",
			input:       nil,
			wantPanic:   true,
		},
		{
			name:        "Empty",
			description: "",
			input: &Expr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
			wantPanic: true,
			want:      "expression not set",
		},
		{
			name:        "Left",
			description: "",
			input: &Expr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Left: testBuildExprTree[*ExprConditional](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
			},
			want: "1",
		},
		{
			name:        "If",
			description: "",
			input: &Expr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				If: &ExprIf{
					Condition: *testBuildExprTree[*ExprLogicalOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Ident:   &Ident{Parts: []string{"foo"}},
					}),
					Left:  nil,
					Right: nil,
				},
			},
			want: `if foo { }`,
		},
		{
			name:        "Switch",
			description: "",
			input: &Expr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Switch: &ExprSwitch{
					Selector: *testBuildExprTree[*ExprLogicalOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Ident:   &Ident{Parts: []string{"foo"}},
					}),
					Cases: nil,
				},
			},
			want: `switch foo { }`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

func TestIf_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		description string
		input       *ExprIf
		wantPanic   bool
		want        string
	}{
		{
			name:      "Nil",
			input:     nil,
			wantPanic: true,
		},
		{
			name:        "Empty",
			description: "Condition can never be nil",
			input:       &ExprIf{},
			want:        "if null { }",
		},
		{
			name: "If with empty body",
			input: &ExprIf{
				Condition: *testBuildExprTree[*ExprLogicalOr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Ident:   &Ident{Parts: []string{"foo"}},
				}),
				Left:  nil,
				Right: nil,
			},
			want: `if foo { }`,
		},
		{
			name: "If",
			input: &ExprIf{
				Condition: *testBuildExprTree[*ExprLogicalOr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Ident:   &Ident{Parts: []string{"foo"}},
				}),
				Left: testBuildExprTree[*Expr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
				Right: nil,
			},
			want: `
if foo {
	1
}`[1:],
		},
		{
			name: "If/Else",
			input: &ExprIf{
				Condition: *testBuildExprTree[*ExprLogicalOr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Ident:   &Ident{Parts: []string{"foo"}},
				}),
				Left: testBuildExprTree[*Expr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
				Right: testBuildExprTree[*Expr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(2), "2"},
				}),
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
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

func TestSwitch_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		description string
		input       *ExprSwitch
		wantPanic   bool
		want        string
	}{
		{
			name:      "Nil",
			input:     nil,
			wantPanic: true,
		},
		{
			name:  "Empty",
			input: &ExprSwitch{},
			want:  `switch null { }`,
		},
		{
			name: "Only Selector",
			input: &ExprSwitch{
				Selector: *testBuildExprTree[*ExprLogicalOr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Ident:   &Ident{Parts: []string{"foo"}},
				}),
				Cases: nil,
			},
			want: `switch foo { }`,
		},
		{
			name: "One single case, no default",
			input: &ExprSwitch{
				Selector: *testBuildExprTree[*ExprLogicalOr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Ident:   &Ident{Parts: []string{"foo"}},
				}),
				Cases: []*ExprCase{
					{
						Conditions: []*ExprLogicalOr{
							testBuildExprTree[*ExprLogicalOr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
								Number:  &ValueNumber{big.NewFloat(1), "1"},
							}),
						},
						Expr: testBuildExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Number:  &ValueNumber{big.NewFloat(2), "2"},
						}),
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
			input: &ExprSwitch{
				Selector: *testBuildExprTree[*ExprLogicalOr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Ident:   &Ident{Parts: []string{"foo"}},
				}),
				Cases: []*ExprCase{
					{
						Conditions: []*ExprLogicalOr{
							testBuildExprTree[*ExprLogicalOr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
								Number:  &ValueNumber{big.NewFloat(1), "1"},
							}),
							testBuildExprTree[*ExprLogicalOr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
								Number:  &ValueNumber{big.NewFloat(2), "2"},
							}),
						},
						Expr: testBuildExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Number:  &ValueNumber{big.NewFloat(3), "3"},
						}),
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
			input: &ExprSwitch{
				Selector: *testBuildExprTree[*ExprLogicalOr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Ident:   &Ident{Parts: []string{"foo"}},
				}),
				Cases: []*ExprCase{
					{
						Conditions: []*ExprLogicalOr{
							testBuildExprTree[*ExprLogicalOr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
								Number:  &ValueNumber{big.NewFloat(1), "1"},
							}),
						},
						Expr: testBuildExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Number:  &ValueNumber{big.NewFloat(2), "2"},
						}),
					},
					{
						Conditions: []*ExprLogicalOr{
							testBuildExprTree[*ExprLogicalOr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
								Number:  &ValueNumber{big.NewFloat(3), "3"},
							}),
						},
						Expr: testBuildExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Number:  &ValueNumber{big.NewFloat(4), "4"},
						}),
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
			input: &ExprSwitch{
				Selector: *testBuildExprTree[*ExprLogicalOr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Ident:   &Ident{Parts: []string{"foo"}},
				}),
				Cases: []*ExprCase{
					{
						Conditions: nil,
						Default:    true,
						Expr: testBuildExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Number:  &ValueNumber{big.NewFloat(3), "3"},
						}),
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
			input: &ExprSwitch{
				Selector: *testBuildExprTree[*ExprLogicalOr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Ident:   &Ident{Parts: []string{"foo"}},
				}),
				Cases: []*ExprCase{
					{
						Conditions: []*ExprLogicalOr{
							testBuildExprTree[*ExprLogicalOr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
								Number:  &ValueNumber{big.NewFloat(1), "1"},
							}),
						},
						Expr: testBuildExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Number:  &ValueNumber{big.NewFloat(2), "2"},
						}),
					},
					{
						Conditions: nil,
						Default:    true,
						Expr: testBuildExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Number:  &ValueNumber{big.NewFloat(3), "3"},
						}),
					},
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
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

func TestCase_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		description string
		input       *ExprCase
		wantPanic   bool
		want        string
	}{
		{
			name:      "Nil",
			input:     nil,
			wantPanic: true,
		},
		{
			name:        "Empty",
			description: "No condition and not default",
			input:       &ExprCase{},
			wantPanic:   true,
			want:        "non-default case statement without condition",
		},
		{
			name: "One condition",
			input: &ExprCase{
				Conditions: []*ExprLogicalOr{
					testBuildExprTree[*ExprLogicalOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
				},
				Default: false,
				Expr: testBuildExprTree[*Expr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(2), "2"},
				}),
			},
			want: `
case 1: {
	2
}`[1:],
		},
		{
			name: "Several condition",
			input: &ExprCase{
				Conditions: []*ExprLogicalOr{
					testBuildExprTree[*ExprLogicalOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
					testBuildExprTree[*ExprLogicalOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number:  &ValueNumber{big.NewFloat(2), "2"},
					}),
				},
				Default: false,
				Expr: testBuildExprTree[*Expr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(3), "3"},
				}),
			},
			want: `
case 1, 2: {
	3
}`[1:],
		},
		{
			name: "Default",
			input: &ExprCase{
				Conditions: nil,
				Default:    true,
				Expr: testBuildExprTree[*Expr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(3), "3"},
				}),
			},
			want: `
default: {
	3
}`[1:],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

func TestConditional_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		description string
		input       *ExprConditional
		wantPanic   bool
		want        string
	}{
		{
			name:      "Nil",
			input:     nil,
			wantPanic: true,
		},
		{
			name: "Only Condition",
			input: &ExprConditional{
				Condition: *testBuildExprTree[*ExprLogicalOr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
			},
			want: "1",
		},
		{
			name: "All expressions but not operators",
			input: &ExprConditional{
				Condition: *testBuildExprTree[*ExprLogicalOr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
				TrueExpr: testBuildExprTree[*ExprLogicalOr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(2), "2"},
				}),
				FalseExpr: testBuildExprTree[*ExprLogicalOr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(3), "3"},
				}),
			},
			want: "1",
		},
		{
			name: "Operator with neither True nor False",
			input: &ExprConditional{
				Condition: *testBuildExprTree[*ExprLogicalOr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
				ConditionOp: true,
			},
			want: "1 ? null : null",
		},
		{
			name: "Only True expression",
			input: &ExprConditional{
				Condition: *testBuildExprTree[*ExprLogicalOr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
				ConditionOp: true,
				TrueExpr: testBuildExprTree[*ExprLogicalOr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(2), "2"},
				}),
			},
			want: "1 ? 2 : null",
		},
		{
			name: "Only False expression",
			input: &ExprConditional{
				Condition: *testBuildExprTree[*ExprLogicalOr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
				ConditionOp: true,
				FalseExpr: testBuildExprTree[*ExprLogicalOr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(3), "3"},
				}),
			},
			want: "1 ? null : 3",
		},
		{
			name:        "All parts",
			description: "Both sides of the condition must be present.",
			input: &ExprConditional{
				Condition: *testBuildExprTree[*ExprLogicalOr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
				ConditionOp: true,
				TrueExpr: testBuildExprTree[*ExprLogicalOr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(2), "2"},
				}),
				FalseExpr: testBuildExprTree[*ExprLogicalOr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(3), "3"},
				}),
			},
			want: "1 ? 2 : 3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

func TestLogicalOr_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		description string
		input       *ExprLogicalOr
		wantPanic   bool
		want        string
	}{
		{
			name:      "Nil",
			input:     nil,
			wantPanic: true,
		},
		{
			name:  "Empty",
			input: &ExprLogicalOr{},
			want:  TokenNull,
		},
		{
			name: "Left",
			input: &ExprLogicalOr{
				Left: *testBuildExprTree[*ExprLogicalAnd](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			input: &ExprLogicalOr{
				Left: ExprLogicalAnd{},
				Op:   OpLogicalOr,
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			input: &ExprLogicalOr{
				Left: *testBuildExprTree[*ExprLogicalAnd](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
				Op: OpLogicalOr,
				Right: testBuildExprTree[*ExprLogicalOr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(2), "2"},
				}),
			},
			want: "1 || 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

func TestLogicalAnd_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		description string
		input       *ExprLogicalAnd
		wantPanic   bool
		want        string
	}{
		{
			name:      "Nil",
			input:     nil,
			wantPanic: true,
		},
		{
			name:  "Empty",
			input: &ExprLogicalAnd{},
			want:  TokenNull,
		},
		{
			name: "Left",
			input: &ExprLogicalAnd{
				Left: *testBuildExprTree[*ExprBitwiseOr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			input: &ExprLogicalAnd{
				Left: ExprBitwiseOr{},
				Op:   OpLogicalAnd,
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			input: &ExprLogicalAnd{
				Left: *testBuildExprTree[*ExprBitwiseOr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
				Op: OpLogicalAnd,
				Right: testBuildExprTree[*ExprLogicalAnd](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(2), "2"},
				}),
			},
			want: "1 && 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

func TestBitwiseOr_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		description string
		input       *ExprBitwiseOr
		wantPanic   bool
		want        string
	}{
		{
			name:      "Nil",
			input:     nil,
			wantPanic: true,
		},
		{
			name:  "Empty",
			input: &ExprBitwiseOr{},
			want:  TokenNull,
		},
		{
			name: "Left",
			input: &ExprBitwiseOr{
				Left: *testBuildExprTree[*ExprBitwiseXor](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			input: &ExprBitwiseOr{
				Left: ExprBitwiseXor{},
				Op:   OpBitwiseOr,
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			input: &ExprBitwiseOr{
				Left: *testBuildExprTree[*ExprBitwiseXor](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
				Op: OpBitwiseOr,
				Right: testBuildExprTree[*ExprBitwiseOr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(2), "2"},
				}),
			},
			want: "1 | 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

func TestBitwiseXor_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		description string
		input       *ExprBitwiseXor
		wantPanic   bool
		want        string
	}{
		{
			name:      "Nil",
			input:     nil,
			wantPanic: true,
		},
		{
			name:  "Empty",
			input: &ExprBitwiseXor{},
			want:  TokenNull,
		},
		{
			name: "Left",
			input: &ExprBitwiseXor{
				Left: *testBuildExprTree[*ExprBitwiseAnd](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			input: &ExprBitwiseXor{
				Left: ExprBitwiseAnd{},
				Op:   OpBitwiseXOr,
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			input: &ExprBitwiseXor{
				Left: *testBuildExprTree[*ExprBitwiseAnd](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
				Op: OpBitwiseXOr,
				Right: testBuildExprTree[*ExprBitwiseXor](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(2), "2"},
				}),
			},
			want: "1 ^ 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

func TestBitwiseAnd_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		description string
		input       *ExprBitwiseAnd
		wantPanic   bool
		want        string
	}{
		{
			name:      "Nil",
			input:     nil,
			wantPanic: true,
		},
		{
			name:  "Empty",
			input: &ExprBitwiseAnd{},
			want:  TokenNull,
		},
		{
			name: "Left",
			input: &ExprBitwiseAnd{
				Left: *testBuildExprTree[*ExprEquality](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			input: &ExprBitwiseAnd{
				Left: ExprEquality{},
				Op:   OpBitwiseAnd,
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			input: &ExprBitwiseAnd{
				Left: *testBuildExprTree[*ExprEquality](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
				Op: OpBitwiseAnd,
				Right: testBuildExprTree[*ExprBitwiseAnd](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(2), "2"},
				}),
			},
			want: "1 & 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

func TestEquality_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		description string
		input       *ExprEquality
		wantPanic   bool
		want        string
	}{
		{
			name:      "Nil",
			input:     nil,
			wantPanic: true,
		},
		{
			name:  "Empty",
			input: &ExprEquality{},
			want:  TokenNull,
		},
		{
			name: "Left",
			input: &ExprEquality{
				Left: *testBuildExprTree[*ExprRelational](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			input: &ExprEquality{
				Left: ExprRelational{},
				Op:   OpEqual,
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			input: &ExprEquality{
				Left: *testBuildExprTree[*ExprRelational](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
				Op: OpEqual,
				Right: testBuildExprTree[*ExprEquality](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(2), "2"},
				}),
			},
			want: "1 == 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

func TestRelational_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		description string
		input       *ExprRelational
		wantPanic   bool
		want        string
	}{
		{
			name:      "Nil",
			input:     nil,
			wantPanic: true,
		},
		{
			name:      "Empty",
			input:     &ExprRelational{},
			wantPanic: false,
			want:      TokenNull,
		},
		{
			name: "Left",
			input: &ExprRelational{
				Left: *testBuildExprTree[*ExprShift](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			input: &ExprRelational{
				Left: ExprShift{},
				Op:   OpMoreOrEqual,
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			input: &ExprRelational{
				Left: *testBuildExprTree[*ExprShift](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
				Op: OpMoreOrEqual,
				Right: testBuildExprTree[*ExprRelational](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(2), "2"},
				}),
			},
			want: "1 >= 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

func TestShift_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		description string
		input       *ExprShift
		wantPanic   bool
		want        string
	}{
		{
			name:      "Nil",
			input:     nil,
			wantPanic: true,
		},
		{
			name:  "Empty",
			input: &ExprShift{},
			want:  TokenNull,
		},
		{
			name: "Left",
			input: &ExprShift{
				Left: *testBuildExprTree[*ExprAdditive](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			input: &ExprShift{
				Left: ExprAdditive{},
				Op:   OpBitwiseShiftRight,
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			input: &ExprShift{
				Left: *testBuildExprTree[*ExprAdditive](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
				Op: OpBitwiseShiftRight,
				Right: testBuildExprTree[*ExprShift](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(2), "2"},
				}),
			},
			want: "1 >> 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

func TestAdditive_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		description string
		input       *ExprAdditive
		wantPanic   bool
		want        string
	}{
		{
			name:      "Nil",
			input:     nil,
			wantPanic: true,
		},
		{
			name:  "Empty",
			input: &ExprAdditive{},
			want:  TokenNull,
		},
		{
			name: "Left",
			input: &ExprAdditive{
				Left: *testBuildExprTree[*ExprMultiplicative](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			input: &ExprAdditive{
				Left: ExprMultiplicative{},
				Op:   OpPlus,
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			input: &ExprAdditive{
				Left: *testBuildExprTree[*ExprMultiplicative](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
				Op: OpPlus,
				Right: testBuildExprTree[*ExprAdditive](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(2), "2"},
				}),
			},
			want: "1 + 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

func TestMultiplicative_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		description string
		input       *ExprMultiplicative
		wantPanic   bool
		want        string
	}{
		{
			name:      "Nil",
			input:     nil,
			wantPanic: true,
		},
		{
			name:  "Empty",
			input: &ExprMultiplicative{},
			want:  TokenNull,
		},
		{
			name: "Left",
			input: &ExprMultiplicative{
				Left: *testBuildExprTree[*ExprUnary](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			input: &ExprMultiplicative{
				Left: *testBuildExprTree[*ExprUnary](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Ident:   &Ident{Parts: []string{"foo"}},
				}),
				Op: OpMultiplication,
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			input: &ExprMultiplicative{
				Left: *testBuildExprTree[*ExprUnary](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
				Op: OpMultiplication,
				Right: testBuildExprTree[*ExprMultiplicative](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(2), "2"},
				}),
			},
			want: "1 * 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

func TestUnary_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		description string
		input       *ExprUnary
		want        string
		wantPanic   bool
	}{
		{
			name:      "Nil",
			input:     nil,
			wantPanic: true,
		},
		{
			name:  "Empty",
			input: &ExprUnary{},
			want:  TokenNull,
		},
		{
			name: "Right",
			input: &ExprUnary{
				Right: *testBuildExprTree[*ExprPostfix](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Ident:   &Ident{Parts: []string{"foo"}},
				}),
			},
			want: "foo",
		},
		{
			name: "Operator",
			input: &ExprUnary{
				Op: OpBitwiseNot,
				Right: *testBuildExprTree[*ExprPostfix](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Ident:   &Ident{Parts: []string{"foo"}},
				}),
			},
			want: "~foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

func TestPostfix_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		description string
		input       *ExprPostfix
		wantPanic   bool
		want        string
	}{
		{
			name:      "Nil",
			input:     nil,
			wantPanic: true,
		},
		{
			name:  "Empty",
			input: &ExprPostfix{},
			want:  TokenNull,
		},
		{
			name: "Left",
			input: &ExprPostfix{
				Left: *testBuildExprTree[*ExprPrimary](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Ident:   &Ident{Parts: []string{"foo"}},
				}),
			},
			want: "foo",
		},
		{
			name: "Both",
			input: &ExprPostfix{
				Left: *testBuildExprTree[*ExprPrimary](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Ident:   &Ident{Parts: []string{"foo"}},
				}),
				Right: testBuildExprTree[*Expr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
			},
			want: "foo[1]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

func TestPrimary_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		description string
		input       *ExprPrimary
		want        string
		wantPanic   bool
	}{
		{
			name:      "Nil",
			input:     nil,
			wantPanic: true,
		},
		{
			name:  "Empty",
			input: &ExprPrimary{},
			want:  TokenNull,
		},
		{
			name: "Sub Expression",
			input: &ExprPrimary{
				SubExpression: testBuildExprTree[*Expr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
			},
			want: "1",
		},
		{
			name: "Value",
			input: &ExprPrimary{
				Value: &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"}},
			},
			want: "1",
		},
		{
			name: "Invocation",
			input: &ExprPrimary{
				Value: &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"}},
				Invocation: &ExprInvocation{
					Ident: Ident{
						Parts: []string{"foo"},
					},
					Monads: []*ExprInvocationParams{{
						Values: []*Expr{
							testBuildExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
								Ident:   &Ident{Parts: []string{"bar"}},
							}),
							testBuildExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
								Ident:   &Ident{Parts: []string{"baz"}},
							}),
						},
					}},
				},
			},
			want: "foo(bar, baz)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

func TestInvocation_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		description string
		input       *ExprInvocation
		want        string
		wantPanic   bool
	}{
		{
			name:      "Nil",
			input:     nil,
			wantPanic: true,
		},
		{
			name:  "Empty",
			input: &ExprInvocation{},
			want:  "()",
		},
		{
			name: "No parameters",
			input: &ExprInvocation{
				Ident: Ident{
					Parts: []string{
						"foo",
					},
				},
			},
			want: "foo()",
		},
		{
			name: "One parameter",
			input: &ExprInvocation{
				Ident: Ident{
					Parts: []string{
						"foo",
					},
				},
				Monads: []*ExprInvocationParams{{
					Values: []*Expr{
						testBuildExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Number:  &ValueNumber{big.NewFloat(1), "1"},
						}),
					},
				}},
			},
			want: "foo(1)",
		},
		{
			name: "Several parameters",
			input: &ExprInvocation{
				Ident: Ident{
					Parts: []string{
						"foo",
					},
				},
				Monads: []*ExprInvocationParams{{
					Values: []*Expr{
						testBuildExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Number:  &ValueNumber{big.NewFloat(1), "1"},
						}),
						testBuildExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Number:  &ValueNumber{big.NewFloat(2), "2"},
						}),
					},
				}},
			},
			want: "foo(1, 2)",
		},
		{
			name: "Monadic invocation",
			input: &ExprInvocation{
				Ident: Ident{
					Parts: []string{
						"foo",
					},
				},
				Monads: []*ExprInvocationParams{
					{
						Values: []*Expr{
							testBuildExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
								Number:  &ValueNumber{big.NewFloat(1), "1"},
							}),
						},
					},
					{
						Values: []*Expr{
							testBuildExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
								Number:  &ValueNumber{big.NewFloat(2), "2"},
							}),
						},
					},
				},
			},
			want: "foo(1)(2)",
		},
		{
			name: "Dot reference on invocation",
			input: &ExprInvocation{
				Ident: Ident{
					Parts: []string{
						"foo",
					},
				},
				Postfix: testBuildExprTree[*ExprPostfix](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Ident:   &Ident{Parts: []string{"bar"}},
				}),
			},
			want: "foo().bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

func TestInvocationParams_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		description string
		input       *ExprInvocationParams
		want        string
		wantPanic   bool
	}{
		{
			name:      "Nil",
			input:     nil,
			wantPanic: true,
		},
		{
			name:  "Empty",
			input: &ExprInvocationParams{},
			want:  "",
		},
		{
			name: "One value",
			input: &ExprInvocationParams{
				Values: []*Expr{
					testBuildExprTree[*Expr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
				},
			},
			want: "foo",
		},
		{
			name: "Two values",
			input: &ExprInvocationParams{
				Values: []*Expr{
					testBuildExprTree[*Expr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
					testBuildExprTree[*Expr](t, &Value{Ident: &Ident{Parts: []string{"bar"}}}),
				},
			},
			want: "foo, bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}
