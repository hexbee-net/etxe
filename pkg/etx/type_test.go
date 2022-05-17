package etx

import (
	"testing"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestType_Parsing(t *testing.T) {
	type args struct {
		Input string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    *Type
	}{
		{
			name: "Empty Enum",
			args: args{
				Input: `type foo enum {}`,
			},
			wantErr: false,
			want: &Type{
				Label:  "foo",
				Enum:   &Enum{},
				Object: nil,
			},
		},
		{
			name: "Empty Object",
			args: args{
				Input: `type foo object {}`,
			},
			wantErr: false,
			want: &Type{
				Label:  "foo",
				Enum:   nil,
				Object: &Object{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := participle.MustBuild(
				&Type{},
				participle.Lexer(lexer.MustStateful(lexRules())),
				participle.Elide(TokenWhitespace),
			)

			res := &Type{}
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
