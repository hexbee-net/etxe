package etx

type Func struct {
	Comments []string `parser:"@Comment*" json:"comments,omitempty"`

	Label      string           `parser:"Func Whitespace @Ident Whitespace?"                                                          json:"label"`
	Parameters []*FuncParameter `parser:"'(' (@@ ( Whitespace? ',' Whitespace? @@ )*)? ')'"                                           json:"parameters,omitempty"`
	Return     *ParameterType   `parser:"Whitespace? @@?"                                                                             json:"return,omitempty"`
	Body       []*FuncStatement `parser:"(Whitespace|NewLine)? BodyStart ( (Whitespace|NewLine)? @@ (Whitespace|NewLine)? )* BodyEnd" json:"body,omitempty"`
}

type FuncParameter struct {
	Label string         `parser:"@Ident"                         json:"label"`
	Type  *ParameterType `parser:"Whitespace? ':' Whitespace? @@" json:"type"`
}

type FuncStatement struct {
	Comments []string  `parser:"@Comment*"                        json:"comments,omitempty"`
	Decl     *FuncDecl `parser:"(   Whitespace? @@ Whitespace?"   json:"decl,omitempty"`
	Expr     *Expr     `parser:"  | Whitespace? @@ Whitespace? )" json:"expr,omitempty"`
}

type FuncDecl struct {
	DeclType string `parser:"@(Const | Val) Whitespace"             json:"decl_type"`
	Label    string `parser:"@Ident"                                json:"label"`
	Type     string `parser:"(Whitespace? ':' Whitespace? @Ident)?" json:"type"`
	Value    *Expr  `parser:"(Whitespace? '=' Whitespace? @@)?"     json:"value"`
}
