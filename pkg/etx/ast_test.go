package etx

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAST_Parsing(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *AST
	}{
		{
			name: "One Attribute - no value",
			input: `
foo`[1:],
			wantErr: false,
			want: &AST{
				Items: []*RootItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Attribute: &Attribute{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Key:     "foo",
						},
					},
				},
			},
		},
		{
			name: "One Attribute - set value",
			input: `
foo = 1`[1:],
			wantErr: false,
			want: &AST{
				Items: []*RootItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Attribute: &Attribute{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Key:     "foo",
							Value: BuildTestExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
								Number: &ValueNumber{
									ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
									Value:   big.NewFloat(1),
									Source:  "1",
								},
							}),
						},
					},
				},
			},
		},
		{
			name: "Two Attributes - no value",
			input: `
foo
bar`[1:],
			wantErr: false,
			want: &AST{
				Items: []*RootItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Attribute: &Attribute{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Key:     "foo",
						},
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 2, Column: 1}},
						Attribute: &Attribute{
							ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 2, Column: 1}},
							Key:     "bar",
						},
					},
				},
			},
		},
		{
			name: "Two Attributes - set values",
			input: `
foo = 1
bar = 2`[1:],
			wantErr: false,
			want: &AST{
				Items: []*RootItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Attribute: &Attribute{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Key:     "foo",
							Value: BuildTestExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
								Number: &ValueNumber{
									ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
									Value:   big.NewFloat(1),
									Source:  "1",
								},
							}),
						},
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 2, Column: 1}},
						Attribute: &Attribute{
							ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 2, Column: 1}},
							Key:     "bar",
							Value: BuildTestExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 14, Line: 2, Column: 7}},
								Number: &ValueNumber{
									ASTNode: ASTNode{Pos: Position{Offset: 14, Line: 2, Column: 7}},
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
			name: "One Decl - set value",
			input: `
val foo = 1`[1:],
			wantErr: false,
			want: &AST{
				Items: []*RootItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Decl: &Decl{
							ASTNode:  ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							DeclType: "val",
							Label:    "foo",
							Type:     nil,
							Value: BuildTestExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 10, Line: 1, Column: 11}},
								Number: &ValueNumber{
									ASTNode: ASTNode{Pos: Position{Offset: 10, Line: 1, Column: 11}},
									Value:   big.NewFloat(1),
									Source:  "1",
								},
							}),
						},
					},
				},
			},
		},
		{
			name: "One Attribute and One Decl",
			input: `
foo = 1
val bar = 2`[1:],
			wantErr: false,
			want: &AST{
				Items: []*RootItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Attribute: &Attribute{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Key:     "foo",
							Value: BuildTestExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
								Number: &ValueNumber{
									ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
									Value:   big.NewFloat(1),
									Source:  "1",
								},
							}),
						},
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 2, Column: 1}},
						Decl: &Decl{
							ASTNode:  ASTNode{Pos: Position{Offset: 8, Line: 2, Column: 1}},
							DeclType: "val",
							Label:    "bar",
							Type:     nil,
							Value: BuildTestExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 18, Line: 2, Column: 11}},
								Number: &ValueNumber{
									ASTNode: ASTNode{Pos: Position{Offset: 18, Line: 2, Column: 11}},
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
			name: "One Decl - set value",
			input: `
val foo = 1`[1:],
			wantErr: false,
			want: &AST{
				Items: []*RootItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Decl: &Decl{
							ASTNode:  ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							DeclType: "val",
							Label:    "foo",
							Type:     nil,
							Value: BuildTestExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 10, Line: 1, Column: 11}},
								Number: &ValueNumber{
									ASTNode: ASTNode{Pos: Position{Offset: 10, Line: 1, Column: 11}},
									Value:   big.NewFloat(1),
									Source:  "1",
								},
							}),
						},
					},
				},
			},
		},
		{
			name:    "One empty Func",
			input:   `def foo() {}`,
			wantErr: false,
			want: &AST{
				Items: []*RootItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Func: &Func{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Label:   "foo",
						},
					},
				},
			},
		},
		{
			name: "Separated comments",
			input: `

// foo

// bar

`[1:],
			wantErr: false,
			want: &AST{
				Items: []*RootItem{
					{EmptyLine: "\n"},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Comment: &Comment{
							ASTNode:    ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							SingleLine: []string{"// foo"},
						},
					},
					{EmptyLine: "\n"},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 3, Column: 1}},
						Comment: &Comment{
							ASTNode:    ASTNode{Pos: Position{Offset: 8, Line: 3, Column: 1}},
							SingleLine: []string{"// bar"},
						},
					},
					{EmptyLine: "\n"},
				},
			},
		},
		{
			name: "Empty lines",
			input: `


`[1:],
			wantErr: false,
			want: &AST{
				Items: []*RootItem{
					{
						ASTNode:   ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						EmptyLine: "\n\n",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, false)
		})
	}
}

