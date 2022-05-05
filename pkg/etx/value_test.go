package etx

import (
	"github.com/alecthomas/participle/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestValue_Parsing(t *testing.T) {
	type args struct {
		Input string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    *Value
	}{
		{
			name: "Bool - true",
			args: args{
				Input: `true`,
			},
			wantErr: false,
			want: &Value{
				Bool: func() *Bool { v := true; return (*Bool)(&v) }(),
			},
		},
		{
			name: "Bool - false",
			args: args{
				Input: `false`,
			},
			wantErr: false,
			want: &Value{
				Bool: func() *Bool { v := false; return (*Bool)(&v) }(),
			},
		},
		{
			name: "Number",
			args: args{
				Input: `-12.34e+2`,
			},
			wantErr: false,
			want: &Value{
				Number: &Number{Float: big.NewFloat(-12.34e+2)},
			},
		},
		{
			name: "String",
			args: args{
				Input: `"hello world"`,
			},
			wantErr: false,
			want: &Value{
				Str: String{
					{Text: `hello world`},
				},
			},
		},
		{
			name: "Ident",
			args: args{
				Input: `var`,
			},
			wantErr: false,
			want: &Value{
				Ident: func() *string { v := "var"; return &v }(),
			},
		},
		{
			name: "List - Empty",
			args: args{
				Input: `[]`,
			},
			wantErr: false,
			want: &Value{
				HaveList: true,
			},
		},
		{
			name: "List - Elements",
			args: args{
				Input: `[ 1, 2 ]`,
			},
			wantErr: false,
			want: &Value{
				HaveList: true,
				List: []*Value{
					{Number: &Number{big.NewFloat(1)}},
					{Number: &Number{big.NewFloat(2)}},
				},
			},
		},
		{
			name: "Map - Empty",
			args: args{
				Input: `{}`,
			},
			wantErr: false,
			want: &Value{
				HaveMap: true,
			},
		},
		{
			name: "Map - Elements",
			args: args{
				Input: `{ foo : bar , baz : qux }`,
			},
			wantErr: false,
			want: &Value{
				HaveMap: true,
				Map: []*MapEntry{
					{
						Key:   &Value{Ident: func() *string { v := "foo"; return &v }()},
						Value: &Value{Ident: func() *string { v := "bar"; return &v }()},
					},
					{
						Key:   &Value{Ident: func() *string { v := "baz"; return &v }()},
						Value: &Value{Ident: func() *string { v := "qux"; return &v }()},
					},
				},
			},
		},
		{
			name: "Heredoc - Empty",
			args: args{
				Input: `<<FOO
FOO`,
			},
			wantErr: false,
			want: &Value{
				HeredocDelimiter: "<<FOO",
				// Heredoc:          func() *string { v := ""; return &v }(),
			},
		},
		{
			name: "Heredoc - Content",
			args: args{
				Input: `<<FOO
bar
FOO`,
			},
			wantErr: false,
			want: &Value{
				HeredocDelimiter: "<<FOO",
				Heredoc:          func() *string { v := "\nbar"; return &v }(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lex()
			parser := participle.MustBuild(&Value{}, participle.Lexer(l))

			res := &Value{}
			err := parser.ParseString("", tt.args.Input, res)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, res)
		})
	}
}
