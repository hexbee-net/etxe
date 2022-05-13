package etx

import (
	"testing"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
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
							Ident: "bool",
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
							Ident: "bool",
						},
					},
					{
						Label: "baz",
						Type: &ParameterType{
							Ident: "number",
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
					Ident: "bool",
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
									{Ident: "int"},
								},
								Return: &ParameterType{
									Ident: "bool",
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
									{Ident: "int"},
								},
								Return: &ParameterType{
									Ident: "bool",
								},
							},
						},
					},
					{
						Label: "baz",
						Type: &ParameterType{
							Func: &FuncSignature{
								Parameters: []*ParameterType{
									{Ident: "bool"},
								},
								Return: &ParameterType{
									Ident: "int",
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
							{Ident: "int"},
						},
						Return: &ParameterType{
							Ident: "bool",
						},
					},
				},
			},
		},

		{
			name: "Ident body, no params, no return",
			args: args{
				Input: `def foo() {a}`,
			},
			wantErr: false,
			want: &Func{
				Label: "foo",
				Body: []*FuncExpr{
					{
						Todo: "a",
					},
				},
			},
		},

		// 		{
		// 			name: "wip",
		// 			args: args{
		// 				Input: `def foo(bar: bool, baz: number) bool {
		// 	body
		// }`,
		// 			},
		// 			wantErr: false,
		// 			want: &Func{
		// 				// Pos:   lexer.Position{Offset: 0, Line: 1, Column: 1},
		// 				Label: "foo",
		// 				Parameters: []*FuncParameter{
		// 					{
		// 						Label: "bar",
		// 						Type: &ParameterType{
		// 							Ident: "bool",
		// 						},
		// 					},
		// 					{
		// 						Label: "baz",
		// 						Type: &ParameterType{
		// 							Ident: "number",
		// 						},
		// 					},
		// 				},
		// 				Return: &ParameterType{
		// 					Ident: "bool",
		// 				},
		// 				Body: &FuncBody{
		// 					Todo: "body",
		// 				},
		// 			},
		// 		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := participle.MustBuild(&Func{}, participle.Lexer(lexer.MustStateful(lexRules())))

			res := &Func{}
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
