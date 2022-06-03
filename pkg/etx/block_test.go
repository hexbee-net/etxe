package etx

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlock_Parsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Block
	}{
		{
			name: "no labels - no body",
			input: `
foo { }`[1:],
			wantErr: false,
			want: &Block{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Name:    "foo",
			},
		},
		{
			name: "one label (ident) - no body",
			input: `
foo bar { }`[1:],
			wantErr: false,
			want: &Block{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Name:    "foo",
				Labels:  []string{"bar"},
			},
		},
		{
			name: "two labels (ident) - no body",
			input: `
foo bar baz { }`[1:],
			wantErr: false,
			want: &Block{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Name:    "foo",
				Labels:  []string{"bar", "baz"},
			},
		},
		{
			name: "one label (string) - no body",
			input: `
foo 'bar' { }`[1:],
			wantErr: false,
			want: &Block{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Name:    "foo",
				Labels:  []string{"bar"},
			},
		},
		{
			name: "two labels (string) - no body",
			input: `
foo "bar" "baz"{ }`[1:],
			wantErr: false,
			want: &Block{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Name:    "foo",
				Labels:  []string{"bar", "baz"},
			},
		},
		{
			name: "no labels - body",
			input: `
foo {
	bar
}`[1:],
			wantErr: false,
			want: &Block{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Name:    "foo",
				Body: []*BlockItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 2, Column: 2}},
						Attribute: &Attribute{
							ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 2, Column: 2}},
							Key:     "bar",
						},
					},
				},
			},
		},
		{
			name: "no labels - body - close on same line",
			input: `
foo {
	bar }`[1:],
			wantErr: false,
			want: &Block{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Name:    "foo",
				Body: []*BlockItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 2, Column: 2}},
						Attribute: &Attribute{
							ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 2, Column: 2}},
							Key:     "bar",
						},
					},
				},
			},
		},
		{
			name: "one label - body",
			input: `
foo bar {
	baz
}`[1:],
			wantErr: false,
			want: &Block{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Name:    "foo",
				Labels:  []string{"bar"},
				Body: []*BlockItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 11, Line: 2, Column: 2}},
						Attribute: &Attribute{
							ASTNode: ASTNode{Pos: Position{Offset: 11, Line: 2, Column: 2}},
							Key:     "baz",
						},
					},
				},
			},
		},

		{
			name: "Single-line comment",
			input: `
// foo
bar { }`[1:],
			wantErr: false,
			want: &Block{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Comment: &Comment{
					ASTNode:    ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					SingleLine: []string{"// foo"},
				},
				Name: "bar",
			},
		},
		{
			name: "Single-line comment - separated",
			input: `
// foo

bar { }`[1:],
			wantErr: true,
		},
		{
			name: "Multi-line comment",
			input: `
/* foo */
bar { }`[1:],
			wantErr: false,
			want: &Block{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Comment: &Comment{
					ASTNode:   ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Multiline: "/* foo */\n",
				},
				Name: "bar",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, true)
		})
	}
}

