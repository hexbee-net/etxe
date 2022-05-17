package etx

import (
	"math/big"
	"testing"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLambda_Parsing(t *testing.T) {
	type args struct {
		Input string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    *Lambda
	}{
		{
			name: "",
			args: args{
				Input: `(x: number) => x + 1`,
			},
			wantErr: false,
			want: &Lambda{
				Parameters: []*LambdaParameter{
					{
						Label: "x",
						Type: &ParameterType{
							Ident: &Ident{Parts: []string{"number"}},
						},
					},
				},
				Expr: testBuildExprTree[*Expr](t,
					&ExprAdditive{
						Left:  testBuildExprTree[*ExprMultiplicative](t, &Value{Ident: &Ident{Parts: []string{"x"}}}),
						Op:    OpPlus,
						Right: testBuildExprTree[*ExprAdditive](t, &Value{Number: &Number{big.NewFloat(1)}}),
					},
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := participle.MustBuild(
				&Lambda{},
				participle.Lexer(lexer.MustStateful(lexRules(), lexer.InitialState(lexerFunc))),
				participle.Elide(TokenWhitespace),
			)

			res := &Lambda{}
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
