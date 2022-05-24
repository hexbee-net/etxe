package etx

import (
	"math/big"
	"testing"
)

func TestDecl_Parsing(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Decl
	}{
		{
			name:    "Const declaration - no value, not type",
			input:   `const foo`,
			wantErr: false,
			want: &Decl{
				ASTNode:  ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				DeclType: "const",
				Label:    "foo",
			},
		},
		{
			name:    "Const declaration - no value, type - no space",
			input:   `const foo:number`,
			wantErr: false,
			want: &Decl{
				ASTNode:  ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				DeclType: "const",
				Label:    "foo",
				Type: &ParameterType{
					ASTNode: ASTNode{Pos: Position{Offset: 10, Line: 1, Column: 11}},
					Ident:   &Ident{Parts: []string{"number"}},
				},
			},
		},
		{
			name:    "Const declaration - no value, type - post space",
			input:   `const foo: number`,
			wantErr: false,
			want: &Decl{
				ASTNode:  ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				DeclType: "const",
				Label:    "foo",
				Type: &ParameterType{
					ASTNode: ASTNode{Pos: Position{Offset: 11, Line: 1, Column: 12}},
					Ident:   &Ident{Parts: []string{"number"}},
				},
			},
		},
		{
			name:    "Const declaration - no value, type - pre and post space",
			input:   `const foo : number`,
			wantErr: false,
			want: &Decl{
				ASTNode:  ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				DeclType: "const",
				Label:    "foo",
				Type: &ParameterType{
					ASTNode: ASTNode{Pos: Position{Offset: 12, Line: 1, Column: 13}},
					Ident:   &Ident{Parts: []string{"number"}},
				},
			},
		},
		{
			name:    "Const declaration - value, type",
			input:   `const foo: number = 1`,
			wantErr: false,
			want: &Decl{
				ASTNode:  ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				DeclType: "const",
				Label:    "foo",
				Type: &ParameterType{
					ASTNode: ASTNode{Pos: Position{Offset: 11, Line: 1, Column: 12}},
					Ident:   &Ident{Parts: []string{"number"}},
				},
				Value: testBuildExprTree[*Expr](t, &Value{
					ASTNode: ASTNode{Pos: Position{Offset: 20, Line: 1, Column: 21}},
					Number:  &ValueNumber{big.NewFloat(1), "1"},
				}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParser(t, tt.input, tt.want, tt.wantErr, true)
		})
	}
}
