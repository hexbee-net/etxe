package etx

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
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
			name:    "Empty body, no params, one return",
			input:   `def foo() bool {}`,
			wantErr: false,
			want: &Func{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
				Return: []*ParameterType{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 10, Line: 1, Column: 11}},
						Ident:   &Ident{Parts: []string{"bool"}},
					},
				},
			},
		},
		{
			name:    "Empty body, no params, two returns",
			input:   `def foo() (bool, number) {}`,
			wantErr: false,
			want: &Func{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
				Return: []*ParameterType{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 11, Line: 1, Column: 12}},
						Ident:   &Ident{Parts: []string{"bool"}},
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 17, Line: 1, Column: 18}},
						Ident:   &Ident{Parts: []string{"number"}},
					},
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
			input:   `def foo() ((int) -> bool) {}`,
			wantErr: false,
			want: &Func{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
				Return: []*ParameterType{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 11, Line: 1, Column: 12}},
						Func: &FuncSignature{
							ASTNode: ASTNode{Pos: Position{Offset: 11, Line: 1, Column: 12}},
							Parameters: []*ParameterType{
								{
									ASTNode: ASTNode{Pos: Position{Offset: 12, Line: 1, Column: 13}},
									Ident:   &Ident{Parts: []string{"int"}},
								},
							},
							Return: ParameterType{
								ASTNode: ASTNode{Pos: Position{Offset: 20, Line: 1, Column: 21}},
								Ident:   &Ident{Parts: []string{"bool"}},
							},
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
							Type: &ParameterType{
								ASTNode: ASTNode{Pos: Position{Offset: 20, Line: 2, Column: 9}},
								Ident:   &Ident{Parts: []string{"number"}},
							},
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
							Type: &ParameterType{
								ASTNode: ASTNode{Pos: Position{Offset: 20, Line: 2, Column: 9}},
								Ident:   &Ident{Parts: []string{"number"}},
							},
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
		{
			name: "Single-line comment",
			input: `
// foo
def foo() {}`[1:],
			wantErr: false,
			want: &Func{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Comment: &Comment{
					ASTNode:    ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					SingleLine: []string{"// foo"},
				},
				Label: "foo",
			},
		},
		{
			name: "Single-line comment - separated",
			input: `
// foo

def foo() {}`[1:],
			wantErr: true,
		},
		{
			name: "Multi-line comment",
			input: `
/* foo */
def foo() {}`[1:],
			wantErr: false,
			want: &Func{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Comment: &Comment{
					ASTNode:   ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Multiline: "/* foo */\n",
				},
				Label: "foo",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, true)
		})
	}
}

func TestFunc_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *Func
		want  *Func
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &Func{},
			want:  &Func{},
		},
		{
			name: "ASTNode",
			input: &Func{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
			want: &Func{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
		},
		{
			name: "Comments",
			input: &Func{
				Comment: &Comment{Multiline: "foo"},
			},
			want: &Func{
				Comment: &Comment{Multiline: "foo"},
			},
		},
		{
			name: "Label",
			input: &Func{
				Label: "foo",
			},
			want: &Func{
				Label: "foo",
			},
		},
		{
			name: "Parameters",
			input: &Func{
				Parameters: []*FuncParameter{
					{Label: "bar"},
				},
			},
			want: &Func{
				Parameters: []*FuncParameter{
					{Label: "bar"},
				},
			},
		},
		{
			name: "Return",
			input: &Func{
				Return: []*ParameterType{
					{Ident: &Ident{Parts: []string{"number"}}},
				},
			},
			want: &Func{
				Return: []*ParameterType{
					{Ident: &Ident{Parts: []string{"number"}}},
				},
			},
		},
		{
			name: "Body",
			input: &Func{
				Body: []*FuncStatement{
					{Expr: testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}})},
				},
			},
			want: &Func{
				Body: []*FuncStatement{
					{Expr: testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}})},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*Func](t, tt.want, tt.input.Clone())
		})
	}
}

