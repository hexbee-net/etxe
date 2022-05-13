package etx

type Func struct {
	Comments []string `parser:"@Comment*" json:"comments,omitempty"`

	Label      string           `parser:"Func Whitespace @Ident Whitespace?" json:"label"`
	Parameters []*FuncParameter `parser:"'(' (@@ ( Whitespace? ',' Whitespace? @@ )*)? ')'" json:"parameters,omitempty"`
	Return     *ParameterType   `parser:"Whitespace? @@?" json:"return,omitempty"`
	Body       []*FuncExpr      `parser:"(Whitespace|NewLine)? BodyStart @@*  BodyEnd" json:"body,omitempty"`
}

type FuncParameter struct {
	Label string         `parser:"@Ident" json:"label"`
	Type  *ParameterType `parser:"Whitespace? ':' Whitespace? @@" json:"type"`
}

type FuncExpr struct {
	Todo string `parser:"(Whitespace|NewLine)? @Ident (Whitespace|NewLine)?" json:"todo"`
}

type FuncDecl struct {
	Comments []string `parser:"@Comment*" json:"comments,omitempty"`

	DeclType string `parser:"@(Const | Val) Whitespace" json:"decl_type"`
	Label    string `parser:"@Ident" json:"label"`
	Type     string `parser:"(Whitespace? ':' Whitespace? @Ident)?" json:"type"`
	Value    *Value `parser:"(Whitespace? '=' Whitespace? @@)?" json:"value"`
}

type FuncReturn struct {
	Comments []string `parser:"@Comment*" json:"comments,omitempty"`
	Expr     *Expr    `parser:"Return @@" json:"expr,omitempty"`
}
