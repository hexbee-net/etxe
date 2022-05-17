package etx

type Func struct {
	Comments   []string         `parser:"@Comment*"                                            json:"comments,omitempty"`
	Label      string           `parser:"Func @Ident "                                         json:"label"`
	Parameters []*FuncParameter `parser:"'(' [ @@ (',' @@)* ] ')'"                            json:"parameters,omitempty"`
	Return     *ParameterType   `parser:"@@?"                                                  json:"return,omitempty"`
	Body       []*FuncStatement `parser:"NewLine? BodyStart ( NewLine? @@ NewLine? )* BodyEnd" json:"body,omitempty"`
}

type FuncParameter struct {
	Label string         `parser:"@Ident" json:"label"`
	Type  *ParameterType `parser:"':' @@" json:"type"`
}

type FuncStatement struct {
	Comments []string  `parser:"@Comment*" json:"comments,omitempty"`
	Decl     *FuncDecl `parser:"(   @@  "  json:"decl,omitempty"`
	Expr     *Expr     `parser:"  | @@ )"  json:"expr,omitempty"`
}

type FuncDecl struct {
	DeclType string `parser:"@(Const | Val)" json:"decl_type"`
	Label    string `parser:"@Ident"         json:"label"`
	Type     string `parser:"[ ':' @Ident ]" json:"type,omitempty"`
	Value    *Expr  `parser:"[ '=' @@     ]" json:"value,omitempty"`
}
