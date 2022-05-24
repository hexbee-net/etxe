package etx

import (
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
				Items: []*Value{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Ident: &Ident{
							Parts: []string{"a"},
						},
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
				Items: []*Value{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Ident: &Ident{
							Parts: []string{"a"},
						},
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
				Items: []*Value{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
						Ident: &Ident{
							Parts: []string{"a"},
						},
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
				Items: []*Value{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Ident: &Ident{
							Parts: []string{"a"},
						},
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
						Ident: &Ident{
							Parts: []string{"b"},
						},
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
				Items: []*Value{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
						Ident: &Ident{
							Parts: []string{"a"},
						},
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 3, Column: 5}},
						Ident: &Ident{
							Parts: []string{"b"},
						},
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
				Items: []*Value{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
						Ident: &Ident{
							Parts: []string{"a"},
						},
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 3, Column: 5}},
						Ident: &Ident{
							Parts: []string{"b"},
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
			name: "Values",
			input: &ValueList{
				Items: []*Value{},
			},
			want: &ValueList{
				Items: []*Value{},
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
				Items: []*Value{},
			},
			want: nil,
		},
		{
			name: "Items",
			input: &ValueList{
				Items: []*Value{
					{},
				},
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
				Items: []*Value{
					{
						Ident: &Ident{
							Parts: []string{"a"},
						},
					},
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
				Items: []*Value{
					{
						Ident: &Ident{
							Parts: []string{"a"},
						},
					},
					{
						Ident: &Ident{
							Parts: []string{"b"},
						},
					},
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
