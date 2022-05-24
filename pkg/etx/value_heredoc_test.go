package etx

import (
	"testing"
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
				Value: nil,
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
				Value: testValPtr(t, "\n"),
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
				Value: testValPtr(t, "\n\nfoo"),
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
				Value: testValPtr(t, "\nfoo\n"),
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
				Value: testValPtr(t, "\nfoo"),
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
				Value: testValPtr(t, "\nfoo"),
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
			name: "Delimiter and Value",
			input: &Heredoc{
				Delimiter: HeredocDelimiter{
					LeadingTabs: false,
					Delimiter:   "FOO",
				},
				Value: testValPtr(t, "bar"),
			},
			want: &Heredoc{
				Delimiter: HeredocDelimiter{
					LeadingTabs: false,
					Delimiter:   "FOO",
				},
				Value: testValPtr(t, "bar"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*Heredoc](t, tt.want, tt.input)
		})
	}
}

func TestHeredoc_String(t *testing.T) {
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
				Value: testValPtr(t, "\n\n"),
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
				Value: testValPtr(t, "\n\nfoo"),
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
				Value: testValPtr(t, "\nfoo\n"),
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
				Value: testValPtr(t, "\nfoo"),
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
				Value: testValPtr(t, "\nfoo"),
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