func TestAST_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		Input *AST
		want  *AST
	}{
		{
			name:  "Nil",
			Input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			Input: &AST{},
			want:  &AST{},
		},
		{
			name: "Item",
			Input: &AST{
				Items: []*RootItem{
					{Attribute: &Attribute{Key: "foo"}},
				},
			},
			want: &AST{
				Items: []*RootItem{
					{Attribute: &Attribute{Key: "foo"}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*AST](t, tt.want, tt.Input.Clone())
		})
	}
}

func TestAST_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *AST
		want  []Node
	}{
		{
			name:  "Empty",
			input: &AST{},
			want:  nil,
		},
		{
			name: "Item",
			input: &AST{
				Items: []*RootItem{
					{Attribute: &Attribute{Key: "foo"}},
				},
			},
			want: []Node{
				&RootItem{Attribute: &Attribute{Key: "foo"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestAST_FormattedString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     *AST
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
			input: &AST{},
			want:  "",
		},
		{
			name: "One item",
			input: &AST{
				Items: []*RootItem{
					{Attribute: &Attribute{Key: "foo"}},
				},
			},
			want: "foo\n",
		},
		{
			name: "Two items",
			input: &AST{
				Items: []*RootItem{
					{Attribute: &Attribute{Key: "foo"}},
					{Attribute: &Attribute{Key: "bar"}},
				},
			},
			want: `
foo

bar
`[1:],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

// /////////////////////////////////////

func TestRootItem_Parsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *RootItem
	}{
		{
			name:    "Decl",
			input:   "val foo",
			wantErr: false,
			want: &RootItem{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Decl: &Decl{
					ASTNode:  ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					DeclType: "val",
					Label:    "foo",
				},
			},
		},
		{
			name:    "Func",
			input:   "def foo() {}",
			wantErr: false,
			want: &RootItem{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Func: &Func{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Label:   "foo",
				},
			},
		},
		{
			name:    "Type",
			input:   "type foo object {}",
			wantErr: false,
			want: &RootItem{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Type: &Type{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Label:   "foo",
					Object: &TypeObject{
						ASTNode: ASTNode{Pos: Position{Offset: 17, Line: 1, Column: 18}},
					},
				},
			},
		},
		{
			name:    "Block",
			input:   "foo {}",
			wantErr: false,
			want: &RootItem{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Block: &Block{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Name:    "foo",
				},
			},
		},
		{
			name:    "Attribute",
			input:   "foo = 1",
			wantErr: false,
			want: &RootItem{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Attribute: &Attribute{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Key:     "foo",
					Value: BuildTestExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
						Number: &ValueNumber{
							ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
							big.NewFloat(1),
							"1",
						},
					}),
				},
			},
		},
		{
			name:    "Single-line Comment",
			input:   "// foo",
			wantErr: false,
			want: &RootItem{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Comment: &Comment{
					ASTNode:    ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					SingleLine: []string{"// foo"},
				},
			},
		},
		{
			name: "EmptyLine",
			input: `
`,
			wantErr: false,
			want: &RootItem{
				ASTNode:   ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				EmptyLine: "\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, true)
		})
	}
}

func TestRootItem_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *RootItem
		want  *RootItem
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &RootItem{},
			want:  &RootItem{},
		},
		{
			name: "ASTNode",
			input: &RootItem{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
			want: &RootItem{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
		},
		{
			name: "Decl",
			input: &RootItem{
				Decl: &Decl{},
			},
			want: &RootItem{
				Decl: &Decl{},
			},
		},
		{
			name: "Func",
			input: &RootItem{
				Func: &Func{},
			},
			want: &RootItem{
				Func: &Func{},
			},
		},
		{
			name: "Type",
			input: &RootItem{
				Type: &Type{},
			},
			want: &RootItem{
				Type: &Type{},
			},
		},
		{
			name: "Block",
			input: &RootItem{
				Block: &Block{},
			},
			want: &RootItem{
				Block: &Block{},
			},
		},
		{
			name: "Attribute",
			input: &RootItem{
				Attribute: &Attribute{},
			},
			want: &RootItem{
				Attribute: &Attribute{},
			},
		},
		{
			name: "Comment",
			input: &RootItem{
				Comment: &Comment{
					ASTNode:   ASTNode{},
					Multiline: "foo",
				},
			},
			want: &RootItem{
				Comment: &Comment{
					ASTNode:   ASTNode{},
					Multiline: "foo",
				},
			},
		},
		{
			name: "EmptyLine",
			input: &RootItem{
				EmptyLine: "\n",
			},
			want: &RootItem{
				EmptyLine: "\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*RootItem](t, tt.want, tt.input.Clone())
		})
	}
}

