package etx

type Var struct {
	Name    string `parser:"'var' @Ident"`
	Type    string `parser:"':' @Ident"`
	Default string `parser:"'=' @String"`
}
