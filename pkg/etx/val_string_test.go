package etx

import (
	"testing"

	"github.com/alecthomas/participle/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestString_Parsing(t *testing.T) {
	type args struct {
		Input string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    String
	}{
		{
			name: "Double quoted",
			args: args{
				Input: `"hello world"`,
			},
			wantErr: false,
			want: String{
				{Text: `hello world`},
			},
		},
		{
			name: "Double quoted with single quote",
			args: args{
				Input: `"hello ' world"`,
			},
			wantErr: false,
			want: String{
				{Text: `hello `},
				{Text: `'`},
				{Text: ` world`},
			},
		},

		{
			name: "Single quoted",
			args: args{
				Input: `'hello world'`,
			},
			wantErr: false,
			want: String{
				{Text: `hello world`},
			},
		},
		{
			name: "Single quoted with double quote",
			args: args{
				Input: `'hello " world'`,
			},
			wantErr: false,
			want: String{
				{Text: `hello `},
				{Text: `"`},
				{Text: ` world`},
			},
		},

		{
			name: "Escaped",
			args: args{
				Input: `"hello \t world"`,
			},
			wantErr: false,
			want: String{
				{Text: `hello `},
				{Escaped: `\t`},
				{Text: ` world`},
			},
		},

		{
			name: "Unicode - Short",
			args: args{
				Input: `"hello \u1234 world"`,
			},
			wantErr: false,
			want: String{
				{Text: `hello `},
				{Unicode: `1234`},
				{Text: ` world`},
			},
		},
		{
			name: "Unicode - Long",
			args: args{
				Input: `"hello \u12345678 world"`,
			},
			wantErr: false,
			want: String{
				{Text: `hello `},
				{Unicode: `12345678`},
				{Text: ` world`},
			},
		},
		{
			name: "Unicode - Invalid",
			args: args{
				Input: `"hello \u123 world"`,
			},
			wantErr: true,
			want:    nil,
		},
		{
			name: "Unicode - Trailing numbers",
			args: args{
				Input: `"hello \u123456 world"`,
			},
			wantErr: false,
			want: String{
				{Text: `hello `},
				{Unicode: `1234`},
				{Text: `56 world`},
			},
		},

		{
			name: "Expression",
			args: args{
				Input: `"hello ${foo} world"`,
			},
			wantErr: false,
			want: String{
				{Text: `hello `},
				{Expr: &Expr{Left: &Terminal{Ident: "foo"}}},
				{Text: ` world`},
			},
		},
		{
			name: "Non-Expression",
			args: args{
				Input: `"hello $${foo} world"`,
			},
			wantErr: false,
			want: String{
				{Text: `hello `},
				{Text: `$${`},
				{Text: `foo} world`},
			},
		},

		{
			name: "Directive",
			args: args{
				Input: `"hello %{foo} world"`,
			},
			wantErr: false,
			want: String{
				{Text: `hello `},
				{Directive: &Expr{Left: &Terminal{Ident: "foo"}}},
				{Text: ` world`},
			},
		},
		{
			name: "Non-Directive",
			args: args{
				Input: `"hello %%{foo} world"`,
			},
			wantErr: false,
			want: String{
				{Text: `hello `},
				{Text: `%%{`},
				{Text: `foo} world`},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type Str struct {
				Str String `parser:"String @@* StringEnd"`
			}
			parser := participle.MustBuild(&Str{}, participle.Lexer(lex()))

			res := &Str{}
			err := parser.ParseString("", tt.args.Input, res)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			want := &Str{
				Str: tt.want,
			}
			assert.Equal(t, want, res)
		})
	}
}

func TestString_String(t *testing.T) {
	type args struct {
		Input String
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    string
	}{
		{
			name: "One Text Fragment",
			args: args{
				Input: String{
					{Text: `hello world`},
				},
			},
			wantErr: false,
			want:    `hello world`,
		},
		{
			name: "Several Text Fragments",
			args: args{
				Input: String{
					{Text: `hello `},
					{Text: `world`},
				},
			},
			wantErr: false,
			want:    `hello world`,
		},
		{
			name: "Text and Escaped",
			args: args{
				Input: String{
					{Text: `hello `},
					{Escaped: `\t`},
					{Text: ` world`},
				},
			},
			wantErr: false,
			want:    `hello \t world`,
		},
		{
			name: "Text and Unicode",
			args: args{
				Input: String{
					{Text: `hello `},
					{Unicode: `1234`},
					{Text: ` world`},
				},
			},
			wantErr: false,
			want:    `hello \u1234 world`,
		},
		{
			name: "Text and Expression",
			args: args{
				Input: String{
					{Text: `hello `},
					{Expr: &Expr{Left: &Terminal{Ident: "foo"}}},
					{Text: ` world`},
				},
			},
			wantErr: false,
			want:    `hello ${foo} world`,
		},
		{
			name: "Text and Directive",
			args: args{
				Input: String{
					{Text: `hello `},
					{Directive: &Expr{Left: &Terminal{Ident: "foo"}}},
					{Text: ` world`},
				},
			},
			wantErr: false,
			want:    `hello %{foo} world`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.args.Input.String())
		})
	}
}
