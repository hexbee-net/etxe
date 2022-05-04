package etx

import (
	"fmt"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/repr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lex()
			parser := participle.MustBuild(&Value{}, participle.Lexer(l))
			fmt.Println(parser.String())

			res := &Value{}
			err := parser.ParseString("", tt.args.Input, res)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			repr.Println(res)

			assert.Equal(t, tt.want, res)
		})
	}
}