func TestRootItem_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *RootItem
		want  []Node
	}{
		{
			name:  "Empty",
			input: &RootItem{},
			want:  nil,
		},
		{
			name: "Decl",
			input: &RootItem{
				Decl: &Decl{},
			},
			want: []Node{
				&Decl{},
			},
		},
		{
			name: "Func",
			input: &RootItem{
				Func: &Func{},
			},
			want: []Node{
				&Func{},
			},
		},
		{
			name: "Type",
			input: &RootItem{
				Type: &Type{},
			},
			want: []Node{
				&Type{},
			},
		},
		{
			name: "Block",
			input: &RootItem{
				Block: &Block{},
			},
			want: []Node{
				&Block{},
			},
		},
		{
			name: "Attribute",
			input: &RootItem{
				Attribute: &Attribute{},
			},
			want: []Node{
				&Attribute{},
			},
		},
		{
			name: "Comment",
			input: &RootItem{
				Comment: &Comment{},
			},
			want: []Node{
				&Comment{},
			},
		},
		{
			name: "EmptyLine",
			input: &RootItem{
				EmptyLine: "\n",
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestRootItem_FormattedString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     *RootItem
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
			input:     &RootItem{},
			wantPanic: true,
		},
		{
			name: "Decl",
			input: &RootItem{
				Decl: &Decl{
					DeclType: "val",
					Label:    "foo",
				},
			},
			want: "val foo",
		},
		{
			name: "Func",
			input: &RootItem{
				Func: &Func{
					Label: "foo",
				},
			},
			want: "def foo() {}",
		},
		{
			name: "Type",
			input: &RootItem{
				Type: &Type{
					Label:  "foo",
					Object: &TypeObject{},
				},
			},
			want: "type foo object {\n}",
		},
		{
			name: "Block",
			input: &RootItem{
				Block: &Block{
					Name: "foo",
				},
			},
			want: "foo {}",
		},
		{
			name: "Attribute",
			input: &RootItem{
				Attribute: &Attribute{
					Key:   "foo",
					Value: BuildTestExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
				},
			},
			want: "foo: 1",
		},
		{
			name: "Comment",
			input: &RootItem{
				Comment: &Comment{SingleLine: []string{"// foo"}},
			},
			want: "// foo\n",
		},
		{
			name: "EmptyLine",
			input: &RootItem{
				EmptyLine: "\n",
			},
			want: "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

// /////////////////////////////////////

func TestParameterType_Parsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *ParameterType
	}{
		{
			name:    "Ident",
			input:   "foo",
			wantErr: false,
			want: &ParameterType{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Ident: &Ident{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Parts:   []string{"foo"},
				},
			},
		},
		{
			name:    "Func",
			input:   "() -> number",
			wantErr: false,
			want: &ParameterType{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Func: &FuncSignature{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Return: ParameterType{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
						Ident: &Ident{
							ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
							Parts:   []string{"number"},
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

func TestParameterType_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ParameterType
		want  *ParameterType
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name: "ASTNode",
			input: &ParameterType{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
			want: &ParameterType{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
		},
		{
			name:  "Empty",
			input: &ParameterType{},
			want:  &ParameterType{},
		},
		{
			name: "Ident",
			input: &ParameterType{
				Ident: &Ident{Parts: []string{"foo"}},
			},
			want: &ParameterType{
				Ident: &Ident{Parts: []string{"foo"}},
			},
		},
		{
			name: "Func",
			input: &ParameterType{
				Func: &FuncSignature{Return: ParameterType{Ident: &Ident{Parts: []string{"foo"}}}},
			},
			want: &ParameterType{
				Func: &FuncSignature{Return: ParameterType{Ident: &Ident{Parts: []string{"foo"}}}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*ParameterType](t, tt.want, tt.input.Clone())
		})
	}
}

func TestParameterType_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ParameterType
		want  []Node
	}{
		{
			name:  "Empty",
			input: &ParameterType{},
			want:  nil,
		},
		{
			name: "Ident",
			input: &ParameterType{
				Ident: &Ident{Parts: []string{"foo"}},
			},
			want: []Node{
				&Ident{Parts: []string{"foo"}},
			},
		},
		{
			name: "Func",
			input: &ParameterType{
				Func: &FuncSignature{Return: ParameterType{Ident: &Ident{Parts: []string{"foo"}}}},
			},
			want: []Node{
				&FuncSignature{Return: ParameterType{Ident: &Ident{Parts: []string{"foo"}}}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestParameterType_FormattedString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     *ParameterType
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
			input:     &ParameterType{},
			wantPanic: true,
		},
		{
			name: "Ident",
			input: &ParameterType{
				Ident: &Ident{Parts: []string{"foo"}},
			},
			want: Ident{Parts: []string{"foo"}}.FormattedString(),
		},
		{
			name: "Func",
			input: &ParameterType{
				Func: &FuncSignature{Return: ParameterType{Ident: &Ident{Parts: []string{"foo"}}}},
			},
			want: FuncSignature{Return: ParameterType{Ident: &Ident{Parts: []string{"foo"}}}}.FormattedString(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

// /////////////////////////////////////

func TestFuncSignature_Parsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *FuncSignature
	}{
		{
			name:    "no parameters",
			input:   "() -> number",
			wantErr: false,
			want: &FuncSignature{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Return: ParameterType{
					ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
					Ident: &Ident{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
						Parts:   []string{"number"},
					},
				},
			},
		},
		{
			name:    "one parameter",
			input:   "(bool) -> number",
			wantErr: false,
			want: &FuncSignature{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Parameters: []*ParameterType{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Ident: &Ident{
							ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
							Parts:   []string{"bool"},
						},
					},
				},
				Return: ParameterType{
					ASTNode: ASTNode{Pos: Position{Offset: 10, Line: 1, Column: 11}},
					Ident: &Ident{
						ASTNode: ASTNode{Pos: Position{Offset: 10, Line: 1, Column: 11}},
						Parts:   []string{"number"},
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

func TestFuncSignature_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *FuncSignature
		want  *FuncSignature
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &FuncSignature{},
			want:  &FuncSignature{},
		},
		{
			name: "ASTNode",
			input: &FuncSignature{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
			want: &FuncSignature{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
		},
		{
			name: "Parameters",
			input: &FuncSignature{
				Parameters: []*ParameterType{
					{Ident: &Ident{Parts: []string{"foo"}}},
					{Ident: &Ident{Parts: []string{"bar"}}},
				},
			},
			want: &FuncSignature{
				Parameters: []*ParameterType{
					{Ident: &Ident{Parts: []string{"foo"}}},
					{Ident: &Ident{Parts: []string{"bar"}}},
				},
			},
		},
		{
			name: "Return",
			input: &FuncSignature{
				Return: ParameterType{Ident: &Ident{Parts: []string{"bar"}}},
			},
			want: &FuncSignature{
				Return: ParameterType{Ident: &Ident{Parts: []string{"bar"}}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*FuncSignature](t, tt.want, tt.input.Clone())
		})
	}
}

func TestFuncSignature_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *FuncSignature
		want  []Node
	}{
		{
			name:  "Empty",
			input: &FuncSignature{},
			want: []Node{
				&ParameterType{},
			},
		},
		{
			name: "Parameters",
			input: &FuncSignature{
				Parameters: []*ParameterType{
					{Ident: &Ident{Parts: []string{"foo"}}},
					{Ident: &Ident{Parts: []string{"bar"}}},
				},
			},
			want: []Node{
				&ParameterType{Ident: &Ident{Parts: []string{"foo"}}},
				&ParameterType{Ident: &Ident{Parts: []string{"bar"}}},
				&ParameterType{},
			},
		},
		{
			name: "Return",
			input: &FuncSignature{
				Return: ParameterType{Ident: &Ident{Parts: []string{"bar"}}},
			},
			want: []Node{
				&ParameterType{Ident: &Ident{Parts: []string{"bar"}}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestFuncSignature_FormattedString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     *FuncSignature
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
			input:     &FuncSignature{},
			wantPanic: true,
		},
		{
			name: "No parameters - return",
			input: &FuncSignature{
				Return: ParameterType{Ident: &Ident{Parts: []string{"foo"}}},
			},
			want: "() -> foo",
		},
		{
			name: "One parameter - return",
			input: &FuncSignature{
				Parameters: []*ParameterType{
					{Ident: &Ident{Parts: []string{"foo"}}},
				},
				Return: ParameterType{Ident: &Ident{Parts: []string{"bar"}}},
			},
			want: "(foo) -> bar",
		},
		{
			name: "Two parameters - return",
			input: &FuncSignature{
				Parameters: []*ParameterType{
					{Ident: &Ident{Parts: []string{"foo"}}},
					{Ident: &Ident{Parts: []string{"bar"}}},
				},
				Return: ParameterType{Ident: &Ident{Parts: []string{"baz"}}},
			},
			want: "(foo, bar) -> baz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}
