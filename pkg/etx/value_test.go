package etx

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValue_Parsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Value
	}{
		{
			name:    "Null",
			input:   `null`,
			wantErr: false,
			want: &Value{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Null:    true,
			},
		},
		{
			name:    "Bool - true",
			input:   `true`,
			wantErr: false,
			want: &Value{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Bool:    func() *ValueBool { v := true; return (*ValueBool)(&v) }(),
			},
		},
		{
			name:    "Bool - false",
			input:   `false`,
			wantErr: false,
			want: &Value{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Bool:    func() *ValueBool { v := false; return (*ValueBool)(&v) }(),
			},
		},
		{
			name:    "Number",
			input:   `12.34e+2`,
			wantErr: false,
			want: &Value{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Number: &ValueNumber{
					Value:  big.NewFloat(12.34e+2),
					Source: `-12.34e+2`,
				},
			},
		},
		{
			name:    "String",
			input:   `"hello world"`,
			wantErr: false,
			want: &Value{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Str: &ValueString{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Fragment: []*StringFragment{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 1, Column: 2}},
							Text:    `hello world`,
						},
					},
				},
			},
		},
		{
			name:    "Ident",
			input:   `var`,
			wantErr: false,
			want: &Value{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Ident: &Ident{
					Parts: []string{"var"},
				},
			},
		},
		{
			name: "Heredoc - Empty",
			input: `
<<FOO
FOO`[1:],
			wantErr: false,
			want: &Value{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Heredoc: &Heredoc{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Delimiter: HeredocDelimiter{
						LeadingTabs: false,
						Delimiter:   "FOO",
					},
					Value: nil,
				},
			},
		},
		{
			name: "Heredoc - No leading tabs",
			input: `
<<FOO
bar
FOO`[1:],
			wantErr: false,
			want: &Value{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Heredoc: &Heredoc{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Delimiter: HeredocDelimiter{
						LeadingTabs: false,
						Delimiter:   "FOO",
					},
					Value: testValPtr(t, "bar\n"),
				},
			},
		},
		{
			name: "Heredoc - Leading tabs",
			input: `
<<-FOO
bar
FOO`[1:],
			wantErr: false,
			want: &Value{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Heredoc: &Heredoc{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Delimiter: HeredocDelimiter{
						LeadingTabs: true,
						Delimiter:   "FOO",
					},
					Value: testValPtr(t, "bar\n"),
				},
			},
		},
		{
			name:    "List - Empty",
			input:   `[]`,
			wantErr: false,
			want: &Value{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				List: &ValueList{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
		},
		{
			name:    "List - Elements",
			input:   `[ 1, 2 ]`,
			wantErr: false,
			want: &Value{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				List: &ValueList{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Items: []*Value{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Number:  &ValueNumber{big.NewFloat(1), "1"},
						},
						{
							ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
							Number:  &ValueNumber{big.NewFloat(2), "2"},
						},
					},
				},
			},
		},
		{
			name:    "Map - Empty",
			input:   `{}`,
			wantErr: false,
			want: &Value{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Map: &ValueMap{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				},
			},
		},
		{
			name:    "Map - Elements",
			input:   `{ foo : bar , baz : qux }`,
			wantErr: false,
			want: &Value{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Map: &ValueMap{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Entries: []*MapEntry{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Key: Value{
								ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
								Ident:   &Ident{Parts: []string{"foo"}},
							},
							Value: Value{
								ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
								Ident:   &Ident{Parts: []string{"bar"}},
							},
						},
						{
							ASTNode: ASTNode{Pos: Position{Offset: 14, Line: 1, Column: 15}},
							Key: Value{
								ASTNode: ASTNode{Pos: Position{Offset: 14, Line: 1, Column: 15}},
								Ident:   &Ident{Parts: []string{"baz"}},
							},
							Value: Value{
								ASTNode: ASTNode{Pos: Position{Offset: 20, Line: 1, Column: 21}},
								Ident:   &Ident{Parts: []string{"qux"}},
							},
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

func TestValue_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *Value
		want  *Value
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &Value{},
			want:  &Value{},
		},
		{
			name: "ASTNode",
			input: &Value{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
			want: &Value{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
		},
		{
			name: "Bool Value",
			input: &Value{
				Bool: testValPtr[ValueBool](t, true),
			},
			want: &Value{
				Bool: testValPtr[ValueBool](t, true),
			},
		},
		{
			name: "Number Value",
			input: &Value{
				Number: &ValueNumber{big.NewFloat(1), "1"},
			},
			want: &Value{
				Number: &ValueNumber{big.NewFloat(1), "1"},
			},
		},
		{
			name: "Str Value",
			input: &Value{
				Str: &ValueString{
					Fragment: []*StringFragment{},
				},
			},
			want: &Value{
				Str: &ValueString{
					Fragment: []*StringFragment{},
				},
			},
		},
		{
			name: "Ident Value",
			input: &Value{
				Ident: &Ident{
					Parts: []string{"foo"},
				},
			},
			want: &Value{
				Ident: &Ident{
					Parts: []string{"foo"},
				},
			},
		},
		{
			name: "Heredoc Value",
			input: &Value{
				Heredoc: &Heredoc{
					Delimiter: HeredocDelimiter{
						Delimiter: "foo",
					},
					Value: testValPtr(t, "bar"),
				},
			},
			want: &Value{
				Heredoc: &Heredoc{
					Delimiter: HeredocDelimiter{
						Delimiter: "foo",
					},
					Value: testValPtr(t, "bar"),
				},
			},
		},
		{
			name: "List Value",
			input: &Value{
				List: &ValueList{
					Items: []*Value{},
				},
			},
			want: &Value{
				List: &ValueList{
					Items: []*Value{},
				},
			},
		},
		{
			name: "Map Value",
			input: &Value{
				Map: &ValueMap{
					Entries: []*MapEntry{},
				},
			},
			want: &Value{
				Map: &ValueMap{
					Entries: []*MapEntry{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clone := tt.input.Clone()
			testCloner[*Value](t, tt.want, clone)
		})
	}
}

func TestValue_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *Value
		want  []Node
	}{
		{
			name:  "Empty",
			input: &Value{},
			want:  nil,
		},
		{
			name: "Bool Value",
			input: &Value{
				Bool: testValPtr[ValueBool](t, true),
			},
			want: []Node{
				func() Node { v := ValueBool(true); return &v }(),
			},
		},
		{
			name: "Number Value",
			input: &Value{
				Number: &ValueNumber{big.NewFloat(1), "1"},
			},
			want: []Node{
				&ValueNumber{big.NewFloat(1), "1"},
			},
		},
		{
			name: "Str Value",
			input: &Value{
				Str: &ValueString{
					Fragment: []*StringFragment{},
				},
			},
			want: []Node{
				&ValueString{
					Fragment: []*StringFragment{},
				},
			},
		},
		{
			name: "Ident Value",
			input: &Value{
				Ident: &Ident{
					Parts: []string{"foo"},
				},
			},
			want: []Node{
				&Ident{
					Parts: []string{"foo"},
				},
			},
		},
		{
			name: "Heredoc Value",
			input: &Value{
				Heredoc: &Heredoc{
					Delimiter: HeredocDelimiter{
						Delimiter: "foo",
					},
					Value: testValPtr(t, "bar"),
				},
			},
			want: []Node{
				&Heredoc{
					Delimiter: HeredocDelimiter{
						Delimiter: "foo",
					},
					Value: testValPtr(t, "bar"),
				},
			},
		},
		{
			name: "List Value",
			input: &Value{
				List: &ValueList{
					Items: []*Value{},
				},
			},
			want: []Node{
				&ValueList{
					Items: []*Value{},
				},
			},
		},
		{
			name: "Map Value",
			input: &Value{
				Map: &ValueMap{
					Entries: []*MapEntry{},
				},
			},
			want: []Node{
				&ValueMap{
					Entries: []*MapEntry{},
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

func TestValue_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		description string
		input       *Value
		wantX       *Value
		wantPanic   bool
		want        string
	}{
		{
			name:      "Nil",
			input:     nil,
			wantPanic: true,
		},
		{
			name:        "Empty",
			description: "At least one field must be initialized",
			input:       &Value{},
			wantPanic:   true,
		},

		{
			name: "Null",
			input: &Value{
				Null: true,
			},
			want: "null",
		},
		{
			name: "Bool Value",
			input: &Value{
				Bool: testValPtr[ValueBool](t, true),
			},
			want: "true",
		},
		{
			name: "Number Value",
			input: &Value{
				Number: &ValueNumber{big.NewFloat(1), "1"},
			},
			want: "1",
		},
		{
			name: "Str Value",
			input: &Value{
				Str: &ValueString{
					Fragment: []*StringFragment{},
				},
			},
			want: `""`,
		},
		{
			name: "Ident Value",
			input: &Value{
				Ident: &Ident{
					Parts: []string{"foo"},
				},
			},
			want: "foo",
		},
		{
			name: "Heredoc Value",
			input: &Value{
				Heredoc: &Heredoc{
					Delimiter: HeredocDelimiter{
						Delimiter: "FOO",
					},
					Value: testValPtr(t, "\nbar"),
				},
			},
			want: `
<<FOO
bar
FOO`[1:],
		},
		{
			name: "List Value",
			input: &Value{
				List: &ValueList{
					Items: []*Value{},
				},
			},
			want: "[]",
		},
		{
			name: "Map Value",
			input: &Value{
				Map: &ValueMap{
					Entries: []*MapEntry{},
				},
			},
			want: "{}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input, tt.description)
		})
	}
}
