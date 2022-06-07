package etx

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComment_Parsing(t *testing.T) {
	t.Parallel()

	type TestStruct struct {
		Comments []*Comment `parser:"(@@)*" json:"comments,omitempty"`
		Value    *Expr      `parser:"[ @@ ]"`
	}
	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *TestStruct
	}{
		{
			name: "One double-slash",
			input: `
// foo`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []*Comment{
					{
						ASTNode:    ASTNode{},
						SingleLine: []string{"// foo"},
					},
				},
			},
		},
		{
			name: "One double-slash - new line",
			input: `
// foo
`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []*Comment{
					{
						ASTNode:    ASTNode{},
						SingleLine: []string{"// foo"},
					},
				},
			},
		},
		{
			name: "One hashtag",
			input: `
# foo`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []*Comment{
					{
						ASTNode:    ASTNode{},
						SingleLine: []string{"# foo"},
					},
				},
			},
		},

		{
			name: "Two double-slash",
			input: `
// foo
// bar
		`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []*Comment{
					{
						ASTNode: ASTNode{},
						SingleLine: []string{
							"// foo",
							"// bar",
						},
					},
				},
			},
		},
		{
			name: "Two hashtags",
			input: `
# foo
# bar
		`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []*Comment{
					{
						ASTNode: ASTNode{},
						SingleLine: []string{
							"# foo",
							"# bar",
						},
					},
				},
			},
		},
		{
			name: "Double-slash hashtags mix",
			input: `
// foo
# bar
		`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []*Comment{
					{
						ASTNode: ASTNode{},
						SingleLine: []string{
							"// foo",
							"# bar",
						},
					},
				},
			},
		},
		{
			name: "Two separated double-slash",
			input: `
// foo

// bar
		`[1:],
			wantErr: true,
		},

		{
			name: "One double-slash before value",
			input: `
// foo
1`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []*Comment{
					{
						ASTNode:    ASTNode{},
						SingleLine: []string{"// foo"},
					},
				},
				Value: BuildTestExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
			},
		},
		{
			name: "Two double-slash before value",
			input: `
// foo
// bar
1`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []*Comment{
					{
						ASTNode:    ASTNode{},
						SingleLine: []string{"// foo", "// bar"},
					},
				},
				Value: BuildTestExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
			},
		},

		{
			name: "One multiline - one line",
			input: `
/* foo */`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []*Comment{
					{
						ASTNode:   ASTNode{},
						Multiline: "/* foo */",
					},
				},
			},
		},
		{
			name: "One multiline - two lines",
			input: `
/* foo
   bar */`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []*Comment{
					{
						ASTNode:   ASTNode{},
						Multiline: "/* foo\n   bar */",
					},
				},
			},
		},
		{
			name: "Two multilines - one line",
			input: `
/* foo */ /* bar */`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []*Comment{
					{
						ASTNode:   ASTNode{},
						Multiline: "/* foo */",
					},
					{
						ASTNode:   ASTNode{},
						Multiline: "/* bar */",
					},
				},
			},
		},
		{
			name: "Two multilines - two line",
			input: `
/* foo */
/* bar */
`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []*Comment{
					{
						ASTNode:   ASTNode{},
						Multiline: "/* foo */\n",
					},
					{
						ASTNode:   ASTNode{},
						Multiline: "/* bar */\n",
					},
				},
			},
		},

		{
			name: "One multiline - one line - with expr on same line",
			input: `
/* foo */ 1`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []*Comment{
					{
						ASTNode:   ASTNode{},
						Multiline: "/* foo */",
					},
				},
				Value: BuildTestExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
			},
		},
		{
			name: "One multiline - two lines - with expr on same line",
			input: `
/* foo
   bar */ 1`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []*Comment{
					{
						ASTNode:   ASTNode{},
						Multiline: "/* foo\n   bar */",
					},
				},
				Value: BuildTestExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
			},
		},
		{
			name: "One multiline - one line - with expr on new line",
			input: `
/* foo */
1`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []*Comment{
					{
						ASTNode:   ASTNode{},
						Multiline: "/* foo */\n",
					},
				},
				Value: BuildTestExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
			},
		},
		{
			name: "One multiline - two lines - with expr on new line",
			input: `
/* foo
   bar */
1`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []*Comment{
					{
						ASTNode:   ASTNode{},
						Multiline: "/* foo\n   bar */\n",
					},
				},
				Value: BuildTestExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, false)
		})
	}
}

func TestComment_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *Comment
		want  *Comment
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &Comment{},
			want:  &Comment{},
		},
		{
			name: "ASTNode",
			input: &Comment{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
			want: &Comment{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
		},
		{
			name: "SingleLine",
			input: &Comment{
				SingleLine: []string{"foo"},
			},
			want: &Comment{
				SingleLine: []string{"foo"},
			},
		},
		{
			name: "Multiline",
			input: &Comment{
				Multiline: "foo",
			},
			want: &Comment{
				Multiline: "foo",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*Comment](t, tt.want, tt.input.Clone())
		})
	}
}

func TestComment_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *Comment
		want  []Node
	}{
		{
			name:  "Empty",
			input: &Comment{},
			want:  nil,
		},
		{
			name: "SingleLine",
			input: &Comment{
				SingleLine: []string{"foo"},
			},
			want: nil,
		},
		{
			name: "Multiline",
			input: &Comment{
				Multiline: "foo",
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

func TestComment_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     *Comment
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
			input: &Comment{},
			want:  "",
		},
		{
			name: "One SingleLine",
			input: &Comment{
				SingleLine: []string{"// foo"},
			},
			want: "// foo\n",
		},
		{
			name: "Two SingleLine",
			input: &Comment{
				SingleLine: []string{"// foo", "# bar"},
			},
			want: "// foo\n# bar\n",
		},
		{
			name: "One MultiLine",
			input: &Comment{
				Multiline: "/* foo */",
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
