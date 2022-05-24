package etx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString_Parsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *ValueString
	}{
		{
			name:    "Double quoted",
			input:   `"hello world"`,
			wantErr: false,
			want: &ValueString{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Fragment: []*StringFragment{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Text:    `hello world`,
					},
				},
			},
		},
		{
			name:    "Double quoted with single quote",
			input:   `"hello ' world"`,
			wantErr: false,
			want: &ValueString{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Fragment: []*StringFragment{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Text:    `hello `,
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
						Text:    `'`,
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
						Text:    ` world`,
					},
				},
			},
		},

		{
			name:    "Single quoted",
			input:   `'hello world'`,
			wantErr: false,
			want: &ValueString{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Fragment: []*StringFragment{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Text:    `hello world`,
					},
				},
			},
		},
		{
			name:    "Single quoted with double quote",
			input:   `'hello " world'`,
			wantErr: false,
			want: &ValueString{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Fragment: []*StringFragment{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Text:    `hello `,
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
						Text:    `"`,
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
						Text:    ` world`,
					},
				},
			},
		},

		{
			name:    "Escaped",
			input:   `"hello \t world"`,
			wantErr: false,
			want: &ValueString{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Fragment: []*StringFragment{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Text:    `hello `,
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
						Escaped: `\t`,
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 1, Column: 10}},
						Text:    ` world`,
					},
				},
			},
		},

		{
			name:    "Unicode - Short",
			input:   `"hello \u1234 world"`,
			wantErr: false,
			want: &ValueString{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Fragment: []*StringFragment{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Text:    `hello `,
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
						Unicode: `1234`,
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 1, Column: 14}},
						Text:    ` world`,
					},
				},
			},
		},
		{
			name:    "Unicode - Long",
			input:   `"hello \u12345678 world"`,
			wantErr: false,
			want: &ValueString{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Fragment: []*StringFragment{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Text:    `hello `,
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
						Unicode: `12345678`,
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 17, Line: 1, Column: 18}},
						Text:    ` world`,
					},
				},
			},
		},
		{
			name:    "Unicode - Invalid",
			input:   `"hello \u123 world"`,
			wantErr: true,
			want:    &ValueString{},
		},
		{
			name:    "Unicode - Trailing numbers",
			input:   `"hello \u123456 world"`,
			wantErr: false,
			want: &ValueString{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Fragment: []*StringFragment{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Text:    `hello `,
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
						Unicode: `1234`,
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 1, Column: 14}},
						Text:    `56 world`,
					},
				},
			},
		},

		{
			name:    "Expression - Only",
			input:   `"${foo}"`,
			wantErr: false,
			want: &ValueString{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Fragment: []*StringFragment{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Expr: testBuildExprTree[*Expr](t,
							&Value{
								ASTNode: ASTNode{Pos: Position{Offset: 3, Line: 1, Column: 4}},
								Ident:   &Ident{Parts: []string{"foo"}},
							}),
					},
				},
			},
		},
		{
			name:    "Expression - In Text",
			input:   `"hello ${foo} world"`,
			wantErr: false,
			want: &ValueString{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Fragment: []*StringFragment{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Text:    `hello `,
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
						Expr: testBuildExprTree[*Expr](t,
							&Value{
								ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 1, Column: 10}},
								Ident:   &Ident{Parts: []string{"foo"}},
							}),
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 1, Column: 14}},
						Text:    ` world`,
					},
				},
			},
		},
		{
			name:    "Non-Expression",
			input:   `"hello $${foo} world"`,
			wantErr: false,
			want: &ValueString{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Fragment: []*StringFragment{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Text:    `hello `,
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
						Text:    `$${`,
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 10, Line: 1, Column: 11}},
						Text:    `foo} world`,
					},
				},
			},
		},

		{
			name:    "Directive",
			input:   `"hello %{foo} world"`,
			wantErr: false,
			want: &ValueString{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Fragment: []*StringFragment{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Text:    `hello `,
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
						Directive: testBuildExprTree[*Expr](t,
							&Value{
								ASTNode: ASTNode{Pos: Position{Offset: 9, Line: 1, Column: 10}},
								Ident:   &Ident{Parts: []string{"foo"}},
							}),
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 1, Column: 14}},
						Text:    ` world`,
					},
				},
			},
		},
		{
			name:    "Non-Directive",
			input:   `"hello %%{foo} world"`,
			wantErr: false,
			want: &ValueString{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Fragment: []*StringFragment{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
						Text:    `hello `,
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 1, Column: 8}},
						Text:    `%%{`,
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 10, Line: 1, Column: 11}},
						Text:    `foo} world`,
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

