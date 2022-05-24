package etx

import (
	"math/big"
	"testing"
)

func TestAST_Parsing(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *AST
	}{
		{
			name: "One Attribute - no value",
			input: `
foo`[1:],
			wantErr: false,
			want: &AST{
				Items: []*RootItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Attribute: &Attribute{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Key:     "foo",
						},
					},
				},
			},
		},
		{
			name: "One Attribute - set value",
			input: `
foo = 1`[1:],
			wantErr: false,
			want: &AST{
				Items: []*RootItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Attribute: &Attribute{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Key:     "foo",
							Value: testBuildExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
								Number:  &ValueNumber{big.NewFloat(1), "1"},
							}),
						},
					},
				},
			},
		},
		{
			name: "Two Attributes - no value",
			input: `
foo
bar`[1:],
			wantErr: false,
			want: &AST{
				Items: []*RootItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Attribute: &Attribute{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Key:     "foo",
						},
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 2, Column: 1}},
						Attribute: &Attribute{
							ASTNode: ASTNode{Pos: Position{Offset: 4, Line: 2, Column: 1}},
							Key:     "bar",
						},
					},
				},
			},
		},
		{
			name: "Two Attributes - set values",
			input: `
foo = 1
bar = 2`[1:],
			wantErr: false,
			want: &AST{
				Items: []*RootItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Attribute: &Attribute{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Key:     "foo",
							Value: testBuildExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
								Number:  &ValueNumber{big.NewFloat(1), "1"},
							}),
						},
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 2, Column: 1}},
						Attribute: &Attribute{
							ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 2, Column: 1}},
							Key:     "bar",
							Value: testBuildExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 14, Line: 2, Column: 7}},
								Number:  &ValueNumber{big.NewFloat(2), "2"},
							}),
						},
					},
				},
			},
		},
		{
			name: "One Decl - set value",
			input: `
val foo = 1`[1:],
			wantErr: false,
			want: &AST{
				Items: []*RootItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Decl: &Decl{
							ASTNode:  ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							DeclType: "val",
							Label:    "foo",
							Type:     nil,
							Value: testBuildExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 10, Line: 1, Column: 11}},
								Number:  &ValueNumber{big.NewFloat(1), "1"},
							}),
						},
					},
				},
			},
		},
		{
			name: "One Attribute and One Decl",
			input: `
foo = 1
val bar = 2`[1:],
			wantErr: false,
			want: &AST{
				Items: []*RootItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Attribute: &Attribute{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Key:     "foo",
							Value: testBuildExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 6, Line: 1, Column: 7}},
								Number:  &ValueNumber{big.NewFloat(1), "1"},
							}),
						},
					},
					{
						ASTNode: ASTNode{Pos: Position{Offset: 8, Line: 2, Column: 1}},
						Decl: &Decl{
							ASTNode:  ASTNode{Pos: Position{Offset: 8, Line: 2, Column: 1}},
							DeclType: "val",
							Label:    "bar",
							Type:     nil,
							Value: testBuildExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 18, Line: 2, Column: 11}},
								Number:  &ValueNumber{big.NewFloat(2), "2"},
							}),
						},
					},
				},
			},
		},
		{
			name: "One Decl - set value",
			input: `
val foo = 1`[1:],
			wantErr: false,
			want: &AST{
				Items: []*RootItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Decl: &Decl{
							ASTNode:  ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							DeclType: "val",
							Label:    "foo",
							Type:     nil,
							Value: testBuildExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 10, Line: 1, Column: 11}},
								Number:  &ValueNumber{big.NewFloat(1), "1"},
							}),
						},
					},
				},
			},
		},
		{
			name:    "One empty Func",
			input:   `def foo() {}`,
			wantErr: false,
			want: &AST{
				Items: []*RootItem{
					{
						ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
						Func: &Func{
							ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
							Label:   "foo",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, true)
		})
	}
}