func TestFunc_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *Func
		want  []Node
	}{
		{
			name:  "Empty",
			input: &Func{},
			want:  nil,
		},
		{
			name: "Comment",
			input: &Func{
				Comment: &Comment{Multiline: "foo"},
			},
			want: []Node{
				&Comment{Multiline: "foo"},
			},
		},
		{
			name: "Label",
			input: &Func{
				Label: "foo",
			},
			want: nil,
		},
		{
			name: "Parameters",
			input: &Func{
				Parameters: []*FuncParameter{
					{Label: "bar"},
				},
			},
			want: []Node{
				&FuncParameter{Label: "bar"},
			},
		},
		{
			name: "Return",
			input: &Func{
				Return: []*ParameterType{
					{Ident: &Ident{Parts: []string{"number"}}},
				},
			},
			want: []Node{
				&ParameterType{Ident: &Ident{Parts: []string{"number"}}},
			},
		},
		{
			name: "Body",
			input: &Func{
				Body: []*FuncStatement{
					{Expr: testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}})},
				},
			},
			want: []Node{
				&FuncStatement{
					Expr: testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
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

func TestFunc_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     *Func
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
			input: &Func{},
			want:  "def () {}",
		},
		{
			name: "Label - no params - no return - no body",
			input: &Func{
				Label: "foo",
			},
			want: "def foo() {}",
		},
		{
			name: "Label - one param - no return - no body",
			input: &Func{
				Label: "foo",
				Parameters: []*FuncParameter{
					{Label: "bar"},
				},
			},
			want: "def foo(bar) {}",
		},
		{
			name: "Label - two params - no return - no body",
			input: &Func{
				Label: "foo",
				Parameters: []*FuncParameter{
					{Label: "bar"},
					{Label: "baz"},
				},
			},
			want: "def foo(bar, baz) {}",
		},
		{
			name: "Label - no params - one return type - no body",
			input: &Func{
				Label: "foo",
				Return: []*ParameterType{
					{Ident: &Ident{Parts: []string{"number"}}},
				},
			},
			want: "def foo() number {}",
		},
		{
			name: "Label - no params - two return types - no body",
			input: &Func{
				Label: "foo",
				Return: []*ParameterType{
					{Ident: &Ident{Parts: []string{"number"}}},
					{Ident: &Ident{Parts: []string{"bool"}}},
				},
			},
			want: "def foo() (number, bool) {}",
		},
		{
			name: "Label - no params - no return - body",
			input: &Func{
				Label: "foo",
				Body: []*FuncStatement{
					{Expr: testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}})},
				},
			},
			want: `
def foo() {
	1
}`[1:],
		},
		{
			name: "Comment",
			input: &Func{
				Comment: &Comment{SingleLine: []string{"// foo"}},
			},
			want: `
// foo
def () {}`[1:],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

// /////////////////////////////////////

func TestFuncParameter_Parsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *FuncParameter
	}{
		{
			name:    "ident - no type",
			input:   "foo",
			wantErr: false,
			want: &FuncParameter{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
			},
		},
		{
			name:    "ident - type",
			input:   "foo: number",
			wantErr: false,
			want: &FuncParameter{
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

func TestFuncParameter_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *FuncParameter
		want  *FuncParameter
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &FuncParameter{},
			want:  &FuncParameter{},
		},
		{
			name: "ASTNode",
			input: &FuncParameter{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
			want: &FuncParameter{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
		},
		{
			name: "Label",
			input: &FuncParameter{
				Label: "foo",
			},
			want: &FuncParameter{
				Label: "foo",
			},
		},
		{
			name: "Type",
			input: &FuncParameter{
				Type: &ParameterType{Ident: &Ident{Parts: []string{"number"}}},
			},
			want: &FuncParameter{
				Type: &ParameterType{Ident: &Ident{Parts: []string{"number"}}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*FuncParameter](t, tt.want, tt.input.Clone())
		})
	}
}

func TestFuncParameter_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *FuncParameter
		want  []Node
	}{
		{
			name:  "Empty",
			input: &FuncParameter{},
			want:  nil,
		},
		{
			name: "Label",
			input: &FuncParameter{
				Label: "foo",
			},
			want: nil,
		},
		{
			name: "Type",
			input: &FuncParameter{
				Type: &ParameterType{Ident: &Ident{Parts: []string{"number"}}},
			},
			want: []Node{
				&ParameterType{Ident: &Ident{Parts: []string{"number"}}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestFuncParameter_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     *FuncParameter
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
			input: &FuncParameter{},
			want:  "",
		},
		{
			name: "Label only",
			input: &FuncParameter{
				Label: "foo",
			},
			want: "foo",
		},
		{
			name: "Label and type",
			input: &FuncParameter{
				Label: "foo",
				Type:  &ParameterType{Ident: &Ident{Parts: []string{"number"}}},
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

// /////////////////////////////////////

func TestFuncStatement_Parsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *FuncStatement
	}{
		{
			name:    "Decl",
			input:   "val foo",
			wantErr: false,
			want: &FuncStatement{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Decl: &FuncDecl{
					ASTNode:  ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					DeclType: "val",
					Label:    "foo",
				},
			},
		},
		{
			name:    "Expr",
			input:   "1",
			wantErr: false,
			want: &FuncStatement{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Expr: testBuildExprTree[*Expr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
			},
		},
		{
			name:    "Single-line comment",
			input:   `// foo`,
			wantErr: false,
			want: &FuncStatement{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Comment: &Comment{
					ASTNode:    ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					SingleLine: []string{"// foo"},
				},
			},
		},
		{
			name:    "Multi-line comment",
			input:   `/* foo */`,
			wantErr: false,
			want: &FuncStatement{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Comment: &Comment{
					ASTNode:   ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Multiline: "/* foo */",
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

func TestFuncStatement_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *FuncStatement
		want  *FuncStatement
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &FuncStatement{},
			want:  &FuncStatement{},
		},
		{
			name: "ASTNode",
			input: &FuncStatement{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
			want: &FuncStatement{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
		},
		{
			name: "Comments",
			input: &FuncStatement{
				Comment: &Comment{Multiline: "foo"},
			},
			want: &FuncStatement{
				Comment: &Comment{Multiline: "foo"},
			},
		},
		{
			name: "Decl",
			input: &FuncStatement{
				Decl: &FuncDecl{
					DeclType: "val",
					Label:    "foo",
					Value:    testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
				},
			},
			want: &FuncStatement{
				Decl: &FuncDecl{
					DeclType: "val",
					Label:    "foo",
					Value:    testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
				},
			},
		},
		{
			name: "Expr",
			input: &FuncStatement{
				Expr: testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
			},
			want: &FuncStatement{
				Expr: testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
			},
		},
		{
			name: "EmptyLine",
			input: &FuncStatement{
				EmptyLine: "\n",
			},
			want: &FuncStatement{
				EmptyLine: "\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*FuncStatement](t, tt.want, tt.input.Clone())
		})
	}
}

func TestFuncStatement_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *FuncStatement
		want  []Node
	}{
		{
			name:  "Empty",
			input: &FuncStatement{},
			want:  nil,
		},
		{
			name: "Comment",
			input: &FuncStatement{
				Comment: &Comment{Multiline: "foo"},
			},
			want: []Node{
				&Comment{Multiline: "foo"},
			},
		},
		{
			name: "Decl",
			input: &FuncStatement{
				Decl: &FuncDecl{
					DeclType: "val",
					Label:    "foo",
					Value:    testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
				},
			},
			want: []Node{
				&FuncDecl{
					DeclType: "val",
					Label:    "foo",
					Value:    testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
				},
			},
		},
		{
			name: "Expr",
			input: &FuncStatement{
				Expr: testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
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

func TestFuncStatement_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     *FuncStatement
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
			input: &FuncStatement{},
			want:  "",
		},
		{
			name: "EmptyLine",
			input: &FuncStatement{
				EmptyLine: "\n",
			},
			want: "\n",
		},
		{
			name: "Decl",
			input: &FuncStatement{
				Decl: &FuncDecl{
					DeclType: "val",
					Label:    "foo",
					Value:    testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
				},
			},
			want: "val foo = 1",
		},
		{
			name: "Expr",
			input: &FuncStatement{
				Expr: testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
			},
			want: "1",
		},
		{
			name: "Comment",
			input: &FuncStatement{
				Comment: &Comment{SingleLine: []string{"// foo"}},
			},
			want: "// foo\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

// /////////////////////////////////////

func TestFuncDecl_Parsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *FuncDecl
	}{
		{
			name:    "const - no type",
			input:   "const foo",
			wantErr: false,
			want: &FuncDecl{
				ASTNode:  ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				DeclType: "const",
				Label:    "foo",
			},
		},
		{
			name:    "val - no type",
			input:   "val foo",
			wantErr: false,
			want: &FuncDecl{
				ASTNode:  ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				DeclType: "val",
				Label:    "foo",
			},
		},
		{
			name:    "unknown - no type",
			input:   "foo bar",
			wantErr: true,
		},

		{
			name:    "type - no value",
			input:   "val foo: number",
			wantErr: false,
			want: &FuncDecl{
				ASTNode:  ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				DeclType: "val",
				Label:    "foo",
				Type: &ParameterType{
					ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 1, Column: 10}},
					Ident:   &Ident{Parts: []string{"number"}},
				},
			},
		},
		{
			name:    "no type - value",
			input:   "val foo = 1",
			wantErr: false,
			want: &FuncDecl{
				ASTNode:  ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				DeclType: "val",
				Label:    "foo",
				Value: testBuildExprTree[*Expr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 10, Line: 1, Column: 11}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
			},
		},
		{
			name:    "type - value",
			input:   "val foo: number = 1",
			wantErr: false,
			want: &FuncDecl{
				ASTNode:  ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				DeclType: "val",
				Label:    "foo",
				Type: &ParameterType{
					ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 1, Column: 10}},
					Ident:   &Ident{Parts: []string{"number"}},
				},
				Value: testBuildExprTree[*Expr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 18, Line: 1, Column: 19}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, true)
		})
	}
}

func TestFuncDecl_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		Input *FuncDecl
		want  *FuncDecl
	}{
		{
			name:  "Nil",
			Input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			Input: &FuncDecl{},
			want:  &FuncDecl{},
		},
		{
			name: "ASTNode",
			Input: &FuncDecl{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
			want: &FuncDecl{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
		},
		{
			name: "DeclType",
			Input: &FuncDecl{
				DeclType: "val",
			},
			want: &FuncDecl{
				DeclType: "val",
			},
		},
		{
			name: "Label",
			Input: &FuncDecl{
				Label: "foo",
			},
			want: &FuncDecl{
				Label: "foo",
			},
		},
		{
			name: "Type",
			Input: &FuncDecl{
				Type: &ParameterType{Ident: &Ident{Parts: []string{"number"}}},
			},
			want: &FuncDecl{
				Type: &ParameterType{Ident: &Ident{Parts: []string{"number"}}},
			},
		},
		{
			name: "Value",
			Input: &FuncDecl{
				Value: testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
			},
			want: &FuncDecl{
				Value: testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*FuncDecl](t, tt.want, tt.Input.Clone())
		})
	}
}

func TestFuncDecl_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *FuncDecl
		want  []Node
	}{
		{
			name:  "Empty",
			input: &FuncDecl{},
			want:  nil,
		},
		{
			name: "DeclType",
			input: &FuncDecl{
				DeclType: "val",
			},
			want: nil,
		},
		{
			name: "Label",
			input: &FuncDecl{
				Label: "foo",
			},
			want: nil,
		},
		{
			name: "Type",
			input: &FuncDecl{
				Type: &ParameterType{Ident: &Ident{Parts: []string{"number"}}},
			},
			want: []Node{
				&ParameterType{Ident: &Ident{Parts: []string{"number"}}},
			},
		},
		{
			name: "Value",
			input: &FuncDecl{
				Value: testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
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

func TestFuncDecl_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     *FuncDecl
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
			input: &FuncDecl{},
			want:  "",
		},
		{
			name: "Label only",
			input: &FuncDecl{
				DeclType: "val",
				Label:    "foo",
			},
			want: "val foo",
		},
		{
			name: "Label and type",
			input: &FuncDecl{
				DeclType: "val",
				Label:    "foo",
				Type:     &ParameterType{Ident: &Ident{Parts: []string{"number"}}},
			},
			want: "val foo: number",
		},
		{
			name: "Label and value",
			input: &FuncDecl{
				DeclType: "val",
				Label:    "foo",
				Value:    testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
			},
			want: "val foo = 1",
		},
		{
			name: "Label, type and value",
			input: &FuncDecl{
				DeclType: "val",
				Label:    "foo",
				Type:     &ParameterType{Ident: &Ident{Parts: []string{"number"}}},
				Value:    testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
			},
			want: "val foo: number = 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}
