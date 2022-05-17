package etx

type Lambda struct {
	Parameters []LambdaParameter `parser:"'(' [ @@ (',' @@)* ] ')' OpLambda" json:"parameters"`
	Expr       Expr              `parser:"@@"                                json:"expr"`
}

type LambdaParameter struct {
	Label string         `parser:"@Ident"     json:"label"`
	Type  *ParameterType `parser:"[ ':' @@ ]" json:"type,omitempty"`
}
