package etx

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLambda_Parsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Lambda
	}{
		{
			name:    "No parameters, simple expression",
			input:   `() => 1`,
			wantErr: false,
			want: &Lambda{
				ASTNode:    ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Parameters: nil,
				Expr: *testBuildExprTree[*Expr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
			},
		},
		{
			name:    "One parameter, no type, simple expression",
			input:   `(x) => 1`,
			wantErr: false,
			want: &Lambda{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Parameters: []*LambdaParameter{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Label:   "x",
					},
				},
				Expr: *testBuildExprTree[*Expr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
			},
		},
		{
			name:    "Two parameters, no type, simple expression",
			input:   `(x, y) => 1`,
			wantErr: false,
			want: &Lambda{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Parameters: []*LambdaParameter{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Label:   "x",
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Label:   "y",
					},
				},
				Expr: *testBuildExprTree[*Expr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 10, Line: 1, Column: 11}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
			},
		},
		{
			name:    "One parameter, type, simple expression",
			input:   `(x: number) => 1`,
			wantErr: false,
			want: &Lambda{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Parameters: []*LambdaParameter{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Label:   "x",
						Type: &ParameterType{
							ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
							Ident:   &Ident{Parts: []string{"number"}},
						},
					},
				},
				Expr: *testBuildExprTree[*Expr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 15, Line: 1, Column: 16}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
			},
		},
		{
			name:    "Two parameters, type on both, simple expression",
			input:   `(x: number, y: string) => 1`,
			wantErr: false,
			want: &Lambda{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Parameters: []*LambdaParameter{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Label:   "x",
						Type: &ParameterType{
							ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
							Ident:   &Ident{Parts: []string{"number"}},
						},
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 12, Line: 1, Column: 13}},
						Label:   "y",
						Type: &ParameterType{
							ASTNode: ASTNode{Pos: Position{Offset: 15, Line: 1, Column: 16}},
							Ident:   &Ident{Parts: []string{"string"}},
						},
					},
				},
				Expr: *testBuildExprTree[*Expr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 26, Line: 1, Column: 27}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
			},
		},
		{
			name:    "Two parameters, type on first, simple expression",
			input:   `(x: number, y) => 1`,
			wantErr: false,
			want: &Lambda{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Parameters: []*LambdaParameter{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Label:   "x",
						Type: &ParameterType{
							ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
							Ident:   &Ident{Parts: []string{"number"}},
						},
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 12, Line: 1, Column: 13}},
						Label:   "y",
					},
				},
				Expr: *testBuildExprTree[*Expr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 18, Line: 1, Column: 19}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
			},
		},
		{
			name:    "Two parameters, type on second, simple expression",
			input:   `(x, y: string) => 1`,
			wantErr: false,
			want: &Lambda{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Parameters: []*LambdaParameter{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Label:   "x",
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
						Label:   "y",
						Type: &ParameterType{
							ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
							Ident:   &Ident{Parts: []string{"string"}},
						},
					},
				},
				Expr: *testBuildExprTree[*Expr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 18, Line: 1, Column: 19}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
			},
		},
		{
			name:    "Complete declaration",
			input:   `(x: number) => x + 1`,
			wantErr: false,
			want: &Lambda{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Parameters: []*LambdaParameter{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Label:   "x",
						Type: &ParameterType{
							ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 1, Column: 5}},
							Ident:   &Ident{Parts: []string{"number"}},
						},
					},
				},
				Expr: *testBuildExprTree[*Expr](t,
					&ExprAdditive{
						ASTNode: ASTNode{Pos: Position{Offset: 15, Line: 1, Column: 16}},
						Left: *testBuildExprTree[*ExprMultiplicative](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 15, Line: 1, Column: 16}},
							Ident:   &Ident{Parts: []string{"x"}},
						}),
						Op: OpPlus,
						Right: testBuildExprTree[*ExprAdditive](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 19, Line: 1, Column: 20}},
							Number:  &ValueNumber{big.NewFloat(1), "1"},
						}),
					},
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, true)
		})
	}
}

