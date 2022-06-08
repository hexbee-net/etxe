package etx

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap_Parsing(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *ValueMap
	}{
		{
			name:    "Empty",
			input:   `{ }`,
			wantErr: false,
			want: &ValueMap{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Items:   nil,
			},
		},
		{
			name:    "One entry - One line",
			input:   `{ a= 1 }`,
			wantErr: false,
			want: &ValueMap{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Items: []*MapItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Key: &MapKey{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Ident: &Ident{
								ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
								Parts:   []string{"a"},
							},
						},
						Value: BuildTestExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
							Number: &ValueNumber{
								ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
								Value:   big.NewFloat(1),
								Source:  `1`,
							},
						}),
					},
				},
			},
		},
		{
			name:    "One entry - One line - Trailing comma",
			input:   `{ a= 1, }`,
			wantErr: false,
			want: &ValueMap{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Items: []*MapItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Key: &MapKey{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Ident: &Ident{
								ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
								Parts:   []string{"a"},
							},
						},
						Value: BuildTestExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
							Number: &ValueNumber{
								ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
								Value:   big.NewFloat(1),
								Source:  `1`,
							},
						}),
					},
				},
			},
		},
		{
			name: "One entry - Linebreaks",
			input: `
{
    a= 1
}`[1:],
			wantErr: false,
			want: &ValueMap{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Items: []*MapItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
						Key: &MapKey{
							ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
							Ident: &Ident{
								ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
								Parts:   []string{"a"},
							},
						},
						Value: BuildTestExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 2, Column: 8}},
							Number: &ValueNumber{
								ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 2, Column: 8}},
								Value:   big.NewFloat(1),
								Source:  `1`,
							},
						}),
					},
				},
			},
		},
		{
			name:    "Two entries - One line",
			input:   `{ a= 1, b= 2 }`,
			wantErr: false,
			want: &ValueMap{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Items: []*MapItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Key: &MapKey{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Ident: &Ident{
								ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
								Parts:   []string{"a"},
							},
						},
						Value: BuildTestExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
							Number: &ValueNumber{
								ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
								Value:   big.NewFloat(1),
								Source:  `1`,
							},
						}),
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
						Key: &MapKey{
							ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
							Ident: &Ident{
								ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
								Parts:   []string{"b"},
							},
						},
						Value: BuildTestExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 11, Line: 1, Column: 12}},
							Number: &ValueNumber{
								ASTNode: ASTNode{Pos: Position{Offset: 11, Line: 1, Column: 12}},
								Value:   big.NewFloat(2),
								Source:  `2`,
							},
						}),
					},
				},
			},
		},
		{
			name: "Two entries - Linebreaks",
			input: `
{
    a= 1,
    b= 2
}`[1:],
			wantErr: false,
			want: &ValueMap{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Items: []*MapItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
						Key: &MapKey{
							ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
							Ident: &Ident{
								ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
								Parts:   []string{"a"},
							},
						},
						Value: BuildTestExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 2, Column: 8}},
							Number: &ValueNumber{
								ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 2, Column: 8}},
								Value:   big.NewFloat(1),
								Source:  `1`,
							},
						}),
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 16, Line: 3, Column: 5}},
						Key: &MapKey{
							ASTNode: ASTNode{Pos: Position{Offset: 16, Line: 3, Column: 5}},
							Ident: &Ident{
								ASTNode: ASTNode{Pos: Position{Offset: 16, Line: 3, Column: 5}},
								Parts:   []string{"b"},
							},
						},
						Value: BuildTestExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 19, Line: 3, Column: 8}},
							Number: &ValueNumber{
								ASTNode: ASTNode{Pos: Position{Offset: 19, Line: 3, Column: 8}},
								Value:   big.NewFloat(2),
								Source:  `2`,
							},
						}),
					},
				},
			},
		},
		{
			name: "Two entries - Linebreaks - Trailing comma",
			input: `
{
    a= 1,
    b= 2,
}`[1:],
			wantErr: false,
			want: &ValueMap{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Items: []*MapItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
						Key: &MapKey{
							ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
							Ident: &Ident{
								ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
								Parts:   []string{"a"},
							},
						},
						Value: BuildTestExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 2, Column: 8}},
							Number: &ValueNumber{
								ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 2, Column: 8}},
								Value:   big.NewFloat(1),
								Source:  `1`,
							},
						}),
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 16, Line: 3, Column: 5}},
						Key: &MapKey{
							ASTNode: ASTNode{Pos: Position{Offset: 16, Line: 3, Column: 5}},
							Ident: &Ident{
								ASTNode: ASTNode{Pos: Position{Offset: 16, Line: 3, Column: 5}},
								Parts:   []string{"b"},
							},
						},
						Value: BuildTestExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 19, Line: 3, Column: 8}},
							Number: &ValueNumber{
								ASTNode: ASTNode{Pos: Position{Offset: 19, Line: 3, Column: 8}},
								Value:   big.NewFloat(2),
								Source:  `2`,
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

func TestMap_Clone(t *testing.T) {
	tests := []struct {
		name  string
		input *ValueMap
		want  *ValueMap
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name: "Empty",
			input: &ValueMap{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
			want: &ValueMap{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
			},
		},
		{
			name: "Values",
			input: &ValueMap{
				Items: []*MapItem{
					{ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}}},
				},
			},
			want: &ValueMap{
				Items: []*MapItem{
					{ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*ValueMap](t, tt.want, tt.input)
		})
	}
}