func TestBlock_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *Block
		want  *Block
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &Block{},
			want:  &Block{},
		},
		{
			name: "ASTNode",
			input: &Block{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
			want: &Block{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
		},
		{
			name: "Comment",
			input: &Block{
				Comment: &Comment{Multiline: "foo"},
			},
			want: &Block{
				Comment: &Comment{Multiline: "foo"},
			},
		},
		{
			name: "Name",
			input: &Block{
				Name: "resource",
			},
			want: &Block{
				Name: "resource",
			},
		},
		{
			name: "Labels",
			input: &Block{
				Labels: []string{"foo"},
			},
			want: &Block{
				Labels: []string{"foo"},
			},
		},
		{
			name: "Body",
			input: &Block{
				Body: []*BlockItem{
					{Attribute: &Attribute{Key: "foo"}},
				},
			},
			want: &Block{
				Body: []*BlockItem{
					{Attribute: &Attribute{Key: "foo"}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*Block](t, tt.want, tt.input.Clone())
		})
	}
}

func TestBlock_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *Block
		want  []Node
	}{
		{
			name:  "Empty",
			input: &Block{},
			want:  nil,
		},
		{
			name: "Comment",
			input: &Block{
				Comment: &Comment{Multiline: "foo"},
			},
			want: []Node{
				&Comment{Multiline: "foo"},
			},
		},
		{
			name: "Name",
			input: &Block{
				Name: "resource",
			},
			want: nil,
		},
		{
			name: "Labels",
			input: &Block{
				Labels: []string{"foo"},
			},
			want: nil,
		},
		{
			name: "Body",
			input: &Block{
				Body: []*BlockItem{
					{Attribute: &Attribute{Key: "foo"}},
				},
			},
			want: []Node{
				&BlockItem{
					Attribute: &Attribute{Key: "foo"},
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

func TestBlock_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     *Block
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
			input: &Block{},
			want:  "",
		},
		{
			name: "Name - no labels - no body",
			input: &Block{
				Name: "resource",
			},
			want: "resource {}",
		},
		{
			name: "Name - labels - no body",
			input: &Block{
				Name:   "resource",
				Labels: []string{"foo"},
			},
			want: `resource "foo" {}`,
		},
		{
			name: "Name - no labels - body",
			input: &Block{
				Name: "resource",
				Body: []*BlockItem{
					{Attribute: &Attribute{Key: "foo"}},
				},
			},
			want: `
resource {
	foo
}`[1:],
		},
		{
			name: "Comment",
			input: &Block{
				Comment: &Comment{SingleLine: []string{"// foo"}},
				Name:    "resource",
			},
			want: `
// foo
resource {}`[1:],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

// /////////////////////////////////////

func TestBlockItem_Parsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *BlockItem
	}{
		{
			name:    "Block",
			input:   "foo {}",
			wantErr: false,
			want: &BlockItem{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Block: &Block{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Name:    "foo",
				},
			},
		},
		{
			name:    "Attribute",
			input:   "foo",
			wantErr: false,
			want: &BlockItem{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Attribute: &Attribute{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Key:     "foo",
				},
			},
		},
		{
			name:    "Attribute with scalar",
			input:   "foo = 1",
			wantErr: false,
			want: &BlockItem{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Attribute: &Attribute{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Key:     "foo",
					Value: testBuildExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
						Number:  &ValueNumber{big.NewFloat(1), "1"},
					}),
				},
			},
		},
		{
			name:    "Attribute with list",
			input:   "foo = [1, 2]",
			wantErr: false,
			want: &BlockItem{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Attribute: &Attribute{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Key:     "foo",
					Value: testBuildExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
						List: &ValueList{
							ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
							Items: []*Expr{
								testBuildExprTree[*Expr](t, &Value{
									ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
									Number:  &ValueNumber{big.NewFloat(1), "1"},
								}),
								testBuildExprTree[*Expr](t, &Value{
									ASTNode: ASTNode{Pos: Position{Offset: 10, Line: 1, Column: 11}},
									Number:  &ValueNumber{big.NewFloat(2), "2"},
								}),
							},
						},
					}),
				},
			},
		},
		{
			name:    "Attribute with map",
			input:   `foo = {"bar" = 1}`,
			wantErr: false,
			want: &BlockItem{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Attribute: &Attribute{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Key:     "foo",
					Value: testBuildExprTree[*Expr](t, &Value{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
						Map: &ValueMap{
							ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
							Entries: []*MapEntry{
								{
									ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
									Key: Value{
										ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
										Str: &ValueString{
											ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
											Fragment: []*StringFragment{
												{
													ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
													Text:    "bar",
												},
											},
										},
									},
									Value: *testBuildExprTree[*Expr](t, &Value{
										ASTNode: ASTNode{Pos: Position{Offset: 15, Line: 1, Column: 16}},
										Number:  &ValueNumber{big.NewFloat(1), "1"},
									}),
								},
							},
						},
					}),
				},
			},
		},

		{
			name:    "Comment",
			input:   "// foo",
			wantErr: false,
			want: &BlockItem{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Comment: &Comment{
					ASTNode:    ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					SingleLine: []string{"// foo"},
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

func TestBlockItem_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *BlockItem
		want  *BlockItem
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &BlockItem{},
			want:  &BlockItem{},
		},
		{
			name: "ASTNode",
			input: &BlockItem{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
			want: &BlockItem{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
		},
		{
			name: "Comment",
			input: &BlockItem{
				Comment: &Comment{Multiline: "foo"},
			},
			want: &BlockItem{
				Comment: &Comment{Multiline: "foo"},
			},
		},
		{
			name: "Attribute",
			input: &BlockItem{
				Attribute: &Attribute{Key: "foo"},
			},
			want: &BlockItem{
				Attribute: &Attribute{Key: "foo"},
			},
		},
		{
			name: "Sub-block",
			input: &BlockItem{
				Block: &Block{Name: "resource"},
			},
			want: &BlockItem{
				Block: &Block{Name: "resource"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*BlockItem](t, tt.want, tt.input.Clone())
		})
	}
}

func TestBlockItem_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *BlockItem
		want  []Node
	}{
		{
			name:  "Empty",
			input: &BlockItem{},
			want:  nil,
		},
		{
			name: "Comment",
			input: &BlockItem{
				Comment: &Comment{Multiline: "foo"},
			},
			want: []Node{
				&Comment{Multiline: "foo"},
			},
		},
		{
			name: "Attribute",
			input: &BlockItem{
				Attribute: &Attribute{Key: "foo"},
			},
			want: []Node{
				&Attribute{Key: "foo"},
			},
		},
		{
			name: "Sub-block",
			input: &BlockItem{
				Block: &Block{Name: "resource"},
			},
			want: []Node{
				&Block{Name: "resource"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestBlockItem_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     *BlockItem
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
			input: &BlockItem{},
			want:  "",
		},
		{
			name: "Attribute",
			input: &BlockItem{
				Attribute: &Attribute{Key: "foo"},
			},
			want: Attribute{Key: "foo"}.String(),
		},
		{
			name: "Sub-block",
			input: &BlockItem{
				Block: &Block{Name: "resource"},
			},
			want: Block{Name: "resource"}.String(),
		},
		{
			name: "Comment",
			input: &BlockItem{
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
