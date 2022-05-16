package etx

type Var struct {
	Name    string `parser:"'var' @Ident" json:"name,omitempty"`
	Type    string `parser:"':' @Ident"   json:"type,omitempty"`
	Default string `parser:"'=' @String"  json:"default,omitempty"`
}
