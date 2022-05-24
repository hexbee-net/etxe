package etx

import (
	"math/big"
	"testing"
)

func TestFunc_Parsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Func
	}{
		{
			name:    "Empty body, no params, no return",
			input:   `def foo() {}`,
			wantErr: false,
			want: &Func{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
			},
		},
		{
			name:    "Empty body, one ident param, no return",
			input:   `def foo(bar: bool) {}`,
			wantErr: false,
			want: &Func{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
				Parameters: []*FuncParameter{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
						Label:   "bar",
						Type: &ParameterType{
							ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 1, Column: 14}},
							Ident:   &Ident{Parts: []string{"bool"}},
						},
					},
				},
			},
		},
		{
			name:    "Empty body, two ident param, no return",
			input:   `def foo(bar: bool, baz: number) {}`,
			wantErr: false,
			want: &Func{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
				Parameters: []*FuncParameter{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
						Label:   "bar",
						Type: &ParameterType{
							ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 1, Column: 14}},
							Ident:   &Ident{Parts: []string{"bool"}},
						},
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 19, Line: 1, Column: 20}},
						Label:   "baz",
						Type: &ParameterType{
							ASTNode: ASTNode{Pos: Position{Offset: 24, Line: 1, Column: 25}},
							Ident:   &Ident{Parts: []string{"number"}},
						},
					},
				},
			},
		},
		{
			name:    "Empty body, no params, ident return",
			input:   `def foo() bool {}`,
			wantErr: false,
			want: &Func{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
				Return: &ParameterType{
					ASTNode: ASTNode{Pos: Position{Offset: 10, Line: 1, Column: 11}},
					Ident:   &Ident{Parts: []string{"bool"}},
				},
			},
		},

		{
			name:    "Empty body, one func param, no return",
			input:   `def foo(bar: (int) -> bool) {}`,
			wantErr: false,
			want: &Func{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
				Parameters: []*FuncParameter{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
						Label:   "bar",
						Type: &ParameterType{
							ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 1, Column: 14}},
							Func: &FuncSignature{
								ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 1, Column: 14}},
								Parameters: []*ParameterType{
									{
										ASTNode: ASTNode{Pos: Position{Offset: 14, Line: 1, Column: 15}},
										Ident:   &Ident{Parts: []string{"int"}},
									},
								},
								Return: ParameterType{
									ASTNode: ASTNode{Pos: Position{Offset: 22, Line: 1, Column: 23}},
									Ident:   &Ident{Parts: []string{"bool"}},
								},
							},
						},
					},
				},
			},
		},
		{
			name:    "Empty body, two func params, no return",
			input:   `def foo(bar: (int) -> bool, baz: (bool) -> int) {}`,
			wantErr: false,
			want: &Func{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
				Parameters: []*FuncParameter{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
						Label:   "bar",
						Type: &ParameterType{
							ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 1, Column: 14}},
							Func: &FuncSignature{
								ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 1, Column: 14}},
								Parameters: []*ParameterType{
									{
										ASTNode: ASTNode{Pos: Position{Offset: 14, Line: 1, Column: 15}},
										Ident:   &Ident{Parts: []string{"int"}},
									},
								},
								Return: ParameterType{
									ASTNode: ASTNode{Pos: Position{Offset: 22, Line: 1, Column: 23}},
									Ident:   &Ident{Parts: []string{"bool"}},
								},
							},
						},
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 28, Line: 1, Column: 29}},
						Label:   "baz",
						Type: &ParameterType{
							ASTNode: ASTNode{Pos: Position{Offset: 33, Line: 1, Column: 34}},
							Func: &FuncSignature{
								ASTNode: ASTNode{Pos: Position{Offset: 33, Line: 1, Column: 34}},
								Parameters: []*ParameterType{
									{
										ASTNode: ASTNode{Pos: Position{Offset: 34, Line: 1, Column: 35}},
										Ident:   &Ident{Parts: []string{"bool"}},
									},
								},
								Return: ParameterType{
									ASTNode: ASTNode{Pos: Position{Offset: 43, Line: 1, Column: 44}},
									Ident:   &Ident{Parts: []string{"int"}},
								},
							},
						},
					},
				},
			},
		},
		{
			name:    "Empty body, no params, func return",
			input:   `def foo() (int) -> bool {}`,
			wantErr: false,
			want: &Func{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
				Return: &ParameterType{
					ASTNode: ASTNode{Pos: Position{Offset: 10, Line: 1, Column: 11}},
					Func: &FuncSignature{
						ASTNode: ASTNode{Pos: Position{Offset: 10, Line: 1, Column: 11}},
						Parameters: []*ParameterType{
							{
								ASTNode: ASTNode{Pos: Position{Offset: 11, Line: 1, Column: 12}},
								Ident:   &Ident{Parts: []string{"int"}},
							},
						},
						Return: ParameterType{
							ASTNode: ASTNode{Pos: Position{Offset: 19, Line: 1, Column: 20}},
							Ident:   &Ident{Parts: []string{"bool"}},
						},
					},
				},
			},
		},

		{
			name: "One Expr statement, no params, no return",
			input: `
def foo() {
	a
}`[1:],
			wantErr: false,
			want: &Func{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
				Body: []*FuncStatement{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 2, Column: 2}},
						Expr: testBuildExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 2, Column: 2}},
							Ident:   &Ident{Parts: []string{"a"}},
						}),
					},
				},
			},
		},
		{
			name: "One val Decl statement, no params, no return",
			input: `
def foo() {
	val a: number = 1
}`[1:],
			wantErr: false,
			want: &Func{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
				Body: []*FuncStatement{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 2, Column: 2}},
						Decl: &FuncDecl{
							ASTNode:  ASTNode{Pos: Position{Offset: 13, Line: 2, Column: 2}},
							DeclType: "val",
							Label:    "a",
							Type:     "number",
							Value: testBuildExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 29, Line: 2, Column: 18}},
								Number:  &ValueNumber{big.NewFloat(1), "1"},
							}),
						},
					},
				},
			},
		},
		{
			name: "Val decl and return expr statements, no params, no return",
			input: `
def foo() {
	val a: number = 1
	a
}`[1:],
			wantErr: false,
			want: &Func{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
				Body: []*FuncStatement{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 2, Column: 2}},
						Decl: &FuncDecl{
							ASTNode:  ASTNode{Pos: Position{Offset: 13, Line: 2, Column: 2}},
							DeclType: "val",
							Label:    "a",
							Type:     "number",
							Value: testBuildExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 29, Line: 2, Column: 18}},
								Number:  &ValueNumber{big.NewFloat(1), "1"},
							}),
						},
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 32, Line: 3, Column: 2}},
						Expr: testBuildExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 32, Line: 3, Column: 2}},
							Ident:   &Ident{Parts: []string{"a"}},
						}),
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
