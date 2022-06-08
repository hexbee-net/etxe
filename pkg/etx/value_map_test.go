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
				Entries: nil,
			},
		},
		{
			name:    "One entry - One line",
			input:   `{ a= 1 }`,
			wantErr: false,
			want: &ValueMap{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Entries: []*MapEntry{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Key: MapKey{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Ident: &Ident{
								ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
								Parts:   []string{"a"},
							},
						},
						Value: *BuildTestExprTree[*Expr](t, &Value{
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
				Entries: []*MapEntry{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Key: MapKey{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Ident: &Ident{
								ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
								Parts:   []string{"a"},
							},
						},
						Value: *BuildTestExprTree[*Expr](t, &Value{
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
				Entries: []*MapEntry{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
						Key: MapKey{
							ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
							Ident: &Ident{
								ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
								Parts:   []string{"a"},
							},
						},
						Value: *BuildTestExprTree[*Expr](t, &Value{
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
				Entries: []*MapEntry{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Key: MapKey{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Ident: &Ident{
								ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
								Parts:   []string{"a"},
							},
						},
						Value: *BuildTestExprTree[*Expr](t, &Value{
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
						Key: MapKey{
							ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
							Ident: &Ident{
								ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
								Parts:   []string{"b"},
							},
						},
						Value: *BuildTestExprTree[*Expr](t, &Value{
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
				Entries: []*MapEntry{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
						Key: MapKey{
							ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
							Ident: &Ident{
								ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
								Parts:   []string{"a"},
							},
						},
						Value: *BuildTestExprTree[*Expr](t, &Value{
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
						Key: MapKey{
							ASTNode: ASTNode{Pos: Position{Offset: 16, Line: 3, Column: 5}},
							Ident: &Ident{
								ASTNode: ASTNode{Pos: Position{Offset: 16, Line: 3, Column: 5}},
								Parts:   []string{"b"},
							},
						},
						Value: *BuildTestExprTree[*Expr](t, &Value{
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
				Entries: []*MapEntry{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
						Key: MapKey{
							ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
							Ident: &Ident{
								ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
								Parts:   []string{"a"},
							},
						},
						Value: *BuildTestExprTree[*Expr](t, &Value{
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
						Key: MapKey{
							ASTNode: ASTNode{Pos: Position{Offset: 16, Line: 3, Column: 5}},
							Ident: &Ident{
								ASTNode: ASTNode{Pos: Position{Offset: 16, Line: 3, Column: 5}},
								Parts:   []string{"b"},
							},
						},
						Value: *BuildTestExprTree[*Expr](t, &Value{
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
				Entries: []*MapEntry{
					{ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}}},
				},
			},
			want: &ValueMap{
				Entries: []*MapEntry{
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
				Entries: []*MapEntry{},
			},
			want: nil,
		},
		{
			name: "Entries",
			input: &ValueMap{
				Entries: []*MapEntry{
					{ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}}},
				},
			},
			want: []Node{
				&MapEntry{ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}}},
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
				Entries: []*MapEntry{
					{
						Key: MapKey{
							Ident: &Ident{Parts: []string{"a"}},
						},
						Value: *BuildTestExprTree[*Expr](t, &Value{
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
				Entries: []*MapEntry{
					{
						Key: MapKey{
							Ident: &Ident{Parts: []string{"a"}},
						},
						Value: *BuildTestExprTree[*Expr](t, &Value{
							Number: &ValueNumber{Value: big.NewFloat(1), Source: `1`},
						}),
					},
					{
						Key: MapKey{
							Ident: &Ident{Parts: []string{"b"}},
						},
						Value: *BuildTestExprTree[*Expr](t, &Value{
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
		input *MapEntry
		want  *MapEntry
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &MapEntry{},
			want:  &MapEntry{},
		},
		{
			name: "Comments",
			input: &MapEntry{
				Comment: &Comment{Multiline: "foo"},
			},
			want: &MapEntry{
				Comment: &Comment{Multiline: "foo"},
			},
		},
		{
			name: "Key",
			input: &MapEntry{
				Key: MapKey{
					ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
				},
			},
			want: &MapEntry{
				Key: MapKey{
					ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
				},
			},
		},
		{
			name: "Value",
			input: &MapEntry{
				Value: Expr{
					ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
				},
			},
			want: &MapEntry{
				Value: Expr{
					ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*MapEntry](t, tt.want, tt.input)
		})
	}
}

func TestMapEntry_Children(t *testing.T) {
	tests := []struct {
		name  string
		input MapEntry
		want  []Node
	}{
		{
			name:  "Empty",
			input: MapEntry{},
			want: []Node{
				&MapKey{},
				&Expr{},
			},
		},
		{
			name: "Comment",
			input: MapEntry{
				Comment: &Comment{Multiline: "foo"},
			},
			want: []Node{
				&Comment{Multiline: "foo"},
				&MapKey{},
				&Expr{},
			},
		},
		{
			name: "Key",
			input: MapEntry{
				Key: MapKey{Ident: &Ident{Parts: []string{"a"}}},
			},
			want: []Node{
				&MapKey{Ident: &Ident{Parts: []string{"a"}}},
				&Expr{},
			},
		},
		{
			name: "Value",
			input: MapEntry{
				Value: *BuildTestExprTree[*Expr](t, &Ident{Parts: []string{"a"}}),
			},
			want: []Node{
				&MapKey{},
				BuildTestExprTree[*Expr](t, &Ident{Parts: []string{"a"}}),
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
		input     *MapEntry
		wantPanic bool
		want      string
	}{
		{
			name:      "Nil",
			input:     nil,
			wantPanic: true,
		},
		{
			name: "Values",
			input: &MapEntry{
				Key: MapKey{
					Ident: &Ident{
						Parts: []string{"a"},
					},
				},
				Value: *BuildTestExprTree[*Expr](t, &Value{
					Number: &ValueNumber{Value: big.NewFloat(1), Source: `1`},
				}),
			},
			want: `a = 1`,
		},
		{
			name: "Comment",
			input: &MapEntry{
				Comment: &Comment{SingleLine: []string{"// foo"}},
				Key: MapKey{
					Ident: &Ident{
						Parts: []string{"a"},
					},
				},
				Value: *BuildTestExprTree[*Expr](t, &Value{
					Number: &ValueNumber{Value: big.NewFloat(1), Source: `1`},
				}),
			},
			want: `
// foo
a = 1`[1:],
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
