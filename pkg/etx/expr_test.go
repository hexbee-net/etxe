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
			name:    "Add 3",
			input:   `1 + 2 + 3`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprAdditive{
					Left: *BuildTestExprTree[*ExprMultiplicative](t, &Value{
						Number: &ValueNumber{
							Value:  big.NewFloat(1),
							Source: "1",
						},
					}),
					Op: "+",
					Right: &ExprAdditive{
						Left: *BuildTestExprTree[*ExprMultiplicative](t, &Value{
							Number: &ValueNumber{
								Value:  big.NewFloat(2),
								Source: "2",
							},
						}),
						Op: "+",
						Right: BuildTestExprTree[*ExprAdditive](t, &Value{
							Number: &ValueNumber{
								Value:  big.NewFloat(3),
								Source: "3",
							},
						}),
					},
				},
			),
		},
		{
			name:    "Add and Multiply",
			input:   `1 + 2 * 3`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprAdditive{
					Left: *BuildTestExprTree[*ExprMultiplicative](t, &Value{
						Number: &ValueNumber{
							Value:  big.NewFloat(1),
							Source: "1",
						},
					}),
					Op: "+",
					Right: BuildTestExprTree[*ExprAdditive](t, &ExprMultiplicative{
						Left: *BuildTestExprTree[*ExprUnary](t, &Value{
							Number: &ValueNumber{
								Value:  big.NewFloat(2),
								Source: "2",
							},
						}),
						Op: "*",
						Right: BuildTestExprTree[*ExprMultiplicative](t, &Value{
							Number: &ValueNumber{
								Value:  big.NewFloat(3),
								Source: "3",
							},
						}),
					}),
				},
			),
		},
		{
			name:    "Multiply and Add",
			input:   `1 * 2 + 3`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprAdditive{
					Left: ExprMultiplicative{
						Left: *BuildTestExprTree[*ExprUnary](t, &Value{
							Number: &ValueNumber{
								Value:  big.NewFloat(1),
								Source: "1",
							},
						}),
						Op: "*",
						Right: BuildTestExprTree[*ExprMultiplicative](t, &Value{
							Number: &ValueNumber{
								Value:  big.NewFloat(2),
								Source: "2",
							},
						}),
					},
					Op: "+",
					Right: BuildTestExprTree[*ExprAdditive](t, &Value{
						Number: &ValueNumber{
							Value:  big.NewFloat(3),
							Source: "3",
						},
					}),
				}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, false)
		})
	}
}

func TestExpr_Parsing_If(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Expr
	}{
		{
			name:    "If with empty body",
			input:   `if foo { }`,
			wantErr: false,
			want: &Expr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				If: &ExprIf{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Condition: *BuildTestExprTree[*ExprLogicalOr](t, &Ident{
						ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
						Parts:   []string{"foo"},
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
					Condition: *BuildTestExprTree[*ExprLogicalOr](t, &Ident{
						ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
						Parts:   []string{"foo"},
					}),
					Left: BuildTestExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 1, Column: 10}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 1, Column: 10}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
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
					Condition: *BuildTestExprTree[*ExprLogicalOr](t, &Ident{
						ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
						Parts:   []string{"foo"},
					}),
					Left: BuildTestExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 1, Column: 10}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 1, Column: 10}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Right: BuildTestExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 20, Line: 1, Column: 21}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 20, Line: 1, Column: 21}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, true)
		})
	}
}

func TestExpr_Parsing_Switch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Expr
	}{
		{
			name: "Switch - Empty",
			input: `
switch foo {}`[1:],
			wantErr: false,
			want: &Expr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Switch: &ExprSwitch{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Selector: *BuildTestExprTree[*ExprLogicalOr](t, &Ident{
						ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
						Parts:   []string{"foo"},
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
					Selector: *BuildTestExprTree[*ExprLogicalOr](t, &Ident{
						ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
						Parts:   []string{"foo"},
					}),
					Cases: []*ExprCase{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 14, Line: 2, Column: 2}},
							Conditions: []*ExprLogicalOr{
								BuildTestExprTree[*ExprLogicalOr](t, &Value{
									ASTNode: ASTNode{Pos: Position{Offset: 19, Line: 2, Column: 7}},
									Number: &ValueNumber{
										ASTNode: ASTNode{Pos: Position{Offset: 19, Line: 2, Column: 7}},
										Value:   big.NewFloat(1),
										Source:  "1",
									},
								}),
							},
							Expr: BuildTestExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 24, Line: 2, Column: 12}},
								Number: &ValueNumber{
									ASTNode: ASTNode{Pos: Position{Offset: 24, Line: 2, Column: 12}},
									Value:   big.NewFloat(2),
									Source:  "2",
								},
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
					Selector: *BuildTestExprTree[*ExprLogicalOr](t, &Ident{
						ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
						Parts:   []string{"foo"},
					}),
					Cases: []*ExprCase{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 14, Line: 2, Column: 2}},
							Conditions: []*ExprLogicalOr{
								BuildTestExprTree[*ExprLogicalOr](t, &Value{
									ASTNode: ASTNode{Pos: Position{Offset: 19, Line: 2, Column: 7}},
									Number: &ValueNumber{
										ASTNode: ASTNode{Pos: Position{Offset: 19, Line: 2, Column: 7}},
										Value:   big.NewFloat(1),
										Source:  "1",
									},
								}),
								BuildTestExprTree[*ExprLogicalOr](t, &Value{
									ASTNode: ASTNode{Pos: Position{Offset: 22, Line: 2, Column: 10}},
									Number: &ValueNumber{
										ASTNode: ASTNode{Pos: Position{Offset: 22, Line: 2, Column: 10}},
										Value:   big.NewFloat(2),
										Source:  "2",
									},
								}),
							},
							Expr: BuildTestExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 27, Line: 2, Column: 15}},
								Number: &ValueNumber{
									ASTNode: ASTNode{Pos: Position{Offset: 27, Line: 2, Column: 15}},
									Value:   big.NewFloat(3),
									Source:  "3",
								},
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
					Selector: *BuildTestExprTree[*ExprLogicalOr](t, &Ident{
						ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
						Parts:   []string{"foo"},
					}),
					Cases: []*ExprCase{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 14, Line: 2, Column: 2}},
							Conditions: []*ExprLogicalOr{
								BuildTestExprTree[*ExprLogicalOr](t, &Value{
									ASTNode: ASTNode{Pos: Position{Offset: 19, Line: 2, Column: 7}},
									Number: &ValueNumber{
										ASTNode: ASTNode{Pos: Position{Offset: 19, Line: 2, Column: 7}},
										Value:   big.NewFloat(1),
										Source:  "1",
									},
								}),
							},
							Expr: BuildTestExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 24, Line: 2, Column: 12}},
								Number: &ValueNumber{
									ASTNode: ASTNode{Pos: Position{Offset: 24, Line: 2, Column: 12}},
									Value:   big.NewFloat(2),
									Source:  "2",
								},
							}),
						},
						{
							ASTNode: ASTNode{Pos: Position{Offset: 29, Line: 3, Column: 2}},
							Conditions: []*ExprLogicalOr{
								BuildTestExprTree[*ExprLogicalOr](t, &Value{
									ASTNode: ASTNode{Pos: Position{Offset: 34, Line: 3, Column: 7}},
									Number: &ValueNumber{
										ASTNode: ASTNode{Pos: Position{Offset: 34, Line: 3, Column: 7}},
										Value:   big.NewFloat(3),
										Source:  "3",
									},
								}),
							},
							Expr: BuildTestExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 39, Line: 3, Column: 12}},
								Number: &ValueNumber{
									ASTNode: ASTNode{Pos: Position{Offset: 39, Line: 3, Column: 12}},
									Value:   big.NewFloat(4),
									Source:  "4",
								},
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
					Selector: *BuildTestExprTree[*ExprLogicalOr](t, &Ident{
						ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
						Parts:   []string{"foo"},
					}),
					Cases: []*ExprCase{
						{
							ASTNode:    ASTNode{Pos: Position{Offset: 14, Line: 2, Column: 2}},
							Conditions: nil,
							Default:    true,
							Expr: BuildTestExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 25, Line: 2, Column: 13}},
								Number: &ValueNumber{
									ASTNode: ASTNode{Pos: Position{Offset: 25, Line: 2, Column: 13}},
									Value:   big.NewFloat(3),
									Source:  "3",
								},
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
					Selector: *BuildTestExprTree[*ExprLogicalOr](t, &Ident{
						ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
						Parts:   []string{"foo"},
					}),
					Cases: []*ExprCase{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 14, Line: 2, Column: 2}},
							Conditions: []*ExprLogicalOr{
								BuildTestExprTree[*ExprLogicalOr](t, &Value{
									ASTNode: ASTNode{Pos: Position{Offset: 19, Line: 2, Column: 7}},
									Number: &ValueNumber{
										ASTNode: ASTNode{Pos: Position{Offset: 19, Line: 2, Column: 7}},
										Value:   big.NewFloat(1),
										Source:  "1",
									},
								}),
							},
							Expr: BuildTestExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 24, Line: 2, Column: 12}},
								Number: &ValueNumber{
									ASTNode: ASTNode{Pos: Position{Offset: 24, Line: 2, Column: 12}},
									Value:   big.NewFloat(2),
									Source:  "2",
								},
							}),
						},
						{
							ASTNode:    ASTNode{Pos: Position{Offset: 29, Line: 3, Column: 2}},
							Conditions: nil,
							Default:    true,
							Expr: BuildTestExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 40, Line: 3, Column: 13}},
								Number: &ValueNumber{
									ASTNode: ASTNode{Pos: Position{Offset: 40, Line: 3, Column: 13}},
									Value:   big.NewFloat(3),
									Source:  "3",
								},
							}),
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, true)
		})
	}
}

