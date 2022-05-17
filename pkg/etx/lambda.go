package etx

type Lambda struct {
	Parameters []*LambdaParameter `parser:"'(' [ @@ (',' @@)* ] ')' OpLambda" json:"parameters,omitempty"`
	Expr       *Expr              `parser:"@@"                                json:"expr,omitempty"`
}

type LambdaParameter struct {
	Label string         `parser:"@Ident"     json:"label"`
	Type  *ParameterType `parser:"[ ':' @@ ]" json:"type"`
}
