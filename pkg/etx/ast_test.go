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

func TestAST_Parsing(t *testing.T) {
	type args struct {
		Input string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    *AST
	}{
		{
			name: "One Attribute - no value",
			args: args{
				Input: `
foo`[1:],
			},
			wantErr: false,
			want: &AST{
				Items: []*Item{
					{
						Attribute: &Attribute{
							Pos: Position{Offset: 0, Line: 1, Column: 1},
							Key: "foo",
						},
					},
				},
			},
		},
		{
			name: "One Attribute - set value",
			args: args{
				Input: `
foo = 1`[1:],
			},
			wantErr: false,
			want: &AST{
				Items: []*Item{
					{
						Attribute: &Attribute{
							Pos:   Position{Offset: 0, Line: 1, Column: 1},
							Key:   "foo",
							Value: testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(1)}}),
						},
					},
				},
			},
		},
		{
			name: "Two Attributes - no value",
			args: args{
				Input: `
foo
bar`[1:],
			},
			wantErr: false,
			want: &AST{
				Items: []*Item{
					{
						Attribute: &Attribute{
							Pos: Position{Offset: 0, Line: 1, Column: 1},
							Key: "foo",
						},
					},
					{
						Attribute: &Attribute{
							Pos: Position{Offset: 4, Line: 2, Column: 1},
							Key: "bar",
						},
					},
				},
			},
		},
		{
			name: "Two Attributes - set values",
			args: args{
				Input: `
foo = 1
bar = 2`[1:],
			},
			wantErr: false,
			want: &AST{
				Items: []*Item{
					{
						Attribute: &Attribute{
							Pos:   Position{Offset: 0, Line: 1, Column: 1},
							Key:   "foo",
							Value: testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(1)}}),
						},
					},
					{
						Attribute: &Attribute{
							Pos:   Position{Offset: 8, Line: 2, Column: 1},
							Key:   "bar",
							Value: testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(2)}}),
						},
					},
				},
			},
		},
		{
			name: "One Decl - set value",
			args: args{
				Input: `
val foo = 1`[1:],
			},
			wantErr: false,
			want: &AST{
				Items: []*Item{
					{
						Decl: &Decl{
							Pos:      Position{Line: 1, Column: 1},
							DeclType: "val",
							Label:    "foo",
							Type:     nil,
							Value:    testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(1)}}),
						},
					},
				},
			},
		},
		{
			name: "One Attribute and One Decl",
			args: args{
				Input: `
foo = 1
val bar = 2`[1:],
			},
			wantErr: false,
			want: &AST{
				Items: []*Item{
					{
						Attribute: &Attribute{
							Pos:   Position{Offset: 0, Line: 1, Column: 1},
							Key:   "foo",
							Value: testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(1)}}),
						},
					},
					{
						Decl: &Decl{
							Pos:      Position{Offset: 8, Line: 2, Column: 1},
							DeclType: "val",
							Label:    "bar",
							Type:     nil,
							Value:    testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(2)}}),
						},
					},
				},
			},
		},
		{
			name: "One Decl - set value",
			args: args{
				Input: `
val foo = 1`[1:],
			},
			wantErr: false,
			want: &AST{
				Items: []*Item{
					{
						Decl: &Decl{
							Pos:      Position{Line: 1, Column: 1},
							DeclType: "val",
							Label:    "foo",
							Type:     nil,
							Value:    testBuildExprTree[*Expr](t, &Value{Number: &Number{big.NewFloat(1)}}),
						},
					},
				},
			},
		},
		{
			name: "One empty Func",
			args: args{
				Input: `def foo() {}`,
			},
			wantErr: false,
			want: &AST{
				Items: []*Item{
					{
						Func: &Func{
							Label: "foo",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.MustStateful(lexRules())

			// parser := participle.MustBuild(&AST{}, participle.Lexer(lexer.MustStateful(lexRules())))
			parser := participle.MustBuild(&AST{}, participle.Lexer(l))

			res := &AST{}
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
