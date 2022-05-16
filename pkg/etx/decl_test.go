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
			name: "Const declaration - no value, not type",
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
			name: "Const declaration - no value, type - no space",
			args: args{
				Input: `const foo:number`,
			},
			wantErr: false,
			want: &Decl{
				Pos:      Position{Line: 1, Column: 1},
				DeclType: "const",
				Label:    "foo",
				Type:     &ParameterType{Ident: &Ident{Parts: []string{"number"}}},
			},
		},
		{
			name: "Const declaration - no value, type - post space",
			args: args{
				Input: `const foo: number`,
			},
			wantErr: false,
			want: &Decl{
				Pos:      Position{Line: 1, Column: 1},
				DeclType: "const",
				Label:    "foo",
				Type:     &ParameterType{Ident: &Ident{Parts: []string{"number"}}},
			},
		},
		{
			name: "Const declaration - no value, type - pre and post space",
			args: args{
				Input: `const foo : number`,
			},
			wantErr: false,
			want: &Decl{
				Pos:      Position{Line: 1, Column: 1},
				DeclType: "const",
				Label:    "foo",
				Type:     &ParameterType{Ident: &Ident{Parts: []string{"number"}}},
			},
		},
		{
			name: "Const declaration - value, type",
			args: args{
				Input: `const foo: number = 1`,
			},
			wantErr: false,
			want: &Decl{
				Pos:      Position{Line: 1, Column: 1},
				DeclType: "const",
				Label:    "foo",
				Type:     &ParameterType{Ident: &Ident{Parts: []string{"number"}}},
				Value:    testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(1)}}),
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