func TestMap_Children(t *testing.T) {
	tests := []struct {
		name  string
		input *ValueMap
		want  []Node
	}{
		{
			name:  "Nil",
			input: &ValueMap{},
			want:  nil,
		},
		{
			name: "Empty",
			input: &ValueMap{
				Items: []*MapItem{},
			},
			want: nil,
		},
		{
			name: "Items",
			input: &ValueMap{
				Items: []*MapItem{
					{ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}}},
				},
			},
			want: []Node{
				&MapItem{ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestMap_FormattedString(t *testing.T) {
	tests := []struct {
		name      string
		input     *ValueMap
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
			input: &ValueMap{},
			want:  `{}`,
		},
		{
			name: "One entry",
			input: &ValueMap{
				Items: []*MapItem{
					{
						Key: &MapKey{
							Ident: &Ident{Parts: []string{"a"}},
						},
						Value: BuildTestExprTree[*Expr](t, &Value{
							Number: &ValueNumber{Value: big.NewFloat(1), Source: `1`},
						}),
					},
				},
			},
			want: `
{
	a: 1,
}`[1:],
		},
		{
			name: "Two entries",
			input: &ValueMap{
				Items: []*MapItem{
					{
						Key: &MapKey{
							Ident: &Ident{Parts: []string{"a"}},
						},
						Value: BuildTestExprTree[*Expr](t, &Value{
							Number: &ValueNumber{Value: big.NewFloat(1), Source: `1`},
						}),
					},
					{
						Key: &MapKey{
							Ident: &Ident{Parts: []string{"b"}},
						},
						Value: BuildTestExprTree[*Expr](t, &Value{
							Number: &ValueNumber{Value: big.NewFloat(1), Source: `2`},
						}),
					},
				},
			},
			want: `
{
	a: 1,
	b: 2,
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

func TestMapEntry_Clone(t *testing.T) {
	tests := []struct {
		name  string
		input *MapItem
		want  *MapItem
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name: "EmptyLine",
			input: &MapItem{
				EmptyLine: "\n",
			},
			want: &MapItem{
				EmptyLine: "\n",
			},
		},
		{
			name:  "Empty",
			input: &MapItem{},
			want:  &MapItem{},
		},
		{
			name: "Key",
			input: &MapItem{
				Key: &MapKey{
					ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
				},
			},
			want: &MapItem{
				Key: &MapKey{
					ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
				},
			},
		},
		{
			name: "Value",
			input: &MapItem{
				Value: &Expr{
					ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
				},
			},
			want: &MapItem{
				Value: &Expr{
					ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
				},
			},
		},
		{
			name: "Comments",
			input: &MapItem{
				Comment: &Comment{Multiline: "foo"},
			},
			want: &MapItem{
				Comment: &Comment{Multiline: "foo"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*MapItem](t, tt.want, tt.input)
		})
	}
}

func TestMapEntry_Children(t *testing.T) {
	tests := []struct {
		name  string
		input MapItem
		want  []Node
	}{
		{
			name:  "Empty",
			input: MapItem{},
			want:  nil,
		},
		{
			name: "EmptyLine",
			input: MapItem{
				EmptyLine: "\n",
			},
			want: nil,
		},
		{
			name: "Key",
			input: MapItem{
				Key: &MapKey{Ident: &Ident{Parts: []string{"a"}}},
			},
			want: []Node{
				&MapKey{Ident: &Ident{Parts: []string{"a"}}},
			},
		},
		{
			name: "Value",
			input: MapItem{
				Value: BuildTestExprTree[*Expr](t, &Ident{Parts: []string{"a"}}),
			},
			want: []Node{
				BuildTestExprTree[*Expr](t, &Ident{Parts: []string{"a"}}),
			},
		},
		{
			name: "Comment",
			input: MapItem{
				Comment: &Comment{Multiline: "foo"},
			},
			want: []Node{
				&Comment{Multiline: "foo"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestMapEntry_FormattedString(t *testing.T) {
	tests := []struct {
		name      string
		input     *MapItem
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
			input:     &MapItem{},
			wantPanic: true,
		},
		{
			name: "Values",
			input: &MapItem{
				Key: &MapKey{
					Ident: &Ident{
						Parts: []string{"a"},
					},
				},
				Value: BuildTestExprTree[*Expr](t, &Value{
					Number: &ValueNumber{Value: big.NewFloat(1), Source: `1`},
				}),
			},
			want: `a = 1`,
		},
		{
			name: "Comment",
			input: &MapItem{
				Comment: &Comment{SingleLine: []string{"// foo"}},
			},
			want: "// foo\n",
		},
		{
			name: "EmptyLine",
			input: &MapItem{
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

func TestMapKey_Clone(t *testing.T) {
	tests := []struct {
		name  string
		input *MapKey
		want  *MapKey
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &MapKey{},
			want:  &MapKey{},
		},
		{
			name: "ASTNode",
			input: &MapKey{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
			want: &MapKey{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
		},
		{
			name: "Ident",
			input: &MapKey{
				Ident: &Ident{
					ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
				},
			},
			want: &MapKey{
				Ident: &Ident{
					ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
				},
			},
		},
		{
			name: "Str",
			input: &MapKey{
				Str: &ValueString{
					ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
				},
			},
			want: &MapKey{
				Str: &ValueString{
					ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*MapKey](t, tt.want, tt.input)
		})
	}
}

func TestMapKey_Children(t *testing.T) {
	tests := []struct {
		name  string
		input MapKey
		want  []Node
	}{
		{
			name:  "Empty",
			input: MapKey{},
			want:  nil,
		},
		{
			name: "Ident",
			input: MapKey{
				Ident: &Ident{
					ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
				},
			},
			want: []Node{
				&Ident{
					ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
				},
			},
		},
		{
			name: "Str",
			input: MapKey{
				Str: &ValueString{
					ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
				},
			},
			want: []Node{
				&ValueString{
					ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
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

func TestMapKey_FormattedString(t *testing.T) {
	tests := []struct {
		name      string
		input     *MapKey
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
			input:     &MapKey{},
			wantPanic: true,
		},
		{
			name: "Ident",
			input: &MapKey{
				Ident: &Ident{
					Parts: []string{"a"},
				},
			},
			want: `a`,
		},
		{
			name: "Str",
			input: &MapKey{
				Str: &ValueString{
					Fragment: []*StringFragment{
						{Text: "a"},
					},
				},
			},
			want: `"a"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}
