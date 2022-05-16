package etx

type Lambda struct {
	Parameters []*LambdaParameter `parser:"'(' (@@ ( Whitespace? ',' Whitespace? @@ )*)? ')' Whitespace? Lambda Whitespace?" json:"parameters,omitempty"`
	Expr       *Expr              `parser:"@@"                                                                               json:"expr,omitempty"`
}

type LambdaParameter struct {
	Label string         `parser:"@Ident"                            json:"label"`
	Type  *ParameterType `parser:"(Whitespace? ':' Whitespace? @@)?" json:"type"`
}
