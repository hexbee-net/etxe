package etx

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestType_Parsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Type
	}{
		{
			name:    "Enum - Empty",
			input:   `type foo enum {}`,
			wantErr: false,
			want: &Type{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
				Enum: &TypeEnum{
					ASTNode: ASTNode{Pos: Position{Offset: 15, Line: 1, Column: 16}},
				},
				Object: nil,
			},
		},
		{
			name: "Enum",
			input: `
type foo enum {
  bar: 1
}`[1:],
			wantErr: false,
			want: &Type{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
				Enum: &TypeEnum{
					ASTNode: ASTNode{Pos: Position{Offset: 15, Line: 1, Column: 16}},
					Items: []*TypeEnumItem{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 18, Line: 2, Column: 3}},
							Label:   "bar",
							Value: *testBuildExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 23, Line: 2, Column: 8}},
								Number: &ValueNumber{
									ASTNode: ASTNode{Pos: Position{Offset: 23, Line: 2, Column: 8}},
									Value:   big.NewFloat(1),
									Source:  "1",
								},
							}),
						},
					},
				},
				Object: nil,
			},
		},

		{
			name:    "Object - Empty",
			input:   `type foo object {}`,
			wantErr: false,
			want: &Type{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
				Enum:    nil,
				Object: &TypeObject{
					ASTNode: ASTNode{Pos: Position{Offset: 17, Line: 1, Column: 18}},
				},
			},
		},
		{
			name: "Object",
			input: `
type foo object {
	foo: number
}`[1:],
			wantErr: false,
			want: &Type{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
				Enum:    nil,
				Object: &TypeObject{
					ASTNode: ASTNode{Pos: Position{Offset: 17, Line: 1, Column: 18}},
					Items: []*TypeObjectItem{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 19, Line: 2, Column: 2}},
							Label:   "foo",
							Type: ParameterType{
								ASTNode: ASTNode{Pos: Position{Offset: 24, Line: 2, Column: 7}},
								Ident: &Ident{
									ASTNode: ASTNode{Pos: Position{Offset: 24, Line: 2, Column: 7}},
									Parts: []string{
										"number",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Single-line comment",
			input: `
// foo
type foo enum {}`[1:],
			wantErr: false,
			want: &Type{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Comment: &Comment{
					ASTNode:    ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					SingleLine: []string{"// foo"},
				},
				Label: "foo",
				Enum: &TypeEnum{
					ASTNode: ASTNode{Pos: Position{Offset: 22, Line: 2, Column: 16}},
				},
				Object: nil,
			},
		},
		{
			name: "Single-line comment - separated",
			input: `
// foo

type foo enum {}`[1:],
			wantErr: true,
		},
		{
			name: "Single-line comment",
			input: `
/* foo */
type foo enum {}`[1:],
			wantErr: false,
			want: &Type{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Comment: &Comment{
					ASTNode:   ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Multiline: "/* foo */\n",
				},
				Label: "foo",
				Enum: &TypeEnum{
					ASTNode: ASTNode{Pos: Position{Offset: 25, Line: 2, Column: 16}},
				},
				Object: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, true)
		})
	}
}

func TestType_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		Input *Type
		want  *Type
	}{
		{
			name:  "Nil",
			Input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			Input: &Type{},
			want:  &Type{},
		},
		{
			name: "ASTNode",
			Input: &Type{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
			want: &Type{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
		},
		{
			name: "Comments",
			Input: &Type{
				Comment: &Comment{Multiline: "foo"},
			},
			want: &Type{
				Comment: &Comment{Multiline: "foo"},
			},
		},
		{
			name: "Label",
			Input: &Type{
				Label: "foo",
			},
			want: &Type{
				Label: "foo",
			},
		},
		{
			name: "Enum",
			Input: &Type{
				Enum: &TypeEnum{Items: []*TypeEnumItem{{Label: "foo"}}},
			},
			want: &Type{
				Enum: &TypeEnum{Items: []*TypeEnumItem{{Label: "foo"}}},
			},
		},
		{
			name: "Object",
			Input: &Type{
				Object: &TypeObject{Items: []*TypeObjectItem{{Label: "foo"}}},
			},
			want: &Type{
				Object: &TypeObject{Items: []*TypeObjectItem{{Label: "foo"}}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*Type](t, tt.want, tt.Input.Clone())
		})
	}
}

func TestType_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *Type
		want  []Node
	}{
		{
			name:  "Empty",
			input: &Type{},
			want:  nil,
		},
		{
			name: "Comment",
			input: &Type{
				Comment: &Comment{Multiline: "foo"},
			},
			want: []Node{
				&Comment{Multiline: "foo"},
			},
		},
		{
			name: "Label",
			input: &Type{
				Label: "foo",
			},
			want: nil,
		},
		{
			name: "Enum",
			input: &Type{
				Enum: &TypeEnum{},
			},
			want: []Node{
				&TypeEnum{},
			},
		},
		{
			name: "Object",
			input: &Type{
				Object: &TypeObject{},
			},
			want: []Node{
				&TypeObject{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestType_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     *Type
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
			input:     &Type{},
			wantPanic: true,
		},
		{
			name: "Enum",
			input: &Type{
				Label: "foo",
				Enum: &TypeEnum{
					Items: []*TypeEnumItem{
						{
							Label: "bar",
							Value: *testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
						},
					},
				},
			},
			want: `
type foo enum {
	bar: 1
}`[1:],
		},
		{
			name: "Object",
			input: &Type{
				Label: "foo",
				Object: &TypeObject{
					Items: []*TypeObjectItem{
						{
							Label: "bar",
							Type:  ParameterType{Ident: &Ident{Parts: []string{"number"}}},
						},
					},
				},
			},
			want: `
type foo object {
	bar: number
}`[1:],
		},
		{
			name: "Comment",
			input: &Type{
				Comment: &Comment{SingleLine: []string{"// foo"}},
				Label:   "bar",
				Enum: &TypeEnum{
					Items: []*TypeEnumItem{
						{
							Label: "baz",
							Value: *testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
						},
					},
				},
			},
			want: `
// foo
type bar enum {
	baz: 1
}`[1:],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

// /////////////////////////////////////

func TestTypeEnum_Parsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *TypeEnum
	}{
		{
			name:    "Empty",
			input:   ``,
			wantErr: false,
			want: &TypeEnum{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Items:   nil,
			},
		},
		{
			name: "One value",
			input: `
foo: 1`,
			wantErr: false,
			want: &TypeEnum{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Items: []*TypeEnumItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 1}},
						Label:   "foo",
						Value: *testBuildExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 6}},
							Number: &ValueNumber{
								ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 6}},
								Value:   big.NewFloat(1),
								Source:  "1",
							},
						}),
					},
				},
			},
		},
		{
			name: "Two values",
			input: `
foo: 1
bar: 2`,
			wantErr: false,
			want: &TypeEnum{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Items: []*TypeEnumItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 1}},
						Label:   "foo",
						Value: *testBuildExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 6}},
							Number: &ValueNumber{
								ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 6}},
								Value:   big.NewFloat(1),
								Source:  "1",
							},
						}),
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 3, Column: 1}},
						Label:   "bar",
						Value: *testBuildExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 3, Column: 6}},
							Number: &ValueNumber{
								ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 3, Column: 6}},
								Value:   big.NewFloat(2),
								Source:  "2",
							},
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

