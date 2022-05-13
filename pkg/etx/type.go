package etx

type Type struct {
	Label string `parser:"'type' Whitespace @Ident"`
}