func TestLambda_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		Input *Lambda
		want  *Lambda
	}{
		{
			name:  "Nil",
			Input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			Input: &Lambda{},
			want:  &Lambda{},
		},
		{
			name: "Parameters",
			Input: &Lambda{
				Parameters: []*LambdaParameter{{Label: "foo"}},
			},
			want: &Lambda{
				Parameters: []*LambdaParameter{{Label: "foo"}},
			},
		},
		{
			name: "Expr",
			Input: &Lambda{
				Expr: *testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
			},
			want: &Lambda{
				Expr: *testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*Lambda](t, tt.want, tt.Input.Clone())
		})
	}
}

func TestLambda_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *Lambda
		want  []Node
	}{
		{
			name:  "Empty",
			input: &Lambda{},
			want: []Node{
				&Expr{},
			},
		},
		{
			name: "Parameters",
			input: &Lambda{
				Parameters: []*LambdaParameter{{Label: "foo"}},
			},
			want: []Node{
				&LambdaParameter{Label: "foo"},
				&Expr{},
			},
		},
		{
			name: "Expr",
			input: &Lambda{
				Expr: *testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
			},
			want: []Node{
				testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestLambda_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     *Lambda
		wantPanic bool
		want      string
	}{
		{
			name:      "Nil",
			input:     nil,
			wantPanic: true,
		},
		{
			name:      "Empty",
			input:     &Lambda{},
			wantPanic: true,
		},
		{
			name: "No Parameters",
			input: &Lambda{
				Parameters: nil,
				Expr:       *testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
			},
			want: "() => 1",
		},
		{
			name: "One Parameter",
			input: &Lambda{
				Parameters: []*LambdaParameter{
					{Label: "foo"},
				},
				Expr: *testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
			},
			want: "(foo) => 1",
		},
		{
			name: "Two Parameters",
			input: &Lambda{
				Parameters: []*LambdaParameter{
					{Label: "foo"},
					{Label: "bar"},
				},
				Expr: *testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
			},
			want: "(foo, bar) => 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

// /////////////////////////////////////

func TestLambdaParameter_Parsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *LambdaParameter
	}{
		{
			name:    "Only label",
			input:   `foo`,
			wantErr: false,
			want: &LambdaParameter{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
			},
		},
		{
			name:    "Label and type",
			input:   `foo: number`,
			wantErr: false,
			want: &LambdaParameter{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
				Type: &ParameterType{
					ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
					Ident:   &Ident{Parts: []string{"number"}},
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

func TestLambdaParameter_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		Input *LambdaParameter
		want  *LambdaParameter
	}{
		{
			name:  "Nil",
			Input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			Input: &LambdaParameter{},
			want:  &LambdaParameter{},
		},
		{
			name: "Label",
			Input: &LambdaParameter{
				Label: "foo",
			},
			want: &LambdaParameter{
				Label: "foo",
			},
		},
		{
			name: "Type",
			Input: &LambdaParameter{
				Type: &ParameterType{
					Ident: &Ident{Parts: []string{"foo"}},
				},
			},
			want: &LambdaParameter{
				Type: &ParameterType{
					Ident: &Ident{Parts: []string{"foo"}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*LambdaParameter](t, tt.want, tt.Input.Clone())
		})
	}
}

func TestLambdaParameter_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *LambdaParameter
		want  []Node
	}{
		{
			name:  "Empty",
			input: &LambdaParameter{},
			want:  nil,
		},
		{
			name:  "Label",
			input: &LambdaParameter{},
			want:  nil,
		},
		{
			name: "Type",
			input: &LambdaParameter{
				Type: &ParameterType{
					Ident: &Ident{Parts: []string{"foo"}},
				},
			},
			want: []Node{
				&ParameterType{
					Ident: &Ident{Parts: []string{"foo"}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestLambdaParameter_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     *LambdaParameter
		wantPanic bool
		want      string
	}{
		{
			name:      "Nil",
			input:     nil,
			wantPanic: true,
		},
		{
			name:  "Empty",
			input: &LambdaParameter{},
			want:  "",
		},
		{
			name: "No Type",
			input: &LambdaParameter{
				Label: "foo",
			},
			want: "foo",
		},
		{
			name: "Typed",
			input: &LambdaParameter{
				Label: "foo",
				Type: &ParameterType{
					Ident: &Ident{Parts: []string{"number"}},
				},
			},
			want: "foo: number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}
