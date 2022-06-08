package etx

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeredoc_Parsing(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Heredoc
	}{
		{
			name: "Empty",
			input: `
<<EOF
EOF`[1:],
			wantErr: false,
			want: &Heredoc{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Delimiter: HeredocDelimiter{
					LeadingTabs: false,
					Delimiter:   "EOF",
				},
				Fragments: nil,
			},
		},
		{
			name: "Only linebreaks",
			input: `
<<EOF

EOF`[1:],
			wantErr: false,
			want: &Heredoc{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Delimiter: HeredocDelimiter{
					LeadingTabs: false,
					Delimiter:   "EOF",
				},
				Fragments: []*HeredocFragment{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 1}},
						Text:    "\n",
					},
				},
			},
		},
		{
			name: "value with leading linebreaks",
			input: `
<<EOF

foo
EOF`[1:],
			wantErr: false,
			want: &Heredoc{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Delimiter: HeredocDelimiter{
					LeadingTabs: false,
					Delimiter:   "EOF",
				},
				Fragments: []*HeredocFragment{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 1}},
						Text:    "\nfoo\n",
					},
				},
			},
		},
		{
			name: "value with trailing linebreaks",
			input: `
<<EOF
foo

EOF`[1:],
			wantErr: false,
			want: &Heredoc{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Delimiter: HeredocDelimiter{
					LeadingTabs: false,
					Delimiter:   "EOF",
				},
				Fragments: []*HeredocFragment{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 1}},
						Text:    "foo\n\n",
					},
				},
			},
		},
		{
			name: "value - no trailing tabs",
			input: `
<<EOF
foo
EOF`[1:],
			wantErr: false,
			want: &Heredoc{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Delimiter: HeredocDelimiter{
					LeadingTabs: false,
					Delimiter:   "EOF",
				},
				Fragments: []*HeredocFragment{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 1}},
						Text:    "foo\n",
					},
				},
			},
		},
		{
			name: "value - trailing tabs",
			input: `
<<-EOF
foo
EOF`[1:],
			wantErr: false,
			want: &Heredoc{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Delimiter: HeredocDelimiter{
					LeadingTabs: true,
					Delimiter:   "EOF",
				},
				Fragments: []*HeredocFragment{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 2, Column: 1}},
						Text:    "foo\n",
					},
				},
			},
		},

		{
			name: "value - expression",
			input: `
<<EOF
foo ${ 1 } bar
EOF`[1:],
			wantErr: false,
			want: &Heredoc{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Delimiter: HeredocDelimiter{
					LeadingTabs: false,
					Delimiter:   "EOF",
				},
				Fragments: []*HeredocFragment{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 1}},
						Text:    "foo ",
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 10, Line: 2, Column: 5}},
						Expr: BuildTestExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 2, Column: 8}},
							Number: &ValueNumber{
								ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 2, Column: 8}},
								Value:   big.NewFloat(1),
								Source:  "1",
							},
						}),
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 16, Line: 2, Column: 11}},
						Text:    " bar\n",
					},
				},
			},
		},

		{
			name: "value - directive",
			input: `
<<EOF
foo %{ 1 } bar
EOF`[1:],
			wantErr: false,
			want: &Heredoc{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Delimiter: HeredocDelimiter{
					LeadingTabs: false,
					Delimiter:   "EOF",
				},
				Fragments: []*HeredocFragment{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 1}},
						Text:    "foo ",
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 10, Line: 2, Column: 5}},
						Directive: BuildTestExprTree[*Expr](t, &Value{
							ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 2, Column: 8}},
							Number: &ValueNumber{
								ASTNode: ASTNode{Pos: Position{Offset: 13, Line: 2, Column: 8}},
								Value:   big.NewFloat(1),
								Source:  "1",
							},
						}),
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 16, Line: 2, Column: 11}},
						Text:    " bar\n",
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

func TestHeredoc_Clone(t *testing.T) {
	tests := []struct {
		name  string
		input *Heredoc
		want  *Heredoc
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &Heredoc{},
			want:  &Heredoc{},
		},
		{
			name: "ASTNode",
			input: &Heredoc{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
			want: &Heredoc{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
		},
		{
			name: "Delimiter and Value",
			input: &Heredoc{
				Delimiter: HeredocDelimiter{
					LeadingTabs: false,
					Delimiter:   "FOO",
				},
				Fragments: []*HeredocFragment{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Text:    "bar",
					},
				},
			},
			want: &Heredoc{
				Delimiter: HeredocDelimiter{
					LeadingTabs: false,
					Delimiter:   "FOO",
				},
				Fragments: []*HeredocFragment{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Text:    "bar",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*Heredoc](t, tt.want, tt.input)
		})
	}
}

