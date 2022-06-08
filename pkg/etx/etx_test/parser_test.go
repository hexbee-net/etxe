package etx_test

import (
	"math/big"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hexbee-net/etxe/pkg/etx"
)

func testJoinNodes[T any](t *testing.T, sep T, elems ...[]T) []T {
	t.Helper()

	switch len(elems) {
	case 0:
		return []T{}
	case 1:
		return elems[0]
	}

	n := len(elems)
	for i := 0; i < len(elems); i++ {
		n += len(elems[i])
	}

	res := make([]T, 0, n)
	res = append(res, elems[0]...)

	for _, t := range elems[1:] {
		res = append(res, sep)
		res = append(res, t...)
	}

	return res
}

func TestParse(t *testing.T) {
	testJoinFileNodes := func(t *testing.T, elems ...[]*etx.RootItem) []*etx.RootItem {
		return testJoinNodes[*etx.RootItem](t, &etx.RootItem{EmptyLine: "\n"}, elems...)
	}

	tests := []struct {
		name    string
		source  string
		want    *etx.AST
		wantErr bool
	}{
		{
			name:   "Comments",
			source: "fixtures/comments.etx",
			want: &etx.AST{
				Items: testJoinFileNodes(t,
					[]*etx.RootItem{
						{Comment: &etx.Comment{
							SingleLine: []string{"// One line double-dashed"},
						}},
					},
					[]*etx.RootItem{
						{Comment: &etx.Comment{
							SingleLine: []string{
								"// Two lines",
								"// double-dashed",
							},
						}},
					},
					[]*etx.RootItem{
						{Comment: &etx.Comment{
							SingleLine: []string{"# One line hashtag"},
						}},
					},
					[]*etx.RootItem{
						{Comment: &etx.Comment{
							SingleLine: []string{
								"# Two lines",
								"# double-dashed",
							},
						}},
					},
					[]*etx.RootItem{
						{Comment: &etx.Comment{
							SingleLine: []string{
								"#  Two lines",
								"// mixed",
							},
						}},
					},
					[]*etx.RootItem{
						{Comment: &etx.Comment{
							Multiline: "/* One multiline on one line */\n",
						}},
					},
					[]*etx.RootItem{
						{Comment: &etx.Comment{
							Multiline: "/* One multiline on\n   two line */\n",
						}},
					},
				),
			},
		},
		{
			name:   "Blocks",
			source: "fixtures/blocks.etx",
			want: &etx.AST{
				Items: testJoinFileNodes(t,
					[]*etx.RootItem{
						{Block: &etx.Block{
							Name: "block_name",
						}},
					},
					[]*etx.RootItem{
						{Block: &etx.Block{
							Name:   "block_name",
							Labels: []string{"label-1", "label-2"},
						}},
					},
					[]*etx.RootItem{
						{Comment: &etx.Comment{
							SingleLine: []string{"# block comment"},
						}},
						{Block: &etx.Block{
							Name: "block_name",
						}},
					},
					[]*etx.RootItem{
						{Comment: &etx.Comment{
							SingleLine: []string{"// block comment"},
						}},
						{Block: &etx.Block{
							Name: "block_name",
						}},
					},
					[]*etx.RootItem{
						{Comment: &etx.Comment{
							Multiline: "/* block comment */\n",
						}},
						{Block: &etx.Block{
							Name: "block_name",
						}},
					},
					[]*etx.RootItem{
						{Comment: &etx.Comment{
							SingleLine: []string{
								"# block comment",
								"# on two lines",
							},
						}},
						{Block: &etx.Block{
							Name: "block_name",
						}},
					},
					[]*etx.RootItem{
						{Comment: &etx.Comment{
							SingleLine: []string{
								"// block comment",
								"// on two lines",
							},
						}},
						{Block: &etx.Block{
							Name: "block_name",
						}},
					},
					[]*etx.RootItem{
						{Comment: &etx.Comment{
							Multiline: "/* block comment */\n",
						}},
						{Comment: &etx.Comment{
							Multiline: "/* on two lines  */\n",
						}},
						{Block: &etx.Block{
							Name: "block_name",
						}},
					},
					[]*etx.RootItem{
						{Comment: &etx.Comment{
							Multiline: "/* block comment\n   on two lines   */\n",
						}},
						{Block: &etx.Block{
							Name: "block_name",
						}},
					},

					[]*etx.RootItem{
						{Block: &etx.Block{
							Name: "block_attributes",
							Body: testJoinNodes[*etx.BlockItem](t,
								&etx.BlockItem{EmptyLine: "\n"},
								[]*etx.BlockItem{
									{Attribute: &etx.Attribute{
										Key: "attr1",
									}},
								},
								[]*etx.BlockItem{
									{Attribute: &etx.Attribute{
										Key: "attr2",
									}},
								},
								[]*etx.BlockItem{
									{Attribute: &etx.Attribute{
										Key:   "attr",
										Value: etx.BuildTestExprTree[*etx.Expr](t, &etx.Value{Null: true}),
									}},
								},
								[]*etx.BlockItem{
									{Attribute: &etx.Attribute{
										Key:   "attr",
										Value: etx.BuildTestExprTree[*etx.Expr](t, &etx.Value{Bool: &etx.ValueBool{Value: true}}),
									}},
									{Attribute: &etx.Attribute{
										Key:   "attr",
										Value: etx.BuildTestExprTree[*etx.Expr](t, &etx.Value{Bool: &etx.ValueBool{Value: false}}),
									}},
								},
								[]*etx.BlockItem{
									{Attribute: &etx.Attribute{
										Key:   "attr",
										Value: etx.BuildTestExprTree[*etx.Expr](t, &etx.Value{Number: &etx.ValueNumber{Value: big.NewFloat(1), Source: "1"}}),
									}},
								},
								[]*etx.BlockItem{
									{Attribute: &etx.Attribute{
										Key: "attr",
										Value: etx.BuildTestExprTree[*etx.Expr](t, &etx.Value{Str: &etx.ValueString{
											Fragment: []*etx.StringFragment{
												{Text: "string"},
											},
										}}),
									}},
									{Attribute: &etx.Attribute{
										Key: "attr",
										Value: etx.BuildTestExprTree[*etx.Expr](t, &etx.Value{Str: &etx.ValueString{
											Fragment: []*etx.StringFragment{
												{Expr: etx.BuildTestExprTree[*etx.Expr](t, &etx.Ident{Parts: []string{"expr"}})},
												{Text: "-string"},
											},
										}}),
									}},
								},
								[]*etx.BlockItem{
									{Attribute: &etx.Attribute{
										Key:   "attr",
										Value: etx.BuildTestExprTree[*etx.Expr](t, &etx.Ident{Parts: []string{"ident"}}),
									}},
								},
								[]*etx.BlockItem{
									{Attribute: &etx.Attribute{
										Key: "attr",
										Value: etx.BuildTestExprTree[*etx.Expr](t, &etx.Value{Heredoc: &etx.Heredoc{
											Delimiter: etx.HeredocDelimiter{
												LeadingTabs: false,
												Delimiter:   "EOF",
											},
											Fragments: []*etx.HeredocFragment{
												{
													Text: "heredoc\n",
												},
											},
										}}),
									}},
								},
								[]*etx.BlockItem{
									{Attribute: &etx.Attribute{
										Key: "attr",
										Value: etx.BuildTestExprTree[*etx.Expr](t, &etx.Value{Heredoc: &etx.Heredoc{
											Delimiter: etx.HeredocDelimiter{
												LeadingTabs: true,
												Delimiter:   "EOF",
											},
											Fragments: []*etx.HeredocFragment{
												{
													Text: "heredoc leading tabs\n",
												},
											},
										}}),
									}},
								},
								[]*etx.BlockItem{
									{Attribute: &etx.Attribute{
										Key: "attr",
										Value: etx.BuildTestExprTree[*etx.Expr](t, &etx.Value{
											List: &etx.ValueList{
												Items: []*etx.ListItem{
													{Value: etx.BuildTestExprTree[*etx.Expr](t, &etx.Value{Number: &etx.ValueNumber{Value: big.NewFloat(1), Source: "1"}})},
													{Value: etx.BuildTestExprTree[*etx.Expr](t, &etx.Value{Number: &etx.ValueNumber{Value: big.NewFloat(2), Source: "2"}})},
												},
											},
										}),
									}},
									{Attribute: &etx.Attribute{
										Key: "attr",
										Value: etx.BuildTestExprTree[*etx.Expr](t, &etx.Value{
											List: &etx.ValueList{
												Items: []*etx.ListItem{
													{Value: etx.BuildTestExprTree[*etx.Expr](t, &etx.Value{Number: &etx.ValueNumber{Value: big.NewFloat(1), Source: "1"}})},
													{Value: etx.BuildTestExprTree[*etx.Expr](t, &etx.Value{Number: &etx.ValueNumber{Value: big.NewFloat(2), Source: "2"}})},
												},
											},
										}),
									}},
								},
								[]*etx.BlockItem{
									{Attribute: &etx.Attribute{
										Key: "attr",
										Value: etx.BuildTestExprTree[*etx.Expr](t, &etx.Value{
											Map: &etx.ValueMap{
												Items: []*etx.MapItem{
													{
														Key: &etx.MapKey{
															Ident: &etx.Ident{Parts: []string{"a"}},
														},
														Value: etx.BuildTestExprTree[*etx.Expr](t, &etx.Value{
															Number: &etx.ValueNumber{Value: big.NewFloat(1), Source: `1`},
														}),
													},
													{
														Key: &etx.MapKey{
															Str: &etx.ValueString{
																Fragment: []*etx.StringFragment{{Text: "b"}},
															},
														},
														Value: etx.BuildTestExprTree[*etx.Expr](t, &etx.Value{
															Number: &etx.ValueNumber{Value: big.NewFloat(2), Source: `2`},
														}),
													},
												},
											},
										}),
									}},
									{Attribute: &etx.Attribute{
										Key: "attr",
										Value: etx.BuildTestExprTree[*etx.Expr](t, &etx.Value{
											Map: &etx.ValueMap{
												Items: []*etx.MapItem{
													{
														Key: &etx.MapKey{
															Ident: &etx.Ident{
																Parts: []string{"a"},
															},
														},
														Value: etx.BuildTestExprTree[*etx.Expr](t, &etx.Value{
															Number: &etx.ValueNumber{Value: big.NewFloat(1), Source: `1`},
														}),
													},
													{
														Key: &etx.MapKey{
															Str: &etx.ValueString{
																Fragment: []*etx.StringFragment{{Text: "b"}},
															},
														},
														Value: etx.BuildTestExprTree[*etx.Expr](t, &etx.Value{
															Number: &etx.ValueNumber{Value: big.NewFloat(2), Source: `2`},
														}),
													},
												},
											},
										}),
									}},
								},
							),
						}},
					},

					[]*etx.RootItem{
						{Block: &etx.Block{
							Name: "block_sub_block",
							Body: []*etx.BlockItem{
								{
									Block: &etx.Block{
										Name: "sub_block",
									},
								},
							},
						}},
					},
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, err := os.Open(tt.source)

			res, err := etx.Parse(reader)
			if tt.wantErr {
				assert.Error(t, err)
				return
			} else {
				require.NoError(t, err)

				l1 := len(tt.want.Items)
				l2 := len(res.Items)
				assert.Equal(t, l1, l2, "lengths should match")

				numberComparer := cmp.Comparer(func(x, y *etx.ValueNumber) bool {
					if x == nil && y == nil {
						return true
					}

					if x == nil || y == nil {
						return false
					}

					return x.Value.Cmp(y.Value) == 0
				})

				posComparer := cmp.Comparer(func(x, y etx.ASTNode) bool {
					return true
				})

				if !cmp.Equal(tt.want, res, numberComparer, posComparer) {
					assert.Fail(t, "Not equal -want +res", cmp.Diff(tt.want, res, numberComparer, posComparer))
				}
			}
		})
	}
}