func TestExpr_Parsing_Ternary(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Expr
	}{
		{
			name:    "Ternary",
			input:   "1 ? 2 : 3",
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprConditional{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Condition: *BuildTestExprTree[*ExprLogicalOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					ConditionOp: true,
					TrueExpr: BuildTestExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
					FalseExpr: BuildTestExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
							Value:   big.NewFloat(3),
							Source:  "3",
						},
					}),
				},
			),
		},
		{
			name:    "Nested Ternary - Left",
			input:   "1 ? 2 ? 3 : 4 : 5",
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprConditional{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Condition: *BuildTestExprTree[*ExprLogicalOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					ConditionOp: true,
					TrueExpr: BuildTestExprTree[*Expr](t,
						&ExprConditional{
							ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
							Condition: *BuildTestExprTree[*ExprLogicalOr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
								Number: &ValueNumber{
									ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
									Value:   big.NewFloat(2),
									Source:  "2",
								},
							}),
							ConditionOp: true,
							TrueExpr: BuildTestExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
								Number: &ValueNumber{
									ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
									Value:   big.NewFloat(3),
									Source:  "3",
								},
							}),
							FalseExpr: BuildTestExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 12, Line: 1, Column: 13}},
								Number: &ValueNumber{
									ASTNode: ASTNode{Pos: Position{Offset: 12, Line: 1, Column: 13}},
									Value:   big.NewFloat(4),
									Source:  "4",
								},
							}),
						},
					),
					FalseExpr: BuildTestExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 16, Line: 1, Column: 17}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 16, Line: 1, Column: 17}},
							Value:   big.NewFloat(5),
							Source:  "5",
						},
					}),
				},
			),
		},
		{
			name:    "Nested Ternary - Right",
			input:   "1 ? 2 : 3 ? 4 : 5",
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprConditional{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Condition: *BuildTestExprTree[*ExprLogicalOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					ConditionOp: true,
					TrueExpr: BuildTestExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
					FalseExpr: BuildTestExprTree[*Expr](t,
						&ExprConditional{
							ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
							Condition: *BuildTestExprTree[*ExprLogicalOr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
								Number: &ValueNumber{
									ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
									Value:   big.NewFloat(3),
									Source:  "3",
								},
							}),
							ConditionOp: true,
							TrueExpr: BuildTestExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 12, Line: 1, Column: 13}},
								Number: &ValueNumber{
									ASTNode: ASTNode{Pos: Position{Offset: 12, Line: 1, Column: 13}},
									Value:   big.NewFloat(4),
									Source:  "4",
								},
							}),
							FalseExpr: BuildTestExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 16, Line: 1, Column: 17}},
								Number: &ValueNumber{
									ASTNode: ASTNode{Pos: Position{Offset: 16, Line: 1, Column: 17}},
									Value:   big.NewFloat(5),
									Source:  "5",
								},
							}),
						},
					),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, true)
		})
	}
}

func TestExpr_Parsing_LogicalOr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Expr
	}{
		{
			name:    "Logical OR - no spaces",
			input:   `1||2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprLogicalOr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprLogicalAnd](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpLogicalOr,
					Right: BuildTestExprTree[*ExprLogicalOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
		{
			name:    "Logical OR - spaces",
			input:   `1 || 2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprLogicalOr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprLogicalAnd](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpLogicalOr,
					Right: BuildTestExprTree[*ExprLogicalOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, true)
		})
	}
}

func TestExpr_Parsing_LogicalAnd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Expr
	}{
		{
			name:    "Logical AND - no spaces",
			input:   `1&&2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprLogicalAnd{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprBitwiseOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpLogicalAnd,
					Right: BuildTestExprTree[*ExprLogicalAnd](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
		{
			name:    "Logical AND - spaces",
			input:   `1 && 2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprLogicalAnd{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprBitwiseOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpLogicalAnd,
					Right: BuildTestExprTree[*ExprLogicalAnd](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, true)
		})
	}
}

func TestExpr_Parsing_BitwiseOr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Expr
	}{
		{
			name:    "Bitwise OR - no spaces",
			input:   `1|2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprBitwiseOr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprBitwiseXor](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpBitwiseOr,
					Right: BuildTestExprTree[*ExprBitwiseOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
		{
			name:    "Bitwise OR - spaces",
			input:   `1 | 2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprBitwiseOr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprBitwiseXor](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpBitwiseOr,
					Right: BuildTestExprTree[*ExprBitwiseOr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, true)
		})
	}
}

func TestExpr_Parsing_BitwiseXOr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Expr
	}{
		{
			name:    "Bitwise XOR - no spaces",
			input:   `1^2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprBitwiseXor{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprBitwiseAnd](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpBitwiseXOr,
					Right: BuildTestExprTree[*ExprBitwiseXor](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
		{
			name:    "Bitwise XOR - spaces",
			input:   `1 ^ 2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprBitwiseXor{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprBitwiseAnd](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpBitwiseXOr,
					Right: BuildTestExprTree[*ExprBitwiseXor](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, true)
		})
	}
}