func TestString_Clone(t *testing.T) {
	t.Parallel()

	type args struct {
	}
	tests := []struct {
		name  string
		input *ValueString
		want  *ValueString
	}{
		{
			name:  "Nil String",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty String",
			input: &ValueString{},
			want:  &ValueString{},
		},
		{
			name: "One Text Fragment",
			input: &ValueString{
				Fragment: []*StringFragment{
					{Text: `hello world`},
				},
			},
			want: &ValueString{
				Fragment: []*StringFragment{
					{Text: `hello world`},
				},
			},
		},
		{
			name: "Several Text Fragments",
			input: &ValueString{
				Fragment: []*StringFragment{
					{Text: `hello `},
					{Text: `world`},
				},
			},
			want: &ValueString{
				Fragment: []*StringFragment{
					{Text: `hello `},
					{Text: `world`},
				},
			},
		},
		{
			name: "Text and Escaped",
			input: &ValueString{
				Fragment: []*StringFragment{
					{Text: `hello `},
					{Escaped: `\t`},
					{Text: ` world`},
				},
			},
			want: &ValueString{
				Fragment: []*StringFragment{
					{Text: `hello `},
					{Escaped: `\t`},
					{Text: ` world`},
				},
			},
		},
		{
			name: "Text and Unicode",
			input: &ValueString{
				Fragment: []*StringFragment{
					{Text: `hello `},
					{Unicode: `1234`},
					{Text: ` world`},
				},
			},
			want: &ValueString{
				Fragment: []*StringFragment{
					{Text: `hello `},
					{Unicode: `1234`},
					{Text: ` world`},
				},
			},
		},
		{
			name: "Text and Expression",
			input: &ValueString{
				Fragment: []*StringFragment{
					{Text: `hello `},
					{Expr: testBuildExprTree[*Expr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}})},
					{Text: ` world`},
				},
			},
			want: &ValueString{
				Fragment: []*StringFragment{
					{Text: `hello `},
					{Expr: testBuildExprTree[*Expr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}})},
					{Text: ` world`},
				},
			},
		},
		{
			name: "Text and Expression",
			input: &ValueString{
				Fragment: []*StringFragment{
					{Text: `hello `},
					{Directive: testBuildExprTree[*Expr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}})},
					{Text: ` world`},
				},
			},
			want: &ValueString{
				Fragment: []*StringFragment{
					{Text: `hello `},
					{Directive: testBuildExprTree[*Expr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}})},
					{Text: ` world`},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*ValueString](t, tt.want, tt.input)
		})
	}
}

func TestString_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *ValueString
		want  []Node
	}{
		{
			name:  "Empty",
			input: &ValueString{},
			want:  nil,
		},
		{
			name: "Fragments",
			input: &ValueString{
				Fragment: []*StringFragment{
					{Text: "foo"},
				},
			},
			want: []Node{
				&StringFragment{
					Text: "foo",
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

func TestString_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     *ValueString
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
			input: &ValueString{},
			want:  `""`,
		},
		{
			name: "Empty String",
			input: &ValueString{
				Fragment: []*StringFragment{
					{Text: ``},
				},
			},
			want: `""`,
		},
		{
			name: "One Text Fragment",
			input: &ValueString{
				Fragment: []*StringFragment{
					{Text: `hello world`},
				},
			},
			want: `"hello world"`,
		},
		{
			name: "Several Text Fragments",
			input: &ValueString{
				Fragment: []*StringFragment{
					{Text: `hello `},
					{Text: `world`},
				},
			},
			want: `"hello world"`,
		},
		{
			name: "Text and Escaped",
			input: &ValueString{
				Fragment: []*StringFragment{
					{Text: `hello `},
					{Escaped: `\t`},
					{Text: ` world`},
				},
			},
			want: `"hello \t world"`,
		},
		{
			name: "Text and Unicode",
			input: &ValueString{
				Fragment: []*StringFragment{
					{Text: `hello `},
					{Unicode: `1234`},
					{Text: ` world`},
				},
			},
			want: `"hello \u1234 world"`,
		},
		{
			name: "Text and Expression",
			input: &ValueString{
				Fragment: []*StringFragment{
					{Text: `hello `},
					{Expr: testBuildExprTree[*Expr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}})},
					{Text: ` world`},
				},
			},
			want: `"hello ${foo} world"`,
		},
		{
			name: "Text and Directive",
			input: &ValueString{
				Fragment: []*StringFragment{
					{Text: `hello `},
					{Directive: testBuildExprTree[*Expr](t, &Value{Ident: &Ident{Parts: []string{"foo"}}})},
					{Text: ` world`},
				},
			},
			want: `"hello %{foo} world"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}
