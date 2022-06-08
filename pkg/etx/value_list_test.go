package etx

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList_Parsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *ValueList
	}{
		{
			name:    "Empty",
			input:   `[ ]`,
			wantErr: false,
			want: &ValueList{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Items:   nil,
			},
		},
		{
			name:    "One item - One line",
			input:   `[ a ]`,
			wantErr: false,
			want: &ValueList{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Items: []*ListItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Value: BuildTestExprTree[*Expr](t, &Ident{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Parts:   []string{"a"},
						}),
					},
				},
			},
		},
		{
			name:    "One item - One line - Trailing comma",
			input:   `[ a, ]`,
			wantErr: false,
			want: &ValueList{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Items: []*ListItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Value: BuildTestExprTree[*Expr](t, &Ident{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Parts:   []string{"a"},
						}),
					},
				},
			},
		},
		{
			name: "One item - Linebreaks",
			input: `
[
    a
]`[1:],
			wantErr: false,
			want: &ValueList{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Items: []*ListItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
						Value: BuildTestExprTree[*Expr](t, &Ident{
							ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
							Parts:   []string{"a"},
						}),
					},
				},
			},
		},
		{
			name:    "Two items - One line",
			input:   `[ a, b ]`,
			wantErr: false,
			want: &ValueList{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Items: []*ListItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Value: BuildTestExprTree[*Expr](t, &Ident{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Parts:   []string{"a"},
						}),
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
						Value: BuildTestExprTree[*Expr](t, &Ident{
							ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
							Parts:   []string{"b"},
						}),
					},
				},
			},
		},
		{
			name: "Two items - Linebreaks",
			input: `
[
    a,
    b
]`[1:],
			wantErr: false,
			want: &ValueList{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Items: []*ListItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
						Value: BuildTestExprTree[*Expr](t, &Ident{
							ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
							Parts:   []string{"a"},
						}),
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 3, Column: 5}},
						Value: BuildTestExprTree[*Expr](t, &Ident{
							ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 3, Column: 5}},
							Parts:   []string{"b"},
						}),
					},
				},
			},
		},
		{
			name: "Two entries - Linebreaks - Trailing comma",
			input: `
[
    a,
    b,
]`[1:],
			wantErr: false,
			want: &ValueList{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Items: []*ListItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
						Value: BuildTestExprTree[*Expr](t, &Ident{
							ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
							Parts:   []string{"a"},
						}),
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 3, Column: 5}},
						Value: BuildTestExprTree[*Expr](t, &Ident{
							ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 3, Column: 5}},
							Parts:   []string{"b"},
						}),
					},
				},
			},
		},
		{
			name:    "One expression - One line",
			input:   `[ 1 + 2 ]`,
			wantErr: false,
			want: &ValueList{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Items: []*ListItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Value: BuildTestExprTree[*Expr](t,
							&ExprAdditive{
								ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
								Left: *BuildTestExprTree[*ExprMultiplicative](t, &Value{
									ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
									Number: &ValueNumber{
										ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
										Value:   big.NewFloat(1),
										Source:  "1",
									},
								}),
								Op: OpPlus,
								Right: BuildTestExprTree[*ExprAdditive](t, &Value{
									ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
									Number: &ValueNumber{
										ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
										Value:   big.NewFloat(2),
										Source:  "2",
									},
								}),
							},
						),
					},
				},
			},
		},
		{
			name:    "Two expressions - One line",
			input:   `[ 1 + 2, 3 - 4 ]`,
			wantErr: false,
			want: &ValueList{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Items: []*ListItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Value: BuildTestExprTree[*Expr](t,
							&ExprAdditive{
								ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
								Left: *BuildTestExprTree[*ExprMultiplicative](t, &Value{
									ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
									Number: &ValueNumber{
										ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
										Value:   big.NewFloat(1),
										Source:  "1",
									},
								}),
								Op: OpPlus,
								Right: BuildTestExprTree[*ExprAdditive](t, &Value{
									ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
									Number: &ValueNumber{
										ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
										Value:   big.NewFloat(2),
										Source:  "2",
									},
								}),
							},
						)},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 1, Column: 10}},
						Value: BuildTestExprTree[*Expr](t,
							&ExprAdditive{
								ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 1, Column: 10}},
								Left: *BuildTestExprTree[*ExprMultiplicative](t, &Value{
									ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 1, Column: 10}},
									Number: &ValueNumber{
										ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 1, Column: 10}},
										Value:   big.NewFloat(3),
										Source:  "3",
									},
								}),
								Op: OpMinus,
								Right: BuildTestExprTree[*ExprAdditive](t, &Value{
									ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 1, Column: 14}},
									Number: &ValueNumber{
										ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 1, Column: 14}},
										Value:   big.NewFloat(4),
										Source:  "4",
									},
								}),
							},
						)},
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

