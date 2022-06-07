package etx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdent_Parsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Ident
	}{
		{
			name:    "one part",
			input:   "foo",
			wantErr: false,
			want: &Ident{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Parts: []string{
					"foo",
				},
			},
		},
		{
			name:    "several parts",
			input:   "foo.bar.baz",
			wantErr: false,
			want: &Ident{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Parts: []string{
					"foo",
					"bar",
					"baz",
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

func TestIdent_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		Input *Ident
		want  *Ident
	}{
		{
			name:  "Nil",
			Input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			Input: &Ident{},
			want:  &Ident{},
		},
		{
			name: "Parts",
			Input: &Ident{
				Parts: []string{"foo", "bar"},
			},
			want: &Ident{
				Parts: []string{"foo", "bar"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*Ident](t, tt.want, tt.Input.Clone())
		})
	}
}

func TestIdent_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *Ident
		want  []Node
	}{
		{
			name:  "Empty",
			input: &Ident{},
			want:  nil,
		},
		{
			name: "Parts",
			input: &Ident{
				Parts: []string{"foo", "bar"},
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

func TestIdent_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input Ident
		want  string
	}{
		{
			name: "One part",
			input: Ident{
				Parts: []string{
					"foo",
				},
			},
			want: "foo",
		},
		{
			name: "Several parts",
			input: Ident{
				Parts: []string{
					"foo",
					"bar",
					"baz",
				},
			},
			want: "foo.bar.baz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, false, tt.want, tt.input)
		})
	}
}
