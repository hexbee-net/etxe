package etx

import (
	"math/big"
	"testing"
)

func TestType_Parsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		want    *Type
	}{
		{
			name:    "Enum - Empty",
			input:   `type foo enum {}`,
			wantErr: false,
			want: &Type{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
				Enum: &TypeEnum{
					ASTNode: ASTNode{Pos: Position{Offset: 15, Line: 1, Column: 16}},
				},
				Object: nil,
			},
		},
		{
			name: "Enum - One valid value",
			input: `
type foo enum {
  bar: 1
}`[1:],
			wantErr: false,
			want: &Type{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
				Enum: &TypeEnum{
					ASTNode: ASTNode{Pos: Position{Offset: 15, Line: 1, Column: 16}},
					Items: []*TypeEnumItem{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 18, Line: 2, Column: 3}},
							Label:   "bar",
							Value: *testBuildExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 23, Line: 2, Column: 8}},
								Number:  &ValueNumber{big.NewFloat(1), "1"},
							}),
						},
					},
				},
				Object: nil,
			},
		},
		{
			name: "Enum - Two valid values",
			input: `
type foo enum {
  bar: 1
  baz: 2
}`[1:],
			wantErr: false,
			want: &Type{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
				Enum: &TypeEnum{
					ASTNode: ASTNode{Pos: Position{Offset: 15, Line: 1, Column: 16}},
					Items: []*TypeEnumItem{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 18, Line: 2, Column: 3}},
							Label:   "bar",
							Value: *testBuildExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 23, Line: 2, Column: 8}},
								Number:  &ValueNumber{big.NewFloat(1), "1"},
							}),
						},
						{
							ASTNode: ASTNode{Pos: Position{Offset: 27, Line: 3, Column: 3}},
							Label:   "baz",
							Value: *testBuildExprTree[*Expr](t, &Value{
								ASTNode: ASTNode{Pos: Position{Offset: 32, Line: 3, Column: 8}},
								Number:  &ValueNumber{big.NewFloat(2), "2"},
							}),
						},
					},
				},
				Object: nil,
			},
		},

		{
			name:    "Object - Empty",
			input:   `type foo object {}`,
			wantErr: false,
			want: &Type{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
				Enum:    nil,
				Object: &TypeObject{
					ASTNode: ASTNode{Pos: Position{Offset: 17, Line: 1, Column: 18}},
				},
			},
		},
		{
			name: "Object - One valid declaration",
			input: `
type foo object {
	foo: number
}`[1:],
			wantErr: false,
			want: &Type{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
				Enum:    nil,
				Object: &TypeObject{
					ASTNode: ASTNode{Pos: Position{Offset: 17, Line: 1, Column: 18}},
					Items: []*TypeObjectItem{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 19, Line: 2, Column: 2}},
							Label:   "foo",
							Type: ParameterType{
								ASTNode: ASTNode{Pos: Position{Offset: 24, Line: 2, Column: 7}},
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
			name: "Object - Two valid declarations",
			input: `
type foo object {
	foo: number
    bar: bool
}`[1:],
			wantErr: false,
			want: &Type{
				ASTNode: ASTNode{Pos: Position{Offset: 0, Line: 1, Column: 1}},
				Label:   "foo",
				Enum:    nil,
				Object: &TypeObject{
					ASTNode: ASTNode{Pos: Position{Offset: 17, Line: 1, Column: 18}},
					Items: []*TypeObjectItem{
						{
							ASTNode: ASTNode{Pos: Position{Offset: 19, Line: 2, Column: 2}},
							Label:   "foo",
							Type: ParameterType{
								ASTNode: ASTNode{Pos: Position{Offset: 24, Line: 2, Column: 7}},
								Ident: &Ident{
									Parts: []string{
										"number",
									},
								},
							},
						},
						{
							ASTNode: ASTNode{Pos: Position{Offset: 35, Line: 3, Column: 5}},
							Label:   "bar",
							Type: ParameterType{
								ASTNode: ASTNode{Pos: Position{Offset: 40, Line: 3, Column: 10}},
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
			testParser(t, tt.input, tt.want, tt.wantErr, true)
		})
	}
}
