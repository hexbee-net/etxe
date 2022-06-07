package etx

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttribute_Parsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Attribute
	}{
		{
			name:    "Key only",
			input:   "foo",
			wantErr: false,
			want: &Attribute{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Key:     "foo",
			},
		},
		{
			name:    "Key and value",
			input:   "foo = 1",
			wantErr: false,
			want: &Attribute{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Key:     "foo",
				Value: testBuildExprTree[*Expr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
					Number:  &ValueNumber{ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}}, big.NewFloat(1), "1"},
				}),
			},
		},
		{
			name: "Single-line comment",
			input: `
// foo
bar`[1:],
			wantErr: false,
			want: &Attribute{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Comment: &Comment{
					ASTNode:    ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					SingleLine: []string{"// foo"},
				},
				Key: "bar",
			},
		},
		{
			name: "Single-line comment - separated",
			input: `
// foo

bar`[1:],
			wantErr: true,
		},
		{
			name: "Multi-line comment",
			input: `
/* foo */
bar`[1:],
			wantErr: false,
			want: &Attribute{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Comment: &Comment{
					ASTNode:   ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Multiline: "/* foo */\n",
				},
				Key: "bar",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, true)
		})
	}
}

func TestAttribute_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *Attribute
		want  *Attribute
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &Attribute{},
			want:  &Attribute{},
		},
		{
			name: "ASTNode",
			input: &Attribute{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
			want: &Attribute{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
		},
		{
			name: "Comment",
			input: &Attribute{
				Comment: &Comment{Multiline: "foo"},
			},
			want: &Attribute{
				Comment: &Comment{Multiline: "foo"},
			},
		},
		{
			name: "Key",
			input: &Attribute{
				Key: "foo",
			},
			want: &Attribute{
				Key: "foo",
			},
		},
		{
			name: "Value",
			input: &Attribute{
				Value: testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
			},
			want: &Attribute{
				Value: testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*Attribute](t, tt.want, tt.input.Clone())
		})
	}
}

func TestAttribute_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *Attribute
		want  []Node
	}{
		{
			name:  "Empty",
			input: &Attribute{},
			want:  nil,
		},
		{
			name: "Comment",
			input: &Attribute{
				Comment: &Comment{Multiline: "foo"},
			},
			want: []Node{
				&Comment{Multiline: "foo"},
			},
		},
		{
			name: "Key",
			input: &Attribute{
				Key: "foo",
			},
			want: nil,
		},
		{
			name: "Value",
			input: &Attribute{
				Value: testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
			},
			want: []Node{
				testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestAttribute_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     *Attribute
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
			input: &Attribute{},
			want:  "",
		},
		{
			name: "Key only",
			input: &Attribute{
				Key: "foo",
			},
			want: "foo",
		},
		{
			name: "Key and value",
			input: &Attribute{
				Key:   "foo",
				Value: testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
			},
			want: "foo: 1",
		},
		{
			name: "Comment",
			input: &Attribute{
				Comment: &Comment{SingleLine: []string{"// foo"}},
				Key:     "bar",
			},
			want: `
// foo
bar`[1:],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}