func TestExpr_Parsing_BitwiseAnd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Expr
	}{
		{
			name:    "Bitwise AND - no spaces",
			input:   `1&2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprBitwiseAnd{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprEquality](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpBitwiseAnd,
					Right: BuildTestExprTree[*ExprBitwiseAnd](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
		{
			name:    "Bitwise AND - spaces",
			input:   `1 & 2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprBitwiseAnd{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprEquality](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpBitwiseAnd,
					Right: BuildTestExprTree[*ExprBitwiseAnd](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, true)
		})
	}
}

func TestExpr_Parsing_Equality(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Expr
	}{
		{
			name:    "Equality - Equal - no spaces",
			input:   `1==2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprEquality{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprRelational](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpEqual,
					Right: BuildTestExprTree[*ExprEquality](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
		{
			name:    "Equality - Equal - spaces",
			input:   `1 == 2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprEquality{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprRelational](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpEqual,
					Right: BuildTestExprTree[*ExprEquality](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
		{
			name:    "Equality - Not Equal - no spaces",
			input:   `1!=2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprEquality{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprRelational](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpNotEqual,
					Right: BuildTestExprTree[*ExprEquality](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
		{
			name:    "Equality - Not Equal - spaces",
			input:   `1 != 2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprEquality{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprRelational](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpNotEqual,
					Right: BuildTestExprTree[*ExprEquality](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, true)
		})
	}
}

func TestExpr_Parsing_Relational(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Expr
	}{
		{
			name:    "Relational - More - no spaces",
			input:   `1>2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprRelational{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprShift](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpMore,
					Right: BuildTestExprTree[*ExprRelational](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
		{
			name:    "Relational - More - spaces",
			input:   `1 > 2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprRelational{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprShift](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpMore,
					Right: BuildTestExprTree[*ExprRelational](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
		{
			name:    "Relational - Less - no spaces",
			input:   `1<2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprRelational{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprShift](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpLess,
					Right: BuildTestExprTree[*ExprRelational](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
		{
			name:    "Relational - Less - spaces",
			input:   `1 < 2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprRelational{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprShift](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpLess,
					Right: BuildTestExprTree[*ExprRelational](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
		{
			name:    "Relational - More Or Equal - no spaces",
			input:   `1>=2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprRelational{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprShift](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpMoreOrEqual,
					Right: BuildTestExprTree[*ExprRelational](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
		{
			name:    "Relational - More Or Equal - spaces",
			input:   `1 >= 2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprRelational{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprShift](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpMoreOrEqual,
					Right: BuildTestExprTree[*ExprRelational](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
		{
			name:    "Relational - Less Or Equal - no spaces",
			input:   `1<=2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprRelational{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprShift](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpLessOrEqual,
					Right: BuildTestExprTree[*ExprRelational](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
		{
			name:    "Relational - Less Or Equal - spaces",
			input:   `1 <= 2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprRelational{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprShift](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpLessOrEqual,
					Right: BuildTestExprTree[*ExprRelational](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, true)
		})
	}
}

func TestExpr_Parsing_Shift(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Expr
	}{
		{
			name:    "Shift - Left - no spaces",
			input:   `1<<2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprShift{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprAdditive](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpBitwiseShiftLeft,
					Right: BuildTestExprTree[*ExprShift](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
		{
			name:    "Shift - Left - spaces",
			input:   `1 << 2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprShift{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprAdditive](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpBitwiseShiftLeft,
					Right: BuildTestExprTree[*ExprShift](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
		{
			name:    "Shift - Right - no spaces",
			input:   `1>>2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprShift{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprAdditive](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpBitwiseShiftRight,
					Right: BuildTestExprTree[*ExprShift](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
		{
			name:    "Shift - Right - spaces",
			input:   `1 >> 2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprShift{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprAdditive](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpBitwiseShiftRight,
					Right: BuildTestExprTree[*ExprShift](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, true)
		})
	}
}

func TestExpr_Parsing_Additive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Expr
	}{
		{
			name:    "Additive - Plus - no spaces",
			input:   `1+2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprAdditive{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprMultiplicative](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpPlus,
					Right: BuildTestExprTree[*ExprAdditive](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
		{
			name:    "Additive - Plus - spaces",
			input:   `1 + 2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprAdditive{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprMultiplicative](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpPlus,
					Right: BuildTestExprTree[*ExprAdditive](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
		{
			name:    "Additive - Minus - no spaces",
			input:   `1-2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprAdditive{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprMultiplicative](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpMinus,
					Right: BuildTestExprTree[*ExprAdditive](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
		{
			name:    "Additive - Minus - spaces",
			input:   `1 - 2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprAdditive{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprMultiplicative](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpMinus,
					Right: BuildTestExprTree[*ExprAdditive](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, true)
		})
	}
}

func TestExpr_Parsing_Multiplicative(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Expr
	}{
		{
			name:    "Multiplicative - Division - no spaces",
			input:   `1/2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprMultiplicative{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprUnary](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpDivision,
					Right: BuildTestExprTree[*ExprMultiplicative](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
		{
			name:    "Multiplicative - Division - spaces",
			input:   `1 / 2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprMultiplicative{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprUnary](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpDivision,
					Right: BuildTestExprTree[*ExprMultiplicative](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
		{
			name:    "Multiplicative - Multiplication - no spaces",
			input:   `1*2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprMultiplicative{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprUnary](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpMultiplication,
					Right: BuildTestExprTree[*ExprMultiplicative](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
		{
			name:    "Multiplicative - Multiplication - spaces",
			input:   `1 * 2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprMultiplicative{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprUnary](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpMultiplication,
					Right: BuildTestExprTree[*ExprMultiplicative](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
		{
			name:    "Multiplicative - Modulo - no spaces",
			input:   `1%2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprMultiplicative{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprUnary](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpModulo,
					Right: BuildTestExprTree[*ExprMultiplicative](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
		{
			name:    "Multiplicative - Modulo - spaces",
			input:   `1 % 2`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprMultiplicative{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Left: *BuildTestExprTree[*ExprUnary](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Op: OpModulo,
					Right: BuildTestExprTree[*ExprMultiplicative](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, true)
		})
	}
}

func TestExpr_Parsing_Unary(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Expr
	}{
		{
			name:    "Unary - Bitwise NOT - no spaces",
			input:   `~1`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprUnary{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Op:      OpBitwiseNot,
					Right: *BuildTestExprTree[*ExprPostfix](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
				},
			),
		},
		{
			name:    "Unary - Bitwise NOT - spaces",
			input:   `~ 1`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprUnary{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Op:      OpBitwiseNot,
					Right: *BuildTestExprTree[*ExprPostfix](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
				},
			),
		},
		{
			name:    "Unary - Logical NOT - no spaces",
			input:   `!1`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprUnary{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Op:      OpLogicalNot,
					Right: *BuildTestExprTree[*ExprPostfix](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
				},
			),
		},
		{
			name:    "Unary - Logical NOT - spaces",
			input:   `! 1`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprUnary{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Op:      OpLogicalNot,
					Right: *BuildTestExprTree[*ExprPostfix](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
				},
			),
		},
		{
			name:    "Unary - Minus - no spaces",
			input:   `-1`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprUnary{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Op:      OpMinus,
					Right: *BuildTestExprTree[*ExprPostfix](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
				},
			),
		},
		{
			name:    "Unary - Minus - spaces",
			input:   `- 1`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprUnary{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Op:      OpMinus,
					Right: *BuildTestExprTree[*ExprPostfix](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
				},
			),
		},
		{
			name:    "Unary - Plus - no spaces",
			input:   `+1`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprUnary{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Op:      OpPlus,
					Right: *BuildTestExprTree[*ExprPostfix](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
				},
			),
		},
		{
			name:    "Unary - Plus - spaces",
			input:   `+ 1`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprUnary{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Op:      OpPlus,
					Right: *BuildTestExprTree[*ExprPostfix](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, true)
		})
	}
}

func TestExpr_Parsing_Postfix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Expr
	}{
		{
			name:    "Index - No spaces",
			input:   `attr[2]`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprPostfix{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Value: *BuildTestExprTree[*ExprPrimary](t, &Ident{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Parts:   []string{"attr"},
					}),
					Index: BuildTestExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},
		{
			name:    "Index - Spaces",
			input:   `attr [ 2 ]`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprPostfix{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Value: *BuildTestExprTree[*ExprPrimary](t, &Ident{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Parts:   []string{"attr"},
					}),
					Index: BuildTestExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
							Value:   big.NewFloat(2),
							Source:  "2",
						},
					}),
				},
			),
		},

		{
			name:    "Single post",
			input:   `attr[1].field`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprPostfix{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Value: *BuildTestExprTree[*ExprPrimary](t, &Ident{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Parts:   []string{"attr"},
					}),
					Index: BuildTestExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Post: BuildTestExprTree[*ExprPostfix](t, &Ident{
						ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
						Parts:   []string{"field"},
					}),
				},
			),
		},
		{
			name:    "Post with post",
			input:   `attr[1].field[2]`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprPostfix{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Value: *BuildTestExprTree[*ExprPrimary](t, &Ident{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Parts:   []string{"attr"},
					}),
					Index: BuildTestExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Post: &ExprPostfix{
						ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
						Value: *BuildTestExprTree[*ExprPrimary](t, &Ident{
							ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
							Parts:   []string{"field"},
						}),
						Index: BuildTestExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 14, Line: 1, Column: 15}},
							Number: &ValueNumber{
								ASTNode: ASTNode{Pos: Position{Offset: 14, Line: 1, Column: 15}},
								Value:   big.NewFloat(2),
								Source:  "2",
							},
						}),
					},
				},
			),
		},
		{
			name:    "Post with post with field",
			input:   `attr[1].field[2].sub`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprPostfix{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Value: *BuildTestExprTree[*ExprPrimary](t, &Ident{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Parts:   []string{"attr"},
					}),
					Index: BuildTestExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
					Post: &ExprPostfix{
						ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
						Value: *BuildTestExprTree[*ExprPrimary](t, &Ident{
							ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
							Parts:   []string{"field"},
						}),
						Index: BuildTestExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 14, Line: 1, Column: 15}},
							Number: &ValueNumber{
								ASTNode: ASTNode{Pos: Position{Offset: 14, Line: 1, Column: 15}},
								Value:   big.NewFloat(2),
								Source:  "2",
							},
						}),
						Post: BuildTestExprTree[*ExprPostfix](t, &Ident{
							ASTNode: ASTNode{Pos: Position{Offset: 17, Line: 1, Column: 18}},
							Parts:   []string{"sub"},
						}),
					},
				},
			),
		},
		{
			name:    "Post with sub",
			input:   `attr[0].field.sub`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprPostfix{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Value: *BuildTestExprTree[*ExprPrimary](t, &Ident{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Parts:   []string{"attr"},
					}),
					Index: BuildTestExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
							Value:   big.NewFloat(0),
							Source:  "0",
						},
					}),
					Post: BuildTestExprTree[*ExprPostfix](t, &Ident{
						ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
						Parts:   []string{"field", "sub"},
					}),
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, true)
		})
	}
}

func TestExpr_Parsing_Primary(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Expr
	}{
		{
			name:    "Value",
			input:   `1`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprPrimary{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Value: BuildTestExprTree[*Value](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
				},
			),
		},

		{
			name:    "SubExpression - no spaces",
			input:   `(1)`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprPrimary{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					SubExpression: BuildTestExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
				},
			),
		},
		{
			name:    "SubExpression - spaces",
			input:   `( 1 )`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t,
				&ExprPrimary{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					SubExpression: BuildTestExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
				},
			),
		},

		{
			name:    "Invocation - No parameters",
			input:   `foo()`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t, &ExprPrimary{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Ident: &Ident{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Parts:   []string{"foo"},
				},
				Monads: []*ExprInvocationParams{{
					ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
				}},
			}),
		},
		{
			name:    "Invocation - Parameters",
			input:   `foo(bar, baz)`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t, &ExprPrimary{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Ident: &Ident{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Parts:   []string{"foo"},
				},
				Monads: []*ExprInvocationParams{{
					ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
					Values: []*Expr{
						BuildTestExprTree[*Expr](t, &Ident{
							ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
							Parts:   []string{"bar"},
						}),
						BuildTestExprTree[*Expr](t, &Ident{
							ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 1, Column: 10}},
							Parts:   []string{"baz"},
						}),
					},
				}},
			}),
		},
		{
			name:    "Invocation - Dot Invocation",
			input:   `foo.bar(baz, qux)`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t, &ExprPrimary{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Ident: &Ident{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Parts: []string{
						"foo",
						"bar",
					},
				},
				Monads: []*ExprInvocationParams{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
						Values: []*Expr{
							BuildTestExprTree[*Expr](t, &Ident{
								ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
								Parts:   []string{"baz"},
							}),
							BuildTestExprTree[*Expr](t, &Ident{
								ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 1, Column: 14}},
								Parts:   []string{"qux"},
							}),
						},
					},
				},
			}),
		},
		{
			name:    "Invocation - Monadic invocation",
			input:   `foo(bar)(baz)`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t, &ExprPrimary{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Ident: &Ident{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Parts:   []string{"foo"},
				},
				Monads: []*ExprInvocationParams{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Values: []*Expr{
							BuildTestExprTree[*Expr](t, &Ident{
								ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
								Parts:   []string{"bar"},
							}),
						},
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 1, Column: 10}},
						Values: []*Expr{
							BuildTestExprTree[*Expr](t, &Ident{
								ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 1, Column: 10}},
								Parts:   []string{"baz"},
							}),
						},
					},
				},
			}),
		},
		{
			name:    "Invocation - Dot reference on invocation",
			input:   `foo().bar`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t, &ExprPrimary{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Ident: &Ident{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Parts: []string{
						"foo",
					},
				},
				Monads: []*ExprInvocationParams{
					{ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}}},
				},
				Post: BuildTestExprTree[*ExprPostfix](t, &Ident{
					ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
					Parts:   []string{"bar"},
				}),
			}),
		},
		{
			name:    "Invocation - Dot reference invocation on invocation",
			input:   `foo().bar()`,
			wantErr: false,
			want: BuildTestExprTree[*Expr](t, &ExprPrimary{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Ident: &Ident{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Parts: []string{
						"foo",
					},
				},
				Monads: []*ExprInvocationParams{
					{ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}}},
				},
				Post: BuildTestExprTree[*ExprPostfix](t, &ExprPrimary{
					ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
					Ident: &Ident{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
						Parts: []string{
							"bar",
						},
					},
					Monads: []*ExprInvocationParams{
						{ASTNode: ASTNode{Pos: Position{Offset: 10, Line: 1, Column: 11}}},
					},
				})},
			),
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
			name:  "Empty",
			input: &Expr{},
			want:  &Expr{},
		},
		{
			name: "ASTNode",
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
				Left: &ExprConditional{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &Expr{
				Left: &ExprConditional{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
		},
		{
			name: "If",
			input: &Expr{
				If: &ExprIf{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &Expr{
				If: &ExprIf{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
		},
		{
			name: "Switch",
			input: &Expr{
				Switch: &ExprSwitch{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &Expr{
				Switch: &ExprSwitch{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
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
			name: "ASTNode",
			input: &ExprIf{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
			want: &ExprIf{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
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
			name: "ASTNode",
			input: &ExprSwitch{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
			want: &ExprSwitch{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
		},
		{
			name: "Selector",
			input: &ExprSwitch{
				Selector: ExprLogicalOr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprSwitch{
				Selector: ExprLogicalOr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
		},
		{
			name: "Cases",
			input: &ExprSwitch{
				Cases: []*ExprCase{{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				}},
			},
			want: &ExprSwitch{
				Cases: []*ExprCase{{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				}},
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
			name: "ASTNode",
			input: &ExprCase{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
			want: &ExprCase{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
		},
		{
			name: "Conditions",
			input: &ExprCase{
				Conditions: []*ExprLogicalOr{{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				}},
			},
			want: &ExprCase{
				Conditions: []*ExprLogicalOr{{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				}},
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
			name: "ASTNode",
			input: &ExprConditional{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
			want: &ExprConditional{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
		},
		{
			name: "Condition",
			input: &ExprConditional{
				Condition: ExprLogicalOr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprConditional{
				Condition: ExprLogicalOr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
		},
		{
			name: "ConditionOp",
			input: &ExprConditional{
				ConditionOp: true,
			},
			want: &ExprConditional{
				ConditionOp: true,
			},
		},
		{
			name: "TrueExpr",
			input: &ExprConditional{
				TrueExpr: &Expr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprConditional{
				TrueExpr: &Expr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
		},
		{
			name: "FalseExpr",
			input: &ExprConditional{
				FalseExpr: &Expr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprConditional{
				FalseExpr: &Expr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
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
			name: "ASTNode",
			input: &ExprLogicalOr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
			want: &ExprLogicalOr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
		},
		{
			name: "Left",
			input: &ExprLogicalOr{
				Left: ExprLogicalAnd{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprLogicalOr{
				Left: ExprLogicalAnd{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
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
				Right: &ExprLogicalOr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprLogicalOr{
				Right: &ExprLogicalOr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
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
			name: "ASTNode",
			input: &ExprLogicalAnd{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
			want: &ExprLogicalAnd{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
		},

		{
			name: "Left",
			input: &ExprLogicalAnd{
				Left: ExprBitwiseOr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprLogicalAnd{
				Left: ExprBitwiseOr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
		},
		{
			name: "Op",
			input: &ExprLogicalAnd{
				Op: OpLogicalAnd,
			},
			want: &ExprLogicalAnd{
				Op: OpLogicalAnd,
			},
		},
		{
			name: "Right",
			input: &ExprLogicalAnd{
				Right: &ExprLogicalAnd{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprLogicalAnd{
				Right: &ExprLogicalAnd{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
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
			name: "ASTNode",
			input: &ExprBitwiseOr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
			want: &ExprBitwiseOr{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
		},
		{
			name: "Left",
			input: &ExprBitwiseOr{
				Left: ExprBitwiseXor{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprBitwiseOr{
				Left: ExprBitwiseXor{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
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
				Right: &ExprBitwiseOr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprBitwiseOr{
				Right: &ExprBitwiseOr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
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
			name: "ASTNode",
			input: &ExprBitwiseXor{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
			want: &ExprBitwiseXor{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
		},
		{
			name: "Left",
			input: &ExprBitwiseXor{
				Left: ExprBitwiseAnd{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprBitwiseXor{
				Left: ExprBitwiseAnd{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
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
				Right: &ExprBitwiseXor{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprBitwiseXor{
				Right: &ExprBitwiseXor{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
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
			name: "ASTNode",
			input: &ExprBitwiseAnd{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
			want: &ExprBitwiseAnd{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
		},
		{
			name: "Left",
			input: &ExprBitwiseAnd{
				Left: ExprEquality{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprBitwiseAnd{
				Left: ExprEquality{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
		},
		{
			name: "Op",
			input: &ExprBitwiseAnd{
				Op: OpBitwiseAnd,
			},
			want: &ExprBitwiseAnd{
				Op: OpBitwiseAnd,
			},
		},
		{
			name: "Right",
			input: &ExprBitwiseAnd{
				Right: &ExprBitwiseAnd{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprBitwiseAnd{
				Right: &ExprBitwiseAnd{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
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
			name: "ASTNode",
			input: &ExprEquality{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
			want: &ExprEquality{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
		},
		{
			name: "Left",
			input: &ExprEquality{
				Left: ExprRelational{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprEquality{
				Left: ExprRelational{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
		},
		{
			name: "Op",
			input: &ExprEquality{
				Op: OpNotEqual,
			},
			want: &ExprEquality{
				Op: OpNotEqual,
			},
		},
		{
			name: "Right",
			input: &ExprEquality{
				Right: &ExprEquality{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprEquality{
				Right: &ExprEquality{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
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
			name: "ASTNode",
			input: &ExprRelational{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
			want: &ExprRelational{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
		},
		{
			name: "Left",
			input: &ExprRelational{
				Left: ExprShift{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprRelational{
				Left: ExprShift{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
		},
		{
			name: "Op",
			input: &ExprRelational{
				Op: OpLessOrEqual,
			},
			want: &ExprRelational{
				Op: OpLessOrEqual,
			},
		},
		{
			name: "Right",
			input: &ExprRelational{
				Right: &ExprRelational{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprRelational{
				Right: &ExprRelational{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
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
			name: "ASTNode",
			input: &ExprShift{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
			want: &ExprShift{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
		},
		{
			name: "Left",
			input: &ExprShift{
				Left: ExprAdditive{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprShift{
				Left: ExprAdditive{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
		},
		{
			name: "Op",
			input: &ExprShift{
				Op: OpBitwiseShiftRight,
			},
			want: &ExprShift{
				Op: OpBitwiseShiftRight,
			},
		},
		{
			name: "Right",
			input: &ExprShift{
				Right: &ExprShift{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprShift{
				Right: &ExprShift{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
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
			name: "ASTNode",
			input: &ExprAdditive{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
			want: &ExprAdditive{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
		},
		{
			name: "Left",
			input: &ExprAdditive{
				Left: ExprMultiplicative{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprAdditive{
				Left: ExprMultiplicative{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
		},
		{
			name: "Op",
			input: &ExprAdditive{
				Op: OpMinus,
			},
			want: &ExprAdditive{
				Op: OpMinus,
			},
		},
		{
			name: "Right",
			input: &ExprAdditive{
				Right: &ExprAdditive{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprAdditive{
				Right: &ExprAdditive{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
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
			name: "ASTNode",
			input: &ExprMultiplicative{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
			want: &ExprMultiplicative{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
		},
		{
			name: "Left",
			input: &ExprMultiplicative{
				Left: ExprUnary{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprMultiplicative{
				Left: ExprUnary{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
		},
		{
			name: "Op",
			input: &ExprMultiplicative{
				Op: OpModulo,
			},
			want: &ExprMultiplicative{
				Op: OpModulo,
			},
		},
		{
			name: "Right",
			input: &ExprMultiplicative{
				Right: &ExprMultiplicative{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprMultiplicative{
				Right: &ExprMultiplicative{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
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
			name: "ASTNode",
			input: &ExprUnary{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
			want: &ExprUnary{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
		},
		{
			name: "Op",
			input: &ExprUnary{
				Op: OpMinus,
			},
			want: &ExprUnary{
				Op: OpMinus,
			},
		},
		{
			name: "Right",
			input: &ExprUnary{
				Right: ExprPostfix{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprUnary{
				Right: ExprPostfix{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
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
			name: "ASTNode",
			input: &ExprPostfix{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
			want: &ExprPostfix{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
		},
		{
			name: "Value",
			input: &ExprPostfix{
				Value: ExprPrimary{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprPostfix{
				Value: ExprPrimary{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
		},
		{
			name: "Index",
			input: &ExprPostfix{
				Index: &Expr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprPostfix{
				Index: &Expr{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
		},
		{
			name: "Post",
			input: &ExprPostfix{
				Post: &ExprPostfix{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprPostfix{
				Post: &ExprPostfix{
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
		input *ExprPrimary
		want  *ExprPrimary
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &ExprPrimary{},
			want:  &ExprPrimary{},
		},
		{
			name: "ASTNode",
			input: &ExprPrimary{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
			want: &ExprPrimary{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
		},
		{
			name: "Sub Expression",
			input: &ExprPrimary{
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
			name: "Ident",
			input: &ExprPrimary{
				Ident: &Ident{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprPrimary{
				Ident: &Ident{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
		},
		{
			name: "Monads",
			input: &ExprPrimary{
				Monads: []*ExprInvocationParams{
					{ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}}},
				},
			},
			want: &ExprPrimary{
				Monads: []*ExprInvocationParams{
					{ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}}},
				},
			},
		},
		{
			name: "Post",
			input: &ExprPrimary{
				Post: &ExprPostfix{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprPrimary{
				Post: &ExprPostfix{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
		},
		{
			name: "Value",
			input: &ExprPrimary{
				Value: &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: &ExprPrimary{
				Value: &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*ExprPrimary](t, tt.want, tt.input.Clone())
		})
	}
}

func TestInvocationParams_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ExprInvocationParams
		want  *ExprInvocationParams
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &ExprInvocationParams{},
			want:  &ExprInvocationParams{},
		},
		{
			name: "ASTNode",
			input: &ExprInvocationParams{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
			want: &ExprInvocationParams{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
		},
		{
			name: "No Values",
			input: &ExprInvocationParams{
				Values: []*Expr{},
			},
			want: &ExprInvocationParams{
				Values: []*Expr{},
			},
		},
		{
			name: "Values",
			input: &ExprInvocationParams{
				Values: []*Expr{
					BuildTestExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
				},
			},
			want: &ExprInvocationParams{
				Values: []*Expr{
					BuildTestExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Number: &ValueNumber{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Value:   big.NewFloat(1),
							Source:  "1",
						},
					}),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*ExprInvocationParams](t, tt.want, tt.input.Clone())
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
				TrueExpr: &Expr{},
			},
			want: []Node{
				&ExprLogicalOr{},
				&Expr{},
			},
		},
		{
			name: "FalseExpr",
			input: &ExprConditional{
				FalseExpr: &Expr{},
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
				Value: ExprPrimary{},
			},
			want: []Node{
				&ExprPrimary{},
			},
		},
		{
			name: "Right",
			input: &ExprPostfix{
				Index: &Expr{},
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
			name: "Ident",
			input: &ExprPrimary{
				Ident: &Ident{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: []Node{
				&Ident{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
		},
		{
			name: "Monads",
			input: &ExprPrimary{
				Monads: []*ExprInvocationParams{
					{ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}}},
				},
			},
			want: []Node{
				&ExprInvocationParams{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
		},
		{
			name: "Post",
			input: &ExprPrimary{
				Post: &ExprPostfix{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
			want: []Node{
				&ExprPostfix{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
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
			input:       &Expr{},
			wantPanic:   true,
			want:        "expression not set",
		},
		{
			name:        "Left",
			description: "",
			input: &Expr{
				Left: BuildTestExprTree[*ExprConditional](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
				}),
			},
			want: "1",
		},
		{
			name:        "If",
			description: "",
			input: &Expr{
				If: &ExprIf{
					Condition: *BuildTestExprTree[*ExprLogicalOr](t, &Ident{Parts: []string{"foo"}}),
					Left:      nil,
					Right:     nil,
				},
			},
			want: `if foo { }`,
		},
		{
			name:        "Switch",
			description: "",
			input: &Expr{
				Switch: &ExprSwitch{
					Selector: *BuildTestExprTree[*ExprLogicalOr](t, &Ident{Parts: []string{"foo"}}),
					Cases:    nil,
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
			wantPanic:   true,
		},
		{
			name: "If with empty body",
			input: &ExprIf{
				Condition: *BuildTestExprTree[*ExprLogicalOr](t, &Ident{Parts: []string{"foo"}}),
				Left:      nil,
				Right:     nil,
			},
			want: `if foo { }`,
		},
		{
			name: "If",
			input: &ExprIf{
				Condition: *BuildTestExprTree[*ExprLogicalOr](t, &Ident{Parts: []string{"foo"}}),
				Left: BuildTestExprTree[*Expr](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
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
				Condition: *BuildTestExprTree[*ExprLogicalOr](t, &Ident{Parts: []string{"foo"}}),
				Left: BuildTestExprTree[*Expr](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
				}),
				Right: BuildTestExprTree[*Expr](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(2),
						Source: "2",
					},
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
			name:      "Empty",
			input:     &ExprSwitch{},
			wantPanic: true,
		},
		{
			name: "Only Selector",
			input: &ExprSwitch{
				Selector: *BuildTestExprTree[*ExprLogicalOr](t, &Ident{Parts: []string{"foo"}}),
				Cases:    nil,
			},
			want: `switch foo { }`,
		},
		{
			name: "One single case, no default",
			input: &ExprSwitch{
				Selector: *BuildTestExprTree[*ExprLogicalOr](t, &Ident{Parts: []string{"foo"}}),
				Cases: []*ExprCase{
					{
						Conditions: []*ExprLogicalOr{
							BuildTestExprTree[*ExprLogicalOr](t, &Value{
								Number: &ValueNumber{
									Value:  big.NewFloat(1),
									Source: "1",
								},
							}),
						},
						Expr: BuildTestExprTree[*Expr](t, &Value{
							Number: &ValueNumber{
								Value:  big.NewFloat(2),
								Source: "2",
							},
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
				Selector: *BuildTestExprTree[*ExprLogicalOr](t, &Ident{Parts: []string{"foo"}}),
				Cases: []*ExprCase{
					{
						Conditions: []*ExprLogicalOr{
							BuildTestExprTree[*ExprLogicalOr](t, &Value{
								Number: &ValueNumber{
									Value:  big.NewFloat(1),
									Source: "1",
								},
							}),
							BuildTestExprTree[*ExprLogicalOr](t, &Value{
								Number: &ValueNumber{
									Value:  big.NewFloat(2),
									Source: "2",
								},
							}),
						},
						Expr: BuildTestExprTree[*Expr](t, &Value{
							Number: &ValueNumber{
								Value:  big.NewFloat(3),
								Source: "3",
							},
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
				Selector: *BuildTestExprTree[*ExprLogicalOr](t, &Ident{Parts: []string{"foo"}}),
				Cases: []*ExprCase{
					{
						Conditions: []*ExprLogicalOr{
							BuildTestExprTree[*ExprLogicalOr](t, &Value{
								Number: &ValueNumber{
									Value:  big.NewFloat(1),
									Source: "1",
								},
							}),
						},
						Expr: BuildTestExprTree[*Expr](t, &Value{
							Number: &ValueNumber{
								Value:  big.NewFloat(2),
								Source: "2",
							},
						}),
					},
					{
						Conditions: []*ExprLogicalOr{
							BuildTestExprTree[*ExprLogicalOr](t, &Value{
								Number: &ValueNumber{
									Value:  big.NewFloat(3),
									Source: "3",
								},
							}),
						},
						Expr: BuildTestExprTree[*Expr](t, &Value{
							Number: &ValueNumber{
								Value:  big.NewFloat(4),
								Source: "4",
							},
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
				Selector: *BuildTestExprTree[*ExprLogicalOr](t, &Ident{Parts: []string{"foo"}}),
				Cases: []*ExprCase{
					{
						Conditions: nil,
						Default:    true,
						Expr: BuildTestExprTree[*Expr](t, &Value{
							Number: &ValueNumber{
								Value:  big.NewFloat(3),
								Source: "3",
							},
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
				Selector: *BuildTestExprTree[*ExprLogicalOr](t, &Ident{Parts: []string{"foo"}}),
				Cases: []*ExprCase{
					{
						Conditions: []*ExprLogicalOr{
							BuildTestExprTree[*ExprLogicalOr](t, &Value{
								Number: &ValueNumber{
									Value:  big.NewFloat(1),
									Source: "1",
								},
							}),
						},
						Expr: BuildTestExprTree[*Expr](t, &Value{
							Number: &ValueNumber{
								Value:  big.NewFloat(2),
								Source: "2",
							},
						}),
					},
					{
						Conditions: nil,
						Default:    true,
						Expr: BuildTestExprTree[*Expr](t, &Value{
							Number: &ValueNumber{
								Value:  big.NewFloat(3),
								Source: "3",
							},
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
					BuildTestExprTree[*ExprLogicalOr](t, &Value{
						Number: &ValueNumber{
							Value:  big.NewFloat(1),
							Source: "1",
						},
					}),
				},
				Default: false,
				Expr: BuildTestExprTree[*Expr](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(2),
						Source: "2",
					},
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
					BuildTestExprTree[*ExprLogicalOr](t, &Value{
						Number: &ValueNumber{
							Value:  big.NewFloat(1),
							Source: "1",
						},
					}),
					BuildTestExprTree[*ExprLogicalOr](t, &Value{
						Number: &ValueNumber{
							Value:  big.NewFloat(2),
							Source: "2",
						},
					}),
				},
				Default: false,
				Expr: BuildTestExprTree[*Expr](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(3),
						Source: "3",
					},
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
				Expr: BuildTestExprTree[*Expr](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(3),
						Source: "3",
					},
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
				Condition: *BuildTestExprTree[*ExprLogicalOr](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
				}),
			},
			want: "1",
		},
		{
			name: "All expressions but not operators",
			input: &ExprConditional{
				Condition: *BuildTestExprTree[*ExprLogicalOr](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
				}),
				TrueExpr: BuildTestExprTree[*Expr](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(2),
						Source: "2",
					},
				}),
				FalseExpr: BuildTestExprTree[*Expr](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(3),
						Source: "3",
					},
				}),
			},
			want: "1",
		},
		{
			name: "Operator with neither True nor False",
			input: &ExprConditional{
				Condition: *BuildTestExprTree[*ExprLogicalOr](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
				}),
				ConditionOp: true,
			},
			want: "1 ? null : null",
		},
		{
			name: "Only True expression",
			input: &ExprConditional{
				Condition: *BuildTestExprTree[*ExprLogicalOr](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
				}),
				ConditionOp: true,
				TrueExpr: BuildTestExprTree[*Expr](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(2),
						Source: "2",
					},
				}),
			},
			want: "1 ? 2 : null",
		},
		{
			name: "Only False expression",
			input: &ExprConditional{
				Condition: *BuildTestExprTree[*ExprLogicalOr](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
				}),
				ConditionOp: true,
				FalseExpr: BuildTestExprTree[*Expr](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(3),
						Source: "3",
					},
				}),
			},
			want: "1 ? null : 3",
		},
		{
			name:        "All parts",
			description: "Both sides of the condition must be present.",
			input: &ExprConditional{
				Condition: *BuildTestExprTree[*ExprLogicalOr](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
				}),
				ConditionOp: true,
				TrueExpr: BuildTestExprTree[*Expr](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(2),
						Source: "2",
					},
				}),
				FalseExpr: BuildTestExprTree[*Expr](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(3),
						Source: "3",
					},
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
			name:      "Empty",
			input:     &ExprLogicalOr{},
			wantPanic: true,
		},
		{
			name: "Left",
			input: &ExprLogicalOr{
				Left: *BuildTestExprTree[*ExprLogicalAnd](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
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
				Left: *BuildTestExprTree[*ExprLogicalAnd](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
				}),
				Op: OpLogicalOr,
				Right: BuildTestExprTree[*ExprLogicalOr](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(2),
						Source: "2",
					},
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
			name:      "Empty",
			input:     &ExprLogicalAnd{},
			wantPanic: true,
		},
		{
			name: "Left",
			input: &ExprLogicalAnd{
				Left: *BuildTestExprTree[*ExprBitwiseOr](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
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
				Left: *BuildTestExprTree[*ExprBitwiseOr](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
				}),
				Op: OpLogicalAnd,
				Right: BuildTestExprTree[*ExprLogicalAnd](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(2),
						Source: "2",
					},
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
			name:      "Empty",
			input:     &ExprBitwiseOr{},
			wantPanic: true,
		},
		{
			name: "Left",
			input: &ExprBitwiseOr{
				Left: *BuildTestExprTree[*ExprBitwiseXor](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
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
				Left: *BuildTestExprTree[*ExprBitwiseXor](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
				}),
				Op: OpBitwiseOr,
				Right: BuildTestExprTree[*ExprBitwiseOr](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(2),
						Source: "2",
					},
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
			name:      "Empty",
			input:     &ExprBitwiseXor{},
			wantPanic: true,
		},
		{
			name: "Left",
			input: &ExprBitwiseXor{
				Left: *BuildTestExprTree[*ExprBitwiseAnd](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
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
				Left: *BuildTestExprTree[*ExprBitwiseAnd](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
				}),
				Op: OpBitwiseXOr,
				Right: BuildTestExprTree[*ExprBitwiseXor](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(2),
						Source: "2",
					},
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
			name:      "Empty",
			input:     &ExprBitwiseAnd{},
			wantPanic: true,
		},
		{
			name: "Left",
			input: &ExprBitwiseAnd{
				Left: *BuildTestExprTree[*ExprEquality](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
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
				Left: *BuildTestExprTree[*ExprEquality](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
				}),
				Op: OpBitwiseAnd,
				Right: BuildTestExprTree[*ExprBitwiseAnd](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(2),
						Source: "2",
					},
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
			name:      "Empty",
			input:     &ExprEquality{},
			wantPanic: true,
		},
		{
			name: "Left",
			input: &ExprEquality{
				Left: *BuildTestExprTree[*ExprRelational](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
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
				Left: *BuildTestExprTree[*ExprRelational](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
				}),
				Op: OpEqual,
				Right: BuildTestExprTree[*ExprEquality](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(2),
						Source: "2",
					},
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
			wantPanic: true,
		},
		{
			name: "Left",
			input: &ExprRelational{
				Left: *BuildTestExprTree[*ExprShift](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
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
				Left: *BuildTestExprTree[*ExprShift](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
				}),
				Op: OpMoreOrEqual,
				Right: BuildTestExprTree[*ExprRelational](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(2),
						Source: "2",
					},
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
			name:      "Empty",
			input:     &ExprShift{},
			wantPanic: true,
		},
		{
			name: "Left",
			input: &ExprShift{
				Left: *BuildTestExprTree[*ExprAdditive](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
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
				Left: *BuildTestExprTree[*ExprAdditive](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
				}),
				Op: OpBitwiseShiftRight,
				Right: BuildTestExprTree[*ExprShift](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(2),
						Source: "2",
					},
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
			name:      "Empty",
			input:     &ExprAdditive{},
			wantPanic: true,
		},
		{
			name: "Left",
			input: &ExprAdditive{
				Left: *BuildTestExprTree[*ExprMultiplicative](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
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
				Left: *BuildTestExprTree[*ExprMultiplicative](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
				}),
				Op: OpPlus,
				Right: BuildTestExprTree[*ExprAdditive](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(2),
						Source: "2",
					},
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
			name:      "Empty",
			input:     &ExprMultiplicative{},
			wantPanic: true,
		},
		{
			name: "Left",
			input: &ExprMultiplicative{
				Left: *BuildTestExprTree[*ExprUnary](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
				}),
			},
			want: "1",
		},
		{
			name:        "Operator without right side",
			description: "Right side can't be nil when an operator is present",
			input: &ExprMultiplicative{
				Left: *BuildTestExprTree[*ExprUnary](t, &Ident{Parts: []string{"foo"}}),
				Op:   OpMultiplication,
			},
			wantPanic: true,
			want:      "operator with <nil> right side",
		},
		{
			name: "Both sides",
			input: &ExprMultiplicative{
				Left: *BuildTestExprTree[*ExprUnary](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
				}),
				Op: OpMultiplication,
				Right: BuildTestExprTree[*ExprMultiplicative](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(2),
						Source: "2",
					},
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
			name:      "Empty",
			input:     &ExprUnary{},
			wantPanic: true,
		},
		{
			name: "Right",
			input: &ExprUnary{
				Right: *BuildTestExprTree[*ExprPostfix](t, &Ident{Parts: []string{"foo"}}),
			},
			want: "foo",
		},
		{
			name: "Operator",
			input: &ExprUnary{
				Op:    OpBitwiseNot,
				Right: *BuildTestExprTree[*ExprPostfix](t, &Ident{Parts: []string{"foo"}}),
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
			name:      "Empty",
			input:     &ExprPostfix{},
			wantPanic: true,
		},
		{
			name: "Value",
			input: &ExprPostfix{
				Value: *BuildTestExprTree[*ExprPrimary](t, &Ident{Parts: []string{"foo"}}),
			},
			want: "foo",
		},
		{
			name: "Value and Index",
			input: &ExprPostfix{
				Value: *BuildTestExprTree[*ExprPrimary](t, &Ident{Parts: []string{"foo"}}),
				Index: BuildTestExprTree[*Expr](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
				}),
			},
			want: "foo[1]",
		},
		{
			name: "Value, Index and Post",
			input: &ExprPostfix{
				Value: *BuildTestExprTree[*ExprPrimary](t, &Ident{Parts: []string{"foo"}}),
				Index: BuildTestExprTree[*Expr](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
				}),
				Post: BuildTestExprTree[*ExprPostfix](t, &Ident{Parts: []string{"bar"}}),
			},
			want: "foo[1].bar",
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
			name:      "Empty",
			input:     &ExprPrimary{},
			wantPanic: true,
		},
		{
			name: "Sub Expression",
			input: &ExprPrimary{
				SubExpression: BuildTestExprTree[*Expr](t, &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1",
					},
				}),
			},
			want: "1",
		},
		{
			name: "Value",
			input: &ExprPrimary{
				Value: &Value{
					Number: &ValueNumber{
						Value:  big.NewFloat(1),
						Source: "1"},
				},
			},
			want: "1",
		},
		{
			name: "Ident",
			input: &ExprPrimary{
				Ident: &Ident{
					Parts: []string{"foo"},
				},
			},
			want: "foo",
		},

		{
			name: "Invocation - No parameters",
			input: &ExprPrimary{
				Ident: &Ident{
					Parts: []string{"foo"},
				},
				Monads: []*ExprInvocationParams{{}},
			},
			want: "foo()",
		},
		{
			name: "Invocation - One parameter",
			input: &ExprPrimary{
				Ident: &Ident{
					Parts: []string{"foo"},
				},
				Monads: []*ExprInvocationParams{{
					Values: []*Expr{
						BuildTestExprTree[*Expr](t, &Value{
							Number: &ValueNumber{
								Value:  big.NewFloat(1),
								Source: "1",
							},
						}),
					},
				}},
			},
			want: "foo(1)",
		},
		{
			name: "Invocation - Several parameters",
			input: &ExprPrimary{
				Ident: &Ident{
					Parts: []string{"foo"},
				},
				Monads: []*ExprInvocationParams{{
					Values: []*Expr{
						BuildTestExprTree[*Expr](t, &Value{
							Number: &ValueNumber{
								Value:  big.NewFloat(1),
								Source: "1",
							},
						}),
						BuildTestExprTree[*Expr](t, &Value{
							Number: &ValueNumber{
								Value:  big.NewFloat(2),
								Source: "2",
							},
						}),
					},
				}},
			},
			want: "foo(1, 2)",
		},
		{
			name: "Invocation - Monadic",
			input: &ExprPrimary{
				Ident: &Ident{
					Parts: []string{"foo"},
				},
				Monads: []*ExprInvocationParams{
					{
						Values: []*Expr{
							BuildTestExprTree[*Expr](t, &Value{
								Number: &ValueNumber{
									Value:  big.NewFloat(1),
									Source: "1",
								},
							}),
						},
					},
					{
						Values: []*Expr{
							BuildTestExprTree[*Expr](t, &Value{
								Number: &ValueNumber{
									Value:  big.NewFloat(2),
									Source: "2",
								},
							}),
						},
					},
				},
			},
			want: "foo(1)(2)",
		},
		{
			name: "Invocation - Postfix",
			input: &ExprPrimary{
				Ident: &Ident{
					Parts: []string{"foo"},
				},
				Monads: []*ExprInvocationParams{{}},
				Post:   BuildTestExprTree[*ExprPostfix](t, &Ident{Parts: []string{"bar"}}),
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
					BuildTestExprTree[*Expr](t, &Ident{Parts: []string{"foo"}}),
				},
			},
			want: "foo",
		},
		{
			name: "Two values",
			input: &ExprInvocationParams{
				Values: []*Expr{
					BuildTestExprTree[*Expr](t, &Ident{Parts: []string{"foo"}}),
					BuildTestExprTree[*Expr](t, &Ident{Parts: []string{"bar"}}),
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
