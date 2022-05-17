package etx

import (
	"math/big"
	"testing"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/repr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFunc_Parsing(t *testing.T) {
	type args struct {
		Input string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    *Func
	}{
		{
			name: "Empty body, no params, no return",
			args: args{
				Input: `def foo() {}`,
			},
			wantErr: false,
			want: &Func{
				Label: "foo",
			},
		},
		{
			name: "Empty body, one ident param, no return",
			args: args{
				Input: `def foo(bar: bool) {}`,
			},
			wantErr: false,
			want: &Func{
				Label: "foo",
				Parameters: []*FuncParameter{
					{
						Label: "bar",
						Type: &ParameterType{
							Ident: &Ident{Parts: []string{"bool"}},
						},
					},
				},
			},
		},
		{
			name: "Empty body, two ident param, no return",
			args: args{
				Input: `def foo(bar: bool, baz: number) {}`,
			},
			wantErr: false,
			want: &Func{
				Label: "foo",
				Parameters: []*FuncParameter{
					{
						Label: "bar",
						Type: &ParameterType{
							Ident: &Ident{Parts: []string{"bool"}},
						},
					},
					{
						Label: "baz",
						Type: &ParameterType{
							Ident: &Ident{Parts: []string{"number"}},
						},
					},
				},
			},
		},
		{
			name: "Empty body, no params, ident return",
			args: args{
				Input: `def foo() bool {}`,
			},
			wantErr: false,
			want: &Func{
				Label: "foo",
				Return: &ParameterType{
					Ident: &Ident{Parts: []string{"bool"}},
				},
			},
		},

		{
			name: "Empty body, one func param, no return",
			args: args{
				Input: `def foo(bar: (int) -> bool) {}`,
			},
			wantErr: false,
			want: &Func{
				Label: "foo",
				Parameters: []*FuncParameter{
					{
						Label: "bar",
						Type: &ParameterType{
							Func: &FuncSignature{
								Parameters: []*ParameterType{
									{Ident: &Ident{Parts: []string{"int"}}},
								},
								Return: &ParameterType{
									Ident: &Ident{Parts: []string{"bool"}},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Empty body, two func params, no return",
			args: args{
				Input: `def foo(bar: (int) -> bool, baz: (bool) -> int) {}`,
			},
			wantErr: false,
			want: &Func{
				Label: "foo",
				Parameters: []*FuncParameter{
					{
						Label: "bar",
						Type: &ParameterType{
							Func: &FuncSignature{
								Parameters: []*ParameterType{
									{Ident: &Ident{Parts: []string{"int"}}},
								},
								Return: &ParameterType{
									Ident: &Ident{Parts: []string{"bool"}},
								},
							},
						},
					},
					{
						Label: "baz",
						Type: &ParameterType{
							Func: &FuncSignature{
								Parameters: []*ParameterType{
									{Ident: &Ident{Parts: []string{"bool"}}},
								},
								Return: &ParameterType{
									Ident: &Ident{Parts: []string{"int"}},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Empty body, no params, func return",
			args: args{
				Input: `def foo() (int) -> bool {}`,
			},
			wantErr: false,
			want: &Func{
				Label: "foo",
				Return: &ParameterType{
					Func: &FuncSignature{
						Parameters: []*ParameterType{
							{Ident: &Ident{Parts: []string{"int"}}},
						},
						Return: &ParameterType{
							Ident: &Ident{Parts: []string{"bool"}},
						},
					},
				},
			},
		},

		{
			name: "One Expr statement, no params, no return",
			args: args{
				Input: `
def foo() {
	a
}`[1:],
			},
			wantErr: false,
			want: &Func{
				Label: "foo",
				Body: []*FuncStatement{
					{
						Expr: testBuildExprTree[*Expr](t, &Value{Ident: &Ident{Parts: []string{"a"}}}),
					},
				},
			},
		},
		{
			name: "One val Decl statement, no params, no return",
			args: args{
				Input: `
def foo() {
	val a: number = 1
}`[1:],
			},
			wantErr: false,
			want: &Func{
				Label: "foo",
				Body: []*FuncStatement{
					{
						Decl: &FuncDecl{
							DeclType: "val",
							Label:    "a",
							Type:     "number",
							Value:    testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(1)}}),
						},
					},
				},
			},
		},
		{
			name: "Val decl and return expr statements, no params, no return",
			args: args{
				Input: `
def foo() {
	val a: number = 1
	a
}`[1:],
			},
			wantErr: false,
			want: &Func{
				Label: "foo",
				Body: []*FuncStatement{
					{
						Decl: &FuncDecl{
							DeclType: "val",
							Label:    "a",
							Type:     "number",
							Value:    testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(1)}}),
						},
					},
					{
						Expr: testBuildExprTree[*Expr](t, &Value{Ident: &Ident{Parts: []string{"a"}}}),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := participle.MustBuild(
				&Func{},
				participle.Lexer(lexer.MustStateful(lexRules())),
				participle.Elide(TokenWhitespace),
			)

			res := &Func{}
			err := parser.ParseString("", tt.args.Input, res)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if !assert.Equal(t, tt.want, res) {
				repr.Println(res)
			}
		})
	}
}
