package etx

type Type struct {
	Comments []string `parser:"@Comment*"               json:"comments,omitempty"`
	Label    string   `parser:"'type' @Ident"           json:"label"`
	Enum     *Enum    `parser:"(   Enum   '{' @@ NewLine? '}' "  json:"enum,omitempty"`
	Object   *Object  `parser:"  | Object '{' @@ NewLine? '}' )" json:"object,omitempty"`
}
type Enum struct {
	Items []EnumItem `parser:"(NewLine @@)*"  json:"items,omitempty"`
}

type EnumItem struct {
	Comments []string `parser:"@Comment*"  json:"comments,omitempty"`
	Label    string   `parser:"@Ident ':'" json:"label"`
	Value    Expr     `parser:"@@"         json:"value"`
}

type Object struct {
	Items []ObjectItem `parser:"(NewLine @@)*" json:"items,omitempty"`
}

type ObjectItem struct {
	Comments []string      `parser:"@Comment*"  json:"comments,omitempty"`
	Label    string        `parser:"@Ident ':'" json:"label"`
	Type     ParameterType `parser:"@@"         json:"type"`
}
