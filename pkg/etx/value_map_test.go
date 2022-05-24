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
			input:   `{ a: 1 }`,
			wantErr: false,
			want: &ValueMap{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Entries: []*MapEntry{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Key: Value{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Ident: &Ident{
								Parts: []string{"a"},
							},
						},
						Value: Value{
							ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
							Number:  &ValueNumber{big.NewFloat(1), `1`},
						},
					},
				},
			},
		},
		{
			name:    "One entry - One line - Trailing comma",
			input:   `{ a: 1, }`,
			wantErr: false,
			want: &ValueMap{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Entries: []*MapEntry{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Key: Value{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Ident: &Ident{
								Parts: []string{"a"},
							},
						},
						Value: Value{
							ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
							Number:  &ValueNumber{big.NewFloat(1), `1`},
						},
					},
				},
			},
		},
		{
			name: "One entry - Linebreaks",
			input: `
{
    a: 1
}`[1:],
			wantErr: false,
			want: &ValueMap{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Entries: []*MapEntry{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
						Key: Value{
							ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
							Ident: &Ident{
								Parts: []string{"a"},
							},
						},
						Value: Value{
							ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 2, Column: 8}},
							Number:  &ValueNumber{big.NewFloat(1), `1`},
						},
					},
				},
			},
		},
		{
			name:    "Two entries - One line",
			input:   `{ a: 1, b: 2 }`,
			wantErr: false,
			want: &ValueMap{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Entries: []*MapEntry{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
						Key: Value{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Ident: &Ident{
								Parts: []string{"a"},
							},
						},
						Value: Value{
							ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
							Number:  &ValueNumber{big.NewFloat(1), `1`},
						},
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
						Key: Value{
							ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
							Ident: &Ident{
								Parts: []string{"b"},
							},
						},
						Value: Value{
							ASTNode: ASTNode{Pos: Position{Offset: 11, Line: 1, Column: 12}},
							Number:  &ValueNumber{big.NewFloat(2), `2`},
						},
					},
				},
			},
		},
		{
			name: "Two entries - Linebreaks",
			input: `
{
    a: 1,
    b: 2
}`[1:],
			wantErr: false,
			want: &ValueMap{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Entries: []*MapEntry{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
						Key: Value{
							ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
							Ident: &Ident{
								Parts: []string{"a"},
							},
						},
						Value: Value{
							ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 2, Column: 8}},
							Number:  &ValueNumber{big.NewFloat(1), `1`},
						},
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 16, Line: 3, Column: 5}},
						Key: Value{
							ASTNode: ASTNode{Pos: Position{Offset: 16, Line: 3, Column: 5}},
							Ident: &Ident{
								Parts: []string{"b"},
							},
						},
						Value: Value{
							ASTNode: ASTNode{Pos: Position{Offset: 19, Line: 3, Column: 8}},
							Number:  &ValueNumber{big.NewFloat(2), `2`},
						},
					},
				},
			},
		},
		{
			name: "Two entries - Linebreaks - Trailing comma",
			input: `
{
    a: 1,
    b: 2,
}`[1:],
			wantErr: false,
			want: &ValueMap{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Entries: []*MapEntry{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
						Key: Value{
							ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 5}},
							Ident: &Ident{
								Parts: []string{"a"},
							},
						},
						Value: Value{
							ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 2, Column: 8}},
							Number:  &ValueNumber{big.NewFloat(1), `1`},
						},
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 16, Line: 3, Column: 5}},
						Key: Value{
							ASTNode: ASTNode{Pos: Position{Offset: 16, Line: 3, Column: 5}},
							Ident: &Ident{
								Parts: []string{"b"},
							},
						},
						Value: Value{
							ASTNode: ASTNode{Pos: Position{Offset: 19, Line: 3, Column: 8}},
							Number:  &ValueNumber{big.NewFloat(2), `2`},
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
			name:  "Empty",
			input: &ValueMap{},
			want:  &ValueMap{},
		},
		{
			name: "Values",
			input: &ValueMap{
				Entries: []*MapEntry{},
			},
			want: &ValueMap{
				Entries: []*MapEntry{},
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
					{},
				},
			},
			want: []Node{
				&MapEntry{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestMap_String(t *testing.T) {
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
						Key: Value{
							Ident: &Ident{
								Parts: []string{"a"},
							},
						},
						Value: Value{
							Number: &ValueNumber{big.NewFloat(1), `1`},
						},
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
						Key: Value{
							Ident: &Ident{
								Parts: []string{"a"},
							},
						},
						Value: Value{
							Number: &ValueNumber{big.NewFloat(1), `1`},
						},
					},
					{
						Key: Value{
							Ident: &Ident{
								Parts: []string{"b"},
							},
						},
						Value: Value{
							Number: &ValueNumber{big.NewFloat(1), `2`},
						},
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
			name: "Values",
			input: &MapEntry{
				Key:   Value{},
				Value: Value{},
			},
			want: &MapEntry{
				Key:   Value{},
				Value: Value{},
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
				&Value{},
				&Value{},
			},
		},
		{
			name: "Key",
			input: MapEntry{
				Key: Value{
					Ident: &Ident{Parts: []string{"a"}},
				},
			},
			want: []Node{
				&Value{
					Ident: &Ident{Parts: []string{"a"}},
				},
				&Value{},
			},
		},
		{
			name: "Value",
			input: MapEntry{
				Value: Value{
					Ident: &Ident{Parts: []string{"a"}},
				},
			},
			want: []Node{
				&Value{},
				&Value{
					Ident: &Ident{Parts: []string{"a"}},
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
