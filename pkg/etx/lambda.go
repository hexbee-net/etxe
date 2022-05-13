package etx

type Lambda struct {
	Parameters []*LambdaParameter `parser:"'(' (@@ ( Whitespace? ',' Whitespace? @@ )*)? ')' Whitespace? Lambda Whitespace?" json:"parameters,omitempty"`
	Expr       *Expr              `parser:"@@"`
}

type LambdaParameter struct {
	Label string         `parser:"@Ident" json:"label"`
	Type  *ParameterType `parser:"(Whitespace? ':' Whitespace? @@)?" json:"type"`
}
