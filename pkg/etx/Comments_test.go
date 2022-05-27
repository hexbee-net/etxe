package etx

import (
	"math/big"
	"testing"
)

func TestCommentNode_Parsing(t *testing.T) {
	t.Parallel()

	type TestStruct struct {
		Comments []string `parser:"(@Comment [ NewLine ])*" json:"comments,omitempty"`
		Value    *Expr    `parser:"[ @@ ]"`
	}
	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *TestStruct
	}{
		{
			name: "None",
			input: `
1`[1:],
			wantErr: false,
			want: &TestStruct{
				Value: testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "2"}}),
			},
		},

		{
			name: "One slashed",
			input: `
// foo`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []string{"// foo"},
				Value:    nil,
			},
		},
		{
			name: "One slashed before value",
			input: `
// foo
1`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []string{"// foo\n"},
				Value:    testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
			},
		},
		{
			name: "Two slashed",
			input: `
// foo
// bar`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []string{"// foo\n", "// bar"},
				Value:    nil,
			},
		},
		{
			name: "Two slashed before value",
			input: `
// foo
// bar
1`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []string{"// foo\n", "// bar\n"},
				Value:    testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
			},
		},

		{
			name: "One hashtag",
			input: `
# foo`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []string{"# foo"},
				Value:    nil,
			},
		},
		{
			name: "Two hashtags",
			input: `
# foo
# bar`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []string{"# foo\n", "# bar"},
				Value:    nil,
			},
		},

		{
			name: "One multiline - one line",
			input: `
/* foo */`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []string{"/* foo */"},
				Value:    nil,
			},
		},
		{
			name: "Two multilines - one line",
			input: `
/* foo */ /* bar */`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []string{"/* foo */", "/* bar */"},
				Value:    nil,
			},
		},
		{
			name: "Two multilines - two line",
			input: `
/* foo */
/* bar */`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []string{"/* foo */", "/* bar */"},
				Value:    nil,
			},
		},
		{
			name: "One multiline - two lines",
			input: `
/* foo
   bar */`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []string{"/* foo\n   bar */"},
				Value:    nil,
			},
		},
		{
			name: "One multiline - one line - with value on same line",
			input: `
/* foo */ 1`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []string{"/* foo */"},
				Value:    testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
			},
		},
		{
			name: "One multiline - two lines - with value on same line",
			input: `
/* foo
   bar */ 1`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []string{"/* foo\n   bar */"},
				Value:    testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
			},
		},
		{
			name: "One multiline - one line - with value on new line",
			input: `
/* foo */
1`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []string{"/* foo */"},
				Value:    testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
			},
		},
		{
			name: "One multiline - two lines - with value on new line",
			input: `
/* foo
   bar */
1`[1:],
			wantErr: false,
			want: &TestStruct{
				Comments: []string{"/* foo\n   bar */"},
				Value:    testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, false)
		})
	}
}
