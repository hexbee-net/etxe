package etx

import (
	"math/big"
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
			name: "Enum - Empty",
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
			name: "Enum - One valid value",
			args: args{
				Input: `
type foo enum {
  bar: 1
}`[1:],
			},
			wantErr: false,
			want: &Type{
				Label: "foo",
				Enum: &Enum{
					Items: []EnumItem{
						{
							Label: "bar",
							Value: *testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(1)}}),
						},
					},
				},
				Object: nil,
			},
		},
		{
			name: "Enum - One invalid value",
			args: args{
				Input: `
type foo enum {
  bar:
}`[1:],
			},
			wantErr: true,
			want: &Type{
				Label: "foo",
				Enum: &Enum{
					Items: []EnumItem{
						{
							Label: "bar",
						},
					},
				},
				Object: nil,
			},
		},
		{
			name: "Enum - Two valid values",
			args: args{
				Input: `
type foo enum {
  bar: 1
  baz: 2
}`[1:],
			},
			wantErr: false,
			want: &Type{
				Label: "foo",
				Enum: &Enum{
					Items: []EnumItem{
						{
							Label: "bar",
							Value: *testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(1)}}),
						},
						{
							Label: "baz",
							Value: *testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(2)}}),
						},
					},
				},
				Object: nil,
			},
		},

		{
			name: "Object - Empty",
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
		{
			name: "Object - One valid declaration",
			args: args{
				Input: `
type foo object {
	foo: number
}`[1:],
			},
			wantErr: false,
			want: &Type{
				Label: "foo",
				Enum:  nil,
				Object: &Object{
					Items: []ObjectItem{
						{
							Label: "foo",
							Type: ParameterType{
								Ident: &Ident{
									Parts: []string{
										"number",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Object - One invalid declaration",
			args: args{
				Input: `
type foo object {
	foo:
}`[1:],
			},
			wantErr: true,
			want: &Type{
				Label: "foo",
				Enum:  nil,
				Object: &Object{
					Items: []ObjectItem{
						{
							Label: "foo",
							Type:  ParameterType{},
						},
					},
				},
			},
		},
		{
			name: "Object - Two valid declarations",
			args: args{
				Input: `
type foo object {
	foo: number
    bar: bool
}`[1:],
			},
			wantErr: false,
			want: &Type{
				Label: "foo",
				Enum:  nil,
				Object: &Object{
					Items: []ObjectItem{
						{
							Label: "foo",
							Type: ParameterType{
								Ident: &Ident{
									Parts: []string{
										"number",
									},
								},
							},
						},
						{
							Label: "bar",
							Type: ParameterType{
								Ident: &Ident{
									Parts: []string{
										"bool",
									},
								},
							},
						},
					},
				},
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