func TestTypeEnum_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		Input *TypeEnum
		want  *TypeEnum
	}{
		{
			name:  "Nil",
			Input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			Input: &TypeEnum{},
			want:  &TypeEnum{},
		},
		{
			name: "ASTNode",
			Input: &TypeEnum{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
			want: &TypeEnum{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
		},
		{
			name: "Items",
			Input: &TypeEnum{
				Items: []*TypeEnumItem{{Label: "foo"}},
			},
			want: &TypeEnum{
				Items: []*TypeEnumItem{{Label: "foo"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*TypeEnum](t, tt.want, tt.Input.Clone())
		})
	}
}

func TestTypeEnum_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *TypeEnum
		want  []Node
	}{
		{
			name:  "Empty",
			input: &TypeEnum{},
			want:  nil,
		},
		{
			name: "Items",
			input: &TypeEnum{
				Items: []*TypeEnumItem{{Label: "foo"}},
			},
			want: []Node{
				&TypeEnumItem{Label: "foo"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestTypeEnum_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     *TypeEnum
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
			input: &TypeEnum{},
			want:  "",
		},
		{
			name: "One Value",
			input: &TypeEnum{
				Items: []*TypeEnumItem{
					{
						Label: "foo",
						Value: *testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
					},
				},
			},
			want: `
foo: 1`,
		},
		{
			name: "Two Values",
			input: &TypeEnum{
				Items: []*TypeEnumItem{
					{
						Label: "foo",
						Value: *testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
					},
					{
						Label: "bar",
						Value: *testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "2"}}),
					},
				},
			},
			want: `
foo: 1
bar: 2`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

// /////////////////////////////////////

func TestTypeEnumItem_Parsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *TypeEnumItem
	}{
		{
			name:    "Value",
			input:   "foo: 1",
			wantErr: false,
			want: &TypeEnumItem{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
				Value: *testBuildExprTree[*Expr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
					Number: &ValueNumber{
						ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
						Value:   big.NewFloat(1),
						Source:  "1",
					},
				}),
			},
		},
		{
			name: "Single-line comment",
			input: `
// foo`[1:],
			wantErr: false,
			want: &TypeEnumItem{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Comment: &Comment{
					ASTNode:    ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					SingleLine: []string{"// foo"},
				},
			},
		},
		{
			name: "Multi-line comment",
			input: `
/* foo */`[1:],
			wantErr: false,
			want: &TypeEnumItem{
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

func TestTypeEnumItem_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		Input *TypeEnumItem
		want  *TypeEnumItem
	}{
		{
			name:  "Nil",
			Input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			Input: &TypeEnumItem{},
			want:  &TypeEnumItem{},
		},
		{
			name: "EmptyLine",
			Input: &TypeEnumItem{
				EmptyLine: "\n",
			},
			want: &TypeEnumItem{
				EmptyLine: "\n",
			},
		},
		{
			name: "ASTNode",
			Input: &TypeEnumItem{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
			want: &TypeEnumItem{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
		},
		{
			name: "Comments",
			Input: &TypeEnumItem{
				Comment: &Comment{Multiline: "foo"},
			},
			want: &TypeEnumItem{
				Comment: &Comment{Multiline: "foo"},
			},
		},
		{
			name: "Label",
			Input: &TypeEnumItem{
				Label: "foo",
			},
			want: &TypeEnumItem{
				Label: "foo",
			},
		},
		{
			name: "Value",
			Input: &TypeEnumItem{
				Value: *testBuildExprTree[*Expr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
			},
			want: &TypeEnumItem{
				Value: *testBuildExprTree[*Expr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*TypeEnumItem](t, tt.want, tt.Input.Clone())
		})
	}
}

func TestTypeEnumItem_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *TypeEnumItem
		want  []Node
	}{
		{
			name:  "Empty",
			input: &TypeEnumItem{},
			want: []Node{
				&Expr{},
			},
		},
		{
			name: "Comment",
			input: &TypeEnumItem{
				Comment: &Comment{Multiline: "foo"},
			},
			want: []Node{
				&Comment{Multiline: "foo"},
				&Expr{},
			},
		},
		{
			name: "Label",
			input: &TypeEnumItem{
				Label: "foo",
			},
			want: []Node{
				&Expr{},
			},
		},
		{
			name: "Value",
			input: &TypeEnumItem{
				Value: *testBuildExprTree[*Expr](t, &Value{
					Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"},
				}),
			},
			want: []Node{
				testBuildExprTree[*Expr](t, &Value{
					Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"},
				}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestTypeEnumItem_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     *TypeEnumItem
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
			input:     &TypeEnumItem{},
			wantPanic: true,
		},
		{
			name: "Value",
			input: &TypeEnumItem{
				Label: "foo",
				Value: *testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
			},
			want: "foo: 1",
		},
		{
			name: "Comment",
			input: &TypeEnumItem{
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

func TestTypeObject_Parsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *TypeObject
	}{
		{
			name:    "Empty",
			input:   ``,
			wantErr: false,
			want: &TypeObject{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Items:   nil,
			},
		},
		{
			name: "One declaration",
			input: `
foo: number`,
			wantErr: false,
			want: &TypeObject{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Items: []*TypeObjectItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 1}},
						Label:   "foo",
						Type: ParameterType{
							ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 6}},
							Ident: &Ident{
								ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 6}},
								Parts: []string{
									"number",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Two declarations",
			input: `
foo: number
bar: bool`,
			wantErr: false,
			want: &TypeObject{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Items: []*TypeObjectItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 1}},
						Label:   "foo",
						Type: ParameterType{
							ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 6}},
							Ident: &Ident{
								ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 6}},
								Parts: []string{
									"number",
								},
							},
						},
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 3, Column: 1}},
						Label:   "bar",
						Type: ParameterType{
							ASTNode: ASTNode{Pos: Position{Offset: 18, Line: 3, Column: 6}},
							Ident: &Ident{
								ASTNode: ASTNode{Pos: Position{Offset: 18, Line: 3, Column: 6}},
								Parts: []string{
									"bool",
								},
							},
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

func TestTypeObject_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		Input *TypeObject
		want  *TypeObject
	}{
		{
			name:  "Nil",
			Input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			Input: &TypeObject{},
			want:  &TypeObject{},
		},
		{
			name: "ASTNode",
			Input: &TypeObject{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
			want: &TypeObject{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
		},
		{
			name: "Items",
			Input: &TypeObject{
				Items: []*TypeObjectItem{{Label: "foo"}},
			},
			want: &TypeObject{
				Items: []*TypeObjectItem{{Label: "foo"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*TypeObject](t, tt.want, tt.Input.Clone())
		})
	}
}

func TestTypeObject_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *TypeObject
		want  []Node
	}{
		{
			name:  "Empty",
			input: &TypeObject{},
			want:  nil,
		},
		{
			name: "Items",
			input: &TypeObject{
				Items: []*TypeObjectItem{{Label: "foo"}},
			},
			want: []Node{
				&TypeObjectItem{Label: "foo"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestTypeObject_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     *TypeObject
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
			input: &TypeObject{},
			want:  "",
		},
		{
			name: "One Type",
			input: &TypeObject{
				Items: []*TypeObjectItem{
					{
						Label: "foo",
						Type:  ParameterType{Ident: &Ident{Parts: []string{"number"}}},
					},
				},
			},
			want: `
foo: number`,
		},
		{
			name: "Two Types",
			input: &TypeObject{
				Items: []*TypeObjectItem{
					{
						Label: "foo",
						Type:  ParameterType{Ident: &Ident{Parts: []string{"number"}}},
					},
					{
						Label: "bar",
						Type:  ParameterType{Ident: &Ident{Parts: []string{"string"}}},
					},
				},
			},
			want: `
foo: number
bar: string`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

// /////////////////////////////////////

func TestTypeObjectItem_Parsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *TypeObjectItem
	}{
		{
			name:    "Item",
			input:   "foo: number",
			wantErr: false,
			want: &TypeObjectItem{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
				Type: ParameterType{
					ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
					Ident: &Ident{
						ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
						Parts: []string{
							"number",
						},
					},
				},
			},
		},
		{
			name: "Single-line comment",
			input: `
// foo`[1:],
			wantErr: false,
			want: &TypeObjectItem{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Comment: &Comment{
					ASTNode:    ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					SingleLine: []string{"// foo"},
				},
			},
		},
		{
			name: "Multi-line comment",
			input: `
/* foo */`[1:],
			wantErr: false,
			want: &TypeObjectItem{
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

func TestTypeObjectItem_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		Input *TypeObjectItem
		want  *TypeObjectItem
	}{
		{
			name:  "Nil",
			Input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			Input: &TypeObjectItem{},
			want:  &TypeObjectItem{},
		},
		{
			name: "EmptyLine",
			Input: &TypeObjectItem{
				EmptyLine: "\n",
			},
			want: &TypeObjectItem{
				EmptyLine: "\n",
			},
		},
		{
			name: "ASTNode",
			Input: &TypeObjectItem{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
			want: &TypeObjectItem{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
		},
		{
			name: "Comments",
			Input: &TypeObjectItem{
				Comment: &Comment{Multiline: "foo"},
			},
			want: &TypeObjectItem{
				Comment: &Comment{Multiline: "foo"},
			},
		},
		{
			name: "Label",
			Input: &TypeObjectItem{
				Label: "foo",
			},
			want: &TypeObjectItem{
				Label: "foo",
			},
		},
		{
			name: "Type",
			Input: &TypeObjectItem{
				Type: ParameterType{Ident: &Ident{Parts: []string{"foo"}}},
			},
			want: &TypeObjectItem{
				Type: ParameterType{Ident: &Ident{Parts: []string{"foo"}}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*TypeObjectItem](t, tt.want, tt.Input.Clone())
		})
	}
}

func TestTypeObjectItem_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *TypeObjectItem
		want  []Node
	}{
		{
			name:  "Empty",
			input: &TypeObjectItem{},
			want: []Node{
				&ParameterType{},
			},
		},
		{
			name: "Comment",
			input: &TypeObjectItem{
				Comment: &Comment{Multiline: "foo"},
			},
			want: []Node{
				&Comment{Multiline: "foo"},
				&ParameterType{},
			},
		},
		{
			name: "Label",
			input: &TypeObjectItem{
				Label: "foo",
			},
			want: []Node{
				&ParameterType{},
			},
		},
		{
			name: "Type",
			input: &TypeObjectItem{
				Type: ParameterType{
					Ident: &Ident{
						Parts: []string{"foo"},
					},
				},
			},
			want: []Node{
				&ParameterType{
					Ident: &Ident{
						Parts: []string{"foo"},
					},
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

func TestTypeObjectItem_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     *TypeObjectItem
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
			input:     &TypeObjectItem{},
			wantPanic: true,
		},
		{
			name: "Value",
			input: &TypeObjectItem{
				Label: "foo",
				Type: ParameterType{
					Ident: &Ident{Parts: []string{"bar"}},
				},
			},
			want: "foo: bar",
		},
		{
			name: "Comment",
			input: &TypeObjectItem{
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
