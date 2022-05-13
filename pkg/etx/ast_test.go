package etx

import (
	"math/big"
	"testing"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecl_Parsing(t *testing.T) {
	type args struct {
		Input string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    *Decl
	}{
		{
			name: "",
			args: args{
				Input: `const foo`,
			},
			wantErr: false,
			want: &Decl{
				Pos:      Position{Line: 1, Column: 1},
				DeclType: "const",
				Label:    "foo",
			},
		},
		{
			name: "",
			args: args{
				Input: `const foo:number`,
			},
			wantErr: false,
			want: &Decl{
				Pos:      Position{Line: 1, Column: 1},
				DeclType: "const",
				Label:    "foo",
				Type:     "number",
			},
		},
		{
			name: "",
			args: args{
				Input: `const foo: number`,
			},
			wantErr: false,
			want: &Decl{
				Pos:      Position{Line: 1, Column: 1},
				DeclType: "const",
				Label:    "foo",
				Type:     "number",
			},
		},
		{
			name: "",
			args: args{
				Input: `const foo : number`,
			},
			wantErr: false,
			want: &Decl{
				Pos:      Position{Line: 1, Column: 1},
				DeclType: "const",
				Label:    "foo",
				Type:     "number",
			},
		},
		{
			name: "",
			args: args{
				Input: `const foo: number = 1`,
			},
			wantErr: false,
			want: &Decl{
				Pos:      Position{Line: 1, Column: 1},
				DeclType: "const",
				Label:    "foo",
				Type:     "number",
				Value:    &Value{Number: &Number{big.NewFloat(1)}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := participle.MustBuild(&Decl{}, participle.Lexer(lexer.MustStateful(lexRules())))

			res := &Decl{}
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
						// Pos:   Position{Offset: 1, Line: 1, Column: 2},
						Label: "x",
						Type: &ParameterType{
							Ident: "number",
						},
					},
				},
				Expr: testBuildExprTree[*Expr](t,
					&ExprAdditive{
						Left:  testBuildExprTree[*ExprMultiplicative](t, &Value{Ident: testValPtr(t, "x")}),
						Op:    TokenOpPlus,
						Right: testBuildExprTree[*ExprAdditive](t, &Value{Number: &Number{big.NewFloat(1)}}),
					},
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := participle.MustBuild(&Lambda{}, participle.Lexer(lexer.MustStateful(lexRules(), lexer.InitialState(lexerFunc))))

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
