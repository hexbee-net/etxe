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
				Bool:    &ValueBool{Value: true},
			},
		},
		{
			name:    "Bool - false",
			input:   `false`,
			wantErr: false,
			want: &Value{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Bool:    &ValueBool{Value: false},
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
					Fragments: nil,
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
					Fragments: []*HeredocFragment{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 2, Column: 1}},
							Text:    "bar\n",
						},
					},
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
					Fragments: []*HeredocFragment{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 7, Line: 2, Column: 1}},
							Text:    "bar\n",
						},
					},
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
					Items: []*ListItem{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Value: BuildTestExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
								Number: &ValueNumber{
									ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
									Value:   big.NewFloat(1),
									Source:  "1",
								},
							}),
						},
						{
							ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
							Value: BuildTestExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
								Number: &ValueNumber{
									ASTNode: ASTNode{Pos: Position{Offset: 5, Line: 1, Column: 6}},
									Value:   big.NewFloat(2),
									Source:  "2",
								},
							}),
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
			input:   `{ foo = bar , baz = qux }`,
			wantErr: false,
			want: &Value{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Map: &ValueMap{
					ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
					Items: []*MapItem{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
							Key: &MapKey{
								ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
								Ident: &Ident{
									ASTNode: ASTNode{Pos: Position{Offset: 2, Line: 1, Column: 3}},
									Parts:   []string{"foo"},
								},
							},
							Value: BuildTestExprTree[*Expr](t, &Ident{
								ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 1, Column: 9}},
								Parts:   []string{"bar"},
							}),
						},
						{
							ASTNode: ASTNode{Pos: Position{Offset: 14, Line: 1, Column: 15}},
							Key: &MapKey{
								ASTNode: ASTNode{Pos: Position{Offset: 14, Line: 1, Column: 15}},
								Ident: &Ident{
									ASTNode: ASTNode{Pos: Position{Offset: 14, Line: 1, Column: 15}},
									Parts:   []string{"baz"},
								},
							},
							Value: BuildTestExprTree[*Expr](t, &Ident{
								ASTNode: ASTNode{Pos: Position{Offset: 20, Line: 1, Column: 21}},
								Parts:   []string{"qux"},
							}),
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
				Bool: &ValueBool{Value: true},
			},
			want: &Value{
				Bool: &ValueBool{Value: true},
			},
		},
		{
			name: "Number Value",
			input: &Value{
				Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"},
			},
			want: &Value{
				Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"},
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
			name: "Heredoc Value",
			input: &Value{
				Heredoc: &Heredoc{
					Delimiter: HeredocDelimiter{
						Delimiter: "foo",
					},
					Fragments: []*HeredocFragment{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Text:    "bar",
						},
					},
				},
			},
			want: &Value{
				Heredoc: &Heredoc{
					Delimiter: HeredocDelimiter{
						Delimiter: "foo",
					},
					Fragments: []*HeredocFragment{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Text:    "bar",
						},
					},
				},
			},
		},
		{
			name: "List Value",
			input: &Value{
				List: &ValueList{
					Items: []*ListItem{},
				},
			},
			want: &Value{
				List: &ValueList{
					Items: []*ListItem{},
				},
			},
		},
		{
			name: "Map Value",
			input: &Value{
				Map: &ValueMap{
					Items: []*MapItem{},
				},
			},
			want: &Value{
				Map: &ValueMap{
					Items: []*MapItem{},
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
				Bool: &ValueBool{Value: true},
			},
			want: []Node{
				&ValueBool{Value: true},
			},
		},
		{
			name: "Number Value",
			input: &Value{
				Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"},
			},
			want: []Node{
				&ValueNumber{Value: big.NewFloat(1), Source: "1"},
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
			name: "Heredoc Value",
			input: &Value{
				Heredoc: &Heredoc{
					Delimiter: HeredocDelimiter{
						Delimiter: "foo",
					},
					Fragments: []*HeredocFragment{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Text:    "bar",
						},
					},
				},
			},
			want: []Node{
				&Heredoc{
					Delimiter: HeredocDelimiter{
						Delimiter: "foo",
					},
					Fragments: []*HeredocFragment{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Text:    "bar",
						},
					},
				},
			},
		},
		{
			name: "List Value",
			input: &Value{
				List: &ValueList{
					Items: []*ListItem{},
				},
			},
			want: []Node{
				&ValueList{
					Items: []*ListItem{},
				},
			},
		},
		{
			name: "Map Value",
			input: &Value{
				Map: &ValueMap{
					Items: []*MapItem{},
				},
			},
			want: []Node{
				&ValueMap{
					Items: []*MapItem{},
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

func TestValue_FormattedString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		description string
		input       *Value
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
				Bool: &ValueBool{Value: true},
			},
			want: "true",
		},
		{
			name: "Number Value",
			input: &Value{
				Number: &ValueNumber{Value: big.NewFloat(1), Source: "1"},
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
			name: "Heredoc Value",
			input: &Value{
				Heredoc: &Heredoc{
					Delimiter: HeredocDelimiter{
						Delimiter: "FOO",
					},
					Fragments: []*HeredocFragment{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Text:    "\nbar",
						},
					},
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
					Items: []*ListItem{},
				},
			},
			want: "[]",
		},
		{
			name: "Map Value",
			input: &Value{
				Map: &ValueMap{
					Items: []*MapItem{},
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