func TestHeredoc_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *Heredoc
		want  []Node
	}{
		{
			name:  "Empty",
			input: &Heredoc{},
			want:  nil,
		},
		{
			name: "Delimiter",
			input: &Heredoc{
				Delimiter: HeredocDelimiter{
					Delimiter: "foo",
				},
			},
			want: nil,
		},
		{
			name: "Fragments",
			input: &Heredoc{
				Fragments: []*HeredocFragment{
					{Text: "foo"},
				},
			},
			want: []Node{
				&HeredocFragment{Text: "foo"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestHeredoc_FormattedString(t *testing.T) {
	tests := []struct {
		name      string
		input     *Heredoc
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
			input:     &Heredoc{},
			wantPanic: true,
		},
		{
			name: "No value",
			input: &Heredoc{
				Delimiter: HeredocDelimiter{
					LeadingTabs: false,
					Delimiter:   "EOF",
				},
			},
			want: `
<<EOF
EOF`[1:],
		},
		{
			name: "Only linebreaks",
			input: &Heredoc{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Delimiter: HeredocDelimiter{
					LeadingTabs: false,
					Delimiter:   "EOF",
				},
				Fragments: []*HeredocFragment{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Text:    "\n\n",
					},
				},
			},
			want: `
<<EOF


EOF`[1:],
		},
		{
			name: "value with leading linebreaks",
			input: &Heredoc{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Delimiter: HeredocDelimiter{
					LeadingTabs: false,
					Delimiter:   "EOF",
				},
				Fragments: []*HeredocFragment{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Text:    "\n\nfoo",
					},
				},
			},
			want: `
<<EOF

foo
EOF`[1:],
		},
		{
			name: "value with trailing linebreaks",
			input: &Heredoc{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Delimiter: HeredocDelimiter{
					LeadingTabs: false,
					Delimiter:   "EOF",
				},
				Fragments: []*HeredocFragment{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Text:    "\nfoo\n",
					},
				},
			},
			want: `
<<EOF
foo

EOF`[1:],
		},
		{
			name: "value - no trailing tabs",
			input: &Heredoc{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Delimiter: HeredocDelimiter{
					LeadingTabs: false,
					Delimiter:   "EOF",
				},
				Fragments: []*HeredocFragment{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Text:    "\nfoo",
					},
				},
			},
			want: `
<<EOF
foo
EOF`[1:],
		},
		{
			name: "value - trailing tabs",
			input: &Heredoc{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Delimiter: HeredocDelimiter{
					LeadingTabs: true,
					Delimiter:   "EOF",
				},
				Fragments: []*HeredocFragment{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Text:    "\nfoo",
					},
				},
			},
			want: `
<<-EOF
foo
EOF`[1:],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}

// /////////////////////////////////////

func TestHeredocFragment_Clone(t *testing.T) {
	tests := []struct {
		name  string
		input *HeredocFragment
		want  *HeredocFragment
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &HeredocFragment{},
			want:  &HeredocFragment{},
		},
		{
			name: "ASTNode",
			input: &HeredocFragment{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
			want: &HeredocFragment{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
		},
		{
			name: "Expr",
			input: &HeredocFragment{
				Expr: BuildTestExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
			},
			want: &HeredocFragment{
				Expr: BuildTestExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
			},
		},
		{
			name: "Directive",
			input: &HeredocFragment{
				Directive: BuildTestExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
			},
			want: &HeredocFragment{
				Directive: BuildTestExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
			},
		},
		{
			name: "Text",
			input: &HeredocFragment{
				Text: "foo",
			},
			want: &HeredocFragment{
				Text: "foo",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*HeredocFragment](t, tt.want, tt.input)
		})
	}

}

func TestHeredocFragment_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *HeredocFragment
		want  []Node
	}{
		{
			name:  "Empty",
			input: &HeredocFragment{},
			want:  nil,
		},
		{
			name: "Expr",
			input: &HeredocFragment{
				Expr: BuildTestExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
			},
			want: []Node{
				BuildTestExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
			},
		},
		{
			name: "Directive",
			input: &HeredocFragment{
				Directive: BuildTestExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
			},
			want: []Node{
				BuildTestExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
			},
		},
		{
			name: "Text",
			input: &HeredocFragment{
				Text: "foo",
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

func TestHeredocFragment_FormattedString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		description string
		input       *HeredocFragment
		wantPanic   bool
		want        string
	}{
		{
			name:      "Nil",
			input:     nil,
			wantPanic: true,
		},
		{
			name:  "Empty",
			input: &HeredocFragment{},
			want:  "",
		},

		{
			name: "Expr",
			input: &HeredocFragment{
				Expr: BuildTestExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
			},
			want: "${ 1 }",
		},
		{
			name: "Directive",
			input: &HeredocFragment{
				Directive: BuildTestExprTree[*Expr](t, &Value{Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
			},
			want: "%{ 1 }",
		},
		{
			name: "Text",
			input: &HeredocFragment{
				Text: "foo",
			},
			want: "foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input, tt.description)
		})
	}
}