func TestList_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ValueList
		want  *ValueList
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &ValueList{},
			want:  &ValueList{},
		},
		{
			name: "ASTNode",
			input: &ValueList{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
			want: &ValueList{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
		},
		{
			name: "Values",
			input: &ValueList{
				Items: []*ListItem{
					{ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}}},
				},
			},
			want: &ValueList{
				Items: []*ListItem{
					{ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*ValueList](t, tt.want, tt.input)
		})
	}
}

func TestList_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ValueList
		want  []Node
	}{
		{
			name:  "Nil",
			input: &ValueList{},
			want:  nil,
		},
		{
			name: "Empty",
			input: &ValueList{
				Items: []*ListItem{},
			},
			want: nil,
		},
		{
			name: "Items",
			input: &ValueList{
				Items: []*ListItem{
					{ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}}},
				},
			},
			want: []Node{
				&ListItem{ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestList_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     *ValueList
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
			input: &ValueList{},
			want:  `[]`,
		},
		{
			name: "One item",
			input: &ValueList{
				Items: []*ListItem{
					{Value: BuildTestExprTree[*Expr](t, &Ident{Parts: []string{"a"}})},
				},
			},
			want: `
[
	a,
]`[1:],
		},
		{
			name: "Two items",
			input: &ValueList{
				Items: []*ListItem{
					{Value: BuildTestExprTree[*Expr](t, &Ident{Parts: []string{"a"}})},
					{Value: BuildTestExprTree[*Expr](t, &Ident{Parts: []string{"b"}})},
				},
			},
			want: `
[
	a,
	b,
]`[1:],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

// /////////////////////////////////////

func TestListItem_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ListItem
		want  *ListItem
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &ListItem{},
			want:  &ListItem{},
		},
		{
			name: "ASTNode",
			input: &ListItem{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
			want: &ListItem{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
		},
		{
			name: "EmptyLine",
			input: &ListItem{
				EmptyLine: "\n",
			},
			want: &ListItem{
				EmptyLine: "\n",
			},
		},
		{
			name: "Value",
			input: &ListItem{
				Value: &Expr{ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}}},
			},
			want: &ListItem{
				Value: &Expr{ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}}},
			},
		},
		{
			name: "Comment",
			input: &ListItem{
				Comment: &Comment{ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}}},
			},
			want: &ListItem{
				Comment: &Comment{ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*ListItem](t, tt.want, tt.input)
		})
	}
}

func TestListItem_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ListItem
		want  []Node
	}{
		{
			name:  "Nil",
			input: &ListItem{},
			want:  nil,
		},
		{
			name:  "Empty",
			input: &ListItem{},
			want:  nil,
		},
		{
			name: "EmptyLine",
			input: &ListItem{
				EmptyLine: "\n",
			},
			want: nil,
		},
		{
			name: "Value",
			input: &ListItem{
				Value: &Expr{ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}}},
			},
			want: []Node{
				&Expr{ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}}},
			},
		},
		{
			name: "Comment",
			input: &ListItem{
				Comment: &Comment{ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}}},
			},
			want: []Node{
				&Comment{ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestListItem_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     *ListItem
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
			input:     &ListItem{},
			wantPanic: true,
		},
		{
			name: "EmptyLine",
			input: &ListItem{
				EmptyLine: "\n",
			},
			want: "\n",
		},
		{
			name: "Value",
			input: &ListItem{
				Value: BuildTestExprTree[*Expr](t, &Ident{Parts: []string{"a"}}),
			},
			want: "a,\n",
		},
		{
			name: "Comment",
			input: &ListItem{
				Comment: &Comment{
					Multiline: "/* foo */",
				},
			},
			want: "/* foo */",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}
