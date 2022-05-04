package etx

type Expr struct {
	Left  *Terminal `parser:"@@"`
	Op    string    `parser:"( @Oper"`
	Right *Terminal `parser:"  @@)?"`
}

func (e *Expr) String() string {
	return "TODO: Implement me"
}

type Terminal struct {
	String *String `parser:"  @@"`
	Ident  string  `parser:"| @Ident"`
}
