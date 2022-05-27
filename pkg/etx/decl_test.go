package etx

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestDecl_Clone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *Decl
		want  *Decl
	}{
		{
			name:  "Nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "Empty",
			input: &Decl{},
			want:  &Decl{},
		},
		{
			name: "ASTNode",
			input: &Decl{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
			want: &Decl{
				ASTNode: ASTNode{Pos: Position{Offset: 1, Line: 2, Column: 3}},
			},
		},
		{
			name: "DeclType",
			input: &Decl{
				DeclType: "val",
			},
			want: &Decl{
				DeclType: "val",
			},
		},
		{
			name: "Label",
			input: &Decl{
				Label: "foo",
			},
			want: &Decl{
				Label: "foo",
			},
		},
		{
			name: "Type",
			input: &Decl{
				Type: &ParameterType{Ident: &Ident{Parts: []string{"number"}}},
			},
			want: &Decl{
				Type: &ParameterType{Ident: &Ident{Parts: []string{"number"}}},
			},
		},
		{
			name: "Value",
			input: &Decl{
				Value: testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
			},
			want: &Decl{
				Value: testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCloner[*Decl](t, tt.want, tt.input.Clone())
		})
	}
}

func TestDecl_Children(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *Decl
		want  []Node
	}{
		{
			name:  "Empty",
			input: &Decl{},
			want:  nil,
		},
		{
			name: "DeclType",
			input: &Decl{
				DeclType: "val",
			},
			want: nil,
		},
		{
			name: "Label",
			input: &Decl{
				Label: "foo",
			},
			want: nil,
		},
		{
			name: "Type",
			input: &Decl{
				Type: &ParameterType{Ident: &Ident{Parts: []string{"number"}}},
			},
			want: []Node{
				&ParameterType{Ident: &Ident{Parts: []string{"number"}}},
			},
		},
		{
			name: "Value",
			input: &Decl{
				Value: testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
			},
			want: []Node{
				testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.Children())
		})
	}
}

func TestDecl_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     *Decl
		wantPanic bool
		want      string
	}{
		{
			name:      "Nil",
			input:     nil,
			wantPanic: true,
		},
		{
			name:  "Empty",
			input: &Decl{},
			want:  "",
		},
		{
			name: "Label - no type - no value",
			input: &Decl{
				DeclType: "val",
				Label:    "foo",
			},
			want: "val foo",
		},
		{
			name: "Label - type - no value",
			input: &Decl{
				DeclType: "val",
				Label:    "foo",
				Type:     &ParameterType{Ident: &Ident{Parts: []string{"number"}}},
			},
			want: "val foo: number",
		},
		{
			name: "Label - no type - value",
			input: &Decl{
				DeclType: "val",
				Label:    "foo",
				Value:    testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
			},
			want: "val foo = 1",
		},
		{
			name: "Label - type - value",
			input: &Decl{
				DeclType: "val",
				Label:    "foo",
				Type:     &ParameterType{Ident: &Ident{Parts: []string{"number"}}},
				Value:    testBuildExprTree[*Expr](t, &Value{Number: &ValueNumber{big.NewFloat(1), "1"}}),
			},
			want: "val foo: number = 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringer(t, tt.wantPanic, tt.want, tt.input)
		})
	}
}
